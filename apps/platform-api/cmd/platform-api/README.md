# platform-api

> The main HTTP API server: a Fiber app exposing every public/end-user bounded context, optionally running the task runner and event-bus relay in-process.

## What it does
Loads the platform configuration, bootstraps `Core` (DB, cache, secret manager, file gateway, event bus) and the full module graph via `bootstrap.NewPlatformAPI`, then builds a Fiber app and registers every module's routes behind JWT auth, API-key, and cache middleware. It serves liveness/readiness probes, applies a Redis-backed rate limiter, and runs Prometheus metrics on a side server. When `TASK_RUNNER_MODE=embedded` it also starts the task runner (`StartAsMain`) and the eventbus transactional-outbox relay in the same process.

## Bootstrap & config
- **Assembler:** `bootstrap.NewCore` + `bootstrap.NewPlatformAPI` — wires the complete module graph (all bounded contexts) plus JWT service and realtime broadcaster.
- **Config:** `config.LoadPlatformAPI()` → `config.PlatformAPIConfig` (embeds `SharedConfig`, `HTTPServer`, `TaskRunner`).
- **Runs as:** HTTP server on `HTTP_SERVER_*` address (Fiber); metrics on a separate `HTTP_SERVER_METRICS_*` server.

## Bounded contexts activated
- common, account, organizations, web3, oauth, plans, payments, subscriptions, quota, billing
- nodeengine, platform, credentials, llm
- workflows, marketplace, library, executions, triggers
- taskrunner, apikeys, support, notifications, ws, webhooks

## Notes
- Registers routes for every module; `TaskRunnerModule` and `APIKeysModule` get the API-key middleware in addition to auth.
- `StartAsMain` always runs; the eventbus relay (`core.EventBus.StartRelay`) runs only in `embedded` task-runner mode.
- Graceful shutdown via `bootstrap.ListenAndWait`: cancels app context, then shuts Fiber, app modules, metrics, and core in order.
