-- Migration: add_authorization_codes_table
-- Created: Wed Sep 10 01:34:45 +03 2026

CREATE TABLE IF NOT EXISTS mcpoauth.authorization_codes (
    id UUID                 PRIMARY KEY,
    grant_id                UUID NOT NULL,
    client_id               VARCHAR(512) NOT NULL,
    code_hash               VARCHAR(255) NOT NULL,
    redirect_uri            TEXT NOT NULL,
    code_challenge          VARCHAR(255) NOT NULL,
    code_challenge_method   VARCHAR(10) NOT NULL,
    resource                TEXT NOT NULL,
    scopes                  TEXT[] NOT NULL DEFAULT '{}',
    expires_at              TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    used_at                 TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL,
    created_at              TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    updated_at              TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    deleted_at              TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_mcpoauth_authorization_codes_code_hash
    ON mcpoauth.authorization_codes(code_hash)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_mcpoauth_authorization_codes_grant
    ON mcpoauth.authorization_codes(grant_id)
    WHERE deleted_at IS NULL;
