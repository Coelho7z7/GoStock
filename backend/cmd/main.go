package main

import (
	"fmt"
	"net/http"
	"os"

	database "gostock/backend/database"
	"gostock/backend/services"
)

func main() {
	if err := prepareProjectDirectory(); err != nil {
		fmt.Println("Erro ao localizar os arquivos do projeto:", err)
		os.Exit(1)
	}

	if err := database.Connect(); err != nil {
		fmt.Println("Erro ao conectar ao banco de dados:", err)
		os.Exit(1)
	}
	defer database.DB.Close()

	if err := database.CreateTables(); err != nil {
		fmt.Println("Erro ao preparar as tabelas do banco de dados:", err)
		os.Exit(1)
	}

	if len(os.Args) > 1 && os.Args[1] == "reset-password" {
		runResetPasswordCommand(os.Args[2:])
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "create-user" {
		runCreateUserCommand(os.Args[2:])
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "rename-user" {
		runRenameUserCommand(os.Args[2:])
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "change-email" {
		runChangeEmailCommand(os.Args[2:])
		return
	}

	if err := services.SeedDefaultUsers(); err != nil {
		fmt.Println("Erro ao criar usuários padrão:", err)
		os.Exit(1)
	}

	registerRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Servidor web disponível na porta", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println("Erro no servidor web:", err)
		os.Exit(1)
	}
}

// runResetPasswordCommand troca a senha de uma conta já existente via
// linha de comando, ex.:
//
//	go run ./backend/cmd reset-password ceo@gmail.com NovaSenha!123
//
// Encerra o processo sem subir o servidor web. Existe porque ainda não há
// uma tela no painel para trocar a senha de um usuário já criado (só na
// criação e no seed inicial).
func runResetPasswordCommand(args []string) {
	if len(args) != 2 {
		fmt.Println("Uso: reset-password <email> <nova-senha>")
		os.Exit(1)
	}

	if err := services.ResetPassword(args[0], args[1]); err != nil {
		fmt.Println("Erro ao trocar a senha:", err)
		os.Exit(1)
	}

	fmt.Println("Senha atualizada com sucesso para", args[0])
}

// runCreateUserCommand cadastra um usuário via linha de comando, ex.:
//
//	go run ./backend/cmd create-user "Nome" email@gmail.com "Senha!123" admin
//
// Permissões aceitas: admin, gerente, basico (nunca "ceo" — reservado a
// ceo@gmail.com e criado apenas pelo seed). Encerra o processo sem subir
// o servidor web.
func runCreateUserCommand(args []string) {
	if len(args) != 4 {
		fmt.Println("Uso: create-user <nome> <email> <senha> <admin|gerente|basico>")
		os.Exit(1)
	}

	if err := services.CreateUserWeb(args[0], args[1], args[2], args[3]); err != nil {
		fmt.Println("Erro ao criar usuário:", err)
		os.Exit(1)
	}

	fmt.Println("Usuário criado com sucesso:", args[1])
}

// runRenameUserCommand troca o nome de exibição de uma conta já existente
// via linha de comando, ex.:
//
//	go run ./backend/cmd rename-user matheus@gmail.com "Novo Nome"
//
// Encerra o processo sem subir o servidor web.
func runRenameUserCommand(args []string) {
	if len(args) != 2 {
		fmt.Println("Uso: rename-user <email> <novo-nome>")
		os.Exit(1)
	}

	if err := services.RenameUser(args[0], args[1]); err != nil {
		fmt.Println("Erro ao renomear usuário:", err)
		os.Exit(1)
	}

	fmt.Println("Nome atualizado com sucesso para", args[0])
}

// runChangeEmailCommand troca o email de uma conta já existente via linha
// de comando, ex.:
//
//	go run ./backend/cmd change-email matheus@gmail.com gerente@gmail.com
//
// Encerra o processo sem subir o servidor web.
func runChangeEmailCommand(args []string) {
	if len(args) != 2 {
		fmt.Println("Uso: change-email <email-atual> <novo-email>")
		os.Exit(1)
	}

	if err := services.ChangeUserEmail(args[0], args[1]); err != nil {
		fmt.Println("Erro ao trocar email:", err)
		os.Exit(1)
	}

	fmt.Println("Email atualizado com sucesso:", args[0], "->", args[1])
}

// registerRoutes conecta cada rota HTTP ao seu handler correspondente
// e configura os servidores de arquivos estáticos (CSS/JS).
func registerRoutes() {
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("frontend/css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("frontend/js"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)

	http.HandleFunc("/dashboard", dashboardHandler)

	http.HandleFunc("/produtos", productHandler)
	http.HandleFunc("/alterar-produto", editProductHandler)

	http.HandleFunc("/estoque", stockHandler)

	http.HandleFunc("/vendas", saleHandler)
	http.HandleFunc("/api/vendas", saleAPIHandler)

	http.HandleFunc("/movimentacoes", movementHandler)

	http.HandleFunc("/usuarios", userHandler)
}
