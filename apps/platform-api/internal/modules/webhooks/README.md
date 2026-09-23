# Webhooks

> Inbound HTTP edge that receives external webhook calls and routes them to the owning bounded context's processor.

## Responsibility
This context is a thin HTTP routing adapter for inbound webhooks. It exposes the public webhook endpoints, extracts the raw request (body, headers, query, path params), and delegates to the processor interface owned by taskrunner/triggers. It holds no domain model and no persistence of its own; it only marshals HTTP in and out.

## Domain
N/A — this module has no `domain/` or `application/` layer.

## HTTP API
- `GET /webhooks/triggers/:source/:token` — receive a workflow-trigger webhook (also handles provider verification handshakes; if the processor returns a `Verification`, its status code and body are sent verbatim). No auth middleware (authenticated by `:token`).
- `POST /webhooks/triggers/:source/:token` — same as above for POST deliveries; body is parsed as JSON into a `map[string]any` payload (empty map on parse failure), and a `VerificationRequest` (method, headers, body, query) is passed for processor-side signature verification.

## Dependencies
- **Bounded contexts (via service interfaces):**
  - `taskrunner/application/webhooks.WebhookProcessor` — processes trigger webhooks (request type from `triggers/application/webhooks.Request`, verification from `nodeengine/domain/adapters.VerificationRequest`).
- **Infrastructure:** none — no DB, cache, or outbound clients. The processor is injected via `Dependencies`.

## Layout
Only `http/` is present (`routes.go` + `triggers/` handler), wired in `module.go`: `Register` passes the injected processor interface straight to `http.RegisterRoutes`. There is no domain, application or persistence layer in this module.
