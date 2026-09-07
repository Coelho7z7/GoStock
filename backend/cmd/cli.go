package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"

	"gostock/backend/models"
	"gostock/backend/services"
	"gostock/backend/ui"
)

// runCLI roda o menu de terminal (login/cadastro + operações de
// estoque), em paralelo ao servidor web.
func runCLI(reader *bufio.Reader) {
	for {
		fmt.Println()
		fmt.Println("========= GOSTOCK =========")
		fmt.Println("1 - Criar conta")
		fmt.Println("2 - Entrar")
		fmt.Println("3 - Sair")

		option, err := readMenuOption(reader, "Escolha uma opção: ")
		if err != nil {
			fmt.Println("Opção inválida.")
			continue
		}

		switch option {
		case 1:
			services.CreateUser(reader)
		case 2:
			user, success := services.Login(reader)
			if success {
				stockMenu(reader, user)
			}
		case 3:
			fmt.Println("Encerrando...")
			return
		default:
			fmt.Println("Opção inválida.")
		}
	}
}

// stockMenu exibe as operações disponíveis para um usuário já
// autenticado no terminal.
func stockMenu(reader *bufio.Reader, user *models.User) {
	for {
		fmt.Println()
		ui.ShowMenu()

		option, err := readMenuOption(reader, "Escolha uma opção: ")
		if err != nil {
			fmt.Println("Opção inválida.")
			continue
		}

		switch option {
		case 1:
			services.CreateProduct(reader, user.ID)
		case 2:
			services.ListProducts()
		case 3:
			services.FindProduct(reader)
		case 4:
			services.DeleteProduct(reader)
		case 5:
			services.UpdateProduct(reader, user.ID)
		case 6:
			services.AddStock(reader, user.ID)
		case 7:
			services.RegisterStockExit(reader, user.ID)
		case 8:
			services.ListMovements()
		case 9:
			fmt.Println("Saindo da conta...")
			return
		default:
			fmt.Println("Opção inválida.")
		}
	}
}

// readMenuOption lê um número inteiro do terminal.
func readMenuOption(reader *bufio.Reader, message string) (int, error) {
	fmt.Print(message)
	text, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return 0, err
	}
	var option int
	if _, scanErr := fmt.Sscanf(text, "%d", &option); scanErr != nil {
		return 0, scanErr
	}
	return option, nil
}
