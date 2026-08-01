# go-packages

Blocknext monorepo of reusable Go packages.

## Installation

```bash
go get github.com/blocknextai/go-packages
```

Requires Go `1.26.5`.

## Packages

| Package | Description |
| --- | --- |
| [`apperror`](./apperror) | Typed application errors (`Validation`, `NotFound`, `Conflict`, `Unauthorized`, etc.) with `Kind`-to-HTTP-status mapping. |
| [`auth/jwt`](./auth/jwt) | JWT access/refresh token generation and validation. |
| [`base64`](./base64) | Base64 encode/decode helpers: standard, URL, and raw-URL variants plus `io.Reader` streaming. |
| [`cache`](./cache) | `Service` abstraction (get/set, incr/decr, exists, keys, expire, atomic semaphores) with a Redis implementation under [`cache/redis`](./cache/redis). |
| [`cast`](./cast) | Loose type coercion helpers: `ToString`, `ToBool`, `ToFloat`, `ToInt64`, `ToSlice`. |
| [`dag`](./dag) | Directed Acyclic Graph with topological sort and priority queue. |
| [`database`](./database) | PostgreSQL helpers: connection pool, `Executor` interface, `TransactionManager` with context propagation, `BaseRepository`, generic `GetOne`/`GetMany`/`GetManyWithTotal` scanners, null scanners, and pq error helpers. |
| [`digest`](./digest) | Digest helpers (e.g. `SHA256Hex`). |
| [`email`](./email) | `EmailSender` abstraction with `Cc`/`Bcc`, per-recipient delivery (`SendSeparately`), and default-from support; SMTP ([`email/smtp`](./email/smtp)), Resend ([`email/resend`](./email/resend)), SendGrid ([`email/sendgrid`](./email/sendgrid)), and log ([`email/log`](./email/log)) implementations. |
| [`fiber/errorhandler`](./fiber/errorhandler) | Fiber `ErrorHandler` that maps `apperror`/Fiber errors to JSON responses, with production-safe messages. |
| [`fiber/middleware/metrics`](./fiber/middleware/metrics) | Prometheus observability middleware for Fiber: request count, latency, response size, and in-flight gauges, plus a `/metrics` endpoint handler. |
| [`fiber/storage/cache`](./fiber/storage/cache) | Fiber `Storage` adapter backed by a `cache.Service` (e.g. for session/limiter stores). |
| [`file`](./file) | File helpers: MIME-type detection, extension/category lookup, filename extraction, and reader-to-bytes. |
| [`hashing`](./hashing) | `Hasher` abstraction with a bcrypt implementation under [`hashing/bcrypt`](./hashing/bcrypt). |
| [`hex`](./hex) | Hex encode/decode helpers. |
| [`httpclient`](./httpclient) | Fluent HTTP client with retry, multipart, and OAuth1 support. |
| [`json`](./json) | `encoding/json` wrapper (`Marshal`/`Unmarshal`, `RawMessage`, `Number`) plus `ArgsToStruct`. |
| [`pkce`](./pkce) | PKCE (RFC 7636) code verifier/challenge helpers. |
| [`rbac`](./rbac) | Role-based access control: permissions, roles, and user/organization role bindings. |
| [`result`](./result) | Generic `Result[T]` wrapper with pagination, search, and metadata. |
| [`secretmanager`](./secretmanager) | `SecretManager` abstraction for encrypt/decrypt. |
| [`url/platform`](./url/platform) | Detect the platform/source of a URL. |
| [`uuid`](./uuid) | UUID generation helpers (`NewV4`, `NewV7`). |
