package services

import (
	"bufio"
	"database/sql"
	"fmt"

	database "gostock/backend/Database"
	"gostock/backend/utils"
)

func AdicionarEstoque(reader *bufio.Reader, usuarioID int) {
	id, err := utils.LerInteiro(reader, "Digite o ID do produto: ")
	if err != nil {
		fmt.Println("ID inválido.")
		return
	}

	produtoID, nome, quantidadeAtual, err := buscarDadosEstoque(id)
	if err != nil {
		fmt.Println("Produto não encontrado.")
		return
	}

	fmt.Println("Produto:", nome)
	fmt.Println("Estoque atual:", quantidadeAtual)

	quantidade := utils.LerQuantidadeValida(reader, "Quantidade que chegou: ")
	if quantidade <= 0 {
		fmt.Println("A quantidade deve ser maior que zero.")
		return
	}

	if err := AdicionarEstoqueWeb(produtoID, quantidade, usuarioID); err != nil {
		fmt.Println("Erro ao atualizar estoque:", err)
		return
	}

	fmt.Println("Estoque atualizado com sucesso!")
}

func AdicionarEstoqueWeb(produtoID int, quantidade int, usuarioID int) error {
	if quantidade <= 0 {
		return fmt.Errorf("a quantidade deve ser maior que zero")
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.Exec(`
		UPDATE produtos
		SET quantidade = quantidade + ?
		WHERE id = ? AND ativo = 1
	`, quantidade, produtoID)
	if err != nil {
		return err
	}

	linhas, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if linhas == 0 {
		return fmt.Errorf("produto não encontrado")
	}

	if err := registrarMovimentacaoTx(tx, produtoID, usuarioID, "ENTRADA", quantidade); err != nil {
		return err
	}

	return tx.Commit()
}

func RegistrarSaida(reader *bufio.Reader, usuarioID int) {
	id, err := utils.LerInteiro(reader, "Digite o ID do produto: ")
	if err != nil {
		fmt.Println("ID inválido.")
		return
	}

	produtoID, nome, quantidadeAtual, err := buscarDadosEstoque(id)
	if err != nil {
		fmt.Println("Produto não encontrado.")
		return
	}

	fmt.Println("Produto:", nome)
	fmt.Println("Estoque atual:", quantidadeAtual)

	quantidade := utils.LerQuantidadeValida(reader, "Quantidade que saiu: ")
	if quantidade <= 0 {
		fmt.Println("A quantidade deve ser maior que zero.")
		return
	}

	if err := RegistrarSaidaWeb(produtoID, quantidade, usuarioID); err != nil {
		if err.Error() == "estoque insuficiente" {
			fmt.Println("Estoque insuficiente.")
			fmt.Println("Estoque disponível:", quantidadeAtual)
			return
		}

		fmt.Println("Erro ao registrar saída:", err)
		return
	}

	fmt.Println("Saída registrada com sucesso!")
}

func RegistrarSaidaWeb(produtoID int, quantidade int, usuarioID int) error {
	if quantidade <= 0 {
		return fmt.Errorf("a quantidade deve ser maior que zero")
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var estoque int

	err = tx.QueryRow(`
		SELECT quantidade
		FROM produtos
		WHERE id = ? AND ativo = 1
	`, produtoID).Scan(&estoque)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("produto não encontrado")
		}

		return err
	}

	if quantidade > estoque {
		return fmt.Errorf("estoque insuficiente")
	}

	_, err = tx.Exec(`
		UPDATE produtos
		SET quantidade = quantidade - ?
		WHERE id = ?
	`, quantidade, produtoID)

	if err != nil {
		return err
	}

	if err := registrarMovimentacaoTx(
		tx,
		produtoID,
		usuarioID,
		"SAIDA",
		quantidade,
	); err != nil {
		return err
	}

	return tx.Commit()
}
func buscarDadosEstoque(produtoID int) (int, string, int, error) {
	var nome string
	var quantidade int

	err := database.DB.QueryRow(`
		SELECT id, nome, quantidade
		FROM produtos
		WHERE id = ?
	`, produtoID).Scan(&produtoID, &nome, &quantidade)

	return produtoID, nome, quantidade, err
}

func registrarMovimentacaoTx(tx *sql.Tx, produtoID int, usuarioID int, tipo string, quantidade int) error {
	_, err := tx.Exec(`
		INSERT INTO movimentacoes
		(produto_id, usuario_id, tipo, quantidade)
		VALUES (?, ?, ?, ?)
	`, produtoID, usuarioID, tipo, quantidade)
	return err
}
