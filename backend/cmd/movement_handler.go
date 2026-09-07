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

// movementHandler exibe o histórico de entradas/saídas/atualizações.
func movementHandler(w http.ResponseWriter, r *http.Request) {
	if _, authenticated := userFromSession(r); !authenticated {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	movements, err := services.GetMovementsWeb()
	if err != nil {
		http.Error(w, "Erro ao buscar movimentações", http.StatusInternalServerError)
		return
	}

	const movementsPerPage = 5
	page, _ := strconv.Atoi(r.URL.Query().Get("pagina"))
	if page < 1 {
		page = 1
	}
	totalPages := (len(movements) + movementsPerPage - 1) / movementsPerPage
	if totalPages < 1 {
		totalPages = 1
	}
	if page > totalPages {
		page = totalPages
	}
	start := (page - 1) * movementsPerPage
	end := start + movementsPerPage
	if end > len(movements) {
		end = len(movements)
	}
	movements = movements[start:end]

	tmpl, err := template.ParseFiles("frontend/html/movements.html")
	if err != nil {
		http.Error(w, "Erro ao carregar movimentações", http.StatusInternalServerError)
		return
	}

	data := struct {
		Movements    []models.Movement
		Page         int
		TotalPages   int
		PreviousPage int
		NextPage     int
		IsAdmin      bool
	}{
		Movements:    movements,
		Page:         page,
		TotalPages:   totalPages,
		PreviousPage: page - 1,
		NextPage:     page + 1,
		IsAdmin:      canViewUsersTab(r),
	}

	var content bytes.Buffer
	if err := tmpl.Execute(&content, data); err != nil {
		log.Println("erro ao renderizar movimentações:", err)
		http.Error(w, "Erro ao renderizar movimentações", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(content.Bytes())
}
