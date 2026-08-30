package main

import (
	"html/template"
	"net/http"

	"gostock/backend/models"
	"gostock/backend/services"
)

// handlerMovimentacoes exibe o histórico de entradas/saídas/atualizações.
func handlerMovimentacoes(w http.ResponseWriter, r *http.Request) {
	if _, ok := usuarioDaSessao(r); !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	movimentacoes, err := services.BuscarMovimentacoesWeb()
	if err != nil {
		http.Error(w, "Erro ao buscar movimentações", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("frontend/html/movimentacoes.html")
	if err != nil {
		http.Error(w, "Erro ao carregar movimentações", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, struct {
		Movimentacoes []models.Movimentacao
	}{movimentacoes}); err != nil {
		http.Error(w, "Erro ao renderizar movimentações", http.StatusInternalServerError)
	}
}
