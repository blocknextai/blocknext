-- Migration: common_create_execution_context_enum
-- Created: Sun Sep 21 19:41:04 +03 2025

-- Create enum type for ExecutionContext
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'execution_context') THEN
        CREATE TYPE execution_context AS ENUM ('workflow', 'library');
    END IF;
END$$;
