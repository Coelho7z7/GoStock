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

## Como executar

```bash
# Clone o repositório
git clone https://github.com/coelho7z7/GoStock.git

# Entre na pasta do projeto
cd GoStock

# Rode o projeto
go run main.go
```

> Ajuste o comando de execução conforme o nome real do arquivo principal e eventuais variáveis de ambiente do projeto.

## Autor

Desenvolvido por [Matheus Henrique Coelho Lopes](https://github.com/coelho7z7).
