-- Migration: organizations_add_users_table
-- Created: Tue Sep  9 02:01:28 +03 2025

-- Create organizations.users table
CREATE TABLE IF NOT EXISTS organizations.users (
    id              UUID PRIMARY KEY,
    organization_id UUID NOT NULL,
    user_id         UUID NOT NULL,
    role            VARCHAR(255) NOT NULL,
    alias           VARCHAR(255),
    created_at      TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    updated_at      TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    deleted_at      TIMESTAMP WITHOUT TIME ZONE
);

-- Add indexes for faster lookup on organization and user
CREATE INDEX IF NOT EXISTS idx_organizations_users_organization_id ON organizations.users(organization_id);
CREATE INDEX IF NOT EXISTS idx_organizations_users_user_id ON organizations.users(user_id);

CREATE UNIQUE INDEX IF NOT EXISTS uq_organizations_users_org_user
    ON organizations.users(organization_id, user_id)
    WHERE deleted_at IS NULL;
