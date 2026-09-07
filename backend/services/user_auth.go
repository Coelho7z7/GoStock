package services

import (
	"bufio"
	"fmt"
	"strings"

	database "gostock/backend/database"
	"gostock/backend/models"
	"gostock/backend/utils"

	"golang.org/x/crypto/bcrypt"
)

func AuthenticateUser(email string, password string) (*models.User, bool) {
	email = strings.ToLower(strings.TrimSpace(email))

	query := `
		SELECT id, nome, senha, role
		FROM usuarios
		WHERE email = ? AND ativo = 1
	`

	var user models.User
	var passwordHash string

	err := database.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&passwordHash,
		&user.Role,
	)

	if err != nil {
		return nil, false
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return nil, false
	}

	return &user, true
}

// Login continua sendo usado pelo sistema do terminal.
func Login(reader *bufio.Reader) (*models.User, bool) {
	email := strings.ToLower(strings.TrimSpace(utils.ReadText(reader, "Email: ")))
	password := utils.ReadText(reader, "Senha: ")

	user, success := AuthenticateUser(email, password)
	if !success {
		fmt.Println("Email ou senha incorretos.")
		return nil, false
	}

	fmt.Println("Login realizado com sucesso!")
	fmt.Println("Bem-vindo,", user.Name)

	return user, true
}
