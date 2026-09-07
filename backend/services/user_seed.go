package services

import (
	"fmt"

	database "gostock/backend/database"

	"golang.org/x/crypto/bcrypt"
)

func SeedDefaultUsers() error {
	users := []struct {
		name     string
		email    string
		password string
		role     string
	}{
		{name: "Matheus", email: "matheus@gmail.com", password: "Dominio12e@", role: "gerente"},
		{name: "Administrador", email: "admin@gmail.com", password: "@admin12e", role: "ceo"},
		{name: "Usuario", email: "usuario@gmail.com", password: "usuario123", role: "basico"},
	}

	for _, user := range users {
		var exists bool
		if err := database.DB.QueryRow(`
			SELECT EXISTS(SELECT 1 FROM usuarios WHERE LOWER(TRIM(email)) = LOWER(TRIM(?)))
		`, user.email).Scan(&exists); err != nil {
			return err
		}

		if user.email == "admin@gmail.com" {
			// A conta CEO é criada apenas se ainda não existir e, se existir,
			// recebe somente a correção de identidade/cargo; a senha existente
			// não é sobrescrita em cada inicialização.
			if !exists {
				hash, err := bcrypt.GenerateFromPassword([]byte(user.password), bcrypt.DefaultCost)
				if err != nil {
					return fmt.Errorf("gerar senha de %s: %w", user.email, err)
				}
				if _, err := database.DB.Exec(`
					INSERT INTO usuarios (nome, email, senha, role)
					VALUES (?, ?, ?, 'ceo')
				`, user.name, user.email, string(hash)); err != nil {
					return fmt.Errorf("inserir usuário %s: %w", user.email, err)
				}
			} else {
				if _, err := database.DB.Exec(`
					UPDATE usuarios
					SET role = 'ceo'
					WHERE LOWER(TRIM(email)) = 'admin@gmail.com'
				`); err != nil {
					return fmt.Errorf("proteger usuário %s: %w", user.email, err)
				}
			}
			continue
		}

		if exists {
			continue
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(user.password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("gerar senha de %s: %w", user.email, err)
		}

		if _, err := database.DB.Exec(`
			INSERT INTO usuarios (nome, email, senha, role)
			VALUES (?, ?, ?, ?)
		`, user.name, user.email, string(hash), user.role); err != nil {
			return fmt.Errorf("inserir usuário %s: %w", user.email, err)
		}
	}

	return nil
}
