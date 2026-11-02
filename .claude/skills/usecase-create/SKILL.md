---
name: usecase-create
description: Use this skill when the user wants to add a new use case — a write-side command (create / update / delete / regenerate / archive …) or a read-side query (get by id / get all / list / search) — to a bounded context in platform-api. Triggers on phrases like "add a `<verb><entity>` command", "scaffold getAll for X", "create a delete endpoint for Y", "add a getById use case", "add a list endpoint". Generates the `<use_case>.go` method on the aggregate's `usecases/<aggregate>.Service` plus its `_command.go` / `_query.go`, `_validator.go`, `_mapper.go`, `_response.go` files, the Fiber handler under `http/<aggregate>/`, and the route in `http/routes.go`.
---

> **Working directory:** all relative paths in this skill (`internal/...`, `cmd/...`, `go build ./...`) are relative to `apps/platform-api/` in the monorepo. Make targets run from the repo root.

# Use Case Create

Adds one use case to an existing aggregate: a method on `internal/modules/<module>/internal/usecases/<aggregate>/service.go`'s `Service`, the sibling files that carry its input / output / validation / mapping, a Fiber handler in `internal/http/<aggregate>/<use_case>.go`, and a route line in `internal/http/routes.go`. Commands and queries share the same shape; the differences are listed per section.

## When to use

- User says "add a `<verb><entity>` command" / "let's add `delete<Entity>`" / "scaffold `getAll<Entities>`" / "add a by-id endpoint for X"
- The module and aggregate already exist. For a new bounded context run `module-create` first.

## Required information

