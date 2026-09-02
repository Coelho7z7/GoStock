package main

import (
	"crypto/sha256"
	"encoding/hex"
	"html/template"
	"net/http"

	database "gostock/backend/Database"
	"gostock/backend/services"
)

// handlerIndex exibe a tela de login (rota "/").
func handlerIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles("frontend/html/index.html")
	if err != nil {
		http.Error(w, "Erro ao carregar tela de login", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, struct {
		Email string
		Erro  string
	}{}); err != nil {
		http.Error(w, "Erro ao renderizar tela de login", http.StatusInternalServerError)
	}
}

// handlerLogin processa o formulário de login.
func handlerLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")
	senha := r.FormValue("senha")

	usuario, sucesso := services.AutenticarUsuario(email, senha)
	if !sucesso {
		tmpl, err := template.ParseFiles("frontend/html/index.html")
		if err != nil {
			http.Error(w, "Erro ao carregar tela de login", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusUnauthorized)
		if err := tmpl.Execute(w, struct {
			Email string
			Erro  string
		}{email, "Email ou senha incorretos."}); err != nil {
			http.Error(w, "Erro ao renderizar tela de login", http.StatusInternalServerError)
		}
		return
	}

	token, err := criarSessao(usuario.ID)
	if err != nil {
		http.Error(w, "Erro ao iniciar sessão", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "sessao",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   60 * 60 * 24 * 30,
	})

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func handlerCadastro(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("frontend/html/cadastro.html")
		if err != nil {
			http.Error(w, "Erro ao carregar cadastro", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
		return
	}
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	if err := services.CadastrarUsuarioWeb(r.FormValue("nome"), r.FormValue("email"), r.FormValue("senha")); err != nil {
		tmpl, loadErr := template.ParseFiles("frontend/html/cadastro.html")
		if loadErr != nil {
			http.Error(w, "Erro ao carregar cadastro", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		tmpl.Execute(w, struct{ Erro string }{Erro: err.Error()})
		return
	}
	http.Redirect(w, r, "/?cadastro=sucesso", http.StatusSeeOther)
}

// handlerLogout apaga a sessão atual e redireciona para o login.
func handlerLogout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("sessao"); err == nil {
		hash := sha256.Sum256([]byte(cookie.Value))
		tokenHash := hex.EncodeToString(hash[:])
		database.DB.Exec("DELETE FROM sessoes WHERE token_hash = ?", tokenHash)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "sessao",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
