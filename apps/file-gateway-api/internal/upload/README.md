# internal/upload

Upload bounded context. Accepts multipart file uploads, validates them against a
pre-defined `UploadRule`, and persists them through the storage provider.

## Scope

- **Rule registry** — every upload endpoint is gated by an `UploadRule` keyed by a
  UUID `uploadId`. A rule defines the max size, allowed MIME types, target folder,
  public/private bucket, and whether the filename is overridden. Rules live in
  [`domain/upload/rules.go`](domain/upload/rules.go) (marketplace covers/gallery/
  documents, AI node assets, web3 token icons, organization icons/covers, support
  files, etc.).
- **Validation** — enforces the rule's size and MIME constraints before any write.
- **Filename resolution** — when `IsOverrideFilename` is set, generates a UUIDv7
  name with an extension derived from the content type; otherwise keeps the client
  filename.
- **Persistence** — uploads to the public or private bucket via the
  `storage.Provider`, returning the resulting URL/key.

## Endpoints (auth-protected)

| Method | Path | Description |
| --- | --- | --- |
| `GET`  | `/upload/:uploadId` | Returns the rule definition for `uploadId` |
| `POST` | `/upload/:uploadId` | Uploads a multipart file validated against the rule |

## Layout

```
domain/upload/       # UploadRule, File, rules registry, validators
domain/errors.go     # Upload-specific domain errors
application/upload/   # Service: GetRule + Upload orchestration
presentation/        # Routes + multipart form parser
module.go            # Wires the service with a storage.Provider; exposes Register()
```

> To add a new upload type, register a rule in `domain/upload/rules.go` — see the
> `upload-create-rule` skill.
