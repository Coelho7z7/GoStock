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

// productHandler exibe a lista de produtos e processa o cadastro de
// um novo produto (POST).
func productHandler(w http.ResponseWriter, r *http.Request) {
	userID, authenticated := userFromSession(r)
	if !authenticated {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	const productsPerPage = 5

	data := struct {
		Products     []models.Product
		Name         string
		Quantity     int
		Price        string
		Search       string
		Order        string
		Message      string
		Error        string
		Page         int
		TotalPages   int
		PreviousPage int
		NextPage     int
		IsAdmin      bool
	}{IsAdmin: canViewUsersTab(r)}

	data.Message = map[string]string{
		"cadastrado": "Produto cadastrado com sucesso.",
		"atualizado": "Produto atualizado com sucesso.",
		"removido":   "Produto removido com sucesso.",
		"entrada":    "Estoque adicionado com sucesso.",
		"saida":      "Saída registrada com sucesso.",
	}[r.URL.Query().Get("sucesso")]

	if r.Method == http.MethodPost {
		if !requireAdmin(w, r) {
			return
		}
		data.Name = strings.TrimSpace(r.FormValue("nome"))
		quantityText := strings.TrimSpace(r.FormValue("quantidade"))
		data.Price = strings.TrimSpace(r.FormValue("preco"))

		quantity, quantityErr := strconv.Atoi(quantityText)
		price, priceErr := strconv.ParseFloat(data.Price, 64)

		if data.Name == "" {
			data.Error = "Informe o nome do produto."
		} else if quantityErr != nil || quantity < 0 {
			data.Error = "Informe uma quantidade válida."
		} else if priceErr != nil || price < 0 {
			data.Error = "Informe um preço válido."
		} else if err := services.CreateProductWeb(
			data.Name,
			quantity,
			price,
			userID,
		); err != nil {
			data.Error = err.Error()
		} else {
			http.Redirect(
				w,
				r,
				"/produtos?sucesso=cadastrado",
				http.StatusSeeOther,
			)
			return
		}

		data.Quantity = quantity
	} else if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("pagina"))
	if page < 1 {
		page = 1
	}

	data.Search = strings.TrimSpace(r.URL.Query().Get("busca"))
	data.Order = r.URL.Query().Get("ordem")
	if data.Order == "" {
		data.Order = "recentes"
	}
	products, total, err := services.PaginatedSortedProducts(
		data.Search,
		page,
		productsPerPage,
		data.Order,
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

	tmpl, err := template.ParseFiles("frontend/html/products.html")
	if err != nil {
		http.Error(
			w,
			"Erro ao carregar produtos",
			http.StatusInternalServerError,
		)
		return
	}

	if data.Error != "" {
		w.WriteHeader(http.StatusBadRequest)
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(
			w,
			"Erro ao renderizar produtos",
			http.StatusInternalServerError,
		)
	}
}

// editProductHandler exibe a tela de edição/remoção de produtos e
// processa as ações de atualizar ou remover (POST).
func editProductHandler(w http.ResponseWriter, r *http.Request) {
	userID, authenticated := userFromSession(r)
	if !authenticated {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	const productsPerPage = 5

	messages := map[string]string{
		"atualizado": "Produto atualizado com sucesso.",
		"removido":   "Produto removido com sucesso.",
	}

	data := struct {
		Products     []models.Product
		Message      string
		Error        string
		Page         int
		TotalPages   int
		PreviousPage int
		NextPage     int
		IsAdmin      bool
	}{
		Message: messages[r.URL.Query().Get("sucesso")],
		IsAdmin: canViewUsersTab(r),
	}

	if r.Method == http.MethodPost {
		if !requireAdmin(w, r) {
			return
		}
		productID, idErr := strconv.Atoi(r.FormValue("produto_id"))

		if idErr != nil {
			data.Error = "Produto inválido."

		} else if r.FormValue("acao") == "remover" {
			opErr := services.DeleteProductWeb(productID)

			if opErr == nil {
				http.Redirect(
					w,
					r,
					"/alterar-produto?sucesso=removido",
					http.StatusSeeOther,
				)
				return
			}

			data.Error = opErr.Error()

		} else if r.FormValue("acao") == "atualizar" {
			name := strings.TrimSpace(r.FormValue("nome"))

			price, priceErr := strconv.ParseFloat(
				strings.TrimSpace(r.FormValue("preco")),
				64,
			)

			if priceErr != nil {
				data.Error = "Informe um preço válido."
			} else {
				opErr := services.UpdateProductWeb(
					productID,
					name,
					price,
					userID,
				)

				if opErr == nil {
					http.Redirect(
						w,
						r,
						"/alterar-produto?sucesso=atualizado",
						http.StatusSeeOther,
					)
					return
				}

				data.Error = opErr.Error()
			}

		} else {
			data.Error = "Ação inválida."
		}

	} else if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Paginação
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
		http.Error(
			w,
			"Erro ao buscar produtos",
			http.StatusInternalServerError,
		)
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

	tmpl, err := template.ParseFiles(
		"frontend/html/edit_product.html",
	)

	if err != nil {
		http.Error(
			w,
			"Erro ao carregar alteração de produto",
			http.StatusInternalServerError,
		)
		return
	}

	if data.Error != "" {
		w.WriteHeader(http.StatusBadRequest)
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(
			w,
			"Erro ao renderizar alteração de produto",
			http.StatusInternalServerError,
		)
	}
}
