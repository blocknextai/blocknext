-- Migration: apikeys_add_api_keys_table
-- Created: Mon Apr 13 00:00:03 +03 2026

CREATE TABLE IF NOT EXISTS apikeys.api_keys (
    id UUID         PRIMARY KEY,
    owner_type      VARCHAR(50) NOT NULL,
    owner_id        UUID NOT NULL,
    name            VARCHAR(255) NOT NULL,
    key_hash        VARCHAR(255) NOT NULL,
    scopes          TEXT[] NOT NULL DEFAULT '{}',
    last_used_at    TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL,
    created_at      TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    updated_at      TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    deleted_at      TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_api_keys_key_hash
    ON apikeys.api_keys(key_hash)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_api_keys_owner
    ON apikeys.api_keys(owner_type, owner_id)
    WHERE deleted_at IS NULL;
