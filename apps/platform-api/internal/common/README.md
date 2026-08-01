# Common

> Shared kernel for the platform: cross-cutting presentation, application, and domain code consumed by every bounded context. This is NOT a business bounded context — it owns no aggregates or persistence of its own.

## Responsibility
Houses the code that does not belong to any single context but is needed by all of them: HTTP auth/authz middleware, request-context helpers, the generic CQRS behavior pipeline, shared value types (owner type, execution context, billing period, etc.), reusable input validators, and factories for shared infrastructure (email sender, password hasher). `module.go` is thin — it wires only the two stateful infrastructure services; most of the package is stateless helpers imported directly.

## What it provides
- **`presentation/auth.AuthMiddleware`** — JWT-based `Authenticate()` / `AuthenticateWebSocket()`, plus `RequireUserPermission` and `RequireOrganizationPermission` (RBAC checks via `commonApplicationAuth` checkers; org scope read from the `:organizationId` path param). Validates sessions through `account` `SessionService`.
- **`presentation/auth.APIKeyMiddleware`** — `Authenticate()` (validates `X-API-Key`, resolves owner type/ID into locals) and `RequireScope(scope)` for per-key capability scopes.
- **`presentation/http`** — Fiber request-context helpers (`SetUserID`/`GetUserID`, session ID, organization ID, WebSocket variants), `RespondPaginated`, and shared error sentinels.
- **`application/cqrs`** — generic mediator-style pipeline: `Handler[C,R]`, `HandlerFunc`, `Behavior`, `Chain`, and `ValidationBehavior` (runs `Validate()` on any `Validatable` command before the handler).
- **`application/auth`** — `UserPermissionChecker` / `OrganizationPermissionChecker` interfaces (implemented elsewhere, injected into the auth middleware).
- **`domain`** — shared value types/enums: `OwnerType`, `ExecutionContext`, `DomainEvent` interface, billing-period type, plus `domain/credential` and `domain/oauth2` shared structs.
- **`validation`** — reusable input validators (`Email`, `Password`).
- **`infrastructure`** — provider-selecting factories `NewEmailSender` (log/SMTP/Resend/SendGrid) and `NewPasswordHasher` (bcrypt). Exposed on `Module` as `EmailSender` and `PasswordHasher`.

## Consumed by
Effectively every context. The auth middleware and `presentation/http` helpers guard and parse requests across all HTTP/WebSocket routes; the `cqrs` pipeline backs all command/query handlers; `domain` value types and `validation` are imported throughout; the `Module`'s `EmailSender`/`PasswordHasher` are wired from the bootstrap entry points (`platform_api`, `mcp_api`, `task_worker`, `webhook_api`).

## Dependencies
- **Infrastructure / external:** go-packages (`apperror`, `auth/jwt`, `rbac`, `email/*`, `hashing/bcrypt`, `result`), Fiber v3 (+ websocket contrib), `google/uuid`. Cross-context imports: `account` `SessionService`, `apikeys` validator/scope types.

## Layout
- `application/auth/` — permission-checker interfaces.
- `application/cqrs/` — handler + behavior pipeline.
- `domain/` — shared enums, the `DomainEvent` interface, and `domain/credential`, `domain/oauth2` subpackages.
- `infrastructure/` — email-sender and password-hasher factories.
- `presentation/auth/` — JWT and API-key middleware.
- `presentation/http/` — context-locals helpers, pagination response, error sentinels.
- `validation/` — email/password validators.
- `module.go` — wires `EmailSender` and `PasswordHasher` only.
