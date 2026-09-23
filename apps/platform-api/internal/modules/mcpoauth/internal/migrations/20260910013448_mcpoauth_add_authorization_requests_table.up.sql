-- Migration: add_authorization_requests_table
-- Created: Wed Sep 10 01:34:45 +03 2026

CREATE TABLE IF NOT EXISTS mcpoauth.authorization_requests (
    id UUID                 PRIMARY KEY,
    client_id               VARCHAR(512) NOT NULL,
    redirect_uri            TEXT NOT NULL,
    scopes                  TEXT[] NOT NULL DEFAULT '{}',
    state                   TEXT DEFAULT NULL,
    code_challenge          VARCHAR(255) NOT NULL,
    code_challenge_method   VARCHAR(10) NOT NULL,
    resource                TEXT NOT NULL,
    status                  VARCHAR(20) NOT NULL,
    user_id                 UUID DEFAULT NULL,
    organization_id         UUID DEFAULT NULL,
    expires_at              TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    resolved_at             TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL,
    created_at              TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    updated_at              TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    deleted_at              TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS idx_mcpoauth_authorization_requests_expires_at
    ON mcpoauth.authorization_requests(expires_at)
    WHERE deleted_at IS NULL;
