# Realtime

> Ephemeral server→client pub/sub broadcaster (Redis-backed). Infrastructure, not a business bounded context. Renamed from `messagebroker`.

## Responsibility
Fan-outs live task/node execution events to connected clients, scoped per organization. It is the ephemeral, server→client half of the project's two-layer messaging model; the durable server→server half is `internal/eventbus`. The two layers share the same Redis instance (the broker also wakes eventbus), but realtime carries best-effort, non-durable WebSocket traffic only.

## What it provides
- `Broadcaster` interface:
  - `PublishTaskEvent(ctx, *taskrunner/domain/task.TaskEvent)` / `PublishNodeEvent(ctx, *taskrunner/domain/node.NodeEvent)` — publish to the `organization:{id}` Redis channel.
  - `Subscribe(ctx, organizationID) (<-chan string, error)` — subscribe to one organization's channel; returns a buffered string channel of JSON payloads, closed on ctx cancel.
  - `Ping(ctx)` / `Close()`.
- `New(config.BrokerOptions) (Broadcaster, error)` — factory; builds the Redis broadcaster when `Type == BrokerTypeRedis`, else `ErrInvalidBrokerType`.

## Used by
Built in the assemblers that serve live updates (`PlatformAPI`, `WebhookAPI`, `TaskWorker`). Publishers: `taskrunner` (receives the `Broadcaster` to emit task/node events). Subscriber: `ws` module (consumes `Subscribe` to stream events to WebSocket clients). Health checks call `Ping`.

## Notes
- Two-layer model: `eventbus` = durable server→server (transactional outbox); `realtime` = ephemeral server→client. Do not route durable work through realtime.
- Config still uses `BROKER_*` env vars (`config.BrokerOptions`), reflecting the shared pub/sub Redis with eventbus wake; the package itself was renamed from `messagebroker`.
- Channel naming is `organization:{uuid}`; delivery is best-effort with a 64-slot buffered channel per subscriber.

## Layout
- `realtime.go` — `Broadcaster` interface.
- `factory.go` — `New` factory + `ErrInvalidBrokerType`.
- `redis/redis_broadcaster.go` — Redis pub/sub implementation (`PoolOptions`, publish/subscribe, marshal-error sentinels).
