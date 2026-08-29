package main

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestExtrairItensVendaFormulario(t *testing.T) {
	ids := []string{"1", "2", "3", "4"}
	quantidades := []string{"", "2", "0", "5"}

	itens := extrairItensVendaFormulario(ids, quantidades)
	if len(itens) != 2 {
		t.Fatalf("esperava 2 itens válidos, recebeu %d", len(itens))
	}

	if itens[0].ProdutoID != 2 || itens[0].Quantidade != 2 {
		t.Fatalf("primeiro item inesperado: %+v", itens[0])
	}

	if itens[1].ProdutoID != 4 || itens[1].Quantidade != 5 {
		t.Fatalf("segundo item inesperado: %+v", itens[1])
	}
}

func TestExtrairItensVendaFormularioSemItensValidos(t *testing.T) {
	ids := []string{"1", "2"}
	quantidades := []string{"", "0"}

	itens := extrairItensVendaFormulario(ids, quantidades)
	if len(itens) != 0 {
		t.Fatalf("esperava 0 itens válidos, recebeu %d", len(itens))
	}
}

func TestItensVendaDoFormulario(t *testing.T) {
	body := strings.NewReader("produto_id=10&produto_id=11&produto_id=12&quantidade=1&quantidade=&quantidade=2")
	req := httptest.NewRequest("POST", "/vendas", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	itens, err := itensVendaDoFormulario(req)
	if err != nil {
		t.Fatalf("não deveria haver erro ao processar o formulário: %v", err)
	}

	if len(itens) != 2 {
		t.Fatalf("esperava 2 itens válidos, recebeu %d", len(itens))
	}

	if itens[0].ProdutoID != 10 || itens[0].Quantidade != 1 {
		t.Fatalf("primeiro item inesperado: %+v", itens[0])
	}

	if itens[1].ProdutoID != 12 || itens[1].Quantidade != 2 {
		t.Fatalf("segundo item inesperado: %+v", itens[1])
	}
}
