-- Migration: add_inbox_entries_table
-- Created: Mon Jun  1 10:00:02 +03 2026

CREATE TABLE IF NOT EXISTS eventbus.inbox_entries (
    handler_key     VARCHAR(255) NOT NULL,
    event_id        UUID NOT NULL,
    processed_at    TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    PRIMARY KEY (handler_key, event_id)
);
