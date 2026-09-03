package services

import (
	"fmt"

	database "gostock/backend/Database"

	"golang.org/x/crypto/bcrypt"
)

func SeedUsuariosPadrao() error {
	usuarios := []struct {
		nome  string
		email string
		senha string
		role  string
	}{
		{nome: "Administrador", email: "admin@gmail.com", senha: "@admin12e", role: "ceo"},
		{nome: "Usuario", email: "usuario@gmail.com", senha: "usuario123", role: "basico"},
	}

	for _, usuario := range usuarios {
		var existe bool
		if err := database.DB.QueryRow(`
			SELECT EXISTS(SELECT 1 FROM usuarios WHERE LOWER(TRIM(email)) = LOWER(TRIM(?)))
		`, usuario.email).Scan(&existe); err != nil {
			return err
		}

		if usuario.email == "admin@gmail.com" {
			// A conta CEO é criada apenas se ainda não existir e, se existir,
			// recebe somente a correção de identidade/cargo; a senha existente
			// não é sobrescrita em cada inicialização.
			if !existe {
				hash, err := bcrypt.GenerateFromPassword([]byte(usuario.senha), bcrypt.DefaultCost)
				if err != nil {
					return fmt.Errorf("gerar senha de %s: %w", usuario.email, err)
				}
				if _, err := database.DB.Exec(`
					INSERT INTO usuarios (nome, email, senha, role)
					VALUES (?, ?, ?, 'ceo')
				`, usuario.nome, usuario.email, string(hash)); err != nil {
					return fmt.Errorf("inserir usuário %s: %w", usuario.email, err)
				}
			} else {
				if _, err := database.DB.Exec(`
					UPDATE usuarios
					SET role = 'ceo'
					WHERE LOWER(TRIM(email)) = 'admin@gmail.com'
				`); err != nil {
					return fmt.Errorf("proteger usuário %s: %w", usuario.email, err)
				}
			}
			continue
		}

		if existe {
			continue
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(usuario.senha), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("gerar senha de %s: %w", usuario.email, err)
		}

		if _, err := database.DB.Exec(`
			INSERT INTO usuarios (nome, email, senha, role)
			VALUES (?, ?, ?, ?)
		`, usuario.nome, usuario.email, string(hash), usuario.role); err != nil {
			return fmt.Errorf("inserir usuário %s: %w", usuario.email, err)
		}
	}

	return nil
}
