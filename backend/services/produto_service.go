package services

import (
	"bufio"
	"fmt"
	"strings"
	"time"

	database "gostock/backend/Database"
	"gostock/backend/models"
	"gostock/backend/utils"
)

func CadastrarProduto(reader *bufio.Reader, usuarioID int) {
	nome := utils.LerNomeValido(reader)
	quantidade := utils.LerQuantidadeValida(reader, "Quantidade: ")
	preco := utils.LerPrecoValido(reader, "Preço: ")

	query := `
		INSERT INTO produtos (nome, preco, quantidade)
		VALUES (?, ?, ?)
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
		FROM produtos
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
		fmt.Println("ID:", produto.ID)
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

func AtualizarProduto(reader *bufio.Reader, usuarioID int) {
	id, err := utils.LerInteiro(
		reader,
		"Digite o ID do produto: ",
	)

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
	`, id).Scan(
		&produtoID,
		&nomeAtual,
		&precoAtual,
	)

	if err != nil {
		fmt.Println("Produto não encontrado.")
		return
	}

	fmt.Println("Produto encontrado.")
	fmt.Println("Nome atual:", nomeAtual)
	fmt.Println("Preço atual:", precoAtual)

	novoNome := utils.LerNomeValido(reader)

	novoPreco := utils.LerPrecoValido(
		reader,
		"Digite o novo preço: ",
	)

	_, err = database.DB.Exec(`
		UPDATE produtos
		SET nome = ?, preco = ?
		WHERE id = ?
	`,
		novoNome,
		novoPreco,
		produtoID,
	)

	if err != nil {
		fmt.Println("Erro ao atualizar produto:", err)
		return
	}

	_, err = database.DB.Exec(`
		INSERT INTO movimentacoes
		(produto_id, usuario_id, tipo, quantidade)
		VALUES (?, ?, ?, ?)
	`,
		produtoID,
		usuarioID,
		"ATUALIZACAO",
		0,
	)

	if err != nil {
		fmt.Println("Erro ao registrar atualização:", err)
		return
	}

	fmt.Println("Produto atualizado com sucesso!")
}

func AdicionarEstoque(reader *bufio.Reader, usuarioID int) {
	id, err := utils.LerInteiro(
		reader,
		"Digite o ID do produto: ",
	)

	if err != nil {
		fmt.Println("ID inválido.")
		return
	}

	var produtoID int
	var nome string
	var quantidadeAtual int

	err = database.DB.QueryRow(`
		SELECT id, nome, quantidade
		FROM produtos
		WHERE id = ?
	`, id).Scan(
		&produtoID,
		&nome,
		&quantidadeAtual,
	)

	if err != nil {
		fmt.Println("Produto não encontrado.")
		return
	}

	fmt.Println("Produto:", nome)
	fmt.Println("Estoque atual:", quantidadeAtual)

	quantidadeAdicionar := utils.LerQuantidadeValida(
		reader,
		"Quantidade que chegou: ",
	)

	if quantidadeAdicionar <= 0 {
		fmt.Println("A quantidade deve ser maior que zero.")
		return
	}

	_, err = database.DB.Exec(`
		UPDATE produtos
		SET quantidade = quantidade + ?
		WHERE id = ?
	`,
		quantidadeAdicionar,
		produtoID,
	)

	if err != nil {
		fmt.Println("Erro ao atualizar estoque:", err)
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
		quantidadeAdicionar,
	)

	if err != nil {
		fmt.Println("Erro ao registrar entrada:", err)
		return
	}

	fmt.Println("Estoque atualizado com sucesso!")
}

func RegistrarVenda(reader *bufio.Reader, usuarioID int) {
	id, err := utils.LerInteiro(
		reader,
		"Digite o ID do produto: ",
	)

	if err != nil {
		fmt.Println("ID inválido.")
		return
	}

	var produtoID int
	var nome string
	var quantidadeAtual int

	err = database.DB.QueryRow(`
		SELECT id, nome, quantidade
		FROM produtos
		WHERE id = ?
	`, id).Scan(
		&produtoID,
		&nome,
		&quantidadeAtual,
	)

	if err != nil {
		fmt.Println("Produto não encontrado.")
		return
	}

	fmt.Println("Produto:", nome)
	fmt.Println("Estoque atual:", quantidadeAtual)

	quantidadeVendida := utils.LerQuantidadeValida(
		reader,
		"Quantidade vendida: ",
	)

	if quantidadeVendida <= 0 {
		fmt.Println("A quantidade deve ser maior que zero.")
		return
	}

	if quantidadeVendida > quantidadeAtual {
		fmt.Println("Estoque insuficiente.")
		fmt.Println("Estoque disponível:", quantidadeAtual)
		return
	}

	_, err = database.DB.Exec(`
		UPDATE produtos
		SET quantidade = quantidade - ?
		WHERE id = ?
	`,
		quantidadeVendida,
		produtoID,
	)

	if err != nil {
		fmt.Println("Erro ao atualizar estoque:", err)
		return
	}

	_, err = database.DB.Exec(`
		INSERT INTO movimentacoes
		(produto_id, usuario_id, tipo, quantidade)
		VALUES (?, ?, ?, ?)
	`,
		produtoID,
		usuarioID,
		"SAIDA",
		quantidadeVendida,
	)

	if err != nil {
		fmt.Println("Erro ao registrar saída:", err)
		return
	}

	fmt.Println("Venda registrada com sucesso!")
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

		if tipo == "ATUALIZACAO" {
			fmt.Println("Quantidade: -")
		} else {
			fmt.Println("Quantidade:", quantidade)
		}

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

		if err := rows.Scan(
			&produto.ID,
			&produto.Nome,
			&produto.Preco,
			&produto.Quantidade,
		); err != nil {
			return nil, err
		}

		if produto.Quantidade == 0 {
			produto.StatusEstoque = "zerado"
		} else if produto.Quantidade <= 5 {
			produto.StatusEstoque = "baixo"
		} else {
			produto.StatusEstoque = "normal"
		}

		produtos = append(produtos, produto)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return produtos, nil
}

func AdicionarEstoqueWeb(produtoID int, quantidade int, usuarioID int) error {
	if quantidade <= 0 {
		return fmt.Errorf("a quantidade deve ser maior que zero")
	}

	_, err := database.DB.Exec(`
		UPDATE produtos
		SET quantidade = quantidade + ?
		WHERE id = ?
	`, quantidade, produtoID)

	if err != nil {
		return err
	}

	_, err = database.DB.Exec(`
		INSERT INTO movimentacoes
		(produto_id, usuario_id, tipo, quantidade)
		VALUES (?, ?, ?, ?)
	`, produtoID, usuarioID, "ENTRADA", quantidade)

	if err != nil {
		return err
	}

	return nil
}
