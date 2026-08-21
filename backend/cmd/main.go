package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"

	database "gostock/backend/Database"
	"gostock/backend/models"
	"gostock/backend/services"
	"gostock/backend/ui"
	"gostock/backend/utils"
	"html/template"
)

type DashboardData struct {
	Usuario  *models.Usuario
	Produtos []models.Produto
}

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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "frontend/index.html")
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("frontend"))))

	http.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("usuario_id")

		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		usuarioID, err := strconv.Atoi(cookie.Value)

		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		usuario, err := services.BuscarUsuarioPorID(usuarioID)

		if err != nil {
			http.Error(w, "Usuário não encontrado", http.StatusInternalServerError)
			return
		}

		produtos, err := services.BuscarTodosProdutos()

		if err != nil {
			http.Error(w, "Erro ao buscar produtos", http.StatusInternalServerError)
			return
		}

		dados := DashboardData{
			Usuario:  usuario,
			Produtos: produtos,
		}

		tmpl, err := template.ParseFiles("frontend/dashboard.html")

		if err != nil {
			fmt.Println("Erro ao carregar dashboard:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, dados)

		if err != nil {
			http.Error(w, "Erro ao renderizar página", http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "POST" {
			fmt.Fprintln(w, "Método não permitido.")
			return
		}

		email := r.FormValue("email")
		senha := r.FormValue("senha")

		usuario, sucesso := services.AutenticarUsuario(email, senha)

		if !sucesso {
			fmt.Fprintln(w, "Email ou senha incorretos.")
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "usuario_id",
			Value: fmt.Sprintf("%d", usuario.ID),
			Path:  "/",
		})

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	})

	go func() {
		fmt.Println("Servidor rodando em http://localhost:8080")

		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			fmt.Println("Erro no servidor:", err)
		}
	}()

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
			services.RegistrarVenda(reader, usuario.ID)

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
