-- Rollback: credentials_add_credentials_table
-- Created: Thu Jan  8 01:48:40 +03 2026

-- Drop credentials.credentials table
DROP TABLE IF EXISTS credentials.credentials;

-- Drop enum type (owner_type is from common module, don't drop it)
DROP TYPE IF EXISTS credential_source_type;
