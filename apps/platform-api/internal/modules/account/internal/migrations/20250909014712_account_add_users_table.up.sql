-- Migration: account_add_users_table
-- Created: Tue Sep  9 01:47:12 +03 2025

-- Create account.users table
CREATE TABLE IF NOT EXISTS account.users (
    id          UUID PRIMARY KEY,
    role        VARCHAR(50) NOT NULL,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    is_banned   BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    updated_at  TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    deleted_at  TIMESTAMP WITHOUT TIME ZONE
);
