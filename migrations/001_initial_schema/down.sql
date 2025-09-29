-- migrations/001_initial_schema/down.sql

-- Drop triggers first (they depend on functions)
DROP TRIGGER IF EXISTS hash_password_trigger ON accounts;
DROP TRIGGER IF EXISTS update_accounts_updated_at ON accounts;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop functions
DROP FUNCTION IF EXISTS hash_password();
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indices
DROP INDEX IF EXISTS idx_accounts_expires;
DROP INDEX IF EXISTS idx_accounts_user_provider;
DROP INDEX IF EXISTS idx_users_email;

-- Drop tables (in reverse dependency order)
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS users;

-- Drop enum types
DROP TYPE IF EXISTS providers;

-- Drop extensions (be careful with this in production)
-- DROP EXTENSION IF EXISTS "pgcrypto";
-- DROP EXTENSION IF EXISTS "uuid-ossp";