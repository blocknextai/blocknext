-- Migration: add_grants_table
-- Created: Wed Sep 10 01:34:45 +03 2026

CREATE TABLE IF NOT EXISTS mcpoauth.grants (
    id UUID             PRIMARY KEY,
    client_id           VARCHAR(512) NOT NULL,
    user_id             UUID NOT NULL,
    organization_id     UUID NOT NULL,
    scopes              TEXT[] NOT NULL DEFAULT '{}',
    resource            TEXT NOT NULL,
    revoked_at          TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL,
    created_at          TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    updated_at          TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    deleted_at          TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS idx_mcpoauth_grants_client_user_organization
    ON mcpoauth.grants(client_id, user_id, organization_id)
    WHERE revoked_at IS NULL AND deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_mcpoauth_grants_user
    ON mcpoauth.grants(user_id)
    WHERE revoked_at IS NULL AND deleted_at IS NULL;
