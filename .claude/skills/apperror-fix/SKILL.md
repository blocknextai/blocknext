---
name: apperror-fix
description: Use this skill when the user wants to convert raw Go errors (`errors.New`, `fmt.Errorf`, plain string returns) in a target to the project's `apperror` Kind+Message+Cause pattern. Triggers on phrases like "convert errors to apperror in X", "fix error handling in Y package", "apperror migration for Z module", "audit non-apperror errors", "introduce sentinels for X". Focused complement to `go-standards-review`'s Rule 9 — this skill does the judgement-aware Kind selection that the review skill leaves to humans.
---

> **Working directory:** all relative paths in this skill (`internal/...`, `cmd/...`, `go build ./...`) are relative to `apps/platform-api/` in the monorepo. Make targets run from the repo root.


# Apperror Fix

Converts raw Go error usage (`errors.New(...)`, `fmt.Errorf(...)`, plain inline strings, ad-hoc `error` returns) in a target file or package to the project's `apperror` pattern: domain sentinels + application sentinels + `WithCause(err)` wrapping.

## When to use

- User says "convert errors to apperror in `<package>`" / "fix error handling in `<file>`" / "apperror migration for `<module>`" / "audit non-apperror errors"
- After a refactor that pulled in older code from a non-apperror codebase
- When a new module skipped the sentinel layer and inlined `errors.New` / `fmt.Errorf` calls

## What this skill does NOT do

- Does not change any other Go-standard rule (no `fmt` migration, no `slog` migration). For broad audits use `go-standards-review`.
- Does not invent error wording. The skill proposes Messages based on context; the user picks.
- Does not auto-pick `Kind` for ambiguous cases. Each conversion's Kind is presented as a choice with a recommended default.

## Required information

