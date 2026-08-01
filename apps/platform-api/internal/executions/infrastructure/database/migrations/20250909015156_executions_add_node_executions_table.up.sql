-- Migration: executions_add_node_executions_table
-- Created: Tue Sep  9 01:51:56 +03 2025

-- Create executions.node_executions table
CREATE TABLE IF NOT EXISTS executions.node_executions (
    id                      UUID PRIMARY KEY,
    task_id                 UUID NOT NULL,
    node_type               VARCHAR(100) NOT NULL,
    node_id                 VARCHAR(100) NOT NULL,
    status                  VARCHAR(50) NOT NULL,
    inputs                  JSONB,
    outputs                 JSONB,
    function_calling_outputs JSONB,
    error_message           TEXT,
    started_at              TIMESTAMP WITHOUT TIME ZONE,
    completed_at            TIMESTAMP WITHOUT TIME ZONE,
    created_at              TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    updated_at              TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    deleted_at              TIMESTAMP WITHOUT TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_node_executions_status ON executions.node_executions(status);

CREATE INDEX IF NOT EXISTS idx_node_executions_task_id ON executions.node_executions(task_id) WHERE deleted_at IS NULL;
