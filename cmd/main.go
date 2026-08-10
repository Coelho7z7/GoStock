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
	produtos := []models.Produto{}
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

		case 2:
			services.ListarProdutos(produtos)

		case 3:
			services.BuscarProduto(produtos, reader)

		case 4:
			services.RemoverProduto(&produtos, reader)

		case 5:
			services.AtualizarProduto(produtos, reader)

		case 6:
			fmt.Println("Até mais!")
			return

		default:
			fmt.Println("Opção inválida.")
		}
	}
}
