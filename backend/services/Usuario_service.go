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

func AutenticarUsuario(email string, senha string) (*models.Usuario, bool) {
	email = strings.ToLower(strings.TrimSpace(email))

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
		return nil, false
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(senhaHash),
		[]byte(senha),
	)

	if err != nil {
		return nil, false
	}

	return &usuario, true
}

// Login continua sendo usado pelo sistema do terminal.
func Login(reader *bufio.Reader) (*models.Usuario, bool) {
	email := strings.ToLower(
		strings.TrimSpace(
			utils.LerTexto(reader, "Email: "),
		),
	)

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
func BuscarUsuarioPorID(id int) (*models.Usuario, error) {
	var usuario models.Usuario

	query := `
		SELECT id, nome, email
		FROM usuarios
		WHERE id = ?
	`

	err := database.DB.QueryRow(
		query,
		id,
	).Scan(
		&usuario.ID,
		&usuario.Nome,
		&usuario.Email,
	)

	if err != nil {
		return nil, err
	}

	return &usuario, nil
}
