-- Migration: 003_create_pages.sql
-- Description: Create pages and page_sections tables
-- Created: 2024-01-01

-- Page status enum
CREATE TYPE page_status AS ENUM ('draft', 'published', 'archived');

-- Pages table
CREATE TABLE IF NOT EXISTS pages (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    site_id         UUID NOT NULL REFERENCES sites(id) ON DELETE CASCADE,
    title           VARCHAR(255) NOT NULL,
    slug            VARCHAR(255) NOT NULL,
    description     TEXT,
    status          page_status NOT NULL DEFAULT 'draft',
    is_homepage     BOOLEAN NOT NULL DEFAULT false,
    -- SEO fields
    seo_title       VARCHAR(255),
    seo_description VARCHAR(500),
    seo_keywords    TEXT,
    -- Open Graph fields
    og_title        VARCHAR(255),
    og_description  VARCHAR(500),
    og_image        TEXT,
    og_type         VARCHAR(50) DEFAULT 'website',
    -- Twitter Card fields
    twitter_title   VARCHAR(255),
    twitter_description VARCHAR(500),
    twitter_image   TEXT,
    twitter_card    VARCHAR(50) DEFAULT 'summary_large_image',
    -- Schema.org / Structured Data
    schema_markup   JSONB,
    -- Custom head tags
    custom_head     TEXT,
    -- Canonical URL
    canonical_url   TEXT,
    -- Robots meta
    robots_meta     VARCHAR(100) DEFAULT 'index, follow',
    -- Template/layout
    template        VARCHAR(100) DEFAULT 'default',
    -- Sort order for menu
    sort_order      INTEGER NOT NULL DEFAULT 0,
    -- Publish scheduling
    published_at    TIMESTAMPTZ,
    -- Metadata for extensibility
    metadata        JSONB DEFAULT '{}',
    created_by      UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_by      UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_pages_slug_site ON pages(slug, site_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_pages_site_id ON pages(site_id);
CREATE INDEX idx_pages_status ON pages(status);
CREATE INDEX idx_pages_is_homepage ON pages(site_id, is_homepage) WHERE is_homepage = true;
CREATE INDEX idx_pages_deleted_at ON pages(deleted_at) WHERE deleted_at IS NULL;

-- Section type enum
CREATE TYPE section_type AS ENUM (
    'hero',
    'features',
    'pricing',
    'testimonials',
    'faq',
    'cta',
    'about',
    'team',
    'gallery',
    'stats',
    'logos',
    'newsletter',
    'contact',
    'video',
    'custom',
    'html'
);

-- Page sections table
CREATE TABLE IF NOT EXISTS page_sections (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    page_id         UUID NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
    name            VARCHAR(255) NOT NULL,
    type            section_type NOT NULL DEFAULT 'custom',
    identifier      VARCHAR(100),  -- unique identifier within page for frontend targeting
    is_visible      BOOLEAN NOT NULL DEFAULT true,
    sort_order      INTEGER NOT NULL DEFAULT 0,
    -- Background settings
    bg_color        VARCHAR(50),
    bg_image        TEXT,
    bg_video        TEXT,
    bg_overlay      BOOLEAN DEFAULT false,
    bg_overlay_color VARCHAR(50),
    bg_overlay_opacity DECIMAL(3,2) DEFAULT 0.5,
    -- Layout settings
    layout          VARCHAR(50) DEFAULT 'default',  -- default, full-width, contained, split
    padding_top     VARCHAR(20) DEFAULT 'md',       -- none, sm, md, lg, xl
    padding_bottom  VARCHAR(20) DEFAULT 'md',
    -- Animation
    animation       VARCHAR(50),  -- fade-in, slide-up, etc.
    -- Custom CSS class
    css_class       VARCHAR(255),
    -- Custom CSS
    custom_css      TEXT,
    -- Metadata for extensibility
    metadata        JSONB DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_page_sections_page_id ON page_sections(page_id);
CREATE INDEX idx_page_sections_type ON page_sections(type);
CREATE INDEX idx_page_sections_sort ON page_sections(page_id, sort_order);
CREATE UNIQUE INDEX idx_page_sections_identifier ON page_sections(page_id, identifier) WHERE identifier IS NOT NULL;

-- Apply triggers
CREATE TRIGGER update_pages_updated_at
    BEFORE UPDATE ON pages
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_page_sections_updated_at
    BEFORE UPDATE ON page_sections
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Insert default homepage
DO $$
DECLARE
    v_site_id UUID;
    v_page_id UUID;
BEGIN
    SELECT id INTO v_site_id FROM sites WHERE slug = 'default' LIMIT 1;
    
    IF v_site_id IS NOT NULL THEN
        INSERT INTO pages (
            site_id, title, slug, description, status, is_homepage,
            seo_title, seo_description, og_title, og_description
        ) VALUES (
            v_site_id,
            'Home',
            'home',
            'Main landing page',
            'published',
            true,
            'My Awesome Product - Best Solution',
            'Discover the best solution for your needs. Fast, reliable, and affordable.',
            'My Awesome Product',
            'The best solution for your needs'
        ) RETURNING id INTO v_page_id;

        -- Insert default sections
        INSERT INTO page_sections (page_id, name, type, identifier, sort_order) VALUES
        (v_page_id, 'Hero Section', 'hero', 'hero', 1),
        (v_page_id, 'Features Section', 'features', 'features', 2),
        (v_page_id, 'Stats Section', 'stats', 'stats', 3),
        (v_page_id, 'Testimonials Section', 'testimonials', 'testimonials', 4),
        (v_page_id, 'Pricing Section', 'pricing', 'pricing', 5),
        (v_page_id, 'FAQ Section', 'faq', 'faq', 6),
        (v_page_id, 'CTA Section', 'cta', 'cta', 7);
    END IF;
END $$;

-- ============================================================
-- ROLLBACK SCRIPT
-- ============================================================
-- DROP TABLE IF EXISTS page_sections CASCADE;
-- DROP TABLE IF EXISTS pages CASCADE;
-- DROP TYPE IF EXISTS section_type;
-- DROP TYPE IF EXISTS page_status;
