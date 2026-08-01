# `cmd/` — Executable Entrypoints

Each subdirectory is a thin `main` package compiled to one binary. The `main.go` files only load config and call a `bootstrap` assembler; all dependency wiring lives in [`internal/bootstrap`](../internal/bootstrap/README.md). Each binary activates a subset of the bounded contexts under [`internal/`](../internal/README.md).

## Processes
| Binary | Type | Purpose |
| --- | --- | --- |
| [`platform-api`](./platform-api/README.md) | HTTP API | Main end-user API; registers every module's routes; optionally runs the embedded task runner + eventbus relay. |
| [`mcp-api`](./mcp-api/README.md) | HTTP API | Model Context Protocol server; exposes nodes as MCP tools, API-key auth. |
| [`webhook-api`](./webhook-api/README.md) | HTTP API | Inbound webhook edge; dispatches workflow-trigger webhooks. |
| [`task-worker`](./task-worker/README.md) | Worker | Consumes queued workflow tasks (only in `queue` task-runner mode). |
| [`event-relay-worker`](./event-relay-worker/README.md) | Worker | Drains the eventbus transactional outbox out-of-process. |
| [`platform-api-migration`](./platform-api-migration/README.md) | Migration | One-shot per-module DB migration runner. |

## Shared wiring
Every binary (except the migration CLI) builds from [`internal/bootstrap`](../internal/bootstrap/README.md): `bootstrap.NewCore` provides the shared infrastructure (DB, cache, secret manager, file gateway, eventbus), and a per-entrypoint assembler (`NewPlatformAPI`, `NewMCPAPI`, `NewWebhookAPI`, `NewTaskWorker`) wires the contexts that process needs. Config is loaded via [`internal/config`](../internal/config/README.md) (`SharedConfig` + per-binary structs). The event model is two-layer: a durable server→server eventbus (relayed embedded in `platform-api` or by `event-relay-worker`) and an ephemeral server→client realtime broadcaster.

## See also
- [`internal/README.md`](../internal/README.md) — bounded context map.
