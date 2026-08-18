package services

import (
	"bufio"
	"fmt"
	"strings"

	database "gostock/Database"
	"gostock/models"
	"gostock/utils"
)

func CadastrarProduto(reader *bufio.Reader) {
	nome := utils.LerNomeValido(reader)
	quantidade := utils.LerQuantidadeValida(reader, "Quantidade:")
	preco := utils.LerPrecoValido(reader, "Preço:")

	query := `
		INSERT INTO produtos (nome, preco, quantidade)
		VALUES (?, ?, ?);
	`

	_, err := database.DB.Exec(query, nome, preco, quantidade)

	if err != nil {
		fmt.Println("Erro ao cadastrar produto:", err)
		return
	}

	fmt.Println("Produto cadastrado com sucesso!")
}

func ListarProdutos() ([]models.Produto, error) {
	query := `
		SELECT id, nome, preco, quantidade
		FROM produtos;
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var produtos []models.Produto

	for rows.Next() {
		var produto models.Produto

		err := rows.Scan(
			&produto.ID,
			&produto.Nome,
			&produto.Preco,
			&produto.Quantidade,
		)

		if err != nil {
			return nil, err
		}

		produtos = append(produtos, produto)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return produtos, nil

}

func BuscarProduto(reader *bufio.Reader) {
	buscar := strings.TrimSpace(utils.LerTexto(reader, "Digite o nome do produto: "))
	encontrado := false

	query := `
		SELECT id, nome, preco, quantidade
		FROM produtos
		WHERE nome LIKE ?
	`

	rows, err := database.DB.Query(query, "%"+buscar+"%")
	if err != nil {
		fmt.Println("Erro ao buscar produto:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var produto models.Produto

		err := rows.Scan(
			&produto.ID,
			&produto.Nome,
			&produto.Preco,
			&produto.Quantidade,
		)

		if err != nil {
			fmt.Println("Erro ao ler produto:", err)
			return
		}

		fmt.Println("Produto encontrado!")
		fmt.Println("Nome:", produto.Nome)
		fmt.Println("Quantidade:", produto.Quantidade)
		fmt.Println("Preço:", produto.Preco)

		encontrado = true
	}

	if !encontrado {
		fmt.Println("Produto não encontrado.")
	}
}

func RemoverProduto(reader *bufio.Reader) {
	id, err := utils.LerInteiro(reader, "Digite o ID do produto: ")

	if err != nil {
		fmt.Println("ID inválido.")
		return
	}
	removido := false
	query := `
	DELETE FROM produtos
	WHERE id = ?
`

	resultado, err := database.DB.Exec(query, id)

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
		removido = true
	}

	if !removido {
		fmt.Println("Produto não encontrado.")
	}
}

func AtualizarProduto(reader *bufio.Reader) {
	id, err := utils.LerInteiro(reader, "Digite o ID do produto: ")

	if err != nil {
		fmt.Println("ID inválido.")
		return
	}

	query := `
		SELECT id
		FROM produtos
		WHERE id = ?
	`

	var produtoID int

	err = database.DB.QueryRow(query, id).Scan(&produtoID)

	if err != nil {
		fmt.Println("Produto não encontrado.")
		return
	}

	fmt.Println("Produto encontrado.")

	novoNome := utils.LerNomeValido(reader)

	novaQuantidade := utils.LerQuantidadeValida(
		reader,
		"Digite sua nova quantidade: ",
	)

	novoPreco := utils.LerPrecoValido(
		reader,
		"Digite seu novo preço: ",
	)

	query = `
		UPDATE produtos
		SET nome = ?, quantidade = ?, preco = ?
		WHERE id = ?
	`

	_, err = database.DB.Exec(
		query,
		novoNome,
		novaQuantidade,
		novoPreco,
		id,
	)

	if err != nil {
		fmt.Println("Erro ao atualizar produto:", err)
		return
	}

	fmt.Println("Produto atualizado com sucesso!")
}
