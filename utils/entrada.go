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

func LerNomeValido(reader *bufio.Reader) string {
	for {
		nome := LerTexto(reader, "Nome:")
		if ValidarNome(nome) {
			return nome
		}
		fmt.Println("O nom não pode ser vázio.")
	}
}

func LerQuantidadeValida(reader *bufio.Reader, mensagem string) int {
	for {
		quantidade, err := LerInteiro(reader, mensagem)

		if err != nil {
			fmt.Println("Quantidade inválida. Tente novamente.")
			continue
		}
		if !ValidarQuantidade(quantidade) {
			fmt.Println("A quantidade não pode ser negativa.")
			continue
		}
		return quantidade
	}
}

func LerPrecoValido(reader *bufio.Reader, mensagem string) float64 {
	for {
		preco, err := LerFloat(reader, mensagem)

		if err != nil {
			fmt.Println("Preço inválido. Tente novamente.")
			continue
		}

		if !ValidarPreco(preco) {
			fmt.Println("O preço não pode ser negativo.")
			continue
		}

		return preco
	}
}
