# internal/config

Centralized, environment-based configuration. Loaded once in `main` and passed
into each module's dependencies — modules never read the environment directly.

## Scope

- **`Load()`** — parses the environment into `ConfigurationOptions` via
  [`caarlos0/env`](https://github.com/caarlos0/env), returning an error on invalid
  input.
- **Option groups** — one file per concern, each bound to an env prefix:

| File | Prefix | Covers |
| --- | --- | --- |
| `api.go` | `API_` | Fiber server tuning, CORS, rate limit, trusted proxies |
| `auth.go` | `AUTH_` | Service key + JWT signing/expiry |
| `cache.go` | `CACHE_` | Cache type + Redis connection/pool |
| `storage.go` | `STORAGE_` | Driver selection + S3/Bunny/Local buckets |
| `download.go` | `DOWNLOAD_` | Remote download size/timeout limits |
| `metrics.go` | `API_METRICS_` | Prometheus metrics endpoint toggle/port |
| `env.go` | `APP_ENV` | App environment helper (`IsProduction`) |

See [`.env.example`](../../.env.example) for the full list of variables.
