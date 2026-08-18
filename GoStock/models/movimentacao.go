package models

import("time"
)

type Movimentacoes struct{
	ID int 
	Usuario string
	Senha int
	Tipo string
	Valor float64
	Data time.Time
}