# Node Engine

> The catalog and runtime registry of workflow nodes (actions), their executors, third-party credentials, function-calling declarations, and MCP server exposure.

## Responsibility

The node engine owns the **definitions** of everything a workflow can do: each node (an action like "send a Slack message" or a system primitive like `condition`), its executor (the code that runs the action), the credentials that authenticate against third-party providers, and the trigger adapters that normalize inbound webhook payloads. At process startup it populates a set of in-memory registries (nodes, executors, credentials, function-calling, MCP servers, adapters) and exposes read services + catalog HTTP endpoints over them. It does **not** own workflow definitions, runs, or scheduling — the `workflows`, `executions`, and `taskrunner` contexts consume this engine's registries (via the exposed services) to resolve and execute nodes. The engine is stateless apart from the singleton registries; there is no persistence layer (`infrastructure/` is wiring + CQRS handler registration only).

## Core concepts

- **Node** (`domain/nodes/node.go`): a declarative action descriptor — `ID`, `Version`, name/description, `InputSchema`/`OutputSchema` (`jsonschema-go`), `Categories`/`SubCategories`/`Tags`, `SupportedCredentials`, `Annotations`, `Credits`, `HasNaturalLanguage`, `Disabled`. Implements `NodeManager` (getter interface).
- **Executor** (`domain/executors/executor.go`): the behavior for a node, keyed by the same ID. `ExecutorManager.ExecuteWithContext(ctx, credentials, data)` is the execution entry point.
- **Credential** (`domain/credentials/credential.go`): a provider auth descriptor — schema, `IsOAuth1`/`IsOAuth2`/`IsSupportPlatform`, `SupportedNodes`. OAuth credentials may implement `RefreshableCredential`. `IsSupportPlatform` is auto-derived at load (host-provided OAuth apps).
- **Adapter** (`domain/adapters/adapter.go`): a trigger adapter that maps a raw webhook payload into a `TriggerContext` (`source`/`sender`/`prompt`/`payload`, exposed as `$trigger.*` variables).
- **FunctionCalling** (`domain/functioncalling/`): an LLM function-declaration generated from a node's input schema via `functioncalling.Generate(node)`; disabled unless the node has `HasNaturalLanguage`.
- **MCP Server** (`domain/mcp/server.go`): groups a provider's nodes as MCP tools (`Tools []NodeManager`) for external MCP clients.
- **Registries**: six singleton registries (`GetRegistry()` / `sync.Once`) — nodes, executors, credentials, functioncalling, mcp, adapters — each with `Register*` / `Get*` / `GetAll*`. Disabled items are skipped at registration.

## Node & credential catalog

Nodes are organized as `nodes/<provider>/<action>/{node.go,executor.go}` with a per-provider `register.go` and shared `helpers/`. There are ~26 providers (`nodes/nodes.go`) spanning LLMs/AI (anthropic, chatgpt, gemini, deepseek, deepl), media generation (elevenlabs, veo, piapi, sunomusic, soundcloud), social/messaging (discord, slack, telegram, whatsapp, x, instagram, facebook, tiktok, linkedin), Google (docs, drive, gmail, sheets, keep, youtube), productivity/data (airtable, notion, coingecko, sendgrid), plus a `system` provider for engine primitives (`condition`, `sleep`, `starter`). Roughly ~80 node actions in total.

A provider's `register.go` wires each action consistently: build the node, build a JSON-schema validator (`jsonschema.New[Input]`), build the executor, then `nodes.RegisterNode` + `executors.RegisterExecutor`, optionally `functioncalling.RegisterFunctionCalling(functioncalling.Generate(node))`, and a single `mcp.RegisterServer` grouping the provider's nodes as MCP tools. A node links to a credential by ID through `Node.SupportedCredentials`, and the credential mirrors the link via `Credential.SupportedNodes`. ~31 credentials live in `credentials/` as `<provider>_(api|oauth2).go`, registered in `credentials/credentials.go` (OAuth2 ones receive the OAuth redirect URL at registration).

## Use cases (application)

Read-side only (queries), each under `application/` as a CQRS handler wrapped with `cqrs.ValidationBehavior`:

- **getallnodes** — list all registered nodes.
- **getallcredentials** / **getcredentialbyid** — list / fetch credential descriptors.
- **getallwebhooksources** — list trigger adapters as webhook sources (uses the webhook URL template).
- **getalltriggervariables** — list available `$trigger.*` variables.

Application also exposes plain services consumed by other contexts: `NodeService`, `CredentialService`, `CredentialProcessor` (masks/strips write-only & OAuth fields), `ExecutorService`, `AdapterService`, `ServerService` (MCP), plus a `jsonschema` validator helper.

## HTTP API

All routes are under `/node-engine` and wrapped with a 5-minute cache middleware (no auth middleware applied here):

- `GET /node-engine/nodes` — all registered nodes.
- `GET /node-engine/credentials` — all credential descriptors.
- `GET /node-engine/credentials/:id` — single credential by ID.
- `GET /node-engine/webhook-sources` — trigger adapters as webhook sources.
- `GET /node-engine/trigger-variables` — available `$trigger.*` variables.

## Events

- **Published / Consumed:** none.

## Dependencies

- **Bounded contexts:** `filegateway` (`FileGateway`, injected into media-capable node registrations); `common` (`cqrs`, `domain/oauth2`, `domain/credential` mask constant). Other contexts (workflows, executions, taskrunner, mcp presentation) consume this engine's exposed services/registries.
- **Infrastructure:** `jsonschema-go` (input/output schemas + validation); `gofiber/v3`; go-packages fiber cache middleware and `json`.

## Extending

Add a node via `nodes/<provider>/<action>/` (a `node.go` + `executor.go`, registered in the provider `register.go`: node + executor + optional function-calling + MCP tool), keeping `SupportedCredentials` in sync with the credential's `SupportedNodes`; the `node-create` / `node-update` conventions automate this. Add a credential as `credentials/<provider>_(api|oauth2).go` and register it in `credentials/credentials.go` (the `credential-create` / `credential-update` conventions cover the registry wiring and cross-references).

## Layout

```
domain/          core types + singleton registries (nodes, executors, credentials, functioncalling, mcp, adapters)
nodes/           <provider>/<action>/{node.go,executor.go}, per-provider register.go + helpers/; nodes.go aggregates Register()
credentials/     <provider>_(api|oauth2).go credential defs; credentials.go registers them
application/      read services + CQRS query handlers (nodes, credentials, adapters, executors, mcp) + jsonschema validator
infrastructure/  infrastructure.go — runs node/credential registration and builds the query Handlers struct
presentation/    fiber route registration (nodes, credentials, webhooks)
module.go        Module wiring: builds services, registers infrastructure, mounts presentation
```
