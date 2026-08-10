package services

import (
	"bufio"
	"fmt"
	"strings"

	"gostock/models"
	"gostock/utils"
)

func CadastrarProduto(produtos *[]models.Produto, reader *bufio.Reader) {
	var nome string

	for {
		nome = utils.LerTexto(reader, "Nome: ")

		if nome == "" {
			fmt.Println("O nome não pode ficar vazio.")
			continue
		}

		break
	}

	var quantidade int

	for {
		var err error

		quantidade, err = utils.LerInteiro(reader, "Quantidade: ")

		if err != nil {
			fmt.Println("Quantidade inválida. Tente novamente.")
			continue
		}
		if quantidade < 0 {
			fmt.Println("A quantidade não pode ser negativa.")
			continue
		}

		break
	}

	var preco float64

	for {
		var err error

		preco, err = utils.LerFloat(reader, "Preço: ")

		if err != nil {
			fmt.Println("Preço inválido. Tente novamente.")
			continue
		}
		if preco < 0 {
			fmt.Println("O preço não pode ser negativa")
			continue
		}

		break
	}

	produto := models.Produto{
		Nome:       nome,
		Quantidade: quantidade,
		Preco:      preco,
	}

	*produtos = append(*produtos, produto)

	fmt.Println("Produto cadastrado com sucesso!")
}

func ListarProdutos(produtos []models.Produto) {
	if len(produtos) == 0 {
		fmt.Println("Nenhum produto cadastrado.")
		return
	}

	for _, produto := range produtos {
		fmt.Println("Nome:", produto.Nome)
		fmt.Println("Quantidade:", produto.Quantidade)
		fmt.Println("Preço:", produto.Preco)
		fmt.Println("----------------------")
	}
}

func BuscarProduto(produtos []models.Produto, reader *bufio.Reader) {
	buscar := utils.LerTexto(reader, "Digite o nome do produto: ")
	encontrado := false

	for _, produto := range produtos {
		if strings.EqualFold(buscar, produto.Nome) {
			fmt.Println("Produto encontrado!")
			fmt.Println("Nome:", produto.Nome)
			fmt.Println("Quantidade:", produto.Quantidade)
			fmt.Println("Preço:", produto.Preco)

			encontrado = true
		}
	}

	if !encontrado {
		fmt.Println("Produto não encontrado.")
	}
}

func RemoverProduto(produtos *[]models.Produto, reader *bufio.Reader) {
	remover := utils.LerTexto(reader, "Digite o nome do produto: ")
	removido := false

	for i, produto := range *produtos {
		if strings.EqualFold(remover, produto.Nome) {
			*produtos = append((*produtos)[:i], (*produtos)[i+1:]...)

			fmt.Println("Produto removido com sucesso.")
			removido = true
			break
		}
	}

	if !removido {
		fmt.Println("Produto não encontrado.")
	}
}

func AtualizarProduto(produtos []models.Produto, reader *bufio.Reader) {
	atualizar := utils.LerTexto(reader, "Digite o nome do produto: ")
	atualizado := false

	for i, produto := range produtos {
		if strings.EqualFold(atualizar, produto.Nome) {
			fmt.Println("Produto encontrado.")

			var novoNome string

			for {
				novoNome = utils.LerTexto(reader, "Digite seu novo nome: ")

				if novoNome == "" {
					fmt.Println("O nome não pode ficar vazio.")
					continue
				}

				break
			}

			var novaQuantidade int

			for {
				var err error

				novaQuantidade, err = utils.LerInteiro(
					reader,
					"Digite sua nova quantidade: ",
				)

				if err != nil {
					fmt.Println("Quantidade inválida. Tente novamente.")
					continue
				}

				if novaQuantidade < 0 {
					fmt.Println("A quantidade não pode ser negativa.")
					continue
				}

				break
			}

			var novoPreco float64

			for {
				var err error

				novoPreco, err = utils.LerFloat(
					reader,
					"Digite seu novo preço: ",
				)

				if err != nil {
					fmt.Println("Preço inválido. Tente novamente.")
					continue
				}

				if novoPreco < 0 {
					fmt.Println("O preço não pode ser negativo.")
					continue
				}

				break
			}

			produtos[i].Nome = novoNome
			produtos[i].Quantidade = novaQuantidade
			produtos[i].Preco = novoPreco

			fmt.Println("Produto atualizado com sucesso!")
			atualizado = true
			break
		}
	}

	if !atualizado {
		fmt.Println("Produto não encontrado.")
	}
}
