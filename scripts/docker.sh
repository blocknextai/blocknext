#!/usr/bin/env bash

set -e

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
COMPOSE="docker compose --project-directory $REPO_ROOT -f $REPO_ROOT/docker-compose.prod.yml"
ALL_PROFILES='COMPOSE_PROFILES=*' 
COMPOSE_DEV="docker compose --project-directory $REPO_ROOT -f $REPO_ROOT/docker-compose.dev.yml"

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
    echo "Dev (builds from ./docker):"
    echo "  dev-build     Build all service images"
    echo "  dev-up        Start services in the background"
    echo "  dev-down      Stop services"
    echo "  dev-logs      Follow service logs"
    echo "  dev-ps        List service status"
    echo "  dev-clean     Remove services and volumes"
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
        env "$ALL_PROFILES" $COMPOSE down
        ;;
    logs)
        $COMPOSE logs -f
        ;;
    ps)
        $COMPOSE ps
        ;;
    clean)
        env "$ALL_PROFILES" $COMPOSE down -v --remove-orphans
        ;;
    dev-build)
        DOCKER_BUILDKIT=1 $COMPOSE_DEV build
        ;;
    dev-up)
        DOCKER_BUILDKIT=1 $COMPOSE_DEV up -d --build
        ;;
    dev-down)
        env "$ALL_PROFILES" $COMPOSE_DEV down
        ;;
    dev-logs)
        $COMPOSE_DEV logs -f
        ;;
    dev-ps)
        $COMPOSE_DEV ps
        ;;
    dev-clean)
        env "$ALL_PROFILES" $COMPOSE_DEV down -v --remove-orphans
        ;;
    *)
        usage
        ;;
esac
