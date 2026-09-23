-- Migration: workflows_add_workflows_table
-- Created: Tue Sep  9 01:55:38 +03 2025

-- Create workflows.workflows table
CREATE TABLE IF NOT EXISTS workflows.workflows (
    id                UUID PRIMARY KEY,
    organization_id   UUID NOT NULL,
    owner_id          UUID NOT NULL,
    title             VARCHAR(255) NOT NULL,
    description       TEXT,
    is_pinned         BOOLEAN NOT NULL DEFAULT FALSE,
    is_active         BOOLEAN NOT NULL DEFAULT TRUE,
    nodes             JSONB,
    edges             JSONB,
    created_at        TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    updated_at        TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    deleted_at        TIMESTAMP WITHOUT TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_workflows_workflows_organization_id
    ON workflows.workflows(organization_id);

CREATE INDEX IF NOT EXISTS idx_workflows_workflows_owner_id
    ON workflows.workflows(owner_id);

CREATE INDEX IF NOT EXISTS idx_workflows_workflows_org_pinned_updated
    ON workflows.workflows(organization_id, is_pinned DESC, updated_at DESC)
    WHERE deleted_at IS NULL;
