package services

import (
	database "gostock/backend/Database"
	"gostock/backend/models"
)

func ExistemUsuarios() bool {
	var quantidade int

	err := database.DB.QueryRow(`SELECT COUNT(*) FROM usuarios`).Scan(&quantidade)
	if err != nil {
		return false
	}

	return quantidade > 0
}

func BuscarUsuarioPorID(id int) (*models.Usuario, error) {
	var usuario models.Usuario

	err := database.DB.QueryRow(`
		SELECT id, nome, email
		FROM usuarios
		WHERE id = ?
	`, id).Scan(
		&usuario.ID,
		&usuario.Nome,
		&usuario.Email,
	)

	if err != nil {
		return nil, err
	}

	return &usuario, nil
}
