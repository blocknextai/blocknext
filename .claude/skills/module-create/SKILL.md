---
name: module-create
description: Use this skill when the user wants to add a new bounded context (module) to platform-api. Triggers on phrases like "create a new module for X", "scaffold a `<module>` bounded context", "let's add the `<module>` module", "I need a new module to handle Y". Generates the full DDD layout (`application/`, `domain/`, `infrastructure/`, `presentation/`, `module.go`), the first aggregate, the schema migration, and the wiring in `internal/bootstrap/platform_api.go` (instantiation) + `cmd/platform-api/main.go` (route registration). Delegates per-endpoint scaffolding to the `command-create`, `query-create`, and `migration-create` skills.
---

> **Working directory:** all relative paths in this skill (`internal/...`, `cmd/...`, `go build ./...`) are relative to `apps/platform-api/` in the monorepo. Make targets run from the repo root.


# Module Create

Scaffolds a new bounded context under `internal/<module>/` with the project's standard DDD layout, creates the first aggregate's domain entity + repository interface + SQL repository, generates the schema-init migration, exposes a `Module` type with `Dependencies`, instantiates it in `internal/bootstrap/platform_api.go`, and registers its routes in `cmd/platform-api/main.go`.

## When to use

- User says "create a new module for `<purpose>`" / "scaffold the `<module>` module" / "add a `<module>` bounded context"
- A new business concept needs its own boundary, distinct from existing modules. (If the concept fits inside an existing module, add a new aggregate there instead — this skill is overkill for that.)

## What this skill does NOT do

This skill produces a **minimal viable module**: one aggregate, no commands or queries, just the structure. It does **not** scaffold individual endpoints. After running this skill, run:
- `migration-create` for the first table (the skill creates the schema-init migration; data tables are separate)
- `command-create` for write-side endpoints (create, update, delete)
- `query-create` for read-side endpoints (get, list)

If the user asks for a full module with endpoints in one go, delegate the per-endpoint work to those skills sequentially after this one is done.

## Required information

