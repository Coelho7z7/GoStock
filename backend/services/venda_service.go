package services

import (
	"database/sql"
	"fmt"

	database "gostock/backend/Database"
)

type ItemVenda struct {
	ProdutoID     int
	Quantidade    int
	PrecoUnitario float64
}

// formasPagamentoValidas lista as formas de pagamento aceitas pelo PDV.
var formasPagamentoValidas = map[string]bool{
	"dinheiro": true,
	"cartao":   true,
	"pix":      true,
	"debito":   true,
}

// normalizarFormaPagamento garante que sempre seja gravado um valor
// conhecido, usando "dinheiro" como padrão quando o valor vier vazio
// ou fora da lista esperada.
func normalizarFormaPagamento(formaPagamento string) string {
	if formasPagamentoValidas[formaPagamento] {
		return formaPagamento
	}
	return "dinheiro"
}

func calcularValorTotalVenda(quantidade int, precoUnitario float64) float64 {
	return float64(quantidade) * precoUnitario
}

func calcularValorTotalItensVenda(itens []ItemVenda) float64 {
	total := 0.0
	for _, item := range itens {
		total += calcularValorTotalVenda(item.Quantidade, item.PrecoUnitario)
	}
	return total
}

// RegistrarVendaCompletaWeb registra uma venda com um ou mais itens,
// debitando o estoque, gravando a venda e a movimentação de saída de
// cada produto dentro de uma única transação.
func RegistrarVendaCompletaWeb(itens []ItemVenda, usuarioID int, formaPagamento string) error {
	if len(itens) == 0 {
		return fmt.Errorf("a venda deve conter pelo menos um item")
	}

	formaPagamento = normalizarFormaPagamento(formaPagamento)

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	valorTotal := 0.0

	for _, item := range itens {
		if item.Quantidade <= 0 {
			return fmt.Errorf("a quantidade deve ser maior que zero")
		}

		var produtoIDAtual int
		var nome string
		var preco float64
		var estoque int

		err = tx.QueryRow(`
			SELECT id, nome, preco, quantidade
			FROM produtos
			WHERE id = ? AND ativo = 1
		`, item.ProdutoID).Scan(&produtoIDAtual, &nome, &preco, &estoque)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("produto não encontrado: %d", item.ProdutoID)
			}
			return err
		}

		if item.Quantidade > estoque {
			return fmt.Errorf("estoque insuficiente para %s", nome)
		}

		valorTotal += calcularValorTotalVenda(item.Quantidade, preco)

		_, err = tx.Exec(`
			UPDATE produtos
			SET quantidade = quantidade - ?
			WHERE id = ?
		`, item.Quantidade, item.ProdutoID)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			INSERT INTO vendas (produto_id, usuario_id, quantidade, valor_unitario, valor_total, forma_pagamento)
			VALUES (?, ?, ?, ?, ?, ?)
		`, item.ProdutoID, usuarioID, item.Quantidade, preco, calcularValorTotalVenda(item.Quantidade, preco), formaPagamento)
		if err != nil {
			return err
		}

		if err := registrarMovimentacaoTx(tx, item.ProdutoID, usuarioID, "SAIDA", item.Quantidade); err != nil {
			return err
		}
	}

	if valorTotal <= 0 {
		return fmt.Errorf("valor da venda deve ser maior que zero")
	}

	return tx.Commit()
}
