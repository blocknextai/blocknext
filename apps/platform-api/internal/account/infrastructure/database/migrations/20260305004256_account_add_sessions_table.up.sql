-- Migration: add_sessions_table
-- Created: Thu Mar 5 00:42:56 +03 2026

CREATE TABLE IF NOT EXISTS account.sessions (
    id                        UUID PRIMARY KEY,
    user_id                   UUID NOT NULL,
    auth_provider             VARCHAR(20) NOT NULL,
    ip_address                VARCHAR(45) NOT NULL DEFAULT '',
    user_agent                TEXT NOT NULL DEFAULT '',
    is_revoked                BOOLEAN NOT NULL DEFAULT FALSE,
    refresh_token_hash        TEXT,
    refresh_token_expires_at  TIMESTAMP WITHOUT TIME ZONE,
    token_family              UUID,
    token_generation          INTEGER NOT NULL DEFAULT 0,
    created_at                TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    updated_at                TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    deleted_at                TIMESTAMP WITHOUT TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_account_sessions_user_id
    ON account.sessions(user_id) WHERE deleted_at IS NULL AND is_revoked = FALSE;

CREATE INDEX IF NOT EXISTS idx_account_sessions_token_family
    ON account.sessions(token_family) WHERE deleted_at IS NULL AND is_revoked = FALSE;
