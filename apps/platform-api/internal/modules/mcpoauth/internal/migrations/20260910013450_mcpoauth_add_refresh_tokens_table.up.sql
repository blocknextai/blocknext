-- Migration: add_refresh_tokens_table
-- Created: Wed Sep 10 01:34:45 +03 2026

CREATE TABLE IF NOT EXISTS mcpoauth.refresh_tokens (
    id UUID         PRIMARY KEY,
    grant_id        UUID NOT NULL,
    client_id       VARCHAR(512) NOT NULL,
    token_hash      VARCHAR(255) NOT NULL,
    scopes          TEXT[] NOT NULL DEFAULT '{}',
    resource        TEXT NOT NULL,
    expires_at      TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    used_at         TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL,
    revoked_at      TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL,
    created_at      TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    updated_at      TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    deleted_at      TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_mcpoauth_refresh_tokens_token_hash
    ON mcpoauth.refresh_tokens(token_hash)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_mcpoauth_refresh_tokens_grant
    ON mcpoauth.refresh_tokens(grant_id)
    WHERE deleted_at IS NULL;
