package main

import (
	"bufio"
	"fmt"
	"os"

	"gostock/models"
	"gostock/services"
	"gostock/ui"
	"gostock/utils"
)

func main() {
	produtos, err := utils.CarregarProdutos()

	if err != nil {
		fmt.Println("Não foi possível carregar os produtos:", err)
		produtos = []models.Produto{}
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		ui.ExibirMenu()

		opcao, err := utils.LerInteiro(reader, "Escolha uma opção: ")
		if err != nil {
			fmt.Println("Opção inválida.")
			continue
		}

		switch opcao {
		case 1:
			services.CadastrarProduto(&produtos, reader)

			err := utils.SalvarProdutos(produtos)
			if err != nil {
				fmt.Println("Erro ao salvar produtos:", err)
			}

		case 2:
			services.ListarProdutos(produtos)

		case 3:
			services.BuscarProduto(produtos, reader)

		case 4:
			services.RemoverProduto(&produtos, reader)

			err := utils.SalvarProdutos(produtos)
			if err != nil {
				fmt.Println("Erro ao salvar produtos:", err)
			}

		case 5:
			services.AtualizarProduto(produtos, reader)

			err := utils.SalvarProdutos(produtos)
			if err != nil {
				fmt.Println("Erro ao salvar produtos:", err)
			}

		case 6:
			fmt.Println("Até mais!")
			return

		default:
			fmt.Println("Opção inválida.")
		}
	}
}
