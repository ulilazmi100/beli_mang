DROP INDEX IF EXISTS idx_users_role;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_username;

DROP TABLE IF EXISTS users;

-- Drop extensions
DROP EXTENSION IF EXISTS "btree_gist";