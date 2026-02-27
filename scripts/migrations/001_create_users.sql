-- Migration: 001_create_users.sql
-- Description: Create users table with roles and auth support
-- Created: 2024-01-01

-- Enable UUID extension (already enabled in Supabase)
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create enum for user roles
CREATE TYPE user_role AS ENUM ('super_admin', 'admin', 'editor');

-- Create enum for user status
CREATE TYPE user_status AS ENUM ('active', 'inactive', 'suspended');

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email           VARCHAR(255) NOT NULL,
    password_hash   VARCHAR(255) NOT NULL,
    full_name       VARCHAR(255) NOT NULL,
    avatar_url      TEXT,
    role            user_role NOT NULL DEFAULT 'editor',
    status          user_status NOT NULL DEFAULT 'active',
    last_login_at   TIMESTAMPTZ,
    last_login_ip   INET,
    failed_attempts INTEGER NOT NULL DEFAULT 0,
    locked_until    TIMESTAMPTZ,
    email_verified  BOOLEAN NOT NULL DEFAULT false,
    metadata        JSONB DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

-- Unique constraint on email (only for non-deleted users)
CREATE UNIQUE INDEX idx_users_email_unique ON users(email) WHERE deleted_at IS NULL;

-- Index for common queries
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_deleted_at ON users(deleted_at) WHERE deleted_at IS NULL;

-- Refresh tokens table
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash  VARCHAR(255) NOT NULL,
    expires_at  TIMESTAMPTZ NOT NULL,
    is_revoked  BOOLEAN NOT NULL DEFAULT false,
    ip_address  INET,
    user_agent  TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at  TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_refresh_tokens_hash ON refresh_tokens(token_hash);
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
CREATE INDEX idx_refresh_tokens_is_revoked ON refresh_tokens(is_revoked) WHERE is_revoked = false;

-- Password reset tokens table
CREATE TABLE IF NOT EXISTS password_reset_tokens (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash  VARCHAR(255) NOT NULL,
    expires_at  TIMESTAMPTZ NOT NULL,
    is_used     BOOLEAN NOT NULL DEFAULT false,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    used_at     TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_password_reset_tokens_hash ON password_reset_tokens(token_hash);
CREATE INDEX idx_password_reset_tokens_user_id ON password_reset_tokens(user_id);

-- Auto-update updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply trigger to users table
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Insert default super admin
-- Password: Admin@123456 (CHANGE IMMEDIATELY IN PRODUCTION!)
-- This hash was generated with bcrypt cost 10 for: Admin@123456
-- To regenerate: SELECT crypt('Admin@123456', gen_salt('bf', 10));
INSERT INTO users (
    email,
    password_hash,
    full_name,
    role,
    status,
    email_verified
) VALUES (
    'admin@example.com',
    crypt('Admin@123456', gen_salt('bf', 12)),
    'Super Admin',
    'super_admin',
    'active',
    true
) ON CONFLICT (email) DO UPDATE SET
    password_hash = EXCLUDED.password_hash,
    role = EXCLUDED.role,
    status = EXCLUDED.status;

-- ============================================================
-- ROLLBACK SCRIPT (run to undo this migration)
-- ============================================================
-- DROP TABLE IF EXISTS password_reset_tokens CASCADE;
-- DROP TABLE IF EXISTS refresh_tokens CASCADE;
-- DROP TABLE IF EXISTS users CASCADE;
-- DROP TYPE IF EXISTS user_status;
-- DROP TYPE IF EXISTS user_role;
-- DROP FUNCTION IF EXISTS update_updated_at_column();
