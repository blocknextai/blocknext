# API Keys

> Owns programmatic-access API keys for organizations: their issuance, scoping, rotation, and validation.

## Responsibility
This context manages the lifecycle of API keys that owners (organizations) use for programmatic access to the platform. It generates keys (`bnx_`-prefixed, stored only as a SHA-256 hash), enforces per-key capability scopes, and exposes a validator that other contexts use to authenticate inbound API-key requests. It does NOT own RBAC user/role permissions (those gate the management endpoints via the common auth middleware) and does NOT decide what each scope ultimately authorizes downstream.

## Domain
- **Aggregates / key types:**
  - `APIKey` — an issued key bound to an owner (`OwnerType` + `OwnerID`), holding the key hash, a `Name`, its `Scopes`, and `LastUsedAt`.
  - `Scope` / `Scopes` — closed capability vocabulary: `workflows:trigger`, `mcp:invoke` (no wildcard).
  - `GeneratedKey` — transient pair of plaintext (returned once) and its hash (persisted).
- **Key rules / invariants:**
  - Name must be non-empty; `OwnerType` valid; `OwnerID` non-nil; key hash non-empty.
  - A key must carry at least one scope, and every scope must be valid (`Scope.IsValid()`).
  - The plaintext key is shown only at creation/regeneration; only the SHA-256 hash is stored.
  - Unique index on `key_hash` (where not soft-deleted).

## Use cases (application)
**Commands**
- `createapikey` — generate a new key for an owner with required scopes; returns the plaintext key once.
- `regenerateapikey` — rotate an existing key's hash, returning a new plaintext key.
- `deleteapikey` — soft-delete a key.

**Queries**
- `getallapikeys` — list an organization's API keys.
- `getscopes` — return the catalog of available scopes (no-auth, cached).

## HTTP API
| Method | Path | Purpose | Auth |
| --- | --- | --- | --- |
| GET | `/organizations/:organizationId/api-keys/` | List org API keys | Authenticate + `ReadAPIKeyPermission` |
| POST | `/organizations/:organizationId/api-keys/` | Create an API key | Authenticate + `CreateAPIKeyPermission` |
| POST | `/organizations/:organizationId/api-keys/:apiKeyId/regenerate` | Rotate a key | Authenticate + `UpdateAPIKeyPermission` |
| DELETE | `/organizations/:organizationId/api-keys/:apiKeyId` | Delete a key | Authenticate + `DeleteAPIKeyPermission` |
| GET | `/api-keys/scopes` | List available scopes | Public (cached 1h) |

## Dependencies
- **Bounded contexts:** None consumed. Exposes `APIKeyValidator` (returns `AuthenticatedAPIKey` with owner + scopes) for the common API-key auth middleware to authenticate inbound requests; `Validate` also updates `LastUsedAt`.
- **Infrastructure:** Postgres table `apikeys.api_keys` (hash + owner indexes); `go-packages/digest` (SHA-256) and `crypto/rand` for key generation.

## Layout
Standard DDD layers: `domain/` (aggregates + rules), `application/` (CQRS handlers), `infrastructure/` (persistence + adapters), `presentation/` (HTTP). Wired in `module.go`.
