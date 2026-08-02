# mcp-api

> A Fiber HTTP server exposing the Model Context Protocol (MCP) surface so external MCP clients can invoke workflow nodes as tools, authenticated by API key.

## What it does
Loads the MCP API config and bootstraps a trimmed module graph via `bootstrap.NewMCPAPI`, then serves a Fiber app that registers only the MCP module routes. Requests are gated by an API-key middleware (built from the apikeys module's `APIKeyValidator`) and a cache middleware; a Redis-backed rate limiter and liveness/readiness probes are also wired. The MCP module exposes node executors as MCP tools, resolving OAuth tokens per invocation.

## Bootstrap & config
- **Assembler:** `bootstrap.NewCore` + `bootstrap.NewMCPAPI` — wires the contexts MCP needs (no workflows/executions/triggers/taskrunner).
- **Config:** `config.LoadMCPAPI()` → `config.MCPAPIConfig` (embeds `SharedConfig`, `HTTPServer`, `MCP`).
- **Runs as:** HTTP server on `HTTP_SERVER_*` address.

## Bounded contexts activated
- common, account, organizations, web3, credentialoauth
- nodeengine, platform, credentials
- apikeys
- mcp

## Notes
- Only `MCPModule.Register` mounts routes; auth is API-key only (no JWT auth middleware).
- MCP server URL is templated from `MCP_SERVER_URL_TEMPLATE`.
- Shutdown via `bootstrap.ListenAndWait` (Fiber + core).
