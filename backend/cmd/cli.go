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

// executarCLI roda o menu de terminal (login/cadastro + operações de
// estoque), em paralelo ao servidor web.
func executarCLI(reader *bufio.Reader) {
	for {
		fmt.Println()
		fmt.Println("========= GOSTOCK =========")
		fmt.Println("1 - Criar conta")
		fmt.Println("2 - Entrar")
		fmt.Println("3 - Sair")

		opcao, err := utilsLerOpcao(reader, "Escolha uma opção: ")
		if err != nil {
			fmt.Println("Opção inválida.")
			continue
		}

		switch opcao {
		case 1:
			services.CadastrarUsuario(reader)
		case 2:
			usuario, sucesso := services.Login(reader)
			if sucesso {
				menuEstoque(reader, usuario)
			}
		case 3:
			fmt.Println("Encerrando...")
			return
		default:
			fmt.Println("Opção inválida.")
		}
	}
}

// menuEstoque exibe as operações disponíveis para um usuário já
// autenticado no terminal.
func menuEstoque(reader *bufio.Reader, usuario *models.Usuario) {
	for {
		fmt.Println()
		ui.ExibirMenu()

		opcao, err := utilsLerOpcao(reader, "Escolha uma opção: ")
		if err != nil {
			fmt.Println("Opção inválida.")
			continue
		}

		switch opcao {
		case 1:
			services.CadastrarProduto(reader, usuario.ID)
		case 2:
			services.ListarProdutos()
		case 3:
			services.BuscarProduto(reader)
		case 4:
			services.RemoverProduto(reader)
		case 5:
			services.AtualizarProduto(reader, usuario.ID)
		case 6:
			services.AdicionarEstoque(reader, usuario.ID)
		case 7:
			services.RegistrarSaida(reader, usuario.ID)
		case 8:
			services.ListarMovimentacoes()
		case 9:
			fmt.Println("Saindo da conta...")
			return
		default:
			fmt.Println("Opção inválida.")
		}
	}
}

// utilsLerOpcao lê um número inteiro do terminal.
func utilsLerOpcao(reader *bufio.Reader, mensagem string) (int, error) {
	fmt.Print(mensagem)
	texto, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return 0, err
	}
	var opcao int
	if _, scanErr := fmt.Sscanf(texto, "%d", &opcao); scanErr != nil {
		return 0, scanErr
	}
	return opcao, nil
}