Confirm with the user (ask only what's missing):

1. **Target** — file path, package directory, or "diff" (changes vs `main`).
2. **Mode** — `report` (default; show all findings as a punch list) or `apply` (rewrite each finding after Kind decision).
3. **Sentinel placement strategy** — for each new sentinel:
   - **Domain sentinel** (`apperror.Validation` / `apperror.NotFound`) → goes in `internal/<module>/domain/<aggregate>/errors.go`
   - **App sentinel** (`apperror.Internal`) → goes in `internal/<module>/application/<aggregate>/errors.go`
   - **Inline anonymous** (rare, prefer named sentinels) → only if the error never bubbles outside the file

## The `apperror` taxonomy

The project uses six Kinds (verified from existing code):

| Kind                      | When to use                                         | HTTP mapping (typical) |
| ------------------------- | --------------------------------------------------- | ---------------------- |
| `apperror.Validation`     | Invalid input shape / value (bad enum, empty name)  | 400                    |
| `apperror.NotFound`       | Entity does not exist (or owned by someone else)    | 404                    |
| `apperror.Forbidden`      | Authenticated but not authorized for the action     | 403                    |
| `apperror.Unauthorized`   | No credentials / invalid credentials                | 401                    |
| `apperror.Conflict`       | Resource state conflicts with operation             | 409                    |
| `apperror.Internal`       | Infra / external-service / unexpected failure       | 500                    |

**Memory:** `project_apperror_pattern` — Kind + Message + Cause via `go-packages/apperror`; presentation maps Kind → HTTP; **no `Code` field**.

## Layering

| Layer            | Sentinel kind                                  | Naming                              |
| ---------------- | ---------------------------------------------- | ----------------------------------- |
| Domain           | `Validation`, `NotFound`, `Conflict`           | `ErrInvalid<Thing>`, `Err<Thing>NotFound`, `Err<Thing>Conflict` |
| Application      | `Internal` (wraps infra), occasionally others  | `ErrFailedTo<Verb><Thing>`          |
| Presentation     | usually doesn't add new sentinels              | reuses domain/app sentinels         |
| Infrastructure   | doesn't define sentinels — translates SQL errors → domain `NotFound` | n/a                  |

**Memory:** `feedback_validation_layering` — sentinels in domain, app `Validate()` runs BEFORE tx for fail-fast.

## Detection

```bash
grep -rn --include='*.go' -E 'errors\.New\(|fmt\.Errorf\(|return .*errors\.New|return .*fmt\.Errorf' <target>
```

For each match, read 5–10 surrounding lines to determine:
- **Layer** (domain / app / presentation / infrastructure)
- **Cause** (is there a wrapped `err`? `%w` in `fmt.Errorf` is the give-away)
- **Semantic Kind** (validation / not-found / internal / etc.)

Also flag:
- `return nil, <some-error-construction>` patterns
- Plain `error` types declared as `type ... error` (rare; usually a sign of pre-apperror code)

## Decision matrix

For each match, propose the conversion:

### Pattern A: `errors.New("...")` with no wrapping

This is a sentinel candidate. Promote to a named `var Err... = apperror.<Kind>("...")` in the appropriate `errors.go` file. Replace the inline call with the sentinel reference.

```go
// Before:
return nil, errors.New("user id is required")

// After (validator.go):
return nil, ErrInvalidUserID
// where errors.go declares:
var ErrInvalidUserID = apperror.Validation("invalid user id")
```

**Recommended Kind by context:**
- Inside a `Validate()` method or a domain `New(...)`/`Update(...)` constructor → `Validation`
- Inside an ownership/existence check → `NotFound` (use the existing `Err<Entity>NotFound` if defined)
- Inside an auth middleware / permission check → `Forbidden` or `Unauthorized`
- Default if unclear → `Internal` (and tell the user)

### Pattern B: `fmt.Errorf("...: %w", err)` with cause

This is a wrap. Use an app-level sentinel + `.WithCause(err)`.

```go
// Before:
return fmt.Errorf("failed to insert user: %w", err)

// After (handler.go inside transaction):
return ErrFailedToCreateUser.WithCause(err)
// where application/<aggregate>/errors.go declares:
var ErrFailedToCreateUser = apperror.Internal("failed to create user")
```

**Recommended Kind:**
- Wrapped repository / SQL / external API failure → `Internal`
- Wrapped JSON marshal/unmarshal failure → `Internal`
- Wrapped encryption/decryption failure → `Internal` (existing pattern: `ErrFailedToEncryptData` / `ErrFailedToDecryptData`)
- Wrapped `Validate()` failure during user input handling → don't wrap; propagate the original Validation sentinel

### Pattern C: `fmt.Errorf("... %s", x)` with no `%w`

This is a formatting + new error case. Two-step:
1. If the message is a known canned shape (e.g., `"%s is required"`), promote to a sentinel without the format — the field name is already known statically at the call site.
2. If the message genuinely needs the runtime value, prefer a sentinel + add the value via `slog.Error(...)` for debugging, OR use `apperror.<Kind>(msg)` with the formatted string as the message.

```go
// Before:
return fmt.Errorf("invalid status %s", status)

// After (preferred when status values are bounded):
return ErrInvalidStatus  // sentinel; details go to slog at the boundary

// After (acceptable when the value must surface to the API):
return apperror.Validation("invalid status: " + status)  // string concat, not fmt
```

The second form uses `+` concat (memory: `feedback_no_fmt_package` — `fmt` not allowed).

### Pattern D: `errors.New("...")` returned from a domain method

If the message is one of: `invalid X`, `X is required`, `X must be ...` — these are domain invariants. Always Validation:

```go
// Before:
if k.Name == "" {
    return errors.New("api key name is required")
}

// After:
if strings.TrimSpace(k.Name) == "" {
    return ErrInvalidAPIKeyName
}
```

### Pattern E: `errors.Is(err, sql.ErrNoRows)` translation in repository

Standard infrastructure pattern. Translate to the domain not-found sentinel:

```go
// Before:
if err == sql.ErrNoRows {
    return nil, errors.New("not found")
}

// After:
if errors.Is(err, sql.ErrNoRows) {
    return nil, <module>Domain<Aggregate>.Err<Entity>NotFound
}
```

`errors.Is` / `errors.As` are PREDICATES — keep them; they're not violations.

## Workflow

### Step 1: Detect

```bash
grep -rn --include='*.go' -E 'errors\.New\(|fmt\.Errorf\(' <target>
```

Filter out:
- `*_test.go` (test ergonomics)
- `errors.Is(...)`, `errors.As(...)`, `errors.Unwrap(...)` (predicates, not constructors)

### Step 2: Classify each match

For each match, read the surrounding code and assign:
- **Layer** (D = domain, A = application, P = presentation, I = infrastructure)
- **Pattern** (A / B / C / D / E from above)
- **Recommended Kind** + **Recommended sentinel name** + **Recommended location**

If you cannot confidently assign, mark the match as `??` and surface to the user.

### Step 3: Locate / create errors.go

For each new sentinel, check whether `internal/<module>/domain/<aggregate>/errors.go` (domain) or `internal/<module>/application/<aggregate>/errors.go` (application) exists.

- If it exists, append the sentinel to the `var (...)` block. Group by Kind: keep all Validation sentinels together, then NotFound, then Conflict (rarely), then Internal in the app file.
- If it does not exist, create it with the standard header:

```go
package <aggregate>

import (
    "github.com/blocknextai/go-packages/apperror"
)

var (
    Err<Name> = apperror.<Kind>("<message>")
)
```

### Step 4: Report

```markdown
## Apperror Fix — `<target>`

### Findings (N)

#### `internal/foo/bar/handler.go:42` — Pattern B (wrapped infra failure)
```go
// Before:
return fmt.Errorf("failed to save user: %w", err)
// After:
return ErrFailedToSaveUser.WithCause(err)
```
- New sentinel: `ErrFailedToSaveUser = apperror.Internal("failed to save user")`
- Location: `internal/foo/application/users/errors.go` (file does not exist; will be created)

#### `internal/foo/domain/users/user.go:73` — Pattern A (validation in constructor)
```go
// Before:
return errors.New("user email is required")
// After:
return ErrInvalidUserEmail
```
- New sentinel: `ErrInvalidUserEmail = apperror.Validation("invalid user email")`
- Location: `internal/foo/domain/users/errors.go` (existing; appending)

### Ambiguous (M) — please confirm Kind
- `internal/foo/bar/baz.go:88` — `errors.New("token mismatch")` — Validation? Unauthorized? Conflict?
```

### Step 5 (apply mode): Edit

For each confirmed conversion:
1. Open or create the target `errors.go` and add the sentinel.
2. Open the call site and replace the `errors.New(...)` / `fmt.Errorf(...)` call with the sentinel reference (or `.WithCause(err)` for wrapped cases).
3. Add or remove imports as needed (`errors` may be removable if only used for `errors.New`; `fmt` likewise).

After every batch of edits:
```bash
gofmt -w <files>
go build ./...
golangci-lint run <package>
```

### Step 6: Verify

Re-run the detection grep on the target. The remaining matches should be either:
- In `*_test.go` (allowed)
- Predicates (`errors.Is` / `errors.As` / `errors.Unwrap`)
- Genuinely-ambiguous cases the user opted to keep as-is

## Critical rules (project memory)

1. **Kind + Message + Cause only.** No `Code` field. Don't introduce error codes — the project's apperror is opinionated against them.
2. **Domain sentinels are Validation / NotFound / Conflict.** App-level sentinels are Internal. Never put `apperror.Internal` sentinels in `domain/<aggregate>/errors.go`.
3. **Wrap with `.WithCause(err)`, not by including the original error in the message.** The Cause is preserved structurally for downstream Sentry / logging.
4. **Don't wrap a `Validation`/`NotFound` with another `Internal`.** If the inner error is already an apperror with the right Kind, propagate it.
5. **Validate runs before tx** (memory: `feedback_validation_layering`). When converting validators, the `errors.New` calls almost always become Validation sentinels in `domain/<aggregate>/errors.go`.
6. **Use `strings.TrimSpace(s) == ""` for empty checks** (memory: `feedback_strings_trimspace_for_empty_check`). Often the bare `s == ""` check appears alongside `errors.New("X is required")` — fix both in the same edit.
7. **Don't reintroduce `fmt`** while fixing apperror. If the original message used `fmt.Errorf("invalid X %s", x)`, prefer a sentinel without the value, OR use `apperror.<Kind>("invalid x: " + x)` with `+` concat.

## What NOT to do

- Do NOT batch this with unrelated refactors. Apperror migration is a focused change with reviewable scope.
- Do NOT pick `Internal` as the default Kind for everything. Default to the right Kind for the layer; `Internal` is for infra failures specifically.
- Do NOT delete original error wording that captures useful detail. The sentinel's Message can absorb the static part; the dynamic part can go to `slog.Error(...)` at the call site.
- Do NOT touch `*_test.go` — test ergonomics warrants raw `errors.New`. Skip them in the audit.
- Do NOT remove `errors.Is` / `errors.As` / `errors.Unwrap`. They're predicates, not constructors.
- Do NOT collapse multiple distinct error conditions into one sentinel. `ErrInvalidName` and `ErrNameTooLong` are separate failures and deserve separate sentinels (the consumer may want to differentiate).
- Do NOT introduce new Kinds beyond the six listed. The mapping to HTTP status is handled by the apperror package; making up new Kinds breaks that contract.
