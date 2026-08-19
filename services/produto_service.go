package services

import (
	"bufio"
	"fmt"
	"strings"
	"time"

	"gostock/database"
	"gostock/models"
	"gostock/utils"
)

func CadastrarProduto(reader *bufio.Reader, usuarioID int) {
	nome := utils.LerNomeValido(reader)
	quantidade := utils.LerQuantidadeValida(reader, "Quantidade:")
	preco := utils.LerPrecoValido(reader, "Preço:")

	query := `
		INSERT INTO produtos (nome, preco, quantidade)
		VALUES (?, ?, ?);
	`

	resultado, err := database.DB.Exec(query, nome, preco, quantidade)

	if err != nil {
		fmt.Println("Erro ao cadastrar produto:", err)
		return
	}

	produtoID, err := resultado.LastInsertId()

	if err != nil {
		fmt.Println("Erro ao obter ID do produto:", err)
		return
	}

	_, err = database.DB.Exec(`
		INSERT INTO movimentacoes
		(produto_id, usuario_id, tipo, quantidade)
		VALUES (?, ?, ?, ?)
	`,
		produtoID,
		usuarioID,
		"ENTRADA",
		quantidade,
	)

	if err != nil {
		fmt.Println("Erro ao registrar movimentação:", err)
		return
	}

	fmt.Println("Produto cadastrado com sucesso!")
}

func ListarProdutos() {
	query := `
		SELECT id, nome, preco, quantidade
		FROM produtos;
	`

	rows, err := database.DB.Query(query)

	if err != nil {
		fmt.Println("Erro ao buscar produtos:", err)
		return
	}

	defer rows.Close()

	encontrou := false

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

		fmt.Println("ID:", produto.ID)
		fmt.Println("Nome:", produto.Nome)
		fmt.Println("Quantidade:", produto.Quantidade)
		fmt.Println("Preço:", produto.Preco)
		fmt.Println("----------------------")

		encontrou = true
	}

	if !encontrou {
		fmt.Println("Nenhum produto cadastrado.")
	}
}

func BuscarProduto(reader *bufio.Reader) {
	buscar := strings.TrimSpace(
		utils.LerTexto(reader, "Digite o nome do produto: "),
	)

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
	id, err := utils.LerInteiro(
		reader,
		"Digite o ID do produto: ",
	)

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

func AtualizarProduto(reader *bufio.Reader, usuarioID int) {
	id, err := utils.LerInteiro(
		reader,
		"Digite o ID do produto: ",
	)

	if err != nil {
		fmt.Println("ID inválido.")
		return
	}

	query := `
		SELECT id, quantidade
		FROM produtos
		WHERE id = ?
	`

	var produtoID int
	var quantidadeAtual int

	err = database.DB.QueryRow(query, id).Scan(
		&produtoID,
		&quantidadeAtual,
	)

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

	diferenca := novaQuantidade - quantidadeAtual

	tipo := ""

	if diferenca > 0 {
		tipo = "ENTRADA"
	} else if diferenca < 0 {
		tipo = "SAIDA"
	}

	quantidadeMovimentada := diferenca

	if quantidadeMovimentada < 0 {
		quantidadeMovimentada = -quantidadeMovimentada
	}

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

	if diferenca != 0 {
		_, err = database.DB.Exec(`
			INSERT INTO movimentacoes
			(produto_id, usuario_id, tipo, quantidade)
			VALUES (?, ?, ?, ?)
		`,
			produtoID,
			usuarioID,
			tipo,
			quantidadeMovimentada,
		)

		if err != nil {
			fmt.Println("Erro ao registrar movimentação:", err)
			return
		}
	}

	fmt.Println("Produto atualizado com sucesso!")
}

func ListarMovimentacoes() {
	query := `
		SELECT
			m.data,
			p.nome,
			u.nome,
			m.tipo,
			m.quantidade
		FROM movimentacoes m
		JOIN produtos p ON p.id = m.produto_id
		JOIN usuarios u ON u.id = m.usuario_id
		ORDER BY m.data DESC
	`

	rows, err := database.DB.Query(query)

	if err != nil {
		fmt.Println("Erro ao buscar movimentações:", err)
		return
	}

	defer rows.Close()

	encontrou := false

	for rows.Next() {
		var data string
		var produto string
		var usuario string
		var tipo string
		var quantidade int

		err := rows.Scan(
			&data,
			&produto,
			&usuario,
			&tipo,
			&quantidade,
		)

		if err != nil {
			fmt.Println("Erro ao ler movimentação:", err)
			return
		}

		dataFormatada, err := time.Parse(
			time.RFC3339,
			data,
		)

		if err != nil {
			fmt.Println("Erro ao formatar data:", err)
			return
		}

		fmt.Println("========== MOVIMENTAÇÃO ==========")
		fmt.Println("Produto:", produto)
		fmt.Println("Usuário:", usuario)
		fmt.Println("Tipo:", tipo)
		fmt.Println("Quantidade:", quantidade)
		fmt.Println(
			"Data:",
			dataFormatada.Local().Format("02/01/2006 15:04:05"),
		)
		fmt.Println("==================================")

		encontrou = true
	}

	if !encontrou {
		fmt.Println("Nenhuma movimentação registrada.")
	}
}