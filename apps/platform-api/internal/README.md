# `internal/` — Bounded Context Map

This service follows Domain-Driven Design: each directory under `internal/` is either a
**bounded context** (a `module.go` wires its `domain` / `application` / `infrastructure` /
`presentation` layers) or a **shared/infrastructure package**. Contexts never reach into each
other's repositories — they integrate only through exported **service interfaces** (synchronous)
or **domain events** over the `eventbus` (asynchronous). See each directory's own `README.md` for
its detailed scope.

## Subdomains

### Identity & Access
| Context | Scope |
| --- | --- |
| [`account`](account/README.md) | User identity & authentication (password, magic-link, OAuth sign-in, crypto-wallet) + per-user settings. |
| [`organizations`](organizations/README.md) | Tenants/workspaces, RBAC members. Auto-provisions a default org on `user.created`. |
| [`apikeys`](apikeys/README.md) | Organization-scoped programmatic-access keys (issuance, scopes, rotation, validation). |
| [`credentialoauth`](credentialoauth/README.md) | User-facing OAuth2 authorization-code (PKCE) flow + token refresh for third-party credentials. |

### Workflow Core
| Context | Scope |
| --- | --- |
| [`workflows`](workflows/README.md) | Workflow definitions (canvas-truth node/edge JSON graph) + AI-assisted generation chat. |
| [`nodeengine`](nodeengine/README.md) | Catalog & runtime registries of nodes (actions), executors, credentials, function-calling, trigger adapters. |
| [`credentials`](credentials/README.md) | User/org-scoped third-party integration secrets, encrypted at rest, consumed by nodes. |
| [`platform`](platform/README.md) | Host-provided "platform" OAuth/API credentials from operator config (no DB). |
| [`triggers`](triggers/README.md) | Persisted schedule (cron) + webhook triggers; resolves inbound webhooks into execution context. |
| [`taskrunner`](taskrunner/README.md) | Executes workflow tasks as DAGs (trigger/dispatch/run/rerun/cancel). |
| [`executions`](executions/README.md) | Records of task runs, per-node executions, and worker task-claim leases. |
| [`llm`](llm/README.md) | Provider-agnostic LLM client (streaming chat + function calling); Gemini + local backends. |

### Channels & Integrations
| Context | Scope |
| --- | --- |
| [`mcp`](mcp/README.md) | Model Context Protocol server exposing `nodeengine` nodes as MCP tools (HTTP adapter, no domain). |
| [`web3`](web3/README.md) | Ethereum signature verification backing the crypto-wallet (MetaMask) login. |
| [`notifications`](notifications/README.md) | Event-driven in-app inbox; fans out domain events from other contexts. |
| [`webhooks`](webhooks/README.md) | Inbound HTTP edge — receives external webhook calls (triggers) and routes to the owning processor. |
| [`ws`](ws/README.md) | Per-org WebSocket fan-out gateway pushing realtime messages to clients. |

### Shared kernel & infrastructure (not business contexts)
| Package | Scope |
| --- | --- |
| [`common`](common/README.md) | Shared kernel: auth/authz middleware, CQRS pipeline, request-context helpers, shared value types. |
| [`eventbus`](eventbus/README.md) | Transactional-outbox event infrastructure (outbox relay + inbox dedupe + in-process typed `Bus`). |
| [`realtime`](realtime/README.md) | Ephemeral server→client Redis broadcaster (pairs with `ws`). |
| [`cache`](cache/README.md) | Config-driven Redis-backed cache factory. |
| [`filegateway`](filegateway/README.md) | HTTP-client gateway to an external file-storage service. |
| [`config`](config/README.md) | Env-var options layer (composition-root config; never decoupled). |
| [`bootstrap`](bootstrap/README.md) | Composition root: shared infra + per-entrypoint assemblers. |

## Integration relationships

Synchronous calls go through exported service interfaces; the dashed lines are asynchronous
domain events relayed by `eventbus`.

```mermaid
flowchart TB
    subgraph identity[Identity & Access]
        account
        organizations
        apikeys
        credentialoauth
    end
    subgraph wf[Workflow Core]
        workflows
        nodeengine
        credentials
        platform
        triggers
        taskrunner
        executions
        llm
    end
    subgraph channels[Channels]
        mcp
        web3
        notifications
        ws
        webhooks
    end

    account -. user.created .-> organizations
    account -. domain events .-> notifications
    organizations -. domain events .-> notifications

    account --> web3
    workflows --> nodeengine
    taskrunner --> nodeengine
    taskrunner --> executions
    taskrunner --> credentialoauth
    taskrunner --> llm
    nodeengine --> platform
    credentials --> nodeengine
    credentials --> platform
    credentialoauth --> credentials
    mcp --> nodeengine
    webhooks --> taskrunner
    triggers --> taskrunner
    taskrunner -. progress .-> ws
```

## Runtime entrypoints

`bootstrap` assembles the contexts into several processes (see `cmd/` + `internal/bootstrap`):

- **PlatformAPI** — the main HTTP API (most contexts' `presentation` routes).
- **MCPAPI** — the Model Context Protocol server (`mcp` + `nodeengine`).
- **WebhookAPI** — the inbound webhook edge (`webhooks` → triggers).
- **TaskWorker** — executes queued workflow tasks (`taskrunner` + `executions`).
- **Cron workers** — scheduled jobs (eventbus outbox relay, trigger scheduling, etc.).

## Conventions

- **No cross-context repository injection** — contexts integrate via service interfaces only.
- **CQRS** — each use case is a handler under `application/<usecase>/`; reads use the mapper convention.
- **Events** — names follow `<context>.<entity>.<action>` (past tense); durable server→server flows
  go through `eventbus`, ephemeral server→client flows through `realtime`/`ws`.
- **New modules** must be registered in the bootstrap wiring and the migration module registry.
