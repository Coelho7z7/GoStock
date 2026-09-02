# GoStock: Sistema de Controle de Estoque em Go

GoStock é um sistema de controle de estoque desenvolvido para simplificar o gerenciamento de produtos, entradas, saídas e movimentações, unindo um backend em Go a uma interface web leve e direta.

## Overview

Pequenos negócios e times operacionais frequentemente perdem tempo e precisão controlando estoque em planilhas ou processos manuais. O GoStock resolve isso oferecendo um sistema centralizado onde é possível cadastrar produtos, acompanhar quantidades, registrar saídas e visualizar o histórico de movimentações em um só lugar.

### Key Features

* **Dashboard:** visão geral e centralizada do estoque
* **Controle de Estoque:** adição de quantidade por produto de forma simples
* **Registro de Saídas:** controle das saídas de produtos do estoque
* **Movimentações:** histórico completo de entradas e saídas
* **Cadastro e Edição:** cadastro de novos produtos e alteração dos já existentes
* **Login seguro:** autenticação com opção de mostrar/ocultar senha

## Architecture

O GoStock é organizado em módulos:

1. **Backend (Go):** lógica de negócio, rotas e regras do sistema
2. **Frontend (HTML/CSS/JS):** interface web, com JS usado principalmente na tela de login (exibir/ocultar senha)
3. **Banco de Dados (SQLite):** persistência leve e local dos dados de estoque

## Requirements

* Go 1.20+ (ajustar conforme a versão usada no projeto)
* SQLite
* Navegador web atualizado para acessar a interface

## Usuários de teste

Estas credenciais existem apenas para desenvolvimento local:

* **Administrador:** `admin@gmail.com` / `Admin123` — role `admin`
* **Usuário:** `usuario@gmail.com` / `usuario123` — role `Usuario`

## Como executar

```bash
# Na raiz do projeto, inicie o servidor HTTP
go run ./backend/cmd
```

Depois, acesse `http://localhost:8080` no navegador. Para usar outra porta local, defina a variável `PORT` antes de iniciar o servidor.

No Railway, a aplicação usa automaticamente a porta fornecida pela variável `PORT` do ambiente.

### Rotas principais

* `/` — Login
* `/login` — Autenticação
* `/cadastro` — Criação de conta
* `/dashboard` — Dashboard
* `/produtos` — Cadastro e listagem de produtos
* `/alterar-produto` — Alteração de produtos
* `/estoque` — Entradas e saídas de estoque
* `/vendas` — PDV e vendas
* `/movimentacoes` — Histórico de movimentações

## Autor

Desenvolvido por [Matheus Henrique Coelho Lopes](https://github.com/coelho7z7).
