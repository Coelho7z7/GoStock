package utils

import "strings"

func ValidarNome(nome string) bool {
	return strings.TrimSpace(nome) != ""
}

func ValidarQuantidade(quantidade int) bool {
	return quantidade >= 0
}

func ValidarPreco(preco float64) bool {
	return preco >= 0
}
