package services

import (
	"bufio"
	"fmt"
	"strings"

	database "gostock/backend/Database"
	"gostock/backend/models"
	"gostock/backend/utils"
)

func ListarProdutos() {
	rows, err := database.DB.Query(`
		SELECT id, nome, preco, quantidade
		FROM produtos
	`)
	if err != nil {
		fmt.Println("Erro ao buscar produtos:", err)
		return
	}
	defer rows.Close()

	encontrou := false

	for rows.Next() {
		var produto models.Produto

		if err := rows.Scan(&produto.ID, &produto.Nome, &produto.Preco, &produto.Quantidade); err != nil {
			fmt.Println("Erro ao ler produto:", err)
			return
		}

		fmt.Println("ID:", produto.ID)
		fmt.Println("Nome:", produto.Nome)
		fmt.Println("Quantidade:", produto.Quantidade)
		fmt.Println("Preço:", produto.Preco)
		fmt.Println("----------------------")
		encontrou = true
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Erro ao percorrer produtos:", err)
		return
	}

	if !encontrou {
		fmt.Println("Nenhum produto cadastrado.")
	}
}

func BuscarProduto(reader *bufio.Reader) {
	buscar := strings.TrimSpace(utils.LerTexto(reader, "Digite o nome do produto: "))

	rows, err := database.DB.Query(`
		SELECT id, nome, preco, quantidade
		FROM produtos
		WHERE nome LIKE ?
	`, "%"+buscar+"%")
	if err != nil {
		fmt.Println("Erro ao buscar produto:", err)
		return
	}
	defer rows.Close()

	encontrado := false

	for rows.Next() {
		var produto models.Produto

		if err := rows.Scan(&produto.ID, &produto.Nome, &produto.Preco, &produto.Quantidade); err != nil {
			fmt.Println("Erro ao ler produto:", err)
			return
		}

		fmt.Println("Produto encontrado!")
		fmt.Println("ID:", produto.ID)
		fmt.Println("Nome:", produto.Nome)
		fmt.Println("Quantidade:", produto.Quantidade)
		fmt.Println("Preço:", produto.Preco)
		encontrado = true
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Erro ao percorrer produtos:", err)
		return
	}

	if !encontrado {
		fmt.Println("Produto não encontrado.")
	}
}

func BuscarTodosProdutos() ([]models.Produto, error) {
	rows, err := database.DB.Query(`
		SELECT id, nome, preco, quantidade
		FROM produtos
		ORDER BY id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var produtos []models.Produto

	for rows.Next() {
		var produto models.Produto

		if err := rows.Scan(&produto.ID, &produto.Nome, &produto.Preco, &produto.Quantidade); err != nil {
			return nil, err
		}

		produto.StatusEstoque = statusDoEstoque(produto.Quantidade)
		produtos = append(produtos, produto)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return produtos, nil
}

func statusDoEstoque(quantidade int) string {
	switch {
	case quantidade == 0:
		return "zerado"
	case quantidade <= 5:
		return "baixo"
	default:
		return "normal"
	}
}

func AtualizarProduto(reader *bufio.Reader, usuarioID int) {
	id, err := utils.LerInteiro(reader, "Digite o ID do produto: ")
	if err != nil {
		fmt.Println("ID inválido.")
		return
	}

	var produtoID int
	var nomeAtual string
	var precoAtual float64

	err = database.DB.QueryRow(`
		SELECT id, nome, preco
		FROM produtos
		WHERE id = ?
	`, id).Scan(&produtoID, &nomeAtual, &precoAtual)

	if err != nil {
		fmt.Println("Produto não encontrado.")
		return
	}

	fmt.Println("Produto encontrado.")
	fmt.Println("Nome atual:", nomeAtual)
	fmt.Println("Preço atual:", precoAtual)

	novoNome := utils.LerNomeValido(reader)
	novoPreco := utils.LerPrecoValido(reader, "Digite o novo preço: ")

	_, err = database.DB.Exec(`
		UPDATE produtos
		SET nome = ?, preco = ?
		WHERE id = ?
	`, novoNome, novoPreco, produtoID)
	if err != nil {
		fmt.Println("Erro ao atualizar produto:", err)
		return
	}

	if err := registrarMovimentacao(produtoID, usuarioID, "ATUALIZACAO", 0); err != nil {
		fmt.Println("Erro ao registrar atualização:", err)
		return
	}

	fmt.Println("Produto atualizado com sucesso!")
}
