package services

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"

	database "gostock/backend/database"

	"golang.org/x/crypto/bcrypt"
)

// seedPassword lê a senha padrão de uma variável de ambiente; se ela não
// estiver definida, gera uma senha aleatória e imprime uma única vez no log
// (nunca fica hardcoded no código-fonte nem versionada no git).
func seedPassword(envVar, label string) (string, error) {
	if pw := os.Getenv(envVar); pw != "" {
		return pw, nil
	}

	buf := make([]byte, 12)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("gerar senha aleatória para %s: %w", label, err)
	}
	pw := base64.RawURLEncoding.EncodeToString(buf)
	fmt.Printf("[seed] %s: variável %s não definida, gerando senha temporária: %s\n", label, envVar, pw)
	return pw, nil
}

func SeedDefaultUsers() error {
	gerentePw, err := seedPassword("SEED_GERENTE_PASSWORD", "Matheus (gerente)")
	if err != nil {
		return err
	}
	ceoPw, err := seedPassword("SEED_CEO_PASSWORD", "Administrador (ceo)")
	if err != nil {
		return err
	}
	usuarioPw, err := seedPassword("SEED_USUARIO_PASSWORD", "Usuario (basico)")
	if err != nil {
		return err
	}

	users := []struct {
		name     string
		email    string
		password string
		role     string
	}{
		{name: "Gerente", email: "gerente@gmail.com", password: gerentePw, role: "gerente"},
		{name: "Administrador", email: "ceo@gmail.com", password: ceoPw, role: "ceo"},
		{name: "Usuario", email: "usuario@gmail.com", password: usuarioPw, role: "basico"},
	}

	for _, user := range users {
		var exists bool
		if err := database.DB.QueryRow(`
			SELECT EXISTS(SELECT 1 FROM usuarios WHERE LOWER(TRIM(email)) = LOWER(TRIM(?)))
		`, user.email).Scan(&exists); err != nil {
			return err
		}

		if user.email == "ceo@gmail.com" {
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
					WHERE LOWER(TRIM(email)) = 'ceo@gmail.com'
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
