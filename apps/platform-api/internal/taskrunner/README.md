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

## Scheduling & execution order

**Task admission.** In `queue` mode, triggered tasks are enqueued to a Redis Stream and consumed **FIFO** by a consumer group (`TASK_RUNNER_QUEUE_REDIS_PREFETCH_COUNT` messages per fetch) — there are no priority classes between tasks. In `embedded` mode the platform-api process dispatches directly. Either way, admission is gated by a per-organization Redis **semaphore**: at most `TASK_RUNNER_MAX_CONCURRENT_TASKS` tasks run concurrently per org. In-process parallelism is bounded by the worker pool (`TASK_RUNNER_WORKER_POOL_SIZE` workers, `TASK_RUNNER_WORKER_POOL_QUEUE_SIZE` backlog; saturation yields `ErrWorkerPoolFull`).

**Node order within a task.** Execution is **dependency-driven**, not list-driven: a node becomes runnable once all of its upstream nodes have finished, and independent branches run **concurrently** (one goroutine per ready node). Disabled nodes are marked `skipped` and count as finished for their dependents. On `rerun_failed`, previously successful nodes are skipped and their stored outputs are reused.

**Deterministic ordering.** At DAG construction, a topological sort (Kahn's algorithm over a min-heap) validates acyclicity and fixes a canonical order: among simultaneously-ready nodes, priority follows node type — `system_starter` (0) → `system_condition` (1) → `system_action` (2) → provider nodes (3) — with node-ID as the tie-breaker.

**Schedules & recovery.** Cron-scheduled triggers fire only on the elected **leader** instance (Redis leader election, `TASK_RUNNER_LEADER_*`). A periodic recovery pass (`TASK_RUNNER_RECOVERY_INTERVAL`) reconciles cron jobs against active triggers and reclaims queue messages whose claim went stale (`TASK_RUNNER_QUEUE_STALE_CLAIM_TIMEOUT`), retrying up to `TASK_RUNNER_QUEUE_MAX_RETRIES` times.

## Prompt layering (function calling)

When function calling is enabled (`FUNCTION_CALLING_ENABLED`) and a node carries natural-language text, the `NodeExecutor` asks the LLM (`llm/functioncalling`) to resolve the node's input parameters. The prompt is assembled from **four layers, in a fixed order**:

| # | Layer | Source | Sent as |
| - | - | - | - |
| 1 | **System instruction** | platform-owned, embedded from `internal/prompts` and loaded once; Gemini provider serves it from content cache | LLM `systemInstruction` role |
| 2 | **Node instruction** (`instruction`) | Authored on the node at design time in the flow editor | `NODE INSTRUCTION:` block in the user turn |
| 3 | **Runtime instruction** (`runtimeInstruction`) | Supplied per node at run time (trigger `RuntimeConfig.Nodes` overlay) | `RUNTIME INSTRUCTION:` block in the user turn |
| 4 | **Runtime prompt** (`runtimePrompt`) | The trigger's free-text prompt (`RuntimeConfig.RuntimePrompt`; also feeds `TriggerContext` for `$trigger.*` references) | `RUNTIME PROMPT (reference only):` block in the user turn |

Semantics of the ordering:
- The platform-owned system instruction defines the resolution contract and always wins — user-supplied layers cannot rewrite it (they live in the user turn, not the system role).
- Within the user turn the blocks are concatenated exactly in the order above: the design-time node instruction sets the node's job, the runtime instruction refines it for this run, and the runtime prompt is explicitly labeled **reference only** — it provides context to draw values from, not directives.
- **Explicit configuration beats LLM inference:** after function calling returns, the node's configured parameters (with `$references` resolved) are merged **over** the LLM-resolved values — any non-empty configured value overrides what the model inferred. Only gaps are filled by the LLM.

The same layering is implemented by both providers (`gemini`, `localllm`); function-calling outputs are also persisted alongside node executions for inspection.

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
