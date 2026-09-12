package services

import (
	"errors"
	"strings"

	database "gostock/backend/database"
	"gostock/backend/utils"

	"golang.org/x/crypto/bcrypt"
)

// ResetPassword troca a senha de uma conta ativa já existente, pelo email.
// Usado pelo comando de linha de comando "reset-password": ainda não existe
// uma tela no painel para trocar a senha de um usuário já criado (só na
// criação e no seed inicial), então isso cobre essa operação enquanto essa
// tela não existe.
func ResetPassword(email, password string) error {
	email = strings.ToLower(strings.TrimSpace(email))
	if !utils.ValidateEmail(email) {
		return errors.New("email inválido")
	}
	if !utils.ValidatePassword(password) {
		return errors.New("a senha deve ter no mínimo 6 caracteres e 1 caractere especial")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("erro ao proteger a senha")
	}

	result, err := database.DB.Exec(`
		UPDATE usuarios SET senha = ?
		WHERE LOWER(TRIM(email)) = ? AND ativo = 1
	`, string(hash), email)
	if err != nil {
		return errors.New("erro ao atualizar a senha")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return errors.New("erro ao atualizar a senha")
	}
	if rows == 0 {
		return errors.New("nenhum usuário ativo encontrado com esse email")
	}

	return nil
}

// RenameUser troca o nome de exibição de uma conta ativa já existente,
// pelo email. Usado pelo comando de linha de comando "rename-user" pelo
// mesmo motivo do ResetPassword: o seed só define o nome na criação da
// conta, nunca em contas que já existem.
func RenameUser(email, name string) error {
	email = strings.ToLower(strings.TrimSpace(email))
	name = strings.TrimSpace(name)

	if !utils.ValidateEmail(email) {
		return errors.New("email inválido")
	}
	if name == "" {
		return errors.New("informe o novo nome")
	}

	result, err := database.DB.Exec(`
		UPDATE usuarios SET nome = ?
		WHERE LOWER(TRIM(email)) = ? AND ativo = 1
	`, name, email)
	if err != nil {
		return errors.New("erro ao atualizar o nome")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return errors.New("erro ao atualizar o nome")
	}
	if rows == 0 {
		return errors.New("nenhum usuário ativo encontrado com esse email")
	}

	return nil
}

// ChangeUserEmail troca o email de uma conta ativa já existente. Usado
// pelo comando "change-email" pelo mesmo motivo do ResetPassword: o seed
// só define o email na criação da conta. O email ceo@gmail.com é
// reservado e não pode ser assumido por outra conta por aqui.
func ChangeUserEmail(currentEmail, newEmail string) error {
	currentEmail = strings.ToLower(strings.TrimSpace(currentEmail))
	newEmail = strings.ToLower(strings.TrimSpace(newEmail))

	if !utils.ValidateEmail(currentEmail) || !utils.ValidateEmail(newEmail) {
		return errors.New("email inválido")
	}
	if newEmail == "ceo@gmail.com" {
		return errors.New("o email ceo@gmail.com é reservado ao CEO")
	}

	result, err := database.DB.Exec(`
		UPDATE usuarios SET email = ?
		WHERE LOWER(TRIM(email)) = ? AND ativo = 1
		AND NOT EXISTS (SELECT 1 FROM usuarios WHERE LOWER(TRIM(email)) = ?)
	`, newEmail, currentEmail, newEmail)
	if err != nil {
		return errors.New("erro ao atualizar o email")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return errors.New("erro ao atualizar o email")
	}
	if rows == 0 {
		return errors.New("nenhum usuário ativo encontrado com esse email, ou o novo email já está em uso")
	}

	return nil
}
