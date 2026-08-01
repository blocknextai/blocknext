# Credentials

> Owns user- and organization-scoped third-party integration secrets (API keys, OAuth data) used by workflow nodes.

## Responsibility
This bounded context stores, encrypts, and retrieves the credentials that workflow nodes need to call external services. It owns the persistence and at-rest encryption (via `SecretManager`) of credential data, the owner scoping (user vs. organization), and the distinction between owner-provided and host-provided platform credentials. It does not define which credential schemas exist — that vocabulary lives in `nodeengine` — nor does it execute nodes.

## Domain
- **Aggregates / key types:**
  - `Credential` — the aggregate: `OwnerType` + `OwnerID` scope, `SourceType`, `Key` (credential type), `Title`, encrypted `Data`; supports `Update` and soft `Delete`.
  - `CredentialInfo` — read projection returned to consumers with decrypted `Data` (`map[string]any`), `Key`, and `SourceType`.
  - `SourceType` — enum `owner` | `platform` (platform = host-provided managed OAuth apps).
- **Key rules / invariants:** owner type must be valid; owner ID, title, key, and data are all required (non-blank). Platform-source credentials must reference a key supported by the platform credential registry. Credential `Data` is always encrypted at rest and decrypted only via the service.

## Use cases (application)
- **CreateCredential** — encrypt and persist a new credential; validates platform-key support for `platform` source.
- **UpdateCredential** — re-encrypt and update an existing credential.
- **DeleteCredential** — soft-delete a credential.
- **GetCredentialByID** — fetch a single credential (scoped to owner).
- **GetAllCredentials** — list credentials for an owner.
- **GetCredentialsForNodes** — list credentials usable by a given set of node IDs.

The exported `CredentialService` adds `GetByIDForOwner` (returns decrypted `CredentialInfo`) and `SaveCredentialForOwner` for cross-context consumers (e.g. OAuth token refresh).

## HTTP API
User-scoped (`/users/me/credentials`) and organization-scoped (`/organizations/:organizationId/credentials`) variants, all behind `Authenticate()` plus the matching user/organization RBAC permission:

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
