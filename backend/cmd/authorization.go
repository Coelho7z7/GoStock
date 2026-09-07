package main

import (
	"html/template"
	"net/http"
	"strings"

	"gostock/backend/services"
)

// isAdmin indica se o usuário logado tem poderes de administrador.
// O CEO (reservado a admin@gmail.com) também conta como admin aqui — ele
// fica acima do administrador na hierarquia, então tudo que um admin pode
// fazer o CEO também pode.
func isAdmin(r *http.Request) bool {
	userID, authenticated := userFromSession(r)
	if !authenticated {
		return false
	}

	user, err := services.GetUserByID(userID)
	if err != nil {
		return false
	}
	role := strings.ToLower(strings.TrimSpace(user.Role))
	return role == "admin" || role == "ceo"
}

// isManager indica se o usuário logado tem a permissão "gerente".
// Um gerente pode cadastrar novos usuários (com permissão básica), mas não
// pode remover usuários nem alterar permissões — isso continua exclusivo
// do administrador.
func isManager(r *http.Request) bool {
	userID, authenticated := userFromSession(r)
	if !authenticated {
		return false
	}

	user, err := services.GetUserByID(userID)
	return err == nil && strings.EqualFold(strings.TrimSpace(user.Role), "gerente")
}

// canViewUsersTab indica se o usuário logado pode acessar a aba de
// administração de usuários (administradores e gerentes).
func canViewUsersTab(r *http.Request) bool {
	return isAdmin(r) || isManager(r)
}

func requireAdmin(w http.ResponseWriter, r *http.Request) bool {
	if isAdmin(r) {
		return true
	}
	tmpl, err := template.ParseFiles("frontend/html/access_denied.html")
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

// requireAdminOrManager bloqueia o acesso de quem não é administrador nem
// gerente — usado na tela de usuários, que gerentes também podem abrir.
func requireAdminOrManager(w http.ResponseWriter, r *http.Request) bool {
	if canViewUsersTab(r) {
		return true
	}
	tmpl, err := template.ParseFiles("frontend/html/access_denied.html")
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
