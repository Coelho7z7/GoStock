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
		{nome: "Administrador", email: "admin@gmail.com", senha: "Admin123", role: "admin"},
		{nome: "Usuario", email: "usuario@gmail.com", senha: "usuario123", role: "Usuario"},
	}

	for _, usuario := range usuarios {
		var existe bool
		if err := database.DB.QueryRow(`
			SELECT EXISTS(SELECT 1 FROM usuarios WHERE email = ?)
		`, usuario.email).Scan(&existe); err != nil {
			return err
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
