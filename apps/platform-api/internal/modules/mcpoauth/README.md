# MCP OAuth

> The OAuth 2.1 authorization server that issues the access tokens MCP clients use against the MCP resource servers (`mcp`, `mcpplatform`).

## Responsibility
Implements the authorization side of the MCP authorization spec: discovery (RFC 8414), client identification through Client ID Metadata Documents and dynamic client registration (RFC 7591), the authorization code flow with mandatory PKCE S256 and resource indicators (RFC 8707), `iss` in authorization responses (RFC 9207), refresh token rotation and revocation (RFC 7009). It owns consent: an authorization request waits until the signed-in user approves it for one of their organizations, and the resulting grant is what codes and refresh tokens belong to. It runs inside `mcp-api`; it does not serve MCP itself.

## Domain
- **Aggregates:**
  - `clients.Client` — a registered client, or one built from a Client ID Metadata Document URL. Secrets are stored as SHA-256 hashes, redirect URIs match exactly.
  - `authorizationrequests.AuthorizationRequest` — a pending consent with the validated redirect URI, scopes, PKCE challenge, resource and state.
  - `grants.Grant` — the durable consent (client + user + organization + scopes + resource). Revoking it revokes every refresh token of the grant.
  - `authorizationcodes.AuthorizationCode`, `refreshtokens.RefreshToken` — single-use credentials, persisted as hashes.
- **Protocol value types** (`domain/oauth2`): scopes, grant types, response types, token endpoint auth methods, the S256 code challenge method, RFC error codes and error response, resource canonicalization, redirect URI rules, opaque secret generation, the access token claims and the endpoint paths.
- **Key rules:**
  - PKCE is required and only `S256` is accepted; only the `code` response type and the `authorization_code` / `refresh_token` grants exist.
  - Replaying an authorization code or a rotated refresh token revokes the whole grant.
  - Approval requires the approving user to be a member of the chosen organization.
  - Errors about the client identity or its redirect URI are never redirected back to the client.

## Use cases
**Commands**
- `clients/registerclient` — dynamic client registration.
- `authorizationrequests/createauthorizationrequest` — validates the authorization request and returns the consent URL, or the client redirect carrying the RFC error.
- `authorizationrequests/approveauthorizationrequest`, `DenyAuthorizationRequest` — resolve consent and return the client redirect (`code` or `access_denied`, always with `state` and `iss`).
- `tokens/exchangetoken` — authorization code and refresh token grants with rotation and replay detection.
- `tokens/revoketoken` — revokes the grant behind a refresh or access token.
- `grants/revokegrant` — the user disconnects an application.

**Queries**
- `metadata/getauthorizationservermetadata`, `authorizationrequests/getauthorizationrequest`, `grants/getallgrants`.

**Shared**
- `clients.ClientResolver` — resolves and authenticates clients.
- `grants.GrantRevoker` — revokes a grant together with its refresh tokens.
- `tokens.AccessTokenSigner` / `AccessTokenValidator` — ports implemented by `tokensigner/` (RFC 9068 `at+jwt`, audience bound to the resource).

## HTTP API
| Method | Path | Purpose | Auth |
| --- | --- | --- | --- |
| GET | `/.well-known/oauth-authorization-server` | Authorization server metadata | Public (cached 1h) |
| POST | `/mcp-oauth/oauth2/register` | Dynamic client registration | Public |
| GET | `/mcp-oauth/oauth2/authorize` | Start the flow, redirect to consent | Public |
| POST | `/mcp-oauth/oauth2/token` | Token endpoint (form encoded) | Client credentials |
| POST | `/mcp-oauth/oauth2/revoke` | Token revocation | Client credentials |
| GET | `/mcp-oauth/authorization-requests/:authorizationRequestId` | Consent screen data | Authenticate |
| POST | `/mcp-oauth/authorization-requests/:authorizationRequestId/approve` | Approve for an organization | Authenticate |
| POST | `/mcp-oauth/authorization-requests/:authorizationRequestId/deny` | Deny | Authenticate |
| GET | `/users/me/mcp-oauth/grants` | List connected applications | Authenticate + `ReadUserPermission` |
| DELETE | `/users/me/mcp-oauth/grants/:grantId` | Disconnect an application | Authenticate + `UpdateUserPermission` |

The register, token and revoke endpoints answer protocol errors with the RFC `{"error", "error_description"}` body; their `http/` packages map domain sentinels to RFC error codes.

## Dependencies
- **Bounded contexts:** `organizations` (`OrganizationUserService` — membership check on approval).
- **Infrastructure:** Postgres schema `mcpoauth` (clients, grants, authorization_requests, authorization_codes, refresh_tokens); the cache service (Client ID Metadata Documents, 15 minutes); `go-packages/httpclient` with a public HTTPS destination guard; `golang-jwt`.

## Layout
Standard DDD layers, wired in `module.go`. `Module.AccessTokenValidator` (implements the kernel `common/auth.AccessTokenValidator`) and `Module.MetadataService` are the only OAuth surface outside this module: the common `AccessTokenMiddleware` authenticates requests through the former, and `mcp` / `mcpplatform` build their protected resource metadata from the latter, so issuer, resource URLs and token internals stay inside this module. Migrations are registered in `platform-api-migration`.
