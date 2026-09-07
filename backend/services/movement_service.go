package services

import (
	"fmt"
	"time"

	database "gostock/backend/database"
	"gostock/backend/models"
)

func GetMovementsWeb() ([]models.Movement, error) {
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

	var movements []models.Movement
	for rows.Next() {
		var movement models.Movement
		var date string

		if err := rows.Scan(
			&movement.ID,
			&movement.ProductID,
			&movement.UserID,
			&date,
			&movement.Product,
			&movement.User,
			&movement.Type,
			&movement.Quantity,
		); err != nil {
			return nil, err
		}

		formattedDate, err := parseMovementDate(date)
		if err != nil {
			return nil, err
		}

		movement.Date = formattedDate
		movement.FormattedDate = formattedDate.Local().Format("02/01/2006")
		movement.FormattedTime = formattedDate.Local().Format("15:04")
		movement.FormattedType = map[string]string{
			"ENTRADA":     "Entrada",
			"SAIDA":       "Saída",
			"ATUALIZACAO": "Atualização",
		}[movement.Type]
		if movement.FormattedType == "" {
			movement.FormattedType = movement.Type
		}
		movements = append(movements, movement)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movements, nil
}

func registerMovement(productID int, userID int, movementType string, quantity int) error {
	_, err := database.DB.Exec(`
		INSERT INTO movimentacoes
		(produto_id, usuario_id, tipo, quantidade)
		VALUES (?, ?, ?, ?)
	`, productID, userID, movementType, quantity)

	return err
}

func ListMovements() {
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

	found := false

	for rows.Next() {
		var date string
		var product string
		var user string
		var movementType string
		var quantity int

		if err := rows.Scan(&date, &product, &user, &movementType, &quantity); err != nil {
			fmt.Println("Erro ao ler movimentação:", err)
			return
		}

		formattedDate, err := parseMovementDate(date)
		if err != nil {
			fmt.Println("Erro ao formatar data:", err)
			return
		}

		fmt.Println("========== MOVIMENTAÇÃO ==========")
		fmt.Println("Produto:", product)
		fmt.Println("Usuário:", user)
		fmt.Println("Tipo:", movementType)

		if movementType == "ATUALIZACAO" {
			fmt.Println("Quantidade: -")
		} else {
			fmt.Println("Quantidade:", quantity)
		}

		fmt.Println("Data:", formattedDate.Local().Format("02/01/2006 15:04:05"))
		fmt.Println("==================================")

		found = true
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Erro ao percorrer movimentações:", err)
		return
	}

	if !found {
		fmt.Println("Nenhuma movimentação registrada.")
	}
}

func parseMovementDate(date string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
	}

	for _, format := range formats {
		if formattedDate, err := time.Parse(format, date); err == nil {
			return formattedDate, nil
		}
	}

	return time.Time{}, fmt.Errorf("data inválida: %s", date)
}
