# file-gateway-api

A Go-based file gateway service built with [Fiber v3](https://gofiber.io/) that provides rule-based file uploads and proxied file downloads from remote URLs.

## Features

- **Rule-based uploads**: Each upload endpoint is gated by a pre-defined `UploadRule` (max size, allowed MIME types, target folder, public/private bucket, filename override).
- **Pluggable storage**: Local disk (default), S3-compatible, or Bunny.net drivers with separate public and private bucket configuration.
- **URL-based downloads**: Streams remote files through the API with size, timeout, and redirect limits.
- **JWT or service-key auth**: Upload and download routes accept either a service-key header (server-to-server) or a JWT minted via `/auth/token`. The mint endpoint itself is unauthenticated so trusted UI clients can obtain short-lived tokens; security relies on per-rule upload constraints + IP rate limiting.
- **Production hardening**: Helmet, CORS, rate limiting, request IDs, panic recovery, and graceful shutdown.

## Requirements

- Docker and Docker Compose
- Go 1.27.0 (only for local builds outside Docker)
- S3-compatible or Bunny.net storage credentials (optional — the default driver stores files on local disk)

## Setup

Configuration lives in the single `.env` file at the monorepo root — run `make setup` there to generate it (secrets included) from `.env.example`, then start the stack as described in the root `README.md`.

The API listens on `http://localhost:3300` by default.

## API

Health probes are exposed at `/livez` and `/readyz`.

### Auth

| Method | Path | Description |
| --- | --- | --- |
| `POST` | `/auth/token` | **Public** — issues a short-lived JWT bound to a random session ID. Intended for trusted UI clients that need a token to call upload/download. IP rate-limited; not gated by the service key. |

Upload and download routes accept either:
- `X-Service-Key: <key>` header (preferred for server-to-server calls)
- `Authorization: Bearer <jwt>` from `/auth/token`

### Upload (auth required)

| Method | Path | Description |
| --- | --- | --- |
| `GET`  | `/upload/:uploadId` | Returns the rule definition for `uploadId` |
| `POST` | `/upload/:uploadId` | Uploads a multipart file validated against the rule |

`uploadId` values are defined in [`internal/upload/domain/upload/rules.go`](internal/upload/domain/upload/rules.go) (e.g. AI node generated assets).

### Download (auth required)

| Method | Path | Description |
| --- | --- | --- |
| `POST` | `/download` | Streams a remote URL back to the client with file metadata headers |

Body: `{ "url": "https://..." }`

## Configuration

Configuration is loaded from environment variables via [`caarlos0/env`](https://github.com/caarlos0/env). See the root [`.env.example`](../../.env.example) for the full list. Key groups:

- `FILE_GATEWAY_API_*` — Fiber server tuning (timeouts, body limit, CORS, rate limit)
- `FILE_GATEWAY_AUTH_*` — Service key + JWT signing
- `FILE_GATEWAY_STORAGE_*` — Storage driver selection + per-driver (local/S3/Bunny) public and private bucket settings
- `FILE_GATEWAY_DOWNLOAD_*` — Limits applied to remote downloads
- `CACHE_*` — Shared Redis cache (same instance as platform-api)

## Project Layout

Each module follows the DDD pattern (`application/`, `domain/`, `infrastructure/` where needed, `presentation/`):

```
cmd/file-gateway-api/                          # Entrypoint, Fiber wiring, graceful shutdown
internal/auth/{domain,infrastructure,presentation}/   # JWT issuance + middleware
internal/upload/{application,domain,presentation}/    # Upload rules, validators, form parser, handlers
internal/download/{application,domain,presentation}/  # Remote URL download service
internal/storage/{domain,infrastructure}/             # Storage Provider interface + S3/Bunny/Local impls
internal/common/domain/                               # LimitedReader, shared errors
internal/cache/infrastructure/                        # Redis-backed cache wiring
internal/config/                                      # Env-based configuration
```
