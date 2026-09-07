package main

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestExtractSaleItemsFromForm(t *testing.T) {
	ids := []string{"1", "2", "3", "4"}
	quantities := []string{"", "2", "0", "5"}

	items := extractSaleItemsFromForm(ids, quantities)
	if len(items) != 2 {
		t.Fatalf("esperava 2 itens válidos, recebeu %d", len(items))
	}

	if items[0].ProductID != 2 || items[0].Quantity != 2 {
		t.Fatalf("primeiro item inesperado: %+v", items[0])
	}

	if items[1].ProductID != 4 || items[1].Quantity != 5 {
		t.Fatalf("segundo item inesperado: %+v", items[1])
	}
}

func TestExtractSaleItemsFromFormNoValidItems(t *testing.T) {
	ids := []string{"1", "2"}
	quantities := []string{"", "0"}

	items := extractSaleItemsFromForm(ids, quantities)
	if len(items) != 0 {
		t.Fatalf("esperava 0 itens válidos, recebeu %d", len(items))
	}
}

func TestSaleItemsFromForm(t *testing.T) {
	body := strings.NewReader("produto_id=10&produto_id=11&produto_id=12&quantidade=1&quantidade=&quantidade=2")
	req := httptest.NewRequest("POST", "/vendas", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	items, err := saleItemsFromForm(req)
	if err != nil {
		t.Fatalf("não deveria haver erro ao processar o formulário: %v", err)
	}

	if len(items) != 2 {
		t.Fatalf("esperava 2 itens válidos, recebeu %d", len(items))
	}

	if items[0].ProductID != 10 || items[0].Quantity != 1 {
		t.Fatalf("primeiro item inesperado: %+v", items[0])
	}

	if items[1].ProductID != 12 || items[1].Quantity != 2 {
		t.Fatalf("segundo item inesperado: %+v", items[1])
	}
}
