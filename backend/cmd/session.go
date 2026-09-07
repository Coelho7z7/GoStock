package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	database "gostock/backend/database"
)

// createSession gera um token aleatório, guarda apenas o hash dele no
// banco (nunca o token em texto puro) e devolve o token para ser
// colocado no cookie do navegador.
func createSession(userID int) (string, error) {
	bytes := make([]byte, 32)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	token := hex.EncodeToString(bytes)

	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	expiresAt := time.Now().Add(30 * 24 * time.Hour)

	_, err := database.DB.Exec(`
		INSERT INTO sessoes (usuario_id, token_hash, expira_em)
		VALUES (?, ?, ?)
	`, userID, tokenHash, expiresAt)

	if err != nil {
		return "", err
	}

	return token, nil
}

// userFromSession lê o cookie de sessão da requisição e retorna o ID
// do usuário logado, se a sessão existir e ainda não tiver expirado.
func userFromSession(r *http.Request) (int, bool) {
	cookie, err := r.Cookie("sessao")

	if err != nil {
		return 0, false
	}

	hash := sha256.Sum256([]byte(cookie.Value))
	tokenHash := hex.EncodeToString(hash[:])

	var userID int
	var role string
	var expiresAt time.Time

	err = database.DB.QueryRow(`
		SELECT s.usuario_id, s.expira_em, u.role
		FROM sessoes s
		JOIN usuarios u ON u.id = s.usuario_id
		WHERE token_hash = ?
	`, tokenHash).Scan(&userID, &expiresAt, &role)

	if err != nil {
		return 0, false
	}

	if time.Now().After(expiresAt) {
		database.DB.Exec(`
			DELETE FROM sessoes
			WHERE token_hash = ?
		`, tokenHash)

		return 0, false
	}

	return userID, true
}
