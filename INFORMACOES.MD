// GoStock: sistema de controle de estoque e PDV.
//
// O código do servidor web fica dividido por assunto:
//   - project.go              -> localizar a raiz do projeto
//   - session.go              -> criação/leitura de sessão de login
//   - handlers_auth.go        -> "/", "/login", "/logout"
//   - handlers_dashboard.go   -> "/dashboard"
//   - handlers_produtos.go    -> "/produtos", "/alterar-produto"
//   - handlers_estoque.go     -> "/estoque"
//   - handlers_vendas.go      -> "/vendas", "/api/vendas"
//   - handlers_movimentacoes.go -> "/movimentacoes"
//   - cli.go                  -> menu de terminal (uso opcional, em paralelo ao servidor)