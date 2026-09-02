package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"gostock/backend/models"
	"gostock/backend/services"
)

// extrairItensVendaFormulario junta as listas paralelas de produto_id e
// quantidade vindas de um formulário HTML em itens de venda válidos,
// descartando linhas com ID ou quantidade inválidos.
func extrairItensVendaFormulario(ids, quantidades []string) []services.ItemVenda {
	capacidade := len(ids)
	if len(quantidades) < capacidade {
		capacidade = len(quantidades)
	}
	itens := make([]services.ItemVenda, 0, capacidade)

	for i, idValor := range ids {
		if i >= len(quantidades) {
			continue
		}

		id, errID := strconv.Atoi(strings.TrimSpace(idValor))
		if errID != nil || id <= 0 {
			continue
		}

		quantidade, errQtd := strconv.Atoi(strings.TrimSpace(quantidades[i]))
		if errQtd != nil || quantidade <= 0 {
			continue
		}

		itens = append(itens, services.ItemVenda{ProdutoID: id, Quantidade: quantidade})
	}

	return itens
}

func itensVendaDoFormulario(r *http.Request) ([]services.ItemVenda, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	ids := r.PostForm["produto_id"]
	quantidades := r.PostForm["quantidade"]
	return extrairItensVendaFormulario(ids, quantidades), nil
}

// handlerVendas exibe o PDV (com busca de produtos) e processa o envio
// clássico do formulário (usado como alternativa caso o JavaScript do
// PDV não rode). O fluxo principal de finalização de venda é feito via
// JavaScript chamando handlerApiVendas.
func handlerVendas(w http.ResponseWriter, r *http.Request) {
	usuarioID, true := usuarioDaSessao(r)
	if !true {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		if !exigirAdmin(w, r) {
			return
		}
		itens, err := itensVendaDoFormulario(r)
		if err != nil {
			http.Error(w, "Dados do formulário inválidos", http.StatusBadRequest)
			return
		}

		if len(itens) == 0 {
			http.Error(w, "Selecione pelo menos um produto com quantidade maior que zero", http.StatusBadRequest)
			return
		}

		formaPagamento := r.FormValue("forma_pagamento")

		if err := services.RegistrarVendaCompletaWeb(itens, usuarioID, formaPagamento); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, "/vendas?sucesso=venda", http.StatusSeeOther)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	produtos, err := services.BuscarTodosProdutos()
	if err != nil {
		http.Error(w, "Erro ao buscar produtos", http.StatusInternalServerError)
		return
	}

	busca := strings.TrimSpace(strings.ToLower(r.URL.Query().Get("busca")))
	if busca != "" {
		produtosFiltrados := make([]models.Produto, 0, len(produtos))
		idBusca, idErr := strconv.Atoi(busca)
		for _, produto := range produtos {
			porID := idErr == nil && produto.ID == idBusca
			porNome := strings.Contains(strings.ToLower(produto.Nome), busca)
			if porID || porNome {
				produtosFiltrados = append(produtosFiltrados, produto)
			}
		}
		produtos = produtosFiltrados
	}

	tmpl, err := template.ParseFiles("frontend/html/vendas.html")
	if err != nil {
		http.Error(w, "Erro ao carregar PDV", http.StatusInternalServerError)
		return
	}

	mensagem := map[string]string{
		"venda": "Venda registrada com sucesso.",
	}[r.URL.Query().Get("sucesso")]

	if err := tmpl.Execute(w, struct {
		Produtos []models.Produto
		Mensagem string
		Busca    string
	}{Produtos: produtos, Mensagem: mensagem, Busca: busca}); err != nil {
		http.Error(w, "Erro ao renderizar PDV", http.StatusInternalServerError)
	}
}

// handlerApiVendas é o endpoint JSON usado pelo JavaScript do PDV para
// finalizar a venda, incluindo a forma de pagamento escolhida no modal.
func handlerApiVendas(w http.ResponseWriter, r *http.Request) {
	usuarioID, true := usuarioDaSessao(r)
	if !true {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	if !exigirAdmin(w, r) {
		return
	}

	var payload struct {
		Items []struct {
			ID         int `json:"id"`
			Quantidade int `json:"quantidade"`
		} `json:"items"`
		Desconto       float64 `json:"desconto"`
		FormaPagamento string  `json:"formaPagamento"`
		ValorRecebido  float64 `json:"valorRecebido"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	itens := make([]services.ItemVenda, 0, len(payload.Items))
	for _, item := range payload.Items {
		if item.ID <= 0 || item.Quantidade <= 0 {
			continue
		}
		itens = append(itens, services.ItemVenda{ProdutoID: item.ID, Quantidade: item.Quantidade})
	}

	if len(itens) == 0 {
		http.Error(w, "Selecione pelo menos um produto com quantidade maior que zero", http.StatusBadRequest)
		return
	}

	venda, err := services.RegistrarVendaCompletaWebDetalhes(itens, usuarioID, payload.FormaPagamento)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	resposta := struct {
		Mensagem string                   `json:"mensagem"`
		Venda    services.VendaRegistrada `json:"venda"`
	}{Mensagem: "Venda registrada com sucesso.", Venda: venda}
	if err := json.NewEncoder(w).Encode(resposta); err != nil {
		http.Error(w, "Erro ao responder venda", http.StatusInternalServerError)
	}
}
