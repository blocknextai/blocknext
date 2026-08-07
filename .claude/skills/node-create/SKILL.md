---
name: node-create
description: Use this skill when the user wants to add a new node (workflow action) to platform-api. Triggers on phrases like "add a node for X", "create a new node", "add a new action under <provider>", "let's add <provider>.<action>". Generates the node definition, executor, registers it (nodes + executors + functioncalling + MCP server), and ensures the credential link is correct.
---

> **Working directory:** all relative paths in this skill (`internal/...`, `cmd/...`, `go build ./...`) are relative to `apps/platform-api/` in the monorepo. Make targets run from the repo root.


# Create Node

Scaffolds a new workflow node under `internal/nodeengine/nodes/<provider>/<action>/`, wires it into the provider's `register.go` (which registers it in 4 places: nodes registry, executors registry, function calling, and MCP server), and verifies the credential it consumes exists.

## When to use

- User says "add an `<action>` node for `<provider>`" / "let's add `<provider>_<action>`" / "create a new node"
- A new third-party action is needed in the workflow editor

## Required information

Confirm with the user (ask only what's missing):

1. **Provider** — e.g., `mistral`, `notion`. Lowercase. Matches an existing folder under `internal/nodeengine/nodes/` if extending, or creates a new one.
2. **Action** — e.g., `chat`, `sendmessage`, `createrecord`. Lowercase, single word. Becomes the directory and the package name.
3. **Display name** — e.g., "Mistral Chat", "Notion Create Page".
4. **Description** — one-line description of what the node does at runtime.
5. **Category** — top-level category. Use one of the existing ones unless the user explicitly wants a new one:
   - `AI`, `Audio`, `Blockchain`, `CRM`, `Database`, `Ecommerce`, `Google Workspace`, `Image`, `Mailing`, `Publishing`, `System`, `Utility`, `Video`
6. **Sub-category** — usually the provider's display name (e.g., "Mistral", "Notion", "CoinGecko").
7. **Tags** — 3–7 lowercase keywords for search.
8. **Credential ID** — which credential this node consumes (e.g., `mistral_api`, `notion_oauth2`). Underscore separator. MUST already exist in `internal/nodeengine/credentials/`. If it doesn't, run the `credential-create` skill first.
9. **Input fields** — names, types, descriptions, optional defaults, optional enums. These become both the JSON schema and the executor input struct.
10. **Output fields** — names, types, descriptions. Shape of the **per-record** output map; the executor returns an array of these and the OutputSchema wraps it accordingly (see template below).
11. **MCP annotations group** — pick one based on what the tool does (see "Annotations classification" below).
12. **HTTP base URL + auth header** — for the helpers client (only if the node calls an HTTP API).

## NodeID format

Node IDs use **underscore separator**: `<provider>_<action>` (e.g., `facebook_publish_post`, `gemini_chat`). The MCP adapter strips the server-ID prefix from the nodeID when forming the tool name within an MCP server — this means a node's nodeID MUST start with `<serverID>_` for the MCP listing to look clean. Stick to the convention.

## File layout

```
internal/nodeengine/nodes/<provider>/
├── helpers/
│   └── client.go        # shared HTTP client; one per provider folder
├── <action>/
│   ├── node.go          # node definition (schema + metadata + annotations)
│   └── executor.go      # ExecuteWithContext logic
└── register.go          # wires nodes + executors + functioncalling + MCP server
```

When **adding a new action to an existing provider**, only create `<action>/node.go` and `<action>/executor.go`, then extend `register.go` (append a new action block AND add the new node to the provider's existing `mcp.RegisterServer(...)` Tools slice). Do NOT recreate `helpers/client.go`.

When **adding a new provider from scratch**, create the entire tree above, then add the provider's `Register()` call to `internal/nodeengine/nodes/nodes.go`.

## Icon and handles

A node states what it looks like and how it connects; nothing is derived from
its name or category.

**Icon** is a brand plus a glyph. `Brand` names the artwork under
`apps/platform/public/assets/icons/brands/<brand>/{light,dark}.svg` — reuse the
provider's existing brand (`gmail`, `slack`, `x`). A provider that has no mark
yet needs one dropped in that folder; no code registers it. `Glyph` names the
action badge under `.../glyphs/<glyph>.svg` — reuse an existing one where it
fits (`send`, `search`, `create`, `list`, `eye`, `trash`, `pencil`, `image`,
`film`, `speaker`, `music`, `chat`, `folder`, `file`, `table`, `record`,
`branch`, `clock`, `play`, `note`, `story`, `upload`, `organize`, `translate`,
`trending`, `target`). Two nodes of the same provider must not share a glyph —
that is the only thing telling them apart in a list.

A system primitive with no brand behind it declares only `Glyph`.

**Handles** are the connection points, and every node lists its own; there is no
default. Almost every node is one in, one out:

```go
Inputs: []nodes.NodeHandle{
	{Key: "in"},
},
Outputs: []nodes.NodeHandle{
	{Key: "out"},
},
```

A node with no way in declares `Inputs: []nodes.NodeHandle{}` (an empty, non-nil
slice — `nil` would mean "not declared"). A node that branches declares one
output per branch with a `Label`, which the canvas prints next to the dot:

```go
Outputs: []nodes.NodeHandle{
	{Key: "true", Label: "True"},
	{Key: "false", Label: "False"},
},
```

Edges record which output they leave from, so a key is part of the saved
workflow — renaming one breaks every flow that uses it.

## Annotations classification

Every node carries `Annotations: nodes.NodeAnnotations{...}` describing how the tool affects its environment. MCP clients (Claude Desktop, Cursor, etc.) use these hints to decide when to ask the user for confirmation. Pick one of these groups:

### Group 1 — Pure read (read-only, idempotent)
Tools that query external state without modifying anything. `gmail_search_emails`, `coingecko_price_monitor`, `*_get_*`, `*_read_*`, `*_search_*`, `*_list_*`.

```go
Annotations: nodes.NodeAnnotations{
    ReadOnly:   true,
    Idempotent: true,
},
```

### Group 2 — Closed-world compute (no external interaction)
Pure local computation, no external API. `qrcode_generate`, `system_starter`, `system_sleep`, `system_condition`.

```go
Annotations: nodes.NodeAnnotations{
    ReadOnly:    true,
    Idempotent:  true,
    Destructive: new(false),
    OpenWorld:   new(false),
},
```

### Group 3 — Non-destructive creators (additive write)
Creates new resources without modifying or deleting existing ones. Publish, send, create, upload, subscribe. `facebook_publish_post`, `gmail_send_email`, `mailchimp_subscribe`, `*_create_*`, `*_upload_*`.

```go
Annotations: nodes.NodeAnnotations{
    Destructive: new(false),
},
```

### Group 4a — LLM / generation (non-deterministic compute, external API)
Calls external API to generate content, doesn't modify user's environment, but non-deterministic. `gemini_chat`, `anthropic_chat`, `gemini_imagen`, `veo_veo3`, `elevenlabs_text_to_speech`, `sunomusic_create_music`.

```go
Annotations: nodes.NodeAnnotations{
    ReadOnly:    true,
    Destructive: new(false),
},
```

### Group 4b — LLM / generation (deterministic)
Same as 4a but deterministic for same inputs. `deepl_translate`.

```go
Annotations: nodes.NodeAnnotations{
    ReadOnly:    true,
    Idempotent:  true,
    Destructive: new(false),
},
```

### Group 5 — Explicit destructive (irreversible / overwrite / delete)
Tools that modify or delete existing state. Always explicit even though `Destructive: new(true)` matches MCP default — documents intent. `ethereum_send_transaction`, `*_delete_*`, `*_update` (when overwriting), `gmail_organize_emails`.

```go
Annotations: nodes.NodeAnnotations{
    Destructive: new(true),
},
```

**Default per MCP spec is `Destructive: true`** — but every node in this project sets annotations explicitly for safety + grep-ability. Never leave `Annotations` unset.

## Templates

### `<action>/node.go`

```go
package <action>

import (
	"github.com/blocknextai/go-packages/json"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type <Provider><Action>Node struct {
	nodes.Node
}

func New<Provider><Action>Node(nodeID string) *<Provider><Action>Node {
	return &<Provider><Action>Node{
		Node: nodes.Node{
			ID:          nodeID,
			Version:     "0.0.1",
			Name:        "<DisplayName>",
			Description: "<one-line description>",
			Icon: nodes.NodeIcon{
				Brand: "<brand>",
				Glyph: "<glyph>",
			},
			Inputs: []nodes.NodeHandle{
				{Key: "in"},
			},
			Outputs: []nodes.NodeHandle{
				{Key: "out"},
			},
			Categories:    []string{"<Category>"},
			SubCategories: []string{"<SubCategory>"},
			Tags: []string{
				"<tag1>",
				"<tag2>",
				// ...
			},
			SupportedCredentials: []string{
				"<provider>_<authtype>",
			},
			InputSchema: &gjs.Schema{
				Type: "object",
				Properties: map[string]*gjs.Schema{
					"<field>": {
						Type:        "string",
						Title:       "<Title>",
						Description: "<description>",
						// Default: json.RawMessage(`"<default>"`), // optional
						// Enum: []any{"a", "b"},                  // optional
					},
					// ...
				},
				Required: []string{
					"<field>",
					// ...
				},
			},
			// OutputSchema MUST wrap in array — executors return []map[string]any,
			// MCP envelopes it as { items: [...] } at runtime.
			OutputSchema: &gjs.Schema{
				Type: "array",
				Items: &gjs.Schema{
					Type: "object",
					Properties: map[string]*gjs.Schema{
						"<field>": {
							Type:        "string",
							Description: "<description>",
						},
						// ...
					},
				},
			},
			Annotations: nodes.NodeAnnotations{
				// Pick one group from "Annotations classification" above.
				Destructive: new(false),
			},
			HasNaturalLanguage: true, // set false for nodes that should NOT appear in LLM function calling
		},
	}
}
```

Rules:
- `Version` starts at `"0.0.1"` for every new node. Bump via SemVer when shape changes: patch for non-breaking fixes, minor for additive input/output, major for breaking schema changes.
- `Icon.Light` and `Icon.Dark` are both set to `nodeID`. Do not hardcode strings.
- `SupportedCredentials` is **always** a multi-line slice literal — even when there is exactly one element. (Project convention; see memory `feedback_supported_credentials_multiline`.)
- For nodes that need no credentials (e.g., `system_sleep`, `qrcode_generate`, `rss_read_feed`), omit `SupportedCredentials` entirely.
- `OutputSchema` always wraps in `{Type: "array", Items: {...}}` — the inner object is the per-record shape. Executors return arrays; this matches reality + MCP wraps in an `items` envelope at runtime.
- `Annotations` is **always set explicitly** — never leave it unset. Pick the right group; if truly destructive, write `Destructive: new(true)` for documentation.
- Numeric fields use `Type: "number"` (not `"integer"` unless the API truly requires it). Numeric defaults are `json.RawMessage(`1024`)` style — raw numbers, no quotes.
- For enum fields, list values in `Enum: []any{...}` and pick a `Default`.
- Use `Minimum: new(0.0)` / `Maximum: new(1.0)` for bounded numerics (see `anthropic/chat/node.go` for the temperature example).

### `<action>/executor.go`

```go
package <action>

import (
	"context"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/<provider>/helpers"
)

var (
	ErrExecutorEmptyResponse = apperror.Internal("empty response")
)

type <Provider><Action>ExecutorInput struct {
	<Field1> string `schema:"<field1>"`
	<Field2> int    `schema:"<field2>"`
	// ...
}

type <Provider><Action>Executor struct {
	executors.Executor
	validator *jsonschema.Validator[<Provider><Action>ExecutorInput]
}

func New<Provider><Action>Executor(
	nodeID string,
	validator *jsonschema.Validator[<Provider><Action>ExecutorInput],
) *<Provider><Action>Executor {
	return &<Provider><Action>Executor{
		Executor: executors.Executor{
			ID: nodeID,
		},
		validator: validator,
	}
}

type SuccessResponse struct {
	// shape of provider's success body
}

type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (e *<Provider><Action>Executor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		credential := nodeEngineDomainCredentials.GetCredentials(credentials, "<provider>_<authtype>")
		// API key credentials:
		apiKey := credential.String("apiKey")
		client := helpers.CreateClient(ctx, apiKey)
		// OAuth2 credentials use a nested object instead — replace the two lines above with:
		//   oauthToken := credential.Object("oauthToken")
		//   accessToken := oauthToken.String("accessToken")
		//   client := helpers.CreateClient(ctx, accessToken)

		results := make([]map[string]any, 0)
		for _, item := range data {
			input, err := e.validator.Parse(item)
			if err != nil {
				return nil, err
			}

			var successResponse SuccessResponse
			var errorResponse ErrorResponse
			response, err := client.Post("/<endpoint>").
				JSONContentType().
				Body(map[string]any{
					"<field>": input.<Field1>,
				}).
				Do(&successResponse, &errorResponse)

			if err != nil {
				return nil, err
			}

			if !response.IsSuccess() {
				return nil, apperror.Internal(errorResponse.Error.Message)
			}

			results = append(results, map[string]any{
				"<output>": successResponse.<Field>,
			})
		}

		return results, nil
	}
}
```

Rules (project memory):
- **HTTP**: use `httpclient` from `helpers.CreateClient(...)`. NEVER import `net/http` directly. NEVER use `fmt.Sprintf` for URL building — use the chained client methods (`.Get`, `.QueryParam`, `.Header`, `.Body`).
- **Initialism casing**: `JSONContentType()` not `JsonContentType()`; `URL` not `Url`; `ID` not `Id`. Memory: `feedback_initialism_rename_autonomous`.
- **Errors**: return `apperror.Internal(...)` (or another apperror Kind) for failures. Sentinel errors live as `var Err... = apperror.Internal(...)` at the top of the file.
- **JSON**: `github.com/blocknextai/go-packages/json` (`json.Marshal` / `json.Unmarshal` / `json.RawMessage`) everywhere — both `node.go` and `executor.go`. Do NOT use stdlib `encoding/json` or a `utils` package. `go-packages/json`'s `RawMessage` is a type alias of the stdlib type, so it satisfies `gjs.Schema` defaults in `node.go` too.
- **Strings**: use `strings.Builder` for hot-path concatenation. `strings.TrimSpace(s) == ""` for empty-string checks.
- **Logging**: use `slog` (structured) — never `log` or `fmt.Print*`.
- **Slices**: prefer `result = append(result, x)` over `result[i] = x`.
- **Iteration**: the executor receives a slice `data []map[string]any` (fan-out from the runtime). Loop and append; do not assume length 1.
- **Context**: respect `ctx.Done()` at the top of `ExecuteWithContext`. Each HTTP call already inherits the context via `httpclient.NewClientBuilder().Context(ctx)`.

### `helpers/client.go` (only when adding a new provider)

```go
package helpers

import (
	"context"

	"github.com/blocknextai/go-packages/httpclient"
)

func CreateClient(ctx context.Context, apiKey string) *httpclient.Client {
	return httpclient.NewClientBuilder().
		Context(ctx).
		BaseURL("<provider base url, e.g. https://api.example.com/v1>").
		Header("Authorization", "Bearer "+apiKey). // adjust per provider auth scheme
		Build()
}
```

If the provider uses a non-Bearer scheme (e.g., Anthropic uses `x-api-key`), swap headers accordingly. Some providers also need a version header (`anthropic-version`).

For OAuth2-authenticated nodes, the `apiKey` parameter is replaced by the access token retrieved via the nested `oauthToken` object in the executor; the helper just takes a `token string` instead.

For nodes that upload/download files via `filegateway`, the helper signature accepts `fileGatewayService filegateway.FileGateway` and the provider's `Register()` is updated to receive it (see `gemini` or `discord`).

### `register.go` (provider-level)

When **creating a new provider**, write `internal/nodeengine/nodes/<provider>/register.go`. The register call wires the node into 4 registries: nodes, executors, function calling, AND MCP server.

```go
package <provider>

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/functioncalling"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/mcp"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/<provider>/<action>"
)

func Register() {
	nodeID := "<provider>"

	<action>NodeID := nodeID + "_<action>"
	<action>Node := <action>.New<Provider><Action>Node(<action>NodeID)
	<action>Validator := jsonschema.New[<action>.<Provider><Action>ExecutorInput](<action>Node.GetInputSchema())
	<action>Executor := <action>.New<Provider><Action>Executor(<action>NodeID, <action>Validator)

	nodes.RegisterNode(<action>Node)
	executors.RegisterExecutor(<action>Executor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(<action>Node))

	mcp.RegisterServer(&mcp.Server{
		ID:          nodeID,
		Name:        "<Provider Display Name>",
		Description: "Tools for <what the provider does>.",
		Version:     "0.0.1",
		Tools: []nodes.NodeManager{
			<action>Node,
		},
	})
}
```

When **adding a new action to an existing provider**:
1. Add a new action block matching the existing pattern (Node + Validator + Executor + 3× Register calls).
2. **Append the new node to the existing `mcp.RegisterServer` Tools slice** — do NOT add a second `RegisterServer` call. There is exactly one MCP server per provider.

### `nodes/nodes.go` (only for new providers)

If the provider is new, append `<provider>.Register()` to the body of `Register(fileGatewayService filegateway.FileGateway)` in `internal/nodeengine/nodes/nodes.go`. Add the corresponding import. The list is intentionally not strictly alphabetized — append at the end.

If the provider needs `filegateway` (file uploads/downloads), its `Register` signature accepts the gateway: `Register(fileGatewayService)` (see `gemini.Register`, `drive.Register`).

## Function calling

`functioncalling.RegisterFunctionCalling(functioncalling.Generate(node))` makes the node available to LLM nodes for tool use. Always include this call. To opt a node out of LLM function calling, set `HasNaturalLanguage: false` on the node — `functioncalling.Generate(node)` checks this internally and yields a disabled FunctionCalling entry that the registry filters out.

## MCP exposure

Every node is automatically exposed as an MCP tool when its provider's `mcp.RegisterServer(...)` includes it in `Tools`. The MCP adapter:
- Reads `node.GetAnnotations()` and forwards to MCP `ToolAnnotations`
- Wraps `OutputSchema` in `{type: "object", properties: {items: <node-output-schema>}}` so SDK accepts it (SDK rejects top-level non-object schemas)
- Resolves `credential:<scope>:<uuid>` references from the tool's `credentials` parameter

No additional work in node.go / executor.go is required to expose to MCP — register.go's `mcp.RegisterServer` call is the only touchpoint.

## Critical rules

1. **Credential must exist** before the node is registered. Verify the credential ID listed in `SupportedCredentials` is present in `internal/nodeengine/credentials/credentials.go`. If not, run the `credential-create` skill first.
2. **`SupportedCredentials` is always multi-line.**
3. **`OutputSchema` always wraps in array** — `{type: "array", items: {type: "object", properties: {...}}}`. Domain truth: executors return `[]map[string]any`.
4. **`Annotations` is always set explicitly** — pick one of 5 groups. Never leave unset.
5. **NodeID uses underscore**: `<provider>_<action>`, not dot.
6. **Don't use `fmt`** for performance and consistency: use `strings.Builder`, `strconv`, `errors.New`, `slog`.
7. **Use `go-packages/json` everywhere** — `json.Marshal` / `json.Unmarshal` in `executor.go`, `json.RawMessage` for `gjs.Schema` defaults in `node.go`. Never stdlib `encoding/json`, never a `utils` package. (Memory: `feedback_no_utils_package`.)
8. **Don't reach into other modules' repositories.** Service interfaces are the only allowed cross-module dependency. Within `nodeengine`, helpers and the domain interfaces are the boundary.
9. **Default to no comments.** Identifier names should explain intent. Add a comment only for non-obvious WHY (provider quirk, undocumented response shape, etc.).

## Final steps

After all files exist:

1. Run `gofmt -w` on every new/modified `.go` file.
2. Run `go build ./...` — must compile.
3. Run `golangci-lint run ./internal/nodeengine/...` — fix issues; lint catches things `go build` misses.
4. Manually verify:
   - The node ID in `Register()` matches the credential's `SupportedNodes` list (in the credential file). If not, append it there.
   - Tags are lowercase and contain no duplicates.
   - `Required` includes every field the runtime should guarantee is populated (the validator fills defaults). Truly optional fields (no default, safe to omit — e.g., `description`) stay out. A field with a `Default` is usually still in `Required`.
   - For new providers: `mcp.RegisterServer(...)` includes ALL of the provider's nodes in `Tools`.
   - For existing providers gaining a new action: the new node was appended to the existing `mcp.RegisterServer` Tools slice.

## What NOT to do

- Do NOT create a separate file per node action variation under one folder; one folder = one action.
- Do NOT register a node from `node.go` or `executor.go` directly — registration only happens in `register.go`.
- Do NOT manually marshal/unmarshal HTTP bodies — let `httpclient.Body(map[string]any{...}).Do(&success, &error)` handle it.
- Do NOT swallow errors with `_ = err`. If an error truly cannot affect outcome (rare), explain why with a one-line comment.
- Do NOT reference `$<something>` patterns or workflow JSON in the node — those are runtime concerns handled by the executor framework, not the node implementation.
- Do NOT add a second `mcp.RegisterServer` call for an existing provider — extend the existing one's Tools slice.
- Do NOT use dot separator in nodeID (`<provider>.<action>` is WRONG); always underscore.
