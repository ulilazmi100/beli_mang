-- Enable necessary extensions
CREATE EXTENSION IF NOT EXISTS "btree_gist";

-- Creating the users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(35),
    password VARCHAR(80),
    email VARCHAR(255),
    role VARCHAR(15),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- CREATE INDEX IF NOT EXISTS idx_users_role ON users (role);