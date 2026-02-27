-- ============================================================
-- Seeder: seed_landing_page.sql
-- Description: Complete seeder for landing page content
-- Includes: Site, Pages, Sections, Content, Navigation,
--           Features, Testimonials, Pricing Plans, FAQs
-- ============================================================

-- This seeder creates a complete, production-ready landing page
-- with all sections populated with realistic content.

DO $$
DECLARE
    v_site_id UUID;
    v_home_page_id UUID;
    v_about_page_id UUID;
    v_pricing_page_id UUID;
    v_contact_page_id UUID;
    v_privacy_page_id UUID;
    v_terms_page_id UUID;
    
    -- Section IDs for home page
    v_hero_section_id UUID;
    v_features_section_id UUID;
    v_stats_section_id UUID;
    v_testimonials_section_id UUID;
    v_pricing_section_id UUID;
    v_faq_section_id UUID;
    v_cta_section_id UUID;
    
    -- Navigation menu IDs
    v_header_menu_id UUID;
    v_footer_menu_id UUID;
    
BEGIN
    -- ─── Get or create default site ───────────────────────────────────────────
    SELECT id INTO v_site_id FROM sites WHERE slug = 'default' LIMIT 1;
    
    IF v_site_id IS NULL THEN
        INSERT INTO sites (name, slug, domain, description, is_active)
        VALUES ('My Awesome Product', 'default', 'localhost', 'The best solution for your needs', true)
        RETURNING id INTO v_site_id;
    END IF;
    
    RAISE NOTICE 'Using site ID: %', v_site_id;
    
    -- ─── Update site settings ─────────────────────────────────────────────────
    INSERT INTO site_settings (site_id, key, value, type, group_name, label, description, is_public, sort_order)
    VALUES
        -- General
        (v_site_id, 'site_title', 'GrowthOS', 'string', 'general', 'Site Title', 'The main title of your site', true, 1),
        (v_site_id, 'site_tagline', 'Scale your business with confidence', 'string', 'general', 'Site Tagline', 'A short tagline for your site', true, 2),
        (v_site_id, 'site_description', 'GrowthOS is the all-in-one platform that helps startups and enterprises scale faster with powerful tools, analytics, and automation.', 'string', 'general', 'Site Description', 'Full description of your site', true, 3),
        (v_site_id, 'contact_email', 'hello@growthOS.io', 'string', 'general', 'Contact Email', 'Primary contact email', true, 4),
        (v_site_id, 'contact_phone', '+1 (555) 123-4567', 'string', 'general', 'Contact Phone', 'Primary contact phone', true, 5),
        (v_site_id, 'contact_address', '100 Innovation Drive, San Francisco, CA 94105', 'string', 'general', 'Contact Address', 'Physical address', true, 6),
        -- SEO
        (v_site_id, 'seo_title', 'GrowthOS - Scale Your Business with Confidence', 'string', 'seo', 'SEO Title', 'Title for search engines (max 60 chars)', true, 10),
        (v_site_id, 'seo_description', 'GrowthOS helps startups and enterprises scale faster with powerful analytics, automation, and collaboration tools. Start free today.', 'string', 'seo', 'SEO Description', 'Description for search engines (max 160 chars)', true, 11),
        (v_site_id, 'seo_keywords', 'growth platform, business analytics, startup tools, enterprise software, automation', 'string', 'seo', 'SEO Keywords', 'Comma-separated keywords', true, 12),
        (v_site_id, 'og_title', 'GrowthOS - Scale Your Business', 'string', 'seo', 'OG Title', 'Open Graph title for social sharing', true, 13),
        (v_site_id, 'og_description', 'The all-in-one platform for scaling your business. Join 10,000+ companies already growing with GrowthOS.', 'string', 'seo', 'OG Description', 'Open Graph description', true, 14),
        (v_site_id, 'og_image', 'https://images.unsplash.com/photo-1551434678-e076c223a692?w=1200&h=630&fit=crop', 'string', 'seo', 'OG Image URL', 'Open Graph image (1200x630px)', true, 15),
        (v_site_id, 'twitter_card', 'summary_large_image', 'string', 'seo', 'Twitter Card Type', 'Twitter card type', true, 16),
        (v_site_id, 'twitter_site', '@growthOS', 'string', 'seo', 'Twitter Site Handle', 'Twitter @handle', true, 17),
        -- Social
        (v_site_id, 'social_twitter', 'https://twitter.com/growthOS', 'string', 'social', 'Twitter URL', 'Twitter profile URL', true, 20),
        (v_site_id, 'social_linkedin', 'https://linkedin.com/company/growthOS', 'string', 'social', 'LinkedIn URL', 'LinkedIn company page URL', true, 21),
        (v_site_id, 'social_github', 'https://github.com/growthOS', 'string', 'social', 'GitHub URL', 'GitHub organization URL', true, 22),
        (v_site_id, 'social_facebook', '', 'string', 'social', 'Facebook URL', 'Facebook page URL', true, 23),
        (v_site_id, 'social_instagram', '', 'string', 'social', 'Instagram URL', 'Instagram profile URL', true, 24),
        -- Analytics
        (v_site_id, 'google_analytics_id', '', 'string', 'analytics', 'Google Analytics ID', 'GA4 Measurement ID (G-XXXXXXXXXX)', false, 30),
        (v_site_id, 'google_tag_manager_id', '', 'string', 'analytics', 'Google Tag Manager ID', 'GTM Container ID (GTM-XXXXXXX)', false, 31),
        -- Appearance
        (v_site_id, 'primary_color', '#6366f1', 'color', 'appearance', 'Primary Color', 'Main brand color', true, 40),
        (v_site_id, 'secondary_color', '#8b5cf6', 'color', 'appearance', 'Secondary Color', 'Secondary brand color', true, 41),
        (v_site_id, 'font_family', 'Inter', 'string', 'appearance', 'Font Family', 'Primary font family', true, 42),
        (v_site_id, 'footer_text', '© 2024 GrowthOS. All rights reserved.', 'string', 'appearance', 'Footer Text', 'Copyright text in footer', true, 43)
    ON CONFLICT (site_id, key) DO UPDATE SET
        value = EXCLUDED.value,
        label = EXCLUDED.label,
        description = EXCLUDED.description,
        is_public = EXCLUDED.is_public;
    
    -- ─── Create Pages ─────────────────────────────────────────────────────────
    
    -- Home Page
    INSERT INTO pages (
        site_id, title, slug, description, status, is_homepage,
        seo_title, seo_description, seo_keywords,
        og_title, og_description,
        robots_meta, sort_order
    ) VALUES (
        v_site_id,
        'Home',
        'home',
        'Main landing page for GrowthOS',
        'published',
        true,
        'GrowthOS - Scale Your Business with Confidence',
        'GrowthOS helps startups and enterprises scale faster with powerful analytics, automation, and collaboration tools.',
        'growth platform, business analytics, startup tools',
        'GrowthOS - Scale Your Business',
        'The all-in-one platform for scaling your business.',
        'index, follow',
        1
    )
    ON CONFLICT DO NOTHING
    RETURNING id INTO v_home_page_id;
    
    IF v_home_page_id IS NULL THEN
        SELECT id INTO v_home_page_id FROM pages WHERE slug = 'home' AND site_id = v_site_id LIMIT 1;
    END IF;
    
    -- About Page
    INSERT INTO pages (
        site_id, title, slug, description, status, is_homepage,
        seo_title, seo_description, robots_meta, sort_order
    ) VALUES (
        v_site_id,
        'About Us',
        'about',
        'Learn about GrowthOS and our mission',
        'published',
        false,
        'About GrowthOS - Our Story and Mission',
        'Learn about GrowthOS, our team, and our mission to help businesses scale with confidence.',
        'index, follow',
        2
    )
    ON CONFLICT DO NOTHING
    RETURNING id INTO v_about_page_id;
    
    IF v_about_page_id IS NULL THEN
        SELECT id INTO v_about_page_id FROM pages WHERE slug = 'about' AND site_id = v_site_id LIMIT 1;
    END IF;
    
    -- Pricing Page
    INSERT INTO pages (
        site_id, title, slug, description, status, is_homepage,
        seo_title, seo_description, robots_meta, sort_order
    ) VALUES (
        v_site_id,
        'Pricing',
        'pricing',
        'Simple, transparent pricing for every team',
        'published',
        false,
        'GrowthOS Pricing - Simple, Transparent Plans',
        'Choose the perfect plan for your team. Start free, upgrade when you need more.',
        'index, follow',
        3
    )
    ON CONFLICT DO NOTHING
    RETURNING id INTO v_pricing_page_id;
    
    IF v_pricing_page_id IS NULL THEN
        SELECT id INTO v_pricing_page_id FROM pages WHERE slug = 'pricing' AND site_id = v_site_id LIMIT 1;
    END IF;
    
    -- Contact Page
    INSERT INTO pages (
        site_id, title, slug, description, status, is_homepage,
        seo_title, seo_description, robots_meta, sort_order
    ) VALUES (
        v_site_id,
        'Contact',
        'contact',
        'Get in touch with the GrowthOS team',
        'published',
        false,
        'Contact GrowthOS - We''re Here to Help',
        'Have questions? Our team is ready to help. Contact us via email, phone, or live chat.',
        'index, follow',
        4
    )
    ON CONFLICT DO NOTHING
    RETURNING id INTO v_contact_page_id;
    
    IF v_contact_page_id IS NULL THEN
        SELECT id INTO v_contact_page_id FROM pages WHERE slug = 'contact' AND site_id = v_site_id LIMIT 1;
    END IF;
    
    -- Privacy Policy Page
    INSERT INTO pages (
        site_id, title, slug, description, status, is_homepage,
        seo_title, robots_meta, sort_order
    ) VALUES (
        v_site_id,
        'Privacy Policy',
        'privacy',
        'GrowthOS Privacy Policy',
        'published',
        false,
        'Privacy Policy - GrowthOS',
        'noindex, follow',
        10
    )
    ON CONFLICT DO NOTHING
    RETURNING id INTO v_privacy_page_id;
    
    IF v_privacy_page_id IS NULL THEN
        SELECT id INTO v_privacy_page_id FROM pages WHERE slug = 'privacy' AND site_id = v_site_id LIMIT 1;
    END IF;
    
    -- Terms of Service Page
    INSERT INTO pages (
        site_id, title, slug, description, status, is_homepage,
        seo_title, robots_meta, sort_order
    ) VALUES (
        v_site_id,
        'Terms of Service',
        'terms',
        'GrowthOS Terms of Service',
        'published',
        false,
        'Terms of Service - GrowthOS',
        'noindex, follow',
        11
    )
    ON CONFLICT DO NOTHING
    RETURNING id INTO v_terms_page_id;
    
    IF v_terms_page_id IS NULL THEN
        SELECT id INTO v_terms_page_id FROM pages WHERE slug = 'terms' AND site_id = v_site_id LIMIT 1;
    END IF;
    
    RAISE NOTICE 'Pages created. Home page ID: %', v_home_page_id;
    
    -- ─── Create Home Page Sections ────────────────────────────────────────────
    
    IF v_home_page_id IS NOT NULL THEN
        -- Delete existing sections for clean seeding
        DELETE FROM page_sections WHERE page_id = v_home_page_id;
        
        -- Hero Section
        INSERT INTO page_sections (page_id, name, type, identifier, is_visible, sort_order, bg_color)
        VALUES (v_home_page_id, 'Hero Section', 'hero', 'hero', true, 1, '#ffffff')
        RETURNING id INTO v_hero_section_id;
        
        -- Features Section
        INSERT INTO page_sections (page_id, name, type, identifier, is_visible, sort_order, bg_color)
        VALUES (v_home_page_id, 'Features Section', 'features', 'features', true, 2, '#f9fafb')
        RETURNING id INTO v_features_section_id;
        
        -- Stats Section
        INSERT INTO page_sections (page_id, name, type, identifier, is_visible, sort_order, bg_color)
        VALUES (v_home_page_id, 'Stats Section', 'stats', 'stats', true, 3, '#6366f1')
        RETURNING id INTO v_stats_section_id;
        
        -- Testimonials Section
        INSERT INTO page_sections (page_id, name, type, identifier, is_visible, sort_order, bg_color)
        VALUES (v_home_page_id, 'Testimonials Section', 'testimonials', 'testimonials', true, 4, '#ffffff')
        RETURNING id INTO v_testimonials_section_id;
        
        -- Pricing Section
        INSERT INTO page_sections (page_id, name, type, identifier, is_visible, sort_order, bg_color)
        VALUES (v_home_page_id, 'Pricing Section', 'pricing', 'pricing', true, 5, '#f9fafb')
        RETURNING id INTO v_pricing_section_id;
        
        -- FAQ Section
        INSERT INTO page_sections (page_id, name, type, identifier, is_visible, sort_order, bg_color)
        VALUES (v_home_page_id, 'FAQ Section', 'faq', 'faq', true, 6, '#ffffff')
        RETURNING id INTO v_faq_section_id;
        
        -- CTA Section
        INSERT INTO page_sections (page_id, name, type, identifier, is_visible, sort_order, bg_color)
        VALUES (v_home_page_id, 'CTA Section', 'cta', 'cta', true, 7, '#6366f1')
        RETURNING id INTO v_cta_section_id;
        
        RAISE NOTICE 'Sections created for home page';
        
        -- ─── Hero Section Content ──────────────────────────────────────────────
        INSERT INTO section_contents (section_id, key, value, type, label, sort_order) VALUES
        (v_hero_section_id, 'badge_text', '🚀 Now with AI-powered insights', 'text', 'Badge Text', 1),
        (v_hero_section_id, 'title', 'Scale Your Business with Confidence', 'text', 'Main Title', 2),
        (v_hero_section_id, 'title_highlight', 'Confidence', 'text', 'Highlighted Word', 3),
        (v_hero_section_id, 'subtitle', 'GrowthOS is the all-in-one platform that helps startups and enterprises grow faster with powerful analytics, automation, and team collaboration tools.', 'text', 'Subtitle', 4),
        (v_hero_section_id, 'cta_primary_text', 'Start Free Trial', 'text', 'Primary CTA Text', 5),
        (v_hero_section_id, 'cta_primary_link', '#pricing', 'link', 'Primary CTA Link', 6),
        (v_hero_section_id, 'cta_secondary_text', 'Watch Demo', 'text', 'Secondary CTA Text', 7),
        (v_hero_section_id, 'cta_secondary_link', '#demo', 'link', 'Secondary CTA Link', 8),
        (v_hero_section_id, 'hero_image', 'https://images.unsplash.com/photo-1551434678-e076c223a692?w=1200&h=675&fit=crop', 'image', 'Hero Image URL', 9),
        (v_hero_section_id, 'hero_image_alt', 'GrowthOS Dashboard Screenshot', 'text', 'Hero Image Alt Text', 10),
        (v_hero_section_id, 'social_proof_text', 'Trusted by 10,000+ teams at companies like Stripe, Notion, and Linear', 'text', 'Social Proof Text', 11)
        ON CONFLICT DO NOTHING;
        
        -- ─── Features Section Content ──────────────────────────────────────────
        INSERT INTO section_contents (section_id, key, value, type, label, sort_order) VALUES
        (v_features_section_id, 'badge_text', 'Features', 'text', 'Badge Text', 1),
        (v_features_section_id, 'title', 'Everything you need to grow', 'text', 'Section Title', 2),
        (v_features_section_id, 'subtitle', 'Powerful features designed for modern teams. From analytics to automation, we have everything you need to scale your business.', 'text', 'Section Subtitle', 3)
        ON CONFLICT DO NOTHING;
        
        -- ─── Stats Section Content ─────────────────────────────────────────────
        INSERT INTO section_contents (section_id, key, value, type, label, sort_order) VALUES
        (v_stats_section_id, 'stat_1_value', '10,000+', 'text', 'Stat 1 Value', 1),
        (v_stats_section_id, 'stat_1_label', 'Active Teams', 'text', 'Stat 1 Label', 2),
        (v_stats_section_id, 'stat_2_value', '99.9%', 'text', 'Stat 2 Value', 3),
        (v_stats_section_id, 'stat_2_label', 'Uptime SLA', 'text', 'Stat 2 Label', 4),
        (v_stats_section_id, 'stat_3_value', '50+', 'text', 'Stat 3 Value', 5),
        (v_stats_section_id, 'stat_3_label', 'Countries', 'text', 'Stat 3 Label', 6),
        (v_stats_section_id, 'stat_4_value', '2.5B+', 'text', 'Stat 4 Value', 7),
        (v_stats_section_id, 'stat_4_label', 'Events Tracked', 'text', 'Stat 4 Label', 8)
        ON CONFLICT DO NOTHING;
        
        -- ─── Testimonials Section Content ──────────────────────────────────────
        INSERT INTO section_contents (section_id, key, value, type, label, sort_order) VALUES
        (v_testimonials_section_id, 'badge_text', 'Testimonials', 'text', 'Badge Text', 1),
        (v_testimonials_section_id, 'title', 'Loved by thousands of teams', 'text', 'Section Title', 2),
        (v_testimonials_section_id, 'subtitle', 'Don''t just take our word for it. Here''s what our customers have to say about GrowthOS.', 'text', 'Section Subtitle', 3)
        ON CONFLICT DO NOTHING;
        
        -- ─── Pricing Section Content ───────────────────────────────────────────
        INSERT INTO section_contents (section_id, key, value, type, label, sort_order) VALUES
        (v_pricing_section_id, 'badge_text', 'Pricing', 'text', 'Badge Text', 1),
        (v_pricing_section_id, 'title', 'Simple, transparent pricing', 'text', 'Section Title', 2),
        (v_pricing_section_id, 'subtitle', 'No hidden fees. No surprises. Start free and upgrade when you need more.', 'text', 'Section Subtitle', 3)
        ON CONFLICT DO NOTHING;
        
        -- ─── FAQ Section Content ───────────────────────────────────────────────
        INSERT INTO section_contents (section_id, key, value, type, label, sort_order) VALUES
        (v_faq_section_id, 'badge_text', 'FAQ', 'text', 'Badge Text', 1),
        (v_faq_section_id, 'title', 'Frequently asked questions', 'text', 'Section Title', 2),
        (v_faq_section_id, 'subtitle', 'Everything you need to know about GrowthOS. Can''t find the answer you''re looking for? Contact our support team.', 'text', 'Section Subtitle', 3)
        ON CONFLICT DO NOTHING;
        
        -- ─── CTA Section Content ───────────────────────────────────────────────
        INSERT INTO section_contents (section_id, key, value, type, label, sort_order) VALUES
        (v_cta_section_id, 'title', 'Ready to scale your business?', 'text', 'CTA Title', 1),
        (v_cta_section_id, 'subtitle', 'Join 10,000+ teams already growing with GrowthOS. Start your free trial today — no credit card required.', 'text', 'CTA Subtitle', 2),
        (v_cta_section_id, 'cta_primary_text', 'Start Free Trial', 'text', 'Primary CTA Text', 3),
        (v_cta_section_id, 'cta_primary_link', '#pricing', 'link', 'Primary CTA Link', 4),
        (v_cta_section_id, 'cta_secondary_text', 'Talk to Sales', 'text', 'Secondary CTA Text', 5),
        (v_cta_section_id, 'cta_secondary_link', '/contact', 'link', 'Secondary CTA Link', 6)
        ON CONFLICT DO NOTHING;
        
    END IF;
    
    -- ─── Features ─────────────────────────────────────────────────────────────
    
    -- Delete existing features for clean seeding
    DELETE FROM features WHERE site_id = v_site_id;
    
    INSERT INTO features (site_id, section_id, title, description, icon, icon_color, is_active, sort_order) VALUES
    (v_site_id, v_features_section_id, 'Real-time Analytics', 'Get instant insights into your business performance with live dashboards and custom reports.', 'bar-chart', '#6366f1', true, 1),
    (v_site_id, v_features_section_id, 'Smart Automation', 'Automate repetitive tasks and workflows to save time and reduce human error.', 'zap', '#8b5cf6', true, 2),
    (v_site_id, v_features_section_id, 'Team Collaboration', 'Work together seamlessly with real-time collaboration tools and shared workspaces.', 'users', '#06b6d4', true, 3),
    (v_site_id, v_features_section_id, 'Enterprise Security', 'Bank-grade security with SOC 2 compliance, SSO, and advanced access controls.', 'shield', '#10b981', true, 4),
    (v_site_id, v_features_section_id, 'API & Integrations', 'Connect with 200+ tools including Slack, Salesforce, HubSpot, and more.', 'plug', '#f59e0b', true, 5),
    (v_site_id, v_features_section_id, '24/7 Support', 'Get help when you need it with our dedicated support team and extensive documentation.', 'headphones', '#ef4444', true, 6),
    (v_site_id, v_features_section_id, 'Global Scale', 'Deploy to 50+ regions worldwide with 99.9% uptime SLA and automatic failover.', 'globe', '#3b82f6', true, 7),
    (v_site_id, v_features_section_id, 'AI-Powered Insights', 'Leverage machine learning to predict trends, identify opportunities, and optimize performance.', 'star', '#f97316', true, 8);
    
    RAISE NOTICE 'Features created';
    
    -- ─── Testimonials ─────────────────────────────────────────────────────────
    
    -- Delete existing testimonials for clean seeding
    DELETE FROM testimonials WHERE site_id = v_site_id;
    
    INSERT INTO testimonials (site_id, section_id, author_name, author_title, author_company, content, rating, is_featured, is_active, sort_order) VALUES
    (v_site_id, v_testimonials_section_id, 'Sarah Chen', 'Head of Growth', 'Stripe', 'GrowthOS transformed how we track and optimize our growth metrics. The real-time analytics are incredibly powerful, and the automation features saved our team 20+ hours per week.', 5, true, true, 1),
    (v_site_id, v_testimonials_section_id, 'Marcus Johnson', 'CTO', 'Notion', 'We evaluated 10+ platforms before choosing GrowthOS. The API integrations are seamless, the security is enterprise-grade, and the support team is exceptional. Highly recommended.', 5, true, true, 2),
    (v_site_id, v_testimonials_section_id, 'Emily Rodriguez', 'VP of Marketing', 'Linear', 'The AI-powered insights have been a game-changer for our marketing campaigns. We''ve seen a 40% improvement in conversion rates since switching to GrowthOS.', 5, false, true, 3),
    (v_site_id, v_testimonials_section_id, 'David Kim', 'Founder & CEO', 'Vercel', 'GrowthOS scales with us perfectly. From 10 users to 10,000, the platform has never let us down. The team collaboration features are best-in-class.', 5, true, true, 4),
    (v_site_id, v_testimonials_section_id, 'Priya Patel', 'Director of Operations', 'Figma', 'Implementation was surprisingly smooth. Our team was up and running in less than a day. The onboarding experience and documentation are excellent.', 4, false, true, 5),
    (v_site_id, v_testimonials_section_id, 'Alex Thompson', 'Engineering Manager', 'Shopify', 'The automation workflows have eliminated so many manual processes. Our team can now focus on strategic work instead of repetitive tasks. ROI was immediate.', 5, false, true, 6);
    
    RAISE NOTICE 'Testimonials created';
    
    -- ─── Pricing Plans ────────────────────────────────────────────────────────
    
    -- Delete existing pricing plans for clean seeding
    DELETE FROM pricing_plans WHERE site_id = v_site_id;
    
    INSERT INTO pricing_plans (
        site_id, section_id, name, description,
        price_monthly, price_yearly, currency,
        cta_text, cta_link,
        is_popular, is_custom, is_active, sort_order,
        features, features_excluded
    ) VALUES
    (
        v_site_id, v_pricing_section_id,
        'Starter', 'Perfect for small teams and startups',
        29, 290, 'USD',
        'Start Free Trial', '/signup?plan=starter',
        false, false, true, 1,
        '["Up to 5 team members", "10,000 events/month", "Basic analytics dashboard", "Email support", "API access", "2 integrations", "7-day data retention"]'::jsonb,
        '["Advanced analytics", "Custom reports", "Priority support", "SSO", "Unlimited integrations"]'::jsonb
    ),
    (
        v_site_id, v_pricing_section_id,
        'Pro', 'For growing teams that need more power',
        79, 790, 'USD',
        'Start Free Trial', '/signup?plan=pro',
        true, false, true, 2,
        '["Up to 25 team members", "500,000 events/month", "Advanced analytics & reports", "Priority email & chat support", "Full API access", "Unlimited integrations", "30-day data retention", "Custom dashboards", "Automation workflows", "A/B testing"]'::jsonb,
        '["SSO", "Dedicated account manager", "Custom SLA"]'::jsonb
    ),
    (
        v_site_id, v_pricing_section_id,
        'Enterprise', 'For large organizations with custom needs',
        NULL, NULL, 'USD',
        'Contact Sales', '/contact?plan=enterprise',
        false, true, true, 3,
        '["Unlimited team members", "Unlimited events", "Custom analytics & reporting", "24/7 dedicated support", "Full API access", "Unlimited integrations", "Unlimited data retention", "Custom dashboards", "Advanced automation", "A/B testing", "SSO & SAML", "Dedicated account manager", "Custom SLA", "On-premise deployment option", "Security audit & compliance"]'::jsonb,
        '[]'::jsonb
    );
    
    RAISE NOTICE 'Pricing plans created';
    
    -- ─── FAQs ─────────────────────────────────────────────────────────────────
    
    -- Delete existing FAQs for clean seeding
    DELETE FROM faqs WHERE site_id = v_site_id;
    
    INSERT INTO faqs (site_id, section_id, question, answer, category, is_active, sort_order) VALUES
    -- General
    (v_site_id, v_faq_section_id, 'What is GrowthOS?', 'GrowthOS is an all-in-one business growth platform that combines analytics, automation, and team collaboration tools. It helps startups and enterprises track performance, automate workflows, and make data-driven decisions to scale faster.', 'General', true, 1),
    (v_site_id, v_faq_section_id, 'How does the free trial work?', 'You can start a 14-day free trial with full access to all Pro features — no credit card required. After the trial, you can choose to upgrade to a paid plan or continue with our free Starter plan.', 'General', true, 2),
    (v_site_id, v_faq_section_id, 'Can I change my plan at any time?', 'Yes! You can upgrade or downgrade your plan at any time. When you upgrade, you''ll be charged the prorated difference. When you downgrade, the change takes effect at the end of your billing cycle.', 'General', true, 3),
    -- Billing
    (v_site_id, v_faq_section_id, 'What payment methods do you accept?', 'We accept all major credit cards (Visa, Mastercard, American Express), PayPal, and bank transfers for annual plans. Enterprise customers can also pay by invoice.', 'Billing', true, 4),
    (v_site_id, v_faq_section_id, 'Is there a discount for annual billing?', 'Yes! When you choose annual billing, you save approximately 17% compared to monthly billing. For example, the Pro plan is $79/month billed monthly, or $790/year (equivalent to $65.83/month) billed annually.', 'Billing', true, 5),
    (v_site_id, v_faq_section_id, 'Do you offer refunds?', 'We offer a 30-day money-back guarantee for all paid plans. If you''re not satisfied for any reason, contact our support team within 30 days of your purchase for a full refund.', 'Billing', true, 6),
    -- Technical
    (v_site_id, v_faq_section_id, 'How do I integrate GrowthOS with my existing tools?', 'GrowthOS integrates with 200+ popular tools including Slack, Salesforce, HubSpot, Zapier, and more. You can set up integrations in minutes from the Integrations page in your dashboard. We also offer a full REST API for custom integrations.', 'Technical', true, 7),
    (v_site_id, v_faq_section_id, 'Is my data secure?', 'Absolutely. GrowthOS is SOC 2 Type II certified and GDPR compliant. We use AES-256 encryption for data at rest and TLS 1.3 for data in transit. We also offer SSO, 2FA, and advanced access controls for Enterprise customers.', 'Technical', true, 8),
    (v_site_id, v_faq_section_id, 'What is your uptime SLA?', 'We guarantee 99.9% uptime for all paid plans. Our infrastructure is distributed across multiple regions with automatic failover. You can check our real-time status at status.growthOS.io.', 'Technical', true, 9),
    -- Support
    (v_site_id, v_faq_section_id, 'What kind of support do you offer?', 'Starter plan customers get email support with a 24-hour response time. Pro plan customers get priority email and live chat support with a 4-hour response time. Enterprise customers get 24/7 dedicated support with a 1-hour response time.', 'Support', true, 10),
    (v_site_id, v_faq_section_id, 'Do you offer onboarding assistance?', 'Yes! All new customers get access to our onboarding resources including video tutorials, documentation, and a getting started guide. Pro and Enterprise customers also get a personalized onboarding call with our customer success team.', 'Support', true, 11);
    
    RAISE NOTICE 'FAQs created';
    
    -- ─── Navigation Menus ─────────────────────────────────────────────────────
    
    -- Delete existing navigation for clean seeding
    DELETE FROM navigation_items WHERE menu_id IN (
        SELECT id FROM navigation_menus WHERE site_id = v_site_id
    );
    DELETE FROM navigation_menus WHERE site_id = v_site_id;
    
    -- Header Navigation
    INSERT INTO navigation_menus (site_id, name, identifier, description, is_active)
    VALUES (v_site_id, 'Header Navigation', 'header', 'Main header navigation menu', true)
    RETURNING id INTO v_header_menu_id;
    
    -- Footer Navigation
    INSERT INTO navigation_menus (site_id, name, identifier, description, is_active)
    VALUES (v_site_id, 'Footer Navigation', 'footer', 'Footer navigation menu', true)
    RETURNING id INTO v_footer_menu_id;
    
    -- Header Navigation Items
    INSERT INTO navigation_items (menu_id, page_id, label, url, target, is_active, sort_order, depth) VALUES
    (v_header_menu_id, v_home_page_id, 'Home', '/', '_self', true, 1, 0),
    (v_header_menu_id, NULL, 'Features', '/#features', '_self', true, 2, 0),
    (v_header_menu_id, v_pricing_page_id, 'Pricing', '/pricing', '_self', true, 3, 0),
    (v_header_menu_id, NULL, 'FAQ', '/#faq', '_self', true, 4, 0),
    (v_header_menu_id, v_about_page_id, 'About', '/about', '_self', true, 5, 0),
    (v_header_menu_id, v_contact_page_id, 'Contact', '/contact', '_self', true, 6, 0);
    
    -- Footer Navigation Items
    INSERT INTO navigation_items (menu_id, page_id, label, url, target, is_active, sort_order, depth) VALUES
    -- Product column
    (v_footer_menu_id, NULL, 'Features', '/#features', '_self', true, 1, 0),
    (v_footer_menu_id, v_pricing_page_id, 'Pricing', '/pricing', '_self', true, 2, 0),
    (v_footer_menu_id, NULL, 'Changelog', '/changelog', '_self', true, 3, 0),
    (v_footer_menu_id, NULL, 'Roadmap', '/roadmap', '_self', true, 4, 0),
    -- Company column
    (v_footer_menu_id, v_about_page_id, 'About Us', '/about', '_self', true, 5, 0),
    (v_footer_menu_id, NULL, 'Blog', '/blog', '_self', true, 6, 0),
    (v_footer_menu_id, NULL, 'Careers', '/careers', '_self', true, 7, 0),
    (v_footer_menu_id, v_contact_page_id, 'Contact', '/contact', '_self', true, 8, 0),
    -- Legal column
    (v_footer_menu_id, v_privacy_page_id, 'Privacy Policy', '/privacy', '_self', true, 9, 0),
    (v_footer_menu_id, v_terms_page_id, 'Terms of Service', '/terms', '_self', true, 10, 0),
    (v_footer_menu_id, NULL, 'Cookie Policy', '/cookies', '_self', true, 11, 0),
    (v_footer_menu_id, NULL, 'Security', '/security', '_self', true, 12, 0);
    
    RAISE NOTICE 'Navigation created';
    
    RAISE NOTICE '✅ Seeding complete! Site ID: %', v_site_id;
    RAISE NOTICE '   Home page ID: %', v_home_page_id;
    RAISE NOTICE '   About page ID: %', v_about_page_id;
    RAISE NOTICE '   Pricing page ID: %', v_pricing_page_id;
    RAISE NOTICE '   Contact page ID: %', v_contact_page_id;
    
END $$;

-- ─── Verify seeding ───────────────────────────────────────────────────────────
SELECT 
    'sites' as table_name, COUNT(*) as count FROM sites WHERE deleted_at IS NULL
UNION ALL
SELECT 'pages', COUNT(*) FROM pages WHERE deleted_at IS NULL
UNION ALL
SELECT 'page_sections', COUNT(*) FROM page_sections
UNION ALL
SELECT 'section_contents', COUNT(*) FROM section_contents
UNION ALL
SELECT 'features', COUNT(*) FROM features WHERE deleted_at IS NULL
UNION ALL
SELECT 'testimonials', COUNT(*) FROM testimonials WHERE deleted_at IS NULL
UNION ALL
SELECT 'pricing_plans', COUNT(*) FROM pricing_plans WHERE deleted_at IS NULL
UNION ALL
SELECT 'faqs', COUNT(*) FROM faqs WHERE deleted_at IS NULL
UNION ALL
SELECT 'navigation_menus', COUNT(*) FROM navigation_menus
UNION ALL
SELECT 'navigation_items', COUNT(*) FROM navigation_items
ORDER BY table_name;
