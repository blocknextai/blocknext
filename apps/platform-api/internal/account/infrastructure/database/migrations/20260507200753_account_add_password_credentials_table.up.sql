-- Migration: add_password_credentials_table
-- Stores per-user password hash for email/password auth. Identity row lives in linked_accounts.

CREATE TABLE IF NOT EXISTS account.password_credentials (
    id               UUID PRIMARY KEY,
    user_id          UUID NOT NULL,
    password_hash    TEXT NOT NULL,
    created_at       TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    updated_at       TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    deleted_at       TIMESTAMP WITHOUT TIME ZONE
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_account_password_credentials_user
    ON account.password_credentials(user_id)
    WHERE deleted_at IS NULL;
