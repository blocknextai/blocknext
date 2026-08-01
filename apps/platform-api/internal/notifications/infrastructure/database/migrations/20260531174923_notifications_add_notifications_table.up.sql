-- Migration: add_notifications_table
-- Created: Sun May 31 17:49:23 +03 2026

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'notification_level') THEN
        CREATE TYPE notification_level AS ENUM ('info', 'success', 'warning', 'error');
    END IF;
END$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'notification_audience_type') THEN
        CREATE TYPE notification_audience_type AS ENUM ('user', 'organization');
    END IF;
END$$;

CREATE TABLE IF NOT EXISTS notifications.notifications (
    id                          UUID PRIMARY KEY,
    type                        VARCHAR(128) NOT NULL,
    level                       notification_level NOT NULL,
    audience_type               notification_audience_type NOT NULL,
    audience_id                 UUID NOT NULL,
    title                       VARCHAR(255) NOT NULL,
    body                        TEXT DEFAULT NULL,
    data                        JSONB DEFAULT NULL,
    action_url                  TEXT DEFAULT NULL,
    created_at                  TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    updated_at                  TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    deleted_at                  TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS idx_notifications_audience
    ON notifications.notifications(audience_type, audience_id)
    WHERE deleted_at IS NULL;
