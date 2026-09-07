package utils

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

func ReadText(reader *bufio.Reader, message string) string {
	fmt.Print(message)

	text, _ := reader.ReadString('\n')

	return strings.TrimSpace(text)
}

func ReadInt(reader *bufio.Reader, message string) (int, error) {
	text := ReadText(reader, message)

	return strconv.Atoi(text)
}

func ReadFloat(reader *bufio.Reader, message string) (float64, error) {
	text := ReadText(reader, message)

	text = strings.ReplaceAll(text, ",", ".")

	return strconv.ParseFloat(text, 64)
}

func ReadValidName(reader *bufio.Reader) string {
	for {
		name := ReadText(reader, "Nome:")
		if ValidateName(name) {
			return name
		}
		fmt.Println("O nome não pode ser vazio.")
	}
}

func ReadValidQuantity(reader *bufio.Reader, message string) int {
	for {
		quantity, err := ReadInt(reader, message)

		if err != nil {
			fmt.Println("Quantidade inválida. Tente novamente.")
			continue
		}
		if !ValidateQuantity(quantity) {
			fmt.Println("A quantidade não pode ser negativa.")
			continue
		}
		return quantity
	}
}

func ReadValidPrice(reader *bufio.Reader, message string) float64 {
	for {
		price, err := ReadFloat(reader, message)

		if err != nil {
			fmt.Println("Preço inválido. Tente novamente.")
			continue
		}

		if !ValidatePrice(price) {
			fmt.Println("O preço não pode ser negativo.")
			continue
		}

		return price
	}
}
