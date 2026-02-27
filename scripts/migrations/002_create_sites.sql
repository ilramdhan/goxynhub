-- Migration: 002_create_sites.sql
-- Description: Create sites and site_settings tables
-- Created: 2024-01-01

-- Sites table (supports multi-site)
CREATE TABLE IF NOT EXISTS sites (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(255) NOT NULL,
    slug            VARCHAR(255) NOT NULL,
    domain          VARCHAR(255),
    description     TEXT,
    logo_url        TEXT,
    favicon_url     TEXT,
    is_active       BOOLEAN NOT NULL DEFAULT true,
    metadata        JSONB DEFAULT '{}',
    created_by      UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_sites_slug_unique ON sites(slug) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_sites_domain_unique ON sites(domain) WHERE deleted_at IS NULL AND domain IS NOT NULL;
CREATE INDEX idx_sites_is_active ON sites(is_active) WHERE is_active = true;
CREATE INDEX idx_sites_deleted_at ON sites(deleted_at) WHERE deleted_at IS NULL;

-- Site settings table (key-value store for flexible settings)
CREATE TABLE IF NOT EXISTS site_settings (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    site_id     UUID NOT NULL REFERENCES sites(id) ON DELETE CASCADE,
    key         VARCHAR(255) NOT NULL,
    value       TEXT,
    value_json  JSONB,
    type        VARCHAR(50) NOT NULL DEFAULT 'string',  -- string, number, boolean, json, html
    group_name  VARCHAR(100) NOT NULL DEFAULT 'general', -- general, seo, social, analytics, appearance
    label       VARCHAR(255),
    description TEXT,
    is_public   BOOLEAN NOT NULL DEFAULT false,  -- whether this setting is exposed to public API
    sort_order  INTEGER NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_site_settings_site_key ON site_settings(site_id, key);
CREATE INDEX idx_site_settings_site_id ON site_settings(site_id);
CREATE INDEX idx_site_settings_group ON site_settings(site_id, group_name);
CREATE INDEX idx_site_settings_is_public ON site_settings(site_id, is_public) WHERE is_public = true;

-- Apply updated_at trigger
CREATE TRIGGER update_sites_updated_at
    BEFORE UPDATE ON sites
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_site_settings_updated_at
    BEFORE UPDATE ON site_settings
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Insert default site
INSERT INTO sites (name, slug, domain, description, is_active)
VALUES ('My Landing Page', 'default', 'localhost', 'Default landing page site', true)
ON CONFLICT DO NOTHING;

-- Insert default site settings
DO $$
DECLARE
    v_site_id UUID;
BEGIN
    SELECT id INTO v_site_id FROM sites WHERE slug = 'default' LIMIT 1;
    
    IF v_site_id IS NOT NULL THEN
        INSERT INTO site_settings (site_id, key, value, type, group_name, label, is_public) VALUES
        -- General settings
        (v_site_id, 'site_title', 'My Awesome Product', 'string', 'general', 'Site Title', true),
        (v_site_id, 'site_tagline', 'The best solution for your needs', 'string', 'general', 'Site Tagline', true),
        (v_site_id, 'site_description', 'A comprehensive solution that helps you achieve your goals faster.', 'string', 'general', 'Site Description', true),
        (v_site_id, 'contact_email', 'contact@example.com', 'string', 'general', 'Contact Email', true),
        (v_site_id, 'contact_phone', '+1 (555) 000-0000', 'string', 'general', 'Contact Phone', true),
        (v_site_id, 'contact_address', '123 Main St, City, Country', 'string', 'general', 'Contact Address', true),
        -- SEO settings
        (v_site_id, 'seo_title', 'My Awesome Product - Best Solution', 'string', 'seo', 'SEO Title', true),
        (v_site_id, 'seo_description', 'Discover the best solution for your needs. Fast, reliable, and affordable.', 'string', 'seo', 'SEO Description', true),
        (v_site_id, 'seo_keywords', 'product, solution, awesome, best', 'string', 'seo', 'SEO Keywords', true),
        (v_site_id, 'og_title', 'My Awesome Product', 'string', 'seo', 'OG Title', true),
        (v_site_id, 'og_description', 'The best solution for your needs', 'string', 'seo', 'OG Description', true),
        (v_site_id, 'og_image', '', 'string', 'seo', 'OG Image URL', true),
        (v_site_id, 'twitter_card', 'summary_large_image', 'string', 'seo', 'Twitter Card Type', true),
        (v_site_id, 'twitter_site', '@myproduct', 'string', 'seo', 'Twitter Site Handle', true),
        (v_site_id, 'canonical_url', '', 'string', 'seo', 'Canonical URL', true),
        -- Social settings
        (v_site_id, 'social_facebook', '', 'string', 'social', 'Facebook URL', true),
        (v_site_id, 'social_twitter', '', 'string', 'social', 'Twitter URL', true),
        (v_site_id, 'social_instagram', '', 'string', 'social', 'Instagram URL', true),
        (v_site_id, 'social_linkedin', '', 'string', 'social', 'LinkedIn URL', true),
        (v_site_id, 'social_youtube', '', 'string', 'social', 'YouTube URL', true),
        (v_site_id, 'social_github', '', 'string', 'social', 'GitHub URL', true),
        -- Analytics settings
        (v_site_id, 'google_analytics_id', '', 'string', 'analytics', 'Google Analytics ID', false),
        (v_site_id, 'google_tag_manager_id', '', 'string', 'analytics', 'Google Tag Manager ID', false),
        (v_site_id, 'facebook_pixel_id', '', 'string', 'analytics', 'Facebook Pixel ID', false),
        -- Appearance settings
        (v_site_id, 'primary_color', '#6366f1', 'string', 'appearance', 'Primary Color', true),
        (v_site_id, 'secondary_color', '#8b5cf6', 'string', 'appearance', 'Secondary Color', true),
        (v_site_id, 'font_family', 'Inter', 'string', 'appearance', 'Font Family', true),
        (v_site_id, 'custom_css', '', 'html', 'appearance', 'Custom CSS', false),
        (v_site_id, 'custom_js', '', 'html', 'appearance', 'Custom JavaScript', false)
        ON CONFLICT DO NOTHING;
    END IF;
END $$;

-- ============================================================
-- ROLLBACK SCRIPT
-- ============================================================
-- DROP TABLE IF EXISTS site_settings CASCADE;
-- DROP TABLE IF EXISTS sites CASCADE;
