-- Migration: 007_create_components.sql
-- Description: Create reusable content component tables (features, testimonials, pricing, FAQs, CTAs, social links)
-- Created: 2024-01-01

-- Features table
CREATE TABLE IF NOT EXISTS features (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    site_id         UUID NOT NULL REFERENCES sites(id) ON DELETE CASCADE,
    section_id      UUID REFERENCES page_sections(id) ON DELETE SET NULL,
    title           VARCHAR(255) NOT NULL,
    description     TEXT,
    icon            VARCHAR(100),           -- icon name (e.g., 'zap', 'shield', 'star')
    icon_color      VARCHAR(50),
    image_url       TEXT,
    image_alt       VARCHAR(255),
    link_url        TEXT,
    link_text       VARCHAR(100),
    is_active       BOOLEAN NOT NULL DEFAULT true,
    sort_order      INTEGER NOT NULL DEFAULT 0,
    metadata        JSONB DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_features_site_id ON features(site_id);
CREATE INDEX idx_features_section_id ON features(section_id);
CREATE INDEX idx_features_sort ON features(site_id, sort_order);
CREATE INDEX idx_features_deleted_at ON features(deleted_at) WHERE deleted_at IS NULL;

-- Testimonials table
CREATE TABLE IF NOT EXISTS testimonials (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    site_id         UUID NOT NULL REFERENCES sites(id) ON DELETE CASCADE,
    section_id      UUID REFERENCES page_sections(id) ON DELETE SET NULL,
    author_name     VARCHAR(255) NOT NULL,
    author_title    VARCHAR(255),           -- e.g., 'CEO at Company'
    author_company  VARCHAR(255),
    author_avatar   TEXT,
    content         TEXT NOT NULL,
    rating          SMALLINT CHECK (rating >= 1 AND rating <= 5),
    source          VARCHAR(100),           -- e.g., 'Twitter', 'G2', 'Capterra'
    source_url      TEXT,
    is_featured     BOOLEAN NOT NULL DEFAULT false,
    is_active       BOOLEAN NOT NULL DEFAULT true,
    sort_order      INTEGER NOT NULL DEFAULT 0,
    metadata        JSONB DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_testimonials_site_id ON testimonials(site_id);
CREATE INDEX idx_testimonials_section_id ON testimonials(section_id);
CREATE INDEX idx_testimonials_is_featured ON testimonials(site_id, is_featured) WHERE is_featured = true;
CREATE INDEX idx_testimonials_deleted_at ON testimonials(deleted_at) WHERE deleted_at IS NULL;

-- Pricing plans table
CREATE TABLE IF NOT EXISTS pricing_plans (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    site_id         UUID NOT NULL REFERENCES sites(id) ON DELETE CASCADE,
    section_id      UUID REFERENCES page_sections(id) ON DELETE SET NULL,
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    price_monthly   DECIMAL(10,2),
    price_yearly    DECIMAL(10,2),
    currency        VARCHAR(10) DEFAULT 'USD',
    price_label     VARCHAR(100),           -- e.g., 'per user/month', 'one-time'
    is_popular      BOOLEAN NOT NULL DEFAULT false,
    is_custom       BOOLEAN NOT NULL DEFAULT false,  -- for 'Contact us' plans
    badge_text      VARCHAR(100),           -- e.g., 'Most Popular', 'Best Value'
    cta_text        VARCHAR(100) DEFAULT 'Get Started',
    cta_link        TEXT,
    features        JSONB DEFAULT '[]',     -- array of feature strings
    features_excluded JSONB DEFAULT '[]',  -- features NOT included (shown with X)
    is_active       BOOLEAN NOT NULL DEFAULT true,
    sort_order      INTEGER NOT NULL DEFAULT 0,
    metadata        JSONB DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_pricing_plans_site_id ON pricing_plans(site_id);
CREATE INDEX idx_pricing_plans_section_id ON pricing_plans(section_id);
CREATE INDEX idx_pricing_plans_deleted_at ON pricing_plans(deleted_at) WHERE deleted_at IS NULL;

-- FAQs table
CREATE TABLE IF NOT EXISTS faqs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    site_id         UUID NOT NULL REFERENCES sites(id) ON DELETE CASCADE,
    section_id      UUID REFERENCES page_sections(id) ON DELETE SET NULL,
    question        TEXT NOT NULL,
    answer          TEXT NOT NULL,
    category        VARCHAR(100),
    is_active       BOOLEAN NOT NULL DEFAULT true,
    sort_order      INTEGER NOT NULL DEFAULT 0,
    metadata        JSONB DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_faqs_site_id ON faqs(site_id);
CREATE INDEX idx_faqs_section_id ON faqs(section_id);
CREATE INDEX idx_faqs_category ON faqs(site_id, category);
CREATE INDEX idx_faqs_deleted_at ON faqs(deleted_at) WHERE deleted_at IS NULL;

-- CTA buttons table (reusable CTAs)
CREATE TABLE IF NOT EXISTS cta_buttons (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    site_id         UUID NOT NULL REFERENCES sites(id) ON DELETE CASCADE,
    section_id      UUID REFERENCES page_sections(id) ON DELETE SET NULL,
    label           VARCHAR(255) NOT NULL,
    url             TEXT NOT NULL,
    target          VARCHAR(20) DEFAULT '_self',
    variant         VARCHAR(50) DEFAULT 'primary',  -- primary, secondary, outline, ghost
    size            VARCHAR(20) DEFAULT 'md',        -- sm, md, lg
    icon            VARCHAR(100),
    icon_position   VARCHAR(10) DEFAULT 'left',      -- left, right
    is_active       BOOLEAN NOT NULL DEFAULT true,
    sort_order      INTEGER NOT NULL DEFAULT 0,
    metadata        JSONB DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_cta_buttons_site_id ON cta_buttons(site_id);
CREATE INDEX idx_cta_buttons_section_id ON cta_buttons(section_id);

-- Social links table
CREATE TABLE IF NOT EXISTS social_links (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    site_id     UUID NOT NULL REFERENCES sites(id) ON DELETE CASCADE,
    platform    VARCHAR(100) NOT NULL,  -- facebook, twitter, instagram, linkedin, etc.
    url         TEXT NOT NULL,
    icon        VARCHAR(100),
    label       VARCHAR(100),
    is_active   BOOLEAN NOT NULL DEFAULT true,
    sort_order  INTEGER NOT NULL DEFAULT 0,
    metadata    JSONB DEFAULT '{}',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_social_links_site_platform ON social_links(site_id, platform);
CREATE INDEX idx_social_links_site_id ON social_links(site_id);

-- Team members table
CREATE TABLE IF NOT EXISTS team_members (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    site_id         UUID NOT NULL REFERENCES sites(id) ON DELETE CASCADE,
    section_id      UUID REFERENCES page_sections(id) ON DELETE SET NULL,
    name            VARCHAR(255) NOT NULL,
    title           VARCHAR(255),
    bio             TEXT,
    avatar_url      TEXT,
    email           VARCHAR(255),
    linkedin_url    TEXT,
    twitter_url     TEXT,
    github_url      TEXT,
    is_active       BOOLEAN NOT NULL DEFAULT true,
    sort_order      INTEGER NOT NULL DEFAULT 0,
    metadata        JSONB DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_team_members_site_id ON team_members(site_id);
CREATE INDEX idx_team_members_section_id ON team_members(section_id);
CREATE INDEX idx_team_members_deleted_at ON team_members(deleted_at) WHERE deleted_at IS NULL;

-- Logo/Partner logos table
CREATE TABLE IF NOT EXISTS logos (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    site_id     UUID NOT NULL REFERENCES sites(id) ON DELETE CASCADE,
    section_id  UUID REFERENCES page_sections(id) ON DELETE SET NULL,
    name        VARCHAR(255) NOT NULL,
    image_url   TEXT NOT NULL,
    alt_text    VARCHAR(255),
    link_url    TEXT,
    is_active   BOOLEAN NOT NULL DEFAULT true,
    sort_order  INTEGER NOT NULL DEFAULT 0,
    metadata    JSONB DEFAULT '{}',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_logos_site_id ON logos(site_id);
CREATE INDEX idx_logos_section_id ON logos(section_id);

-- Apply triggers
CREATE TRIGGER update_features_updated_at BEFORE UPDATE ON features FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_testimonials_updated_at BEFORE UPDATE ON testimonials FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_pricing_plans_updated_at BEFORE UPDATE ON pricing_plans FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_faqs_updated_at BEFORE UPDATE ON faqs FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_cta_buttons_updated_at BEFORE UPDATE ON cta_buttons FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_social_links_updated_at BEFORE UPDATE ON social_links FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_team_members_updated_at BEFORE UPDATE ON team_members FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_logos_updated_at BEFORE UPDATE ON logos FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Insert default data
DO $$
DECLARE
    v_site_id UUID;
    v_features_section_id UUID;
    v_testimonials_section_id UUID;
    v_pricing_section_id UUID;
    v_faq_section_id UUID;
BEGIN
    SELECT id INTO v_site_id FROM sites WHERE slug = 'default' LIMIT 1;
    
    SELECT ps.id INTO v_features_section_id 
    FROM page_sections ps JOIN pages p ON p.id = ps.page_id 
    WHERE p.slug = 'home' AND ps.identifier = 'features' LIMIT 1;
    
    SELECT ps.id INTO v_testimonials_section_id 
    FROM page_sections ps JOIN pages p ON p.id = ps.page_id 
    WHERE p.slug = 'home' AND ps.identifier = 'testimonials' LIMIT 1;
    
    SELECT ps.id INTO v_pricing_section_id 
    FROM page_sections ps JOIN pages p ON p.id = ps.page_id 
    WHERE p.slug = 'home' AND ps.identifier = 'pricing' LIMIT 1;
    
    SELECT ps.id INTO v_faq_section_id 
    FROM page_sections ps JOIN pages p ON p.id = ps.page_id 
    WHERE p.slug = 'home' AND ps.identifier = 'faq' LIMIT 1;

    IF v_site_id IS NOT NULL THEN
        -- Default features
        INSERT INTO features (site_id, section_id, title, description, icon, sort_order) VALUES
        (v_site_id, v_features_section_id, 'Lightning Fast', 'Built for performance. Load times under 100ms with our optimized infrastructure.', 'zap', 1),
        (v_site_id, v_features_section_id, 'Secure by Default', 'Enterprise-grade security with end-to-end encryption and SOC 2 compliance.', 'shield', 2),
        (v_site_id, v_features_section_id, 'Easy Integration', 'Connect with 100+ tools you already use. REST API and webhooks included.', 'plug', 3),
        (v_site_id, v_features_section_id, 'Real-time Analytics', 'Track every interaction with detailed analytics and custom dashboards.', 'bar-chart', 4),
        (v_site_id, v_features_section_id, 'Team Collaboration', 'Work together seamlessly with roles, permissions, and real-time updates.', 'users', 5),
        (v_site_id, v_features_section_id, '24/7 Support', 'Our expert team is always available to help you succeed.', 'headphones', 6);

        -- Default testimonials
        INSERT INTO testimonials (site_id, section_id, author_name, author_title, author_company, content, rating, is_featured, sort_order) VALUES
        (v_site_id, v_testimonials_section_id, 'Sarah Johnson', 'CEO', 'TechCorp', 'This platform transformed how we manage our digital presence. The CMS is incredibly intuitive and the performance is outstanding.', 5, true, 1),
        (v_site_id, v_testimonials_section_id, 'Michael Chen', 'CTO', 'StartupXYZ', 'We migrated from a legacy system and the difference is night and day. Setup took less than an hour and the team was up and running immediately.', 5, true, 2),
        (v_site_id, v_testimonials_section_id, 'Emily Rodriguez', 'Marketing Director', 'GrowthCo', 'The ability to update landing page content without touching code has been a game-changer for our marketing team.', 5, false, 3);

        -- Default pricing plans
        INSERT INTO pricing_plans (site_id, section_id, name, description, price_monthly, price_yearly, is_popular, cta_text, cta_link, features, sort_order) VALUES
        (v_site_id, v_pricing_section_id, 'Starter', 'Perfect for individuals and small projects', 9.00, 90.00, false, 'Get Started', '/signup?plan=starter', '["1 website", "5 pages", "10GB storage", "Basic analytics", "Email support"]', 1),
        (v_site_id, v_pricing_section_id, 'Pro', 'For growing businesses and teams', 29.00, 290.00, true, 'Start Free Trial', '/signup?plan=pro', '["5 websites", "Unlimited pages", "50GB storage", "Advanced analytics", "Priority support", "Custom domain", "API access"]', 2),
        (v_site_id, v_pricing_section_id, 'Enterprise', 'For large organizations with custom needs', NULL, NULL, false, 'Contact Sales', '/contact', '["Unlimited websites", "Unlimited pages", "500GB storage", "Custom analytics", "24/7 dedicated support", "SLA guarantee", "Custom integrations", "SSO/SAML"]', 3);

        -- Default FAQs
        INSERT INTO faqs (site_id, section_id, question, answer, sort_order) VALUES
        (v_site_id, v_faq_section_id, 'How do I get started?', 'Simply sign up for a free account, choose your plan, and you can have your landing page live in minutes. No technical knowledge required.', 1),
        (v_site_id, v_faq_section_id, 'Can I change my plan later?', 'Yes, you can upgrade or downgrade your plan at any time. Changes take effect immediately and billing is prorated.', 2),
        (v_site_id, v_faq_section_id, 'Is there a free trial?', 'Yes! All plans come with a 14-day free trial. No credit card required to start.', 3),
        (v_site_id, v_faq_section_id, 'What payment methods do you accept?', 'We accept all major credit cards (Visa, MasterCard, American Express), PayPal, and bank transfers for annual plans.', 4),
        (v_site_id, v_faq_section_id, 'Can I use my own domain?', 'Yes, you can connect your custom domain on Pro and Enterprise plans. We provide free SSL certificates for all custom domains.', 5),
        (v_site_id, v_faq_section_id, 'How secure is my data?', 'We take security seriously. All data is encrypted at rest and in transit. We are SOC 2 Type II certified and GDPR compliant.', 6);

        -- Default social links
        INSERT INTO social_links (site_id, platform, url, icon, label, sort_order) VALUES
        (v_site_id, 'twitter', 'https://twitter.com/myproduct', 'twitter', 'Twitter', 1),
        (v_site_id, 'linkedin', 'https://linkedin.com/company/myproduct', 'linkedin', 'LinkedIn', 2),
        (v_site_id, 'github', 'https://github.com/myproduct', 'github', 'GitHub', 3);
    END IF;
END $$;

-- ============================================================
-- ROLLBACK SCRIPT
-- ============================================================
-- DROP TABLE IF EXISTS logos CASCADE;
-- DROP TABLE IF EXISTS team_members CASCADE;
-- DROP TABLE IF EXISTS social_links CASCADE;
-- DROP TABLE IF EXISTS cta_buttons CASCADE;
-- DROP TABLE IF EXISTS faqs CASCADE;
-- DROP TABLE IF EXISTS pricing_plans CASCADE;
-- DROP TABLE IF EXISTS testimonials CASCADE;
-- DROP TABLE IF EXISTS features CASCADE;
