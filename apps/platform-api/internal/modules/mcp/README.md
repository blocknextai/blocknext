# MCP

> Hosts every BlockNext MCP server over the Model Context Protocol: the nodeengine tool servers and the servers other contexts provide, such as the platform account server from `mcpplatform`.

## Responsibility
This context is an HTTP adapter, not a bounded context with its own domain. It collects servers from `ServerProvider`s, mounts one stateless streamable MCP endpoint per server, serves a single catalog of all of them with their RFC 9728 protected resource metadata, and applies each server's authentication policy. Its own provider turns `nodeengine` servers into MCP servers whose tools resolve credential references and run node executors. It owns no aggregates or database tables.

## Use cases
- **Servers** (`application/servers`) — the provider contract. `ServerProvider.GetAllServers()` returns `Server`s carrying the catalog fields, the supported `AuthMethods` (`apiKey`, `oauth`), the OAuth `Scopes` that grant access and the built `mcpsdk.Server`. Server ids must be unique across providers.
- **NodeServers** (`application/nodeservers`) — provider for the nodeengine servers: API key or OAuth, scope `mcp:invoke`.
- **GetAllServers** (`usecases/servers/getallservers.go`) — catalog of every provided server; URLs are built from `ServerURLTemplate`.
- **GetProtectedResourceMetadata** (`usecases/servers/getprotectedresourcemetadata.go`) — RFC 9728 document for a server endpoint, advertising that server's scopes.
- **Adapter** (`application/adapter`) — builds an `mcpsdk.Server` per nodeengine server: registers each node as an MCP tool, augments the input schema with a required `credentials` object, wraps outputs under `items`, and on invocation extracts owner from request headers, resolves credentials and runs the executor.
- **History** (`application/history`) — records every finished tool call through `executions`' `ToolInvocationService` (source `mcp`) with the credential-free parameters, credential references and outputs; `executions` publishes the realtime event. Failures to record are logged, never surfaced to the caller.
- **CredentialResolver** (`application/credentialresolver`) — parses `credential:organization:<uuid>` references, enforces the owner is an organization, rejects a credential whose catalog key differs from the one the tool asks for, and regenerates OAuth tokens as needed.

## HTTP API
- `GET /servers` — every MCP server (nodeengine and platform); cached 5 min; no auth.
- `GET /.well-known/oauth-protected-resource/:serverId/mcp` — RFC 9728 protected resource metadata naming the `mcpoauth` authorization server and the server's scopes; no auth.
- `POST|GET|DELETE /:serverId/mcp` — streamable MCP endpoint per server (stateless). Every server accepts an MCP OAuth 2.1 bearer token whose audience is that server and that carries one of its scopes; servers that support `apiKey` also accept an organization API key (`X-API-Key` + `mcp:invoke` scope). A request without valid credentials gets a `401` carrying the `WWW-Authenticate: Bearer resource_metadata=…` challenge so MCP clients can discover the authorization server. Owner is taken from the auth-injected owner-type/owner-id headers either way.

## Dependencies
- **Bounded contexts:** mcpoauth (`MetadataService` for resource URLs and protected resource metadata), nodeengine (`mcp.ServerService`, `executors.ExecutorService`), credentials (`CredentialService` for the credential type check), credentialoauth (`CredentialOAuthTokenRegenerateService`), executions (`ToolInvocationService` for call history), apikeys (`ScopeMCPInvoke` for the required scope). Other contexts plug in through `ServerProviders` (`mcpplatform`).
- **Infrastructure:** `modelcontextprotocol/go-sdk` (streamable HTTP MCP servers); Fiber cache middleware; API-key and access-token auth middlewares. No database.

## Layout
Adapter-style module (no `domain/` layer): `application/` (servers, nodeservers, adapter, credentialresolver, history), `usecases/` (`services.go` collects providers and builds the per-server streamable handlers; `servers/` holds the catalog and metadata queries), `http/` (routes, per-server auth), wired in `module.go`.
