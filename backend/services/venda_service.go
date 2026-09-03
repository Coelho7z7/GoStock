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

type ItemVendaRegistrada struct {
	Nome       string  `json:"nome"`
	Quantidade int     `json:"quantidade"`
	Preco      float64 `json:"preco"`
	Subtotal   float64 `json:"subtotal"`
}

type VendaRegistrada struct {
	Numero         int                   `json:"numero"`
	FormaPagamento string                `json:"formaPagamento"`
	Total          float64               `json:"total"`
	Itens          []ItemVendaRegistrada `json:"itens"`
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
	_, err := RegistrarVendaCompletaWebDetalhes(itens, usuarioID, formaPagamento)
	return err
}

func RegistrarVendaCompletaWebDetalhes(itens []ItemVenda, usuarioID int, formaPagamento string) (VendaRegistrada, error) {
	var venda VendaRegistrada

	if len(itens) == 0 {
		return venda, fmt.Errorf("a venda deve conter pelo menos um item")
	}

	formaPagamento = normalizarFormaPagamento(formaPagamento)
	venda.FormaPagamento = formaPagamento

	tx, err := database.DB.Begin()
	if err != nil {
		return venda, err
	}
	defer tx.Rollback()

	valorTotal := 0.0

	for _, item := range itens {
		if item.Quantidade <= 0 {
			return venda, fmt.Errorf("a quantidade deve ser maior que zero")
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
				return venda, fmt.Errorf("produto não encontrado: %d", item.ProdutoID)
			}
			return venda, err
		}

		if item.Quantidade > estoque {
			return venda, fmt.Errorf("estoque insuficiente para %s", nome)
		}

		valorTotal += calcularValorTotalVenda(item.Quantidade, preco)

		_, err = tx.Exec(`
			UPDATE produtos
			SET quantidade = quantidade - ?
			WHERE id = ?
		`, item.Quantidade, item.ProdutoID)
		if err != nil {
			return venda, err
		}

		resultado, err := tx.Exec(`
			INSERT INTO vendas (produto_id, usuario_id, quantidade, valor_unitario, valor_total, forma_pagamento)
			VALUES (?, ?, ?, ?, ?, ?)
		`, item.ProdutoID, usuarioID, item.Quantidade, preco, calcularValorTotalVenda(item.Quantidade, preco), formaPagamento)
		if err != nil {
			return venda, err
		}
		if venda.Numero == 0 {
			id, idErr := resultado.LastInsertId()
			if idErr != nil {
				return venda, idErr
			}
			venda.Numero = int(id)
		}
		venda.Itens = append(venda.Itens, ItemVendaRegistrada{
			Nome: nome, Quantidade: item.Quantidade, Preco: preco,
			Subtotal: calcularValorTotalVenda(item.Quantidade, preco),
		})

		if err := registrarMovimentacaoTx(tx, item.ProdutoID, usuarioID, "SAIDA", item.Quantidade); err != nil {
			return venda, err
		}
	}

	if valorTotal <= 0 {
		return venda, fmt.Errorf("valor da venda deve ser maior que zero")
	}

	venda.Total = valorTotal
	return venda, tx.Commit()
}
