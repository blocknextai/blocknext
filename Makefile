.DEFAULT_GOAL := help

##@ general

.PHONY: help
help: ## show this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nusage:\n  make \033[36m<target>\033[0m\n"} \
		/^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-24s\033[0m %s\n", $$1, $$2 } \
		/^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) }' $(MAKEFILE_LIST)

.PHONY: setup
setup: ## create .env with generated secrets (all you need to run the stack)
	@./scripts/setup.sh

.PHONY: setup-dev
setup-dev: ## setup + install Go and Bun dependencies for local development
	@./scripts/setup-dev.sh

##@ go

.PHONY: go-fmt
go-fmt: ## format source files in every Go module
	@./scripts/go.sh fmt

.PHONY: go-lint
go-lint: ## run golangci-lint in every Go module
	@./scripts/go.sh lint

.PHONY: go-fix
go-fix: ## run go fix in every Go module
	@./scripts/go.sh fix

.PHONY: go-build
go-build: ## compile every Go module
	@./scripts/go.sh build

.PHONY: go-test
go-test: ## run tests in every Go module
	@./scripts/go.sh test

.PHONY: go-update-deps
go-update-deps: ## update and tidy dependencies in every Go module
	@./scripts/go.sh update-deps

##@ migrations

.PHONY: migration-create
migration-create: ## create a timestamped migration pair (name=<n> module=<m>)
	@./scripts/migration.sh create --name=$(name) --module=$(module)

.PHONY: migration-up
migration-up: ## apply migrations (module=<m> dry-run=1)
	@./scripts/migration.sh up $(if $(module),--module=$(module)) $(if $(dry-run),--dry-run)

.PHONY: migration-down
migration-down: ## roll back migrations (module=<m> dry-run=1)
	@./scripts/migration.sh down $(if $(module),--module=$(module)) $(if $(dry-run),--dry-run)

.PHONY: migration-status
migration-status: ## show migration state (module=<m> direction=up|down)
	@./scripts/migration.sh status $(if $(module),--module=$(module)) $(if $(dry-run),--dry-run) $(if $(direction),--direction=$(direction))

##@ docker (prod: ghcr.io/blocknextai/* images)

.PHONY: docker-pull
docker-pull: ## pull the published images
	@./scripts/docker.sh pull

.PHONY: docker-up
docker-up: ## start services in the background
	@./scripts/docker.sh up

.PHONY: docker-down
docker-down: ## stop services
	@./scripts/docker.sh down

.PHONY: docker-logs
docker-logs: ## follow service logs
	@./scripts/docker.sh logs

.PHONY: docker-ps
docker-ps: ## list service status
	@./scripts/docker.sh ps

.PHONY: docker-clean
docker-clean: ## remove services and volumes
	@./scripts/docker.sh clean

##@ local docker (builds from ./docker)

.PHONY: local-docker-build
local-docker-build: ## build all service images
	@./scripts/docker.sh local-build

.PHONY: local-docker-up
local-docker-up: ## start services in the background
	@./scripts/docker.sh local-up

.PHONY: local-docker-down
local-docker-down: ## stop services
	@./scripts/docker.sh local-down

.PHONY: local-docker-logs
local-docker-logs: ## follow service logs
	@./scripts/docker.sh local-logs

.PHONY: local-docker-ps
local-docker-ps: ## list service status
	@./scripts/docker.sh local-ps

.PHONY: local-docker-clean
local-docker-clean: ## remove services and volumes
	@./scripts/docker.sh local-clean
