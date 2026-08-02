# Triggers

> Owns persisted workflow triggers (schedule + webhook) and resolves inbound webhooks into runnable context.

## Responsibility
This context owns the `triggers.triggers` table: schedule (cron) and webhook triggers attached to a workflow. It manages their lifecycle (update, delete, regenerate webhook token) and resolves an inbound webhook (token + provider adapter) into the execution context plus a `TriggerContext`. It does not execute tasks — it exposes services that the `taskrunner` context consumes.

## Domain
- **Aggregates / key types:**
  - `triggers.Trigger` — organization, `ExecutionContext`, `ContextItemID`, `Type`, optional `CronPattern`/`Timezone`, hashed webhook token + encrypted secret, `RuntimeConfig`, `IsActive`.
  - `triggers.TriggerType` — `schedule`, `webhook`.
  - `triggers.RuntimeConfig` — per-node runtime overrides applied at execution time.
- **Key rules / invariants:**
  - `UpdateSchedule` requires `Type == schedule` (`ErrNotScheduleTrigger`); `UpdateWebhook` requires `Type == webhook` (`ErrNotWebhookTrigger`).
  - Webhook token stored only as a hash; unique per active trigger (`ErrWebhookTokenTaken`).
  - Sentinels: `ErrTriggerNotFound`, `ErrUnsupportedType`, plus webhook-domain `ErrTriggerInactive` / `ErrFailedToDecryptSecret`.

## Trigger types & the adapter pattern

**Four ways a run starts — only two are persisted.** Taskrunner's `TaskTriggerType` covers `manual`, `schedule`, `webhook`, and `api` (plus reruns), but this context stores only the standing ones: `triggers.TriggerType` is exactly `schedule` and `webhook`. Manual and API runs are fire-and-forget; schedule/webhook runs are replayed from a persisted `Trigger` row.

| Type | Persisted? | Fired by | `TriggerContext` source |
| --- | --- | --- | --- |
| `manual` | no | User in the editor (`POST .../task-runner/trigger`, JWT + org permission) | Built from the request's `runtimePrompt`, if any |
| `schedule` | yes (`cron_pattern`, `timezone`) | Cron on the elected **taskrunner leader**; jobs reconciled periodically from `TriggerService.GetAllActive` | `RuntimeConfig.RuntimePrompt` (source `schedule`) |
| `webhook` | yes (hashed token, encrypted secret) | `webhook-api` on `GET/POST /triggers/:source/:token` | Provider **adapter** normalizes the inbound payload |
| `api` | no | API key with `ScopeWorkflowsTrigger` (`POST /task-runner/trigger/:workflowId`) | Built from the request's `runtimePrompt`, if any |

**Trigger registration rides the trigger endpoint.** Taskrunner's `ExecuteTask` doubles as the trigger factory: a request with type `schedule` + a cron pattern, or type `webhook` without an inbound payload, does not execute anything — it calls this context's `TriggerService.Create` to persist the trigger. For webhooks that generates a 32-byte random token whose **plaintext is returned exactly once**; only its SHA-256 hash is stored (`application/triggers/token.go`), unique per active trigger.

**Webhook URL & resolution.** Public webhook URLs come from the `WEBHOOK_TRIGGER_URL_TEMPLATE` env (default shape `${WEBHOOK_API_BASE_URL}/triggers/{source}/{token}`); nodeengine renders the `{source}` list for the UI, and the standalone `webhook-api` binary serves the route. On an inbound call, this context's `WebhookResolver` (consumed by taskrunner's `WebhookProcessor`) hashes the presented token, looks the trigger up by hash (`GetByWebhookTokenHash`), and rejects inactive triggers — an invalid token reveals nothing.

**The adapter pattern.** Each webhook-capable provider registers a `TriggerAdapter` (Discord, Slack, Telegram, WhatsApp — `internal/nodeengine/nodes/<provider>/webhook/adapter.go`) in a registry keyed by the URL's `:source` segment. `Adapt(raw)` normalizes the provider-specific payload into the uniform `TriggerContext{source, sender, prompt, payload}` — e.g. Telegram maps `message.text` → `prompt` and the username → `sender` — which is what nodes consume via `$trigger.*`. Adapters that also implement `WebhookVerifier` get two extra hooks, both using the trigger's secret (stored encrypted, decrypted via `secretmanager` at resolve time): `Verify` answers provider handshakes (Slack `url_verification` challenge, WhatsApp `hub.challenge`) and `ValidateSignature` checks the request's HMAC signature before any task is launched (Slack and WhatsApp implement it).

**RuntimeConfig: per-trigger execution overlay.** `RuntimeConfig` (`domain/triggers/runtime_config.go`) is `{runtimePrompt string, nodes []dag.Node}`, captured at trigger creation and stored in the `runtime_config` JSONB column. At fire time taskrunner's `ApplyRuntimeConfig` merges it over the freshly resolved workflow graph: `runtimePrompt` is stamped onto every node (and becomes `TriggerContext.Prompt` for schedule triggers, feeding `$trigger.prompt`), while the `nodes` overlay overrides `runtimeInstruction`, `credentials`, and `settings` per node ID — so one workflow can run with different prompts, credentials, and settings per trigger, without ever mutating the saved workflow document.

## Use cases (application)
- **UpdateTrigger** (command) — update active state, schedule/webhook fields, and runtime config.
- **DeleteTrigger** (command) — soft-delete a trigger.
- **RegenerateWebhookToken** (command) — issue a new plaintext token and store its hash.
- **GetAllTriggers** (query) — list triggers, enriched via the workflow service.
- **TriggerService** — `GetAllActive` and `Create` (with webhook token generation); consumed by other contexts.
- **WebhookResolver** — resolves token hash → trigger, verifies/validates provider signature, adapts payload to `TriggerContext`.

## HTTP API
- `GET /organizations/:organizationId/triggers/` — list (auth; `ReadTriggersPermission`).
- `PATCH /organizations/:organizationId/triggers/:triggerId` — update (auth; `UpdateTriggersPermission`).
- `POST /organizations/:organizationId/triggers/:triggerId/regenerate-webhook-token` — regenerate token (auth; `UpdateTriggersPermission`).
- `DELETE /organizations/:organizationId/triggers/:triggerId` — delete (auth; `DeleteTriggersPermission`).

## Events
- **Published:** none — `taskrunner` picks up schedule-trigger changes by periodically reconciling cron jobs against `TriggerService.GetAllActive`.

## Dependencies
- **Bounded contexts:** `workflows` (read enrichment), `nodeengine` adapters (webhook verify/adapt).
- **Infrastructure:** Postgres `triggers.triggers` table (schema `triggers`, cron/webhook columns, `runtime_config` JSONB, unique webhook-token-hash index), `secretmanager` (encrypt/decrypt webhook secret), transaction manager.

## Layout
Standard DDD layers present, wired in `module.go`.
