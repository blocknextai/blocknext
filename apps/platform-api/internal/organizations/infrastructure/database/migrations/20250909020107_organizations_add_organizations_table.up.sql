-- Migration: organizations_add_organizations_table
-- Created: Tue Sep  9 02:01:07 +03 2025

-- Create organizations.organizations table
CREATE TABLE IF NOT EXISTS organizations.organizations (
    id          UUID PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    description TEXT,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    updated_at  TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    deleted_at  TIMESTAMP WITHOUT TIME ZONE
);
