---
name: go-standards-review
description: Use this skill when the user wants to audit a Go file, package, or branch diff for project-specific Go standard violations in platform-api. Triggers on phrases like "review this file for go standards", "check the codebase for fmt usage", "audit my changes against project conventions", "scan for net/http usage", "fix go-standards violations". Greps the target for violations of every memory-tracked Go rule (no fmt, slog over log, strings.Builder, TrimSpace, httpclient, utils for JSON, append over index assign, new() over &, gofmt + lint, repository pagination/search param order, presentation pagination/search constructor calls, initialism casing in identifiers vs struct tag content), reports them as a punch list, and optionally applies fixes.
---

> **Working directory:** all relative paths in this skill (`internal/...`, `cmd/...`, `go build ./...`) are relative to `apps/platform-api/` in the monorepo. Make targets run from the repo root.


# Go Standards Review

Audits a target — a single file, a package directory, or the diff of the current branch — for violations of every Go convention captured in this project's memory. Reports violations as `file:line — rule — suggested fix`. Optionally applies the fix when it is mechanical.

## When to use

- User says "review this file for go standards" / "check for fmt usage" / "audit the diff" / "scan `<package>` for violations"
- After a large refactor, before committing, when sanity-checking a new module or aggregate
- When migrating older code that pre-dates a project convention

## Required information

