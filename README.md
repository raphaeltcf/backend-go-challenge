# 🚀 Backend Go Challenge

Serviço backend para processamento de pedidos via fila de eventos, desenvolvido em Go com Clean Architecture.

*******

Tabelas de conteúdo
1. [Pré-requisitos](#prerequisitos)
2. [Colocando para funcionar](#funcionando)
3. [Endpoints](#endpoints)
4. [Features](#features)
5. [Arquitetura](#arquitetura)
6. [Melhorias futuras](#melhorias)
7. [Feito utilizando](#built)

*******

<div id='prerequisitos'/>

## 🚀 Começo

Estas instruções permitirão que você obtenha uma cópia de trabalho do projeto em sua máquina local para fins de desenvolvimento e teste.

### 📋 Pré-requisitos

Antes de começar, você precisará ter as seguintes ferramentas instaladas em sua máquina:
[Git](https://git-scm.com) e
[Go 1.26+](https://golang.org/) ou
[Docker](https://www.docker.com/)


Também é bom ter um editor para trabalhar com o código como [VSCode](https://code.visualstudio.com/)

*******

<div id='funcionando'/>

## 🎲 Colocando para funcionar
```bash
# Clone o repositório
$ git clone https://github.com/raphaeltcf/backend-go-challenge
$ cd backend-go-challenge
```

### Localmente
```bash
# Instala as dependências
$ go mod download

# Roda a aplicação
$ go run cmd/main.go

# Roda os testes
$ go test ./...
```

### Com Docker
```bash
# Build — roda os testes automaticamente antes de compilar
$ docker build -t backend-go-challenge .

# Roda o container
$ docker run -p 8080:8080 backend-go-challenge
```

*******

<div id='endpoints'/>

## 🌐 Endpoints

| Endpoint | Método | Descrição |
|---|---|---|
| /health | GET | Health check do serviço |
| /metrics | GET | Métricas de processamento |

### Exemplo `/metrics`
```json
{
  "failure_rate": "16.67%",
  "total_failed": 0,
  "total_invalid": 2,
  "total_processed": 10
}
```

*******

<div id='features'/>

## ✅ Features

- [x] Consumo de eventos via fila em memória simulando SQS
- [x] Validação de pedidos
- [x] Processamento com retry automático
- [x] Persistência com SQLite
- [x] Logs estruturados em JSON
- [x] Processamento concorrente com worker pool
- [x] Resiliência com retry em caso de falha
- [x] Clean Architecture
- [x] Context propagado entre camadas
- [x] Idempotência — pedidos duplicados não são reprocessados
- [x] Métricas simples com failure rate
- [x] GET /health
- [x] Correlation ID nos logs
- [x] Testes unitários e E2E
- [x] Docker multi-stage com testes automáticos no build

*******

<div id='arquitetura'/>

## 🏛 Arquitetura

O projeto segue **Clean Architecture** com separação em três camadas:

**Domain** — entidades e regras de negócio puras, sem dependência externa.

**Application** — casos de uso e interfaces. O `ProcessOrderUseCase` orquestra validação, idempotência, processamento com retry e persistência.

**Infra** — implementações concretas: SQLite, channel Go como fila, worker pool, slog para logs.

### Fluxo de processamento
```
Gerador → Fila (channel) → Worker Pool → Use Case → SQLite
                                              ↓
                               Validação → Idempotência → Retry → Persistência
```

### Estrutura do projeto
```
backend-go-challenge/
├── cmd/
│   └── main.go                        → entrypoint e injeção de dependências
├── internal/
│   ├── domain/
│   │   └── order.go                   → entidade Order e validação
│   ├── application/
│   │   ├── dto/
│   │   │   ├── order_input.go         → DTO de entrada
│   │   │   └── order_output.go        → DTO de saída
│   │   ├── interfaces.go              → contratos OrderRepository e OrderQueue
│   │   └── process_order.go           → caso de uso principal
│   └── infra/
│       ├── correlation/
│       │   └── correlation.go         → propagação de correlation ID
│       ├── logger/
│       │   └── logger.go              → logger JSON com slog
│       ├── metrics/
│       │   └── metrics.go             → contadores e failure rate
│       ├── queue/
│       │   └── memory_queue.go        → fila em memória e gerador
│       ├── storage/
│       │   ├── connection.go          → conexão SQLite e migration
│       │   └── order_repository.go    → implementação do repositório
│       └── worker/
│           └── pool.go                → worker pool com graceful shutdown
├── Dockerfile
├── go.mod
└── README.md
```

*******

<div id='melhorias'/>

## 🔮 Melhorias futuras

- Endpoints REST para consulta de pedidos
- Métricas de latência por pedido
- Integração real com AWS SQS
- Autenticação nos endpoints HTTP
- Testes de carga para validar comportamento do worker pool sob alta demanda
- Endpoint `POST /orders` para envio manual de pedidos
- Substituir SQLite por um banco de dados de produção como PostgreSQL ou DynamoDB 
*******

<div id='built'/>

## 🛠️ Feito utilizando

<img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/go/go-original.svg" width="40" height="40" /> <img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/docker/docker-original.svg" width="40" height="40" /> <img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/sqlite/sqlite-original.svg" width="40" height="40" />