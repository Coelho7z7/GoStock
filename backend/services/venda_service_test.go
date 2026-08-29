package services

import (
	"math"
	"testing"
)

func TestCalcularValorTotalVenda(t *testing.T) {
	valor := calcularValorTotalVenda(3, 19.9)
	if math.Abs(valor-59.7) > 0.0001 {
		t.Fatalf("valor esperado 59.7, recebido %.2f", valor)
	}
}

func TestCalcularValorTotalItensVenda(t *testing.T) {
	itens := []ItemVenda{
		{ProdutoID: 1, Quantidade: 2, PrecoUnitario: 25.5},
		{ProdutoID: 2, Quantidade: 1, PrecoUnitario: 12.5},
	}

	valor := calcularValorTotalItensVenda(itens)
	if math.Abs(valor-63.5) > 0.0001 {
		t.Fatalf("valor esperado 63.5, recebido %.2f", valor)
	}
}
