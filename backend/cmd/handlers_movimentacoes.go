package main

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"gostock/backend/models"
	"gostock/backend/services"
)

// handlerMovimentacoes exibe o histórico de entradas/saídas/atualizações.
func handlerMovimentacoes(w http.ResponseWriter, r *http.Request) {
	if _, true := usuarioDaSessao(r); !true {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	movimentacoes, err := services.BuscarMovimentacoesWeb()
	if err != nil {
		http.Error(w, "Erro ao buscar movimentações", http.StatusInternalServerError)
		return
	}

	const movimentacoesPorPagina = 5
	pagina, _ := strconv.Atoi(r.URL.Query().Get("pagina"))
	if pagina < 1 {
		pagina = 1
	}
	totalPaginas := (len(movimentacoes) + movimentacoesPorPagina - 1) / movimentacoesPorPagina
	if totalPaginas < 1 {
		totalPaginas = 1
	}
	if pagina > totalPaginas {
		pagina = totalPaginas
	}
	inicio := (pagina - 1) * movimentacoesPorPagina
	fim := inicio + movimentacoesPorPagina
	if fim > len(movimentacoes) {
		fim = len(movimentacoes)
	}
	movimentacoes = movimentacoes[inicio:fim]

	tmpl, err := template.ParseFiles("frontend/html/movimentacoes.html")
	if err != nil {
		http.Error(w, "Erro ao carregar movimentações", http.StatusInternalServerError)
		return
	}

	dados := struct {
		Movimentacoes  []models.Movimentacao
		Pagina         int
		TotalPaginas   int
		PaginaAnterior int
		PaginaProxima  int
		EhAdmin        bool
	}{
		Movimentacoes:  movimentacoes,
		Pagina:         pagina,
		TotalPaginas:   totalPaginas,
		PaginaAnterior: pagina - 1,
		PaginaProxima:  pagina + 1,
		EhAdmin:        usuarioEhAdmin(r),
	}

	var conteudo bytes.Buffer
	if err := tmpl.Execute(&conteudo, dados); err != nil {
		log.Println("erro ao renderizar movimentações:", err)
		http.Error(w, "Erro ao renderizar movimentações", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(conteudo.Bytes())
}
