# platform-api-migration

> A one-shot CLI that runs golang-migrate database migrations per module, each tracked in its own `<module>_migrations` table.

## What it does
Parses flags, connects to Postgres from `DATABASE_*` env vars, and applies (or rolls back) the `.up.sql` / `.down.sql` migrations found under each module's `infrastructure/database/migrations` directory. Modules are processed in the order of a hardcoded `modules` registry; `up` runs them forward, `down` rolls back one step in reverse order. A `-dry-run` mode reports current version, dirty state, and pending counts without applying anything. It runs once and exits.

## Bootstrap & config
- **Assembler:** none — does not use `bootstrap.NewCore`; opens its own `database.NewDB` connection (pool size 1/1).
- **Config:** parses `config.DatabaseOptions` directly with env prefix `DATABASE_` (requires `DATABASE_HOST` and `DATABASE_NAME`).
- **Runs as:** one-shot command.

## Flags
- `-direction` — `up` (default) or `down`.
- `-module` — restrict to a single module (must be in the registry); empty = all.
- `-dry-run` — report migration status without applying.

## Notes
- **Module registry requirement:** the hardcoded `modules` slice (common, account, organizations, executions, support, workflows, triggers, web3, marketplace, subscriptions, quota, payments, library, credentials, billing, apikeys, notifications, eventbus) is the source of truth. A new module's migrations will NOT run unless it is added to this slice.
- Each module uses a dedicated `<module>_migrations` tracking table.
- Migration source path differs in Docker (`/app/migrations/<module>`) vs local (`internal/<module>/infrastructure/database/migrations`), auto-detected via `/app/migrations`.
- Exits non-zero on failure; logs structured slog with `component=migration`.
