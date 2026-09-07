package database

import (
	"database/sql"

	_ "gosqlite.org"
)

var DB *sql.DB

func Connect() error {
	var err error

	DB, err = sql.Open("sqlite", "backend/data/gostock.db")
	if err != nil {
		return err
	}
	DB.SetMaxOpenConns(1)

	if err := DB.Ping(); err != nil {
		return err
	}

	_, err = DB.Exec("PRAGMA foreign_keys = ON")
	return err
}

func CreateTables() error {
	query := `
		CREATE TABLE IF NOT EXISTS produtos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			nome TEXT NOT NULL,
			preco REAL NOT NULL,
			quantidade INTEGER NOT NULL,
			ativo INTEGER NOT NULL DEFAULT 1
		);
		CREATE TABLE IF NOT EXISTS usuarios (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nome TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		senha TEXT NOT NULL,
		role TEXT NOT NULL DEFAULT 'basico',
		ativo INTEGER NOT NULL DEFAULT 1
		);

	CREATE TABLE IF NOT EXISTS movimentacoes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		produto_id INTEGER NOT NULL,
		usuario_id INTEGER NOT NULL,
		tipo TEXT NOT NULL,
		quantidade INTEGER NOT NULL,
		data DATETIME DEFAULT CURRENT_TIMESTAMP,

		FOREIGN KEY (produto_id) REFERENCES produtos(id),
		FOREIGN KEY (usuario_id) REFERENCES usuarios(id)
);

    CREATE TABLE IF NOT EXISTS vendas (
       id INTEGER PRIMARY KEY AUTOINCREMENT,
       produto_id INTEGER NOT NULL,
       usuario_id INTEGER NOT NULL,
       quantidade INTEGER NOT NULL,
       valor_unitario REAL NOT NULL,
       valor_total REAL NOT NULL,
       forma_pagamento TEXT NOT NULL DEFAULT 'dinheiro',
       data DATETIME DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (produto_id) REFERENCES produtos(id),
    FOREIGN KEY (usuario_id) REFERENCES usuarios(id)
);

	CREATE TABLE IF NOT EXISTS sessoes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    usuario_id INTEGER NOT NULL,
    token_hash TEXT NOT NULL UNIQUE,
    expira_em DATETIME NOT NULL,
    criado_em DATETIME DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (usuario_id) REFERENCES usuarios(id)
);
	`

	_, err := DB.Exec(query)
	if err != nil {
		return err
	}

	var activeColumn int
	err = DB.QueryRow(`
		SELECT COUNT(*)
		FROM pragma_table_info('produtos')
		WHERE name = 'ativo'
	`).Scan(&activeColumn)
	if err != nil {
		return err
	}
	if activeColumn == 0 {
		if _, err = DB.Exec(`ALTER TABLE produtos ADD COLUMN ativo INTEGER NOT NULL DEFAULT 1`); err != nil {
			return err
		}
	}

	var roleColumn int
	err = DB.QueryRow(`
		SELECT COUNT(*)
		FROM pragma_table_info('usuarios')
		WHERE name = 'role'
	`).Scan(&roleColumn)
	if err != nil {
		return err
	}
	if roleColumn == 0 {
		if _, err = DB.Exec(`ALTER TABLE usuarios ADD COLUMN role TEXT NOT NULL DEFAULT 'basico'`); err != nil {
			return err
		}
	}

	var userActiveColumn int
	err = DB.QueryRow(`
		SELECT COUNT(*)
		FROM pragma_table_info('usuarios')
		WHERE name = 'ativo'
	`).Scan(&userActiveColumn)
	if err != nil {
		return err
	}
	if userActiveColumn == 0 {
		if _, err = DB.Exec(`ALTER TABLE usuarios ADD COLUMN ativo INTEGER NOT NULL DEFAULT 1`); err != nil {
			return err
		}
	}

	// O CEO é uma identidade reservada: somente admin@gmail.com pode possuir esse cargo.
	// Isso corrige o banco a cada inicialização, mesmo que alguém tenha mexido direto nele.
	if _, err = DB.Exec(`
		UPDATE usuarios
		SET role = CASE
			WHEN LOWER(TRIM(email)) = 'admin@gmail.com' THEN 'ceo'
			WHEN LOWER(TRIM(role)) = 'ceo' THEN 'basico'
			ELSE role
		END
	`); err != nil {
		return err
	}

	var paymentMethodColumn int
	err = DB.QueryRow(`
		SELECT COUNT(*)
		FROM pragma_table_info('vendas')
		WHERE name = 'forma_pagamento'
	`).Scan(&paymentMethodColumn)
	if err != nil {
		return err
	}
	if paymentMethodColumn == 0 {
		_, err = DB.Exec(`ALTER TABLE vendas ADD COLUMN forma_pagamento TEXT NOT NULL DEFAULT 'dinheiro'`)
	}

	return err
}
