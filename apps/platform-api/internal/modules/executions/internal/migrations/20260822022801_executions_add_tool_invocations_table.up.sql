-- Migration: add_tool_invocations_table
-- Created: Sat Aug 22 02:28:01 +03 2026

CREATE TABLE IF NOT EXISTS executions.tool_invocations (
    id              UUID PRIMARY KEY,
    organization_id UUID NOT NULL,
    api_key_id      UUID,
    source          VARCHAR(20) NOT NULL,
    tool_id         VARCHAR(255) NOT NULL,
    status          VARCHAR(20) NOT NULL,
    parameters      JSONB,
    credentials     JSONB,
    outputs         JSONB,
    error_message   TEXT,
    started_at      TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    completed_at    TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    created_at      TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    updated_at      TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    deleted_at      TIMESTAMP WITHOUT TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_tool_invocations_org_started
    ON executions.tool_invocations(organization_id, started_at DESC)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_tool_invocations_tool_id
    ON executions.tool_invocations(tool_id);

CREATE INDEX IF NOT EXISTS idx_tool_invocations_api_key_id
    ON executions.tool_invocations(api_key_id);

CREATE INDEX IF NOT EXISTS idx_tool_invocations_deleted_at
    ON executions.tool_invocations(deleted_at);
