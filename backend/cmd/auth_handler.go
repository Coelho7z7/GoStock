package main

import (
	"crypto/sha256"
	"encoding/hex"
	"html/template"
	"net/http"

	database "gostock/backend/database"
	"gostock/backend/services"
)

// indexHandler exibe a tela de login (rota "/").
func indexHandler(w http.ResponseWriter, r *http.Request) {
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
		Error string
	}{}); err != nil {
		http.Error(w, "Erro ao renderizar tela de login", http.StatusInternalServerError)
	}
}

// loginHandler processa o formulário de login.
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("senha")

	user, success := services.AuthenticateUser(email, password)
	if !success {
		tmpl, err := template.ParseFiles("frontend/html/index.html")
		if err != nil {
			http.Error(w, "Erro ao carregar tela de login", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusUnauthorized)
		if err := tmpl.Execute(w, struct {
			Email string
			Error string
		}{email, "Email ou senha incorretos."}); err != nil {
			http.Error(w, "Erro ao renderizar tela de login", http.StatusInternalServerError)
		}
		return
	}

	token, err := createSession(user.ID)
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

// logoutHandler apaga a sessão atual e redireciona para o login.
func logoutHandler(w http.ResponseWriter, r *http.Request) {
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
