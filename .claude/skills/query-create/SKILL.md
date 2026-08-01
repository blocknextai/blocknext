---
name: query-create
description: Use this skill when the user wants to add a new read-side CQRS query (Get / GetAll / List / search) to a bounded context in platform-api. Triggers on phrases like "add a `get<entity>` query", "scaffold getAll for X", "create a list endpoint for Y", "add a getById handler". Generates query/handler/mapper/response files following the project's "GET endpoints use mapper convention", wires them through `infrastructure.go` and the presentation layer.
---

> **Working directory:** all relative paths in this skill (`internal/...`, `cmd/...`, `go build ./...`) are relative to `apps/platform-api/` in the monorepo. Make targets run from the repo root.


# Query Create

Scaffolds a new CQRS query under `internal/<module>/application/<aggregate>/get<thing>/`, wires it into the module's `infrastructure.go`, adds the Fiber HTTP handler in `presentation/<aggregate>/get_<thing>.go`, and registers the GET route in `presentation/presentation.go`.

## When to use

- User says "add a `get<entity>` query" / "scaffold getAll" / "list endpoint for X" / "by-id endpoint for Y"
- A read-side action is needed (no DB writes, no transaction).
- For write-side operations (create/update/delete) use the `command-create` skill.

## Required information

Confirm with the user (ask only what's missing):

1. **Module** — bounded context (e.g., `apikeys`, `credentials`, `organizations`).
2. **Aggregate** — sub-folder under `application/` (e.g., `apikeys`, `organizationusers`).
3. **Query type** — one of:
   - **By-ID lookup** — single entity by primary key (e.g., `getcredentialbyid`)
   - **List by owner** — paginated + searchable list scoped to an owner (e.g., `getallapikeys`)
   - **List with filter** — list with extra filter params beyond owner (e.g., `getallorganizationusers` filtered by role)
   - **Singleton fetch** — one entity by context (e.g., `getorganizationme` — current user's org membership)
4. **Folder name** — `get<thing>` lowercase no separator. Common patterns: `get<entity>byid`, `getall<entities>`, `get<thing>byname`.
5. **Query fields** — owner scope (`OwnerType` / `OwnerID`), entity ID for by-id, filters, pagination, search.
6. **Response shape** — fields exposed to the API. Usually a subset of the domain entity (no internal IDs, no encrypted data).
7. **Domain dependencies** — repository (always required), plus services if cross-context data must be joined (e.g., `secretManager.Decrypt`, `nodeEngineCredentialService.GetCredentialByID`).
8. **Auth & route** — URL path (always `GET`), RBAC permission constant from `go-packages/rbac`, user vs organization scope (or both).

## File layout

```
internal/<module>/application/<aggregate>/get<thing>/
├── query.go      # request struct (no methods, no Validate)
├── handler.go    # Handler struct + New + Handle (no tx)
├── mapper.go     # pure Map<Entity>(s)ToResponse functions — domain → response
└── response.go   # Response struct(s) with json tags
```

Plus three wiring touch-points (same as command-create):
- `internal/<module>/infrastructure/infrastructure.go` — add to `Handlers` struct + instantiate.
- `internal/<module>/presentation/<aggregate>/get_<thing>.go` — Fiber HTTP handler. **snake_case filename.**
- `internal/<module>/presentation/presentation.go` — register `GET` route under the right group with the right RBAC.

## Naming rules

| Thing                      | Convention                                    | Example                              |
| -------------------------- | --------------------------------------------- | ------------------------------------ |
| Folder / package           | `get<thing>` lowercase                        | `getallapikeys`, `getcredentialbyid` |
| Query struct               | `Get<Thing>Query`                             | `GetAllAPIKeysQuery`                 |
| Handler struct             | `Handler` (package name carries the context)  | `getallapikeys.Handler`             |
| Constructor                | `New`                                         | `getallapikeys.New`                 |
| Response struct            | `Get<Thing>Response`                          | `GetAllAPIKeysResponse`              |
| Item response (list)       | `<Entity>Response`                            | `APIKeyResponse`                     |
| Mapper (single)            | `Map<Entity>ToResponse`                       | `MapCredentialToResponse`            |
| Mapper (slice)             | `Map<Entities>ToResponse` (plural entity)     | `MapAPIKeysToResponse`               |
| Presentation file          | `get_<thing>.go` snake_case                   | `get_all_api_keys.go`                |
| Presentation handler ctor  | `NewGet<Owner><Thing>Handler`                 | `NewGetAllUserAPIKeysHandler`        |

## Templates

### `query.go`

**List query (paginated + searchable):**

```go
package get<thing>

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	"github.com/google/uuid"
)

type Get<Thing>Query struct {
	OwnerType  commonDomain.OwnerType
	OwnerID    uuid.UUID
	Search     resultPkg.SearchRequest      // ALWAYS before Pagination
	Pagination resultPkg.PaginationRequest
	// ... extra filters
}
```

**By-id query** — replace pagination/search with the entity ID:

```go
type Get<Entity>ByIDQuery struct {
	OwnerType  commonDomain.OwnerType
	OwnerID    uuid.UUID
	<Entity>ID uuid.UUID
}
```

Rules:
- **Field order is enforced**: `Search` field always appears before `Pagination`. Repository signatures and handler call sites mirror this. (Memory: `feedback_offset_limit_param_order`.)
- No `Validate()` method on queries — read paths fail at the repository layer if inputs are wrong, and ownership checks happen in the handler.
- No JSON tags — query is internal; HTTP binding is in the presentation request struct.
- `resultPkg.PaginationRequest` and `resultPkg.SearchRequest` are the only sanctioned types. Their `.Normalize()` methods normalize the values (offset clamp, limit defaults/max, query trim) — call `.Normalize()` in the presentation layer, not here.

### `handler.go`

**List query:**

```go
package get<thing>

import (
	"context"

	<module>Domain<Aggregate> "github.com/blocknextai/platform-api/internal/<module>/domain/<aggregate>"
)

type Handler struct {
	<aggregate>Repository <module>Domain<Aggregate>.<Entity>Repository
}

func New(
	<aggregate>Repository <module>Domain<Aggregate>.<Entity>Repository,
) *Handler {
	return &Handler{
		<aggregate>Repository: <aggregate>Repository,
	}
}

func (h *Handler) Handle(ctx context.Context, query *Get<Thing>Query) (*Get<Thing>Response, error) {
	entities, totalCount, err := h.<aggregate>Repository.GetAllByOwner(
		ctx,
		query.OwnerType,
		query.OwnerID,
		query.Search.Query,
		query.Pagination.Offset,
		query.Pagination.Limit,
	)
	if err != nil {
		return nil, err
	}

	return &Get<Thing>Response{
		Items:      Map<Entities>ToResponse(entities),
		TotalCount: totalCount,
	}, nil
}
```

**By-id query (with ownership guard):**

```go
func (h *Handler) Handle(ctx context.Context, query *Get<Entity>ByIDQuery) (*Get<Entity>ByIDResponse, error) {
	entity, err := h.<aggregate>Repository.GetByID(ctx, query.<Entity>ID)
	if err != nil {
		return nil, <module>Application<Aggregate>.ErrFailedToGet<Entity>ByID.WithCause(err)
	}

	if entity.OwnerType != query.OwnerType || entity.OwnerID != query.OwnerID {
		return nil, <module>Domain<Aggregate>.Err<Entity>NotFound
	}

	return Map<Entity>ToResponse(entity), nil
}
```

Rules:
- **Repository call argument order is fixed**: `(ctx, ...domainKeys, searchQuery, offset, limit)` — search FIRST, then offset, then limit. Mirrors the Query struct field order. (Memory: `feedback_offset_limit_param_order`.)
- **Pass `query.Search.Query` / `query.Pagination.Offset` / `query.Pagination.Limit` directly** — no intermediate variable extraction. The values are already normalized by the constructors that ran in the presentation layer.
- **No `transactionManager` on query handlers.** Reads don't need a transaction; the constructor takes only the repository and any required services. (Adding a tx wrapper to a read path is dead weight.)
- **Ownership guard returns NotFound, not Forbidden** — same as commands; do not leak existence.
- **Repository errors get wrapped only when the failure is "internal".** A real `not found` from the repo is the domain `Err<Entity>NotFound` and is returned bare. A network/SQL failure wraps with `apperror.WithCause(err)` via the app-level sentinel.
- **No repository injection across bounded contexts.** Cross-module reads use the other module's service interface. (Memory: `feedback_no_repository_injection_across_contexts`.)
- **Use `new(value)` not `&value`** for pointer construction (memory: `feedback_go_new_keyword`).
- **Decrypt / re-shape side-effects** that aren't a database hit (e.g., decrypting credential data, joining LLM-generated metadata) belong in the handler, not the mapper. The mapper stays pure.

### `mapper.go`

The mapper is the project's signature query convention (memory: `feedback_get_endpoints_use_mapper`). It exists for a reason: it isolates the domain → response shape so handler logic stays focused on data fetching.

**Slice mapper (list):**

```go
package get<thing>

import (
	<module>Domain<Aggregate> "github.com/blocknextai/platform-api/internal/<module>/domain/<aggregate>"
)

func Map<Entities>ToResponse(entities []*<module>Domain<Aggregate>.<Entity>) []*<Entity>Response {
	responses := make([]*<Entity>Response, 0, len(entities))
	for _, e := range entities {
		responses = append(responses, &<Entity>Response{
			ID:        e.ID,
			Name:      e.Name,
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		})
	}
	return responses
}
```

**Single mapper (by-id, often takes processed data alongside the entity):**

```go
func Map<Entity>ToResponse(entity *<module>Domain<Aggregate>.<Entity>, processedData map[string]any) *Get<Entity>ByIDResponse {
	return &Get<Entity>ByIDResponse{
		ID:        entity.ID,
		Name:      entity.Name,
		Data:      processedData,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}
```

Rules:
- **Mapper is a pure function** — no DB calls, no service calls, no `context.Context` parameter. If you need data from elsewhere, fetch it in the handler and pass it as an argument.
- **Slice mappers preallocate**: `make([]*<Entity>Response, 0, len(entities))` then `append`. Do NOT use `make([]*<Entity>Response, len(entities))` + `result[i] = ...`. (Memory: `feedback_use_append_not_index_assign`.)
- **Loop variable is short** (`e`, `k`, `c`) inside mappers — these are throwaway iterators, not aggregate names.
- **Do not inline mapping logic in the handler.** Even for a one-field response, write the mapper. The convention is the value, not the line count.

### `response.go`

```go
package get<thing>

import (
	"time"

	"github.com/google/uuid"
)

type <Entity>Response struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

type Get<Thing>Response struct {
	Items      []*<Entity>Response
	TotalCount int64
}
```

Rules:
- **Item type has json tags**, list wrapper does NOT — the wrapper is consumed by the presentation layer and split into `result.Items` + pagination envelope. JSON tags on the wrapper would be unused.
- **Use `*time.Time` + `omitempty`** for nullable timestamps (`LastUsedAt`, `DeletedAt`, etc.); `time.Time` for always-present ones.
- **For by-id responses** there's no list wrapper; the response IS the JSON-tagged struct returned directly.
- **Do NOT include sensitive fields** — encrypted blobs, internal foreign keys (e.g. raw `OwnerID` if not relevant to the consumer), debug fields. The mapper is the right place to enforce this filter.

## Wiring

### 1. `internal/<module>/infrastructure/infrastructure.go`

`RegisterInfrastructure` takes a single `RegisterInfrastructureDeps` struct (NOT positional params); query handlers read `deps.<Aggregate>Repository` and skip `deps.TransactionManager`.

The `Handlers` struct field holds the `cqrs.Handler[*Query, *Response]` interface (NOT the concrete `*Handler`), and the constructor is wrapped with `cqrs.ValidationBehavior(...)` — queries go through the same behavior pipeline as commands for uniformity (the runtime `Validatable` check is simply a no-op since queries have no `Validate()`). Import `"github.com/blocknextai/platform-api/internal/common/application/cqrs"`.

```go
type Handlers struct {
	GetAll<Entities> cqrs.Handler[*getall<entities>.GetAll<Entities>Query, *getall<entities>.GetAll<Entities>Response]  // ← new
	// ... existing
}

type RegisterInfrastructureDeps struct {
	TransactionManager database.TransactionManager

	<Aggregate>Repository <module>Domain<Aggregate>.<Entity>Repository
}

func RegisterInfrastructure(deps RegisterInfrastructureDeps) *Handlers {
	return &Handlers{
		GetAll<Entities>: cqrs.ValidationBehavior(getall<entities>.New(deps.<Aggregate>Repository)),  // ← new (no tx manager)
	}
}
```

`cqrs.ValidationBehavior` infers its type parameters from the concrete handler (Go 1.26) — no explicit type args needed. `cqrs` is a generic **behavior pipeline**: `ValidationBehavior` is one `cqrs.Behavior[C, R]` decorator; future behaviors (logging, metrics) compose by nesting or via `cqrs.Chain(handler, ...)` (see `command-create` for details). Note the asymmetry inside the constructor call: command handlers receive `deps.TransactionManager`, query handlers do not. If a query handler needs additional services (e.g., `secretManager`), add them as fields on `RegisterInfrastructureDeps` and on the module's `Dependencies` struct in `module.go`.

> **Unmigrated modules:** if the target module's `infrastructure.go` still stores concrete `*Get<Thing>QueryHandler` types (not yet migrated to `cqrs.Handler`), match the module's existing convention OR migrate the whole module to the `cqrs.ValidationBehavior` + `New`/`Handler` pattern first. Do not leave a module half-migrated. (Memory: `project_cqrs_validation_behavior`.)

### 2. `internal/<module>/presentation/<aggregate>/get_<thing>.go`

**List handler (with pagination + search):**

```go
package <aggregate>

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/<module>/application/<aggregate>/get<thing>"
)

type GetAllUser<Entities>Request struct {
	resultPkg.SearchRequest      // ALWAYS before PaginationRequest
	resultPkg.PaginationRequest
}

func NewGetAllUser<Entities>Handler(handler cqrs.Handler[*getall<entities>.GetAll<Entities>Query, *getall<entities>.GetAll<Entities>Response]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(GetAllUser<Entities>Request)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		searchRequest := request.SearchRequest.Normalize()
		paginationRequest := request.PaginationRequest.Normalize()

		result, err := handler.Handle(c.RequestCtx(), &getall<entities>.GetAll<Entities>Query{
			OwnerType:  commonDomain.OwnerTypeUser,
			OwnerID:    userID,
			Search:     searchRequest,
			Pagination: paginationRequest,
		})
		if err != nil {
			return err
		}

		return commonHTTP.RespondPaginated(c, result.Items, result.TotalCount, paginationRequest)
	}
}
```

**Organization variant** — same body, swap `OwnerTypeUser → OwnerTypeOrganization`, take `OrganizationID` from a `uri:"organizationId"` field:

```go
type GetAllOrganization<Entities>Request struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	resultPkg.SearchRequest
	resultPkg.PaginationRequest
}
```

**By-id handler (no pagination):**

```go
type GetUser<Entity>ByIDRequest struct {
	<Entity>ID uuid.UUID `uri:"<entity>Id"`
}

func NewGetUser<Entity>ByIDHandler(handler cqrs.Handler[*get<entity>byid.Get<Entity>ByIDQuery, *get<entity>byid.Get<Entity>ByIDResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		userID := commonHTTP.GetUserID(c)

		request := new(GetUser<Entity>ByIDRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := handler.Handle(c.RequestCtx(), &get<entity>byid.Get<Entity>ByIDQuery{
			OwnerType:  commonDomain.OwnerTypeUser,
			OwnerID:    userID,
			<Entity>ID: request.<Entity>ID,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
```

Rules:
- **Embedded field order is enforced**: `resultPkg.SearchRequest` ALWAYS embedded before `resultPkg.PaginationRequest`. Mirrors Query struct field order. (Memory: `feedback_offset_limit_param_order`.)
- **`.Normalize()` is mandatory**: `request.SearchRequest.Normalize()` and `request.PaginationRequest.Normalize()` MUST be called after `Bind().All()` and before passing to the query handler. `Normalize()` does all normalization (offset clamp, limit default/max, query trim) and returns a fresh `SearchRequest` / `PaginationRequest`. Bypassing it sends raw user input straight to the repository. (Memory: `feedback_pagination_search_normalization`.)
- **Pass the `Normalize()` results into the query**, not the embedded fields directly. Keep the normalized `paginationRequest` around — `commonHTTP.RespondPaginated` needs it.
- **Status code:** always `200 OK` on success (Fiber's default).
- **Response envelope:** lists use the `commonHTTP.RespondPaginated(c, result.Items, result.TotalCount, paginationRequest)` helper (it builds the pagination metadata + envelope internally — do NOT hand-roll `NewPagination` / `WithPagination`). By-id uses `resultPkg.Ok(result)` (single arg, no message — commands use `resultPkg.WithMessage(...)`).
- **List endpoints pass `result.Items`** to `RespondPaginated` so pagination is applied to items, not the wrapper struct.
- **Path param naming:** `<entity>Id` in `uri:"..."` matches the URL path (`:apiKeyId`, `:credentialId`, `:organizationId`). Consistent with how routes are registered.
- **`c.Bind().All(request)`** binds path + query + body. For pure GETs there's no body, but `.All()` is the project default.
- **Don't apply Initialism rules to struct tag contents**: the Go identifier is `OrganizationID`, but the URI tag stays `uri:"organizationId"` (camelCase). Same rule for `query:"..."`. (Memory: `feedback_no_initialism_in_struct_tags`.)

### 3. `internal/<module>/presentation/presentation.go`

Authentication and authorization are SEPARATE composable middlewares. Chain `authMiddleware.Authenticate()` first, then the scope-specific permission middleware:

```go
// user-scoped route (RegisterUser<Aggregate>Presentation, group /users/me/...)
<aggregate>RouterGroup.Get(
	"/",  // or "/:<entity>Id" for by-id
	authMiddleware.Authenticate(),
	authMiddleware.RequireUserPermission(rbac.Read<Entity>Permission),
	<aggregate>.NewGetAllUser<Entities>Handler(handlers.GetAll<Entities>),
)

// organization-scoped route (RegisterOrganization<Aggregate>Presentation, group /organizations/:organizationId/...)
<aggregate>RouterGroup.Get(
	"/",  // or "/:<entity>Id" for by-id
	authMiddleware.Authenticate(),
	authMiddleware.RequireOrganizationPermission(rbac.Read<Entity>Permission),
	<aggregate>.NewGetAllOrganization<Entities>Handler(handlers.GetAll<Entities>),
)
```

Add the route inside the appropriate `RegisterUser<Aggregate>Presentation` and/or `RegisterOrganization<Aggregate>Presentation` function. Most queries are exposed under both — the handler is shared, only the presentation handler constructor differs.

**Scope is decided by the route group, not the permission name:** routes under `/organizations/:organizationId/...` use `RequireOrganizationPermission` (it resolves the org from the `:organizationId` path param); routes under `/users/me/...` use `RequireUserPermission`. The same permission constant can legitimately appear in both scopes. An org-scoped route MUST carry `:organizationId` in its path. (Memory: `project_authn_authz_middleware_split`.)

If the RBAC permission constant doesn't exist (`Read<Entity>Permission`), check `go-packages/rbac` first; do not invent permission names.

## Critical rules (project memory)

1. **GET endpoints use the mapper convention** — separate `mapper.go` with pure functions. No inline construction in `handler.go`. (Memory: `feedback_get_endpoints_use_mapper`.)
2. **Repository pagination/search param order is `(searchQuery, offset, limit)`** — search FIRST. Query struct fields are `Search` then `Pagination`. Presentation request struct embeds `SearchRequest` then `PaginationRequest`. (Memory: `feedback_offset_limit_param_order`.)
3. **`.Normalize()` is required in the presentation layer** — embedding the request types is for binding only; `request.PaginationRequest.Normalize()` / `request.SearchRequest.Normalize()` perform the normalization and return fresh values. Never pass the un-normalized embedded `request.PaginationRequest` / `request.SearchRequest` directly to the query. (Memory: `feedback_pagination_search_normalization`.)
4. **Use `append`, not index assignment** when building slices in mappers. (Memory: `feedback_use_append_not_index_assign`.)
5. **No transaction on read paths.** Query handlers don't take `transactionManager`.
6. **No repository injection across bounded contexts.** Cross-module reads go through service interfaces. (Memory: `feedback_no_repository_injection_across_contexts`.)
7. **`apperror` for all errors.** Domain not-found is `apperror.NotFound(...)`. App-level wraps use `apperror.Internal(...).WithCause(err)`. (Memory: `project_apperror_pattern`.)
8. **Don't use `fmt`.** Use `strings.Builder`, `strconv`, `slog`. (Memory: `feedback_no_fmt_package`.)
9. **`json.Marshal` / `json.Unmarshal` from `github.com/blocknextai/go-packages/json`** for JSON — NOT stdlib `encoding/json`, NOT a `utils` package. (Memory: `feedback_no_utils_package`.)
10. **`new(value)` not `&value`** for pointer literals. (Memory: `feedback_go_new_keyword`.)
11. **`slog` for logging**, never `log` or `fmt.Println`. (Memory: `feedback_use_slog_not_log`.)

## Final steps

After all files exist:

1. `gofmt -w` on every new + modified file.
2. `go build ./...` — must compile.
3. `golangci-lint run ./internal/<module>/...` — fix issues.
4. Verify the new app-level sentinels (`ErrFailedToGet<Entity>ByID`, etc.) are present in `internal/<module>/application/<aggregate>/errors.go`. Add them if missing.
5. Confirm the route is reachable: scan `presentation/presentation.go` for the new `Get` call and the right RBAC permission.
6. If the response includes new fields, double-check they are not sensitive (encrypted blobs, internal IDs).

## What NOT to do

- Do NOT write `result[i] = MapXToResponse(items[i])` in mappers — use `append` (memory rule).
- Do NOT inline the mapping logic in the handler. The mapper file exists; use it even for trivial responses.
- Do NOT add `Validate()` to queries. Validation is for write paths.
- Do NOT add a `transactionManager` to a query handler. Reads don't need transactions in this project.
- Do NOT return `*<module>Domain<Aggregate>.<Entity>` from a handler. Always go through the mapper to a `*<Thing>Response`. The domain entity stays inside the application boundary.
- Do NOT skip the ownership guard on by-id queries when the entity has an owner. Returning someone else's entity is a security bug, not just a bug.
- Do NOT pass `context.Context` to mapper functions. Mappers are pure.
- Do NOT order params as `(offset, limit, searchQuery)` or `(limit, offset, ...)`. The standard is `(searchQuery, offset, limit)` — see Critical rules #2.
- Do NOT pass the un-normalized embedded `request.PaginationRequest` / `request.SearchRequest` directly to the query handler. `.Normalize()` MUST be called on each for normalization — see Critical rules #3.
- Do NOT hand-roll the list envelope with `resultPkg.NewPagination` / `resultPkg.WithPagination` / `resultPkg.Ok(items, "", ...)`. Lists return `commonHTTP.RespondPaginated(c, result.Items, result.TotalCount, paginationRequest)`.
- Do NOT use `paginationRequest.GetOffset()` / `.GetLimit()` — those methods don't exist on `resultPkg.PaginationRequest`. Access fields directly: `.Offset` / `.Limit`.
- Do NOT embed `resultPkg.PaginationRequest` before `resultPkg.SearchRequest` in the presentation request struct. The embedded order is `SearchRequest` first.
