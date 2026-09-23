-- Rollback: account_add_linked_accounts_table
-- Created: Tue Sep  9 01:49:11 +03 2025

-- Drop account.linked_accounts table if it exists
DROP TABLE IF EXISTS account.linked_accounts;
