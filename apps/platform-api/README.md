# platform-api

A modern, scalable API service built with Go, following Domain-Driven Design (DDD) and Clean Architecture principles. The API provides AI-powered workflow automation capabilities.

## 🚀 Quick Start

### Prerequisites
- Docker, Docker Compose and Make
- Go 1.26.5 (optional, for local development)

### 🐳 Docker

The full stack is orchestrated from the monorepo root — see the root `README.md`.

### ⚙️ Configuration

Copy `.env.example` to `.env` and adjust values before running.

## 🏗️ Architecture

### Module Structure

The codebase follows a modular monolith architecture with clear separation of concerns:

```
internal/
├── [module]/
│   ├── application/     # Use cases, command/query handlers
│   ├── domain/          # Business entities and logic
│   ├── infrastructure/  # External dependencies (DB, APIs)
│   └── presentation/    # HTTP controllers and routing
```

### Modules

`account`, `apikeys`, `billing`, `bootstrap`, `cache`, `common`, `config`, `credentials`, `eventbus`, `executions`, `filegateway`, `library`, `llm`, `marketplace`, `mcp`, `nodeengine`, `notifications`, `oauth`, `organizations`, `payments`, `plans`, `platform`, `quota`, `realtime`, `subscriptions`, `support`, `taskrunner`, `triggers`, `web3`, `webhooks`, `workflows`, `ws`.

### Binaries

The `cmd/` directory contains the entry points for each runnable service:

| Binary | Purpose |
| --- | --- |
| `platform-api` | Main HTTP / WebSocket API server |
| `platform-api-migration` | Standalone migration runner |
| `webhook-api` | Webhook ingestion / processing service |
| `mcp-api` | Model Context Protocol (MCP) server |
| `task-worker` | Redis-based background task consumer |
| `event-relay-worker` | Drains the transactional outbox (eventbus relay) |
| `payment-confirmation-worker` | Polls and confirms pending payments |
| `payment-session-cleanup-worker` | Expires stale payment sessions |
| `recurring-credit-worker` | Grants recurring plan credits |
| `subscription-expiration-worker` | Handles subscription expiry transitions |

## 🛠️ Key Technologies

- **Language**: Go 1.26.5
- **Web Framework**: [Fiber v3](https://github.com/gofiber/fiber)
- **Database**: PostgreSQL 18 via `database/sql` + `lib/pq` (no ORM)
- **Migrations**: [golang-migrate](https://github.com/golang-migrate/migrate)
- **Cache**: Redis ([go-redis/v9](https://github.com/redis/go-redis))
- **Eventing**: Redis-backed two-layer model — `eventbus` (durable transactional outbox, server→server) + `realtime` (ephemeral pub/sub broadcast, server→client)
- **WebSocket**: `gofiber/contrib/websocket`
- **Authentication**: JWT (`golang-jwt/jwt/v5`)
- **Configuration**: `caarlos0/env`
- **Logging**: `log/slog` (standard library)
- **Scheduling**: `robfig/cron`
- **Observability**: OpenTelemetry
- **Blockchain**: `go-ethereum`, IPFS (`go-cid`, `multibase`)
- **AI / LLM**: Gemini and local LLM providers (function-calling, streaming chat)
- **MCP**: Model Context Protocol server (`modelcontextprotocol/go-sdk`)
- **Container**: Docker & Docker Compose
- **Hot Reload**: [Air](https://github.com/air-verse/air) (development)
