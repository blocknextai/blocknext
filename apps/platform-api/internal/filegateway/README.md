# FileGateway

> HTTP client abstraction over an external file-storage service. Infrastructure, not a business bounded context.

## Responsibility
Provides a small interface for uploading and downloading files through an external file-gateway service, authenticated with a service key. It wraps `go-packages/httpclient` (retrying, JSON-configured) and exposes streaming download and multipart upload; it does not implement storage itself.

## What it provides
- `FileGateway` interface:
  - `DownloadFile(ctx, url) (*DownloadResult, error)` — POSTs `/download`, streams the body back, and reads `X-Filename` / `X-Type` / `X-Ext` / `Content-Type` / `Content-Length` headers into the result.
  - `UploadFile(ctx, uploadID, fileName, reader) (*UploadResult, error)` — POSTs a multipart form to `/upload/{uploadID}`, returns the stored file URL.
- `NewFileGateway(baseURL, authServiceKey)` — builds the gateway with an `httpclient` configured for the base URL, `X-Service-Key` header, JSON content type, and retry (3 attempts).
- `DownloadResult` (filename, type, ext, content-type, size, `DataReader io.ReadCloser`) and `UploadResult` (`URL`).
- `ErrDownloadFailed`, `ErrReadFailed` (`apperror.Internal`).

## Used by
`bootstrap.NewCore` constructs it (`config.FileGatewayOptions`) as `Core.FileGateway`; passed into modules that handle files — `web3` and `nodeengine` (which forwards it to executors needing file IO).

## Notes
- Uses `go-packages/httpclient`, not `net/http`, per project convention.
- Download returns an open `io.ReadCloser`; the caller owns closing it.

## Layout
- `file_gateway.go` — `fileGateway` impl + `NewFileGateway`.
- `types.go` — `FileGateway` interface, `DownloadResult`, `UploadResult`.
- `errors.go` — sentinel errors.
