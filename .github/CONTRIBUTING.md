# Contributing to BlockNext

Thanks for your interest in contributing! This document explains how the
repository is organized and how to get a working development environment.

## Repository layout

```
apps/
  platform-api/       # Main API — Go modular monolith (DDD + Clean Architecture)
  file-gateway-api/   # File upload/download gateway — Go
  platform/           # End-user UI — React + Vite (Bun workspace)
packages/
  go-packages/        # Shared Go libraries
docker/               # Dockerfiles per app
scripts/              # Everything the Makefile targets run
```

Go modules are tied together with the root `go.work`; the JS side uses Bun
workspaces + Turborepo. All configuration lives in the single root `.env`
(see `.env.example`).

## Getting started

Requirements: Docker + Docker Compose, Make. (Go 1.26+ and [Bun](https://bun.sh)
are only needed for host-side development — see below.)

```bash
make setup            # creates .env with generated secrets
make local-docker-up  # builds and starts the full stack from source
```

The UI is served on http://localhost:4000; APIs on 3000 (platform),
3100 (mcp), 3200 (webhook), 3300 (file gateway). Run `make help` for every
available target.

To work on the code directly on the host — linters, formatters, the UI dev
server — install the Go and Bun dependencies too:

```bash
make setup-dev        # setup + go mod download + bun install
```

For a faster UI feedback loop, stop the `platform` container and run the dev
server on the host:

```bash
cd apps/platform && bun run dev
```

## Making changes

### Go

- Read `apps/platform-api/CLAUDE.md` and the `README.md` inside each
  `internal/<module>/` directory before changing a bounded context — the
  architecture conventions (CQRS layout, cross-context service interfaces,
  repository patterns) are documented there and enforced in review.
- Format and lint before pushing: `make go-fmt && make go-lint`.
- `make go-build` is the fastest correctness check across all modules.
- Database changes go through `make migration-create name=<n> module=<m>`;
  never hand-write migration files.

### UI

- Feature code follows the `service ↔ use-<feature>` hook pattern; pages do
  not hold state themselves.
- `bun run lint` and `bun run build` must pass.

## Commit messages

Use [Conventional Commits](https://www.conventionalcommits.org/):

```
feat(workflows): add cron trigger timezone support
fix(file-gateway): reject uploads above the rule size limit
```

## Pull requests

- Keep PRs focused — one logical change per PR.
- Describe what changed and why; link related issues.
- Make sure `make go-build`, `make go-lint` and the UI build pass locally.

## Reporting bugs & security issues

- Bugs and feature discussions: GitHub issues.
- Security vulnerabilities: **never** open a public issue — see
  [SECURITY.md](SECURITY.md).

## Code of Conduct

By participating you agree to abide by our
[Code of Conduct](CODE_OF_CONDUCT.md).

## License

By contributing, you agree that your contributions will be licensed under the
[Apache License 2.0](../LICENSE).
