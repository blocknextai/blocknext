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

## OAuth flow & token refresh

**Connect flow (authorization code + PKCE).** The user connects a third-party account in three steps, all provider parameters resolved server-side:

| Step | What happens | Hardening |
| --- | --- | --- |
| 1. Auth URL | `POST …/credential-oauth/oauth2/auth` loads the credential (merging the host's client ID/secret for `platform`-source credentials), reads `authUrl`/`scope`/`authentication`/`authQueryParameters` from the credential schema's **hidden** objects in nodeengine, and assembles the provider URL | Provider endpoints come from server-side schema, never from the client; PKCE S256 challenge (`go-packages/pkce`) on every flow |
| 2. State | A 32-byte random base64url state ID is issued; the payload (`CredentialID`, `OwnerType`, `OwnerID`, `CodeVerifier`) lives only in Redis under `credentialoauth:oauth2:state:` with `CREDENTIAL_OAUTH_STATE_TTL` | State is opaque — the callback carries no identity or credential data; TTL bounds the window |
| 3. Callback | Provider redirects to `GET /credential-oauth/oauth2/callback` (`CREDENTIAL_OAUTH_OAUTH2_REDIRECT_URL`); the state is consumed via Redis `GetAndDelete` (single-use — replay/CSRF gets `ErrInvalidState`), the code is exchanged at the schema's `tokenUrl` with the stored `code_verifier`, and the token is saved through the credentials context | Client secret sent via Basic auth header or form body per the schema's `authentication` method; tokens are AES-256-GCM encrypted at rest by the credentials context; for `platform` credentials only the `oauthToken` is persisted — the host's client secret never enters the user's row |

**Configured OAuth2 providers** (credential definitions in `nodeengine/credentials/`): Facebook, Gmail, Google Docs, Google Drive, Google Keep, Google Sheets, Instagram, LinkedIn, Notion, Slack, SoundCloud, TikTok, X, YouTube — the Google family shares one `google_oauth2` platform app.

**Automatic refresh at execution time.** Before a node runs, the taskrunner's credential processor resolves every credential reference through `RegenerateTokenIfNeeded`:

- **When** — only if the stored token has a refresh token *and* expires within a 5-minute buffer (`Token.NeedsRefresh`); otherwise the stored data is returned untouched. Non-OAuth2 credentials skip refresh entirely (platform-source ones get the host data merged in).
- **Single-flight refresh** — a per-credential Redis semaphore (`credentialoauth:oauth2:refresh:<id>`, 30s TTL) serializes refreshes across workers; contenders wait up to 10s (200ms backoff) and then read the winner's freshly stored token instead of refreshing again. The winner re-reads and re-checks expiry after acquiring the lock, so providers that **rotate refresh tokens** are never burned twice.
- **Rotation-aware persistence** — if the provider returns a new refresh token it replaces the old one (rotation is logged); otherwise the existing refresh token is kept alongside the new access token. The updated token is written back encrypted via `CredentialService.SaveCredentialForOwner`.
- **On failure** — a provider `invalid_grant` maps to `ErrRefreshTokenInvalid`: the taskrunner logs "credential needs re-authentication" and the node execution errors out — the user must reconnect via the flow above. Any other provider error surfaces as `ErrTokenRefreshFailed` and also fails the execution; stale tokens are never silently used.
- **Extension point** — a credential type may implement `RefreshableCredential` to override refresh; all current providers use the generic `refresh_token` grant against the schema's `tokenUrl`.

## Use cases (application)
- **`authurl`** (command) — builds the provider authorization URL: loads credential schema, generates a PKCE challenge, issues a `State` into the cache, assembles the URL.
- **`exchangecode`** (query) — callback handler: consumes the state, exchanges the authorization code for tokens at the provider's `tokenUrl`, and saves them via the credentials context.
- **`StateStore`** — cache-backed `Issue`/`Consume` of opaque state IDs (`oauth:oauth2:state:`).
- **`regenerate.OAuthTokenRegenerateService`** — `RegenerateTokenIfNeeded` refreshes an OAuth2 token when needed, guarded by a Redis semaphore lock (`oauth:oauth2:refresh:`) to handle refresh-token rotation; consumed by the taskrunner context.

## HTTP API
- `POST /organizations/:organizationId/oauth/auth` — build auth URL (`Authenticate()` + `RequireOrganizationPermission(UpdateOrganizationCredentialsPermission)`).
- `GET /oauth2/callback` — exchange `code`+`state` for tokens (no auth; the provider redirects here).

## Dependencies
- **Bounded contexts:** credentials (`CredentialService`), nodeengine (`CredentialService` schema/hidden objects), platform (`PlatformCredentialService` for host-provided client credentials); taskrunner consumes `OAuthTokenRegenerateService`.
- **Infrastructure:** cache (Redis) for state + refresh locks; outbound HTTP (`httpclient`) form POSTs to provider token endpoints; PKCE helper. No DB tables, no events.

## Layout
Standard DDD layers present, wired in `module.go`.
