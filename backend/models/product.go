package models

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"nome"`
	Price       float64 `json:"preco"`
	Quantity    int     `json:"quantidade"`
	StockStatus string  `json:"status_estoque"`
}
