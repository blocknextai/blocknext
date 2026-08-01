-- Rollback: common_create_owner_type_enum
-- Created: Tue Sep  9 17:59:20 +03 2025

-- Drop enum type
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'owner_type') THEN
        DROP TYPE owner_type;
    END IF;
END$$;
