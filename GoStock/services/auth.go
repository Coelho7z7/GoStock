package services

import (
	"bufio"
	"fmt"
	"gostock/Database"
	"gostock/models"
	"strings"

	"gostock/utils"

	"golang.org/x/crypto/bcrypt"
)

func CadastrarUsuario(reader *bufio.Reader) {
	nome := utils.LerTexto(reader, "Nome: ")

	var email string

	for {
		email = strings.ToLower(
			strings.TrimSpace(
				utils.LerTexto(reader, "Email: "),
			),
		)

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

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(senha),
		bcrypt.DefaultCost,
	)

	if err != nil {
		fmt.Println("Erro ao proteger senha:", err)
		return
	}

	query := `
		INSERT INTO usuarios (nome, email, senha)
		VALUES (?, ?, ?)
	`

	_, err = database.DB.Exec(
		query,
		nome,
		email,
		string(hash),
	)

	if err != nil {
		fmt.Println("Esse email já está cadastrado.")
		return
	}

	fmt.Println("Usuário cadastrado com sucesso.")
}

func ExistemUsuarios() bool {
	var quantidade int

	query := `SELECT COUNT(*) FROM usuarios`

	err := database.DB.QueryRow(query).Scan(&quantidade)

	if err != nil {
		fmt.Println("Erro ao verificar usuários:", err)
		return false
	}

	return quantidade > 0
}

func Login(reader *bufio.Reader) (*models.Usuario, bool) {
	email := strings.ToLower(
		strings.TrimSpace(
			utils.LerTexto(reader, "Email: "),
		),
	)

	senha := utils.LerTexto(reader, "Senha: ")

	query := `
		SELECT id, nome, senha
		FROM usuarios
		WHERE email = ?
	`

	var usuario models.Usuario
	var senhaHash string

	err := database.DB.QueryRow(
		query,
		email,
	).Scan(
		&usuario.ID,
		&usuario.Nome,
		&senhaHash,
	)

	if err != nil {
		fmt.Println("Email ou senha incorretos.")
		return nil, false
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(senhaHash),
		[]byte(senha),
	)

	if err != nil {
		fmt.Println("Email ou senha incorretos.")
		return nil, false
	}

	fmt.Println("Login realizado com sucesso!")
	fmt.Println("Bem-vindo,", usuario.Nome)

	return &usuario, true
}
