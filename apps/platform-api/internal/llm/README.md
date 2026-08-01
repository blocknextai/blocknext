# LLM

> Provider-agnostic LLM client abstraction for streaming chat and function calling. Infrastructure module with a `module.go`, but no DDD layers — it is a provider/client wrapper, not a business bounded context.

## Responsibility
Defines provider-neutral interfaces for the two LLM use cases the platform needs and selects a concrete backend at startup based on config. Consumers depend only on the abstract interfaces (`streamingchat.Provider`, `functioncalling.FunctionCallingService`) and never touch a vendor SDK directly. Either capability can be disabled via config, in which case `Module` exposes a `nil` provider.

## What it provides
- **`streamingchat.Provider`** — `StreamChat(ctx, systemInstruction, messages) (<-chan Chunk, error)`; streams assistant output as `Chunk{Type, Content, Error}` over `[]Message{Role, Content}`. Used for interactive workflow generation.
- **`functioncalling.FunctionCallingService`** — `ExecuteWithContext(ctx, data, functionDeclarations) ([]map[string]any, error)`; runs an LLM with declared functions over input data. Used for node/runtime function-calling.
- **`cache.Cache`** — `Ensure(ctx, instruction) string`; provider-side prompt/instruction caching contract (implemented for Gemini).
- **`Module`** (via `NewModule`) — exposes `StreamingChatProvider` and `FunctionCallingService`, each constructed from config (`WorkflowsGenerationOptions`, `FunctionCallingOptions`); returns `ErrInvalidStreamingChatProvider` / `ErrInvalidFunctionCallingProvider` for unknown providers.

## Supported providers
- **Gemini** (`providers/gemini`) — `streamingchat`, `functioncalling`, and a `cache` implementation; HTTP client built with go-packages `httpclient`.
- **Local LLM** (`providers/localllm`) — OpenAI-compatible base-URL backend with `streamingchat` and `functioncalling` implementations.

Selection is by config enum (`GeminiWorkflowsGenerationProvider` / `LocalLLMWorkflowsGenerationProvider`, and the function-calling equivalents).

## Consumed by
- **`workflows`** — `generation/chat` service and the `generation/send_message` presentation handler use `StreamingChatProvider` for AI workflow generation.
- **`taskrunner`** — `node_executor` uses `FunctionCallingService` for runtime function calling.
Both are wired through their respective `module.go` from the LLM `Module`.

## Dependencies
- **Infrastructure / external:** go-packages (`apperror`, `httpclient`, `json`), Google Gemini and any OpenAI-compatible local endpoint; configured via `config.WorkflowsGenerationOptions` and `config.FunctionCallingOptions`.

## Layout
- `streamingchat/` — `Provider` interface, `Message`, `Chunk`/`ChunkType`.
- `functioncalling/` — `FunctionCallingService` interface, input type, errors.
- `cache/` — `Cache` interface for instruction caching.
- `providers/gemini/` — Gemini client + `streamingchat`, `functioncalling`, `cache` implementations.
- `providers/localllm/` — local/OpenAI-compatible `streamingchat` and `functioncalling` implementations.
- `module.go` — config-driven provider selection and wiring.
