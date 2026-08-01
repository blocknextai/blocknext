# Task Runner

> Executes workflow tasks as DAGs of nodes — triggering, dispatching, running, retrying, cancelling, and scheduling them.

## Responsibility
This context owns the runtime execution of a task: it resolves a workflow's DAG, dispatches it to a worker pool (embedded or Redis-queue), runs each node, enforces the configured concurrency limit (`TASK_RUNNER_MAX_CONCURRENT_TASKS`), and streams live task/node progress over the realtime broadcaster. It also drives cron-scheduled and webhook-driven executions and recovery of stale/claimed tasks. It does not own workflow definitions, trigger persistence, or execution history — those live in the `workflows`, `triggers`, and `executions` contexts.

## Domain
- **Aggregates / key types:**
  - `task.Task` — an in-flight execution: organization, `ExecutionContext`, `ContextItemID`, the `dag.DAG`, status, node-execution ID map, previous outputs, and `TriggerContext`.
  - `task.TaskTriggerType` — `manual`, `rerun_all`, `rerun_failed`, `schedule`, `webhook`, `api`.
  - `task.TaskEvent` / `node.NodeEvent` — realtime progress payloads (status, error, duration, outputs).
  - `taskqueue.TaskQueue`, `domain/taskrunner/*` — service interfaces (TaskService, TaskDispatcher, WorkerPool, CronService, SemaphoreManager, LeaderRunner, etc.).
- **Key rules / invariants:**
  - `TriggerContext` must be populated for manual/schedule/api/webhook triggers so `$trigger.*` resolves (built from `RuntimePrompt` / adapter payload); passing nil silently breaks resolution.
  - Sentinels: `ErrTaskAlreadyCompleted`, `ErrTaskNotFailed`, `ErrWorkerPoolFull`, `ErrTaskCancelled`, etc.

## Use cases (application)
- **TriggerTask** (command) — start a new task execution for a workflow/context.
- **RerunAll** (command) — re-run all nodes of a prior task.
- **RerunFailed** (command) — re-run only the failed nodes of a task.
- **CancelTask** (command) — cancel an in-flight task.
- **contextresolver** — resolves the DAG (nodes/edges) from the workflow.
- **webhooks.WebhookProcessor** — resolves an inbound webhook to a trigger and launches a `webhook` task.
- **taskrunner/** — execution engine internals: worker pool, dispatcher, executor, lifecycle/coordinator, credential/data processors, cron, recovery, concurrency-limit resolver, event publisher.

## HTTP API
- `POST /organizations/:organizationId/task-runner/trigger` — trigger a task (auth; `TriggerTaskPermission`).
- `POST /organizations/:organizationId/task-runner/rerun-all` — re-run all (auth; `RetryTaskPermission`).
- `POST /organizations/:organizationId/task-runner/rerun-failed` — re-run failed (auth; `RetryTaskPermission`).
- `POST /organizations/:organizationId/task-runner/cancel` — cancel (auth; `CancelTaskPermission`).
- `POST /task-runner/trigger/:workflowId` — trigger via API key (`APIKeyMiddleware` + `ScopeWorkflowsTrigger`).

## Events
- **Published:** `task.TaskEvent` and `node.NodeEvent` via the realtime `Broadcaster` (ephemeral server→client progress, not durable eventbus events).
- **Consumed:** none — schedule-trigger changes are picked up by the recovery pass, which periodically reconciles cron jobs against `TriggerService.GetAllActive`.

## Dependencies
- **Bounded contexts:** `workflows`, `executions` (TaskExecution/TaskClaim/NodeExecution services), `credentialoauth` (token regenerate), `triggers` (TriggerService + WebhookResolver), `llm/functioncalling`, `nodeengine` adapters.
- **Infrastructure:** Redis (task queue, leader election, semaphore), realtime `Broadcaster`, `cache.Service`. No DB tables of its own (state lives in `executions`).

## Layout
Standard DDD layers present, wired in `module.go`.
