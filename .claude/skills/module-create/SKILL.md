---
name: module-create
description: Use this skill when the user wants to add a new bounded context (module) to platform-api. Triggers on phrases like "create a new module for X", "scaffold a `<module>` bounded context", "let's add the `<module>` module", "I need a new module to handle Y". Generates the full module layout (`module.go`, `contract/`, `internal/domain`, `internal/application`, `internal/usecases`, `internal/postgres`, `internal/migrations`, `internal/http`), the first aggregate, the schema migration, the migration registries, and the wiring in `internal/bootstrap/platform_api.go` (instantiation) + `cmd/platform-api/main.go` (route registration). Delegates per-endpoint scaffolding to the `usecase-create` and `migration-create` skills.
---

> **Working directory:** all relative paths in this skill (`internal/...`, `cmd/...`, `go build ./...`) are relative to `apps/platform-api/` in the monorepo. Make targets run from the repo root.

# Module Create

Scaffolds a new bounded context under `internal/modules/<module>/` with the project's standard layout, creates the first aggregate's domain entity + repository interface + Postgres repository, an empty `usecases` service, empty route registration, the schema-init migration, exposes a `Module` type with `Dependencies`, instantiates it in `internal/bootstrap/platform_api.go`, and registers its routes in `cmd/platform-api/main.go`.

## When to use

- User says "create a new module for `<purpose>`" / "scaffold the `<module>` module" / "add a `<module>` bounded context"
- A new business concept needs its own boundary, distinct from existing modules. (If the concept fits inside an existing module, add a new aggregate there instead — this skill is overkill for that.)

## What this skill does NOT do

This skill produces a **minimal viable module**: one aggregate, no use cases, just the structure. It does **not** scaffold individual endpoints. After running this skill, run:
- `migration-create` for the first table (this skill creates the schema-init migration; data tables are separate)
- `usecase-create` for every endpoint (commands and queries)

If the user asks for a full module with endpoints in one go, delegate the per-endpoint work to those skills sequentially after this one is done.

## Required information

