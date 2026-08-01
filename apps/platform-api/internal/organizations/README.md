# Organizations

> Tenants/workspaces and their members (RBAC roles).

## Responsibility
This context owns organizations (multi-tenant workspaces) and their memberships with RBAC roles. Organizations are created explicitly by users, enforces organization-scoped permissions for the rest of the system, and exposes the membership service consumed by other contexts.

## Domain
- **Aggregates / key types:**
  - `organizations.Organization` — a tenant/workspace (`Title`, `Description`, `IsVerified`).
  - `organizationusers.OrganizationUser` — membership linking a user to an organization with a `Role` and `Alias`.
- **Key rules / invariants:**
  - `Organization`: `Title` required; `IsVerified` is immutable via update.
  - `OrganizationUser`: `Role` must be a valid RBAC organization role; membership unique per (org, user).

## Use cases (application)
- **Organizations:** `createorganization`, `updateorganization`, `deleteorganization` (commands); `getorganizationbyid`, `getallorganizations` (queries).
- **Organization users:** `createorganizationuser`, `updateorganizationuserinfo`, `updateorganizationuserrole`, `deleteorganizationuser` (commands); `getallorganizationusers`, `getorganizationuserbyuserid`, `getorganizationme`, `getroles` (queries).
- **`auth.OrganizationPermissionChecker`** — resolves a user's org role (cache-first) and checks `rbac.HasOrganizationPermission`; consumed by middleware across contexts.

## HTTP API
Base `/organizations`. All routes use `Authenticate()` + an RBAC permission, except `/roles` (cached, no auth).

- `GET /` — list (`RequireUserPermission(ReadOrganizationPermission)`)
- `GET /roles` — list roles (cached 5m, no auth)
- `POST /` — create (`RequireUserPermission(CreateOrganizationPermission)`)
- `GET /:organizationId` · `GET /:organizationId/me` — read (`RequireOrganizationPermission(ReadOrganizationPermission)`)
- `PUT /:organizationId` — update (`UpdateOrganizationPermission`)
- `DELETE /:organizationId` — delete (`DeleteOrganizationPermission`)
- `/:organizationId/users` — `GET /`, `POST /`, `GET /:userId`, `PUT /:userId/info`, `PUT /:userId/role`, `DELETE /:userId` (org-user permissions)

## Events
- **Published:** `organizations.organization_user.created` (on owner membership and member add), `organizations.organization_user.role_changed`.
- **Consumed:** none.

## Dependencies
- **Bounded contexts:** account (`UserService`, `LinkedAccountService`).
- **Infrastructure:** Postgres schema `organizations` (tables `organizations.organizations`, `organizations.users`); cache (permission lookups); eventbus + inbox idempotency. No external APIs.

## Layout
Standard DDD layers present, wired in `module.go`.
