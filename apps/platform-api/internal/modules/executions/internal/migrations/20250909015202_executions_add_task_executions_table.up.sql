-- Migration: executions_add_task_executions_table
-- Created: Tue Sep  9 01:52:02 +03 2025

-- Create executions.task_executions table
CREATE TABLE IF NOT EXISTS executions.task_executions (
    id                   UUID PRIMARY KEY,
    organization_id      UUID NOT NULL,
    triggered_by_user_id UUID,
    flow_trigger_id      UUID,
    execution_context    execution_context NOT NULL,
    context_item_id      UUID NOT NULL,
    status               VARCHAR(50) NOT NULL,
    execution_type       VARCHAR(20) NOT NULL DEFAULT 'manual',
    error_message        TEXT,
    nodes                JSONB,
    edges                JSONB,
    started_at           TIMESTAMP WITHOUT TIME ZONE,
    completed_at         TIMESTAMP WITHOUT TIME ZONE,
    created_at           TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    updated_at           TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    deleted_at           TIMESTAMP WITHOUT TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_task_executions_status
    ON executions.task_executions(status);

CREATE INDEX IF NOT EXISTS idx_task_executions_organization_id
    ON executions.task_executions(organization_id);

CREATE INDEX IF NOT EXISTS idx_task_executions_triggered_by_user_id
    ON executions.task_executions(triggered_by_user_id);

CREATE INDEX IF NOT EXISTS idx_task_executions_execution_context
    ON executions.task_executions(execution_context);

CREATE INDEX IF NOT EXISTS idx_task_executions_context_item_id
    ON executions.task_executions(context_item_id);

CREATE INDEX IF NOT EXISTS idx_task_executions_execution_type
    ON executions.task_executions(execution_type);

CREATE INDEX IF NOT EXISTS idx_task_executions_org_updated
    ON executions.task_executions(organization_id, updated_at DESC)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_task_executions_context_and_item
    ON executions.task_executions(execution_context, context_item_id)
    WHERE deleted_at IS NULL;
