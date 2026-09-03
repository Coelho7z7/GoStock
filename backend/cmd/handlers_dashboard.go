package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	database "gostock/backend/Database"
	"gostock/backend/models"
	"gostock/backend/services"
)

type DashboardData struct {
	Usuario    *models.Usuario
	Produtos   []models.Produto
	Atividades []models.Movimentacao
	Grafico    []GraficoDia
	Resumo     ResumoData
}

type GraficoDia struct {
	Data       string
	Valor      float64
	Percentual int
}

type ResumoData struct {
	TotalEstoque       int
	TotalVendas        int
	Faturamento        float64
	TotalMovimentacoes int
}

// handlerDashboard monta a visão geral do sistema.
func handlerDashboard(w http.ResponseWriter, r *http.Request) {
	usuarioID, true := usuarioDaSessao(r)
	if !true {
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
		log.Println("erro em BuscarTodosProdutos (dashboard):", err) // linha temporária
		http.Error(w, "Erro ao buscar produtos", http.StatusInternalServerError)
		return
	}

	resumo, err := carregarResumo()
	if err != nil {
		http.Error(w, "Erro ao carregar resumo", http.StatusInternalServerError)
		return
	}

	atividades, err := services.BuscarMovimentacoesWeb()
	if err != nil {
		http.Error(w, "Erro ao carregar atividades", http.StatusInternalServerError)
		return
	}
	if len(atividades) > 5 {
		atividades = atividades[:5]
	}

	grafico, err := carregarGraficoFaturamento()
	if err != nil {
		http.Error(w, "Erro ao carregar gráfico", http.StatusInternalServerError)
		return
	}

	dados := DashboardData{
		Usuario:    usuario,
		Produtos:   produtos,
		Atividades: atividades,
		Grafico:    grafico,
		Resumo:     resumo,
	}

	tmpl, err := template.New("dashboard.html").Funcs(template.FuncMap{
		"moedaBR": func(valor float64) string {
			texto := strconv.FormatFloat(valor, 'f', 2, 64)
			partes := strings.Split(texto, ".")
			return fmt.Sprintf("%s,%s", partes[0], partes[1])
		},
	}).ParseFiles("frontend/html/dashboard.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, dados); err != nil {
		http.Error(w, "Erro ao renderizar página", http.StatusInternalServerError)
		return

	}

}

func carregarGraficoFaturamento() ([]GraficoDia, error) {
	valores := make(map[string]float64)
	rows, err := database.DB.Query(`
		SELECT date(data), COALESCE(SUM(valor_total), 0)
		FROM vendas
		WHERE date(data) >= date('now', '-6 days')
		GROUP BY date(data)
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var data string
		var valor float64
		if err := rows.Scan(&data, &valor); err != nil {
			return nil, err
		}
		valores[data] = valor
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	maximo := 0.0
	for valor := range valores {
		if valores[valor] > maximo {
			maximo = valores[valor]
		}
	}
	resultado := make([]GraficoDia, 0, 7)
	agora := time.Now()
	for indice := 6; indice >= 0; indice-- {
		data := agora.AddDate(0, 0, -indice)
		chave := data.Format("2006-01-02")
		valor := valores[chave]
		percentual := 0
		if maximo > 0 {
			percentual = int(valor / maximo * 100)
		}
		if valor > 0 && percentual < 8 {
			percentual = 8
		}
		resultado = append(resultado, GraficoDia{Data: data.Format("02/01"), Valor: valor, Percentual: percentual})
	}
	return resultado, nil
}

// carregarResumo calcula os números exibidos nos cartões do dashboard.
//
// TotalVendas soma as unidades vendidas a partir da tabela "vendas"
// (vendas de fato, com preço e forma de pagamento) — e não da tabela
// "movimentacoes", que também registra saídas manuais de estoque
// (ex.: produto danificado) que não são vendas.
func carregarResumo() (ResumoData, error) {
	var resumo ResumoData
	err := database.DB.QueryRow(`SELECT COALESCE(SUM(quantidade), 0) FROM produtos`).Scan(&resumo.TotalEstoque)
	if err != nil {
		return resumo, err
	}
	err = database.DB.QueryRow(`SELECT COALESCE(SUM(quantidade), 0) FROM vendas`).Scan(&resumo.TotalVendas)
	if err != nil {
		return resumo, err
	}
	err = database.DB.QueryRow(`SELECT COALESCE(SUM(valor_total), 0) FROM vendas`).Scan(&resumo.Faturamento)
	if err != nil {
		return resumo, err
	}
	err = database.DB.QueryRow(`SELECT COUNT(*) FROM movimentacoes`).Scan(&resumo.TotalMovimentacoes)
	return resumo, err
}
