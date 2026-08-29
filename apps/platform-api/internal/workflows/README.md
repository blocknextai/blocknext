# Workflows

> Owns workflow definitions (the canvas-truth node/edge JSON graph) plus the AI-assisted generation chat used to author them.

## Responsibility
This context is the system of record for workflow definitions scoped to an organization: the canvas graph (nodes + edges), its metadata, and CRUD/duplication lifecycle. It also owns the workflow-generation feature — chat sessions and messages that drive an LLM to produce workflow JSON. It enriches read responses with cross-context info (owner identity). It does not execute workflows (that is taskrunner/nodeengine) nor store credentials.

## Domain
- **Aggregates / key types:**
  - `workflows.Workflow` — the workflow definition: `OrganizationID`, `OwnerID`, `Title`, `Description`, `IsPinned`, `IsActive`, and the graph `Nodes []dag.Node` / `Edges []dag.Edge` (stored as JSONB). The `dag.Node`/`dag.Edge` shape is the canvas truth: nodes carry full `parameters`/`settings`, edges pair with in-node `$references`; there is no credentials field on the workflow (credentials are resolved at run time).
  - `generation/sessions.GenerationSession` — an AI authoring session (`OrganizationID`, `UserID`, `Title`); soft-deletable.
  - `generation/messages.GenerationMessage` — a chat turn within a session (`SessionID`, `Role`, `Content`, `Metadata`).
- **Key rules / invariants:**
  - `Workflow` requires non-nil `OrganizationID` and `OwnerID` and a non-blank `Title` (sentinels `ErrOrganizationIDIsRequired`, `ErrOwnerIDIsRequired`, `ErrTitleIsRequired`; also `ErrWorkflowNotFound`, `ErrTitleTooLong`).
  - `GenerationSession` requires a non-blank title; `GenerationMessage` requires a non-nil `SessionID`.

## The canvas model

**The document is the canvas.** A workflow persists exactly two structural columns — `nodes JSONB` and `edges JSONB` — holding `[]dag.Node` / `[]dag.Edge` (`packages/go-packages/dag`). What the editor renders is byte-for-byte what the runtime executes; there is no separate "compiled" representation. Each `dag.Node` is self-contained: `id` (canvas id), `nodeId` (catalog id, e.g. `gemini.imagen`), optional `instruction` (design-time prose for function calling), full `parameters`, per-node `settings` (`maxRetries`, `retryDelay`, `timeout`, `continueOnError`, `disabled`), `credentials`, and `position`. `runtimeInstruction`/`runtimePrompt` exist on the struct but are injected at run time by the trigger's `RuntimeConfig` overlay — never authored into the stored document.

**Edges declare flow; `$references` carry the data.** An edge `source → target` alone moves nothing — the target node must also contain a `$reference` to the source's output in its `parameters` (or `instruction`). Both halves are required: an edge without a matching reference is an unresolvable dependency, and a reference without an edge is a broken flow. References are plain strings, valid in **any** string field, and resolved by regex substitution in taskrunner's `DataProcessor` against the upstream node's stored outputs:

| Pattern | Meaning |
| --- | --- |
| `$<nodeId>_<canvasId>.<field.path>` | Field from the referenced node's output for the current item (node key = catalog `nodeId` + `_` + canvas `id`, e.g. `$gemini.imagen_2.images`) |
| `$<key>[0].<field>` / `$<key>[*].<field>` | Explicit index / all items collected into an array |
| `$<key>.first().<field>` / `.last()` / `.get(n)` | Positional access over the node's output list |
| `$input.<field.path>` | The item feeding this node at the same position, without naming the node it comes from; only defined when the node has exactly one incoming edge |
| `$trigger.source` / `.sender` / `.prompt` / `.payload` | Fields of the run's `TriggerContext` (webhook adapter output or runtime prompt) |

**Credentials are referenced, never stored.** The workflow JSON carries no secrets: a node's `credentials` map holds opaque reference strings of the form `credential:<user\|organization>:<uuid>` (`internal/common/domain/credential`). At execution time taskrunner's `CredentialProcessor` parses the reference, resolves the owner scope, and fetches (refreshing OAuth tokens if needed) the actual credential material — which therefore lives only in the credentials store, and a duplicated or exported workflow leaks nothing.

**Deliberate simplicity.** Control flow is intentionally minimal: the only system nodes are `starter`, `condition`, `sleep` and the canvas-only `annotation`. There is **no loop, no sub-workflow, and no generic HTTP node type**, and the graph must be acyclic — `dag.New` runs a topological sort and rejects cycles (`ErrCycleDetected`).

