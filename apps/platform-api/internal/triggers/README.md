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
