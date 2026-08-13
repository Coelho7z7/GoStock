package utils

import (
	"strings"
	"unicode"
)

func ValidarNome(nome string) bool {
	return strings.TrimSpace(nome) != ""
}

func ValidarQuantidade(quantidade int) bool {
	return quantidade >= 0
}

func ValidarPreco(preco float64) bool {
	return preco >= 0
}

func ValidarEmail(email string) bool {
	email = strings.TrimSpace(email)

	if !strings.HasSuffix(email, "@gmail.com") {
		return false
	}

	if strings.Count(email, "@") != 1 {
		return false
	}

	if strings.HasPrefix(email, "@") {
		return false
	}

	return true
}

func ValidarSenha(senha string) bool {
	if len(senha) < 6 {
		return false
	}

	for _, caractere := range senha {
		if !unicode.IsLetter(caractere) && !unicode.IsDigit(caractere) {
			return true
		}
	}

	return false
}
