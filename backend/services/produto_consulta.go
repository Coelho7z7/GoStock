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
		WHERE ativo = 1
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

		produto.StatusEstoque = statusDoEstoque(produto.Quantidade)
		produtos = append(produtos, produto)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return produtos, nil
}
func AtualizarProdutoWeb(produtoID int, nome string, preco float64, usuarioID int) error {
	nome = strings.TrimSpace(nome)
	if !utils.ValidarNome(nome) {
		return fmt.Errorf("o nome do produto é obrigatório")
	}
	if !utils.ValidarPreco(preco) {
		return fmt.Errorf("o preço não pode ser negativo")
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	resultado, err := tx.Exec(`
		UPDATE produtos
		SET nome = ?, preco = ?
		WHERE id = ? AND ativo = 1
	`, nome, preco, produtoID)
	if err != nil {
		return err
	}

	linhas, err := resultado.RowsAffected()
	if err != nil {
		return err
	}
	if linhas == 0 {
		return fmt.Errorf("produto não encontrado")
	}

	if err := registrarMovimentacaoTx(tx, produtoID, usuarioID, "ATUALIZACAO", 0); err != nil {
		return err
	}

	return tx.Commit()
}

func RemoverProdutoWeb(produtoID int) error {
	resultado, err := database.DB.Exec(`
		UPDATE produtos
		SET ativo = 0
		WHERE id = ? AND ativo = 1
	`, produtoID)
	if err != nil {
		return err
	}

	linhas, err := resultado.RowsAffected()
	if err != nil {
		return err
	}
	if linhas == 0 {
		return fmt.Errorf("produto não encontrado")
	}

	return nil
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
		WHERE id = ? AND ativo = 1
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

// ProdutosPaginados busca uma página de produtos ativos, opcionalmente
// filtrando por nome (busca parcial). pagina começa em 1.
func ProdutosPaginados(busca string, pagina int, porPagina int) ([]models.Produto, int, error) {
	return ProdutosPaginadosOrdenados(busca, pagina, porPagina, "recentes")
}

func ProdutosPaginadosOrdenados(busca string, pagina int, porPagina int, ordem string) ([]models.Produto, int, error) {
	if pagina < 1 {
		pagina = 1
	}
	if porPagina < 1 {
		porPagina = 10
	}

	busca = strings.TrimSpace(busca)
	filtroNome := "%" + busca + "%"

	var total int
	if err := database.DB.QueryRow(`
		SELECT COUNT(*)
		FROM produtos
		WHERE ativo = 1 AND nome LIKE ?
	`, filtroNome).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (pagina - 1) * porPagina

	ordenacao := "id DESC"
	switch ordem {
	case "nome":
		ordenacao = "nome COLLATE NOCASE ASC"
	case "preco":
		ordenacao = "preco ASC"
	case "estoque":
		ordenacao = "quantidade ASC"
	}
	rows, err := database.DB.Query(`
		SELECT id, nome, preco, quantidade
		FROM produtos
		WHERE ativo = 1 AND nome LIKE ?
		ORDER BY `+ordenacao+`
		LIMIT ? OFFSET ?
	`, filtroNome, porPagina, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var produtos []models.Produto

	for rows.Next() {
		var produto models.Produto

		if err := rows.Scan(&produto.ID, &produto.Nome, &produto.Preco, &produto.Quantidade); err != nil {
			return nil, 0, err
		}

		produto.StatusEstoque = statusDoEstoque(produto.Quantidade)
		produtos = append(produtos, produto)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return produtos, total, nil
}
