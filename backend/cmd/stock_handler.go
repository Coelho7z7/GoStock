package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"gostock/backend/models"
	"gostock/backend/services"
)

func stockHandler(w http.ResponseWriter, r *http.Request) {
	userID, authenticated := userFromSession(r)
	if !authenticated {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	const productsPerPage = 5

	data := struct {
		Products     []models.Product
		Message      string
		Error        string
		PreviousPage int
		NextPage     int
		TotalPages   int
		Page         int
		IsAdmin      bool
	}{IsAdmin: canViewUsersTab(r)}

	messages := map[string]string{
		"entrada":    "Estoque adicionado com sucesso.",
		"saida":      "Saída registrada com sucesso.",
		"cadastrado": "Produto cadastrado com sucesso.",
	}

	if r.Method == http.MethodPost {
		if !requireAdmin(w, r) {
			return
		}
		productID, idErr := strconv.Atoi(r.FormValue("produto_id"))
		quantity, qtyErr := strconv.Atoi(strings.TrimSpace(r.FormValue("quantidade")))
		action := r.FormValue("acao")

		switch {
		case idErr != nil:
			data.Error = "Produto inválido."
		case qtyErr != nil || quantity <= 0:
			data.Error = "Informe uma quantidade válida."
		case action == "entrada":
			if err := services.AddStockWeb(productID, quantity, userID); err != nil {
				data.Error = err.Error()
			} else {
				http.Redirect(w, r, "/estoque?sucesso=entrada", http.StatusSeeOther)
				return
			}
		case action == "saida":
			if err := services.RegisterStockExitWeb(productID, quantity, userID); err != nil {
				data.Error = err.Error()
			} else {
				http.Redirect(w, r, "/estoque?sucesso=saida", http.StatusSeeOther)
				return
			}
		default:
			data.Error = "Ação inválida."
		}
	} else if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	page, _ := strconv.Atoi(r.URL.Query().Get("pagina"))
	if page < 1 {
		page = 1
	}

	products, total, err := services.PaginatedProducts(
		"",
		page,
		productsPerPage,
	)
	if err != nil {
		log.Println("erro em PaginatedProducts:", err)
		http.Error(w, "Erro ao buscar produtos", http.StatusInternalServerError)
		return
	}

	data.Products = products
	data.Page = page

	data.TotalPages = (total + productsPerPage - 1) / productsPerPage
	if data.TotalPages < 1 {
		data.TotalPages = 1
	}

	data.PreviousPage = page - 1
	data.NextPage = page + 1

	data.Message = messages[r.URL.Query().Get("sucesso")]

	tmpl, err := template.ParseFiles("frontend/html/stock.html")
	if err != nil {
		http.Error(w, "Erro ao carregar estoque", http.StatusInternalServerError)
		return
	}

	if data.Error != "" {
		w.WriteHeader(http.StatusBadRequest)
	}

	if err := tmpl.Execute(w, data); err != nil {
		println("O Erro é: ", err.Error())
		http.Error(
			w,
			"Erro ao renderizar estoque de produto:",
			http.StatusInternalServerError,
		)
	}
}
