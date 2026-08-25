# Credentials

> Owns organization-scoped third-party integration secrets (API keys, OAuth data) used by workflow nodes.

## Responsibility
This bounded context stores, encrypts, and retrieves the credentials that workflow nodes need to call external services. It owns the persistence and at-rest encryption (via `SecretManager`) of credential data, the organization ownership of every credential, and the distinction between owner-provided and host-provided platform credentials. It does not define which credential schemas exist — that vocabulary lives in `nodeengine` — nor does it execute nodes.

## Domain
- **Aggregates / key types:**
  - `Credential` — the aggregate: `OwnerType` + `OwnerID` scope, `SourceType`, `Key` (credential type), `Title`, encrypted `Data`; supports `Update` and soft `Delete`.
  - `CredentialInfo` — read projection returned to consumers with decrypted `Data` (`map[string]any`), `Key`, and `SourceType`.
  - `SourceType` — enum `owner` | `platform` (platform = host-provided managed OAuth apps).
- **Key rules / invariants:** owner type must be valid; owner ID, title, key, and data are all required (non-blank). Platform-source credentials must reference a key supported by the platform credential registry. Credential `Data` is always encrypted at rest and decrypted only via the service.

## Credential lifecycle & security

- **Encryption at rest (AES-256-GCM)** — credential `Data` is never stored as plaintext. `go-packages/secretmanager` JSON-marshals the value, encrypts it with AES-256-GCM under a hex-encoded 32-byte key (`SECRET_MANAGER_SECRET_KEY`, `notEmpty`-validated at boot), prepends a fresh random nonce per encryption, and stores the base64 ciphertext in `credentials.credentials.data`. Decryption is authenticated — a wrong key or tampered ciphertext fails GCM verification instead of yielding garbage.
- **Two credential sources** — `SourceType` distinguishes who supplies the secret material:

  | Source | Secret material | Validation |
  | --- | --- | --- |
  | `owner` | User/organization pastes their own API key or OAuth client into the schema form | Schema-driven; `Data` required non-blank |
  | `platform` | Host operator ships a managed OAuth app / API key via `PLATFORM_CREDENTIALS_*` env vars (base64 JSON, loaded once at boot by the platform config loader) | `CreateCredential` rejects any `platform`-source key not present in the platform registry (`ErrPlatformCredentialNotSupported`); `Data` may start as `{}` |

  For `platform` credentials the host's client ID/secret never enters the user's row: at use time the loader's data is merged **under** the credential's own data, and only the acquired OAuth token is ever written back. A credential type advertises platform support (`IsSupportPlatform`) automatically when its `PlatformID` has host config — it is derived, never hardcoded.
- **Owner scoping on every query** — credentials belong to `OwnerType` `user` or `organization` (+ `OwnerID`); the repository reads by `(id, ownerType, ownerID)` (`GetByIDAndOwner`), so a leaked credential ID from another tenant resolves to not-found. HTTP routes mirror the scope (`/users/me/…`, `/organizations/:organizationId/…`) behind `Authenticate()` + the matching RBAC permission.
- **Read APIs never return secrets** — list and by-nodes responses (`CredentialResponse`) carry metadata only (`id`, `title`, `key`, `uiKey`, `sourceType`, timestamps) — no `data` field exists in the struct. `GET /…/:credentialId` decrypts but then runs the data through the nodeengine `CredentialProcessor`: every schema property marked `WriteOnly` (e.g. `clientSecret`) is replaced with the `CredentialMaskValue` sentinel, and for OAuth credentials the whole `oauthToken` object is collapsed to a boolean "connected" flag.
- **Mask-aware updates** — `UpdateCredential` decrypts the existing data and merges the request over it, skipping any field whose value still equals the mask sentinel (`ShouldSkipField`). The UI can therefore echo a masked form straight back without wiping stored secrets; only genuinely changed fields overwrite.
- **Plaintext only at execution time** — decrypted values leave this context solely through the exported `CredentialService.GetByIDForOwner` (`CredentialInfo` with decrypted `map[string]any`), consumed by `credentialoauth` (token refresh) and the taskrunner when a node actually executes. No HTTP endpoint returns a decrypted secret.

## Use cases (application)
- **CreateCredential** — encrypt and persist a new credential; validates platform-key support for `platform` source.
- **UpdateCredential** — re-encrypt and update an existing credential.
- **DeleteCredential** — soft-delete a credential.
- **GetCredentialByID** — fetch a single credential (scoped to owner).
- **GetAllCredentials** — list credentials for an owner.
- **GetCredentialsForNodes** — list credentials usable by a given set of node IDs.

The exported `CredentialService` adds `GetByIDForOwner` (returns decrypted `CredentialInfo`) and `SaveCredentialForOwner` for cross-context consumers (e.g. OAuth token refresh).

## HTTP API
All under `/organizations/:organizationId/credentials`, behind `Authenticate()` plus the matching organization RBAC permission:

- `GET    /…/credentials` — list credentials (Read permission).
- `GET    /…/credentials/by-nodes` — credentials filtered by node IDs (Read permission).
- `GET    /…/credentials/:credentialId` — get one credential (Read permission).
- `POST   /…/credentials` — create credential (Create permission).
- `PUT    /…/credentials/:credentialId` — update credential (Update permission).
- `DELETE /…/credentials/:credentialId` — delete credential (Delete permission).

## Dependencies
- **Bounded contexts:** `nodeengine` (`CredentialProcessor`, `CredentialService` for node/credential metadata), `platform` (`PlatformCredentialService` for host-provided platform credential support).
- **Infrastructure:** Postgres table `credentials.credentials` (with `credential_source_type` and `owner_type` enums); `SecretManager` for at-rest encryption/decryption; `TransactionManager`.

## Layout
Standard DDD layers present, wired in `module.go`.
