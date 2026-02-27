-- Migration: 004_create_content.sql
-- Description: Create section_contents table for flexible key-value content storage
-- Created: 2024-01-01

-- Content type enum
CREATE TYPE content_type AS ENUM (
    'text',
    'html',
    'markdown',
    'image',
    'video',
    'link',
    'button',
    'color',
    'number',
    'boolean',
    'json',
    'file'
);

-- Section contents table (flexible key-value content per section)
CREATE TABLE IF NOT EXISTS section_contents (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    section_id      UUID NOT NULL REFERENCES page_sections(id) ON DELETE CASCADE,
    key             VARCHAR(255) NOT NULL,  -- e.g., 'title', 'subtitle', 'cta_text', 'cta_link', 'image'
    value           TEXT,                   -- text/html/markdown/url values
    value_json      JSONB,                  -- for complex structured values
    type            content_type NOT NULL DEFAULT 'text',
    label           VARCHAR(255),           -- human-readable label for admin UI
    description     TEXT,                   -- help text for admin UI
    placeholder     VARCHAR(255),           -- placeholder for admin UI input
    is_required     BOOLEAN NOT NULL DEFAULT false,
    sort_order      INTEGER NOT NULL DEFAULT 0,
    -- For image/file types
    alt_text        VARCHAR(255),
    width           INTEGER,
    height          INTEGER,
    -- For link/button types
    link_url        TEXT,
    link_target     VARCHAR(20) DEFAULT '_self',  -- _self, _blank
    -- Metadata for extensibility
    metadata        JSONB DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_section_contents_section_key ON section_contents(section_id, key);
CREATE INDEX idx_section_contents_section_id ON section_contents(section_id);
CREATE INDEX idx_section_contents_type ON section_contents(type);
CREATE INDEX idx_section_contents_sort ON section_contents(section_id, sort_order);

-- Apply trigger
CREATE TRIGGER update_section_contents_updated_at
    BEFORE UPDATE ON section_contents
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Insert default content for default page sections
DO $$
DECLARE
    v_hero_id UUID;
    v_features_id UUID;
    v_stats_id UUID;
    v_testimonials_id UUID;
    v_pricing_id UUID;
    v_faq_id UUID;
    v_cta_id UUID;
BEGIN
    -- Get section IDs
    SELECT ps.id INTO v_hero_id 
    FROM page_sections ps 
    JOIN pages p ON p.id = ps.page_id 
    WHERE p.slug = 'home' AND ps.identifier = 'hero' LIMIT 1;
    
    SELECT ps.id INTO v_features_id 
    FROM page_sections ps 
    JOIN pages p ON p.id = ps.page_id 
    WHERE p.slug = 'home' AND ps.identifier = 'features' LIMIT 1;
    
    SELECT ps.id INTO v_stats_id 
    FROM page_sections ps 
    JOIN pages p ON p.id = ps.page_id 
    WHERE p.slug = 'home' AND ps.identifier = 'stats' LIMIT 1;
    
    SELECT ps.id INTO v_cta_id 
    FROM page_sections ps 
    JOIN pages p ON p.id = ps.page_id 
    WHERE p.slug = 'home' AND ps.identifier = 'cta' LIMIT 1;

    -- Hero section content
    IF v_hero_id IS NOT NULL THEN
        INSERT INTO section_contents (section_id, key, value, type, label, sort_order) VALUES
        (v_hero_id, 'badge_text', '🚀 New Release v2.0', 'text', 'Badge Text', 1),
        (v_hero_id, 'title', 'Build Something Amazing', 'text', 'Main Title', 2),
        (v_hero_id, 'title_highlight', 'Amazing', 'text', 'Highlighted Word in Title', 3),
        (v_hero_id, 'subtitle', 'The most powerful platform to create, manage, and scale your digital presence. Start free, grow unlimited.', 'text', 'Subtitle', 4),
        (v_hero_id, 'cta_primary_text', 'Get Started Free', 'text', 'Primary CTA Text', 5),
        (v_hero_id, 'cta_primary_link', '#pricing', 'link', 'Primary CTA Link', 6),
        (v_hero_id, 'cta_secondary_text', 'Watch Demo', 'text', 'Secondary CTA Text', 7),
        (v_hero_id, 'cta_secondary_link', '#demo', 'link', 'Secondary CTA Link', 8),
        (v_hero_id, 'hero_image', '', 'image', 'Hero Image URL', 9),
        (v_hero_id, 'hero_image_alt', 'Product screenshot', 'text', 'Hero Image Alt Text', 10),
        (v_hero_id, 'social_proof_text', 'Trusted by 10,000+ teams worldwide', 'text', 'Social Proof Text', 11),
        (v_hero_id, 'video_url', '', 'video', 'Demo Video URL', 12)
        ON CONFLICT DO NOTHING;
    END IF;

    -- Features section content
    IF v_features_id IS NOT NULL THEN
        INSERT INTO section_contents (section_id, key, value, type, label, sort_order) VALUES
        (v_features_id, 'badge_text', 'Features', 'text', 'Badge Text', 1),
        (v_features_id, 'title', 'Everything you need to succeed', 'text', 'Section Title', 2),
        (v_features_id, 'subtitle', 'Powerful features designed to help you build faster and scale smarter.', 'text', 'Section Subtitle', 3)
        ON CONFLICT DO NOTHING;
    END IF;

    -- Stats section content
    IF v_stats_id IS NOT NULL THEN
        INSERT INTO section_contents (section_id, key, value, type, label, sort_order) VALUES
        (v_stats_id, 'stat_1_value', '10,000+', 'text', 'Stat 1 Value', 1),
        (v_stats_id, 'stat_1_label', 'Active Users', 'text', 'Stat 1 Label', 2),
        (v_stats_id, 'stat_2_value', '99.9%', 'text', 'Stat 2 Value', 3),
        (v_stats_id, 'stat_2_label', 'Uptime SLA', 'text', 'Stat 2 Label', 4),
        (v_stats_id, 'stat_3_value', '50+', 'text', 'Stat 3 Value', 5),
        (v_stats_id, 'stat_3_label', 'Countries', 'text', 'Stat 3 Label', 6),
        (v_stats_id, 'stat_4_value', '24/7', 'text', 'Stat 4 Value', 7),
        (v_stats_id, 'stat_4_label', 'Support', 'text', 'Stat 4 Label', 8)
        ON CONFLICT DO NOTHING;
    END IF;

    -- CTA section content
    IF v_cta_id IS NOT NULL THEN
        INSERT INTO section_contents (section_id, key, value, type, label, sort_order) VALUES
        (v_cta_id, 'title', 'Ready to get started?', 'text', 'CTA Title', 1),
        (v_cta_id, 'subtitle', 'Join thousands of teams already using our platform. No credit card required.', 'text', 'CTA Subtitle', 2),
        (v_cta_id, 'cta_primary_text', 'Start Free Trial', 'text', 'Primary CTA Text', 3),
        (v_cta_id, 'cta_primary_link', '/signup', 'link', 'Primary CTA Link', 4),
        (v_cta_id, 'cta_secondary_text', 'Contact Sales', 'text', 'Secondary CTA Text', 5),
        (v_cta_id, 'cta_secondary_link', '/contact', 'link', 'Secondary CTA Link', 6)
        ON CONFLICT DO NOTHING;
    END IF;
END $$;

-- ============================================================
-- ROLLBACK SCRIPT
-- ============================================================
-- DROP TABLE IF EXISTS section_contents CASCADE;
-- DROP TYPE IF EXISTS content_type;