**Branching is per item.** A node that routes implements `executors.BranchingExecutor`, returning the indexes of the items that leave through each output handle; `condition` returns them under `true` and `false`. The runner stores each branch separately (`<nodeKey>#<handle>`) together with the original index of every item it kept, so a consumer reads the branch its own edge leaves from, and a `$reference` to a node *upstream* of the branch still resolves to the matching item rather than to position `i` of the filtered list. A handle with no items reaches nothing: its children are marked `skipped`, and the skip walks on to any descendant whose every parent is skipped — a join with one live parent still runs. An edge with no `sourceHandle` is not routed, so flows authored before handles existed keep their old behaviour.

## AI workflow generation

Workflows can be authored by chat: the generation feature (env-gated via `WORKFLOWS_GENERATION_ENABLED`, provider `gemini` or `local` selected by `WORKFLOWS_GENERATION_PROVIDER`) turns a conversation into canvas JSON.

- **Session replay** — `chat.ChatService.SendMessage` persists the user message, then replays the session's full message history (both roles) to the LLM, so the model iterates on its own prior workflow drafts across turns.
- **Live catalog injection** — the system instruction is a platform-owned markdown file (embedded from `internal/prompts`) with three placeholders — `{AVAILABLE_NODES_JSON}`, `{AVAILABLE_CREDENTIALS_JSON}`, `{AVAILABLE_TRIGGER_VARIABLES_JSON}` — filled at request time by context builders from the live nodeengine registries. The model can only compose nodes that actually exist in this deployment.
- **Strict output contract** — the prompt pins the exact canvas shape: a single ` ```json:workflow ` block containing `{nodes, edges}` arrays, mandatory `system_starter` node first, edge/`$reference` parity checked per edge, and an explicit ban on emitting a `credentials` field (the runtime resolves credentials at execution). It also enforces a cost rule: prefer static `parameters` (zero-cost substitution) over `instruction` prose, since every instruction triggers a function-calling LLM call on each run.
- **SSE streaming** — the presentation handler streams the reply as `text/event-stream` (`event: message` chunks, `event: done` / `event: error` terminators); the completed assistant reply is saved back to the session.
- **Scope guard** — the system prompt restricts the assistant to workflow building only and refuses off-topic requests with a fixed sentence.

## Use cases (application)
Workflows (`application/workflows/*`):
- `createworkflow` — create a workflow.
- `updateworkflow` / `deleteworkflow` — edit graph + metadata / soft-delete (transactional).
- `duplicateworkflow` — clone an existing workflow.
- `getallworkflows` / `getworkflowbyid` — read with enrichment: owner alias/verification and linked accounts.
- `getworkflowforrun` — return the workflow graph plus resolved node and credential schemas (from nodeengine) needed to execute it.

Generation (`application/generation/*`):
- `sessions/{createsession,updatesession,deletesession,getallsessions,getallsessionmessages}` — CRUD over generation sessions and their messages.
- `chat.ChatService.SendMessage` — persists the user message, replays session history, builds the system instruction (injecting available nodes / credentials / trigger-variable schemas), streams the LLM response back, and saves the assistant reply on completion.
- Context builders `nodeschema`, `credentialschema`, `triggervariables` — render the catalog placeholders embedded in the generation system prompt.

## HTTP API
All routes under `/organizations/:organizationId/workflows`, each `Authenticate()` + `RequireOrganizationPermission(...)`:
- `GET /` — list workflows (ReadWorkflow).
- `POST /` — create (CreateWorkflow).
- `GET /:workflowId` — get by id (ReadWorkflow).
- `PATCH /:workflowId` — update (UpdateWorkflow).
- `DELETE /:workflowId` — delete (DeleteWorkflow).
- `POST /:workflowId/duplicate` — duplicate (CreateWorkflow).
- `GET /:workflowId/run` — get run-ready graph + schemas (ReadWorkflow).

Generation, under `.../generation/sessions`:
- `GET /` (ReadWorkflowGenerationSession), `POST /` (CreateWorkflowGenerationSession), `PATCH /:sessionId` (Update...), `DELETE /:sessionId` (Delete...).
- `GET /:sessionId/messages` (ReadWorkflowGenerationMessage), `POST /:sessionId/messages` (CreateWorkflowGenerationMessage) — streams the AI reply via `chatService`.

## Events
N/A — this module does not publish or consume durable domain events.

## Dependencies
- **Bounded contexts (via service interfaces):** `organizations` (org users), `account` (users + linked accounts for owner info), `nodeengine` (node, credential, adapter schemas — for run + generation context), `llm/streamingchat` (generation provider).
- **Infrastructure:** Postgres schema `workflows` — tables `workflows` (graph as JSONB `nodes`/`edges`), `generation_sessions`, `generation_messages`. `database.TransactionManager` for write commands. `config.WorkflowsOptions` supplies the generation system instruction.

## Layout
Full DDD: `domain/` (workflows + generation/{sessions,messages}), `application/` (CQRS handlers + generation chat/context builders), `infrastructure/` (repositories, migrations), `presentation/` (workflows + generation routes), wired in `module.go`. `module.go` exposes `WorkflowService` for cross-context consumers.
