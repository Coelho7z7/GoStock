# GoStock: Sistema de Controle de Estoque em Go

GoStock é um sistema de controle de estoque desenvolvido para simplificar o gerenciamento de produtos, entradas, saídas e movimentações, unindo um backend em Go a uma interface web leve e direta.

## Overview

Pequenos negócios e times operacionais frequentemente perdem tempo e precisão controlando estoque em planilhas ou processos manuais. O GoStock resolve isso oferecendo um sistema centralizado onde é possível cadastrar produtos, acompanhar quantidades, registrar saídas, realizar vendas e visualizar o histórico de movimentações em um só lugar.

### Key Features

* **Dashboard:** visão geral e centralizada do estoque

* **Controle de Estoque:** adição e controle de quantidade por produto

* **Registro de Saídas:** controle das saídas de produtos do estoque

* **Vendas e PDV:** registro de vendas e gerenciamento das operações do PDV

* **Movimentações:** histórico de entradas, saídas e alterações

* **Cadastro e Edição:** cadastro de novos produtos e alteração dos já existentes

* **Login e Permissões:** autenticação de usuários e controle de acesso

## Architecture

O GoStock é organizado em módulos:

1. **Backend (Go):** lógica de negócio, rotas, autenticação e regras do sistema

2. **Frontend (HTML/CSS/JS):** interface web e interações do sistema

3. **Banco de Dados (SQLite):** persistência dos dados do sistema

## Requirements

* Go 1.26+

* SQLite

* Navegador web atualizado para acessar a interface

## Usuários de teste

Estas credenciais existem apenas para desenvolvimento local:

* **Administrador:** `***` / `***` — role `admin`

* **Usuário:** `usuario@gmail.com` / `usuario123` — role `usuario`

## Como executar

```bash
go run ./backend/cmd
```

Depois, acesse `http://localhost:8080` no navegador.

Para usar outra porta local, defina a variável `PORT` antes de iniciar o servidor.

No Railway, a aplicação utiliza automaticamente a porta fornecida pela variável `PORT` do ambiente.

### Rotas principais

* `/` — Login

* `/login` — Autenticação

* `/logout` — Encerramento da sessão

* `/dashboard` — Dashboard

* `/produtos` — Cadastro e listagem de produtos

* `/alterar-produto` — Alteração de produtos

* `/estoque` — Controle de estoque

* `/vendas` — PDV e vendas

* `/movimentacoes` — Histórico de movimentações

## Autor

Desenvolvido por [Matheus Henrique Coelho Lopes](https://github.com/coelho7z7).
