package main

import (
	"fmt"
	"net/http"
	"os"

	database "gostock/backend/Database"
	"gostock/backend/services"
)

func main() {
	if err := prepararDiretorioProjeto(); err != nil {
		fmt.Println("Erro ao localizar os arquivos do projeto:", err)
		os.Exit(1)
	}

	if err := database.Conectar(); err != nil {
		fmt.Println("Erro ao conectar ao banco de dados:", err)
		os.Exit(1)
	}
	defer database.DB.Close()

	if err := database.CriarTabelas(); err != nil {
		fmt.Println("Erro ao preparar as tabelas do banco de dados:", err)
		os.Exit(1)
	}
	if err := services.SeedUsuariosPadrao(); err != nil {
		fmt.Println("Erro ao criar usuários padrão:", err)
		os.Exit(1)
	}

	registrarRotas()

	porta := os.Getenv("PORT")
	if porta == "" {
		porta = "8080"
	}

	fmt.Println("Servidor web disponível na porta", porta)
	if err := http.ListenAndServe(":"+porta, nil); err != nil {
		fmt.Println("Erro no servidor web:", err)
		os.Exit(1)
	}
}

// registrarRotas conecta cada rota HTTP ao seu handler correspondente
// e configura os servidores de arquivos estáticos (CSS/JS).
func registrarRotas() {
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("frontend/css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("frontend/js"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))

	http.HandleFunc("/", handlerIndex)
	http.HandleFunc("/login", handlerLogin)
	http.HandleFunc("/logout", handlerLogout)

	http.HandleFunc("/dashboard", handlerDashboard)

	http.HandleFunc("/produtos", handlerProdutos)
	http.HandleFunc("/alterar-produto", handlerAlterarProduto)

	http.HandleFunc("/estoque", handlerEstoque)

	http.HandleFunc("/vendas", handlerVendas)
	http.HandleFunc("/api/vendas", handlerApiVendas)

	http.HandleFunc("/movimentacoes", handlerMovimentacoes)

	http.HandleFunc("/usuarios", handlerUsuarios)
}
