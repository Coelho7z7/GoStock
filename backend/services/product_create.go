package services

import (
	"bufio"
	"fmt"
	"strings"

	database "gostock/backend/database"
	"gostock/backend/utils"
)

func CreateProductWeb(name string, quantity int, price float64, userID int) error {
	name = strings.TrimSpace(name)
	if !utils.ValidateName(name) {
		return fmt.Errorf("o nome do produto é obrigatório")
	}
	if !utils.ValidateQuantity(quantity) {
		return fmt.Errorf("a quantidade não pode ser negativa")
	}
	if !utils.ValidatePrice(price) {
		return fmt.Errorf("o preço não pode ser negativo")
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.Exec(`
		INSERT INTO produtos (nome, quantidade, preco)
		VALUES (?, ?, ?)
	`, name, quantity, price)
	if err != nil {
		return err
	}

	productID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	if err := registerMovementTx(tx, int(productID), userID, "ENTRADA", quantity); err != nil {
		return err
	}

	return tx.Commit()
}

func CreateProduct(reader *bufio.Reader, userID int) {
	name := utils.ReadValidName(reader)
	quantity := utils.ReadValidQuantity(reader, "Quantidade: ")
	price := utils.ReadValidPrice(reader, "Preço: ")

	tx, err := database.DB.Begin()
	if err != nil {
		fmt.Println("Erro ao iniciar cadastro:", err)
		return
	}
	defer tx.Rollback()

	result, err := tx.Exec(`
		INSERT INTO produtos (nome, preco, quantidade)
		VALUES (?, ?, ?)
	`, name, price, quantity)

	if err != nil {
		fmt.Println("Erro ao cadastrar produto:", err)
		return
	}

	productID, err := result.LastInsertId()
	if err != nil {
		fmt.Println("Erro ao obter ID do produto:", err)
		return
	}

	if err := registerMovementTx(tx, int(productID), userID, "ENTRADA", quantity); err != nil {
		fmt.Println("Erro ao registrar movimentação:", err)
		return
	}

	if err := tx.Commit(); err != nil {
		fmt.Println("Erro ao confirmar cadastro:", err)
		return
	}

	fmt.Println("Produto cadastrado com sucesso!")
}

func DeleteProduct(reader *bufio.Reader) {
	id, err := utils.ReadInt(reader, "Digite o ID do produto: ")
	if err != nil {
		fmt.Println("ID inválido.")
		return
	}

	result, err := database.DB.Exec(`
		DELETE FROM produtos
		WHERE id = ?
	`, id)

	if err != nil {
		fmt.Println("Erro ao remover produto:", err)
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		fmt.Println("Erro ao verificar remoção:", err)
		return
	}

	if rows > 0 {
		fmt.Println("Produto removido com sucesso.")
	} else {
		fmt.Println("Produto não encontrado.")
	}

}
