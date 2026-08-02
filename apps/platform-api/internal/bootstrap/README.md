# Bootstrap

> Composition root — wires shared infrastructure plus the DDD modules into per-entrypoint assemblies. Infrastructure, not a business bounded context.

## Responsibility
Owns the application's object graph. `NewCore` builds the shared infrastructure once (DB, cache, transaction manager, secret manager, file gateway, event bus); each entrypoint assembler then instantiates exactly the modules it needs against that `Core`. Also provides the Fiber app factory and graceful shutdown plumbing (the Prometheus metrics server wiring is present but commented out). No business logic lives here — it is pure wiring driven by `internal/config`.

## What it provides
- `Core` / `NewCore(cfg *config.SharedConfig, ...Option)` — shared infra: `DB`, `CacheService`, `TransactionManager`, `SecretManager`, `FileGateway`, `EventBus`; `Health` / `Shutdown`. `WithDBPool` overrides connection-pool sizing.
- `NewPlatformAPI(core, *config.PlatformAPIConfig)` — full HTTP API: every module (account, workflows, taskrunner, ws, webhooks, …) plus JWT service and realtime `Broadcaster`.
- `NewMCPAPI(core, *config.MCPAPIConfig)` — MCP server surface (account/credentials/nodeengine + `mcp.Module`).
- `NewWebhookAPI(core, *config.WebhookAPIConfig)` — webhook ingestion (trigger webhook processor, taskrunner).
- `NewTaskWorker(core, *config.TaskWorkerConfig)` — task-execution worker (taskrunner + broadcaster, no HTTP modules).
- `NewFiber(appName, config.HTTPServerOptions, isProduction)` — Fiber app with requestid/logging/CORS/helmet/metrics/recovery middleware; `ListenAndWait` runs the server and blocks on shutdown.
- `StartMetrics(...)` — separate Fiber app exposing `/metrics`.
- `WaitForShutdown(timeout, ...fn)` — SIGINT/SIGTERM handler running shutdown steps with a deadline.

## Used by
The `cmd/*` entrypoints (`platform-api`, `mcp-api`, `webhook-api`, `task-worker`, `event-relay-worker`). Each `main.go` loads its config, calls `NewCore`, then the matching assembler, registers routes/runs the worker, and defers `Shutdown`.

## Notes
- Each entrypoint re-declares its own module graph; module construction is duplicated across assemblers by design (only the needed slice is built per process).
- `Core.Shutdown` closes cache then DB; assembler `Shutdown` closes module-owned resources (taskrunner, broadcaster, ws) first.

## Layout
- `core.go` — shared infra assembly + health/shutdown.
- `platform_api.go`, `mcp_api.go`, `webhook_api.go`, `task_worker.go` — server/worker assemblers.
- `fiber.go` — Fiber factory + `ListenAndWait`; `metrics.go` — metrics server (commented out); `shutdown.go` — signal-driven `WaitForShutdown`.
