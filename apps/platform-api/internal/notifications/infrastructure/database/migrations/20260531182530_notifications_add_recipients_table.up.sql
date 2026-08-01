-- Migration: add_recipients_table
-- Created: Sun May 31 18:25:30 +03 2026

CREATE TABLE IF NOT EXISTS notifications.recipients (
    id              UUID PRIMARY KEY,
    notification_id UUID NOT NULL,
    user_id         UUID NOT NULL,
    read_at         TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL,
    seen_at         TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL,
    created_at      TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    updated_at      TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    deleted_at      TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_recipients_notification_user
    ON notifications.recipients(notification_id, user_id)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_recipients_inbox
    ON notifications.recipients(user_id, created_at DESC)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_recipients_unread
    ON notifications.recipients(user_id)
    WHERE deleted_at IS NULL AND read_at IS NULL;
