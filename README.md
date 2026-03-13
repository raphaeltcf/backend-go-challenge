# 🚀 Backend Go Challenge

Serviço backend para processamento de pedidos via fila de eventos, desenvolvido em Go com Clean Architecture.

*******

Tabelas de conteúdo
1. [Pré-requisitos](#prerequisitos)
2. [Colocando para funcionar](#funcionando)
3. [Endpoints](#endpoints)
4. [Como simular pedidos](#simulando)
5. [Exemplo de logs](#logs)
6. [Testes](#testes)
7. [Features](#features)
8. [Arquitetura](#arquitetura)
9. [Melhorias futuras](#melhorias)
10. [Feito utilizando](#built)

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

<div id='simulando'/>

## 📦 Como simular pedidos

O serviço possui um gerador automático de pedidos que inicia junto com a aplicação. A cada **500ms** um novo pedido é gerado e enviado para a fila de processamento. A cada **5 pedidos**, 1 pedido inválido é gerado automaticamente para simular cenários de erro.

Você pode acompanhar o processamento em tempo real pelos logs:
```bash
go run cmd/main.go
```

E verificar as métricas acumuladas pelo endpoint:
```bash
curl http://localhost:8080/metrics
```

Para verificar se o serviço está saudável:
```bash
curl http://localhost:8080/health
```

*******

<div id='logs'/>

### 📝 Exemplo de logs
```json
{"time":"2026-03-13T12:38:13.209008365-03:00","level":"INFO","msg":"order enqueued","order_id":"order-1773416293208980288","correlation_id":"corr-1773416293208980288"}
{"time":"2026-03-13T12:38:13.20910405-03:00","level":"INFO","msg":"receiving order","order_id":"order-1773416293208980288","status":"pending","correlation_id":"corr-1773416293208980288"}
{"time":"2026-03-13T12:38:13.209154546-03:00","level":"INFO","msg":"order processed successfully","order_id":"order-1773416293208980288","status":"processed","correlation_id":"corr-1773416293208980288"}
```

Pedido inválido:
```json
{"time":"2026-03-13T12:38:13.209008365-03:00","level":"INFO","msg":"generating invalid order","correlation_id":"corr-1773416293208980288"}
{"time":"2026-03-13T12:38:13.209069617-03:00","level":"INFO","msg":"order enqueued","order_id":"","correlation_id":"corr-1773416293208980288"}
{"time":"2026-03-13T12:38:13.209154546-03:00","level":"ERROR","msg":"invalid order","order_id":"","status":"failed","correlation_id":"corr-1773416293208980288","error":"order amount must be greater than zero"}
```

Graceful shutdown:
```json
{"time":"2026-03-13T12:38:19.545228888-03:00","level":"INFO","msg":"Shutting down order processing system"}
{"time":"2026-03-13T12:38:19.545298591-03:00","level":"INFO","msg":"worker stopped","worker_id":0}
{"time":"2026-03-13T12:38:19.545315633-03:00","level":"INFO","msg":"worker stopped","worker_id":4}
```

Métricas:
```bash
$ curl http://localhost:8080/metrics
{"failure_rate":"16.67%","total_failed":0,"total_invalid":2,"total_processed":10}

$ curl http://localhost:8080/metrics
{"failure_rate":"15.38%","total_failed":0,"total_invalid":2,"total_processed":11}

$ curl http://localhost:8080/metrics
{"failure_rate":"14.29%","total_failed":0,"total_invalid":2,"total_processed":12}
```

Health check:
```bash
$ curl http://localhost:8080/health
{"status":"ok"}
```


*******
<div id='testes'/>
  
## 🧪 Testes

### Rodando os testes
```bash
# Todos os testes
$ go test ./...

# Com verbose para ver cada teste
$ go test ./... -v

# Só os testes unitários
$ go test ./internal/application/...

# Só o teste E2E
$ go test ./internal/infra/...
```

### Testes unitários

Testam o use case `ProcessOrderUseCase` de forma isolada usando um mock do repositório — sem banco de dados real.

Cenários cobertos:
- ✅ Pedido válido deve retornar `status: processed`
- ✅ Pedido com `order_id` vazio deve retornar erro
- ✅ Pedido com `amount` zero deve retornar erro
- ✅ Pedido com `amount` negativo deve retornar erro
- ✅ Pedido com `user_id` vazio deve retornar erro

### Teste E2E

Testa o fluxo completo da aplicação usando SQLite em memória (`:memory:`) — sem deixar arquivos no disco.

Fluxo testado:
```
Enqueue → Worker Pool → Use Case → SQLite → FindByID → status: processed
```

### Testes no Docker

O Dockerfile roda os testes automaticamente antes de compilar. Se algum teste falhar, o build para e o container não é criado:
```bash
$ docker build -t backend-go-challenge .
# [tester] RUN go test ./...   ← testes rodam aqui
# [builder] RUN go build ...   ← só compila se os testes passarem
```


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
