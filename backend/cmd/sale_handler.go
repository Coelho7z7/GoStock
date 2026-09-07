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

// extractSaleItemsFromForm combines the parallel product_id and quantity
// lists coming from an HTML form into valid sale items, discarding rows
// with an invalid ID or quantity.
func extractSaleItemsFromForm(ids, quantities []string) []services.SaleItem {
	capacity := len(ids)
	if len(quantities) < capacity {
		capacity = len(quantities)
	}
	items := make([]services.SaleItem, 0, capacity)

	for i, idValue := range ids {
		if i >= len(quantities) {
			continue
		}

		id, idErr := strconv.Atoi(strings.TrimSpace(idValue))
		if idErr != nil || id <= 0 {
			continue
		}

		quantity, qtyErr := strconv.Atoi(strings.TrimSpace(quantities[i]))
		if qtyErr != nil || quantity <= 0 {
			continue
		}

		items = append(items, services.SaleItem{ProductID: id, Quantity: quantity})
	}

	return items
}

func saleItemsFromForm(r *http.Request) ([]services.SaleItem, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	ids := r.PostForm["produto_id"]
	quantities := r.PostForm["quantidade"]
	return extractSaleItemsFromForm(ids, quantities), nil
}

// saleHandler renders the POS (with product search) and processes the
// classic form submission (used as a fallback in case the POS
// JavaScript doesn't run). The main sale-finalization flow goes through
// JavaScript calling saleAPIHandler.
func saleHandler(w http.ResponseWriter, r *http.Request) {
	userID, authenticated := userFromSession(r)
	if !authenticated {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		items, err := saleItemsFromForm(r)
		if err != nil {
			http.Error(w, "Dados do formulário inválidos", http.StatusBadRequest)
			return
		}

		if len(items) == 0 {
			http.Error(w, "Selecione pelo menos um produto com quantidade maior que zero", http.StatusBadRequest)
			return
		}

		paymentMethod := r.FormValue("forma_pagamento")

		if err := services.RegisterSaleWeb(items, userID, paymentMethod); err != nil {
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

	products, err := services.GetAllProducts()
	if err != nil {
		http.Error(w, "Erro ao buscar produtos", http.StatusInternalServerError)
		return
	}

	search := strings.TrimSpace(strings.ToLower(r.URL.Query().Get("busca")))
	if search != "" {
		filteredProducts := make([]models.Product, 0, len(products))
		searchID, idErr := strconv.Atoi(search)
		for _, product := range products {
			byID := idErr == nil && product.ID == searchID
			byName := strings.Contains(strings.ToLower(product.Name), search)
			if byID || byName {
				filteredProducts = append(filteredProducts, product)
			}
		}
		products = filteredProducts
	}

	tmpl, err := template.ParseFiles("frontend/html/sales.html")
	if err != nil {
		http.Error(w, "Erro ao carregar PDV", http.StatusInternalServerError)
		return
	}

	message := map[string]string{
		"venda": "Venda registrada com sucesso.",
	}[r.URL.Query().Get("sucesso")]

	if err := tmpl.Execute(w, struct {
		Products []models.Product
		Message  string
		Search   string
		IsAdmin  bool
	}{Products: products, Message: message, Search: search, IsAdmin: canViewUsersTab(r)}); err != nil {
		http.Error(w, "Erro ao renderizar PDV", http.StatusInternalServerError)
	}
}

// saleAPIHandler is the JSON endpoint used by the POS JavaScript to
// finalize a sale, including the payment method chosen in the modal.
func saleAPIHandler(w http.ResponseWriter, r *http.Request) {
	userID, authenticated := userFromSession(r)
	if !authenticated {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var payload struct {
		Items []struct {
			ID       int `json:"id"`
			Quantity int `json:"quantidade"`
		} `json:"items"`
		Discount       float64 `json:"desconto"`
		PaymentMethod  string  `json:"formaPagamento"`
		AmountReceived float64 `json:"valorRecebido"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	items := make([]services.SaleItem, 0, len(payload.Items))
	for _, item := range payload.Items {
		if item.ID <= 0 || item.Quantity <= 0 {
			continue
		}
		items = append(items, services.SaleItem{ProductID: item.ID, Quantity: item.Quantity})
	}

	if len(items) == 0 {
		http.Error(w, "Selecione pelo menos um produto com quantidade maior que zero", http.StatusBadRequest)
		return
	}

	sale, err := services.RegisterSaleWebDetails(items, userID, payload.PaymentMethod)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Message string                  `json:"mensagem"`
		Sale    services.RegisteredSale `json:"venda"`
	}{Message: "Venda registrada com sucesso.", Sale: sale}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Erro ao responder venda", http.StatusInternalServerError)
	}
}
