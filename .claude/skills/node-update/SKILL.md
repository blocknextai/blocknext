---
name: node-update
description: Use this skill when the user wants to modify an existing workflow node under `internal/nodeengine/nodes/<provider>/<action>/`. Triggers on phrases like "update the X node", "add a new input field to Y", "change Z's category / tags / annotations", "add a credential to the node", "disable a node", "update the node's HTTP endpoint", "rename a node field". Edits node.go and/or executor.go in place, keeps the input struct + validator in sync with the schema, and propagates cross-references (credentials' `SupportedNodes`, MCP server's Tools slice).
---

> **Working directory:** all relative paths in this skill (`internal/...`, `cmd/...`, `go build ./...`) are relative to `apps/platform-api/` in the monorepo. Make targets run from the repo root.


# Node Update

Edits an existing workflow node under `internal/nodeengine/nodes/<provider>/<action>/` (`node.go` and/or `executor.go`) and keeps the input schema, validator, executor input struct, credential cross-references, and MCP server consistent.

## When to use

- User says "update the `<provider>_<action>` node" / "add a `<field>` input to X" / "change Y's category" / "add `<credential>` to node's supported credentials" / "rename a field in Z node" / "disable the X node" / "change Y's annotations"
- For brand-new nodes use `node-create` instead.

## Required information

