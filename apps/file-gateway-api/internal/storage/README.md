# internal/storage

Storage bounded context. Defines the storage abstraction used by the upload
context and selects a concrete driver at startup.

## Scope

- **Provider interface** — a single `Provider` contract: `Upload` (to the public or
  private bucket) plus `HealthCheckPublic` / `HealthCheckPrivate` probes used by the
  readiness endpoint.
- **Driver selection** — `NewModule` builds the provider from `STORAGE_*`
  configuration, choosing one of the drivers below; an unknown driver returns
  `ErrUnsupportedStorageDriver`.
- **Public/private separation** — every driver is configured with independent
  public and private buckets.

## Drivers

| Driver | Package | Notes |
| --- | --- | --- |
| `s3`    | `infrastructure/s3`    | S3-compatible object storage |
| `bunny` | `infrastructure/bunny` | Bunny.net storage + pull zone |
| `local` | `infrastructure/local` | Local filesystem (public served at `/uploads/*`) |

## Layout

```
domain/storage/      # Provider interface + UploadResult
infrastructure/s3/    # S3-compatible implementation
infrastructure/bunny/ # Bunny.net implementation
infrastructure/local/ # Local filesystem implementation
module.go            # buildProvider() driver factory; exposes Module.Provider
```
