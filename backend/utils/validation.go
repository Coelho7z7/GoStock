package utils

import (
	"strings"
	"unicode"
)

func ValidateName(name string) bool {
	return strings.TrimSpace(name) != ""
}

func ValidateQuantity(quantity int) bool {
	return quantity >= 0
}

func ValidatePrice(price float64) bool {
	return price >= 0
}

func ValidateEmail(email string) bool {
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

func ValidatePassword(password string) bool {
	if len(password) < 6 {
		return false
	}

	for _, char := range password {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) {
			return true
		}
	}

	return false
}
