package database

import (
	"database/sql"

	_ "gosqlite.org"
)

var DB *sql.DB

func Conectar() error {
	var err error

	DB, err = sql.Open("sqlite", "data/gostock.db")
	if err != nil {
		return err
	}

	return DB.Ping()
}

func CriarTabelas() error {
	query := `
		CREATE TABLE IF NOT EXISTS produtos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			nome TEXT NOT NULL,
			preco REAL NOT NULL,
			quantidade INTEGER NOT NULL
		);
	`

	_, err := DB.Exec(query)

	return err
}
