#!/usr/bin/env bash

set -e

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

cd "$REPO_ROOT"

if [ -f .env ]; then
    echo ".env already exists — leaving it untouched"
    exit 0
fi

cp .env.example .env

for token in $(grep -oE 'REPLACE_ME_OPENSSL_[A-Z_]+' .env | sort -u); do
    secret="$(openssl rand -hex 32)"
    sed "s/$token/$secret/g" .env > .env.tmp && mv .env.tmp .env
done

if grep -q REPLACE_ME_OPENSSL .env; then
    echo "warning: unreplaced secret placeholders remain in .env"
    exit 1
fi

echo ".env created with generated secrets — run 'make docker-up' to start the stack"
