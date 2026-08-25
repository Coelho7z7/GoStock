package models

import "time"

type Movimentacao struct {
	ID         int
	ProdutoID  int
	UsuarioID  int
	Tipo       string
	Quantidade int
	Data       time.Time
}
