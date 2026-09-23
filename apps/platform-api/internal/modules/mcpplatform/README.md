# MCP Platform

> The platform's own MCP server: exposes the caller's BlockNext account so an MCP client can read and manage the platform on the user's behalf. Runs inside `mcp-api`.

## Responsibility
An HTTP adapter, not a bounded context with its own domain. It builds one MCP server whose tools read the signed-in user's profile and organizations, list workflows, triggers, executions and credentials, and turn triggers on or off. It always acts on the organization the access token was issued for and never returns credential secrets. It implements `mcp`'s `ServerProvider`, so `mcp` mounts, lists and protects it like every other MCP server. Every tool call is recorded through `executions`' `ToolInvocationService` with source `platform` (input, output, status, timing; no API key, no credentials); failures to record are logged, never surfaced. The server and its tools carry the `platform` brand icon (tools add a glyph each), which the UI expects under `public/assets/icons/brands/platform/{light,dark}.svg`. It owns no aggregates or tables.

## Tools
| Tool | Scope | Purpose |
| --- | --- | --- |
| `platform_get_me` | `platform:read` | The signed-in user, the active organization and every organization the user belongs to. |
| `platform_list_workflows` / `platform_get_workflow` | `platform:read` | Workflow catalog and one workflow with its canvas nodes. |
| `platform_list_executions` | `platform:read` | Recent workflow executions. |
| `platform_list_tool_invocations` | `platform:read` | History of tool calls made through the MCP servers (tool, status, error, timing). |
| `platform_list_credentials` | `platform:read` | Connected integrations (ids and titles only). |
| `platform_list_triggers` | `platform:read` | Schedule and webhook triggers (no webhook secrets). |
| `platform_set_trigger_active` | `platform:write` | Turn a trigger on or off; an inactive trigger starts no runs. |

## Use cases
- **Server** (`application/server`) — `NewServerProvider` builds the `mcpsdk.Server`, registering every tool together with its scope, and returns it with its catalog, `authMethods: ["oauth"]` and scopes `platform:read`, `platform:write`.

## HTTP API
Served by `mcp`:
- `POST|GET|DELETE /platform/mcp` — streamable MCP endpoint (stateless). Requires an MCP OAuth access token with `platform:read` or `platform:write`; each tool re-checks its own scope.
- Listed in `GET /servers`; metadata at `GET /.well-known/oauth-protected-resource/platform/mcp`.

## Dependencies
- **Bounded contexts:** `account` (`UserService`), `organizations` (`OrganizationService`, `OrganizationUserService`), `workflows` (`WorkflowService` via `workflows.NewServices`), `triggers` (`TriggerService` via `triggers.NewServices`), `executions` (`TaskExecutionService`, `ToolInvocationService`), `credentials` (`CredentialService`), `quota` (`CreditTransactionService`), `mcp` (the `servers.ServerProvider` contract it implements).
- **Infrastructure:** `modelcontextprotocol/go-sdk`. No database.

## Layout
Provider-only module: `application/server`, exposed as `Module.ServerProvider` in `module.go`; no routes of its own.