Confirm with the user (ask only what's missing):

1. **Node ID** — `<provider>_<action>` (e.g., `anthropic_chat`, `linear_createissue`). Underscore separator. Files live at `internal/nodeengine/nodes/<provider>/<action>/`.
2. **Change type** — one or more of:
   - **Add input field** (schema + validator input struct + executor uses)
   - **Remove input field** (schema + struct + every executor read)
   - **Rename input field** (schema key + struct field + `schema:"..."` tag + executor reads)
   - **Add output field** (output schema's inner `Items` properties + executor's `results = append(results, map{...})` blocks)
   - **Change input default / enum** (no struct change required if type unchanged)
   - **Change metadata** (Name, Description, Icon, Categories, SubCategories, Tags)
   - **Add / remove a supported credential** (`SupportedCredentials` slice)
   - **Change annotations** (re-classify the node into one of the 5 groups; see `node-create` skill)
   - **Disable** (`Disabled: true`) — node hidden everywhere (nodes registry, executors registry, MCP server tools list, function calling)
   - **Toggle natural language** (`HasNaturalLanguage: true/false`) — controls function calling exposure
   - **Change HTTP endpoint / API base URL** (helpers/client.go or executor's `client.Post("/...")`)
3. **Field details** (for add / change) — name (camelCase JSON), type, Title, Description, Default, Enum, Min/Max bounds, required-or-optional.
4. **Output field details** (for add) — name (camelCase JSON), type, Description. Added to `OutputSchema.Items.Properties`, not at the top level (OutputSchema is `{Type: "array", Items: {...}}`).

## Workflow

### Step 1: Locate and read

```
internal/nodeengine/nodes/<provider>/<action>/
├── node.go        # schema + metadata + annotations
└── executor.go    # input struct + validator + Execute logic
```

Read both before editing. The provider's `register.go` and `helpers/client.go` are upstream from per-action edits — don't touch them unless the change actually requires it (e.g., HTTP base URL change → `helpers/client.go`; new action added → `register.go`).

### Step 2: Apply the change — `node.go`

Use **`Edit`**, not `Write`. Apply the smallest possible diff per change.

#### OutputSchema reminder

OutputSchema is **always** array-wrapped:

```go
OutputSchema: &gjs.Schema{
    Type: "array",
    Items: &gjs.Schema{
        Type: "object",
        Properties: map[string]*gjs.Schema{
            "field1": {...},
            "field2": {...},
        },
    },
},
```

When adding/removing/renaming output fields, edit inside `Items.Properties`, not at the top level.

#### Add an input field

Three edits in lockstep:
1. Add property to `InputSchema.Properties`.
2. Append to `InputSchema.Required` if it has no `Default` and is not optional.
3. Add the field to the executor's `<Provider><Action>ExecutorInput` struct with the matching `schema:"<jsonName>"` tag.

```go
// node.go — InputSchema.Properties:
"newField": {
    Type:        "string",
    Title:       "New Field",
    Description: "<description>",
    Default:     json.RawMessage(`""`),
},
```

```go
// executor.go — Input struct:
type <Provider><Action>ExecutorInput struct {
    // ...existing fields...
    NewField string `schema:"newField"`
}
```

If the field has no default and is not optional, add it to `Required` AND to the validator-checked input struct (the validator will fail at runtime if the field is absent). For optional fields with sensible Go zero-value defaults, omit from `Required`.

#### Remove an input field

Four edits:
1. Remove from `InputSchema.Properties`.
2. Remove from `InputSchema.Required`.
3. Remove from the executor's input struct.
4. Remove every read of the field from `executor.go`'s body (the `input.<Field>` references).

If the executor uses the field in HTTP body construction or branching logic, the removal is semantic — surface the change site for user decision rather than auto-deleting.

#### Rename an input field

Four edits:
1. Rename the key in `InputSchema.Properties`.
2. Rename in `InputSchema.Required` if listed.
3. Rename the struct field AND its `schema:"..."` tag.
4. Replace every `input.<OldName>` reference in `executor.go` with `input.<NewName>`.

The `schema:"..."` tag must match the JSON key in the schema — they're tied via the validator.

#### Add an output field

Two edits:
1. Add property to `OutputSchema.Items.Properties` in `node.go` (NOT `OutputSchema.Properties` — output is array-wrapped).
2. Add `"newField": value` to **every** `results = append(results, map[string]any{...})` block in `executor.go`.

If the value isn't already computed, the executor needs new logic — surface this to the user.

#### Change defaults / enums

One edit in `node.go` — defaults and enums are runtime metadata, no struct/validator change needed:

```go
"model": {
    Type:    "string",
    Title:   "Model",
    Enum:    []any{"new-model-v1", "new-model-v2"},     // ← updated
    Default: json.RawMessage(`"new-model-v1"`),         // ← updated
},
```

Verify the Go type still matches the new values (e.g., changing from `int` to `float64` in the enum requires a struct field type change).

#### Change metadata

`Name`, `Description`, `Icon`, `Categories`, `SubCategories`, `Tags` are independent. Tags must remain lowercase; categories must use one of the existing top-level values (memory: `node-create` skill lists them).

#### Add / remove a supported credential

```go
SupportedCredentials: []string{
    "anthropic_api",
    "anthropic_bedrock",  // ← new
},
```

Always multi-line, even with one element (memory: `feedback_supported_credentials_multiline`). Underscore separator in credential IDs.

After adding, verify the credential file exists. After removing, verify no executor still reads from that credential ID:

```bash
# Confirm the credential exists:
grep -n 'ID:.*"anthropic_bedrock"' internal/nodeengine/credentials/

# After removing — confirm executor no longer reads it:
grep -n 'GetCredentials(credentials, "<removed_credential>")' internal/nodeengine/nodes/<provider>/<action>/executor.go
```

#### Change annotations

Re-classify the node into one of the 5 groups (see `node-create` skill's "Annotations classification" section). Direct edit of `Annotations: nodes.NodeAnnotations{...}` block. Don't leave annotations unset.

```go
// Before (Group 3 — non-destructive creator):
Annotations: nodes.NodeAnnotations{
    Destructive: new(false),
},

// After re-classification to Group 5 (explicit destructive — overwrite):
Annotations: nodes.NodeAnnotations{
    Destructive: new(true),
},
```

Surface the implication to the user: changing destructive/readOnly hints affects how LLMs/MCP clients treat the tool (confirmation prompts, etc.).

#### Disable

```go
Disabled: true,  // node hidden + skipped from registration
```

Disabled nodes are filtered out by all four registries: `nodes.RegisterNode`, `executors.RegisterExecutor`, and `mcp.RegisterServer` (which calls `filterEnabledTools` on the Tools slice). The function calling registry also skips disabled nodes.

The flag is independent of any "coming soon" UI state — there is no separate `IsComingSoon` field; presentation handles that.

#### Toggle natural language (function calling exposure)

```go
HasNaturalLanguage: false,  // hide from LLM function calling
```

When `HasNaturalLanguage: false`, `functioncalling.Generate(node)` marks the entry as disabled and the registry filters it out — no separate edit to `register.go` needed. The node still appears in MCP tools and the workflow editor, just not in LLM tool catalogs.

### Step 3: Apply the change — `executor.go`

For input/output schema changes propagated to the executor, the work is mechanical (add/remove struct fields, add/remove `results[]map` keys). For semantic changes (HTTP endpoint, request body shape, response parsing), surface each change site rather than auto-rewriting.

If the executor's HTTP flow needs to change (different endpoint, different auth header, response shape), confirm the new shape with the user before editing.

### Step 4: Cross-reference propagation

| Change                                | Cross-check                                                                       |
| ------------------------------------- | --------------------------------------------------------------------------------- |
| Add input field                       | Add to executor input struct in same commit                                       |
| Remove input field                    | Remove every `input.<Field>` read in executor.go                                  |
| Rename input field                    | Update struct field + `schema:"..."` tag + every executor read                    |
| Add output field                      | Edit inside `OutputSchema.Items.Properties` (not top level); add to every `results = append(results, ...)` in executor.go |
| Change annotations                    | Pick correct group from `node-create` skill; surface UX implication                |
| Add supported credential              | Verify credential file exists; consider adding node to credential's `SupportedNodes` |
| Remove supported credential           | Verify executor no longer reads from that credential                              |
| Disable node                          | Tell user the node will disappear from UI, function calling, AND MCP server tools list |
| Toggle HasNaturalLanguage             | No `register.go` edit needed — function calling registry auto-filters             |
| Rename node ID (rare; breaks clients) | Update credential `SupportedNodes` references; warn re: stored workflows; tool name in MCP also changes |

#### Update credential's `SupportedNodes` (optional but recommended)

When a node gains or loses a credential, the credential's `SupportedNodes` list often needs the inverse update. The skill should surface the suggested edit:

```bash
grep -n 'SupportedNodes' internal/nodeengine/credentials/<provider>_<authtype>.go
```

Show the user the credential file's current list and ask whether to add/remove the node ID. The list is metadata used by the UI for credential pickers — keeping it consistent prevents UI surprise.

### Step 5: Final verification

```bash
gofmt -w internal/nodeengine/nodes/<provider>/<action>/
gofmt -w <any other modified files>
go build ./...
golangci-lint run ./internal/nodeengine/...
```

All three must be clean.

## Critical rules (project memory)

1. **Schema and Input struct must stay in sync.** Adding a schema property without the matching struct field means the validator silently ignores it. Adding a struct field without the schema means the field never gets populated. The validator is the bridge — both ends must match.
2. **`SupportedCredentials` is always multi-line** (memory: `feedback_supported_credentials_multiline`).
3. **Tags are lowercase, no duplicates.**
4. **`OutputSchema` is always array-wrapped.** Output fields go inside `OutputSchema.Items.Properties`, never `OutputSchema.Properties` directly.
5. **`Annotations` is always set explicitly** — one of 5 groups (see `node-create` skill).
6. **`go-packages/json` everywhere** — `json.RawMessage` from `github.com/blocknextai/go-packages/json` for `gjs.Schema` defaults in node.go (its `RawMessage` is a stdlib type alias, so it works); never stdlib `encoding/json`.
7. **`json.Marshal` / `json.Unmarshal` from `go-packages/json` in executor.go** — never `encoding/json`, never a `utils` package. (Memory: `feedback_no_utils_package`.)
8. **`httpclient` only in helpers, never `net/http` directly.**
9. **`JSONContentType()` (initialism casing)** — never `JsonContentType()`.
10. **`apperror.Internal(...)` for executor errors,** not `fmt.Errorf` or `errors.New`.
11. **Disabling a node makes it invisible to:** UI, function calling, AND MCP server tools list. Confirm intentional.

## What NOT to do

- Do NOT use `Write` to rewrite either file — apply Edits surgically.
- Do NOT bump the file's package, ID, or directory. Renaming a node ID is a much larger change (every workflow that uses it breaks, every MCP client that references it breaks); a separate skill should handle that, not this one.
- Do NOT auto-rewrite the executor body for semantic changes (HTTP endpoint, request body, response shape). Surface each change site for user decision.
- Do NOT add a new tag without checking duplicates against the existing list — duplicate tags are a lint failure waiting to happen.
- Do NOT touch `register.go` or `helpers/client.go` for changes that are scoped to one action's `node.go` / `executor.go`. Changes to those files are upstream and affect all sibling actions.
- Do NOT modify `nodes.go` — the global Register list is for adding/removing whole providers, not individual actions.
- Do NOT add output fields at the top level of `OutputSchema` — they go inside `Items.Properties`.
- Do NOT leave `Annotations` unset — explicit classification is the project standard.
- Do NOT reference fields that don't exist on the Node struct (`IsComingSoon`, `FunctionCallingDisabled` — neither exists; use `Disabled` or `HasNaturalLanguage` instead).
