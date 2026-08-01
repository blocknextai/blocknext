# EventBus

> Transactional-outbox event infrastructure: reliably persists, relays, and deduplicates domain events across the platform.

## Responsibility
This is shared infrastructure, not a business bounded context. It guarantees that domain events published inside a database transaction are durably stored (outbox) and later delivered exactly-once-effectively to in-process subscribers, with retries, dead-lettering, and consumer-side idempotency (inbox). Producers enqueue events; an in-process `Bus` fans out to typed subscribers; a background relay drains the outbox. It has no HTTP/presentation layer.

## Domain
- **Aggregates / key types:**
  - `OutboxMessage` — a persisted event with `EventName`, JSON `Payload`, `Status`, `Attempts`, `NextAttemptAt`, `LockedAt`; state machine via `MarkProcessed` / `Reschedule` / `MarkDead`.
  - `InboxEntry` — dedupe record keyed by (`HandlerKey`, `EventID`), recording that a specific consumer already processed a given event.
  - `Status` — enum `pending` | `processing` | `processed` | `dead`.
- **Key rules / invariants:** status transitions only valid from `processing`; event name and payload are required. Inbox primary key (`handler_key`, `event_id`) enforces per-consumer once-only processing.

## Capabilities (application)
- **publishing.PublisherService.Enqueue** — marshals a `DomainEvent` and inserts a `pending` outbox message (called within the producer's transaction).
- **relay.RelayService** — background loop: claims pending messages with SKIP-LOCKED batching, dispatches via the `Bus`, marks processed, reschedules with exponential backoff on failure, dead-letters after `MaxAttempts`, and a reclaim loop that frees messages stuck in `processing` past `StuckTimeout`.
- **idempotency.InboxService.EnsureOnce** — wraps a consumer's work so it runs at most once per (handler, event); event ID is propagated through context.
- **Bus** — in-process registry: `Subscribe[T]` registers typed handlers + a JSON decoder per event name; `Dispatch` decodes the payload and fans out, joining handler errors.

## Events
- **Published:** none of its own — it carries every other context's `commonDomain.DomainEvent` (event names follow `<context>.<entity>.<action>`).
- **Consumed:** none of its own — it dispatches to subscribers registered by other modules via `Subscribe`.

## Dependencies
- **Bounded contexts:** none directly; it is consumed by all event-producing/consuming modules.
- **Infrastructure:** Postgres schema `eventbus` with tables `outbox_messages` (enum `eventbus_outbox_message_status`, partial indexes on due/locked rows) and `inbox_entries`; relay tuning via `config.EventBusOptions` (poll interval, batch size, max attempts, backoff, stuck timeout, reclaim interval).

## Layout
Standard DDD layers present (no presentation), wired in `module.go`; relay started via `StartRelay`.
