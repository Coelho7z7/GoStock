package main

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"gostock/backend/models"
	"gostock/backend/services"
)

// handlerProdutos exibe a lista de produtos e processa o cadastro de
// um novo produto (POST).
func handlerProdutos(w http.ResponseWriter, r *http.Request) {
	usuarioID, ok := usuarioDaSessao(r)
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	dados := struct {
		Produtos   []models.Produto
		Nome       string
		Quantidade int
		Preco      string
		Mensagem   string
		Erro       string
	}{}

	dados.Mensagem = map[string]string{
		"cadastrado": "Produto cadastrado com sucesso.",
		"atualizado": "Produto atualizado com sucesso.",
		"removido":   "Produto removido com sucesso.",
		"entrada":    "Estoque adicionado com sucesso.",
		"saida":      "Saída registrada com sucesso.",
	}[r.URL.Query().Get("sucesso")]

	if r.Method == http.MethodPost {
		dados.Nome = strings.TrimSpace(r.FormValue("nome"))
		quantidadeTexto := strings.TrimSpace(r.FormValue("quantidade"))
		dados.Preco = strings.TrimSpace(r.FormValue("preco"))

		quantidade, quantidadeErr := strconv.Atoi(quantidadeTexto)
		preco, precoErr := strconv.ParseFloat(dados.Preco, 64)

		if dados.Nome == "" {
			dados.Erro = "Informe o nome do produto."
		} else if quantidadeErr != nil || quantidade < 0 {
			dados.Erro = "Informe uma quantidade válida."
		} else if precoErr != nil || preco < 0 {
			dados.Erro = "Informe um preço válido."
		} else if err := services.CadastrarProdutoWeb(
			dados.Nome,
			quantidade,
			preco,
			usuarioID,
		); err != nil {
			dados.Erro = err.Error()
		} else {
			http.Redirect(w, r, "/produtos?sucesso=cadastrado", http.StatusSeeOther)
			return
		}

		dados.Quantidade = quantidade
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

	tmpl, err := template.ParseFiles("frontend/html/produtos.html")
	if err != nil {
		http.Error(w, "Erro ao carregar produtos", http.StatusInternalServerError)
		return
	}

	if dados.Erro != "" {
		w.WriteHeader(http.StatusBadRequest)
	}

	if err := tmpl.Execute(w, dados); err != nil {
		http.Error(w, "Erro ao renderizar produtos", http.StatusInternalServerError)
	}
}

// handlerAlterarProduto exibe a tela de edição/remoção de produtos e
// processa as ações de atualizar ou remover (POST).
func handlerAlterarProduto(w http.ResponseWriter, r *http.Request) {
	usuarioID, ok := usuarioDaSessao(r)
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	mensagens := map[string]string{
		"atualizado": "Produto atualizado com sucesso.",
		"removido":   "Produto removido com sucesso.",
	}
	dados := struct {
		Produtos []models.Produto
		Mensagem string
		Erro     string
	}{Mensagem: mensagens[r.URL.Query().Get("sucesso")]}

	if r.Method == http.MethodPost {
		produtoID, idErr := strconv.Atoi(r.FormValue("produto_id"))
		dados.Erro = "Produto inválido."
		if idErr == nil && r.FormValue("acao") == "remover" {
			opErr := services.RemoverProdutoWeb(produtoID)
			if opErr == nil {
				http.Redirect(w, r, "/produtos?sucesso=removido", http.StatusSeeOther)
				return
			}
			dados.Erro = opErr.Error()
		} else if idErr == nil && r.FormValue("acao") == "atualizar" {
			nome := strings.TrimSpace(r.FormValue("nome"))
			preco, precoErr := strconv.ParseFloat(strings.TrimSpace(r.FormValue("preco")), 64)
			if precoErr != nil {
				dados.Erro = "Informe um preço válido."
			} else {
				opErr := services.AtualizarProdutoWeb(produtoID, nome, preco, usuarioID)
				if opErr == nil {
					http.Redirect(w, r, "/produtos?sucesso=atualizado", http.StatusSeeOther)
					return
				}
				dados.Erro = opErr.Error()
			}
		}
	} else if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	dados.Produtos, _ = services.BuscarTodosProdutos()
	tmpl, err := template.ParseFiles("frontend/html/alterar_produto.html")
	if err != nil {
		http.Error(w, "Erro ao carregar alteração de produto", http.StatusInternalServerError)
		return
	}
	if dados.Erro != "" {
		w.WriteHeader(http.StatusBadRequest)
	}
	if err := tmpl.Execute(w, dados); err != nil {
		http.Error(w, "Erro ao renderizar alteração de produto", http.StatusInternalServerError)
	}
}
