package services

import (
	database "gostock/backend/Database"
	"gostock/backend/models"
)

func BuscarUsuarioPorID(id int) (*models.Usuario, error) {
	var usuario models.Usuario

	err := database.DB.QueryRow(`
		SELECT id, nome, email, role
		FROM usuarios
		WHERE id = ?
	`, id).Scan(
		&usuario.ID,
		&usuario.Nome,
		&usuario.Email,
		&usuario.Role,
	)

	if err != nil {
		return nil, err
	}

	return &usuario, nil
}
