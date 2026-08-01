---
name: credential-update
description: Use this skill when the user wants to modify an existing credential definition under `internal/nodeengine/credentials/`. Triggers on phrases like "update the X credential", "add a new field to Y credential", "extend OAuth2 scopes for Z", "add `<node>` to credential's supported nodes", "rename a credential field", "deprecate / disable a credential". Locates the credential file, applies the schema/field/scope/SupportedNodes change, and propagates cross-references (executors that read renamed fields, nodes' `SupportedCredentials` that need updating).
---

> **Working directory:** all relative paths in this skill (`internal/...`, `cmd/...`, `go build ./...`) are relative to `apps/platform-api/` in the monorepo. Make targets run from the repo root.


# Credential Update

Edits an existing credential definition in `internal/nodeengine/credentials/<provider>_<authtype>.go` and keeps cross-referenced files (executors, nodes, registry) consistent.

## When to use

- User says "update `<provider>` credential" / "add `apiSecret` to X credential" / "extend Linear OAuth2 scopes" / "add `linear.createissue` to credential's supported nodes" / "rename `apiKey` to `accessToken`" / "deprecate the `<provider>` credential"
- For brand-new credentials use `credential-create` instead.

## Required information

Confirm with the user (ask only what's missing):

1. **Credential ID** — the `<provider>_<authtype>` ID (e.g., `mistral_api`, `gmail_oauth2`). The file is `internal/nodeengine/credentials/<provider>_<authtype>.go`.
2. **Change type** — one or more of:
   - **Add field** to the JSON schema (with type, title, description, optional default, optional WriteOnly)
   - **Remove field** from the schema (and update every executor that reads it)
   - **Rename field** (schema + every executor's `credential.String("<old>")` → `"<new>"`)
   - **Change OAuth2 scope** (the `<provider>Scope` const at the top of the file)
   - **Update OAuth2 endpoints** (`<provider>AuthorizationURL` / `<provider>AccessTokenURL`)
   - **Add / remove a supported node** (`SupportedNodes` slice)
   - **Update display metadata** (Name, Description, Icon)
   - **Disable / enable** the credential (`Disabled: true/false`)
3. **Field details** (for add / change) — name, JSON schema type (`string`, `number`, `boolean`, `object`), Title, Description, Default value (if any), WriteOnly flag (mandatory for secrets).
4. **OAuth2 scope changes** — the new scope string (space-separated for OAuth2, comma-separated NOT used here). Confirm whether it's additive or a replacement; additive scopes still require user re-consent in some providers.

## Workflow

### Step 1: Locate and read the file

`internal/nodeengine/credentials/<provider>_<authtype>.go`. Read it end-to-end before editing. The file's structure:

```
package credentials

import (...)

const (
    <provider>AuthorizationURL = "..."   // OAuth2 only
    <provider>AccessTokenURL   = "..."   // OAuth2 only
    <provider>Scope            = "..."   // OAuth2 only
)

func New<Provider><AuthType>Credential() *domain.Credential {
    return &domain.Credential{
        ID, PlatformID, Name, Description, Icon,
        Schema: &gjs.Schema{ Properties: ..., Required: ... },
        SupportedNodes: &[]string{...},
    }
}
```

### Step 2: Apply the change

Use **`Edit`**, not `Write` — overwriting the whole file risks wiping unrelated lines. Apply the smallest possible diff per change.

#### Add a field

In `Schema.Properties`, add a new key with its `*gjs.Schema`. In `Required`, append the new field name unless it's truly optional (no default and safe to omit).

```go
Properties: map[string]*gjs.Schema{
    "apiKey": { ... },
    "apiSecret": {                                      // ← new
        Type:        "string",
        Title:       "API Secret",
        Description: "<provider> API secret.",
        WriteOnly:   true,
    },
},
Required: []string{
    "apiKey",
    "apiSecret",                                        // ← new
},
```

**Rules:**
- Secret fields (`*Secret`, `*Password`, `*PrivateKey`, `clientSecret`) MUST set `WriteOnly: true`. Otherwise the value leaks back to the UI.
- New required fields **break existing credentials**. If the credential is in production, surface this — usually a "with default" approach is needed (add field as optional + default + backfill, then make required in a later change).
- Use `github.com/blocknextai/go-packages/json` for `Default: json.RawMessage(...)` — its `RawMessage` is a stdlib type alias that satisfies `gjs.Schema.Default`. Do NOT import stdlib `encoding/json` (the real credentials all import `go-packages/json`; see `credential-create` skill).

#### Remove a field

Delete the property + the `Required` entry. Then propagate:

```bash
# Find every executor that reads the removed field:
grep -rn 'credential.String("<removed_field>")\|credential.Object("<removed_field>")' internal/nodeengine/nodes/
```

Each match is a real change site — the executor must stop reading it. Either replace with another field, fall back to a default, or delete the read entirely. Surface each match for user decision; don't auto-rewrite.

#### Rename a field

Rename in the `Properties` key, the `Required` slice, and propagate to every executor read. Use a search-and-replace via `Edit` with `replace_all` ONLY on the `credential.String("<old>")` / `credential.Object("<old>")` patterns — never blanket replace `<old>` since it might be a substring of unrelated identifiers.

```bash
grep -rn '"<old_field>"' internal/nodeengine/nodes/<provider>/
```

Inspect each result before applying.

#### Change OAuth2 scope

Edit the `<provider>Scope` const. If it's additive, existing tokens DO NOT automatically gain the new scope — they were issued for the old set. The user (or an admin) must re-authorize. Surface this prominently:

> Scope change is not retroactive. Existing access tokens were issued for the previous scope set; re-authorization is required for the new scope to take effect.

#### Add / remove a supported node

```go
SupportedNodes: &[]string{
    "linear.createissue",
    "linear.searchissues",  // ← new
},
```

When adding, verify the node ID exists or will be added in the same change. When removing, verify no node currently lists this credential in its `SupportedCredentials`.

```bash
# Verify the node ID exists:
grep -rn 'ID:.*"linear.searchissues"\|nodeID + ".searchissues"' internal/nodeengine/nodes/

# When removing — verify no node still uses this credential:
grep -rn 'SupportedCredentials.*"<credential_id>"' internal/nodeengine/nodes/
```

#### Update display metadata

`Name`, `Description`, `Icon` are independent and non-breaking.

#### Disable / enable

```go
Disabled: true,  // ← add field if not present
```

Disabled credentials are NOT registered (`RegisterCredential` skips them). All nodes referencing this credential in `SupportedCredentials` will silently lose their auth pathway. Surface this and check downstream:

```bash
grep -rn 'SupportedCredentials.*"<credential_id>"' internal/nodeengine/nodes/
```

### Step 3: Cross-reference propagation

After every edit, run the cross-check greps and report any inconsistencies. Do NOT silently leave them out.

| Change                            | Cross-check                                                                      |
| --------------------------------- | -------------------------------------------------------------------------------- |
| Field rename                      | `grep -rn 'credential.String("<old>")\|credential.Object("<old>")' internal/nodeengine/nodes/` |
| Field removal                     | Same as rename, but each match is a deletion candidate                           |
| `SupportedNodes` add              | Verify the node ID exists                                                        |
| `SupportedNodes` remove           | Verify no node lists this credential in `SupportedCredentials`                   |
| `Disabled: true`                  | List every node still listing this credential — they will break at runtime       |
| Scope change (OAuth2)             | Inform user re-auth is required for existing tokens                              |

### Step 4: Final verification

```bash
gofmt -w internal/nodeengine/credentials/<provider>_<authtype>.go
gofmt -w <any executor files modified>
go build ./...
golangci-lint run ./internal/nodeengine/...
```

All three must be clean.

## Critical rules (project memory)

1. **Never set `IsSupportPlatform` manually** — auto-derived by `PlatformConfigLoader`. If you saw old code that set it, that's debt — leave it removed.
2. **`IsOAuth1` / `IsOAuth2` are auto-set** from the credential ID substring matching. Do NOT set them manually.
3. **`WriteOnly: true` on every secret-bearing field.** Without it, the value leaks back to the UI.
4. **Use `go-packages/json` for `gjs.Schema.Default`** — `json.RawMessage` from `github.com/blocknextai/go-packages/json` (a stdlib type alias); never stdlib `encoding/json`, never a `utils` package. (Memory: `feedback_no_utils_package`.)
5. **Field renames break runtime contracts** — every executor that reads the credential needs the matching change in the same commit. Don't ship a partial rename.
6. **OAuth2 scope changes require re-authorization** — tokens are issued per-scope; existing tokens don't auto-upgrade.
7. **Disabling a credential silently kills nodes** that depend on it. Before disabling, confirm there are no live `SupportedCredentials` references.

## What NOT to do

- Do NOT use `Write` to overwrite the whole file — it risks losing unrelated lines that may have been added since you last read.
- Do NOT auto-rewrite executor `credential.String(...)` calls during a rename — surface each match for user decision; the field rename may have semantic implications (e.g., `apiKey` → `accessToken` reflects an OAuth2 migration, not a cosmetic rename).
- Do NOT remove a field from `Required` to "make it optional" without considering existing stored credentials — they were saved with the old required set; loosening here is fine, tightening is not.
- Do NOT add a new required field with no default to a deployed credential. Existing rows fail validation. Stage as: optional + default + backfill, then required.
- Do NOT modify `credentials.go` registry order during an update — registration is alphabetical; updates don't change the order.
- Do NOT touch `internal/credentials/` (the platform-side credentials module) from this skill — that's a different bounded context. This skill only operates on `internal/nodeengine/credentials/`.
