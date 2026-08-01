# Platform

> Owns host-provided ("platform") OAuth/API credentials that the host operator supplies for end users, exposed read-only so the UI knows which integrations the platform pre-configures.

## Responsibility
This context loads operator-supplied credential configs (base64-encoded, keyed by node-engine platform ID), matches them against the node-engine credential catalog, and serves them read-only. These are the self-hosted equivalent of n8n cloud's managed OAuth apps ("use platform credentials" toggle). The module stores nothing in a database; it builds an in-memory map at startup from config and the node-engine credential registry, and flags those credentials as platform-supported.

## Domain
- **Aggregates / key types:**
  - `platformcredentials.PlatformCredential` — a credential identified by ID with an opaque `Data map[string]any`.
- **Key rules / invariants:**
  - Lookup of an unknown credential ID returns `ErrPlatformCredentialNotFound`.

## Use cases (application)
- **GetAllPlatformCredentials** (query) — list all configured platform credentials.
- **GetPlatformCredentialByID** (query) — fetch one by ID.
- **PlatformCredentialService** — read accessor over the loaded credentials, backed by `PlatformConfigLoader`.

## HTTP API
- `GET /platform/credentials` — list platform credentials (auth: `ReadPlatformCredentialPermission`, cached 5m).
- `GET /platform/credentials/:id` — get one by ID (auth: `ReadPlatformCredentialPermission`, cached 5m).

## Events
None published or consumed.

## Dependencies
- **Bounded contexts:** `nodeengine` (`CredentialService` — its credential catalog drives which platform IDs are loaded; the loader calls `SetIsSupportPlatform(true)` on matched credentials).
- **Infrastructure:** none — no DB tables; credentials come from operator config (`CredentialConfigs`, base64-decoded `PLATFORM_CREDENTIALS_*` host-provided OAuth apps), held in memory.

## Layout
Standard DDD layers present, wired in `module.go`.
