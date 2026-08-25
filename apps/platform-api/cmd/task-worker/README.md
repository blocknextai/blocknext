# task-worker

> A background worker that consumes queued workflow tasks from the task queue and executes them; runs only when the task runner is in `queue` mode.

## What it does
Loads the task-worker config and exits early unless `TASK_RUNNER_MODE=queue`. Otherwise it bootstraps the workflow-execution module graph via `bootstrap.NewTaskWorker` and calls `TaskRunnerModule.StartAsWorker`, which pulls tasks from the configured queue (Redis stream) and runs them through the node executors. It wires a realtime broadcaster so execution progress can be pushed to clients, and blocks until shutdown.

## Bootstrap & config
- **Assembler:** `bootstrap.NewCore` + `bootstrap.NewTaskWorker` — wires the full workflow/trigger/execution stack needed to run tasks; exposes only the `TaskRunnerModule` (and broadcaster).
- **Config:** `config.LoadTaskWorker()` → `config.TaskWorkerConfig` (embeds `SharedConfig`, `TaskRunner` under `TASK_RUNNER_`).
- **Runs as:** background worker loop (no HTTP server).

## Bounded contexts activated
- common, account, organizations, credentialoauth
- nodeengine, platform, credentials, llm
- workflows, executions, triggers
- taskrunner

## Notes
- No-op exit if `cfg.TaskRunner.Mode != queue` (in `embedded` mode the platform-api process runs the task runner instead).
- Graceful shutdown (`WaitForShutdown`, `TASK_RUNNER_SHUTDOWN_TIMEOUT`): cancels context, shuts the task runner + broadcaster, then core.
