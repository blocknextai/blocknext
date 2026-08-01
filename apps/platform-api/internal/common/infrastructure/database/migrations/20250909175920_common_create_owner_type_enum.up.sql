-- Migration: common_create_owner_type_enum
-- Created: Tue Sep  9 17:59:20 +03 2025

-- Create enum type for OwnerType
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'owner_type') THEN
        CREATE TYPE owner_type AS ENUM ('organization', 'user');
    END IF;
END$$;
