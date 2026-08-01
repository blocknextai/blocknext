#!/usr/bin/env bash

set -e

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
GO_DIRS="packages/go-packages apps/platform-api apps/file-gateway-api"

usage() {
    echo "Usage: $0 <command>"
    echo ""
    echo "Commands:"
    echo "  fmt           Format source files in every Go module"
    echo "  lint          Run golangci-lint in every Go module"
    echo "  fix           Run go fix in every Go module"
    echo "  build         Compile every Go module"
    echo "  test          Run tests in every Go module"
    echo "  update-deps   Update and tidy dependencies in every Go module"
    echo ""
    exit 1
}

run() {
    for dir in $GO_DIRS; do
        (cd "$REPO_ROOT/$dir" && "$@")
    done
}

case $1 in
    fmt)
        run go fmt ./...
        ;;
    lint)
        run golangci-lint run ./...
        ;;
    fix)
        run go fix ./...
        ;;
    build)
        run go build ./...
        ;;
    test)
        run go test ./... -v
        ;;
    update-deps)
        run sh -c 'go get -u ./... && go mod tidy'
        ;;
    *)
        usage
        ;;
esac
