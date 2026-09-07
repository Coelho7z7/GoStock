package services

import (
	"math"
	"testing"
)

func TestCalculateSaleTotal(t *testing.T) {
	value := calculateSaleTotal(3, 19.9)
	if math.Abs(value-59.7) > 0.0001 {
		t.Fatalf("valor esperado 59.7, recebido %.2f", value)
	}
}

func TestCalculateSaleItemsTotal(t *testing.T) {
	items := []SaleItem{
		{ProductID: 1, Quantity: 2, UnitPrice: 25.5},
		{ProductID: 2, Quantity: 1, UnitPrice: 12.5},
	}

	value := calculateSaleItemsTotal(items)
	if math.Abs(value-63.5) > 0.0001 {
		t.Fatalf("valor esperado 63.5, recebido %.2f", value)
	}
}

func TestNormalizePaymentMethod(t *testing.T) {
	cases := map[string]string{
		"pix":      "pix",
		"cartao":   "cartao",
		"debito":   "debito",
		"dinheiro": "dinheiro",
		"":         "dinheiro",
		"bitcoin":  "dinheiro",
		"PIX":      "dinheiro",
	}

	for input, expected := range cases {
		if result := normalizePaymentMethod(input); result != expected {
			t.Fatalf("normalizePaymentMethod(%q) = %q, esperado %q", input, result, expected)
		}
	}
}
