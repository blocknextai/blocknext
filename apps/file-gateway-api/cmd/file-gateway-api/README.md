# cmd/file-gateway-api

Application entrypoint. This is the only `main` package and is responsible for
wiring the bounded contexts together — it contains no business logic of its own.

## Scope

- **Configuration** — loads `config.ConfigurationOptions` from the environment and
  fails fast on error.
- **Fiber app construction** — builds the HTTP server with the configured timeouts,
  body limits, buffers, trusted-proxy settings, JSON codec, and error handler.
- **Global middleware** — request IDs, CORS, Helmet, panic recovery, and IP-based
  rate limiting (bypassed for service-key callers).
- **Module composition** — instantiates the `cache`, `storage`, `auth`, `upload`,
  and `download` modules and registers their routes. Upload/download are mounted
  under an auth-protected router group; `/auth/token` stays public.
- **Health probes** — `/livez` (liveness) and `/readyz` (readiness, which pings
  storage public/private buckets and the cache).
- **Static serving** — exposes `/uploads/*` only when the local storage driver is
  active.
- **Metrics server** — Prometheus wiring is present in the code but commented out.
- **Lifecycle** — starts the listener(s) and performs graceful shutdown on
  `SIGINT`/`SIGTERM`, closing the cache connection last.

## Dependencies

Depends on every `internal/*` module's public `Module`/`NewModule` surface plus
`internal/config`. Nothing depends on this package.
