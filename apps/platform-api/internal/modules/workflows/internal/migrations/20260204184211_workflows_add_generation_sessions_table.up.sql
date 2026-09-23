-- Migration: add_generation_sessions_table
-- Created: Wed Feb  4 18:42:11 +03 2026

-- Create generation sessions table
CREATE TABLE IF NOT EXISTS workflows.generation_sessions (
    id              UUID PRIMARY KEY,
    organization_id UUID,
    user_id         UUID,
    title           VARCHAR(255) NOT NULL,
    created_at      TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    updated_at      TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    deleted_at      TIMESTAMP WITHOUT TIME ZONE
);

-- Create index
CREATE INDEX IF NOT EXISTS idx_generation_sessions_user_id ON workflows.generation_sessions(user_id) WHERE deleted_at IS NULL;
