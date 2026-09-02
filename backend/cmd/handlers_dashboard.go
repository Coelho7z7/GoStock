package main

import (
	"html/template"
	"log"
	"net/http"

	database "gostock/backend/Database"
	"gostock/backend/models"
	"gostock/backend/services"
)

type DashboardData struct {
	Usuario    *models.Usuario
	Produtos   []models.Produto
	Atividades []models.Movimentacao
	Resumo     ResumoData
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
	if len(atividades) > 8 {
		atividades = atividades[:8]
	}

	dados := DashboardData{
		Usuario:    usuario,
		Produtos:   produtos,
		Atividades: atividades,
		Resumo:     resumo,
	}

	tmpl, err := template.ParseFiles("frontend/html/dashboard.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, dados); err != nil {
		http.Error(w, "Erro ao renderizar página", http.StatusInternalServerError)
		return

	}

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
