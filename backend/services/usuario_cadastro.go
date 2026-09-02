package services

import (
	"bufio"
	"fmt"
	"strings"

	database "gostock/backend/Database"
	"gostock/backend/utils"

	"golang.org/x/crypto/bcrypt"
)

func CadastrarUsuario(reader *bufio.Reader) {
	nome := utils.LerTexto(reader, "Nome: ")

	var email string

	for {
		email = strings.ToLower(strings.TrimSpace(utils.LerTexto(reader, "Email: ")))

		if utils.ValidarEmail(email) {
			break
		}

		fmt.Println("Email inválido. Use um endereço @gmail.com.")
	}

	var senha string

	for {
		senha = utils.LerTexto(reader, "Senha: ")

		if utils.ValidarSenha(senha) {
			break
		}

		fmt.Println("A senha deve ter no mínimo 6 caracteres e 1 caractere especial.")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Erro ao proteger senha:", err)
		return
	}

	_, err = database.DB.Exec(`
		INSERT INTO usuarios (nome, email, senha)
		VALUES (?, ?, ?)
	`, nome, email, string(hash))

	if err != nil {
		fmt.Println("Esse email já está cadastrado.")
		return
	}

	fmt.Println("Usuário cadastrado com sucesso.")
}

func CadastrarUsuarioWeb(nome, email, senha string) error {
	nome = strings.TrimSpace(nome)
	email = strings.ToLower(strings.TrimSpace(email))
	if !utils.ValidarNome(nome) {
		return fmt.Errorf("o nome é obrigatório")
	}
	if !utils.ValidarEmail(email) {
		return fmt.Errorf("use um endereço @gmail.com válido")
	}
	if !utils.ValidarSenha(senha) {
		return fmt.Errorf("a senha deve ter no mínimo 6 caracteres e 1 caractere especial")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = database.DB.Exec(`
		INSERT INTO usuarios (nome, email, senha, role)
		VALUES (?, ?, ?, 'basico')
	`, nome, email, string(hash))
	if err != nil {
		return fmt.Errorf("não foi possível criar a conta: email já cadastrado")
	}
	return nil
}
