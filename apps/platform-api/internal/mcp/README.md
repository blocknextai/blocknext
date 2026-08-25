# MCP

> Exposes the platform's nodeengine servers/tools over the Model Context Protocol (MCP), letting external MCP clients discover servers and invoke node executors as tools.

## Responsibility
This context is an HTTP adapter, not a bounded context with its own domain. It takes the MCP servers defined by `nodeengine` and builds one streamable MCP endpoint per server, each exposing that server's nodes as MCP tools. When a tool is called it resolves credential references and runs the node's executor. It also serves a public listing of available servers. It owns no aggregates or database tables.

## How a node becomes an MCP tool

**Boot-time build, zero per-node code.** At startup `infrastructure.RegisterInfrastructure` walks every server the nodeengine providers registered via `mcp.RegisterServer` and calls `adapter.Build` on each: one `mcpsdk.Server` per provider, wrapped in a **stateless** `NewStreamableHTTPHandler` and mounted at `/:serverId/mcp`. Adding a node to a provider's MCP `Tools` slice is the entire integration — the tool surface is derived, never hand-written. Build fails fast if a tool has no matching executor (`ErrExecutorNotFound`), so a half-registered node cannot ship. The `mcp-api` binary serves this on its own port with a trimmed module graph (no workflows/executions/taskrunner); `GET /servers` advertises each endpoint by substituting `{serverId}` into `MCP_SERVER_URL_TEMPLATE`.

**Per-tool translation** (`application/adapter/schema.go`), all derived from the node descriptor:

| MCP tool field | Derived from |
| - | - |
| `Name` | Node ID (e.g. `slack_send_message`) |
| `Title` / `Description` | Node `Name` / `Description` |
| `Annotations` | `NodeAnnotations` mapped 1:1 to `ToolAnnotations` (read-only/destructive/idempotent/open-world hints) |
| `InputSchema` | Node `InputSchema` **augmented** with a required `credentials` object |
| `OutputSchema` | Node `OutputSchema` (an array of records) wrapped as `{type: "object", properties: {items: <original>}, required: ["items"]}` |

**Credentials are references, never secrets.** For each entry in the node's `SupportedCredentials`, the injected `credentials` object gains a required string property with pattern `^credential:organization:[0-9a-fA-F-]+$`. MCP clients pass opaque references to credentials already stored on the platform; no secret material ever crosses the MCP boundary inbound.

**Invocation path.**
1. The API-key middleware authenticates the call (`Authenticate()` + `RequireScope(ScopeMCPInvoke)` from the `apikeys` module) and injects owner-type/owner-id headers; the tool handler refuses to run without them (`ErrMissingOwner`) and rejects any owner that is not an organization (`ErrOrganizationOwnerRequired`).
2. Arguments are split: the `credentials` object becomes references, the rest becomes the tool's data record.
3. `credentialresolver` parses each reference (`common/domain/credential.ParseReference`), enforces that the reference scope matches the API key's owner type (`ErrCredentialScopeMismatch`), then calls `RegenerateTokenIfNeeded`, so OAuth tokens are refreshed transparently per invocation.
4. The node's executor runs with the standard contract, the single argument object passed as a one-item data list: `ExecuteWithContext(ctx, credentials, []map[string]any{data})`.
5. The result is returned twice per the MCP spec: `StructuredContent: {items: [...]}` (matching the wrapped output schema) plus one JSON `TextContent` block per output record. Executor failures come back as `isError` tool results, not protocol errors.

## Use cases (application)
- **GetAllServers** (`application/getallservers`) — return the catalog of MCP servers (built from `ServerURLTemplate` + `ServerService`).
- **Adapter** (`application/adapter`) — builds an `mcpsdk.Server` per nodeengine server: registers each node as an MCP tool, augments the input schema with a required `credentials` object, wraps outputs under `items`, and on invocation extracts owner from request headers, resolves credentials, then runs the executor.
- **History** (`application/history`) — records every finished tool call as a `ToolInvocation` (source `mcp`) carrying the credential-free parameters, credential references and outputs; failures to record are logged, never surfaced to the caller.
- **CredentialResolver** (`application/credentialresolver`) — parses `credential:organization:<uuid>` references, enforces the owner is an organization, and regenerates OAuth tokens as needed.

## HTTP API
- `GET /servers` — list MCP servers; cached 5 min; no auth.
- `POST|GET|DELETE /:serverId/mcp` — streamable MCP endpoint per server (stateless). Requires API-key `Authenticate()` + `RequireScope(ScopeMCPInvoke)`. Owner is taken from the auth-injected owner-type/owner-id headers.

## Dependencies
- **Bounded contexts:** nodeengine (`mcp.ServerService`, `executors.ExecutorService`), credentialoauth (`CredentialOAuthTokenRegenerateService`), executions (`ToolInvocationService` for call history), apikeys (`ScopeMCPInvoke` for the required scope).
- **Infrastructure:** `modelcontextprotocol/go-sdk` (streamable HTTP MCP servers); Fiber cache middleware; API-key auth middleware. No database.

## Layout
Adapter-style module (no `domain/` layer): `application/` (adapter, credentialresolver, getallservers), `infrastructure/` (builds per-server handlers), `presentation/` (routes), wired in `module.go`.
