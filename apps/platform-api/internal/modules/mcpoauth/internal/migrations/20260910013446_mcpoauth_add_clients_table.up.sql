-- Migration: add_clients_table
-- Created: Wed Sep 10 01:34:45 +03 2026

CREATE TABLE IF NOT EXISTS mcpoauth.clients (
    id UUID                     PRIMARY KEY,
    client_id                   VARCHAR(512) NOT NULL,
    client_secret_hash          VARCHAR(255) DEFAULT NULL,
    name                        VARCHAR(255) NOT NULL,
    redirect_uris               TEXT[] NOT NULL DEFAULT '{}',
    grant_types                 TEXT[] NOT NULL DEFAULT '{}',
    response_types              TEXT[] NOT NULL DEFAULT '{}',
    scopes                      TEXT[] NOT NULL DEFAULT '{}',
    token_endpoint_auth_method  VARCHAR(50) NOT NULL,
    logo_uri                    TEXT DEFAULT NULL,
    client_uri                  TEXT DEFAULT NULL,
    created_at                  TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    updated_at                  TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC') NOT NULL,
    deleted_at                  TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_mcpoauth_clients_client_id
    ON mcpoauth.clients(client_id)
    WHERE deleted_at IS NULL;
