-- Migration: account_add_socials_table
-- Created: Tue Sep  9 01:49:28 +03 2025

-- Create account.socials table
CREATE TABLE IF NOT EXISTS account.socials (
    id         UUID PRIMARY KEY,
    user_id    UUID NOT NULL,
    url        TEXT NOT NULL,
    platform   VARCHAR(255) NOT NULL,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    deleted_at TIMESTAMP WITHOUT TIME ZONE
);

-- Create index on user_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_account_socials_user_id ON account.socials(user_id);

-- Create index on platform for faster lookups
CREATE INDEX IF NOT EXISTS idx_account_socials_platform ON account.socials(platform);
