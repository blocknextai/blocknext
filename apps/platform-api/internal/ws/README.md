# WebSocket (ws)

> Per-organization WebSocket fan-out gateway that pushes realtime messages from the pub/sub layer to connected browser clients.

## Responsibility
This context owns the server-side WebSocket endpoint and the in-memory hub that groups live connections into per-organization rooms. It subscribes each room to the `realtime.Broadcaster` (Redis pub/sub) channel for that organization and fans incoming messages out to every connection in the room. It is push-only: client inbound frames are read solely to drive ping/pong keepalive and connection lifecycle, not interpreted as commands. It holds no persistence.

## Domain
- **Aggregates / key types:**
  - `connections.Connection` — a single WebSocket connection (ID, `UserID`, `OrganizationID`, buffered `Send` channel of 256, `Close()` via `sync.Once`).
  - `rooms.Room` — connections grouped by `OrganizationID`; `AddConnection`/`RemoveConnection`, `Broadcast`, `Stop` (cancels subscription), `CloseAllConnections`, with an optional `maxConnections` cap.
- **Key rules / invariants:**
  - One room per organization; a room with a `maxConnections` cap rejects (and closes) connections beyond it.
  - Broadcast is non-blocking: if a connection's send buffer is full the message is dropped (logged), never blocking the room.
  - When a room's last connection is removed it is stopped (subscription cancelled) and deleted from the hub.

## Use cases (application)
- `hub.HubService` — in-memory hub (`Register`, `Unregister`, `Shutdown`). On first connection per org it creates a room and starts a goroutine subscribed to the broadcaster channel; on shutdown it cancels all subscriptions and closes all connections.

## HTTP API
- `GET /organizations/:organizationId/ws` — upgrade to a WebSocket. Auth: `AuthenticateWebSocket()` + `RequireOrganizationPermission(ReadOrganizationPermission)`; rejects non-upgrade requests with `426 Upgrade Required`. The handler maintains ping/pong keepalive (54s ping period, 60s pong wait, 10s write deadline) and registers/unregisters the connection with the hub.

## Events
N/A — does not publish or consume durable domain events. It consumes the ephemeral `realtime` pub/sub stream per organization.

## Dependencies
- **Bounded contexts:** none (no cross-context service interfaces).
- **Infrastructure:** `realtime.Broadcaster` (Redis pub/sub, subscribed per organization channel); `common/presentation/auth.AuthMiddleware` for the upgrade handshake. No database.

## Layout
`application/` (hub), `domain/` (connections, rooms), and `presentation/` (websockets) are present and wired in `module.go`; there is no `infrastructure/` layer. `module.go` also exposes `Shutdown()` for graceful teardown of the hub.
