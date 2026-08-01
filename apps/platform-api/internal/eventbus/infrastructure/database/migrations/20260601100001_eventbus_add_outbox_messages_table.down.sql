-- Rollback: add_outbox_messages_table
-- Created: Mon Jun  1 10:00:01 +03 2026

DROP TABLE IF EXISTS eventbus.outbox_messages;
DROP TYPE IF EXISTS eventbus_outbox_message_status;
