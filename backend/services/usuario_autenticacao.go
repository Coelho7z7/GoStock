package services

import (
	"bufio"
	"fmt"
	"strings"

	database "gostock/backend/Database"
	"gostock/backend/models"
	"gostock/backend/utils"

	"golang.org/x/crypto/bcrypt"
)

func AutenticarUsuario(email string, senha string) (*models.Usuario, bool) {
	email = strings.ToLower(strings.TrimSpace(email))

	query := `
		SELECT id, nome, senha
		FROM usuarios
		WHERE email = ?
	`

	var usuario models.Usuario
	var senhaHash string

	err := database.DB.QueryRow(query, email).Scan(
		&usuario.ID,
		&usuario.Nome,
		&senhaHash,
	)

	if err != nil {
		return nil, false
	}

	if err := bcrypt.CompareHashAndPassword([]byte(senhaHash), []byte(senha)); err != nil {
		return nil, false
	}

	return &usuario, true
}

// Login continua sendo usado pelo sistema do terminal.
func Login(reader *bufio.Reader) (*models.Usuario, bool) {
	email := strings.ToLower(strings.TrimSpace(utils.LerTexto(reader, "Email: ")))
	senha := utils.LerTexto(reader, "Senha: ")

	usuario, sucesso := AutenticarUsuario(email, senha)
	if !sucesso {
		fmt.Println("Email ou senha incorretos.")
		return nil, false
	}

	fmt.Println("Login realizado com sucesso!")
	fmt.Println("Bem-vindo,", usuario.Nome)

	return usuario, true
}
