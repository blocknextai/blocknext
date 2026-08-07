---
name: credential-create
description: Use this skill when the user wants to add a new credential (third-party integration auth) to platform-api. Triggers on phrases like "add a credential for X", "create a new credential", "register a new API key type", "add OAuth2 for X". Generates the credential file under internal/nodeengine/credentials/ and wires it into the registry.
---

> **Working directory:** all relative paths in this skill (`internal/...`, `cmd/...`, `go build ./...`) are relative to `apps/platform-api/` in the monorepo. Make targets run from the repo root.


# Create Credential

Scaffolds a new credential for a third-party provider in the platform-api `nodeengine` module, wires it into the registry, and ensures the supporting nodes list is consistent.

## When to use

- User says "add a credential for `<provider>`" / "register OAuth2 for X" / "create a new credential type"
- A new node needs an auth type that does not yet exist (e.g., adding a Mistral chat node requires `mistral_api` first)

## Required information

Before generating files, confirm with the user (ask only what's missing):

1. **Provider name** — e.g., `mistral`, `notion`, `linear`. Lowercase, single word.
2. **Auth type** — one of:
   - `api` — single API key (or token + secret)
   - `oauth2` — standard OAuth2 authorization code flow
   - `oauth1` — OAuth1 (rare; only if explicitly requested)
   - `wallet` — Web3 wallet (used by `ethereum.wallet`)
3. **Display name** — human-friendly label (e.g., "Mistral", "Notion")
4. **Description** — one-line description of what the credential is for
5. **Schema fields** — for `api`: which fields (`apiKey`, `apiSecret`, `accountId`, etc). For `oauth2`: authorization URL, access token URL, scope string.
6. **Supported nodes** (optional) — list of node IDs that will consume this credential (e.g., `["mistral.chat"]`). Can be left empty/added later.
7. **PlatformID** — defaults to the credential ID itself. Use a shared platform ID (e.g., `google.oauth2`) only when multiple credentials share the same OAuth client/cloud project. Ask the user if this should be shared.

## File layout

Single file: `internal/nodeengine/credentials/<provider>_<authtype>.go`

Examples:
- `anthropic_api.go` → `NewAnthropicApiCredential`
- `gmail_oauth2.go` → `NewGmailOAuth2Credential`
- `ethereum_wallet.go` → `NewEthereumWalletCredential`

Naming rules:
- File: `<provider>_<authtype>.go` (snake_case)
- Function: `New<Provider><AuthType>Credential` (PascalCase). For OAuth2/OAuth1 the suffix is `OAuth2` / `OAuth1`, not `Oauth2`.
- Credential `ID`: `<provider>_<authtype>` (e.g., `mistral_api`, `notion_oauth2`)

## Icon

A credential wears its provider's plain brand mark; the glyph badge belongs to
individual nodes, so `CredentialIcon` carries only `Brand`. The value names the
artwork under `apps/platform/public/assets/icons/brands/<brand>/{light,dark}.svg`
— reuse the provider's existing brand, or drop the two files in for a provider
that has none. Nothing registers a brand in code.

## Templates

### API key credential

```go
package credentials

import (
	domain "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

func New<Provider>ApiCredential() *domain.Credential {
	return &domain.Credential{
		ID:          "<provider>_api",
		PlatformID:  "<provider>_api",
		Name:        "<DisplayName>",
		Description: "<one-line description>",
		Icon: domain.CredentialIcon{
			Brand: "<brand>",
		},
		Schema: &gjs.Schema{
			Type: "object",
			Properties: map[string]*gjs.Schema{
				"apiKey": {
					Type:        "string",
					Title:       "API Key",
					Description: "<provider description, e.g. 'API key from your dashboard.'>",
					WriteOnly:   true,
				},
			},
			Required: []string{
				"apiKey",
			},
		},
		SupportedNodes: &[]string{
			// "<provider>_<action>", ...
		},
	}
}
```

Rules:
- Secret-bearing fields (`apiKey`, `apiSecret`, `clientSecret`, `password`, `privateKey`) MUST set `WriteOnly: true`.
- `Required` lists every property in `Properties`; do not omit non-secret fields.
- `SupportedNodes` is `*[]string` — always a pointer. Empty pointer-to-slice is fine if no nodes exist yet.

### OAuth2 credential

```go
package credentials

import (
	"github.com/blocknextai/go-packages/json"
	domain "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

const (
	<provider>AuthorizationURL = "<auth url>"
	<provider>AccessTokenURL   = "<token url>"
	<provider>Scope            = "<space-separated scopes>"
)

func New<Provider>OAuth2Credential(redirectURL string) *domain.Credential {
	return &domain.Credential{
		ID:          "<provider>_oauth2",
		PlatformID:  "<provider>_oauth2", // or a shared platform like "google_oauth2"
		Name:        "<DisplayName>",
		Description: "<one-line description>",
		Icon: domain.CredentialIcon{
			Brand: "<brand>",
		},
		Schema: &gjs.Schema{
			Type: "object",
			Properties: map[string]*gjs.Schema{
				"authUrl": {
					Type:    "string",
					Default: json.RawMessage(`"` + <provider>AuthorizationURL + `"`),
					Extra: map[string]any{
						"hidden": true,
					},
				},
				"tokenUrl": {
					Type:    "string",
					Default: json.RawMessage(`"` + <provider>AccessTokenURL + `"`),
					Extra: map[string]any{
						"hidden": true,
					},
				},
				"redirectUrl": {
					Type:        "string",
					Title:       "Redirect URL",
					Description: "OAuth2 redirect URL configured in your <provider> app.",
					Default:     json.RawMessage(`"` + redirectURL + `"`),
					ReadOnly:    true,
					Extra: map[string]any{
						"copyable": true,
					},
				},
				"scope": {
					Type:    "string",
					Default: json.RawMessage(`"` + <provider>Scope + `"`),
					Extra: map[string]any{
						"hidden": true,
					},
				},
				"authentication": {
					Type:    "string",
					Default: json.RawMessage(`"body"`), // or "header" if provider requires it
					Extra: map[string]any{
						"hidden": true,
					},
				},
				"clientId": {
					Type:        "string",
					Title:       "Client ID",
					Description: "<provider> OAuth2 client ID.",
				},
				"clientSecret": {
					Type:        "string",
					Title:       "Client Secret",
					Description: "<provider> OAuth2 client secret.",
					WriteOnly:   true,
				},
			},
			Required: []string{
				"authUrl",
				"tokenUrl",
				"redirectUrl",
				"scope",
				"authentication",
				"clientId",
				"clientSecret",
			},
		},
		SupportedNodes: &[]string{
			// "<provider>_<action>", ...
		},
	}
}
```

Rules specific to OAuth2:
- The constructor takes `redirectURL string` as a parameter — the registry passes `oauth2RedirectURL` (read from config) into it. Do NOT read `os.Getenv(...)` and do NOT `json.Marshal` the env var. The `redirectUrl` field's `Default` is `json.RawMessage(\`"\` + redirectURL + \`"\`)`.
- All hidden fields (`authUrl`, `tokenUrl`, `scope`, `authentication`) MUST appear in `Required` and have `Extra: {"hidden": true}`. There is NO `grantType` field — the OAuth2 flow does not use one. (Some Google-family credentials add an `authQueryParameters` hidden field for `access_type=offline&prompt=consent`; add it only when the provider needs offline refresh.)
- `redirectUrl` is `ReadOnly: true` with `Extra: {"copyable": true}` — user copies it into the provider's app config.
- `authentication` is typically `"body"`. Use `"header"` only if the provider explicitly requires Basic-auth-style token exchange.
- Scopes go in a `const` at the top of the file; multiple scopes are space-separated in a single string.

## Critical rules (project memory)

1. **Never set `IsSupportPlatform` manually.** The `PlatformConfigLoader` derives it from `PlatformID` + bundled platform JSON at startup. Do not add `IsSupportPlatform: true` to the literal, do not call `SetIsSupportPlatform`, do not import a `ApplyPlatformSupport` helper (it was removed).
2. **`IsOAuth1` / `IsOAuth2` are auto-set** from substring matching on the ID (`oauth2` or `oauth1` in lowercased ID) inside `RegisterCredential`. Do NOT set them manually.
3. **Use `github.com/blocknextai/go-packages/json` for `json.RawMessage`** — `gjs.Schema.Default` needs a `json.RawMessage`, and `go-packages/json` exposes it as a type alias of the stdlib type, so it satisfies the jsonschema library while keeping the project-standard import. Do NOT import stdlib `encoding/json` here (the real credentials all import `go-packages/json`).
4. **Use `WriteOnly: true` on every secret** (apiKey, apiSecret, clientSecret, password, privateKey). Otherwise the value leaks back to the UI.

## Wiring

After creating the file, register the credential in `internal/nodeengine/credentials/credentials.go`'s `Register(oauth2RedirectURL string)` function:

1. Open `credentials.go`.
2. Add the registration line in the `Register(...)` body:
   - **OAuth2** credentials take the redirect URL: `domain.RegisterCredential(New<Provider>OAuth2Credential(oauth2RedirectURL))`
   - **API key / wallet / OAuth1** credentials take no args: `domain.RegisterCredential(New<Provider>ApiCredential())`
3. Keep the list **alphabetically sorted by function name** (the existing list is sorted; preserve it).

## Final steps

After writing the file and updating `credentials.go`:

1. Run `gofmt -w internal/nodeengine/credentials/<provider>_<authtype>.go internal/nodeengine/credentials/credentials.go`
2. Run `go build ./...` to verify the package compiles.
3. Run `golangci-lint run ./internal/nodeengine/credentials/...` to catch lint issues.
4. If you also added `SupportedNodes: &[]string{"<id>"}`, double-check those node IDs exist or will be added in the same change. A non-existent node ID here is harmless at runtime but is a documentation lie.

## What NOT to do

- Do NOT create a separate `register.go` per credential — credentials are registered centrally in `credentials.go`.
- Do NOT add `// TODO` placeholders for fields the user did not specify; ask instead.
- Do NOT add HTTP client helpers in this file — those live in `internal/nodeengine/nodes/<provider>/helpers/client.go` and are concerns of the node, not the credential.
- Do NOT use `fmt.Sprintf` for the `Default` JSON values; the existing pattern is string concatenation with `+` which is fine for cold init code.