Confirm with the user (ask only what's missing):

1. **Module name** — lowercase, single word, plural. Examples: `apikeys`, `credentials`, `organizations`, `webhooks`. Avoid hyphens, underscores, and singular nouns. The module name is reused as:
   - Go package name
   - Postgres schema name
   - Directory name under `internal/modules/`
2. **Aggregate name** — the first aggregate root inside the module. Often matches the module name (e.g., `apikeys` module owns the `apikeys` aggregate). For modules with multiple aggregates (e.g., `organizations` owns `organizations`, `organizationusers`, `organizationsocials`), pick the primary one.
3. **Entity name** — singular Go type for the aggregate root (e.g., `APIKey`, `Credential`, `Organization`). PascalCase, acronyms uppercase.
4. **Entity fields** — domain attributes (besides the standard `BaseEntity` triplet of ID + timestamps). For each: name, Go type, whether it's nullable.
5. **Cross-module dependencies** — does this module need data or services from other modules? Repositories from other modules are NEVER injected; only the other module's **service interface**, imported from its `contract/` package, can come in via `Dependencies`. Common examples: `accountContract.UserService`, `quotaContract.PlanLimitService`.
6. **Public surface** — does this module expose a service to other modules? If yes, `internal/application/<aggregate>/service.go` defines it, `contract/` re-exports it, and the `Module` struct gets a public `<X>Service` field; if no, all fields stay private and `contract/` is omitted.
7. **HTTP routes** — does the module expose endpoints? (Almost always yes.) If so, what's the URL prefix (e.g., `/users/me/<entity>`, `/organizations/:organizationId/<entity>`, or both)?

## Directory layout

```
internal/modules/<module>/
├── module.go                          # Dependencies + Module + NewModule + Register — the only entry point for bootstrap
├── contract/                          # only when other modules consume this one: type aliases + interfaces, leaf package
│   └── <entity>.go                    # one file per aggregate, snake_case singular (linked_account.go): domain aliases first, then the service alias
└── internal/                          # sealed by the compiler; nothing outside the module can import it
    ├── domain/<aggregate>/
    │   ├── <entity>.go                # aggregate root: BaseEntity + New + Update + Delete + validateThenReturn
    │   ├── repository.go              # interface only
    │   └── errors.go                  # apperror.Validation / NotFound / Conflict sentinels
    ├── application/<aggregate>/
    │   ├── errors.go                  # apperror.Internal sentinels (wrapped repo failures)
    │   └── service.go                 # only when the module exposes a service to other modules
    ├── usecases/
    │   ├── services.go                # Services + ServiceDependencies + NewServices
    │   └── <aggregate>/service.go     # Service + NewService; use-case files added by usecase-create
    ├── postgres/<aggregate>/
    │   └── repository.go              # SQL implementation (database.BaseRepository + BuildQuery)
    ├── migrations/                    # <TS>_<module>_<name>.{up,down}.sql, generated by make migration-create
    └── http/
        ├── routes.go                  # RegisterRoutes + Register<Scope><Aggregate>Routes (initially empty bodies)
        └── <aggregate>/               # fiber handlers, added by usecase-create
```

Plus changes outside the module:
- `internal/bootstrap/platform_api.go` — import + field + instantiate.
- `cmd/platform-api/main.go` — `app.<Module>Module.Register(...)`.
- `cmd/platform-api-migration/main.go` — `modules` slice entry.
- `<repo-root>/docker/platform-api/migration.Dockerfile` — `COPY` line for `internal/modules/<module>/internal/migrations`.

## Naming conventions

| Thing                       | Convention                                  | Example                              |
| --------------------------- | ------------------------------------------- | ------------------------------------ |
| Module dir / package        | `<module>` lowercase plural                 | `apikeys`, `credentials`             |
| Postgres schema             | identical to module name                    | `apikeys`, `credentials`             |
| Aggregate dir / package     | `<aggregate>` lowercase plural              | `apikeys`, `organizationusers`       |
| Domain entity type          | `<Entity>` PascalCase singular              | `APIKey`, `Credential`               |
| Repository interface        | `<Entity>Repository`                        | `APIKeyRepository`                   |
| Repository impl type        | `<Entity>Repository` (same — different pkg) | `APIKeyRepository`                   |
| Repository ctor             | `New<Entity>Repository`                     | `NewAPIKeyRepository`                |
| Service interface           | `<Aggregate>Service`                        | `CredentialService`                  |
| Use-case service            | `Service` / `NewService` in `usecases/<aggregate>` | `apikeysUseCases.Service`     |
| Module dependencies struct  | `Dependencies` (always, no prefix)          | `apikeys.Dependencies`               |
| Module struct               | `Module` (always, no prefix)                | `apikeys.Module`                     |
| Migration filename prefix   | `<TS>_<module>_<change>`                    | `..._apikeys_init_schema`            |

Import aliases are verbose by convention: `<module>Domain<Aggregate>`, `<module>Application<Aggregate>`, `<module>Postgres<Aggregate>`, `<module>UseCases` (module-level `usecases`), `<aggregate>UseCases` (aggregate-level), `<module>HTTP`; cross-module always `<module>Contract "…/internal/modules/<module>/contract"`.

## Templates

### `internal/domain/<aggregate>/<entity>.go`

```go
package <aggregate>

import (
	"strings"
	"time"

	"github.com/blocknextai/go-packages/database"
	bnuuid "github.com/blocknextai/go-packages/uuid"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	"github.com/google/uuid"
)

type <Entity> struct {
	database.BaseEntity

	OwnerType commonDomain.OwnerType
	OwnerID   uuid.UUID
	Name      string
}

func New(
	ownerType commonDomain.OwnerType,
	ownerID uuid.UUID,
	name string,
) (*<Entity>, error) {
	now := time.Now().UTC()

	entity := &<Entity>{
		BaseEntity: database.BaseEntity{
			ID:        bnuuid.NewV7(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		OwnerType: ownerType,
		OwnerID:   ownerID,
		Name:      name,
	}

	return entity.validateThenReturn()
}

func (e *<Entity>) Update(name string) (*<Entity>, error) {
	e.Name = name
	e.UpdatedAt = time.Now().UTC()
	return e.validateThenReturn()
}

func (e *<Entity>) Delete() (*<Entity>, error) {
	now := time.Now().UTC()
	e.UpdatedAt = now
	e.DeletedAt = new(now)
	return e.validateThenReturn()
}

func (e *<Entity>) validateThenReturn() (*<Entity>, error) {
	if strings.TrimSpace(e.Name) == "" {
		return nil, ErrInvalid<Entity>Name
	}

	if !e.OwnerType.IsValid() {
		return nil, ErrInvalidOwnerType
	}

	if e.OwnerID == uuid.Nil {
		return nil, ErrInvalidOwnerID
	}

	return e, nil
}
```

Rules:
- **Embed `database.BaseEntity`** for ID + CreatedAt + UpdatedAt + DeletedAt. Don't redeclare these.
- **`bnuuid.NewV7()`** generates the ID — UUIDv7 is time-ordered. Don't use `uuid.New()` (v4) or Postgres-side defaults.
- **All timestamps are `time.Now().UTC()`**. The DB stays dumb; UTC is enforced in the app.
- **`validateThenReturn()` is the canonical pattern**: every state transition ends with `return e.validateThenReturn()`.
- **Soft delete via `DeletedAt`** — set to `new(time.Now().UTC())`.
- **Field order matters**: `id`, business fields, then `created_at`, `updated_at`, `deleted_at`. The `columns` const, the scan, `Create`, and the migration all mirror it.

### `internal/domain/<aggregate>/repository.go`

```go
package <aggregate>

import (
	"context"

	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	"github.com/google/uuid"
)

type <Entity>Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*<Entity>, error)
	GetAllByOwner(ctx context.Context, ownerType commonDomain.OwnerType, ownerID uuid.UUID, searchQuery string, offset int, limit int) ([]*<Entity>, int64, error)
	Create(ctx context.Context, entity *<Entity>) error
	Update(ctx context.Context, entity *<Entity>) error
	Delete(ctx context.Context, entity *<Entity>) error
}
```

Rules:
- **The interface lives in `domain/`, the implementation in `postgres/`.**
- **Methods are minimal at first.** Add `GetByName`, `GetAllByStatus`, etc. when the first use case needs them.
- **Pagination/search param order is fixed**: `(ctx, ...domainKeys, searchQuery string, offset int, limit int)`.
- **`Create` / `Update` / `Delete` take the full entity**, not just IDs.

### `internal/domain/<aggregate>/errors.go`

```go
package <aggregate>

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	Err<Entity>NotFound    = apperror.NotFound("<entity> not found")
	ErrInvalid<Entity>Name = apperror.Validation("invalid <entity> name")
	ErrInvalidOwnerType    = apperror.Validation("invalid owner type")
	ErrInvalidOwnerID      = apperror.Validation("invalid owner id")
)
```

Domain sentinels use `apperror.Validation` / `apperror.NotFound` / `apperror.Conflict`. Lowercase messages, no trailing punctuation, no `Code` field.

### `internal/application/<aggregate>/errors.go`

```go
package <aggregate>

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrFailedToCreate<Entity>  = apperror.Internal("failed to create <entity>")
	ErrFailedToUpdate<Entity>  = apperror.Internal("failed to update <entity>")
	ErrFailedToDelete<Entity>  = apperror.Internal("failed to delete <entity>")
	ErrFailedToGet<Entity>ByID = apperror.Internal("failed to get <entity> by id")
)
```

App-level sentinels are `apperror.Internal(...)`, used with `.WithCause(err)` from use cases when a repository or external call fails.

### `internal/application/<aggregate>/service.go` (only when the module exposes a service)

If no other module needs to call into this one, **skip this file and `contract/` entirely.** Otherwise:

```go
package <aggregate>

import (
	"context"

	<module>Domain<Aggregate> "github.com/blocknextai/platform-api/internal/modules/<module>/internal/domain/<aggregate>"
	"github.com/google/uuid"
)

type <Aggregate>Service interface {
	GetByID(ctx context.Context, id uuid.UUID) (*<module>Domain<Aggregate>.<Entity>, error)
}

type <aggregate>Service struct {
	<aggregate>Repository <module>Domain<Aggregate>.<Entity>Repository
}

func New<Aggregate>Service(
	<aggregate>Repository <module>Domain<Aggregate>.<Entity>Repository,
) <Aggregate>Service {
	return &<aggregate>Service{
		<aggregate>Repository: <aggregate>Repository,
	}
}

func (s *<aggregate>Service) GetByID(ctx context.Context, id uuid.UUID) (*<module>Domain<Aggregate>.<Entity>, error) {
	return s.<aggregate>Repository.GetByID(ctx, id)
}
```

Rules:
- **The interface is the public surface**, the struct is private. Only the constructor is exported.
- **Service methods are coarse-grained** — what other modules legitimately need, not repo passthroughs.
- **Domain types cross the boundary only as parameters / results of this interface**, and are re-exported from `contract/` alongside it.

### `contract/` (only when the module exposes a service)

```go
// contract/<entity>.go
package contract

import (
	<module>Application<Aggregate> "github.com/blocknextai/platform-api/internal/modules/<module>/internal/application/<aggregate>"
	<module>Domain<Aggregate> "github.com/blocknextai/platform-api/internal/modules/<module>/internal/domain/<aggregate>"
)

type <Entity> = <module>Domain<Aggregate>.<Entity>

type <Aggregate>Service = <module>Application<Aggregate>.<Aggregate>Service
```

Rules:
- **Type aliases and interfaces only** — no wiring, no logic, no constructors. `contract/` is a leaf: it imports only its own module's `internal/` packages.
- **Do not widen a contract** with sentinels, constants or helpers just because a consumer wants them; the consumer owns those locally. Re-export a constant only when it is part of the service interface's vocabulary (`ScopeMCPInvoke`).
- **Consumers import `<module>Contract` and nothing else** from this module. The goal is that the module can be extracted into its own service by swapping `contract/` for a client.

### `internal/postgres/<aggregate>/repository.go`

```go
package <aggregate>

import (
	"context"
	"database/sql"

	"github.com/blocknextai/go-packages/database"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	<module>Domain<Aggregate> "github.com/blocknextai/platform-api/internal/modules/<module>/internal/domain/<aggregate>"
	"github.com/google/uuid"
)

const (
	tableName = "<module>.<table>"
	columns   = "id, owner_type, owner_id, name, created_at, updated_at, deleted_at"
)

var (
	queryGetByID = database.BuildQuery(`
		SELECT `, columns, `
		FROM `, tableName, `
		WHERE
			id = $1
			AND deleted_at IS NULL
		LIMIT 1
	`)

	queryGetAllByOwner = database.BuildQuery(`
		SELECT `, columns, `, COUNT(*) OVER() AS total
		FROM `, tableName, `
		WHERE
			owner_type = $1
			AND owner_id = $2
			AND deleted_at IS NULL
			AND ($3 = '' OR name ILIKE '%' || $3 || '%')
		ORDER BY created_at DESC
		LIMIT $4 OFFSET $5
	`)

	queryCreate = database.BuildQuery(`
		INSERT INTO `, tableName, ` (`, columns, `)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`)

	queryUpdate = database.BuildQuery(`
		UPDATE `, tableName, `
		SET
			name = $2,
			updated_at = $3
		WHERE
			id = $1
			AND deleted_at IS NULL
	`)

	queryDelete = database.BuildQuery(`
		UPDATE `, tableName, `
		SET
			updated_at = $2,
			deleted_at = $3
		WHERE
			id = $1
			AND deleted_at IS NULL
	`)
)

type <Entity>Repository struct {
	database.BaseRepository
}

func New<Entity>Repository(db *sql.DB) <module>Domain<Aggregate>.<Entity>Repository {
	return &<Entity>Repository{
		BaseRepository: database.NewBaseRepository(db),
	}
}

func (r *<Entity>Repository) scan(row interface{ Scan(dest ...any) error }, extras ...any) (*<module>Domain<Aggregate>.<Entity>, error) {
	var e <module>Domain<Aggregate>.<Entity>
	var deletedAt sql.NullTime

	dest := []any{
		&e.ID,
		&e.OwnerType,
		&e.OwnerID,
		&e.Name,
		&e.CreatedAt,
		&e.UpdatedAt,
		&deletedAt,
	}
	dest = append(dest, extras...)

	if err := row.Scan(dest...); err != nil {
		return nil, err
	}

	e.DeletedAt = database.ScanNullTime(deletedAt)

	return &e, nil
}

func (r *<Entity>Repository) getOne(ctx context.Context, query string, args ...any) (*<module>Domain<Aggregate>.<Entity>, error) {
	return database.GetOne(ctx, r.Executor(ctx), query, r.scan, <module>Domain<Aggregate>.Err<Entity>NotFound, args...)
}

func (r *<Entity>Repository) getManyWithTotal(ctx context.Context, query string, args ...any) ([]*<module>Domain<Aggregate>.<Entity>, int64, error) {
	return database.GetManyWithTotal(ctx, r.Executor(ctx), query, r.scan, args...)
}

func (r *<Entity>Repository) exec(ctx context.Context, query string, args ...any) error {
	return database.Exec(ctx, r.Executor(ctx), query, args...)
}

func (r *<Entity>Repository) execWithRowCheck(ctx context.Context, query string, args ...any) error {
	return database.ExecWithRowCheck(ctx, r.Executor(ctx), query, <module>Domain<Aggregate>.Err<Entity>NotFound, args...)
}

func (r *<Entity>Repository) GetByID(ctx context.Context, id uuid.UUID) (*<module>Domain<Aggregate>.<Entity>, error) {
	return r.getOne(ctx, queryGetByID, id)
}

func (r *<Entity>Repository) GetAllByOwner(ctx context.Context, ownerType commonDomain.OwnerType, ownerID uuid.UUID, searchQuery string, offset int, limit int) ([]*<module>Domain<Aggregate>.<Entity>, int64, error) {
	return r.getManyWithTotal(ctx, queryGetAllByOwner, ownerType, ownerID, searchQuery, limit, offset)
}

func (r *<Entity>Repository) Create(ctx context.Context, entity *<module>Domain<Aggregate>.<Entity>) error {
	return r.exec(ctx, queryCreate,
		entity.ID,
		entity.OwnerType,
		entity.OwnerID,
		entity.Name,
		entity.CreatedAt,
		entity.UpdatedAt,
		entity.DeletedAt,
	)
}

func (r *<Entity>Repository) Update(ctx context.Context, entity *<module>Domain<Aggregate>.<Entity>) error {
	return r.execWithRowCheck(ctx, queryUpdate, entity.ID, entity.Name, entity.UpdatedAt)
}

func (r *<Entity>Repository) Delete(ctx context.Context, entity *<module>Domain<Aggregate>.<Entity>) error {
	return r.execWithRowCheck(ctx, queryDelete, entity.ID, entity.UpdatedAt, entity.DeletedAt)
}
```

Rules:
- **`database.BaseRepository`** provides `Executor(ctx)` (transaction-aware via context). Never hold or call `*sql.DB` directly.
- **`database.GetOne` / `GetManyWithTotal` / `Exec` / `ExecWithRowCheck`** are the query helpers; `GetOne` and `ExecWithRowCheck` translate no-row cases into the domain not-found sentinel you pass in.
- **`COUNT(*) OVER() AS total`** in list queries feeds `GetManyWithTotal` — one round trip for items + total. `scan` accepts `extras` so the total column can be appended to the destination list.
- **`scan` is duplicated per repository on purpose** — it references the entity's specific columns.
- **Unique-violation mapping**: `database.IsUniqueViolationOn(err, "<index>")` → domain `apperror.Conflict` sentinel, when the table has a unique index.
- **Soft delete is an UPDATE.** All queries filter `deleted_at IS NULL`.

### `internal/usecases/<aggregate>/service.go`

```go
package <aggregate>

import (
	"github.com/blocknextai/go-packages/database"
	<module>Domain<Aggregate> "github.com/blocknextai/platform-api/internal/modules/<module>/internal/domain/<aggregate>"
)

type Service struct {
	<aggregate>Repository <module>Domain<Aggregate>.<Entity>Repository
	transactionManager    database.TransactionManager
}

func NewService(
	<aggregate>Repository <module>Domain<Aggregate>.<Entity>Repository,
	transactionManager database.TransactionManager,
) *Service {
	return &Service{
		<aggregate>Repository: <aggregate>Repository,
		transactionManager:    transactionManager,
	}
}
```

Holds the union of dependencies every use case of the aggregate needs. `usecase-create` adds one `<usecase>.go` method per endpoint here and extends the fields when a use case needs more.

### `internal/usecases/services.go`

```go
package usecases

import (
	"github.com/blocknextai/go-packages/database"
	<module>Domain<Aggregate> "github.com/blocknextai/platform-api/internal/modules/<module>/internal/domain/<aggregate>"
	<aggregate>UseCases "github.com/blocknextai/platform-api/internal/modules/<module>/internal/usecases/<aggregate>"
)

type Services struct {
	<Aggregate> *<aggregate>UseCases.Service
}

type ServiceDependencies struct {
	TransactionManager database.TransactionManager

	<Aggregate>Repository <module>Domain<Aggregate>.<Entity>Repository
}

func NewServices(deps ServiceDependencies) *Services {
	return &Services{
		<Aggregate>: <aggregate>UseCases.NewService(deps.<Aggregate>Repository, deps.TransactionManager),
	}
}
```

`ServiceDependencies` lists the core infra dep (`TransactionManager`) first, then a blank line, then repositories and cross-module services. This file and `module.go` are the only places a module's wiring lives.

### `internal/http/routes.go`

```go
package http

import (
	commonAuth "github.com/blocknextai/platform-api/internal/common/auth"
	<module>UseCases "github.com/blocknextai/platform-api/internal/modules/<module>/internal/usecases"
	"github.com/gofiber/fiber/v3"
)

func RegisterOrganization<Aggregate>Routes(
	router fiber.Router,
	authMiddleware *commonAuth.AuthMiddleware,
	useCases *<module>UseCases.Services,
) {
	router.Group("/organizations/:organizationId/<aggregate>")
}

func RegisterRoutes(
	router fiber.Router,
	authMiddleware *commonAuth.AuthMiddleware,
	useCases *<module>UseCases.Services,
) {
	RegisterOrganization<Aggregate>Routes(router, authMiddleware, useCases)
}
```

Initially empty groups; `usecase-create` adds route lines (`authMiddleware.Authenticate()` → `RequireOrganizationPermission` / `RequireUserPermission` → `<aggregate>.New…Handler(useCases.<Aggregate>)`). Add `RegisterUser<Aggregate>Routes` (group `/users/me/<aggregate>`) when the module serves user-scoped routes. Assign the group to `<aggregate>RouterGroup` when the first route lands.

### `module.go`

```go
package <module>

import (
	"database/sql"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/platform-api/internal/common/auth"
	<module>HTTP "github.com/blocknextai/platform-api/internal/modules/<module>/internal/http"
	<module>Postgres<Aggregate> "github.com/blocknextai/platform-api/internal/modules/<module>/internal/postgres/<aggregate>"
	<module>UseCases "github.com/blocknextai/platform-api/internal/modules/<module>/internal/usecases"
	"github.com/gofiber/fiber/v3"
)

type Dependencies struct {
	DB                 *sql.DB
	TransactionManager database.TransactionManager
}

type Module struct {
	useCases *<module>UseCases.Services
}

func NewModule(deps Dependencies) *Module {
	repository := <module>Postgres<Aggregate>.New<Entity>Repository(deps.DB)
	useCases := <module>UseCases.NewServices(<module>UseCases.ServiceDependencies{
		TransactionManager: deps.TransactionManager,

		<Aggregate>Repository: repository,
	})
	return &Module{
		useCases: useCases,
	}
}

func (m *Module) Register(router fiber.Router, authMiddleware *auth.AuthMiddleware) {
	<module>HTTP.RegisterRoutes(router, authMiddleware, m.useCases)
}
```

When the module exposes a service, `NewModule` also builds it (`<module>Application<Aggregate>.New<Aggregate>Service(repository)`) and stores it in a public field `<Aggregate>Service <module>Application<Aggregate>.<Aggregate>Service` that bootstrap passes into consumers' `Dependencies`.

Rules:
- **`Dependencies` lists every external input.** `DB` and `TransactionManager` are the baseline; cross-module `<module>Contract` service interfaces, secret manager, options structs extend it. Cross-module deps go after a blank line.
- **`Module` fields are private by default.** A field becomes public only when another module's `NewModule` reads it. No getters, no speculative exposure.
- **No two-phase init.** Construct fully in `NewModule`. The only sanctioned exception is the existing `workflowsModule.SetMarketplaceWorkflowService(...)` read-enrichment cycle; do not add another.
- **`NewModule` returns `*Module`** (add an `error` return only when construction can genuinely fail).
- **`Register(router, authMiddleware, ...)`** may also take `cacheMiddleware` / `apiKeyMiddleware`; match a similar existing module.

### Schema-init migration

Use the `migration-create` skill (or run `make migration-create` directly):

```bash
make migration-create name=init_schema module=<module>
```

`.up.sql`:

```sql
CREATE SCHEMA IF NOT EXISTS <module>;
```

`.down.sql`:

```sql
DROP SCHEMA IF EXISTS <module> CASCADE;
```

This is the only migration this skill produces. The first table migration (`<TS>_<module>_add_<tables>_table`) is a separate `migration-create` invocation.

**A new module's migrations silently never run unless it is added to both registries:**
1. the `modules` slice in `cmd/platform-api-migration/main.go`
2. the `COPY --from=builder /build/internal/modules/<module>/internal/migrations /app/migrations/<module>/` line in `<repo-root>/docker/platform-api/migration.Dockerfile`

## Wiring

- **`internal/bootstrap/platform_api.go`** — instantiation inside `NewPlatformAPI`.
- **`cmd/platform-api/main.go`** — route registration off the `app` (`*bootstrap.PlatformAPI`).

### `internal/bootstrap/platform_api.go` — three edits

1. **Import** `"github.com/blocknextai/platform-api/internal/modules/<module>"` (alphabetical).
2. **Field** on the `PlatformAPI` struct: `<Module>Module *<module>.Module`, column-aligned with the existing fields.
3. **Instantiate** inside `NewPlatformAPI`, after every module whose public service it consumes and before any module that consumes it:

```go
<module>Module := <module>.NewModule(<module>.Dependencies{
	DB:                 core.DB,
	TransactionManager: core.TransactionManager,
})
```

Cross-module services come from the producer module's public field (`UserService: accountModule.UserService`). If `NewModule` returns an error, use the early-return form. Add `<Module>Module: <module>Module,` to the returned struct literal.

### `cmd/platform-api/main.go` — register routes

```go
app.<Module>Module.Register(fiberApp, authMiddleware)
```

Add `cacheMiddleware` / `apiKeyMiddleware` arguments when the module's `Register` takes them.

## Critical rules

1. **One module per turn.** Scaffold, verify it builds, then move on.
2. **Sealed `internal/`.** Nothing outside the module imports it; other contexts consume `contract/` only, which re-exports type aliases and interfaces. `contract/` stays a leaf.
3. **Private by default, public only when consumed.** No getters, no aliases in bootstrap.
4. **No two-phase init, no lazy suppliers.** Redesign the dependency direction instead.
5. **No repository injection across bounded contexts.** Cross-module dependency = `<module>Contract` service interface.
6. **The wiring chain is fixed:** `bootstrap` → `NewModule(Dependencies)` → `usecases.NewServices(ServiceDependencies)` → `http.RegisterRoutes(router, authMiddleware, useCases)`.
7. **`apperror` everywhere.** Domain: Validation / NotFound / Conflict. App: Internal + `.WithCause`. No `errors.New` / `fmt.Errorf` in business code.
8. **`bnuuid.NewV7()`** for IDs, **`time.Now().UTC()`** for timestamps, **`go-packages/json`** for JSON, **`slog`** for logging, **`strings.TrimSpace(s) == ""`** for empty checks, **`new(value)`** over `&value`, no `fmt`.
9. **Migrations register twice** (`cmd/platform-api-migration/main.go` + `<repo-root>/docker/platform-api/migration.Dockerfile`) or they never run.

## Final steps

1. `gofmt -w` on every new file + `internal/bootstrap/platform_api.go` + `cmd/platform-api/main.go`.
2. `go build ./...` — must compile. Common failures:
   - Missing import / struct field / instantiation / returned-literal entry in `platform_api.go`
   - Missing `app.<Module>Module.Register(...)` in `cmd/platform-api/main.go`
   - A `Dependencies` field referencing a contract service that doesn't exist yet
   - Repository scan column order not matching the domain struct / table
3. `golangci-lint run ./internal/modules/<module>/...` — must be clean.
4. Confirm the module appears in both migration registries.
5. (Optional, requires `GITHUB_TOKEN` + user confirmation) `make migration-up module=<module>`.
6. Hand off: "The module scaffold is ready. Run `migration-create` for the first table, then `usecase-create` for each endpoint."

## What NOT to do

- Do NOT add multiple aggregates in this skill. One aggregate per initial scaffold.
- Do NOT pre-populate `usecases/<aggregate>/` with use-case files or `http/<aggregate>/` with handlers. Those come from `usecase-create` and contain real code.
- Do NOT create `contract/` or `application/<aggregate>/service.go` unless another module is genuinely about to consume this one.
- Do NOT put anything but type aliases and interfaces in `contract/`.
- Do NOT inject `database.TransactionManager` into the application service — only the use-case `Service` needs it.
- Do NOT add `// TODO` placeholders. Either the scaffold compiles (empty route groups are fine) or it does not.
- Do NOT skip the wiring or the migration registries. An unwired module is dead code; an unregistered migration never runs.
- Do NOT run `make migration-up` autonomously.
- Do NOT introduce new kernel helpers (`internal/common/...`) as part of this skill.
