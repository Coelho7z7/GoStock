package main

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"gostock/backend/models"
	"gostock/backend/services"
)
func handlerEstoque(w http.ResponseWriter, r *http.Request) {
	usuarioID, ok := usuarioDaSessao(r)
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	dados := struct {
		Produtos []models.Produto
		Mensagem string
		Erro     string
	}{}

	mensagens := map[string]string{
		"entrada":    "Estoque adicionado com sucesso.",
		"saida":      "Saída registrada com sucesso.",
		"cadastrado": "Produto cadastrado com sucesso.",
	}

	if r.Method == http.MethodPost {
		produtoID, idErr := strconv.Atoi(r.FormValue("produto_id"))
		quantidade, qtdErr := strconv.Atoi(strings.TrimSpace(r.FormValue("quantidade")))
		acao := r.FormValue("acao")

		switch {
		case idErr != nil:
			dados.Erro = "Produto inválido."
		case qtdErr != nil || quantidade <= 0:
			dados.Erro = "Informe uma quantidade válida."
		case acao == "entrada":
			if err := services.AdicionarEstoqueWeb(produtoID, quantidade, usuarioID); err != nil {
				dados.Erro = err.Error()
			} else {
				http.Redirect(w, r, "/estoque?sucesso=entrada", http.StatusSeeOther)
				return
			}
		case acao == "saida":
			if err := services.RegistrarSaidaWeb(produtoID, quantidade, usuarioID); err != nil {
				dados.Erro = err.Error()
			} else {
				http.Redirect(w, r, "/estoque?sucesso=saida", http.StatusSeeOther)
				return
			}
		default:
			dados.Erro = "Ação inválida."
		}
	} else if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	produtos, err := services.BuscarTodosProdutos()
	if err != nil {
		http.Error(w, "Erro ao buscar produtos", http.StatusInternalServerError)
		return
	}
	dados.Produtos = produtos
	dados.Mensagem = mensagens[r.URL.Query().Get("sucesso")]

	tmpl, err := template.ParseFiles("frontend/html/estoque.html")
	if err != nil {
		http.Error(w, "Erro ao carregar estoque", http.StatusInternalServerError)
		return
	}

	if dados.Erro != "" {
		w.WriteHeader(http.StatusBadRequest)
	}

	if err := tmpl.Execute(w, dados); err != nil {
		http.Error(w, "Erro ao renderizar estoque", http.StatusInternalServerError)
	}
}
