# Executions

> Owns the record of workflow task runs, their per-node executions, and the worker claim/lease that drives distributed processing.

## Responsibility
This bounded context persists and reads the history of workflow executions: each task run, the individual node executions within it, and the claim record workers use to lease a task. It owns task/node execution lifecycle state (status, timing, inputs/outputs, errors) and the distributed claiming protocol (claim, heartbeat, release, stale reclaim, retry counting). It does not schedule or run nodes — execution engines/workers call its services to create and update records and to lease work.

## Domain
- **Aggregates / key types:**
  - `TaskExecution` — a single workflow run: organization scope, optional `TriggeredByUserID` / `FlowTriggerID`, `ExecutionContext`, `ContextItemID`, `Status`, `ExecutionType`, the `dag.Node`/`dag.Edge` snapshot, and start/complete timing.
  - `NodeExecution` — one node's run within a task: `NodeType`, `NodeID`, `Status`, `Inputs`/`Outputs`/`FunctionCallingOutputs`, error, timing.
  - `TaskClaim` — worker lease over a task: `ClaimedBy`, `ClaimedAt`, `RetryCount`; supports `Claim`/`Release`/`Heartbeat`/`ForceRelease`/retry counting with staleness checks.
  - `ExecutionType` — enum `manual` | `schedule` | `webhook` | `api`.
- **Key rules / invariants:** a claim can only be claimed by one worker; release/heartbeat require ownership; retry count is non-negative; node executions require task ID, node type, node ID, and status. Creating a `TaskExecution` also creates its `TaskClaim` in the same transaction.

## Use cases (application)
- **GetAllTaskExecutions** — paginated/searchable list of an organization's executions.
- **GetTaskExecutionByID** — fetch a single execution (organization-scoped).
- **DeleteTaskExecution** — soft-delete one execution.
- **BulkDeleteTaskExecutions** — soft-delete multiple executions.

Service interfaces (used by workers/engines, not HTTP): `TaskExecutionService` (create/update/query by id/status), `NodeExecutionService`, `TaskClaimService` (claim, heartbeat, release, release-stale, retry counting).

## HTTP API
All under `/organizations/:organizationId/executions`, behind `Authenticate()` plus organization RBAC:

- `GET    /executions` — list executions (ReadTaskExecution permission).
- `GET    /executions/:executionId` — get one execution (ReadTaskExecution permission).
- `DELETE /executions/bulk` — bulk delete (DeleteTaskExecution permission).
- `DELETE /executions/:executionId` — delete one (DeleteTaskExecution permission).

## Dependencies
- **Bounded contexts:** `organizations` (organization-user resolution), `workflows`, `account` (users, linked accounts) — consumed via service interfaces for enrichment/resolution.
- **Infrastructure:** Postgres schema `executions` with tables `task_executions`, `node_executions`, `task_claims` (unique per task execution); `TransactionManager`.

## Layout
Standard DDD layers present, wired in `module.go`.
