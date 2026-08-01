# Credential OAuth

> User-facing OAuth2 authorization-code (PKCE) flow and token-refresh service for third-party credentials.
>
> Named `credentialoauth` (not `oauth`) to distinguish it from `mcpoauth`: this context is an OAuth **client** that acquires third-party credentials, while `mcpoauth` is the OAuth **authorization server** that grants MCP clients access to this platform.

## Responsibility
This context owns the interactive OAuth2 connect flow for credentials: it builds a provider authorization URL, handles the redirect callback, exchanges the authorization code for tokens, and saves them onto the owning credential. It also provides a token-regeneration service that refreshes OAuth2 tokens on demand for other contexts. It has no persistence of its own — transient state and refresh locks live in cache; tokens are stored via the credentials context.

## Domain
- **Aggregates / key types:** (`domain/oauth2`, value types/DTOs)
  - `State` — cache-persisted callback payload (`CredentialID`, `OwnerType`, `OwnerID`, PKCE `CodeVerifier`).
  - `AuthenticationMethod` (`header`/`body`) and `GrantType` (`authorization_code`/`refresh_token`/`pkce`) enums.
  - `TokenResponse` / `TokenErrorResponse` — provider token-endpoint DTOs; `TokenResponse.ToToken()` maps to the shared `common/domain/oauth2.Token`.

## Use cases (application)
- **`authurl`** (command) — builds the provider authorization URL: loads credential schema, generates a PKCE challenge, issues a `State` into the cache, assembles the URL.
- **`exchangecode`** (query) — callback handler: consumes the state, exchanges the authorization code for tokens at the provider's `tokenUrl`, and saves them via the credentials context.
- **`StateStore`** — cache-backed `Issue`/`Consume` of opaque state IDs (`oauth:oauth2:state:`).
- **`regenerate.OAuthTokenRegenerateService`** — `RegenerateTokenIfNeeded` refreshes an OAuth2 token when needed, guarded by a Redis semaphore lock (`oauth:oauth2:refresh:`) to handle refresh-token rotation; consumed by the taskrunner context.

## HTTP API
- `POST /users/me/oauth/auth` — build auth URL (`Authenticate()` + `RequireUserPermission(UpdateUserCredentialsPermission)`).
- `POST /organizations/:organizationId/oauth/auth` — build auth URL (`Authenticate()` + `RequireOrganizationPermission(UpdateOrganizationCredentialsPermission)`).
- `GET /oauth2/callback` — exchange `code`+`state` for tokens (no auth; the provider redirects here).

## Dependencies
- **Bounded contexts:** credentials (`CredentialService`), nodeengine (`CredentialService` schema/hidden objects), platform (`PlatformCredentialService` for host-provided client credentials); taskrunner consumes `OAuthTokenRegenerateService`.
- **Infrastructure:** cache (Redis) for state + refresh locks; outbound HTTP (`httpclient`) form POSTs to provider token endpoints; PKCE helper. No DB tables, no events.

## Layout
Standard DDD layers present, wired in `module.go`.
