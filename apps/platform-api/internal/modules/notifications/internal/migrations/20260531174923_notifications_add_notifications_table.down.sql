-- Rollback: add_notifications_table
-- Created: Sun May 31 17:49:23 +03 2026

DROP TABLE IF EXISTS notifications.notifications;

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'notification_audience_type') THEN
        DROP TYPE notification_audience_type;
    END IF;
END$$;

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'notification_level') THEN
        DROP TYPE notification_level;
    END IF;
END$$;
