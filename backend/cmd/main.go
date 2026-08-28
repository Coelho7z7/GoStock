package main

import (
	"bufio"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	database "gostock/backend/Database"
	"gostock/backend/models"
	"gostock/backend/services"
	"gostock/backend/ui"
	"gostock/backend/utils"
)

type DashboardData struct {
	Usuario  *models.Usuario
	Produtos []models.Produto
	Resumo   ResumoData
}

type ResumoData struct {
	TotalEstoque       int
	TotalVendas        int
	TotalMovimentacoes int
}

func main() {
	if err := prepararDiretorioProjeto(); err != nil {
		fmt.Println("Erro ao localizar a raiz do projeto:", err)
		return
	}

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

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("frontend/css"))))

	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("frontend/js"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		tmpl, err := template.ParseFiles("frontend/html/index.html")
		if err != nil {
			http.Error(w, "Erro ao carregar tela de login", http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, struct {
			Email string
			Erro  string
		}{}); err != nil {
			http.Error(w, "Erro ao renderizar tela de login", http.StatusInternalServerError)
		}
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("frontend"))))

	http.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		usuarioID, ok := usuarioDaSessao(r)
		if !ok {
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

		resumo, err := carregarResumo()
		if err != nil {
			http.Error(w, "Erro ao carregar resumo", http.StatusInternalServerError)
			return
		}

		dados := DashboardData{
			Usuario:  usuario,
			Produtos: produtos,
			Resumo:   resumo,
		}

		tmpl, err := template.ParseFiles("frontend/html/dashboard.html")

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

	http.HandleFunc("/estoque", func(w http.ResponseWriter, r *http.Request) {
		if _, ok := usuarioDaSessao(r); !ok {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		produtos, err := services.BuscarTodosProdutos()
		if err != nil {
			http.Error(w, "Erro ao buscar produtos", http.StatusInternalServerError)
			return
		}
		tmpl, err := template.ParseFiles("frontend/html/estoque.html")
		if err != nil {
			http.Error(w, "Erro ao carregar estoque", http.StatusInternalServerError)
			return
		}
		mensagens := map[string]string{
			"entrada":    "Estoque adicionado com sucesso.",
			"saida":      "Saída registrada com sucesso.",
			"cadastrado": "Produto cadastrado com sucesso.",
		}
		if err := tmpl.Execute(w, struct {
			Produtos []models.Produto
			Mensagem string
		}{produtos, mensagens[r.URL.Query().Get("sucesso")]}); err != nil {
			http.Error(w, "Erro ao renderizar estoque", http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/produtos", func(w http.ResponseWriter, r *http.Request) {
		usuarioID, ok := usuarioDaSessao(r)
		if !ok {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		dados := struct {
			Produtos   []models.Produto
			Nome       string
			Quantidade int
			Preco      string
			Mensagem   string
			Erro       string
		}{}

		dados.Mensagem = map[string]string{
			"cadastrado": "Produto cadastrado com sucesso.",
			"atualizado": "Produto atualizado com sucesso.",
			"removido":   "Produto removido com sucesso.",
			"entrada":    "Estoque adicionado com sucesso.",
			"saida":      "Saída registrada com sucesso.",
		}[r.URL.Query().Get("sucesso")]

		if r.Method == http.MethodPost {
			dados.Nome = strings.TrimSpace(r.FormValue("nome"))
			quantidadeTexto := strings.TrimSpace(r.FormValue("quantidade"))
			dados.Preco = strings.TrimSpace(r.FormValue("preco"))

			quantidade, quantidadeErr := strconv.Atoi(quantidadeTexto)
			preco, precoErr := strconv.ParseFloat(dados.Preco, 64)

			if dados.Nome == "" {
				dados.Erro = "Informe o nome do produto."
			} else if quantidadeErr != nil || quantidade < 0 {
				dados.Erro = "Informe uma quantidade válida."
			} else if precoErr != nil || preco < 0 {
				dados.Erro = "Informe um preço válido."
			} else if err := services.CadastrarProdutoWeb(
				dados.Nome,
				quantidade,
				preco,
				usuarioID,
			); err != nil {
				dados.Erro = err.Error()
			} else {
				http.Redirect(w, r, "/produtos?sucesso=cadastrado", http.StatusSeeOther)
				return
			}

			dados.Quantidade = quantidade
		} else if r.Method != http.MethodGet {
			w.Header().Set("Allow", "GET, POST")
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		produtos, err := services.BuscarTodosProdutos()
		if err != nil {
			http.Error(w, "Erro ao buscar produtos", http.StatusInternalServerError)
			return
		}

		dados.Produtos = produtos

		tmpl, err := template.ParseFiles("frontend/html/produtos.html")
		if err != nil {
			http.Error(w, "Erro ao carregar produtos", http.StatusInternalServerError)
			return
		}

		if dados.Erro != "" {
			w.WriteHeader(http.StatusBadRequest)
		}

		if err := tmpl.Execute(w, dados); err != nil {
			http.Error(w, "Erro ao renderizar produtos", http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/alterar-produto", func(w http.ResponseWriter, r *http.Request) {
		usuarioID, ok := usuarioDaSessao(r)
		if !ok {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		mensagens := map[string]string{
			"atualizado": "Produto atualizado com sucesso.",
			"removido":   "Produto removido com sucesso.",
		}
		dados := struct {
			Produtos []models.Produto
			Mensagem string
			Erro     string
		}{Mensagem: mensagens[r.URL.Query().Get("sucesso")]}

		if r.Method == http.MethodPost {
			produtoID, idErr := strconv.Atoi(r.FormValue("produto_id"))
			dados.Erro = "Produto inválido."
			if idErr == nil && r.FormValue("acao") == "remover" {
				opErr := services.RemoverProdutoWeb(produtoID)
				if opErr == nil {
					http.Redirect(w, r, "/produtos?sucesso=removido", http.StatusSeeOther)
					return
				}
				dados.Erro = opErr.Error()
			} else if idErr == nil && r.FormValue("acao") == "atualizar" {
				nome := strings.TrimSpace(r.FormValue("nome"))
				preco, precoErr := strconv.ParseFloat(strings.TrimSpace(r.FormValue("preco")), 64)
				if precoErr != nil {
					dados.Erro = "Informe um preço válido."
				} else {
					opErr := services.AtualizarProdutoWeb(produtoID, nome, preco, usuarioID)
					if opErr == nil {
						http.Redirect(w, r, "/produtos?sucesso=atualizado", http.StatusSeeOther)
						return
					}
					dados.Erro = opErr.Error()
				}
			}
		} else if r.Method != http.MethodGet {
			w.Header().Set("Allow", "GET, POST")
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		dados.Produtos, _ = services.BuscarTodosProdutos()
		tmpl, err := template.ParseFiles("frontend/html/alterar_produto.html")
		if err != nil {
			http.Error(w, "Erro ao carregar alteração de produto", http.StatusInternalServerError)
			return
		}
		if dados.Erro != "" {
			w.WriteHeader(http.StatusBadRequest)
		}
		if err := tmpl.Execute(w, dados); err != nil {
			http.Error(w, "Erro ao renderizar alteração de produto", http.StatusInternalServerError)
		}
	})
	http.HandleFunc("/vendas", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/estoque", http.StatusSeeOther)
	})

	http.HandleFunc("/saidas", func(w http.ResponseWriter, r *http.Request) {
		if _, ok := usuarioDaSessao(r); !ok {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		produtos, err := services.BuscarTodosProdutos()
		if err != nil {
			http.Error(w, "Erro ao buscar produtos", http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles("frontend/html/saidas.html")
		if err != nil {
			http.Error(w, "Erro ao carregar saídas", http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, struct {
			Produtos []models.Produto
		}{produtos}); err != nil {
			http.Error(w, "Erro ao renderizar saídas", http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/movimentacoes", func(w http.ResponseWriter, r *http.Request) {
		if _, ok := usuarioDaSessao(r); !ok {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		movimentacoes, err := services.BuscarMovimentacoesWeb()
		if err != nil {
			http.Error(w, "Erro ao buscar movimentações", http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles("frontend/html/movimentacoes.html")
		if err != nil {
			http.Error(w, "Erro ao carregar movimentações", http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, struct {
			Movimentacoes []models.Movimentacao
		}{movimentacoes}); err != nil {
			http.Error(w, "Erro ao renderizar movimentações", http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/estoque/adicionar", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}
		usuarioID, ok := usuarioDaSessao(r)
		if !ok {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		produtoID, err := strconv.Atoi(r.FormValue("produto_id"))
		quantidade, quantidadeErr := strconv.Atoi(r.FormValue("quantidade"))
		if err != nil || quantidadeErr != nil {
			http.Error(w, "Dados inválidos", http.StatusBadRequest)
			return
		}
		if err := services.AdicionarEstoqueWeb(produtoID, quantidade, usuarioID); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, "/produtos?sucesso=entrada", http.StatusSeeOther)
	})

	http.HandleFunc("/estoque/saida", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}
		usuarioID, ok := usuarioDaSessao(r)
		if !ok {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		produtoID, err := strconv.Atoi(r.FormValue("produto_id"))
		quantidade, quantidadeErr := strconv.Atoi(r.FormValue("quantidade"))
		if err != nil || quantidadeErr != nil {
			http.Error(w, "Dados inválidos", http.StatusBadRequest)
			return
		}
		if err := services.RegistrarSaidaWeb(produtoID, quantidade, usuarioID); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, "/produtos?sucesso=saida", http.StatusSeeOther)
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		email := r.FormValue("email")
		senha := r.FormValue("senha")

		usuario, sucesso := services.AutenticarUsuario(email, senha)
		if !sucesso {
			tmpl, err := template.ParseFiles("frontend/html/index.html")
			if err != nil {
				http.Error(w, "Erro ao carregar tela de login", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusUnauthorized)
			if err := tmpl.Execute(w, struct {
				Email string
				Erro  string
			}{email, "Email ou senha incorretos."}); err != nil {
				http.Error(w, "Erro ao renderizar tela de login", http.StatusInternalServerError)
			}
			return
		}

		token, err := criarSessao(usuario.ID)
		if err != nil {
			http.Error(w, "Erro ao iniciar sessão", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "sessao",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   60 * 60 * 24 * 30,
		})

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	})

	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		if cookie, err := r.Cookie("sessao"); err == nil {
			hash := sha256.Sum256([]byte(cookie.Value))
			tokenHash := hex.EncodeToString(hash[:])
			database.DB.Exec("DELETE FROM sessoes WHERE token_hash = ?", tokenHash)
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "sessao",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	go func() {
		fmt.Println("Servidor rodando em http://127.0.0.1:8081")

		err = http.ListenAndServe("127.0.0.1:8081", nil)

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

func prepararDiretorioProjeto() error {
	diretorio, err := os.Getwd()
	if err != nil {
		return err
	}

	for {
		if _, err := os.Stat(filepath.Join(diretorio, "frontend", "html", "index.html")); err == nil {
			return os.Chdir(diretorio)
		}
		projetoAninhado := filepath.Join(diretorio, "GoStock")
		if _, err := os.Stat(filepath.Join(projetoAninhado, "frontend", "html", "index.html")); err == nil {
			return os.Chdir(projetoAninhado)
		}

		parent := filepath.Dir(diretorio)
		if parent == diretorio {
			return fmt.Errorf("frontend/html/index.html não encontrado")
		}
		diretorio = parent
	}
}

func criarSessao(usuarioID int) (string, error) {
	bytes := make([]byte, 32)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	token := hex.EncodeToString(bytes)

	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	expiraEm := time.Now().Add(30 * 24 * time.Hour)

	_, err := database.DB.Exec(`
		INSERT INTO sessoes (usuario_id, token_hash, expira_em)
		VALUES (?, ?, ?)
	`, usuarioID, tokenHash, expiraEm)

	if err != nil {
		return "", err
	}

	return token, nil
}

func usuarioDaSessao(r *http.Request) (int, bool) {
	cookie, err := r.Cookie("sessao")

	if err != nil {
		return 0, false
	}

	hash := sha256.Sum256([]byte(cookie.Value))
	tokenHash := hex.EncodeToString(hash[:])

	var usuarioID int
	var expiraEm time.Time

	err = database.DB.QueryRow(`
		SELECT usuario_id, expira_em
		FROM sessoes
		WHERE token_hash = ?
	`, tokenHash).Scan(&usuarioID, &expiraEm)

	if err != nil {
		return 0, false
	}

	if time.Now().After(expiraEm) {
		database.DB.Exec(`
			DELETE FROM sessoes
			WHERE token_hash = ?
		`, tokenHash)

		return 0, false
	}

	return usuarioID, true
}

func carregarResumo() (ResumoData, error) {
	var resumo ResumoData
	err := database.DB.QueryRow(`SELECT COALESCE(SUM(quantidade), 0) FROM produtos`).Scan(&resumo.TotalEstoque)
	if err != nil {
		return resumo, err
	}
	err = database.DB.QueryRow(`SELECT COALESCE(SUM(quantidade), 0) FROM movimentacoes WHERE tipo = 'SAIDA'`).Scan(&resumo.TotalVendas)
	if err != nil {
		return resumo, err
	}
	err = database.DB.QueryRow(`SELECT COUNT(*) FROM movimentacoes`).Scan(&resumo.TotalMovimentacoes)
	return resumo, err
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
