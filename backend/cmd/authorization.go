package main

import (
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
	return err == nil && strings.EqualFold(strings.TrimSpace(usuario.Role), "admin")
}

func exigirAdmin(w http.ResponseWriter, r *http.Request) bool {
	if usuarioEhAdmin(r) {
		return true
	}
	http.Error(w, "Acesso permitido apenas para administradores", http.StatusForbidden)
	return false
}
