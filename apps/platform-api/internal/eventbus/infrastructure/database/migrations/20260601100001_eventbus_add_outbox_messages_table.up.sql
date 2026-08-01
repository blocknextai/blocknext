-- Migration: add_outbox_messages_table
-- Created: Mon Jun  1 10:00:01 +03 2026

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'eventbus_outbox_message_status') THEN
        CREATE TYPE eventbus_outbox_message_status AS ENUM ('pending', 'processing', 'processed', 'dead');
    END IF;
END$$;

CREATE TABLE IF NOT EXISTS eventbus.outbox_messages (
    id                  UUID PRIMARY KEY,
    event_name          VARCHAR(255) NOT NULL,
    payload             JSONB NOT NULL,
    occurred_at         TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    status              eventbus_outbox_message_status NOT NULL DEFAULT 'pending',
    attempts            INTEGER NOT NULL DEFAULT 0,
    next_attempt_at     TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    locked_at           TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL,
    last_error          TEXT DEFAULT NULL,
    processed_at        TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS idx_eventbus_outbox_messages_due
    ON eventbus.outbox_messages(next_attempt_at, occurred_at)
    WHERE status = 'pending';

CREATE INDEX IF NOT EXISTS idx_eventbus_outbox_messages_locked
    ON eventbus.outbox_messages(locked_at)
    WHERE status = 'processing';
