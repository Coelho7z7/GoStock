package services

import (
	"database/sql"
	"fmt"

	database "gostock/backend/database"
)

type SaleItem struct {
	ProductID int
	Quantity  int
	UnitPrice float64
}

type RegisteredSaleItem struct {
	Name     string  `json:"nome"`
	Quantity int     `json:"quantidade"`
	Price    float64 `json:"preco"`
	Subtotal float64 `json:"subtotal"`
}

type RegisteredSale struct {
	Number        int                   `json:"numero"`
	PaymentMethod string                `json:"formaPagamento"`
	Total         float64               `json:"total"`
	Items         []RegisteredSaleItem  `json:"itens"`
}

// validPaymentMethods lists the payment methods accepted by the POS.
var validPaymentMethods = map[string]bool{
	"dinheiro": true,
	"cartao":   true,
	"pix":      true,
	"debito":   true,
}

// normalizePaymentMethod ensures a known value is always stored, using
// "dinheiro" as the default when the value comes empty or outside the
// expected list.
func normalizePaymentMethod(paymentMethod string) string {
	if validPaymentMethods[paymentMethod] {
		return paymentMethod
	}
	return "dinheiro"
}

func calculateSaleTotal(quantity int, unitPrice float64) float64 {
	return float64(quantity) * unitPrice
}

func calculateSaleItemsTotal(items []SaleItem) float64 {
	total := 0.0
	for _, item := range items {
		total += calculateSaleTotal(item.Quantity, item.UnitPrice)
	}
	return total
}

// RegisterSaleWeb registers a sale with one or more items, debiting the
// stock, recording the sale and the stock-out movement for each product
// within a single transaction.
func RegisterSaleWeb(items []SaleItem, userID int, paymentMethod string) error {
	_, err := RegisterSaleWebDetails(items, userID, paymentMethod)
	return err
}

func RegisterSaleWebDetails(items []SaleItem, userID int, paymentMethod string) (RegisteredSale, error) {
	var sale RegisteredSale

	if len(items) == 0 {
		return sale, fmt.Errorf("a venda deve conter pelo menos um item")
	}

	paymentMethod = normalizePaymentMethod(paymentMethod)
	sale.PaymentMethod = paymentMethod

	tx, err := database.DB.Begin()
	if err != nil {
		return sale, err
	}
	defer tx.Rollback()

	totalValue := 0.0

	for _, item := range items {
		if item.Quantity <= 0 {
			return sale, fmt.Errorf("a quantidade deve ser maior que zero")
		}

		var currentProductID int
		var name string
		var price float64
		var stock int

		err = tx.QueryRow(`
			SELECT id, nome, preco, quantidade
			FROM produtos
			WHERE id = ? AND ativo = 1
		`, item.ProductID).Scan(&currentProductID, &name, &price, &stock)
		if err != nil {
			if err == sql.ErrNoRows {
				return sale, fmt.Errorf("produto não encontrado: %d", item.ProductID)
			}
			return sale, err
		}

		if item.Quantity > stock {
			return sale, fmt.Errorf("estoque insuficiente para %s", name)
		}

		totalValue += calculateSaleTotal(item.Quantity, price)

		_, err = tx.Exec(`
			UPDATE produtos
			SET quantidade = quantidade - ?
			WHERE id = ?
		`, item.Quantity, item.ProductID)
		if err != nil {
			return sale, err
		}

		result, err := tx.Exec(`
			INSERT INTO vendas (produto_id, usuario_id, quantidade, valor_unitario, valor_total, forma_pagamento)
			VALUES (?, ?, ?, ?, ?, ?)
		`, item.ProductID, userID, item.Quantity, price, calculateSaleTotal(item.Quantity, price), paymentMethod)
		if err != nil {
			return sale, err
		}
		if sale.Number == 0 {
			id, idErr := result.LastInsertId()
			if idErr != nil {
				return sale, idErr
			}
			sale.Number = int(id)
		}
		sale.Items = append(sale.Items, RegisteredSaleItem{
			Name: name, Quantity: item.Quantity, Price: price,
			Subtotal: calculateSaleTotal(item.Quantity, price),
		})

		if err := registerMovementTx(tx, item.ProductID, userID, "SAIDA", item.Quantity); err != nil {
			return sale, err
		}
	}

	if totalValue <= 0 {
		return sale, fmt.Errorf("valor da venda deve ser maior que zero")
	}

	sale.Total = totalValue
	return sale, tx.Commit()
}
