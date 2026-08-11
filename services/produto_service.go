package services

import (
	"bufio"
	"fmt"
	"strings"

	"gostock/models"
	"gostock/utils"
)

func proximoID(produtos []models.Produto) int {
	maiorID := 0

	for _, produto := range produtos {
		if produto.ID > maiorID {
			maiorID = produto.ID
		}
	}
	return maiorID + 1
}

func CadastrarProduto(produtos *[]models.Produto, reader *bufio.Reader) {
	nome := utils.LerNomeValido(reader)
	quantidade := utils.LerQuantidadeValida(reader, "Quantidade:")
	preco := utils.LerPrecoValido(reader, "Preço:")

	produto := models.Produto{
		ID:         proximoID(*produtos),
		Nome:       nome,
		Quantidade: quantidade,
		Preco:      preco,
	}
	*produtos = append(*produtos, produto)
	err := utils.SalvarProdutos(*produtos)

	if err != nil {
		fmt.Println("Erro ao salvar produtos:", err)
		return
	}
	fmt.Println("Produto cadastrado com sucesso!")
}

func ListarProdutos(produtos []models.Produto) {
	if len(produtos) == 0 {
		fmt.Println("Nenhum produto cadastrado.")
		return
	}

	for _, produto := range produtos {
		fmt.Println("ID:", produto.ID)
		fmt.Println("Nome:", produto.Nome)
		fmt.Println("Quantidade:", produto.Quantidade)
		fmt.Println("Preço:", produto.Preco)
		fmt.Println("----------------------")
	}
}

func BuscarProduto(produtos []models.Produto, reader *bufio.Reader) {
	buscar := strings.TrimSpace(utils.LerTexto(reader, "Digite o nome do produto: "))
	encontrado := false

	for _, produto := range produtos {
		if strings.Contains(strings.ToLower(produto.Nome), strings.ToLower(buscar)) {
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
	id, err := utils.LerInteiro(reader, "Digite o ID do produto: ")

	if err != nil {
		fmt.Println("ID inválido.")
		return
	}
	removido := false

	for i, produto := range *produtos {
		if produto.ID == id {
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
	id, err := utils.LerInteiro(reader, "Digite o ID do produto: ")

	if err != nil {
		fmt.Println("ID inválido.")
		return
	}

	atualizado := false

	for i, produto := range produtos {
		if produto.ID == id {
			fmt.Println("Produto encontrado.")

			novoNome := utils.LerNomeValido(reader)

			novaQuantidade := utils.LerQuantidadeValida(
				reader,
				"Digite sua nova quantidade: ",
			)

			novoPreco := utils.LerPrecoValido(
				reader,
				"Digite seu novo preço: ",
			)

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
