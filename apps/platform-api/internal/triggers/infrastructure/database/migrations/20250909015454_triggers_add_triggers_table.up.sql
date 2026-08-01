-- Migration: triggers_add_triggers_table
-- Created: Tue Sep  9 01:54:54 +03 2025

-- Create triggers.triggers table
CREATE TABLE IF NOT EXISTS triggers.triggers (
    id                   UUID PRIMARY KEY,
    organization_id      UUID,
    triggered_by_user_id UUID,
    execution_context    execution_context NOT NULL,
    context_item_id      UUID NOT NULL,
    type                 VARCHAR(50) NOT NULL,
    cron_pattern         VARCHAR(100),
    timezone             VARCHAR(50) DEFAULT 'UTC',
    is_active            BOOLEAN NOT NULL DEFAULT TRUE,
    webhook_token_hash   VARCHAR(255) DEFAULT NULL,
    webhook_secret       TEXT DEFAULT NULL,
    runtime_config       JSONB,
    created_at           TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    updated_at           TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    deleted_at           TIMESTAMP WITHOUT TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_triggers_context_item_id
    ON triggers.triggers(context_item_id);

CREATE INDEX IF NOT EXISTS idx_triggers_execution_context
    ON triggers.triggers(execution_context);

CREATE UNIQUE INDEX IF NOT EXISTS uq_triggers_webhook_token_hash
    ON triggers.triggers(webhook_token_hash)
    WHERE webhook_token_hash IS NOT NULL AND deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_triggers_organization_id
    ON triggers.triggers(organization_id)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_triggers_is_active
    ON triggers.triggers(is_active)
    WHERE deleted_at IS NULL;
