GoStock

BACKEND

cmd/
- main.go — servidor
- project.go — projeto
- auth_handler.go — login
- authorization.go — permissões
- session.go — sessões
- dashboard_handler.go — dashboard
- product_handler.go — produtos
- stock_handler.go — estoque
- sale_handler.go — vendas
- movement_handler.go — movimentações
- user_handler.go — usuários
- cli.go — terminal

services/
- product_create.go — cadastro de produtos
- product_query.go — consultas de produtos
- stock_service.go — estoque
- movement_service.go — movimentações
- sale_service.go — vendas
- user_create.go — cadastro de usuários
- user_auth.go — autenticação
- user_query.go — consultas de usuários
- user_management.go — gerenciamento de usuários
- user_seed.go — usuários iniciais

models/
- product.go
- movement.go
- user.go

database/
- database.go

utils/
- input.go
- validation.go


FRONTEND

html/
- index.html — login
- dashboard.html — dashboard
- products.html — produtos
- edit_product.html — editar produto
- stock.html — estoque
- sales.html — vendas
- movements.html — movimentações
- users.html — usuários
- access_denied.html — acesso negado

css/
- base.css
- dashboard.css
- pos.css
- theme.css
- demais arquivos de estilo

js/
- app.js


BANCO

backend/data/gostock.db

SQLite.


TESTES

backend/cmd/main_test.go
backend/services/sale_service_test.go

Executar:

go test ./...