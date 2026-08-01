---
name: command-create
description: Use this skill when the user wants to add a new write-side CQRS command (Create / Update / Delete / regenerate / etc.) to a bounded context in platform-api. Triggers on phrases like "add a `<verb><entity>` command", "create a new command for X", "let's add the delete X handler", "scaffold an update command for X". Generates command/validator/handler/response files, wires them through `infrastructure.go`, the presentation layer handler, and the route registration.
---

> **Working directory:** all relative paths in this skill (`internal/...`, `cmd/...`, `go build ./...`) are relative to `apps/platform-api/` in the monorepo. Make targets run from the repo root.


# Command Create

Scaffolds a new CQRS command under `internal/<module>/application/<aggregate>/<verb><entity>/`, wires it into the module's `infrastructure.go`, adds the Fiber HTTP handler in `presentation/<aggregate>/<verb>_<entity>.go`, and registers the route in `presentation/presentation.go`.

## When to use

- User says "add a `<verb><entity>` command for `<module>`" / "let's add `delete<Entity>`" / "create the update handler for X"
- A write-side action is needed: create, update, delete, regenerate, archive, transfer, etc.
- For read-side (GET) endpoints use the `query-create` skill instead.

## Required information

Confirm with the user (ask only what's missing):

1. **Module** — bounded context that owns the aggregate (e.g., `apikeys`, `organizations`, `credentials`, `workflows`). Must already exist under `internal/`.
2. **Aggregate** — sub-folder under `application/` that groups commands for one entity (e.g., `apikeys`, `organizations`, `organizationusers`). Usually plural, lowercase, single word.
3. **Verb** — `create`, `update`, `delete`, `regenerate`, `transfer`, `archive`, etc. Lowercase.
4. **Entity** — singular noun matching the aggregate root (e.g., `apikey`, `organization`, `credential`).
5. **Command fields** — what the caller must provide (name, ID, OwnerType/OwnerID, payload). Their Go types.
6. **Response shape** — what the handler returns. Often `ID + a few echoed fields`. Delete commands typically return a Response with just the deleted entity's `ID` (and the presentation layer responds `200 OK` with the standard envelope).
7. **Domain dependencies** — the repository (always required), plus any cross-context services (e.g., `secretManager`, `nodeEngineCredentialService`). Identified by the existing handlers in the same module.
8. **Auth & route** — HTTP verb (`POST` / `PATCH` / `DELETE`), URL path, RBAC permission constant from `go-packages/rbac`, and whether it lives under user (`/users/me/...`), organization (`/organizations/:organizationId/...`), or both groups.

## File layout

```
internal/<module>/application/<aggregate>/<verb><entity>/
├── command.go     # request struct (no methods)
├── validator.go   # Validate() method on the command — sentinels only
├── handler.go     # Handler struct + New + Handle (no manual Validate call)
└── response.go    # Response struct (delete commands include this too — typically just ID)
```

Plus three wiring touch-points outside the command folder:
- `internal/<module>/infrastructure/infrastructure.go` — add to `Handlers` struct + instantiate.
- `internal/<module>/presentation/<aggregate>/<verb>_<entity>.go` — Fiber HTTP handler. **Filename uses snake_case.**
- `internal/<module>/presentation/presentation.go` — register route under the right group with the right RBAC.

## Naming rules

| Thing                      | Convention                                                  | Example                          |
| -------------------------- | ----------------------------------------------------------- | -------------------------------- |
| Folder / package           | `<verb><entity>` lowercase, no separator                    | `createapikey`, `updatecredential` |
| Command struct             | `<Verb><Entity>Command` PascalCase                          | `CreateAPIKeyCommand`            |
| Handler struct             | `Handler` (package name carries the context — no stutter)   | `createapikey.Handler`           |
| Constructor                | `New`                                                       | `createapikey.New`               |
| Response struct            | `<Verb><Entity>Response`                                    | `CreateAPIKeyResponse`           |
| Presentation file          | `<verb>_<entity>.go` snake_case                             | `create_api_key.go`              |
| Presentation handler ctor  | `New<Owner><Verb><Entity>Handler`                           | `NewCreateUserAPIKeyHandler`     |

Acronyms (`API`, `URL`, `ID`, `UUID`) stay uppercase in Go identifiers; lowercase in folder names. So `createapikey/` package contains `CreateAPIKeyCommand`.

## Templates

### `command.go`

```go
package <verb><entity>

import (
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	"github.com/google/uuid"
)

type <Verb><Entity>Command struct {
	OwnerType commonDomain.OwnerType
	OwnerID   uuid.UUID
	// ... request-specific fields
}
```

Rules:
- No methods on the command struct except `Validate()` (lives in `validator.go`).
- For update/delete commands, the entity ID is a `uuid.UUID` field (typically `ID`).
- For commands scoped to a user vs organization, include `OwnerType` and `OwnerID` so the same handler serves both presentation paths.
- Do NOT add JSON tags — the command is internal; HTTP binding happens in the presentation request struct.

### `validator.go`

```go
package <verb><entity>

import (
	"strings"

	<module>Domain<Aggregate> "github.com/blocknextai/platform-api/internal/<module>/domain/<aggregate>"
	"github.com/google/uuid"
)

const (
	MaxNameLength = 255
)

func (c *<Verb><Entity>Command) Validate() error {
	if !c.OwnerType.IsValid() {
		return <module>Domain<Aggregate>.ErrInvalidOwnerType
	}

	if c.OwnerID == uuid.Nil {
		return <module>Domain<Aggregate>.ErrInvalidOwnerID
	}

	if strings.TrimSpace(c.Name) == "" {
		return <module>Domain<Aggregate>.ErrInvalid<Entity>Name
	}

	if len(strings.TrimSpace(c.Name)) > MaxNameLength {
		return <module>Domain<Aggregate>.Err<Entity>NameTooLong
	}

	return nil
}
```

Rules (project memory):
- Empty-string check is `strings.TrimSpace(s) == ""`, never `s == ""`.
- Length checks use `len(strings.TrimSpace(s)) > N`.
- Returned errors are domain sentinels declared in `internal/<module>/domain/<aggregate>/errors.go` using `apperror.Validation(...)`. If the right sentinel doesn't exist yet, add it there.
- `Validate()` is invoked automatically by the `cqrs.ValidationBehavior` wrapper (applied in `infrastructure.go`) **before** the transaction — NOT called in the handler body. Keep this method regardless of any similar HTTP-layer check; it is the single source of validation rules. (Memory: `project_cqrs_validation_behavior`.)
- Constants for max lengths live as `const` at the top of `validator.go`.

### `handler.go`

```go
package <verb><entity>

import (
	"context"

	"github.com/blocknextai/go-packages/database"
	<module>Application<Aggregate> "github.com/blocknextai/platform-api/internal/<module>/application/<aggregate>"
	<module>Domain<Aggregate> "github.com/blocknextai/platform-api/internal/<module>/domain/<aggregate>"
)

type Handler struct {
	<aggregate>Repository <module>Domain<Aggregate>.<Entity>Repository
	transactionManager    database.TransactionManager
}

func New(
	<aggregate>Repository <module>Domain<Aggregate>.<Entity>Repository,
	transactionManager database.TransactionManager,
) *Handler {
	return &Handler{
		<aggregate>Repository: <aggregate>Repository,
		transactionManager:    transactionManager,
	}
}

func (h *Handler) Handle(ctx context.Context, command *<Verb><Entity>Command) (*<Verb><Entity>Response, error) {
	var response *<Verb><Entity>Response

	err := h.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		// CREATE pattern: build the aggregate and persist it.
		entity, err := <module>Domain<Aggregate>.New(command.OwnerType, command.OwnerID, command.Name)
		if err != nil {
			return err
		}

		if err := h.<aggregate>Repository.Create(txCtx, entity); err != nil {
			return <module>Application<Aggregate>.ErrFailedToCreate<Entity>.WithCause(err)
		}

		response = &<Verb><Entity>Response{
			ID:   entity.ID,
			Name: entity.Name,
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}
```

For **update / delete**, the body inside the transaction is different:

```go
entity, err := h.<aggregate>Repository.GetByID(txCtx, command.<Entity>ID)
if err != nil {
	return <module>Domain<Aggregate>.Err<Entity>NotFound
}

// Ownership guard — return NotFound (not Forbidden) to avoid leaking existence.
if entity.OwnerType != command.OwnerType || entity.OwnerID != command.OwnerID {
	return <module>Domain<Aggregate>.Err<Entity>NotFound
}

// UPDATE: call entity.Update(...) which returns the new value (immutable update style)
updatedEntity, err := entity.Update(command.Name)
if err != nil {
	return err
}
if err := h.<aggregate>Repository.Update(txCtx, updatedEntity); err != nil {
	return <module>Application<Aggregate>.ErrFailedToUpdate<Entity>.WithCause(err)
}

// DELETE: call the domain Delete() method (marks the entity deleted), then pass the
// returned entity to the repo. Repos consistently take the FULL entity, not the ID.
deletedEntity, err := entity.Delete()
if err != nil {
	return err
}
if err := h.<aggregate>Repository.Delete(txCtx, deletedEntity); err != nil {
	return <module>Application<Aggregate>.ErrFailedToDelete<Entity>.WithCause(err)
}
```

Rules (project memory):
- **Do NOT call `command.Validate()` inside the handler.** Validation is centralized: handlers are wrapped at wiring time by `cqrs.ValidationBehavior(...)`, whose `Handle` runtime-checks `Validatable` and runs `Validate()` before delegating to the inner handler (so it still runs before the transaction). The `Validate()` method itself STAYS in `validator.go` — it's the single source of validation rules; only its invocation moved out of the handler body. (Memory: `project_cqrs_validation_behavior`, `feedback_validation_layering`.)
- **Ownership guard returns NotFound, not Forbidden** — leaking existence is worse than the slight UX cost.
- **Repository errors get wrapped** with `apperror.WithCause(err)` via app-level sentinels (`ErrFailedToCreate<Entity>`, `ErrFailedToUpdate<Entity>`, etc.) declared in `internal/<module>/application/<aggregate>/errors.go`. Add new ones there if missing.
- **Domain "not found" / "invalid" errors are returned directly** (no wrapping) — they already carry the right `apperror.Kind`.
- **Constructors do not take repositories from other bounded contexts.** If you need data from another module, accept its **service interface**, not its repository. (Memory: `feedback_no_repository_injection_across_contexts`.)
- **Use `new(value)` not `&value`** for pointer construction (memory: `feedback_go_new_keyword`).

### `response.go`

```go
package <verb><entity>

import (
	"github.com/google/uuid"
)

type <Verb><Entity>Response struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	// ...
}
```

For delete commands, the Response struct is still present and typically holds just the deleted entity's `ID` (some modules echo extra fields — match the existing delete handler in the same module). Presentation returns `200 OK` with the standard `resultPkg.Ok(result, resultPkg.WithMessage("<entity> deleted"))` envelope.

## Wiring

### 1. `internal/<module>/infrastructure/infrastructure.go`

Add the new handler to the `Handlers` struct and instantiate it in `RegisterInfrastructure`. Order matches struct field order; alphabetical within the struct is fine but match the existing module's convention (most modules list in feature-importance order, not strictly alphabetical).

`RegisterInfrastructure` takes a single `RegisterInfrastructureDeps` struct (NOT positional params); read deps off it as `deps.TransactionManager`, `deps.<Aggregate>Repository`, etc.

The `Handlers` struct fields hold the `cqrs.Handler[*Command, *Response]` interface (NOT the concrete `*Handler`), and each constructor is wrapped with `cqrs.ValidationBehavior(...)` so every handler goes through the validation behavior. Import `"github.com/blocknextai/platform-api/internal/common/application/cqrs"`.

```go
type Handlers struct {
	Create<Entity> cqrs.Handler[*createentity.Create<Entity>Command, *createentity.Create<Entity>Response]
	// ... existing
	<Verb><Entity>  cqrs.Handler[*<verb><entity>.<Verb><Entity>Command, *<verb><entity>.<Verb><Entity>Response]  // ← new
}

type RegisterInfrastructureDeps struct {
	TransactionManager database.TransactionManager

	<Aggregate>Repository <module>Domain<Aggregate>.<Entity>Repository
}

func RegisterInfrastructure(deps RegisterInfrastructureDeps) *Handlers {
	return &Handlers{
		Create<Entity>: cqrs.ValidationBehavior(createentity.New(deps.<Aggregate>Repository, deps.TransactionManager)),
		<Verb><Entity>:  cqrs.ValidationBehavior(<verb><entity>.New(deps.<Aggregate>Repository, deps.TransactionManager)),  // ← new
	}
}
```

`cqrs.ValidationBehavior` infers its type parameters from the concrete handler (Go 1.26), so no explicit type args are needed at the call site. The wrapping is uniform — apply it to query handlers too (see `query-create`), not just commands.

`cqrs` is a generic **behavior pipeline**, not a validation-only helper: `ValidationBehavior` is one `cqrs.Behavior[C, R]` (a `func(next Handler[C, R]) Handler[C, R]` decorator). Validation is currently the only behavior; future cross-cutting behaviors (logging, metrics, timing) are added as more `cqrs.Behavior` values and composed by nesting — `cqrs.Logging(cqrs.ValidationBehavior(handler))` — or via `cqrs.Chain(handler, behaviorA, behaviorB)`. Use `cqrs.ValidationBehavior(handler)` (nested form) so type params infer.

> **Unmigrated modules:** if the target module's `infrastructure.go` still stores concrete `*<Verb><Entity>CommandHandler` types (not yet migrated to `cqrs.Handler`), match the module's existing convention OR migrate the whole module to the `cqrs.ValidationBehavior` + `New`/`Handler` pattern first. Do not leave a module half-migrated. (Memory: `project_cqrs_validation_behavior`.)

If your handler has additional dependencies (secret manager, other services), add them as fields on `RegisterInfrastructureDeps` AND on the module's `Dependencies` struct in `module.go`. If they aren't there yet, this is a real change; surface it to the user and confirm before adding cross-module wiring.

### 2. `internal/<module>/presentation/<aggregate>/<verb>_<entity>.go`

```go
package <aggregate>

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/<module>/application/<aggregate>/<verb><entity>"
)

type <Verb>User<Entity>Request struct {
	Name string `json:"name"`
}

func New<Verb>User<Entity>Handler(handler cqrs.Handler[*<verb><entity>.<Verb><Entity>Command, *<verb><entity>.<Verb><Entity>Response]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(<Verb>User<Entity>Request)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		result, err := handler.Handle(c.RequestCtx(), &<verb><entity>.<Verb><Entity>Command{
			OwnerType: commonDomain.OwnerTypeUser,
			OwnerID:   userID,
			Name:      request.Name,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusCreated).JSON(resultPkg.Ok(result, resultPkg.WithMessage("<entity> <pasttense>")))
	}
}

type <Verb>Organization<Entity>Request struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	Name           string    `json:"name"`
}

func New<Verb>Organization<Entity>Handler(handler cqrs.Handler[*<verb><entity>.<Verb><Entity>Command, *<verb><entity>.<Verb><Entity>Response]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(<Verb>Organization<Entity>Request)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := handler.Handle(c.RequestCtx(), &<verb><entity>.<Verb><Entity>Command{
			OwnerType: commonDomain.OwnerTypeOrganization,
			OwnerID:   request.OrganizationID,
			Name:      request.Name,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusCreated).JSON(resultPkg.Ok(result, resultPkg.WithMessage("<entity> <pasttense>")))
	}
}
```

Rules:
- Status codes: `201 Created` for create, `200 OK` for update / regenerate / delete. (A few legacy endpoints return `204 No Content` for delete — match the module's existing convention if there is one; otherwise default to `200 OK` with the standard envelope.)
- The presentation Request struct has JSON tags (request body) and/or `uri:`/`query:` tags (path/query params).
- `c.Bind().All(request)` binds path + query + body. If only one source is needed, `c.Bind().Body(request)` etc. work, but `.All()` is the project default.
- `commonHTTP.GetUserID(c)` is the only sanctioned way to read the authenticated user; do not pull from headers manually.
- `resultPkg.Ok(result, resultPkg.WithMessage("<message>"))` is the standard envelope — the message is passed via the `resultPkg.WithMessage(...)` option, NOT as a bare string. Match wording to existing handlers in the module.
- Path parameters use `uuid.UUID` and bind via `uri:"<param>"`.

### 3. `internal/<module>/presentation/presentation.go`

Authentication and authorization are SEPARATE composable middlewares. Every authenticated route chains `authMiddleware.Authenticate()` first, then the scope-specific permission middleware:

```go
// user-scoped route (RegisterUser<Aggregate>Presentation, group /users/me/...)
<aggregate>RouterGroup.Post(  // or .Patch / .Delete
	"/",  // or "/:<entityId>"
	authMiddleware.Authenticate(),
	authMiddleware.RequireUserPermission(rbac.<Verb><Entity>Permission),
	<aggregate>.New<Verb>User<Entity>Handler(handlers.<Verb><Entity>),
)

// organization-scoped route (RegisterOrganization<Aggregate>Presentation, group /organizations/:organizationId/...)
<aggregate>RouterGroup.Post(
	"/",  // or "/:<entityId>"
	authMiddleware.Authenticate(),
	authMiddleware.RequireOrganizationPermission(rbac.<Verb><Entity>Permission),
	<aggregate>.New<Verb>Organization<Entity>Handler(handlers.<Verb><Entity>),
)
```

Add the route inside the appropriate `RegisterUser<Aggregate>Presentation` and/or `RegisterOrganization<Aggregate>Presentation` function. Most modules expose both user and organization-scoped routes — if the new command applies to both, register it in both (the handler is shared).

**Scope is decided by the route group, not the permission name:** routes under `/organizations/:organizationId/...` use `RequireOrganizationPermission` (it resolves the org from the `:organizationId` path param); routes under `/users/me/...` use `RequireUserPermission`. The same permission constant can legitimately appear in both scopes. An org-scoped route MUST carry `:organizationId` in its path. A route that needs authentication but no permission check uses `authMiddleware.Authenticate()` alone (no permission middleware). (Memory: `project_authn_authz_middleware_split`.)

If the RBAC permission constant doesn't exist in `go-packages/rbac` yet, that's a real cross-package change — flag it to the user; don't invent permission names.

## Critical rules (project memory)

1. **Validation is centralized via `cqrs.ValidationBehavior`, not called in the handler.** Keep the `Validate()` method in `validator.go` (single source of rules), but do NOT call `command.Validate()` in the handler body — the `cqrs.ValidationBehavior(...)` wrapper applied in `infrastructure.go` runs it (runtime `Validatable` check) before the inner handler / transaction. Do NOT delete the `Validate()` method to dedupe with HTTP-layer parsing. (Memory: `project_cqrs_validation_behavior`, `feedback_validation_layering`.)
2. **No repository injection across bounded contexts.** Cross-module dependencies use the other module's service interface. (Memory: `feedback_no_repository_injection_across_contexts`.)
3. **Use `apperror`, not `errors.New` / `fmt.Errorf`.** Domain errors: `apperror.Validation` / `apperror.NotFound`. App-level wraps: `apperror.Internal(...).WithCause(err)`. (Memory: `project_apperror_pattern`.)
4. **No `fmt` package.** Use `strings.Builder` for hot concat, `strconv` for number formatting, `errors.New` is fine but prefer apperror sentinels. (Memory: `feedback_no_fmt_package`.)
5. **`json.Marshal` / `json.Unmarshal` from `github.com/blocknextai/go-packages/json`** for JSON in handlers if any — NOT stdlib `encoding/json`, NOT a `utils` package. (Memory: `feedback_no_utils_package`.)
6. **`new(value)` over `&value`** for pointer literals. (Memory: `feedback_go_new_keyword`.)
7. **`strings.TrimSpace(s) == ""`** for empty checks, not `s == ""`. (Memory: `feedback_strings_trimspace_for_empty_check`.)
8. **`slog` for any logging**, never `log` or `fmt.Println`. (Memory: `feedback_use_slog_not_log`.)
9. **Initialisms in Go identifiers stay uppercase** — `userID` not `userId`, `commonHTTP` not `commonHttp`, `APIKeyID` not `ApiKeyId`. (Memory: `feedback_initialism_rename_autonomous`.)
10. **Initialism rule does NOT apply to struct tag contents** — Go field is `OrganizationID`, but the tag stays `uri:"organizationId"` (camelCase). Same for `json`, `query`, `db`, `schema` tags. (Memory: `feedback_no_initialism_in_struct_tags`.)
11. **Repos take the full entity for `Create` / `Update` / `Delete`**, never just the ID. The DELETE flow is: load → ownership guard → `entity.Delete()` (domain method) → `repo.Delete(deletedEntity)`.

## Final steps

After all files exist:

1. `gofmt -w` on every new + modified file.
2. `go build ./...` — must compile.
3. `golangci-lint run ./internal/<module>/...` — fix issues.
4. Verify the new domain sentinels (if any were added) are present in `internal/<module>/domain/<aggregate>/errors.go`.
5. If new app-level sentinels were added, they live in `internal/<module>/application/<aggregate>/errors.go`. Most modules already have this file; add it if missing, with the same pattern.
6. Confirm the route is reachable: scan `presentation/presentation.go` for the new path and the right RBAC permission.

## What NOT to do

- Do NOT add `// TODO` placeholders for fields the user did not specify; ask instead.
- Do NOT validate inside `entity.New(...)` AND in `command.Validate()` for the same rule — domain constructors enforce invariants, command validators enforce request shape. They are distinct layers; a small overlap is fine, full duplication is not.
- Do NOT call `repository.Create` outside `transactionManager.ExecuteInTransaction(...)`. All write paths go through the tx manager so cross-aggregate consistency is preserved.
- Do NOT add an HTTP handler that constructs a domain entity directly. Presentation builds a command, calls the handler, returns the response. The mapping `request → command → entity` is strictly one-way per layer.
- Do NOT skip wiring updates. A command that's not in `Handlers` and not on a route is dead code.
- Do NOT call `repository.Delete(txCtx, command.ID)` — repos take the full entity. Always go: GetByID → ownership guard → `entity.Delete()` → `repo.Delete(deletedEntity)`. (See Critical rule 11.)
- Do NOT name the local `userID` variable as `userId`, or import the http helper as `commonHttp`. The project initialism convention is `userID` and `commonHTTP`. (See Critical rule 9.)
- Do NOT "fix" the struct tag content from `uri:"organizationId"` → `uri:"organizationID"` while applying initialism renames — the tag string is consumed by Fiber/JSON which expect camelCase. Only the Go identifier changes. (See Critical rule 10.)
