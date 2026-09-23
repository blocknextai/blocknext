-- Migration: account_add_nonces_table
-- Created: Tue Sep  9 01:48:06 +03 2025

-- Create account.nonces table
CREATE TABLE IF NOT EXISTS account.nonces (
    id                   UUID PRIMARY KEY,
    auth_provider        VARCHAR(50) NOT NULL,
    provider_id          TEXT,
    nonce                VARCHAR(255) NOT NULL,
    code_verifier        VARCHAR(255) NOT NULL,
    code_challenge       VARCHAR(255) NOT NULL,
    code_challenge_method VARCHAR(50) NOT NULL,
    expires_at           TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    created_at           TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    updated_at           TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    deleted_at           TIMESTAMP WITHOUT TIME ZONE
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_account_nonces_nonce
    ON account.nonces(nonce)
    WHERE deleted_at IS NULL;
