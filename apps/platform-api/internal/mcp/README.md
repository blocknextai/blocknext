# MCP

> Exposes the platform's nodeengine servers/tools over the Model Context Protocol (MCP), letting external MCP clients discover servers and invoke node executors as tools.

## Responsibility
This context is an HTTP adapter, not a bounded context with its own domain. It takes the MCP servers defined by `nodeengine` and builds one streamable MCP endpoint per server, each exposing that server's nodes as MCP tools. When a tool is called it resolves credential references and runs the node's executor. It also serves a public listing of available servers. It owns no aggregates or database tables.

## Use cases (application)
- **GetAllServers** (`application/getallservers`) — return the catalog of MCP servers (built from `ServerURLTemplate` + `ServerService`).
- **Adapter** (`application/adapter`) — builds an `mcpsdk.Server` per nodeengine server: registers each node as an MCP tool, augments the input schema with a required `credentials` object, wraps outputs under `items`, and on invocation extracts owner from request headers, resolves credentials, then runs the executor.
- **CredentialResolver** (`application/credentialresolver`) — parses `credential:(organization|user):<uuid>` references, enforces the credential scope matches the owner type, and regenerates OAuth tokens as needed.

## HTTP API
- `GET /servers` — list MCP servers; cached 5 min; no auth.
- `POST|GET|DELETE /:serverId/mcp` — streamable MCP endpoint per server (stateless). Requires API-key `Authenticate()` + `RequireScope(ScopeMCPInvoke)`. Owner is taken from the auth-injected owner-type/owner-id headers.

## Dependencies
- **Bounded contexts:** nodeengine (`mcp.ServerService`, `executors.ExecutorService`), credentialoauth (`CredentialOAuthTokenRegenerateService`), apikeys (`ScopeMCPInvoke` for the required scope).
- **Infrastructure:** `modelcontextprotocol/go-sdk` (streamable HTTP MCP servers); Fiber cache middleware; API-key auth middleware. No database.

## Layout
Adapter-style module (no `domain/` layer): `application/` (adapter, credentialresolver, getallservers), `infrastructure/` (builds per-server handlers), `presentation/` (routes), wired in `module.go`.
