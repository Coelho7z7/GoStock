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

// handlerUsuarios exibe a lista de usuários cadastrados e processa a
// criação de novos usuários, a alteração de permissão e a remoção.
// Toda a tela é restrita a administradores — inclusive o GET.
func handlerUsuarios(w http.ResponseWriter, r *http.Request) {
	usuarioID, autenticado := usuarioDaSessao(r)
	if !autenticado {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if !exigirAdmin(w, r) {
		return
	}

	const usuariosPorPagina = 8

	dados := struct {
		Usuarios       []models.Usuario
		UsuarioID      int
		Busca          string
		Nome           string
		Email          string
		Role           string
		Mensagem       string
		Erro           string
		Pagina         int
		TotalPaginas   int
		PaginaAnterior int
		PaginaProxima  int
	}{UsuarioID: usuarioID}

	dados.Mensagem = map[string]string{
		"criado":     "Usuário criado com sucesso.",
		"atualizado": "Permissão atualizada com sucesso.",
		"removido":   "Usuário removido com sucesso.",
	}[r.URL.Query().Get("sucesso")]

	if r.Method == http.MethodPost {
		switch r.FormValue("acao") {

		case "criar":
			dados.Nome = strings.TrimSpace(r.FormValue("nome"))
			dados.Email = strings.TrimSpace(r.FormValue("email"))
			dados.Role = r.FormValue("role")
			senha := r.FormValue("senha")

			if err := services.CriarUsuarioWeb(dados.Nome, dados.Email, senha, dados.Role); err != nil {
				dados.Erro = err.Error()
			} else {
				http.Redirect(w, r, "/usuarios?sucesso=criado", http.StatusSeeOther)
				return
			}

		case "alterar_permissao":
			alvoID, idErr := strconv.Atoi(r.FormValue("usuario_id"))
			novaRole := r.FormValue("role")

			if idErr != nil {
				dados.Erro = "Usuário inválido."
			} else if alvoID == usuarioID {
				dados.Erro = "Você não pode alterar a sua própria permissão."
			} else if err := services.AtualizarRoleUsuarioWeb(alvoID, novaRole); err != nil {
				dados.Erro = err.Error()
			} else {
				http.Redirect(w, r, "/usuarios?sucesso=atualizado", http.StatusSeeOther)
				return
			}

		case "remover":
			alvoID, idErr := strconv.Atoi(r.FormValue("usuario_id"))

			if idErr != nil {
				dados.Erro = "Usuário inválido."
			} else if alvoID == usuarioID {
				dados.Erro = "Você não pode remover a sua própria conta."
			} else if err := services.RemoverUsuarioWeb(alvoID); err != nil {
				dados.Erro = err.Error()
			} else {
				http.Redirect(w, r, "/usuarios?sucesso=removido", http.StatusSeeOther)
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

	pagina, _ := strconv.Atoi(r.URL.Query().Get("pagina"))
	if pagina < 1 {
		pagina = 1
	}
	dados.Busca = strings.TrimSpace(r.URL.Query().Get("busca"))

	usuarios, total, err := services.ListarUsuariosPaginados(dados.Busca, pagina, usuariosPorPagina)
	if err != nil {
		log.Println("erro em ListarUsuariosPaginados:", err)
		http.Error(w, "Erro ao buscar usuários", http.StatusInternalServerError)
		return
	}

	dados.Usuarios = usuarios
	dados.Pagina = pagina
	dados.TotalPaginas = (total + usuariosPorPagina - 1) / usuariosPorPagina
	if dados.TotalPaginas < 1 {
		dados.TotalPaginas = 1
	}
	dados.PaginaAnterior = pagina - 1
	dados.PaginaProxima = pagina + 1

	tmpl, err := template.ParseFiles("frontend/html/usuarios.html")
	if err != nil {
		http.Error(w, "Erro ao carregar usuários", http.StatusInternalServerError)
		return
	}

	if dados.Erro != "" {
		w.WriteHeader(http.StatusBadRequest)
	}

	if err := tmpl.Execute(w, dados); err != nil {
		http.Error(w, "Erro ao renderizar usuários", http.StatusInternalServerError)
	}
}
