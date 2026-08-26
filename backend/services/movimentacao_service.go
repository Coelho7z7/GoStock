package services

import (
	"fmt"
	"time"

	database "gostock/backend/Database"
	"gostock/backend/models"
)

func BuscarMovimentacoesWeb() ([]models.Movimentacao, error) {
	rows, err := database.DB.Query(`
		SELECT
			m.id,
			m.produto_id,
			m.usuario_id,
			m.data,
			p.nome,
			u.nome,
			m.tipo,
			m.quantidade
		FROM movimentacoes m
		JOIN produtos p ON p.id = m.produto_id
		JOIN usuarios u ON u.id = m.usuario_id
		ORDER BY m.data DESC, m.id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movimentacoes []models.Movimentacao
	for rows.Next() {
		var movimentacao models.Movimentacao
		var data string

		if err := rows.Scan(
			&movimentacao.ID,
			&movimentacao.ProdutoID,
			&movimentacao.UsuarioID,
			&data,
			&movimentacao.Produto,
			&movimentacao.Usuario,
			&movimentacao.Tipo,
			&movimentacao.Quantidade,
		); err != nil {
			return nil, err
		}

		dataFormatada, err := parseDataMovimentacao(data)
		if err != nil {
			return nil, err
		}

		movimentacao.Data = dataFormatada
		movimentacao.DataFormatada = dataFormatada.Local().Format("02/01/2006")
		movimentacao.HoraFormatada = dataFormatada.Local().Format("15:04")
		movimentacao.TipoFormatado = map[string]string{
			"ENTRADA":     "Entrada",
			"SAIDA":       "Saída",
			"ATUALIZACAO": "Atualização",
		}[movimentacao.Tipo]
		if movimentacao.TipoFormatado == "" {
			movimentacao.TipoFormatado = movimentacao.Tipo
		}
		movimentacoes = append(movimentacoes, movimentacao)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movimentacoes, nil
}

func registrarMovimentacao(produtoID int, usuarioID int, tipo string, quantidade int) error {
	_, err := database.DB.Exec(`
		INSERT INTO movimentacoes
		(produto_id, usuario_id, tipo, quantidade)
		VALUES (?, ?, ?, ?)
	`, produtoID, usuarioID, tipo, quantidade)

	return err
}

func ListarMovimentacoes() {
	rows, err := database.DB.Query(`
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
	`)
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

		if err := rows.Scan(&data, &produto, &usuario, &tipo, &quantidade); err != nil {
			fmt.Println("Erro ao ler movimentação:", err)
			return
		}

		dataFormatada, err := parseDataMovimentacao(data)
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

		fmt.Println("Data:", dataFormatada.Local().Format("02/01/2006 15:04:05"))
		fmt.Println("==================================")

		encontrou = true
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Erro ao percorrer movimentações:", err)
		return
	}

	if !encontrou {
		fmt.Println("Nenhuma movimentação registrada.")
	}
}

func parseDataMovimentacao(data string) (time.Time, error) {
	formatos := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
	}

	for _, formato := range formatos {
		if dataFormatada, err := time.Parse(formato, data); err == nil {
			return dataFormatada, nil
		}
	}

	return time.Time{}, fmt.Errorf("data inválida: %s", data)
}
