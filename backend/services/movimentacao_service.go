package services

import (
	"fmt"
	"time"

	database "gostock/backend/Database"
)

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
