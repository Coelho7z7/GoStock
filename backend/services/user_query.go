package services

import (
	database "gostock/backend/database"
	"gostock/backend/models"
)

func GetUserByID(id int) (*models.User, error) {
	var user models.User

	err := database.DB.QueryRow(`
		SELECT id, nome, email, role, ativo
		FROM usuarios
		WHERE id = ?
	`, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Role,
		&user.Active,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
