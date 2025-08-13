# Golang SBA com Hazelcast — Projeto de Treino

Este repositório é um projeto de treino que explora Space-Based Architecture (SBA) usando Hazelcast como middleware distribuído e Go (Golang) para o backend. O objetivo é praticar padrões de alta concorrência, escalabilidade horizontal e processamento de eventos em tempo real.

## Visão Geral
- **Stack**: Go 1.22+, Hazelcast 5.x, Docker/Docker Compose, Kubernetes (manifests), Swagger/OpenAPI.
- **Arquitetura**: Estilo SBA com cluster Hazelcast; serviço HTTP em Go expondo APIs e documentação.
- **Estado**: Em desenvolvimento — ver `plan.md` para o roadmap detalhado.

## Como Executar
### 1) Subir o cluster Hazelcast (Docker Compose)
```bash
docker compose -f deployments/docker-compose.yml up -d
```

### 2) Rodar a aplicação em Go
```bash
go run ./cmd/auction
```

Por padrão o serviço sobe em `:8080` e tenta conectar em `127.0.0.1:5701`.

## Documentação da API (Swagger)
- UI: `http://localhost:8080/swagger/`
- JSON: `http://localhost:8080/swagger/doc.json`

A UI e o JSON são servidos pela própria aplicação (sem serviço externo de swagger-ui).

## Variáveis de Ambiente
- `APP_HTTP_ADDR` (default `:8080`): endereço HTTP do servidor.
- `HZ_CLUSTERNAME` (default `auction-cluster`): nome do cluster Hazelcast.
- `HZ_MEMBERS` (default `127.0.0.1:5701`): lista de membros Hazelcast, separada por vírgula.
- `HZ_SKIP` (opcional, `1` para usar fake in-memory durante desenvolvimento/teste).

Exemplo:
```bash
APP_HTTP_ADDR=":8080" HZ_CLUSTERNAME="auction-cluster" HZ_MEMBERS="127.0.0.1:5701" go run ./cmd/auction
```

## Endpoints Básicos
- `GET /health` — healthcheck simples.
- `GET /swagger/` — UI do Swagger.
- `GET /swagger/doc.json` — especificação OpenAPI gerada.

## Estrutura do Projeto (POD)
- `cmd/auction` — binário/API HTTP e handlers.
- `internal/auction` — domínio de leilões (serviços, modelos, testes).
- `internal/platform/hazelcast` — cliente/abstrações de Hazelcast.
- `internal/docs` — integração com Swagger/OpenAPI.
- `deployments` — Docker, Docker Compose e manifests de Kubernetes.

## Próximos Passos
Consulte o plano em [`plan.md`](./plan.md) para backlog, fases e metas (inclui testes, observabilidade e segurança).

## Referências
- Hazelcast: [`https://hazelcast.com/`](https://hazelcast.com/)
- Go: [`https://go.dev/`](https://go.dev/)


