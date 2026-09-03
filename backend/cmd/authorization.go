package main

import (
	"html/template"
	"net/http"
	"strings"

	"gostock/backend/services"
)

func usuarioEhAdmin(r *http.Request) bool {
	usuarioID, autenticado := usuarioDaSessao(r)
	if !autenticado {
		return false
	}

	usuario, err := services.BuscarUsuarioPorID(usuarioID)
	if err != nil {
		return false
	}

	role := strings.ToLower(strings.TrimSpace(usuario.Role))
	return role == "admin" || role == "ceo"
}

func exigirAdmin(w http.ResponseWriter, r *http.Request) bool {
	if usuarioEhAdmin(r) {
		return true
	}
	tmpl, err := template.ParseFiles("frontend/html/acesso_negado.html")
	if err != nil {
		http.Error(w, "Acesso negado", http.StatusForbidden)
		return false
	}
	w.WriteHeader(http.StatusForbidden)
	if err := tmpl.Execute(w, nil); err != nil {
		return false
	}
	return false
}
