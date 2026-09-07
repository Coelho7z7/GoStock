package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	database "gostock/backend/database"
	"gostock/backend/models"
	"gostock/backend/services"
)

type DashboardData struct {
	User       *models.User
	Products   []models.Product
	Activities []models.Movement
	Chart      []ChartDay
	Summary    SummaryData
}

type ChartDay struct {
	Date       string
	Value      float64
	Percentage int
}

type SummaryData struct {
	TotalStock    int
	TotalSales    int
	Revenue       float64
	TotalMovements int
}

// dashboardHandler monta a visão geral do sistema.
func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	userID, authenticated := userFromSession(r)
	if !authenticated {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	user, err := services.GetUserByID(userID)
	if err != nil {
		http.Error(w, "Usuário não encontrado", http.StatusInternalServerError)
		return
	}

	products, err := services.GetAllProducts()
	if err != nil {
		log.Println("erro em GetAllProducts (dashboard):", err) // linha temporária
		http.Error(w, "Erro ao buscar produtos", http.StatusInternalServerError)
		return
	}

	summary, err := loadSummary()
	if err != nil {
		http.Error(w, "Erro ao carregar resumo", http.StatusInternalServerError)
		return
	}

	activities, err := services.GetMovementsWeb()
	if err != nil {
		http.Error(w, "Erro ao carregar atividades", http.StatusInternalServerError)
		return
	}
	if len(activities) > 5 {
		activities = activities[:5]
	}

	chart, err := loadRevenueChart()
	if err != nil {
		http.Error(w, "Erro ao carregar gráfico", http.StatusInternalServerError)
		return
	}

	data := DashboardData{
		User:       user,
		Products:   products,
		Activities: activities,
		Chart:      chart,
		Summary:    summary,
	}

	tmpl, err := template.New("dashboard.html").Funcs(template.FuncMap{
		"formatBRL": func(value float64) string {
			text := strconv.FormatFloat(value, 'f', 2, 64)
			parts := strings.Split(text, ".")
			return fmt.Sprintf("%s,%s", parts[0], parts[1])
		},
	}).ParseFiles("frontend/html/dashboard.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Erro ao renderizar página", http.StatusInternalServerError)
		return

	}

}

func loadRevenueChart() ([]ChartDay, error) {
	values := make(map[string]float64)
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
		var date string
		var value float64
		if err := rows.Scan(&date, &value); err != nil {
			return nil, err
		}
		values[date] = value
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	max := 0.0
	for value := range values {
		if values[value] > max {
			max = values[value]
		}
	}
	result := make([]ChartDay, 0, 7)
	now := time.Now()
	for i := 6; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		key := date.Format("2006-01-02")
		value := values[key]
		percentage := 0
		if max > 0 {
			percentage = int(value / max * 100)
		}
		if value > 0 && percentage < 8 {
			percentage = 8
		}
		result = append(result, ChartDay{Date: date.Format("02/01"), Value: value, Percentage: percentage})
	}
	return result, nil
}

// loadSummary calcula os números exibidos nos cartões do dashboard.
//
// TotalSales soma as unidades vendidas a partir da tabela "vendas"
// (vendas de fato, com preço e forma de pagamento) — e não da tabela
// "movimentacoes", que também registra saídas manuais de estoque
// (ex.: produto danificado) que não são vendas.
func loadSummary() (SummaryData, error) {
	var summary SummaryData
	err := database.DB.QueryRow(`SELECT COALESCE(SUM(quantidade), 0) FROM produtos`).Scan(&summary.TotalStock)
	if err != nil {
		return summary, err
	}
	err = database.DB.QueryRow(`SELECT COALESCE(SUM(quantidade), 0) FROM vendas`).Scan(&summary.TotalSales)
	if err != nil {
		return summary, err
	}
	err = database.DB.QueryRow(`SELECT COALESCE(SUM(valor_total), 0) FROM vendas`).Scan(&summary.Revenue)
	if err != nil {
		return summary, err
	}
	err = database.DB.QueryRow(`SELECT COUNT(*) FROM movimentacoes`).Scan(&summary.TotalMovements)
	return summary, err
}
