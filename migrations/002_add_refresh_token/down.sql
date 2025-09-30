-- migrations/002_add_refresh_token/down.sql
-- Created at: 2025-09-30 10:42:57

-- Drop trigger
DROP TRIGGER IF EXISTS update_refresh_tokens_updated_at ON refresh_tokens;

-- Drop indexes
DROP INDEX IF EXISTS idx_refresh_tokens_is_active;
DROP INDEX IF EXISTS idx_refresh_tokens_expires_at;
DROP INDEX IF EXISTS idx_refresh_tokens_user_id;

-- Drop table
DROP TABLE IF EXISTS refresh_tokens;
