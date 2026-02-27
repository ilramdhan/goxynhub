-- Migration: 008_create_audit.sql
-- Description: Create audit logs table for tracking admin actions
-- Created: 2024-01-01

-- Audit action enum
CREATE TYPE audit_action AS ENUM (
    'create',
    'update',
    'delete',
    'restore',
    'login',
    'logout',
    'login_failed',
    'password_change',
    'publish',
    'unpublish',
    'upload',
    'download'
);

-- Audit logs table
CREATE TABLE IF NOT EXISTS audit_logs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID REFERENCES users(id) ON DELETE SET NULL,
    user_email      VARCHAR(255),           -- denormalized for historical record
    user_role       VARCHAR(50),            -- denormalized for historical record
    action          audit_action NOT NULL,
    resource_type   VARCHAR(100) NOT NULL,  -- e.g., 'page', 'section', 'user', 'media'
    resource_id     UUID,
    resource_name   VARCHAR(255),           -- human-readable name for display
    old_values      JSONB,                  -- previous state (for updates/deletes)
    new_values      JSONB,                  -- new state (for creates/updates)
    ip_address      INET,
    user_agent      TEXT,
    request_id      VARCHAR(100),           -- for correlating with request logs
    site_id         UUID REFERENCES sites(id) ON DELETE SET NULL,
    metadata        JSONB DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Audit logs are append-only, no updates needed
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_resource ON audit_logs(resource_type, resource_id);
CREATE INDEX idx_audit_logs_site_id ON audit_logs(site_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at DESC);
CREATE INDEX idx_audit_logs_user_created ON audit_logs(user_id, created_at DESC);

-- Migrations tracking table (for managing migration state)
CREATE TABLE IF NOT EXISTS schema_migrations (
    version     VARCHAR(255) PRIMARY KEY,
    applied_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    description TEXT
);

-- Record this migration
INSERT INTO schema_migrations (version, description) VALUES
('001', 'Create users and auth tables'),
('002', 'Create sites and site settings'),
('003', 'Create pages and sections'),
('004', 'Create section contents'),
('005', 'Create media library'),
('006', 'Create navigation'),
('007', 'Create component tables'),
('008', 'Create audit logs')
ON CONFLICT DO NOTHING;

-- ============================================================
-- ROLLBACK SCRIPT
-- ============================================================
-- DROP TABLE IF EXISTS schema_migrations CASCADE;
-- DROP TABLE IF EXISTS audit_logs CASCADE;
-- DROP TYPE IF EXISTS audit_action;
