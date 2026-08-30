package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	database "gostock/backend/Database"
)

// criarSessao gera um token aleatório, guarda apenas o hash dele no
// banco (nunca o token em texto puro) e devolve o token para ser
// colocado no cookie do navegador.
func criarSessao(usuarioID int) (string, error) {
	bytes := make([]byte, 32)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	token := hex.EncodeToString(bytes)

	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	expiraEm := time.Now().Add(30 * 24 * time.Hour)

	_, err := database.DB.Exec(`
		INSERT INTO sessoes (usuario_id, token_hash, expira_em)
		VALUES (?, ?, ?)
	`, usuarioID, tokenHash, expiraEm)

	if err != nil {
		return "", err
	}

	return token, nil
}

// usuarioDaSessao lê o cookie de sessão da requisição e retorna o ID
// do usuário logado, se a sessão existir e ainda não tiver expirado.
func usuarioDaSessao(r *http.Request) (int, bool) {
	cookie, err := r.Cookie("sessao")

	if err != nil {
		return 0, false
	}

	hash := sha256.Sum256([]byte(cookie.Value))
	tokenHash := hex.EncodeToString(hash[:])

	var usuarioID int
	var expiraEm time.Time

	err = database.DB.QueryRow(`
		SELECT usuario_id, expira_em
		FROM sessoes
		WHERE token_hash = ?
	`, tokenHash).Scan(&usuarioID, &expiraEm)

	if err != nil {
		return 0, false
	}

	if time.Now().After(expiraEm) {
		database.DB.Exec(`
			DELETE FROM sessoes
			WHERE token_hash = ?
		`, tokenHash)

		return 0, false
	}

	return usuarioID, true
}
