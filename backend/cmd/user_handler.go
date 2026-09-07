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

// userHandler exibe a lista de usuários cadastrados e processa a
// criação de novos usuários, a alteração de permissão e a remoção.
// A tela é acessível a administradores e gerentes; remover e alterar
// permissão continuam restritos a administradores (checado abaixo).
func userHandler(w http.ResponseWriter, r *http.Request) {
	userID, authenticated := userFromSession(r)
	if !authenticated {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if !requireAdminOrManager(w, r) {
		return
	}

	admin := isAdmin(r)

	const usersPerPage = 8

	data := struct {
		Users        []models.User
		UserID       int
		IsAdmin      bool
		Search       string
		Name         string
		Email        string
		Role         string
		Message      string
		Error        string
		Page         int
		TotalPages   int
		PreviousPage int
		NextPage     int
	}{UserID: userID, IsAdmin: admin}

	data.Message = map[string]string{
		"criado":     "Usuário criado com sucesso.",
		"atualizado": "Permissão atualizada com sucesso.",
		"removido":   "Usuário removido com sucesso.",
	}[r.URL.Query().Get("sucesso")]

	if r.Method == http.MethodPost {
		switch r.FormValue("acao") {

		case "criar":
			data.Name = strings.TrimSpace(r.FormValue("nome"))
			data.Email = strings.TrimSpace(r.FormValue("email"))
			data.Role = r.FormValue("role")
			password := r.FormValue("senha")

			// Gerente só pode cadastrar usuários com permissão básica —
			// não pode criar outro admin/gerente por aqui.
			if !admin {
				data.Role = "basico"
			}

			if err := services.CreateUserWeb(data.Name, data.Email, password, data.Role); err != nil {
				data.Error = err.Error()
			} else {
				http.Redirect(w, r, "/usuarios?sucesso=criado", http.StatusSeeOther)
				return
			}

		case "alterar_permissao":
			if !admin {
				data.Error = "Apenas administradores podem alterar permissões."
				break
			}

			targetID, idErr := strconv.Atoi(r.FormValue("usuario_id"))
			newRole := r.FormValue("role")

			if idErr != nil {
				data.Error = "Usuário inválido."
			} else if targetID == userID {
				data.Error = "Você não pode alterar a sua própria permissão."
			} else if err := services.UpdateUserRoleWeb(targetID, newRole); err != nil {
				data.Error = err.Error()
			} else {
				http.Redirect(w, r, "/usuarios?sucesso=atualizado", http.StatusSeeOther)
				return
			}

		case "remover":
			if !admin {
				data.Error = "Apenas administradores podem remover usuários."
				break
			}

			targetID, idErr := strconv.Atoi(r.FormValue("usuario_id"))

			if idErr != nil {
				data.Error = "Usuário inválido."
			} else if targetID == userID {
				data.Error = "Você não pode remover a sua própria conta."
			} else if err := services.DeleteUserWeb(targetID); err != nil {
				data.Error = err.Error()
			} else {
				http.Redirect(w, r, "/usuarios?sucesso=removido", http.StatusSeeOther)
				return
			}

		default:
			data.Error = "Ação inválida."
		}
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

	users, total, err := services.ListPaginatedUsers(data.Search, page, usersPerPage)
	if err != nil {
		log.Println("erro em ListPaginatedUsers:", err)
		http.Error(w, "Erro ao buscar usuários", http.StatusInternalServerError)
		return
	}

	data.Users = users
	data.Page = page
	data.TotalPages = (total + usersPerPage - 1) / usersPerPage
	if data.TotalPages < 1 {
		data.TotalPages = 1
	}
	data.PreviousPage = page - 1
	data.NextPage = page + 1

	tmpl, err := template.ParseFiles("frontend/html/users.html")
	if err != nil {
		http.Error(w, "Erro ao carregar usuários", http.StatusInternalServerError)
		return
	}

	if data.Error != "" {
		w.WriteHeader(http.StatusBadRequest)
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Erro ao renderizar usuários", http.StatusInternalServerError)
	}
}
