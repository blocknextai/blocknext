-- Migration: add_user_preferences_table
-- Created: Thu Mar  5 00:42:53 +03 2026

CREATE TABLE IF NOT EXISTS account.user_preferences (
    id          UUID PRIMARY KEY,
    user_id     UUID NOT NULL UNIQUE,
    theme_mode  VARCHAR(10) NOT NULL DEFAULT 'system',
    theme_color VARCHAR(20) NOT NULL DEFAULT 'default',
    language    VARCHAR(10) NOT NULL DEFAULT 'en',
    created_at  TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    updated_at  TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    deleted_at  TIMESTAMP WITHOUT TIME ZONE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_account_user_preferences_user_id
    ON account.user_preferences(user_id) WHERE deleted_at IS NULL;
