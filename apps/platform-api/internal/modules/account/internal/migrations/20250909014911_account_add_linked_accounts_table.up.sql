-- Migration: account_add_linked_accounts_table
-- Created: Tue Sep  9 01:49:11 +03 2025

-- Create account.linked_accounts table
CREATE TABLE IF NOT EXISTS account.linked_accounts (
    id            UUID PRIMARY KEY,
    user_id       UUID NOT NULL,
    auth_provider VARCHAR(50) NOT NULL,
    provider_id   VARCHAR(255) NOT NULL,
    identifier    VARCHAR(255) NOT NULL,
    display_name  VARCHAR(255),
    is_primary    BOOLEAN NOT NULL DEFAULT FALSE,
    is_verified   BOOLEAN NOT NULL DEFAULT FALSE,
    created_at    TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    updated_at    TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    deleted_at    TIMESTAMP WITHOUT TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_account_linked_accounts_user_id
    ON account.linked_accounts(user_id);

CREATE UNIQUE INDEX IF NOT EXISTS uq_account_linked_accounts_provider_lookup
    ON account.linked_accounts(auth_provider, provider_id)
    WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS uq_account_linked_accounts_user_provider
    ON account.linked_accounts(user_id, auth_provider)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_account_linked_accounts_identifier
    ON account.linked_accounts(identifier)
    WHERE deleted_at IS NULL;