Confirm with the user (ask only what's missing):

1. **Module name** — lowercase, single word, plural. Examples: `apikeys`, `credentials`, `organizations`, `webhooks`. Avoid hyphens, underscores, and singular nouns. The module name is reused as:
   - Go package name
   - Postgres schema name
   - Directory name under `internal/`
2. **Aggregate name** — the first aggregate root inside the module. Often matches the module name (e.g., `apikeys` module owns the `apikeys` aggregate). For modules with multiple aggregates (e.g., `organizations` owns `organizations`, `organizationusers`, `organizationsocials`), pick the primary one.
3. **Entity name** — singular Go type for the aggregate root (e.g., `APIKey`, `Credential`, `Organization`). PascalCase, acronyms uppercase.
4. **Entity fields** — domain attributes (besides the standard `BaseEntity` triplet of ID + timestamps). For each: name, Go type, whether it's nullable.
5. **Cross-module dependencies** — does this module need data or services from other modules? Repositories from other modules are NEVER injected (memory: `feedback_no_repository_injection_across_contexts`); only the other module's **service interface** can come in via `Dependencies`. Common examples: `account.UserService`, `organizations.OrganizationService`.
6. **Public surface** — does this module expose a service to other modules? If yes, the `Module` struct gets a public `<X>Service` field; if no, all fields stay private. (Memory: `feedback_module_field_pattern`.)
7. **HTTP routes** — does the module expose endpoints? (Almost always yes.) If so, what's the URL prefix (e.g., `/users/me/<entity>`, `/organizations/:organizationId/<entity>`, or both)?

## Directory layout

```
internal/<module>/
├── application/
│   ├── <aggregate>/
│   │   ├── service.go          # only if module exposes a service interface to other modules
│   │   └── errors.go           # app-level apperror.Internal sentinels (wrapped repo failures)
├── domain/
│   ├── <aggregate>/
│   │   ├── <entity>.go         # aggregate root: BaseEntity + New + Update + Delete + validateThenReturn
│   │   ├── repository.go       # interface only — no implementation
│   │   └── errors.go           # apperror.Validation / apperror.NotFound sentinels
├── infrastructure/
│   ├── database/
│   │   └── migrations/         # populated by migration-create skill (NOT this skill, beyond the schema-init)
│   ├── <aggregate>/
│   │   └── repository.go       # SQL implementation (database.BaseRepository + BuildQuery)
│   └── infrastructure.go       # Handlers struct + RegisterInfrastructure (initially empty)
├── presentation/
│   ├── <aggregate>/            # populated by command-create / query-create (initially empty)
│   └── presentation.go         # RegisterPresentation (initially empty body)
└── module.go                   # Module + Dependencies + NewModule + Register
```

Plus two changes outside the module:
- `cmd/platform-api/main.go` — import + instantiate + register routes.
- `internal/<module>/infrastructure/database/migrations/<TS>_<module>_init_schema.up.sql` — generated via `make migration-create`.

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
| Module dependencies struct  | `Dependencies` (always, no prefix)          | `apikeys.Dependencies`               |
| Module struct               | `Module` (always, no prefix)                | `apikeys.Module`                     |
| Migration filename prefix   | `<TS>_<module>_<change>`                    | `..._apikeys_init_schema`            |

## Templates

### `domain/<aggregate>/<entity>.go`

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
	// ... other domain fields
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
- **Embed `database.BaseEntity`** (from `github.com/blocknextai/go-packages/database`) for ID + CreatedAt + UpdatedAt + DeletedAt. Don't redeclare these.
- **`bnuuid.NewV7()`** (alias of `github.com/blocknextai/go-packages/uuid`) generates the ID — UUIDv7 is time-ordered, which the project standardises on. Don't use `uuid.New()` (v4) or Postgres-side defaults. (`github.com/google/uuid` is still imported for the `uuid.UUID` type and `uuid.Nil`.)
- **All timestamps are `time.Now().UTC()`**. The migrations enforce UTC at the DB layer; the domain enforces it at the application layer.
- **`validateThenReturn()` is the canonical pattern**: every state transition (`New`, `Update`, `Delete`, etc.) ends with `return e.validateThenReturn()`. The function checks every invariant and returns the entity on success. This is how the domain enforces invariants — do NOT scatter validation across constructors.
- **Soft delete via `DeletedAt`** — set to `new(time.Now().UTC())`, never `nil` reset. Deletion is reversible at the data layer; the application chooses to ignore soft-deleted rows in queries.
- **`new(value)` over `&value`** for pointer literals (memory: `feedback_go_new_keyword`).
- **`strings.TrimSpace(s) == ""`** for empty-string checks (memory: `feedback_strings_trimspace_for_empty_check`).

### `domain/<aggregate>/repository.go`

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
- **Repository lives in `domain/`**, the interface is the contract; the implementation lives in `infrastructure/`. This is the dependency inversion that lets the application depend on domain only.
- **Methods are minimal at first.** Add new methods (`GetByName`, `GetAllByStatus`, etc.) when the first command/query needs them, not preemptively.
- **`GetAllByOwner` returns `(items, total, error)`** — pagination is two-call style (data + count) because most lists need totals for the UI.
- **Pagination/search param order is fixed**: `(ctx, ...domainKeys, searchQuery string, offset int, limit int)` — search FIRST. Always include `searchQuery string` even if the first version of the SQL ignores it; future search support drops in without a signature break. (Memory: `feedback_offset_limit_param_order`.)
- **`Create` / `Update` / `Delete` take the full entity**, not just IDs. The repo trusts the entity is already validated (constructors enforce that).

### `domain/<aggregate>/errors.go`

```go
package <aggregate>

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	Err<Entity>NotFound      = apperror.NotFound("<entity> not found")
	ErrInvalid<Entity>Name   = apperror.Validation("invalid <entity> name")
	ErrInvalidOwnerType      = apperror.Validation("invalid owner type")
	ErrInvalidOwnerID        = apperror.Validation("invalid owner id")
)
```

Rules:
- **Domain sentinels use `apperror.Validation` or `apperror.NotFound`** — these map to the right HTTP status codes via the presentation layer's apperror handler.
- **Lowercase messages**, no trailing punctuation.
- **No `Code` field** — the project's apperror pattern is Kind + Message + Cause only (memory: `project_apperror_pattern`).

### `application/<aggregate>/errors.go`

```go
package <aggregate>

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrFailedToCreate<Entity>     = apperror.Internal("failed to create <entity>")
	ErrFailedToUpdate<Entity>     = apperror.Internal("failed to update <entity>")
	ErrFailedToDelete<Entity>     = apperror.Internal("failed to delete <entity>")
	ErrFailedToGet<Entity>ByID    = apperror.Internal("failed to get <entity> by id")
)
```

Rules:
- **App-level sentinels are `apperror.Internal(...)`** — wrappers for repository / external-service failures. The application calls them with `.WithCause(err)` to attach the underlying error.
- **Distinct from domain sentinels.** Domain says "this is invalid" (Validation) or "this doesn't exist" (NotFound); app says "the infra layer broke" (Internal).

### `application/<aggregate>/service.go` (only when module exposes a service)

If no other module needs to call into this one, **skip this file entirely.** Otherwise:

```go
package <aggregate>

import (
	"context"

	"github.com/blocknextai/platform-api/internal/<module>/domain/<aggregate>"
	"github.com/google/uuid"
)

type <Aggregate>Service interface {
	GetByID(ctx context.Context, id uuid.UUID) (*<aggregate>.<Entity>, error)
}

type <aggregate>Service struct {
	<aggregate>Repository <aggregate>.<Entity>Repository
}

func New<Aggregate>Service(
	<aggregate>Repository <aggregate>.<Entity>Repository,
) <Aggregate>Service {
	return &<aggregate>Service{
		<aggregate>Repository: <aggregate>Repository,
	}
}

func (s *<aggregate>Service) GetByID(ctx context.Context, id uuid.UUID) (*<aggregate>.<Entity>, error) {
	return s.<aggregate>Repository.GetByID(ctx, id)
}
```

Rules:
- **The interface is the public surface**, the struct is private (lowercase). Only the constructor is public.
- **Service methods are coarse-grained.** They expose what other modules legitimately need, not arbitrary repo passthroughs.
- **Returning a domain entity across module boundaries is fine** — domain types are part of the contract. Returning a repository is not (memory: `feedback_no_repository_injection_across_contexts`).

### `infrastructure/<aggregate>/repository.go`

```go
package <aggregate>

import (
	"context"
	"database/sql"
	"errors"

	"github.com/blocknextai/go-packages/database"
	<module>Domain<Aggregate> "github.com/blocknextai/platform-api/internal/<module>/domain/<aggregate>"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
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
		SELECT `, columns, `
		FROM `, tableName, `
		WHERE
			owner_type = $1
			AND owner_id = $2
			AND ($3 = '' OR name ILIKE '%' || $3 || '%')
			AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $4 OFFSET $5
	`)

	queryCountByOwner = database.BuildQuery(`
		SELECT COUNT(*)
		FROM `, tableName, `
		WHERE
			owner_type = $1
			AND owner_id = $2
			AND ($3 = '' OR name ILIKE '%' || $3 || '%')
			AND deleted_at IS NULL
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

func (r *<Entity>Repository) scan(row interface{ Scan(dest ...any) error }) (*<module>Domain<Aggregate>.<Entity>, error) {
	var e <module>Domain<Aggregate>.<Entity>
	var deletedAt sql.NullTime

	err := row.Scan(
		&e.ID,
		&e.OwnerType,
		&e.OwnerID,
		&e.Name,
		&e.CreatedAt,
		&e.UpdatedAt,
		&deletedAt,
	)
	if err != nil {
		return nil, err
	}

	if deletedAt.Valid {
		e.DeletedAt = new(deletedAt.Time)
	}

	return &e, nil
}

func (r *<Entity>Repository) getOne(ctx context.Context, query string, args ...any) (*<module>Domain<Aggregate>.<Entity>, error) {
	row := r.GetExecutor(ctx).QueryRowContext(ctx, query, args...)
	e, err := r.scan(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, <module>Domain<Aggregate>.Err<Entity>NotFound
		}
		return nil, err
	}
	return e, nil
}

func (r *<Entity>Repository) getMany(ctx context.Context, query string, args ...any) ([]*<module>Domain<Aggregate>.<Entity>, error) {
	rows, err := r.GetExecutor(ctx).QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*<module>Domain<Aggregate>.<Entity>
	for rows.Next() {
		e, err := r.scan(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, e)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *<Entity>Repository) execWithRowCheck(ctx context.Context, query string, args ...any) error {
	result, err := r.GetExecutor(ctx).ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return <module>Domain<Aggregate>.Err<Entity>NotFound
	}

	return nil
}

func (r *<Entity>Repository) GetByID(ctx context.Context, id uuid.UUID) (*<module>Domain<Aggregate>.<Entity>, error) {
	return r.getOne(ctx, queryGetByID, id)
}

func (r *<Entity>Repository) GetAllByOwner(ctx context.Context, ownerType commonDomain.OwnerType, ownerID uuid.UUID, searchQuery string, offset int, limit int) ([]*<module>Domain<Aggregate>.<Entity>, int64, error) {
	var total int64
	err := r.GetExecutor(ctx).QueryRowContext(ctx, queryCountByOwner, ownerType, ownerID, searchQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []*<module>Domain<Aggregate>.<Entity>{}, 0, nil
	}

	result, err := r.getMany(ctx, queryGetAllByOwner, ownerType, ownerID, searchQuery, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return result, total, nil
}

func (r *<Entity>Repository) Create(ctx context.Context, entity *<module>Domain<Aggregate>.<Entity>) error {
	_, err := r.GetExecutor(ctx).ExecContext(ctx, queryCreate,
		entity.ID,
		entity.OwnerType,
		entity.OwnerID,
		entity.Name,
		entity.CreatedAt,
		entity.UpdatedAt,
		entity.DeletedAt,
	)
	return err
}

func (r *<Entity>Repository) Update(ctx context.Context, entity *<module>Domain<Aggregate>.<Entity>) error {
	return r.execWithRowCheck(ctx, queryUpdate, entity.ID, entity.Name, entity.UpdatedAt)
}

func (r *<Entity>Repository) Delete(ctx context.Context, entity *<module>Domain<Aggregate>.<Entity>) error {
	return r.execWithRowCheck(ctx, queryDelete, entity.ID, entity.UpdatedAt, entity.DeletedAt)
}
```

Rules:
- **`database.BaseRepository`** provides `GetExecutor(ctx)` which returns the right executor (transaction-aware via context propagation). NEVER call `r.db.QueryContext(...)` directly.
- **`database.BuildQuery(...)`** is the project's query helper. Its variadic-string form composes table names + columns + literal SQL fragments without runtime concatenation.
- **`scan` / `getOne` / `getMany` / `execWithRowCheck`** are duplicated per repository on purpose — they reference the entity's specific columns. Don't try to abstract them into common.
- **`sql.ErrNoRows` → domain not-found sentinel.** The repo translates the SQL-layer signal into a domain error so the application layer doesn't need to import `database/sql`.
- **`execWithRowCheck` returns `Err<Entity>NotFound` on zero rows affected** — UPDATE / DELETE on a missing or already-deleted row is treated as not-found, not as silent success.
- **Soft delete is an UPDATE.** `Delete` sets `deleted_at`, it doesn't `DELETE FROM`. All queries filter `deleted_at IS NULL`.

### `infrastructure/infrastructure.go`

```go
package infrastructure

import (
	"github.com/blocknextai/go-packages/database"
	<module>Domain<Aggregate> "github.com/blocknextai/platform-api/internal/<module>/domain/<aggregate>"
)

type Handlers struct {
	// Populated by command-create / query-create skills as endpoints are added.
}

type RegisterInfrastructureDeps struct {
	TransactionManager database.TransactionManager

	<Aggregate>Repository <module>Domain<Aggregate>.<Entity>Repository
}

func RegisterInfrastructure(deps RegisterInfrastructureDeps) *Handlers {
	return &Handlers{}
}
```

Initially empty. The `command-create` and `query-create` skills add `Handlers` fields and instantiations as endpoints are scaffolded. Even an empty `Handlers` struct is correct — it's the boundary the presentation layer plugs into. When those skills add endpoints, the `Handlers` fields are typed as `cqrs.Handler[*Command, *Response]` interfaces and each handler is wrapped with `cqrs.ValidationBehavior(...)` (from `internal/common/application/cqrs`) so validation runs uniformly via the behavior decorator — they add the `cqrs` import then. (Memory: `project_cqrs_validation_behavior`.)

`RegisterInfrastructure` takes a single `RegisterInfrastructureDeps` struct (NOT positional params) — the core infra dep (`TransactionManager database.TransactionManager`, from `go-packages/database`) goes first, then a blank line, then repositories and cross-module services. If the module needs additional dependencies (secret manager, cross-module services), add them to `RegisterInfrastructureDeps` AND to the module's `Dependencies` struct. They stay unused until the first endpoint reads them.

### `presentation/presentation.go`

```go
package presentation

import (
	<module>Infrastructure "github.com/blocknextai/platform-api/internal/<module>/infrastructure"
	commonPresentationAuth "github.com/blocknextai/platform-api/internal/common/presentation/auth"
	"github.com/gofiber/fiber/v3"
)

func RegisterPresentation(
	router fiber.Router,
	authMiddleware *commonPresentationAuth.AuthMiddleware,
	handlers *<module>Infrastructure.Handlers,
) {
	// Routes are added by command-create / query-create skills.
}
```

Same shape as `infrastructure.go` — initially empty, fills up as endpoints are scaffolded. Sub-functions for user vs organization route groups (`RegisterUser<Aggregate>Presentation`, `RegisterOrganization<Aggregate>Presentation`) are introduced when the first endpoint that needs them is added.

### `module.go`

```go
package <module>

import (
	"database/sql"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/platform-api/internal/common/presentation/auth"
	<module>Infrastructure "github.com/blocknextai/platform-api/internal/<module>/infrastructure"
	<module>Infrastructure<Aggregate> "github.com/blocknextai/platform-api/internal/<module>/infrastructure/<aggregate>"
	<module>Presentation "github.com/blocknextai/platform-api/internal/<module>/presentation"
	"github.com/gofiber/fiber/v3"
)

type Dependencies struct {
	DB                 *sql.DB
	TransactionManager database.TransactionManager
	// ... cross-module service interfaces if any
}

type Module struct {
	// Public fields ONLY when consumed by another module.
	// Add (e.g.) `<Aggregate>Service <module>Application<Aggregate>.<Aggregate>Service` here only if some other module's Dependencies needs it.

	handlers *<module>Infrastructure.Handlers
}

func NewModule(deps Dependencies) *Module {
	repository := <module>Infrastructure<Aggregate>.New<Entity>Repository(deps.DB)
	handlers := <module>Infrastructure.RegisterInfrastructure(<module>Infrastructure.RegisterInfrastructureDeps{
		TransactionManager: deps.TransactionManager,

		<Aggregate>Repository: repository,
	})
	return &Module{
		handlers: handlers,
	}
}

func (m *Module) Register(router fiber.Router, authMiddleware *auth.AuthMiddleware) {
	<module>Presentation.RegisterPresentation(router, authMiddleware, m.handlers)
}
```

Rules (project memory):
- **`Dependencies` struct lists every external input.** DB and `TransactionManager` are the baseline; cross-module service interfaces, secret manager, etc. extend it as needed.
- **`Module` struct fields are private by default.** A field becomes public ONLY when it is read in another module's `NewModule` call. (Memory: `feedback_module_field_pattern`.) Do not preemptively expose services "in case someone needs them later." No getters, no aliases.
- **No two-phase init.** If you find yourself wanting to construct the module in one call and "wire it up" in a second, you're hitting a cycle — break it via composition, events, or rethinking the dependency direction. Two-phase init is debt, not a pattern. (Memory: `feedback_module_two_phase_init`.)
- **No lazy supplier / decorator wrapping** to break cycles. (Memory: `feedback_lazy_supplier_rejected`.)
- **`NewModule` returns `*Module`**, not an interface. The caller is `cmd/platform-api/main.go` and a few other entry points; they all use the concrete type.
- **`Register(router, authMiddleware, ...)`** signature can vary slightly per module (some take `cacheMiddleware`, some don't). Match the existing modules' styles when there's a similar one.

### Schema-init migration

Use the `migration-create` skill (or run `make migration-create` directly) to generate the first migration:

```bash
make migration-create name=init_schema module=<module>
```

Then fill the generated `.up.sql`:

```sql
-- Migration: init_schema
-- Created: ...

CREATE SCHEMA IF NOT EXISTS <module>;
```

And `.down.sql`:

```sql
-- Rollback: init_schema
-- Created: ...

DROP SCHEMA IF EXISTS <module> CASCADE;
```

This is the only migration this skill produces. The first table migration (`<TS>_<module>_add_<tables>_table`) is a separate `migration-create` invocation — defer to that skill.

## Wiring

Wiring is split across **two files**:
- **`internal/bootstrap/platform_api.go`** — instantiation (the modules are constructed here in `NewPlatformAPI`, NOT in `cmd/platform-api/main.go`).
- **`cmd/platform-api/main.go`** — route registration off the `app` (`*bootstrap.PlatformAPI`) returned by `NewPlatformAPI`.

### `internal/bootstrap/platform_api.go` — three edits

**1. Import** (alphabetically placed in the import block):

```go
"github.com/blocknextai/platform-api/internal/<module>"
```

**2. Add a field to the `PlatformAPI` struct** (`<Module>Module *<module>.Module`), matching the column-aligned style of the existing fields.

**3. Instantiate inside `NewPlatformAPI`.** The existing instantiations follow a rough dependency order — modules that produce services come before modules that consume them. Place the new `<module>Module := <module>.NewModule(...)` right before the first module that depends on it (or near the end if nothing depends on it yet). `deps` come from `core` (`core.DB`, `core.TransactionManager`, `core.EventBus.Bus`, `core.CacheService`, ...) and `shared`/`cfg` for options:

```go
<module>Module := <module>.NewModule(<module>.Dependencies{
	DB:                 core.DB,
	TransactionManager: core.TransactionManager,
	// + any cross-module services this module needs (e.g. accountModule.UserService)
})
```

If `NewModule` returns an error (some modules do — e.g., `account`, `web3`, `llm`, `workflows`), use the early-return pattern (`NewPlatformAPI` returns `(*PlatformAPI, error)`):

```go
<module>Module, err := <module>.NewModule(<module>.Dependencies{...})
if err != nil {
	return nil, err
}
```

For the minimal scaffold this skill produces, `NewModule` does NOT return an error. Only add an error return if construction can genuinely fail (e.g., connecting to an external broker).

**Also add the field to the returned struct literal** at the bottom of `NewPlatformAPI` (`<Module>Module: <module>Module,`).

### `cmd/platform-api/main.go` — register routes

Add a `app.<Module>Module.Register(...)` line to the block of `Register` calls (the `app` variable is the `*bootstrap.PlatformAPI`). Match the existing line style:

```go
app.<Module>Module.Register(fiberApp, authMiddleware)
```

If the module needs `cacheMiddleware`, match modules like `account`, `organizations`, `plans`:

```go
app.<Module>Module.Register(fiberApp, authMiddleware, cacheMiddleware)
```

## Critical rules (project memory)

1. **Module migration one at a time.** This is part of a larger refactor — apply the module pattern to one module per turn, verify it builds, then move on. (Memory: `feedback_module_migration_pace`.)
2. **Private by default, public only when consumed.** No speculative exposure. No getters. No `main.go`-side aliases. (Memory: `feedback_module_field_pattern`.)
3. **No two-phase init.** Construct fully in `NewModule`. If a cycle forces this, redesign. (Memory: `feedback_module_two_phase_init`.)
4. **No lazy supplier / decorator pattern.** Rejected as too ugly; don't propose it. (Memory: `feedback_lazy_supplier_rejected`.)
5. **No repository injection across bounded contexts.** Cross-module dependency = service interface, period. (Memory: `feedback_no_repository_injection_across_contexts`.)
6. **`apperror` everywhere.** Domain: Validation/NotFound. App: Internal + WithCause. No `errors.New` / `fmt.Errorf` in business code. (Memory: `project_apperror_pattern`.)
7. **Validation runs via the `cqrs.ValidationBehavior` decorator, not in the handler.** When commands are added later, the `Validate()` method lives in `validator.go` but is NOT called in `Handle`; the `cqrs.ValidationBehavior(...)` wrapper applied in `infrastructure.go` runs it (runtime `Validatable` check) before the inner handler / transaction. Handler types and constructors use the idiomatic `Handler` / `New` form (not `<Verb><Entity>CommandHandler` / `New...`), since the package name carries the context. (Memory: `project_cqrs_validation_behavior`, `feedback_validation_layering`, `feedback_event_handler_new_naming`.)
8. **`bnuuid.NewV7()`** (`go-packages/uuid`) for IDs. **`time.Now().UTC()`** for timestamps. **`json.Marshal`/`Unmarshal`** from `go-packages/json` (NOT `encoding/json`, NOT a `utils` package). **`slog`** for logging. **`strings.TrimSpace(s) == ""`** for empty checks. **`new(value)`** over `&value`.
9. **Don't use `fmt`.** (Memory: `feedback_no_fmt_package`.)

## Final steps

After the scaffold + migration + wiring:

1. `gofmt -w` on every new file + `internal/bootstrap/platform_api.go` + `cmd/platform-api/main.go`.
2. `go build ./...` — must compile. If it fails, the most common causes are:
   - Forgot to add the import / struct field / instantiation / returned-literal entry in `internal/bootstrap/platform_api.go`
   - Forgot the `app.<Module>Module.Register(...)` line in `cmd/platform-api/main.go`
   - Cross-module dependency in `Dependencies` references a service that doesn't exist yet
   - Repository scan column order doesn't match the table's column order
3. `golangci-lint run ./internal/<module>/...` — must be clean.
4. (Optional, requires user confirmation) `make migration-up module=<module>` to apply the schema-init migration.
5. Hand off to the user with the next steps spelled out:
   - "The module scaffold is ready. To add the first table, run the `migration-create` skill. To add the first endpoint, run `command-create` (write) or `query-create` (read)."

## What NOT to do

- Do NOT add multiple aggregates in this skill. One aggregate per module on initial scaffold; additional aggregates are separate work.
- Do NOT pre-populate `application/<aggregate>/` with empty `create<entity>/`, `get<entity>/` folders. Those come from `command-create` / `query-create` and contain real code.
- Do NOT expose `Module.Repository<Entity>` as a public field unless another module is genuinely about to consume it. Speculative exposure violates the visibility rule.
- Do NOT inject `database.TransactionManager` into the service or query handlers — only command handlers need it. Putting it in the wrong place is a tell that the boundary is confused.
- Do NOT add `// TODO` placeholders in the scaffold. Either the scaffold compiles and runs (empty Handlers + empty routes is fine) or it does not — half-done is worse than missing.
- Do NOT skip the wiring. A module that is not instantiated in `internal/bootstrap/platform_api.go` AND registered in `cmd/platform-api/main.go` is dead code; the user has no way to reach it.
- Do NOT run `make migration-up` autonomously even though the schema-init is a "small" change. The DB mutation rule applies regardless of size.
- Do NOT introduce new common helpers (`internal/common/...`) as part of this skill. If the user asks for one, that's a separate change with its own justification.
