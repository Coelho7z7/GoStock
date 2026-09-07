package services

import (
	"bufio"
	"database/sql"
	"fmt"

	database "gostock/backend/database"
	"gostock/backend/utils"
)

func AddStock(reader *bufio.Reader, userID int) {
	id, err := utils.ReadInt(reader, "Digite o ID do produto: ")
	if err != nil {
		fmt.Println("ID inválido.")
		return
	}

	productID, name, currentStock, err := getStockData(id)
	if err != nil {
		fmt.Println("Produto não encontrado.")
		return
	}

	fmt.Println("Produto:", name)
	fmt.Println("Estoque atual:", currentStock)

	quantity := utils.ReadValidQuantity(reader, "Quantidade que chegou: ")
	if quantity <= 0 {
		fmt.Println("A quantidade deve ser maior que zero.")
		return
	}

	if err := AddStockWeb(productID, quantity, userID); err != nil {
		fmt.Println("Erro ao atualizar estoque:", err)
		return
	}

	fmt.Println("Estoque atualizado com sucesso!")
}

func AddStockWeb(productID int, quantity int, userID int) error {
	if quantity <= 0 {
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
	`, quantity, productID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("produto não encontrado")
	}

	if err := registerMovementTx(tx, productID, userID, "ENTRADA", quantity); err != nil {
		return err
	}

	return tx.Commit()
}

func RegisterStockExit(reader *bufio.Reader, userID int) {
	id, err := utils.ReadInt(reader, "Digite o ID do produto: ")
	if err != nil {
		fmt.Println("ID inválido.")
		return
	}

	productID, name, currentStock, err := getStockData(id)
	if err != nil {
		fmt.Println("Produto não encontrado.")
		return
	}

	fmt.Println("Produto:", name)
	fmt.Println("Estoque atual:", currentStock)

	quantity := utils.ReadValidQuantity(reader, "Quantidade que saiu: ")
	if quantity <= 0 {
		fmt.Println("A quantidade deve ser maior que zero.")
		return
	}

	if err := RegisterStockExitWeb(productID, quantity, userID); err != nil {
		if err.Error() == "estoque insuficiente" {
			fmt.Println("Estoque insuficiente.")
			fmt.Println("Estoque disponível:", currentStock)
			return
		}

		fmt.Println("Erro ao registrar saída:", err)
		return
	}

	fmt.Println("Saída registrada com sucesso!")
}

func RegisterStockExitWeb(productID int, quantity int, userID int) error {
	if quantity <= 0 {
		return fmt.Errorf("a quantidade deve ser maior que zero")
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var stock int

	err = tx.QueryRow(`
		SELECT quantidade
		FROM produtos
		WHERE id = ? AND ativo = 1
	`, productID).Scan(&stock)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("produto não encontrado")
		}

		return err
	}

	if quantity > stock {
		return fmt.Errorf("estoque insuficiente")
	}

	_, err = tx.Exec(`
		UPDATE produtos
		SET quantidade = quantidade - ?
		WHERE id = ?
	`, quantity, productID)

	if err != nil {
		return err
	}

	if err := registerMovementTx(
		tx,
		productID,
		userID,
		"SAIDA",
		quantity,
	); err != nil {
		return err
	}

	return tx.Commit()
}
func getStockData(productID int) (int, string, int, error) {
	var name string
	var quantity int

	err := database.DB.QueryRow(`
		SELECT id, nome, quantidade
		FROM produtos
		WHERE id = ?
	`, productID).Scan(&productID, &name, &quantity)

	return productID, name, quantity, err
}

func registerMovementTx(tx *sql.Tx, productID int, userID int, movementType string, quantity int) error {
	_, err := tx.Exec(`
		INSERT INTO movimentacoes
		(produto_id, usuario_id, tipo, quantidade)
		VALUES (?, ?, ?, ?)
	`, productID, userID, movementType, quantity)
	return err
}
