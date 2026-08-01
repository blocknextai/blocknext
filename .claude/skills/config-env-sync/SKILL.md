---
name: config-env-sync
description: Use this skill when the user has edited (or is about to edit) a configuration struct under `internal/config/` and needs `<repo-root>/.env.example` synchronised in the same change. Triggers on phrases like "I added a new env var", "sync .env.example", "I changed a config option", "audit config drift", "what's missing from .env.example". Walks the Go config struct tree, computes every full env-var name from nested `envPrefix` tags, diffs that set against `<repo-root>/.env.example`, and reports/applies the missing additions, orphan removals, and naming mismatches.
---

> **Working directory:** all relative paths in this skill (`internal/...`, `cmd/...`, `go build ./...`) are relative to `apps/platform-api/` in the monorepo. Make targets run from the repo root.


# Config Env Sync

Keeps `internal/config/**.go` and `<repo-root>/.env.example` in lockstep. Project memory mandates that any `env` / `envPrefix` tag change must update `<repo-root>/.env.example` in the same turn (memory: `feedback_env_example_sync`); this skill audits and applies that sync.

## When to use

- User edited an `*Options` struct under `internal/config/` (added a field, renamed a tag, removed a field)
- User says "sync env example" / "what's missing from .env.example" / "audit config drift"
- After running `command-create` / `module-create` if the new code reads a configuration value
- Before committing changes that touch `internal/config/`

## What this skill produces

A diff between the **struct-derived env var set** and the **set declared in `<repo-root>/.env.example`**, followed by:

- **Missing** entries — present in struct, absent from `<repo-root>/.env.example`. Skill adds them with sensible placeholder values.
- **Orphan** entries — present in `<repo-root>/.env.example`, no longer in any struct. Skill removes them after confirming with the user.
- **Naming mismatches** — typos / drift between the struct's tag and the `<repo-root>/.env.example` line. Skill renames in place.

## Required information

