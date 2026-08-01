# event-relay-worker

> A standalone background worker that drains the eventbus transactional outbox, relaying durably-stored domain events to subscribers.

## What it does
Loads the platform config and bootstraps `Core` plus the full platform module graph via `bootstrap.NewPlatformAPI` (so every event subscriber is registered), then calls `core.EventBus.StartRelay(ctx)` to run the SKIP-LOCKED outbox relay loop. It does not serve HTTP or register routes; it exists to dispatch the transactional outbox out-of-process, an alternative to the embedded relay that `platform-api` runs in `embedded` task-runner mode.

## Bootstrap & config
- **Assembler:** `bootstrap.NewCore` + `bootstrap.NewPlatformAPI` — wires the full module graph so all event subscribers exist when the relay dispatches; only the eventbus relay is started.
- **Config:** `config.LoadPlatformAPI()` → `config.PlatformAPIConfig`.
- **Runs as:** background worker loop (no HTTP server).

## Bounded contexts activated
- The complete platform-api module graph (same set `bootstrap.NewPlatformAPI` builds), registered so their eventbus inbox subscribers are wired; only `core.EventBus.StartRelay` is invoked.

## Notes
- Two-layer event model: this worker drives the durable server→server eventbus relay; ephemeral server→client realtime broadcast is separate.
- Shutdown (`WaitForShutdown`, `TASK_RUNNER_SHUTDOWN_TIMEOUT`): cancels the relay context, closes the broadcaster, then core.
