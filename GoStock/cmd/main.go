package main

import (
	"bufio"
	"fmt"
	"os"

	database "gostock/Database"
	"gostock/models"
	"gostock/services"
	"gostock/ui"
	"gostock/utils"
)

func main() {

	err := database.Conectar()
	if err != nil {
		fmt.Println("Erro ao conectar ao banco:", err)
		return
	}

	defer database.DB.Close()

	err = database.CriarTabelas()
	if err != nil {
		fmt.Println("Erro ao criar tabelas:", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("========= GoStock =========")
		fmt.Println("1 - Entrar")
		fmt.Println("2 - Criar conta")
		fmt.Println("3 - Sair")

		opcao, err := utils.LerInteiro(reader, "Escolha uma opção: ")

		if err != nil {
			fmt.Println("Opção inválida.")
			continue
		}

		switch opcao {

		case 1:
			usuario, sucesso := services.Login(reader)

			if sucesso {
				MenuEstoque(reader, usuario)
			}
		case 2:
			services.CadastrarUsuario(reader)

		case 3:
			fmt.Println("Até mais!")
			return

		default:
			fmt.Println("Opção inválida.")
		}
	}
}

func MenuEstoque(reader *bufio.Reader, usuario *models.Usuario) {

	for {
		ui.ExibirMenu()

		opcao, err := utils.LerInteiro(reader, "Escolha uma opção: ")

		if err != nil {
			fmt.Println("Opção inválida.")
			continue
		}

		switch opcao {

		case 1:
			services.CadastrarProduto(reader)

		case 2:
			services.ListarProdutos()

		case 3:
			services.BuscarProduto(reader)

		case 4:
			services.RemoverProduto(reader)

		case 5:
			services.AtualizarProduto(reader)

		case 6:
			fmt.Println("Saindo da conta...")
			return

		default:
			fmt.Println("Opção inválida.")
		}
	}
}
