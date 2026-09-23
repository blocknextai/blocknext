-- Migration: add_verification_tokens_table
-- Single table for email_verify and password_reset purposes. Token plaintext is never stored — only its hash.

CREATE TABLE IF NOT EXISTS account.verification_tokens (
    id            UUID PRIMARY KEY,
    user_id       UUID NOT NULL,
    purpose       VARCHAR(32) NOT NULL,
    token_hash    TEXT NOT NULL,
    email         VARCHAR(255) NOT NULL,
    expires_at    TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    created_at    TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    updated_at    TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    deleted_at    TIMESTAMP WITHOUT TIME ZONE
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_account_verification_tokens_hash
    ON account.verification_tokens (token_hash)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_account_verification_tokens_user_purpose
    ON account.verification_tokens (user_id, purpose)
    WHERE deleted_at IS NULL;
