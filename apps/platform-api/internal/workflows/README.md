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
