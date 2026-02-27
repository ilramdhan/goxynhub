-- Migration: 005_create_media.sql
-- Description: Create media library table
-- Created: 2024-01-01

-- Media type enum
CREATE TYPE media_type AS ENUM ('image', 'video', 'document', 'audio', 'other');

-- Media table
CREATE TABLE IF NOT EXISTS media (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    site_id         UUID NOT NULL REFERENCES sites(id) ON DELETE CASCADE,
    name            VARCHAR(255) NOT NULL,
    original_name   VARCHAR(255) NOT NULL,
    file_path       TEXT NOT NULL,          -- path in storage bucket
    public_url      TEXT NOT NULL,          -- public accessible URL
    thumbnail_url   TEXT,                   -- thumbnail for images/videos
    type            media_type NOT NULL DEFAULT 'image',
    mime_type       VARCHAR(100) NOT NULL,
    file_size       BIGINT NOT NULL,        -- in bytes
    width           INTEGER,               -- for images/videos
    height          INTEGER,               -- for images/videos
    duration        INTEGER,               -- for videos/audio (in seconds)
    alt_text        VARCHAR(255),
    caption         TEXT,
    tags            TEXT[],                -- array of tags for filtering
    folder          VARCHAR(255) DEFAULT '/',  -- virtual folder path
    is_used         BOOLEAN NOT NULL DEFAULT false,  -- whether media is referenced anywhere
    metadata        JSONB DEFAULT '{}',
    uploaded_by     UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_media_site_id ON media(site_id);
CREATE INDEX idx_media_type ON media(type);
CREATE INDEX idx_media_folder ON media(site_id, folder);
CREATE INDEX idx_media_tags ON media USING GIN(tags);
CREATE INDEX idx_media_deleted_at ON media(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_media_uploaded_by ON media(uploaded_by);

-- Apply trigger
CREATE TRIGGER update_media_updated_at
    BEFORE UPDATE ON media
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ============================================================
-- ROLLBACK SCRIPT
-- ============================================================
-- DROP TABLE IF EXISTS media CASCADE;
-- DROP TYPE IF EXISTS media_type;
