# internal/download

Download bounded context. Proxies a remote URL through the API, streaming the
response body back to the caller with file-metadata headers.

## Scope

- **Remote fetch** — given `{ "url": "https://..." }`, fetches the resource and
  streams it to the client without buffering the whole body in memory.
- **Safety limits** — enforces a max response size, request timeout, and redirect
  limits sourced from `DOWNLOAD_*` configuration. Oversized bodies are cut off via
  the shared `common.LimitedReader`.
- **Metadata headers** — sets `Content-Type`, `Content-Disposition`,
  `Content-Length`, and `X-Filename` / `X-Type` / `X-Ext`, with header values
  sanitized against injection.

## Endpoint (auth-protected)

| Method | Path | Description |
| --- | --- | --- |
| `POST` | `/download` | Streams a remote URL back to the client |

## Layout

```
domain/download/     # Service interface + FileInfo value type
domain/errors.go     # Download-specific domain errors
application/download/ # Service implementation (HTTP fetch, limits)
presentation/        # /download route + response headers
module.go            # Wires the service with DownloadOptions; exposes Register()
```