Confirm with the user (ask only what's missing):

1. **Target** — one of:
   - **File**: an absolute or repo-relative `.go` path
   - **Package / directory**: anything under `internal/...`
   - **Diff**: the changes on the current branch vs. `main` (`git diff main...HEAD --name-only -- '*.go'`)
   - **Whole repo**: `internal/...` (slow; use only when explicitly asked)
2. **Mode** — `report` (default) or `fix`. In `fix` mode, the skill applies mechanical rewrites for unambiguous cases and leaves the rest in the report. Always show the report first, even in fix mode, so the user can confirm.

## Rules and detection

Each rule below has: a one-line summary, the memory entry it ties to, a grep pattern to find candidates, and a suggested fix. Run the patterns from the repository root; the `pcre2` engine (`-P`) is fine where shown.

### 1. No `fmt` package

**Memory:** `feedback_no_fmt_package` — `fmt` severely degrades performance; use `strings.Builder`, `strconv`, `errors.New`, `slog`.

**Detect:**
```bash
grep -rn --include='*.go' -E '"fmt"|fmt\.' <target>
```

**Allowed exceptions** — only flag if not one of these:
- `fmt` in a `*_test.go` file (test ergonomics; still prefer alternatives but not blocking)
- `fmt.Errorf` is **never** acceptable here — convert to `apperror.<Kind>(...)` or `errors.New(...)` per case

**Suggested fix:**
- `fmt.Sprintf("foo %s bar", x)` → `strings.Builder` (hot path) or `"foo " + x + " bar"` (cold init code)
- `fmt.Sprintf("%d", n)` → `strconv.Itoa(n)` (or `strconv.FormatInt(n, 10)`)
- `fmt.Sprintf("%f", f)` → `strconv.FormatFloat(f, 'f', -1, 64)`
- `fmt.Println(...)` → `slog.Info(...)` with structured kv pairs
- `fmt.Errorf("...: %w", err)` → `apperror.Internal("...").WithCause(err)` (preferred) or `errors.New("...")` chained

### 2. Use `slog`, never `log`

**Memory:** `feedback_use_slog_not_log` — structured logging is project standard.

**Detect:**
```bash
grep -rn --include='*.go' -E '^\s*"log"$|^\s*log\.' <target>
```

**Suggested fix:**
- `log.Printf("foo %s", x)` → `slog.Info("foo", "value", x)` (kv pairs, structured)
- `log.Fatal(err)` → `slog.Error("fatal", "error", err); os.Exit(1)`
- Add `import "log/slog"`; remove `import "log"`

### 3. `strings.Builder` for hot string concat

**Memory:** `feedback_string_builder_for_concat` — `+` is fine for cold paths; loops/hot paths need Builder.

**Detect:** look for `+=` on string variables inside `for` loops, or repeated `s = s + ...` patterns.
```bash
grep -rn --include='*.go' -P -B2 -A0 '\+= *"|s\s*=\s*s\s*\+' <target> | grep -B2 'for '
```

This is heuristic — manual review of matches is needed. Cold-path string concat with `+` is fine; only flag `+=` inside loops.

**Suggested fix:** rewrite the hot loop with `strings.Builder`:
```go
var b strings.Builder
b.Grow(estimatedSize)
for ... {
    b.WriteString(part)
}
result := b.String()
```

### 4. Empty-string check with `TrimSpace`

**Memory:** `feedback_strings_trimspace_for_empty_check` — `strings.TrimSpace(s) == ""` not bare `s == ""`.

**Detect:**
```bash
grep -rn --include='*.go' -P '\b\w+\s*==\s*""' <target>
```

**Allowed exceptions:**
- Comparing a slice/map element to `""` after a known typed conversion (e.g., `creds["apiKey"] == ""` where the source is already trimmed)
- `if val == ""` immediately after a `strings.TrimSpace` assignment

Not always mechanical — manual review per match.

**Suggested fix:**
```go
if strings.TrimSpace(s) == "" { ... }
```

### 5. `httpclient` over `net/http`

**Memory:** `feedback_use_httpclient_not_nethttp` — `httpclient.NewClientBuilder` for all outbound HTTP.

**Detect:**
```bash
grep -rn --include='*.go' -E '"net/http"|http\.(Get|Post|NewRequest|Client\b|DefaultClient)' <target>
```

**Allowed exceptions:**
- `internal/common/presentation/...` and any presentation-layer code that defines Fiber types (Fiber re-exports `http.Status*` constants — keep those)
- `*_test.go` may use `httptest`
- Webhook receivers and middleware that operate on incoming requests (these consume `*fiber.Ctx`, not `*http.Request`)

**Suggested fix:**
```go
client := httpclient.NewClientBuilder().
    Context(ctx).
    BaseURL("https://...").
    Header("Authorization", "Bearer "+token).
    Build()
response, err := client.Get("/path").Do(&successResponse, &errorResponse)
```

### 6. `go-packages/json` for JSON, not `encoding/json` (and no `utils` package)

**Memory:** `feedback_no_utils_package` — use `github.com/blocknextai/go-packages/json` (`json.Marshal` / `json.Unmarshal` / `json.RawMessage`); `encoding/json` only for library-required `json.RawMessage`. There is no `utils` package — it was removed; trivial stdlib wrappers live in domain packages (`go-packages/json`, `go-packages/idgen`, `go-packages/convert`, etc.).

**Detect:**
```bash
grep -rn --include='*.go' -E '"encoding/json"|json\.(Marshal|Unmarshal|NewDecoder|NewEncoder)' <target>
```

**Allowed exceptions:**
- `gjs.Schema` defaults that require `json.RawMessage` (e.g., node `node.go` files under `internal/nodeengine/nodes/...` and credential `*.go` files under `internal/nodeengine/credentials/`)
- Any direct `json.RawMessage` field in a struct that is consumed by an external library that requires the stdlib type — keep the import scoped to that field

**Suggested fix:**
```go
import (
   "github.com/blocknextai/go-packages/json"
)

bytes, err := json.Marshal(value)
err = json.Unmarshal(bytes, &target)
```

For `json.RawMessage` field types in domain/infra code, use `go-packages/json`'s `json.RawMessage` (a type alias of the stdlib type) by importing `go-packages/json` instead of `encoding/json`.

### 7. `append` over index assign

**Memory:** `feedback_use_append_not_index_assign` — `result = append(result, x)` over `result[i] = x` for slice building.

**Detect:** look for `make([]T, len(src))` followed by `result[i] = ...` in a `for i := range src` loop.
```bash
grep -rn --include='*.go' -P 'make\(\[\][^)]+, len\(' <target>
```

Then manually inspect each match for the index-assign pattern. A `make([]T, 0, len(src))` (preallocated to length zero with capacity) followed by `append` is the correct pattern.

**Suggested fix:**
```go
// Before:
result := make([]T, len(src))
for i, v := range src {
    result[i] = mapFn(v)
}

// After:
result := make([]T, 0, len(src))
for _, v := range src {
    result = append(result, mapFn(v))
}
```

### 8. `new(value)` over `&value`

**Memory:** `feedback_go_new_keyword` — Go 1.26+ idiom in this project.

**Detect:** `&` immediately followed by a literal value or a local variable being captured for nullable purposes.
```bash
grep -rn --include='*.go' -P '\b\w+\s*=\s*&[a-z][\w.]*' <target>
```

This pattern is heuristic; many `&entity` uses are correct (passing a pointer to an existing struct). Flag only the single-statement "create-and-take-pointer" cases:
```go
now := time.Now().UTC()
e.DeletedAt = &now   // ← flag, prefer new(now)
```

**Suggested fix:**
```go
e.DeletedAt = new(time.Now().UTC())
// or, when reusing:
e.DeletedAt = new(now)
```

### 9. `apperror` instead of `errors.New` / `fmt.Errorf`

**Memory:** `project_apperror_pattern` — Kind+Message+Cause via go-packages/apperror; presentation maps Kind→HTTP; no Code field.

**Detect:**
```bash
grep -rn --include='*.go' -E 'errors\.New\(|fmt\.Errorf\(' <target>
```

**Allowed exceptions:**
- `errors.Is(err, ...)` and `errors.As(...)` are unrelated (these are predicates, not constructors)
- Sentinel `errors.New(...)` in `*_test.go`
- `errors.New(...)` for narrow internal control-flow errors that never bubble to HTTP (rare; prefer apperror for consistency)

**Suggested fix:**
```go
// Domain layer (validation / not-found):
var ErrInvalidName = apperror.Validation("invalid name")
var ErrUserNotFound = apperror.NotFound("user not found")

// Application layer (wrapped infra failures):
var ErrFailedToCreateUser = apperror.Internal("failed to create user")
// usage: return ErrFailedToCreateUser.WithCause(err)
```

Map domain errors to the right `apperror` Kind:
- Invalid input → `apperror.Validation`
- Missing entity → `apperror.NotFound`
- Auth/permission → `apperror.Forbidden` or `apperror.Unauthorized`
- Infra/wrapped → `apperror.Internal`

### 10. `gofmt -w` after every edit

**Memory:** `feedback_run_gofmt_after_edits` — mandatory after Go edits.

**Detect:** can't grep for this — it's a workflow rule. Apply unconditionally at the end of the review.

```bash
gofmt -l <target>          # list files needing formatting
gofmt -w <target>          # apply formatting
```

If `-l` returns any path, those files need formatting. Either run `-w` (in fix mode) or include them in the report.

### 11. `golangci-lint run` after changes

**Memory:** `feedback_run_golangci_lint` — to catch issues build misses.

```bash
golangci-lint run <target_package>
```

Always part of the final report. Lint output is included verbatim in the punch list.

### 12. Repository pagination/search param order

**Memory:** `feedback_offset_limit_param_order` — repo signatures must be `(ctx, ...domainKeys, searchQuery string, offset int, limit int)`. Query structs must order `Search` field before `Pagination`.

**Detect:**
```bash
# Bad single-line repo signatures (offset before search, or limit before offset)
grep -rn --include='*.go' -E 'offset int, *limit int, *searchQuery string|limit int, *offset int' <target>

# Bad Query struct field order (Pagination before Search)
grep -rn --include='*.go' -B0 -A0 -E '^\s*Pagination\s+resultPkg\.PaginationRequest$' <target> \
  | while read line; do
      file="${line%%:*}"; lineno="${line#*:}"; lineno="${lineno%%:*}"
      # check next line for Search field
      sed -n "$((lineno+1))p" "$file" | grep -q 'Search ' && echo "$file:$lineno — Pagination appears before Search"
    done
```

**Allowed exceptions:**
- Functions that take only `offset, limit` with no `searchQuery` (e.g. `GetAllPublishers`) — already in the canonical order.
- Functions that take only `limit` and have semantic reasons (e.g. a worker batch-size parameter, not a pagination window).

**Suggested fix (manual — touches signature + every call site):**
```go
// Before
GetAllByOwner(ctx, ownerType, ownerID, offset int, limit int, searchQuery string)
// After
GetAllByOwner(ctx, ownerType, ownerID, searchQuery string, offset int, limit int)
```
Update the matching domain interface, infrastructure impl, and every handler call site in the same change.

### 13. Pagination/Search `.Normalize()` required in presentation

**Memory:** `feedback_pagination_search_normalization` — `request.PaginationRequest.Normalize()` / `request.SearchRequest.Normalize()` MUST be called in presentation handlers; never pass the un-normalized embedded `request.PaginationRequest` / `request.SearchRequest` directly to the query. (The old `NewPaginationRequest` / `NewSearchRequest` constructors no longer exist — `.Normalize()` replaced them.)

**Detect:**
```bash
# Direct passthrough to a Query — search/pagination embedded fields used without .Normalize()
grep -rn --include='*.go' -E 'Search:\s*request\.SearchRequest|Pagination:\s*request\.PaginationRequest' <target>
```

**Allowed exceptions:**
- None. `.Normalize()` does the normalization (offset clamp, limit default/max, query trim). Skipping it sends raw user input straight to SQL.

**Suggested fix (mechanical):**
```go
// Before
result, err := handler.Handle(ctx, &foo.Query{
    Search:     request.SearchRequest,
    Pagination: request.PaginationRequest,
})

// After
searchRequest := request.SearchRequest.Normalize()
paginationRequest := request.PaginationRequest.Normalize()

result, err := handler.Handle(ctx, &foo.Query{
    Search:     searchRequest,
    Pagination: paginationRequest,
})
```

For paginated list responses, return via the `commonHTTP.RespondPaginated(c, result.Items, result.TotalCount, paginationRequest)` helper — not a hand-rolled `resultPkg.NewPagination` / `WithPagination` envelope.

### 14. Initialisms uppercase in Go identifiers

**Memory:** `feedback_initialism_rename_autonomous` — Go identifiers use uppercase initialisms (`ID`, `URL`, `API`, `HTTP`, `JSON`, `OAuth2`).

**Detect:**
```bash
# Common offenders: variable / field / param / package-alias casing
grep -rn --include='*.go' -E '\b(userId|orgId|apiKeyId|credentialId|workflowId|sessionId|linkedAccountId|organizationId|tenantId)\b' <target>
grep -rn --include='*.go' -E '\b(commonHttp|commonHttpV2|httpUrl|baseUrl|callbackUrl|redirectUrl)\b' <target>
grep -rn --include='*.go' -E '"github\.com/[^"]+/http"\s+commonHttp' <target>
```

**Allowed exceptions:**
- Local variables that match a third-party library's parameter name (rare).

**Suggested fix (mechanical):**
```go
userId  := commonHttp.GetUserID(c)   // ❌
userID  := commonHTTP.GetUserID(c)   // ✅
```
Rename touches the import alias, every reference, and any local that shadows it. Run `gofmt` after to clean up alignment.

### 15. Don't apply initialism rules to struct tag CONTENT

**Memory:** `feedback_no_initialism_in_struct_tags` — Go field becomes `OrganizationID`, but `json:"organizationId,omitempty"` stays unchanged. Same for `db`, `env`, `uri`, `query`, `schema` tag values.

**Detect:**
```bash
# Tag content with uppercase initialisms — likely an over-eager rename
grep -rn --include='*.go' -E '`[a-z]+:"[a-z]*ID(,| |"|\")|`[a-z]+:"[a-z]*URL(,| |"|\")|`[a-z]+:"[a-z]*API(,| |"|\")' <target>
```

**Allowed exceptions:**
- Tag contents that intentionally encode an external API's casing (e.g., a third-party header `"X-API-Key"`). These are rare and obvious from context.

**Suggested fix (mechanical):**
```go
// Before
OrganizationID uuid.UUID `json:"organizationID" db:"organization_id" uri:"organizationID"`
// After
OrganizationID uuid.UUID `json:"organizationId" db:"organization_id" uri:"organizationId"`
```
Only the Go identifier follows the initialism rule; tag values keep their original casing (camelCase for json/uri/query, snake_case for db/env, etc.).

## Workflow

### Step 1: Resolve target

If the user gave a file path, use it. If a directory, glob `*.go`. If "diff", run:
```bash
git diff main...HEAD --name-only -- '*.go'
```
Filter to only existing files (deleted ones don't apply).

### Step 2: Run all detectors in parallel

Each rule's grep pattern can be run independently. Issue all greps in a single message for efficiency. Collect raw results.

### Step 3: Apply allowed-exception filtering

For each match, check whether it falls under that rule's allowed-exception list. The exceptions are listed inline above; do not invent new ones — if a match looks suspicious but plausibly correct, flag it for human review rather than silently dropping it.

### Step 4: Report

Format as a punch list:

```markdown
## Go Standards Review — `<target>`

### Rule violations (N)

#### Rule 1: No `fmt` package (3 violations)
- `internal/foo/bar.go:42` — `fmt.Sprintf("user %s", id)` → use `"user " + id` or `strings.Builder`
- `internal/foo/baz.go:17` — `fmt.Errorf(...)` → convert to `apperror.Internal(...).WithCause(err)`
- ...

#### Rule 5: `httpclient` over `net/http` (1 violation)
- `internal/qux/client.go:8` — imports `net/http` → migrate to `httpclient.NewClientBuilder()`

### Workflow checks
- gofmt: 2 files need formatting (`internal/foo/bar.go`, `internal/foo/baz.go`)
- golangci-lint: 0 issues
```

Group by rule; under each rule list `file:line — code excerpt — suggested fix`. Include rule number for cross-referencing.

If there are zero violations and lint is clean, report a single line: `No violations. gofmt clean. golangci-lint: 0 issues.`

### Step 5 (optional, fix mode): Apply mechanical fixes

For each violation that has a deterministic mechanical rewrite, apply via `Edit`. The mechanical-fix list:

- **Rule 4 (TrimSpace):** safe — wrap with `strings.TrimSpace(...)`. Add the `strings` import if missing.
- **Rule 8 (`new(value)`):** safe when the value is a freshly-bound local variable on the immediately preceding line.
- **Rule 10 (gofmt):** always safe.
- **Rule 6 (`encoding/json` → `go-packages/json`):** safe for `json.Marshal` / `json.Unmarshal` calls (swap the import to `github.com/blocknextai/go-packages/json`); **NOT** safe for `json.RawMessage` field types — surface those for manual review.
- **Rule 13 (Pagination/Search `.Normalize()`):** safe — insert the two `request.SearchRequest.Normalize()` / `request.PaginationRequest.Normalize()` lines and replace the direct passthroughs in the query literal.
- **Rule 14 (Initialism Go identifiers):** safe to rename `userId` → `userID`, `commonHttp` → `commonHTTP`, etc., across the target. The compiler will catch any miss.
- **Rule 15 (Initialism in tag content):** safe to revert `json:"organizationID"` → `json:"organizationId"` — undoing an over-eager Rule-14 sweep.

For everything else (Rule 1's `fmt.Sprintf`, Rule 5's `net/http`, Rule 9's `errors.New`/`fmt.Errorf`, Rule 7's index-assign, Rule 3's hot loop, Rule 12's signature reorder which touches every call site), the fix requires judgement. Leave these in the report.

After fixes, re-run gofmt + golangci-lint. Show the diff summary (`X files changed, Y insertions, Z deletions`).

## Critical rules

1. **Never apply a non-mechanical fix in `fix` mode.** Judgement-call rewrites belong to a human review.
2. **Never silence violations by editing the rule pattern.** If a match is a real exception, document it under that rule's "Allowed exceptions" list — but do that as a separate skill update, not inline during a review.
3. **Always run `gofmt` and `golangci-lint` last,** even in report mode. They surface issues independent of the rule list.
4. **Don't expand the rule set without memory backing.** Every rule in this skill traces to a memory entry. New rules added to the project's standards should land in MEMORY.md first, then propagate here.
5. **Don't flag `_test.go` files for the same rules as production code.** Test ergonomics warrants `fmt.Errorf`, `errors.New`, `httptest`, etc. The skill's allowed-exception lists already note this; respect them.

## Final steps

After any review (regardless of mode):

1. Print the punch list.
2. If in `fix` mode: re-run `gofmt -l` to verify clean, then `golangci-lint run` to verify clean.
3. Hand back to the user with: "N violations remain after fixes — manual review needed for rules X, Y, Z."

## What NOT to do

- Do NOT add new rules to this skill from "general Go best practices." Every rule is project-specific and tied to a memory entry. If a rule isn't in MEMORY.md, it doesn't belong here.
- Do NOT auto-fix `fmt.Sprintf` → `strings.Builder` mechanically. Choosing between `+` (cold) and `Builder` (hot) requires reading the surrounding code.
- Do NOT auto-fix `errors.New` / `fmt.Errorf` → `apperror.<Kind>`. Picking the right Kind requires reading the call site.
- Do NOT run on the entire `internal/...` tree without explicit user request — large targets produce noise. Default to a single file or package.
- Do NOT modify `.golangci.yml` to silence lint rules during a review. If a lint rule is genuinely wrong, that's a separate change with its own discussion.
- Do NOT skip the gofmt/lint final pass even when the report is empty — these catch class-of-issue, not pattern-match.
