-- Migration: credentials_add_credentials_table
-- Created: Thu Jan  8 01:48:40 +03 2026

-- Create enum type for credential_source_type if it does not exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'credential_source_type') THEN
        CREATE TYPE credential_source_type AS ENUM ('owner', 'platform');
    END IF;
END$$;

-- Create credentials.credentials table
CREATE TABLE IF NOT EXISTS credentials.credentials (
    id          UUID PRIMARY KEY,
    owner_type  owner_type NOT NULL,
    owner_id    UUID NOT NULL,
    source_type credential_source_type NOT NULL DEFAULT 'owner',
    key         VARCHAR(255) NOT NULL,
    title       VARCHAR(255) NOT NULL,
    data        TEXT NOT NULL,
    created_at  TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    updated_at  TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    deleted_at  TIMESTAMP WITHOUT TIME ZONE
);

-- Create index on owner_type and owner_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_credentials_owner ON credentials.credentials(owner_type, owner_id);

-- Create index on source_type for filtering by source
CREATE INDEX IF NOT EXISTS idx_credentials_source_type ON credentials.credentials(source_type);

-- Create index on deleted_at for soft delete queries
CREATE INDEX IF NOT EXISTS idx_credentials_deleted_at ON credentials.credentials(deleted_at);
