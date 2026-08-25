# Config

> The composition-root options layer — env-var-bound structs that every assembler reads to wire infrastructure and modules. Infrastructure, not a business bounded context.

## Responsibility
Defines all configuration as nested Go structs tagged with `env` / `envPrefix`, and loads them from the environment via `caarlos0/env`. Each process entrypoint has its own top-level config that embeds the shared config. Config is intentionally the single options layer the composition root depends on — modules receive plain option values, never the config package.

## What it provides
- `SharedConfig` + `LoadShared()` — infra common to every process: `Database`, `Cache`, `Broker`, `EventBus`, `SecretManager`, `FileGateway`, `JWT`, `PlatformUI`, `Platform`, `Webhook`, `Workflows`, `FunctionCalling`, `CredentialOAuth`, `MCP`, `Auth`, `EmailSender`, plus `AppEnv`. `LoadShared` also post-processes file-backed fields (system-instruction files).
- Per-entrypoint configs, each embedding `*SharedConfig` (tagged `env:"-"`) plus its own options and a `Load*()` constructor:
  - `PlatformAPIConfig` / `LoadPlatformAPI` — `HTTPServer` (`PLATFORM_API_`), `WebSocket`, `TaskRunner`.
  - `MCPAPIConfig` / `LoadMCPAPI` — `HTTPServer` (`MCP_API_`), `MCP`.
  - `WebhookAPIConfig` / `LoadWebhookAPI` — `HTTPServer` (`WEBHOOK_API_`), `TaskRunner`.
  - `TaskWorkerConfig` / `LoadTaskWorker` — `TaskRunner`.
- Option leaf structs (one file per concern): `CacheOptions`/`BrokerOptions` (Redis address + pool), `FileGatewayOptions`, `HTTPServerOptions`, `MetricsOptions`, `JWTOptions`, etc., plus `AppEnv` with `IsProduction`/`IsDevelopment`.

## Used by
`bootstrap` exclusively — `LoadShared` feeds `NewCore`; the per-entrypoint `Load*` feeds the matching assembler. Modules consume only the unwrapped option values passed through their `Dependencies`.

## Notes
- Full env-var names are the concatenation of nested `envPrefix` tags plus the leaf `env` tag (e.g. `CACHE_REDIS_ADDRESS`, `PLATFORM_API_READ_TIMEOUT`). Per-entrypoint `HTTPServer` reuses one struct under different prefixes (`PLATFORM_API_`, `MCP_API_`, `WEBHOOK_API_`).
- Cache and broker share the same Redis-options shape but are distinct env trees (`CACHE_REDIS_*` vs `BROKER_REDIS_*`); the broker is the shared realtime/eventbus-wake Redis.
- `.env.example` must be kept in sync with these structs — any `env` / `envPrefix` change requires updating it in the same change.
- Do not "decouple" modules from config; config is the options layer of the composition root and is depended on only there.

## Layout
- `shared.go` — `SharedConfig` + `LoadShared` (the canonical env tree).
- `platform_api.go`, `mcp_api.go`, `webhook_api.go`, `worker.go` — per-entrypoint configs and loaders.
- One file per options group: `cache.go`, `broker.go`, `file_gateway.go`, `http_server.go`, `metrics.go`, `database.go`, `jwt.go`, `auth.go`, `workflows.go`, … ; `env.go` — `AppEnv`.
