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
