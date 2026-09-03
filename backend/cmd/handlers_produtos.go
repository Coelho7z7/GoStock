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

// handlerProdutos exibe a lista de produtos e processa o cadastro de
// um novo produto (POST).
func handlerProdutos(w http.ResponseWriter, r *http.Request) {
	usuarioID, true := usuarioDaSessao(r)
	if !true {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	const produtosPorPagina = 5

	dados := struct {
		Produtos       []models.Produto
		Nome           string
		Quantidade     int
		Preco          string
		Busca          string
		Ordem          string
		Mensagem       string
		Erro           string
		Pagina         int
		TotalPaginas   int
		PaginaAnterior int
		PaginaProxima  int
		EhAdmin        bool
	}{EhAdmin: usuarioEhAdmin(r)}

	dados.Mensagem = map[string]string{
		"cadastrado": "Produto cadastrado com sucesso.",
		"atualizado": "Produto atualizado com sucesso.",
		"removido":   "Produto removido com sucesso.",
		"entrada":    "Estoque adicionado com sucesso.",
		"saida":      "Saída registrada com sucesso.",
	}[r.URL.Query().Get("sucesso")]

	if r.Method == http.MethodPost {
		if !exigirAdmin(w, r) {
			return
		}
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
			http.Redirect(
				w,
				r,
				"/produtos?sucesso=cadastrado",
				http.StatusSeeOther,
			)
			return
		}

		dados.Quantidade = quantidade
	} else if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	pagina, _ := strconv.Atoi(r.URL.Query().Get("pagina"))
	if pagina < 1 {
		pagina = 1
	}

	dados.Busca = strings.TrimSpace(r.URL.Query().Get("busca"))
	dados.Ordem = r.URL.Query().Get("ordem")
	if dados.Ordem == "" {
		dados.Ordem = "recentes"
	}
	produtos, total, err := services.ProdutosPaginadosOrdenados(
		dados.Busca,
		pagina,
		produtosPorPagina,
		dados.Ordem,
	)
	if err != nil {
		log.Println("erro em ProdutosPaginados:", err)
		http.Error(w, "Erro ao buscar produtos", http.StatusInternalServerError)
		return
	}

	dados.Produtos = produtos
	dados.Pagina = pagina

	dados.TotalPaginas = (total + produtosPorPagina - 1) / produtosPorPagina
	if dados.TotalPaginas < 1 {
		dados.TotalPaginas = 1
	}

	dados.PaginaAnterior = pagina - 1
	dados.PaginaProxima = pagina + 1

	tmpl, err := template.ParseFiles("frontend/html/produtos.html")
	if err != nil {
		http.Error(
			w,
			"Erro ao carregar produtos",
			http.StatusInternalServerError,
		)
		return
	}

	if dados.Erro != "" {
		w.WriteHeader(http.StatusBadRequest)
	}

	if err := tmpl.Execute(w, dados); err != nil {
		http.Error(
			w,
			"Erro ao renderizar produtos",
			http.StatusInternalServerError,
		)
	}
}

// handlerAlterarProduto exibe a tela de edição/remoção de produtos e
// processa as ações de atualizar ou remover (POST).
func handlerAlterarProduto(w http.ResponseWriter, r *http.Request) {
	usuarioID, true := usuarioDaSessao(r)
	if !true {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	const produtosPorPagina = 5

	mensagens := map[string]string{
		"atualizado": "Produto atualizado com sucesso.",
		"removido":   "Produto removido com sucesso.",
	}

	dados := struct {
		Produtos       []models.Produto
		Mensagem       string
		Erro           string
		Pagina         int
		TotalPaginas   int
		PaginaAnterior int
		PaginaProxima  int
		EhAdmin        bool
	}{
		Mensagem: mensagens[r.URL.Query().Get("sucesso")],
		EhAdmin:  usuarioEhAdmin(r),
	}

	if r.Method == http.MethodPost {
		if !exigirAdmin(w, r) {
			return
		}
		produtoID, idErr := strconv.Atoi(r.FormValue("produto_id"))

		if idErr != nil {
			dados.Erro = "Produto inválido."

		} else if r.FormValue("acao") == "remover" {
			opErr := services.RemoverProdutoWeb(produtoID)

			if opErr == nil {
				http.Redirect(
					w,
					r,
					"/alterar-produto?sucesso=removido",
					http.StatusSeeOther,
				)
				return
			}

			dados.Erro = opErr.Error()

		} else if r.FormValue("acao") == "atualizar" {
			nome := strings.TrimSpace(r.FormValue("nome"))

			preco, precoErr := strconv.ParseFloat(
				strings.TrimSpace(r.FormValue("preco")),
				64,
			)

			if precoErr != nil {
				dados.Erro = "Informe um preço válido."
			} else {
				opErr := services.AtualizarProdutoWeb(
					produtoID,
					nome,
					preco,
					usuarioID,
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

				dados.Erro = opErr.Error()
			}

		} else {
			dados.Erro = "Ação inválida."
		}

	} else if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Paginação
	pagina, _ := strconv.Atoi(r.URL.Query().Get("pagina"))

	if pagina < 1 {
		pagina = 1
	}

	produtos, total, err := services.ProdutosPaginados(
		"",
		pagina,
		produtosPorPagina,
	)

	if err != nil {
		log.Println("erro em ProdutosPaginados:", err)
		http.Error(
			w,
			"Erro ao buscar produtos",
			http.StatusInternalServerError,
		)
		return
	}

	dados.Produtos = produtos
	dados.Pagina = pagina

	dados.TotalPaginas = (total + produtosPorPagina - 1) / produtosPorPagina

	if dados.TotalPaginas < 1 {
		dados.TotalPaginas = 1
	}

	dados.PaginaAnterior = pagina - 1
	dados.PaginaProxima = pagina + 1

	tmpl, err := template.ParseFiles(
		"frontend/html/alterar_produto.html",
	)

	if err != nil {
		http.Error(
			w,
			"Erro ao carregar alteração de produto",
			http.StatusInternalServerError,
		)
		return
	}

	if dados.Erro != "" {
		w.WriteHeader(http.StatusBadRequest)
	}

	if err := tmpl.Execute(w, dados); err != nil {
		http.Error(
			w,
			"Erro ao renderizar alteração de produto",
			http.StatusInternalServerError,
		)
	}
}
