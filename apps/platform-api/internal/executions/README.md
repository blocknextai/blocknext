# Executions

> Owns the record of workflow task runs, their per-node executions, the worker claim/lease that drives distributed processing, and the log of single tools invoked directly outside any workflow.

## Responsibility
This bounded context persists and reads the history of workflow executions: each task run, the individual node executions within it, and the claim record workers use to lease a task. It owns task/node execution lifecycle state (status, timing, inputs/outputs, errors) and the distributed claiming protocol (claim, heartbeat, release, stale reclaim, retry counting). It does not schedule or run nodes — execution engines/workers call its services to create and update records and to lease work.

## Domain
- **Aggregates / key types:**
  - `TaskExecution` — a single workflow run: organization scope, optional `TriggeredByUserID` / `FlowTriggerID`, `ExecutionContext`, `ContextItemID`, `Status`, `ExecutionType`, the `dag.Node`/`dag.Edge` snapshot, and start/complete timing.
  - `NodeExecution` — one node's run within a task: `NodeType`, `NodeID`, `Status`, `Inputs`/`Outputs`/`FunctionCallingOutputs`, error, timing.
  - `TaskClaim` — worker lease over a task: `ClaimedBy`, `ClaimedAt`, `RetryCount`; supports `Claim`/`Release`/`Heartbeat`/`ForceRelease`/retry counting with staleness checks.
  - `ExecutionType` — enum `manual` | `schedule` | `webhook` | `api`.
  - `ToolInvocation` — one node run directly, outside any task: organization scope, `Source` (which surface invoked it, `mcp` today), `ToolID`, `Status`, `Parameters`, credential references, `Outputs`, error, start/complete timing. Write-once; no claim, no DAG, no rerun.
  - `Source` — enum `mcp`. A row belongs here only when exactly one node ran, no task/DAG was involved, and an external or interactive caller initiated it; anything else is a `TaskExecution`.
- **Key rules / invariants:** a claim can only be claimed by one worker; release/heartbeat require ownership; retry count is non-negative; node executions require task ID, node type, node ID, and status. Creating a `TaskExecution` also creates its `TaskClaim` in the same transaction. A `ToolInvocation` requires organization ID, a valid source, a tool ID and a valid status.

## Use cases (application)
- **GetAllTaskExecutions** — paginated/searchable list of an organization's executions.
- **GetTaskExecutionByID** — fetch a single execution (organization-scoped).
- **DeleteTaskExecution** — soft-delete one execution.
- **BulkDeleteTaskExecutions** — soft-delete multiple executions.

Service interfaces (used by workers/engines, not HTTP): `TaskExecutionService` (create/update/query by id/status), `NodeExecutionService`, `TaskClaimService` (claim, heartbeat, release, release-stale, retry counting), `ToolInvocationService` (`Record` — one insert, consumed by `mcp`).

- **GetAllToolInvocations** — paginated/searchable list of an organization's direct tool runs.
- **GetToolInvocationByID** — fetch one, with parameters and outputs.

## HTTP API
All under `/organizations/:organizationId/executions`, behind `Authenticate()` plus organization RBAC:

- `GET    /executions` — list executions (ReadTaskExecution permission).
- `GET    /executions/:executionId` — get one execution (ReadTaskExecution permission).
- `DELETE /executions/bulk` — bulk delete (DeleteTaskExecution permission).
- `DELETE /executions/:executionId` — delete one (DeleteTaskExecution permission).
- `GET    /tool-invocations` — list direct tool runs (ReadTaskExecution permission).
- `GET    /tool-invocations/:toolInvocationId` — get one (ReadTaskExecution permission).

## Dependencies
- **Bounded contexts:** `organizations` (organization-user resolution), `workflows`, `account` (users, linked accounts) — consumed via service interfaces for enrichment/resolution.
- **Infrastructure:** Postgres schema `executions` with tables `task_executions`, `node_executions`, `task_claims` (unique per task execution); `TransactionManager`.

## Layout
Standard DDD layers present, wired in `module.go`.
