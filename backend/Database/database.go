package database

import (
	"database/sql"

	_ "gosqlite.org"
)

var DB *sql.DB

func Conectar() error {
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

func CriarTabelas() error {
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
		role TEXT NOT NULL DEFAULT 'basico'
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

	var colunaAtivo int
	err = DB.QueryRow(`
		SELECT COUNT(*)
		FROM pragma_table_info('produtos')
		WHERE name = 'ativo'
	`).Scan(&colunaAtivo)
	if err != nil {
		return err
	}
	if colunaAtivo == 0 {
		if _, err = DB.Exec(`ALTER TABLE produtos ADD COLUMN ativo INTEGER NOT NULL DEFAULT 1`); err != nil {
			return err
		}
	}

	var colunaRole int
	err = DB.QueryRow(`
		SELECT COUNT(*)
		FROM pragma_table_info('usuarios')
		WHERE name = 'role'
	`).Scan(&colunaRole)
	if err != nil {
		return err
	}
	if colunaRole == 0 {
		if _, err = DB.Exec(`ALTER TABLE usuarios ADD COLUMN role TEXT NOT NULL DEFAULT 'basico'`); err != nil {
			return err
		}
	}

	var colunaFormaPagamento int
	err = DB.QueryRow(`
		SELECT COUNT(*)
		FROM pragma_table_info('vendas')
		WHERE name = 'forma_pagamento'
	`).Scan(&colunaFormaPagamento)
	if err != nil {
		return err
	}
	if colunaFormaPagamento == 0 {
		_, err = DB.Exec(`ALTER TABLE vendas ADD COLUMN forma_pagamento TEXT NOT NULL DEFAULT 'dinheiro'`)
	}

	return err
}
