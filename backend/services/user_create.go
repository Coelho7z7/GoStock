package services

import (
	"bufio"
	"fmt"
	"strings"

	database "gostock/backend/database"
	"gostock/backend/utils"

	"golang.org/x/crypto/bcrypt"
)

func CreateUser(reader *bufio.Reader) {
	name := utils.ReadText(reader, "Nome: ")

	var email string

	for {
		email = strings.ToLower(strings.TrimSpace(utils.ReadText(reader, "Email: ")))

		if utils.ValidateEmail(email) {
			break
		}

		fmt.Println("Email inválido. Use um endereço @gmail.com.")
	}

	var password string

	for {
		password = utils.ReadText(reader, "Senha: ")

		if utils.ValidatePassword(password) {
			break
		}

		fmt.Println("A senha deve ter no mínimo 6 caracteres e 1 caractere especial.")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Erro ao proteger senha:", err)
		return
	}

	_, err = database.DB.Exec(`
		INSERT INTO usuarios (nome, email, senha)
		VALUES (?, ?, ?)
	`, name, email, string(hash))

	if err != nil {
		fmt.Println("Esse email já está cadastrado.")
		return
	}

	fmt.Println("Usuário cadastrado com sucesso.")
}
