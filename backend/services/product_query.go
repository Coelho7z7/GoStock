package services

import (
	"bufio"
	"fmt"
	"strings"

	database "gostock/backend/database"
	"gostock/backend/models"
	"gostock/backend/utils"
)

func ListProducts() {
	rows, err := database.DB.Query(`
		SELECT id, nome, preco, quantidade
		FROM produtos
	`)
	if err != nil {
		fmt.Println("Erro ao buscar produtos:", err)
		return
	}
	defer rows.Close()

	found := false

	for rows.Next() {
		var product models.Product

		if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Quantity); err != nil {
			fmt.Println("Erro ao ler produto:", err)
			return
		}

		fmt.Println("ID:", product.ID)
		fmt.Println("Nome:", product.Name)
		fmt.Println("Quantidade:", product.Quantity)
		fmt.Println("Preço:", product.Price)
		fmt.Println("----------------------")
		found = true
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Erro ao percorrer produtos:", err)
		return
	}

	if !found {
		fmt.Println("Nenhum produto cadastrado.")
	}
}

func FindProduct(reader *bufio.Reader) {
	search := strings.TrimSpace(utils.ReadText(reader, "Digite o nome do produto: "))

	rows, err := database.DB.Query(`
		SELECT id, nome, preco, quantidade
		FROM produtos
		WHERE nome LIKE ?
	`, "%"+search+"%")
	if err != nil {
		fmt.Println("Erro ao buscar produto:", err)
		return
	}
	defer rows.Close()

	found := false

	for rows.Next() {
		var product models.Product

		if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Quantity); err != nil {
			fmt.Println("Erro ao ler produto:", err)
			return
		}

		fmt.Println("Produto encontrado!")
		fmt.Println("ID:", product.ID)
		fmt.Println("Nome:", product.Name)
		fmt.Println("Quantidade:", product.Quantity)
		fmt.Println("Preço:", product.Price)
		found = true
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Erro ao percorrer produtos:", err)
		return
	}

	if !found {
		fmt.Println("Produto não encontrado.")
	}
}

func GetAllProducts() ([]models.Product, error) {
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

	var products []models.Product

	for rows.Next() {
		var product models.Product

		if err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.Quantity,
		); err != nil {
			return nil, err
		}

		product.StockStatus = stockStatus(product.Quantity)
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
func UpdateProductWeb(productID int, name string, price float64, userID int) error {
	name = strings.TrimSpace(name)
	if !utils.ValidateName(name) {
		return fmt.Errorf("o nome do produto é obrigatório")
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
		UPDATE produtos
		SET nome = ?, preco = ?
		WHERE id = ? AND ativo = 1
	`, name, price, productID)
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

	if err := registerMovementTx(tx, productID, userID, "ATUALIZACAO", 0); err != nil {
		return err
	}

	return tx.Commit()
}

func DeleteProductWeb(productID int) error {
	result, err := database.DB.Exec(`
		UPDATE produtos
		SET ativo = 0
		WHERE id = ? AND ativo = 1
	`, productID)
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

	return nil
}

func stockStatus(quantity int) string {
	switch {
	case quantity == 0:
		return "empty"
	case quantity <= 5:
		return "low"
	default:
		return "normal"
	}
}

func UpdateProduct(reader *bufio.Reader, userID int) {
	id, err := utils.ReadInt(reader, "Digite o ID do produto: ")
	if err != nil {
		fmt.Println("ID inválido.")
		return
	}

	var productID int
	var currentName string
	var currentPrice float64

	err = database.DB.QueryRow(`
		SELECT id, nome, preco
		FROM produtos
		WHERE id = ? AND ativo = 1
	`, id).Scan(&productID, &currentName, &currentPrice)

	if err != nil {
		fmt.Println("Produto não encontrado.")
		return
	}

	fmt.Println("Produto encontrado.")
	fmt.Println("Nome atual:", currentName)
	fmt.Println("Preço atual:", currentPrice)

	newName := utils.ReadValidName(reader)
	newPrice := utils.ReadValidPrice(reader, "Digite o novo preço: ")

	_, err = database.DB.Exec(`
		UPDATE produtos
		SET nome = ?, preco = ?
		WHERE id = ?
	`, newName, newPrice, productID)
	if err != nil {
		fmt.Println("Erro ao atualizar produto:", err)
		return
	}

	if err := registerMovement(productID, userID, "ATUALIZACAO", 0); err != nil {
		fmt.Println("Erro ao registrar atualização:", err)
		return
	}

	fmt.Println("Produto atualizado com sucesso!")
}

// PaginatedProducts fetches a page of active products, optionally
// filtering by name (partial match). page starts at 1.
func PaginatedProducts(search string, page int, perPage int) ([]models.Product, int, error) {
	return PaginatedSortedProducts(search, page, perPage, "recentes")
}

func PaginatedSortedProducts(search string, page int, perPage int, order string) ([]models.Product, int, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	search = strings.TrimSpace(search)
	nameFilter := "%" + search + "%"

	var total int
	if err := database.DB.QueryRow(`
		SELECT COUNT(*)
		FROM produtos
		WHERE ativo = 1 AND nome LIKE ?
	`, nameFilter).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage

	sorting := "id DESC"
	switch order {
	case "nome":
		sorting = "nome COLLATE NOCASE ASC"
	case "preco":
		sorting = "preco ASC"
	case "estoque":
		sorting = "quantidade ASC"
	}
	rows, err := database.DB.Query(`
		SELECT id, nome, preco, quantidade
		FROM produtos
		WHERE ativo = 1 AND nome LIKE ?
		ORDER BY `+sorting+`
		LIMIT ? OFFSET ?
	`, nameFilter, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []models.Product

	for rows.Next() {
		var product models.Product

		if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Quantity); err != nil {
			return nil, 0, err
		}

		product.StockStatus = stockStatus(product.Quantity)
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
