package services

import (
	"bufio"
	"fmt"

	database "gostock/backend/Database"
	"gostock/backend/utils"
)

func CadastrarProduto(reader *bufio.Reader, usuarioID int) {
	nome := utils.LerNomeValido(reader)
	quantidade := utils.LerQuantidadeValida(reader, "Quantidade: ")
	preco := utils.LerPrecoValido(reader, "Preço: ")

	tx, err := database.DB.Begin()
	if err != nil {
		fmt.Println("Erro ao iniciar cadastro:", err)
		return
	}
	defer tx.Rollback()

	resultado, err := tx.Exec(`
		INSERT INTO produtos (nome, preco, quantidade)
		VALUES (?, ?, ?)
	`, nome, preco, quantidade)

	if err != nil {
		fmt.Println("Erro ao cadastrar produto:", err)
		return
	}

	produtoID, err := resultado.LastInsertId()
	if err != nil {
		fmt.Println("Erro ao obter ID do produto:", err)
		return
	}

	if err := registrarMovimentacaoTx(tx, int(produtoID), usuarioID, "ENTRADA", quantidade); err != nil {
		fmt.Println("Erro ao registrar movimentação:", err)
		return
	}

	if err := tx.Commit(); err != nil {
		fmt.Println("Erro ao confirmar cadastro:", err)
		return
	}

	fmt.Println("Produto cadastrado com sucesso!")
}

func RemoverProduto(reader *bufio.Reader) {
	id, err := utils.LerInteiro(reader, "Digite o ID do produto: ")
	if err != nil {
		fmt.Println("ID inválido.")
		return
	}

	resultado, err := database.DB.Exec(`
		DELETE FROM produtos
		WHERE id = ?
	`, id)

	if err != nil {
		fmt.Println("Erro ao remover produto:", err)
		return
	}

	linhas, err := resultado.RowsAffected()
	if err != nil {
		fmt.Println("Erro ao verificar remoção:", err)
		return
	}

	if linhas > 0 {
		fmt.Println("Produto removido com sucesso.")
	} else {
		fmt.Println("Produto não encontrado.")
	}

}
