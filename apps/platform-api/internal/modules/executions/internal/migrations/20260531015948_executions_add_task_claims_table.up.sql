-- Migration: add_task_claims_table
-- Created: Sun May 31 01:59:48 +03 2026

CREATE TABLE IF NOT EXISTS executions.task_claims (
    id                UUID PRIMARY KEY,
    task_execution_id UUID NOT NULL UNIQUE,
    claimed_by        VARCHAR(255),
    claimed_at        TIMESTAMP WITHOUT TIME ZONE,
    retry_count       INTEGER NOT NULL DEFAULT 0,
    created_at        TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    updated_at        TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    deleted_at        TIMESTAMP WITHOUT TIME ZONE
    );

CREATE INDEX IF NOT EXISTS idx_task_claims_claimed_at
    ON executions.task_claims(claimed_at)
    WHERE claimed_at IS NOT NULL;
