#!/usr/bin/env bash

set -e

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
COMPOSE="docker compose --project-directory $REPO_ROOT -f $REPO_ROOT/docker-compose.prod.yml"
COMPOSE_LOCAL="docker compose --project-directory $REPO_ROOT -f $REPO_ROOT/docker-compose.local.yml"

usage() {
    echo "Usage: $0 <command>"
    echo ""
    echo "Prod (ghcr.io/blocknextai/* images):"
    echo "  pull          Pull the published images"
    echo "  up            Start services in the background"
    echo "  down          Stop services"
    echo "  logs          Follow service logs"
    echo "  ps            List service status"
    echo "  clean         Remove services and volumes"
    echo ""
    echo "Local (builds from ./docker):"
    echo "  local-build   Build all service images"
    echo "  local-up      Start services in the background"
    echo "  local-down    Stop services"
    echo "  local-logs    Follow service logs"
    echo "  local-ps      List service status"
    echo "  local-clean   Remove services and volumes"
    echo ""
    exit 1
}

case $1 in
    pull)
        $COMPOSE pull
        ;;
    up)
        $COMPOSE up -d
        ;;
    down)
        $COMPOSE down
        ;;
    logs)
        $COMPOSE logs -f
        ;;
    ps)
        $COMPOSE ps
        ;;
    clean)
        $COMPOSE down -v --remove-orphans
        ;;
    local-build)
        DOCKER_BUILDKIT=1 $COMPOSE_LOCAL build
        ;;
    local-up)
        DOCKER_BUILDKIT=1 $COMPOSE_LOCAL up -d --build
        ;;
    local-down)
        $COMPOSE_LOCAL down
        ;;
    local-logs)
        $COMPOSE_LOCAL logs -f
        ;;
    local-ps)
        $COMPOSE_LOCAL ps
        ;;
    local-clean)
        $COMPOSE_LOCAL down -v --remove-orphans
        ;;
    *)
        usage
        ;;
esac