Confirm with the user (ask only what's missing):

1. **Module** — bounded context under `internal/modules/` (e.g. `apikeys`, `credentials`, `organizations`).
2. **Aggregate** — package under `internal/usecases/` (e.g. `apikeys`, `organizationusers`). Plural, lowercase, single word.
3. **Kind** — **command** (write: create / update / delete / regenerate / transfer …, runs in a transaction) or **query** (read: by-id / list by owner / list with filter / singleton fetch, no transaction).
4. **Use case name** — `<Verb><Entity>` for commands (`CreateAPIKey`, `DeleteCredential`), `Get<Thing>` for queries (`GetAllAPIKeys`, `GetCredentialByID`).
5. **Input fields** — owner scope (`OwnerType` / `OwnerID`), entity id, payload, filters, search + pagination for lists. Their Go types.
6. **Response shape** — commands usually echo `ID` + a few fields (delete returns a response with just `ID`, or an empty struct — match the module); queries return a subset of the entity, never encrypted blobs or internal ids.
7. **Dependencies** — the repository (always) plus any cross-context service the method needs. Check `service.go` first: if the `Service` already holds it, nothing to wire.
8. **Auth & route** — HTTP verb, path, RBAC constant from `go-packages/rbac`, and scope: user (`/users/me/...`), organization (`/organizations/:organizationId/...`), or both.

## File layout

```
internal/modules/<module>/internal/usecases/<aggregate>/
├── service.go                  # Service struct + NewService — union of deps for every use case of the aggregate (exists already)
└── <use_case>.go               # ONE file per use case, snake_case of the method, top to bottom: input struct → Validate() → local sentinels → response struct(s) → func (s *Service) <UseCase>(...) → mappers
```

Plus two touch-points:
- `internal/modules/<module>/internal/http/<aggregate>/<use_case>.go` — Fiber handler, **snake_case filename**.
- `internal/modules/<module>/internal/http/routes.go` — route line in `RegisterUser<Aggregate>Routes` / `RegisterOrganization<Aggregate>Routes`.

Only when the method needs a dependency the `Service` does not hold yet: `service.go` (field + ctor param), `internal/usecases/services.go` (`ServiceDependencies` field + `NewService(...)` call), and `module.go` (`Dependencies` field + `ServiceDependencies` literal) — for a cross-module service that is a real wiring change, confirm with the user first.

## Naming rules

| Thing                  | Convention                                              | Example                                   |
| ---------------------- | ------------------------------------------------------- | ----------------------------------------- |
| Use-case file          | snake_case of the method                                | `create_api_key.go`, `get_all_api_keys.go` |
| Method                 | `<UseCase>` PascalCase, initialisms uppercase           | `CreateAPIKey`, `GetCredentialByID`       |
| Input struct           | `<UseCase>Command` / `<UseCase>Query`                   | `CreateAPIKeyCommand`, `GetAllAPIKeysQuery` |
| Response struct        | `<UseCase>Response`                                     | `CreateAPIKeyResponse`                    |
| List item response     | `<Entity>Response`                                      | `APIKeyResponse`                          |
| Mapper                 | `Map<Entity>ToResponse` / `Map<Entities>ToResponse`     | `MapAPIKeysToResponse`                    |
| HTTP file              | `<use>_<case>.go` snake_case                            | `create_api_key.go`, `get_all_api_keys.go` |
| HTTP handler ctor      | `New<Verb><Owner><Thing>Handler`                        | `NewCreateOrganizationAPIKeyHandler`      |
| HTTP request struct    | `<Verb><Owner><Thing>Request`                           | `GetAllOrganizationAPIKeysRequest`        |

All use cases of an aggregate share one package, so every type name must be unique inside it. When a use-case-local helper type would collide, prefix it with the use case name (`GetWorkflowByIDMarketplaceWorkflow`). Import aliases: aggregate package `<aggregate>UseCases "…/internal/usecases/<aggregate>"`, module package `<module>UseCases "…/internal/usecases"`, domain `<module>Domain<Aggregate> "…/internal/domain/<aggregate>"`.

## Templates

### Command struct — top of `<use_case>.go` (commands)

```go
package <aggregate>

import (
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	"github.com/google/uuid"
)

type <UseCase>Command struct {
	OwnerType commonDomain.OwnerType
	OwnerID   uuid.UUID
	Name      string
}
```

- No methods besides `Validate()`, no JSON tags — HTTP binding lives in the handler's request struct.
- Update / delete carry the entity id as `<Entity>ID uuid.UUID`.
- `OwnerType` + `OwnerID` let one method serve both user and organization routes.

### `Validate()` — right after the command struct (commands, when rules exist)

```go
package <aggregate>

import (
	"strings"

	<module>Domain<Aggregate> "github.com/blocknextai/platform-api/internal/modules/<module>/internal/domain/<aggregate>"
	"github.com/google/uuid"
)

const <UseCase>MaxNameLength = 255

func (c *<UseCase>Command) Validate() error {
	if !c.OwnerType.IsValid() {
		return <module>Domain<Aggregate>.ErrInvalidOwnerType
	}
	if c.OwnerID == uuid.Nil {
		return <module>Domain<Aggregate>.ErrInvalidOwnerID
	}
	if strings.TrimSpace(c.Name) == "" {
		return <module>Domain<Aggregate>.ErrInvalid<Entity>Name
	}
	if len(strings.TrimSpace(c.Name)) > <UseCase>MaxNameLength {
		return <module>Domain<Aggregate>.Err<Entity>NameTooLong
	}
	return nil
}
```

- Empty check is `strings.TrimSpace(s) == ""`; length check is `len(strings.TrimSpace(s)) > N`.
- Returns domain sentinels from `internal/domain/<aggregate>/errors.go` (`apperror.Validation(...)`); add missing ones there.
- The validator exists even when the domain constructor re-checks the same rule — that duplication is intentional. Queries have no validator.

### Query struct — top of `<use_case>.go` (queries)

```go
package <aggregate>

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	"github.com/google/uuid"
)

type Get<Thing>Query struct {
	OwnerType  commonDomain.OwnerType
	OwnerID    uuid.UUID
	Search     resultPkg.SearchRequest
	Pagination resultPkg.PaginationRequest
}
```

By-id queries replace `Search` / `Pagination` with `<Entity>ID uuid.UUID`. `Search` always precedes `Pagination`; `.Normalize()` is called in the HTTP handler, not here.

### `<use_case>.go`

**Command (transactional):**

```go
package <aggregate>

import (
	"context"

	<module>Application<Aggregate> "github.com/blocknextai/platform-api/internal/modules/<module>/internal/application/<aggregate>"
	<module>Domain<Aggregate> "github.com/blocknextai/platform-api/internal/modules/<module>/internal/domain/<aggregate>"
)

func (s *Service) <UseCase>(ctx context.Context, command *<UseCase>Command) (*<UseCase>Response, error) {
	if err := command.Validate(); err != nil {
		return nil, err
	}

	var entity *<module>Domain<Aggregate>.<Entity>

	err := s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		created, err := <module>Domain<Aggregate>.New(command.OwnerType, command.OwnerID, command.Name)
		if err != nil {
			return err
		}

		if err := s.<aggregate>Repository.Create(txCtx, created); err != nil {
			return <module>Application<Aggregate>.ErrFailedToCreate<Entity>.WithCause(err)
		}

		entity = created
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &<UseCase>Response{ID: entity.ID, Name: entity.Name}, nil
}
```

Update / delete body inside the transaction:

```go
entity, err := s.<aggregate>Repository.GetByID(txCtx, command.<Entity>ID)
if err != nil {
	return <module>Domain<Aggregate>.Err<Entity>NotFound
}
if entity.OwnerType != command.OwnerType || entity.OwnerID != command.OwnerID {
	return <module>Domain<Aggregate>.Err<Entity>NotFound
}

updated, err := entity.Update(command.Name)      // or: deleted, err := entity.Delete()
if err != nil {
	return err
}
if err := s.<aggregate>Repository.Update(txCtx, updated); err != nil {   // or Delete(txCtx, deleted)
	return <module>Application<Aggregate>.ErrFailedToUpdate<Entity>.WithCause(err)
}
```

**Query (no transaction):**

```go
func (s *Service) Get<Thing>(ctx context.Context, query *Get<Thing>Query) (*Get<Thing>Response, error) {
	entities, totalCount, err := s.<aggregate>Repository.GetAllByOwner(
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

By-id: fetch, ownership guard returning `Err<Entity>NotFound`, then `return Map<Entity>ToResponse(entity), nil`.

Rules:
- `command.Validate()` is the **first statement** whenever the command has a `Validate()` method — before any transaction. There is no dispatcher, pipeline or `Handler` type.
- Every write goes through `s.transactionManager.ExecuteInTransaction`; reads never do.
- Ownership guard returns NotFound, not Forbidden — do not leak existence.
- Repository / infra failures wrap with app sentinels from `internal/application/<aggregate>/errors.go` (`apperror.Internal(...)`, `.WithCause(err)`); domain sentinels propagate bare.
- Repos take the full entity for `Create` / `Update` / `Delete`; delete flow is load → guard → `entity.Delete()` → `repo.Delete(deleted)`.
- Repository list argument order is `(ctx, ...keys, searchQuery, offset, limit)`; pass `query.Search.Query` / `query.Pagination.Offset` / `query.Pagination.Limit` directly.
- Side effects that are not the mapping itself (decrypting, joining another service's data) happen in the method and are passed into the mapper as arguments.
- Private helpers shared by several use cases of the aggregate become unexported methods on the same `Service`.
- Dependencies come from `Service` fields only; never inject another module's repository — use its `<module>Contract` service interface.

### Mappers — bottom of `<use_case>.go` (queries)

```go
package <aggregate>

import (
	<module>Domain<Aggregate> "github.com/blocknextai/platform-api/internal/modules/<module>/internal/domain/<aggregate>"
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

- Pure: no `context.Context`, no DB or service calls. Extra data arrives as parameters (`Map<Entity>ToResponse(entity, processedData)`).
- Preallocate with `make(..., 0, len)` + `append`, never index assignment.
- Always write the mapper, even for one field — the method never builds a query response inline.

### Response struct(s) — before the Service method

```go
package <aggregate>

import (
	"time"

	"github.com/google/uuid"
)

type <Entity>Response struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Get<Thing>Response struct {
	Items      []*<Entity>Response
	TotalCount int64
}
```

Command responses and by-id responses are a single JSON-tagged struct. List wrappers (`Items` / `TotalCount`) carry no tags — the HTTP handler splits them into the paginated envelope. Nullable timestamps are `*time.Time` with `omitempty`.

## Wiring

### 1. `internal/http/<aggregate>/<use_case>.go`

**Command handler:**

```go
package <aggregate>

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	<aggregate>UseCases "github.com/blocknextai/platform-api/internal/modules/<module>/internal/usecases/<aggregate>"
)

type <Verb>Organization<Entity>Request struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	Name           string    `json:"name"`
}

func New<Verb>Organization<Entity>Handler(service *<aggregate>UseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(<Verb>Organization<Entity>Request)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := service.<UseCase>(c.RequestCtx(), &<aggregate>UseCases.<UseCase>Command{
			OwnerType: commonDomain.OwnerTypeOrganization,
			OwnerID:   request.OrganizationID,
			Name:      request.Name,
		})
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusCreated).JSON(resultPkg.Ok(result, resultPkg.WithMessage("<entity> created")))
	}
}
```

The user variant drops `OrganizationID`, uses `commonDomain.OwnerTypeUser` and `OwnerID: commonHTTP.GetUserID(c)`.

**List query handler:**

```go
type GetAllOrganization<Entities>Request struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	resultPkg.SearchRequest
	resultPkg.PaginationRequest
}

func NewGetAllOrganization<Entities>Handler(service *<aggregate>UseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(GetAllOrganization<Entities>Request)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		searchRequest := request.SearchRequest.Normalize()
		paginationRequest := request.PaginationRequest.Normalize()

		result, err := service.GetAll<Entities>(c.RequestCtx(), &<aggregate>UseCases.GetAll<Entities>Query{
			OwnerType:  commonDomain.OwnerTypeOrganization,
			OwnerID:    request.OrganizationID,
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

By-id handlers bind `<Entity>ID uuid.UUID \`uri:"<entity>Id"\`` and return `c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))`.

Rules:
- Status: `201 Created` for create, `200 OK` for everything else (match the module if it uses `204` for delete).
- `c.Bind().All(request)` is the default binding; `commonHTTP.GetUserID(c)` is the only way to read the authenticated user.
- Embed `SearchRequest` before `PaginationRequest`; call `.Normalize()` on both and pass the results — never the raw embedded fields. Lists respond via `commonHTTP.RespondPaginated`, not a hand-rolled envelope.
- Initialisms: Go identifiers are `OrganizationID` / `userID` / `commonHTTP`; tag contents stay camelCase (`uri:"organizationId"`).
- The handler builds the command / query and returns the result. It never touches domain types.

### 2. `internal/http/routes.go`

```go
// inside RegisterOrganization<Aggregate>Routes(router, authMiddleware, useCases *<module>UseCases.Services)
<aggregate>RouterGroup.Post(
	"/",
	authMiddleware.Authenticate(),
	authMiddleware.RequireOrganizationPermission(rbac.<Verb><Entity>Permission),
	<aggregate>.New<Verb>Organization<Entity>Handler(useCases.<Aggregate>),
)

// inside RegisterUser<Aggregate>Routes — group /users/me/<aggregate>
<aggregate>RouterGroup.Get(
	"/",
	authMiddleware.Authenticate(),
	authMiddleware.RequireUserPermission(rbac.Read<Entity>Permission),
	<aggregate>.NewGetAllUser<Entities>Handler(useCases.<Aggregate>),
)
```

- `authMiddleware.Authenticate()` first, then the scope-specific permission middleware. Scope follows the route group: `/organizations/:organizationId/...` uses `RequireOrganizationPermission`, `/users/me/...` uses `RequireUserPermission`. Authenticated-but-unpermissioned routes use `Authenticate()` alone.
- `useCases.<Aggregate>` is the `*<aggregate>UseCases.Service` field on `usecases.Services`; the same use case serves both scopes through two handler constructors.
- If the `rbac` constant does not exist, flag it — do not invent permission names.

## Critical rules

1. **`command.Validate()` first, before the transaction.** Keep `Validate()` as the single source of request rules even when the domain re-validates.
2. **Queries: mapper convention, no validator, no transaction.**
3. **Repository param order `(searchQuery, offset, limit)`**; struct field order `Search` then `Pagination`; HTTP embeds `SearchRequest` then `PaginationRequest`; `.Normalize()` is mandatory.
4. **No repository injection across bounded contexts** — `<module>Contract` service interfaces only.
5. **`apperror` everywhere**: domain `Validation` / `NotFound` / `Conflict`, app `Internal` + `.WithCause(err)`; no `errors.New` / `fmt.Errorf` in business code.
6. **No `fmt`**, `slog` for logging, `go-packages/json` for JSON, `new(value)` over `&value`, `append` over index assignment, `strings.TrimSpace(s) == ""` for empty checks.
7. **Initialisms uppercase in identifiers, camelCase inside struct tags.**

## Final steps

1. `gofmt -w` on every new and modified file.
2. `go build ./...` — must compile.
3. `golangci-lint run ./internal/modules/<module>/...` — fix issues.
4. Confirm any new sentinels landed in `internal/domain/<aggregate>/errors.go` (domain) or `internal/application/<aggregate>/errors.go` (app).
5. Confirm the route is reachable: scan `internal/http/routes.go` for the new path and RBAC permission.
6. For queries, re-check the response for sensitive fields.

## What NOT to do

- Do NOT introduce a `Handler` type, a `Handle` method, a dispatcher or a behavior pipeline — the fiber handler calls the `Service` method directly.
- Do NOT put the use case in its own package; it is a method in the aggregate's `usecases/<aggregate>` package.
- Do NOT skip `command.Validate()` or move the rules into the HTTP layer.
- Do NOT add `Validate()` or `transactionManager` usage to queries.
- Do NOT build query responses inline or pass `context.Context` to a mapper.
- Do NOT call `repository.Delete(txCtx, command.ID)` — repos take the full entity after `entity.Delete()`.
- Do NOT add `// TODO` placeholders for unspecified fields; ask instead.
- Do NOT leave the use case unrouted — a method with no handler and no route line is dead code.
