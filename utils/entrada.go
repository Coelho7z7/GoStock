package utils

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

func LerTexto(reader *bufio.Reader, mensagem string) string {
	fmt.Print(mensagem)

	texto, _ := reader.ReadString('\n')

	return strings.TrimSpace(texto)
}

func LerInteiro(reader *bufio.Reader, mensagem string) (int, error) {
	texto := LerTexto(reader, mensagem)

	return strconv.Atoi(texto)
}

func LerFloat(reader *bufio.Reader, mensagem string) (float64, error) {
	texto := LerTexto(reader, mensagem)

	texto = strings.ReplaceAll(texto, ",", ".")

	return strconv.ParseFloat(texto, 64)
}
