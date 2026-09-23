-- Migration: add_generation_messages_table
-- Created: Wed Feb 4 18:42:15 +03 2026

-- Create generation messages table
CREATE TABLE IF NOT EXISTS workflows.generation_messages (
    id         UUID PRIMARY KEY,
    session_id UUID NOT NULL,
    role       VARCHAR(50) NOT NULL,
    content    TEXT NOT NULL DEFAULT '',
    metadata   JSONB,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    deleted_at TIMESTAMP WITHOUT TIME ZONE
);

-- Create index
CREATE INDEX IF NOT EXISTS idx_generation_messages_session_id ON workflows.generation_messages(session_id) WHERE deleted_at IS NULL;
