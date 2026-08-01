# file-gateway-api

A Go-based file gateway service built with [Fiber v3](https://gofiber.io/) that provides rule-based file uploads to S3-compatible storage and proxied file downloads from remote URLs.

## Features

- **Rule-based uploads**: Each upload endpoint is gated by a pre-defined `UploadRule` (max size, allowed MIME types, target folder, public/private bucket, filename override).
- **S3 storage**: Pluggable storage layer with separate public and private bucket configuration.
- **URL-based downloads**: Streams remote files through the API with size, timeout, and redirect limits.
- **JWT or service-key auth**: Upload and download routes accept either a service-key header (server-to-server) or a JWT minted via `/auth/token`. The mint endpoint itself is unauthenticated so trusted UI clients can obtain short-lived tokens; security relies on per-rule upload constraints + IP rate limiting.
- **Production hardening**: Helmet, CORS, rate limiting, request IDs, panic recovery, and graceful shutdown.

## Requirements

- Docker and Docker Compose
- Go 1.26+ (only for local builds outside Docker)
- S3-compatible storage credentials (public + private buckets)

## Setup

1. Copy the example environment file and fill in the values:
   ```bash
   cp .env.example .env
   ```
   Generate the secrets with:
   ```bash
   openssl rand -hex 32   # AUTH_SERVICE_KEY
   openssl rand -hex 32   # AUTH_JWT_SECRET
   ```

2. Start the stack from the monorepo root — see the root `README.md`.

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

### Upload (JWT required)

| Method | Path | Description |
| --- | --- | --- |
| `GET`  | `/upload/:uploadId` | Returns the rule definition for `uploadId` |
| `POST` | `/upload/:uploadId` | Uploads a multipart file validated against the rule |

`uploadId` values are defined in [`internal/upload/domain/upload/rules.go`](internal/upload/domain/upload/rules.go) (e.g. AI node generated assets).

### Download (JWT required)

| Method | Path | Description |
| --- | --- | --- |
| `POST` | `/download` | Streams a remote URL back to the client with file metadata headers |

Body: `{ "url": "https://..." }`

## Configuration

Configuration is loaded from environment variables via [`caarlos0/env`](https://github.com/caarlos0/env). See [`.env.example`](.env.example) for the full list. Key groups:

- `API_*` — Fiber server tuning (timeouts, body limit, CORS, rate limit)
- `AUTH_*` — Service key + JWT signing
- `STORAGE_S3_PUBLIC_*` / `STORAGE_S3_PRIVATE_*` — S3 bucket credentials
- `DOWNLOAD_*` — Limits applied to remote downloads

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