Confirm with the user (ask only what's missing):

1. **Mode** — `report` (default; show diff only) or `apply` (modify `<repo-root>/.env.example`).
2. **Scope** — by default, audit every file under `internal/config/`. If the user names a single area (e.g., "just `workflows.go`"), narrow to that struct tree.

## How env vars are derived

The config tree has **three roots** (one per binary), each defined in its own file under `internal/config/`:

- `SharedConfig` (`internal/config/shared.go`) — the bulk of the env vars. Its fields carry the top-level prefixes directly (`DATABASE_`, `WORKFLOWS_`, `OAUTH_`, ...). There is NO outer prefix above them.
- `PlatformAPIConfig` (`internal/config/platform_api.go`) — `Shared *SharedConfig` (tagged `env:"-"`, loaded separately via `LoadShared()`) plus app-specific options: `HTTPServer` (`PLATFORM_API_`), `WebSocket` (`WEBSOCKET_`), `TaskRunner` (`TASK_RUNNER_`).
- `MCPAPIConfig` (`internal/config/mcp_api.go`) — `Shared` plus `HTTPServer` (`MCP_API_`), `MCP` (`MCP_`).
- `WebhookAPIConfig` (`internal/config/webhook_api.go`) — `Shared` plus `HTTPServer` (`WEBHOOK_API_`), `TaskRunner` (`TASK_RUNNER_`).

The full env var set is the **union** across all three roots. `Shared` is parsed once and reused, so its vars are covered by walking `SharedConfig` a single time; then add each app config's own (non-`Shared`) sub-trees. Every field has either:

- `env:"<NAME>"` — leaf string env var (NAME may be `"-"` for programmatic-only fields; `Shared *SharedConfig` itself is `env:"-"` and is walked as its own root, not concatenated)
- `envPrefix:"<PREFIX>_"` — sub-struct prefix (recurses into the nested type)

The full env var name is the **concatenation** of all prefixes from root to leaf, plus the leaf `env` tag.

Example trace for `WORKFLOWS_GENERATION_GEMINI_TEMPERATURE` (rooted in `SharedConfig`):

```
SharedConfig
  └── Workflows  envPrefix:"WORKFLOWS_"  → WorkflowsOptions
        └── Generation  envPrefix:"GENERATION_"  → WorkflowsGenerationOptions
              └── Gemini  envPrefix:"GEMINI_"  → WorkflowsGenerationGeminiOptions
                    └── Temperature  env:"TEMPERATURE"  →  WORKFLOWS_GENERATION_GEMINI_TEMPERATURE
```

**Rules for derivation:**
- Concatenate prefixes left-to-right; each prefix already ends with an underscore (matches existing files).
- Leaf field with `env:"-"` is **excluded** from the env var set (it's set programmatically in `Load()`, e.g., `SystemInstruction`).
- Leaf field with NO tag at all is **excluded** (the `caarlos0/env` library skips it).
- `envPrefix:"-"` doesn't exist in this codebase; if found, treat as a malformed tag and surface it.

## Workflow

### Step 1: Walk the config tree

Walk all three roots: `SharedConfig` (`shared.go`) once, then the non-`Shared` fields of `PlatformAPIConfig` (`platform_api.go`), `MCPAPIConfig` (`mcp_api.go`), and `WebhookAPIConfig` (`webhook_api.go`). For each field:
- Read its tags.
- If `env:"<NAME>"` and NAME != `-`: emit `<accumulated_prefix><NAME>` as a set member.
- If `env:"-"` on a `*SharedConfig` field: skip it here (it's already walked as its own root — do NOT concatenate the app prefix onto shared vars).
- If `envPrefix:"<P>"`: open the field's named type (defined in another `internal/config/*.go` file), recurse with `accumulated_prefix + P`.

The output is a sorted slice of full env var names with their Go types (so we can pick reasonable placeholders).

**Implementation hint:** rather than parsing AST, `grep` the config files. Each struct field declaration is on its own line:
```bash
grep -rEn '`(env|envPrefix):' internal/config/
```

The grep gives `<file>:<line>:<field-line-with-tag>`. Group by struct (file boundaries roughly), then resolve the named-type references manually by looking up the right file. For deep trees this is a small recursive read, not a parse.

### Step 2: Parse `<repo-root>/.env.example`

Treat each line as `KEY=VALUE` if it matches `^[A-Z][A-Z0-9_]*=`. Ignore:
- Empty lines
- Comment-only lines (`# ...`)
- Lines that do not match the `KEY=` pattern (these are likely shell-style continuations or `${VAR}` substitutions)

Build a map `key → { line_number, value }`.

### Step 3: Diff

- **Missing** = `struct_keys − env_example_keys`
- **Orphan** = `env_example_keys − struct_keys`.

### Step 4: Generate sensible placeholders for missing keys

When adding to `<repo-root>/.env.example`, pick a placeholder by Go type:

| Go type                     | Default placeholder                          |
| --------------------------- | -------------------------------------------- |
| `string`                    | empty (`<KEY>=`) unless the name suggests a value (see below) |
| `bool`                      | `false` (or `true` if the field name implies enabled-by-default, e.g., `Enabled`) |
| `int` / `int32` / `int64`   | `0`                                          |
| `float32` / `float64`       | `0.0`                                        |
| `time.Duration`             | `30s` (a non-zero default; `0s` breaks many things) |
| Custom type with `const`    | First declared const value (e.g., `WorkflowsGenerationProvider` → `gemini`) |

**Name-based heuristics for strings:**
- `*Url`, `*URL` → `http://localhost:<plausible-port>` (or empty if no port hint)
- `*Host` → `0.0.0.0` for server hosts, `localhost` for client hosts (heuristic: ports 1-1023 + the field is in an `Api*Options` → server; otherwise client)
- `*Port` → empty (the user picks)
- `*ApiKey`, `*Secret`, `*Token`, `*Password` → `REPLACE_WITH_<purpose>` plus a `# SECURITY:` comment line above with the openssl recipe (see "Security comments" below)
- `*Model` → `gemini-flash-latest` for Gemini-related fields, empty otherwise
- `*ConnectionStrings*Default` → `'host=db port=5432 user=postgres password=REPLACE_WITH_SECURE_PASSWORD dbname=postgres sslmode=prefer'`

When in doubt, leave the value empty and let the user fill it in. **Do not invent secrets, URLs, or production-looking defaults.**

### Step 5: Decide insertion location

`<repo-root>/.env.example` groups entries by section, with one blank line between sections. The section is identified by the leading prefix (e.g., `API_*`, `DB_*`, `WORKFLOWS_*`). For each missing key:

1. Strip everything after the first `_` to get the section prefix (`WORKFLOWS_GENERATION_X` → section `WORKFLOWS`).
2. Find the existing block of lines starting with that prefix in `<repo-root>/.env.example`. Insert the new line at the end of that block, before the trailing blank line.
3. If no block exists yet (entirely new section), append a new block at the end of `<repo-root>/.env.example`: blank line + new entries + (no trailing blank line).
4. Keep entries within a section in **the order they appear in the struct tree**, not alphabetical. Existing entries reflect this — match it.

### Step 6: Security comments

For secret-bearing keys (`*_SECRET`, `*_API_KEY`, `*_TOKEN`, `*_PASSWORD`), prepend a `# SECURITY:` comment line that mirrors the existing pattern in `<repo-root>/.env.example`:

```bash
# SECURITY: Generate with: openssl rand -hex 32 (must be 32+ characters)
JWT_SECRET=REPLACE_WITH_32_CHAR_SECRET_MIN
```

Inspect the existing `<repo-root>/.env.example` for the recipe used by similar secrets and reuse the wording. Don't invent new wording.

### Step 7: Report

Format:

```markdown
## Config Env Sync — `<scope>`

### Missing in .env.example (N)
- `WORKFLOWS_GENERATION_GEMINI_TIMEOUT` — `time.Duration` — placeholder: `30s`
- `WEBHOOK_TRIGGER_RETRY_COUNT` — `int` — placeholder: `0`

### Orphan in .env.example (M)
- `OAUTH_LEGACY_REDIRECT_URL` (line 87) — no matching struct field; safe to remove?

### Naming mismatches (K)
- struct says `WORKFLOWS_GENERATION_GEMINI_MAX_OUTPUT_TOKENS` but .env.example has `WORKFLOWS_GENERATION_GEMINI_MAX_TOKENS` (line 73)

### Summary
- Missing: N
- Orphan: M (requires user confirmation before removal)
- Mismatched: K
```

In `apply` mode, perform the additions/renames automatically; **always ask the user before removing orphan entries** — they may be valid env vars consumed elsewhere (e.g., by tooling scripts rather than the Go app).

### Step 8: Verification

After applying changes:

1. Re-run the diff. The output should be empty.
2. Show the user the modified `<repo-root>/.env.example` lines (`git diff .env.example`).
3. Remind them to update their local `.env` if the change affects their dev setup — the skill does NOT touch `.env` (it's gitignored and may contain real secrets).

## Critical rules (project memory)

1. **Sync in the same turn as the struct edit.** Memory `feedback_env_example_sync` is explicit: any `env` / `envPrefix` tag change must update `<repo-root>/.env.example` immediately. If the user is in the middle of a struct edit, run this skill before they switch context.
2. **Never auto-remove orphan entries without confirmation.** They might be tooling vars or recently-deleted struct fields the user intends to bring back.
3. **Never write real secrets or production URLs into `<repo-root>/.env.example`.** Always use `REPLACE_WITH_<x>` style placeholders + `# SECURITY:` comments where applicable.
4. **Match existing section ordering.** New entries go inside the right prefix-block, not alphabetically appended.
5. **`env:"-"` fields are not config drift.** They're populated programmatically (e.g., `SystemInstruction` is loaded from a file in `Load()`); excluding them from `<repo-root>/.env.example` is correct.

## Final steps

After applying the sync:

1. `git diff .env.example` — show the user the change.
2. If the user used the `module-create` skill (config-only changes from `module.go` adding a new dep don't affect this skill, but new `internal/config/<area>.go` files do): the audit is critical here, since a brand-new struct produces brand-new env vars.
3. Suggest the user re-source their `.env` if their local copy is out of date: `source .env`.
4. If `make migration-up` or any other script depends on the new var, tell the user explicitly — they may need to set it before re-running.

## What NOT to do

- Do NOT touch `.env` — it's the developer's local file with real secrets, gitignored, never modified by tooling.
- Do NOT reorder existing `<repo-root>/.env.example` entries to alphabetise. The current order matches the struct's declaration order, which has signal value (related vars stay together).
- Do NOT remove the `# SECURITY:` comments above secret keys when editing nearby lines.
- Do NOT add a placeholder that looks like a real value (e.g., `JWT_SECRET=mysecretkey`). Use `REPLACE_WITH_<x>` style strictly.
- Do NOT silently drop env vars whose Go type you can't infer. Flag them; let the user decide.
- Do NOT batch this with other refactors. The sync is its own commit so reviewers can spot config drift in isolation.
- Do NOT add new envPrefix tags as a way to "make the env var look nicer." Tag changes propagate to deployed environments — keep them stable unless the rename is part of an explicit refactor.
