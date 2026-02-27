-- Migration: 009_seed_landing_page.sql
-- Description: Complete seeder for landing page content
-- This creates a full, production-ready landing page with all sections populated
-- Created: 2024-01-01

-- ============================================================
-- SITE SETTINGS UPDATE
-- ============================================================
DO $$
DECLARE
    v_site_id UUID;
BEGIN
    SELECT id INTO v_site_id FROM sites WHERE slug = 'default' LIMIT 1;
    
    IF v_site_id IS NOT NULL THEN
        -- Update general settings
        INSERT INTO site_settings (site_id, key, value, type, group_name, label, is_public, sort_order) VALUES
        (v_site_id, 'site_title', 'LandingCMS', 'string', 'general', 'Site Title', true, 1),
        (v_site_id, 'site_tagline', 'Build stunning landing pages in minutes', 'string', 'general', 'Site Tagline', true, 2),
        (v_site_id, 'site_description', 'The most powerful CMS for creating dynamic, high-converting landing pages. No code required.', 'string', 'general', 'Site Description', true, 3),
        (v_site_id, 'contact_email', 'hello@landingcms.io', 'string', 'general', 'Contact Email', true, 4),
        (v_site_id, 'contact_phone', '+1 (555) 000-0000', 'string', 'general', 'Contact Phone', true, 5),
        (v_site_id, 'contact_address', 'San Francisco, CA', 'string', 'general', 'Contact Address', true, 6),
        -- SEO
        (v_site_id, 'seo_title', 'LandingCMS - Build Dynamic Landing Pages', 'string', 'seo', 'SEO Title', true, 1),
        (v_site_id, 'seo_description', 'Create, manage, and optimize your landing pages with our powerful CMS. No coding required. Start free today.', 'string', 'seo', 'SEO Description', true, 2),
        (v_site_id, 'seo_keywords', 'landing page, cms, no-code, website builder, marketing', 'string', 'seo', 'SEO Keywords', true, 3),
        (v_site_id, 'og_title', 'LandingCMS - Build Dynamic Landing Pages', 'string', 'seo', 'OG Title', true, 4),
        (v_site_id, 'og_description', 'Create stunning landing pages without code. Manage all content from one powerful dashboard.', 'string', 'seo', 'OG Description', true, 5),
        (v_site_id, 'twitter_card', 'summary_large_image', 'string', 'seo', 'Twitter Card Type', true, 6),
        (v_site_id, 'twitter_site', '@landingcms', 'string', 'seo', 'Twitter Site Handle', true, 7),
        -- Social
        (v_site_id, 'social_twitter', 'https://twitter.com/landingcms', 'string', 'social', 'Twitter URL', true, 1),
        (v_site_id, 'social_linkedin', 'https://linkedin.com/company/landingcms', 'string', 'social', 'LinkedIn URL', true, 2),
        (v_site_id, 'social_github', 'https://github.com/landingcms', 'string', 'social', 'GitHub URL', true, 3),
        -- Appearance
        (v_site_id, 'primary_color', '#6366f1', 'color', 'appearance', 'Primary Color', true, 1),
        (v_site_id, 'secondary_color', '#8b5cf6', 'color', 'appearance', 'Secondary Color', true, 2),
        (v_site_id, 'font_family', 'Inter', 'string', 'appearance', 'Font Family', true, 3)
        ON CONFLICT (site_id, key) DO UPDATE SET value = EXCLUDED.value;
    END IF;
END $$;

-- ============================================================
-- NAVIGATION MENUS
-- ============================================================
DO $$
DECLARE
    v_site_id UUID;
    v_header_menu_id UUID;
    v_footer_menu_id UUID;
    v_footer_product_id UUID;
    v_footer_company_id UUID;
    v_footer_legal_id UUID;
    v_home_page_id UUID;
BEGIN
    SELECT id INTO v_site_id FROM sites WHERE slug = 'default' LIMIT 1;
    SELECT id INTO v_home_page_id FROM pages WHERE slug = 'home' AND site_id = v_site_id LIMIT 1;
    
    IF v_site_id IS NOT NULL THEN
        -- Get or create header menu
        SELECT id INTO v_header_menu_id FROM navigation_menus WHERE site_id = v_site_id AND identifier = 'header' LIMIT 1;
        IF v_header_menu_id IS NULL THEN
            INSERT INTO navigation_menus (site_id, name, identifier, description)
            VALUES (v_site_id, 'Header Navigation', 'header', 'Main header navigation')
            RETURNING id INTO v_header_menu_id;
        END IF;
        
        -- Clear existing header items
        DELETE FROM navigation_items WHERE menu_id = v_header_menu_id;
        
        -- Insert header navigation items
        INSERT INTO navigation_items (menu_id, page_id, label, url, target, sort_order, depth) VALUES
        (v_header_menu_id, v_home_page_id, 'Home', '/', '_self', 1, 0),
        (v_header_menu_id, NULL, 'Features', '/#features', '_self', 2, 0),
        (v_header_menu_id, NULL, 'Pricing', '/#pricing', '_self', 3, 0),
        (v_header_menu_id, NULL, 'Testimonials', '/#testimonials', '_self', 4, 0),
        (v_header_menu_id, NULL, 'FAQ', '/#faq', '_self', 5, 0),
        (v_header_menu_id, NULL, 'Blog', '/blog', '_self', 6, 0);
        
        -- Get or create footer menu
        SELECT id INTO v_footer_menu_id FROM navigation_menus WHERE site_id = v_site_id AND identifier = 'footer' LIMIT 1;
        IF v_footer_menu_id IS NULL THEN
            INSERT INTO navigation_menus (site_id, name, identifier, description)
            VALUES (v_site_id, 'Footer Navigation', 'footer', 'Footer navigation links')
            RETURNING id INTO v_footer_menu_id;
        END IF;
        
        -- Clear existing footer items
        DELETE FROM navigation_items WHERE menu_id = v_footer_menu_id;
        
        -- Insert footer navigation items
        INSERT INTO navigation_items (menu_id, label, url, target, sort_order, depth) VALUES
        (v_footer_menu_id, 'Privacy Policy', '/privacy', '_self', 1, 0),
        (v_footer_menu_id, 'Terms of Service', '/terms', '_self', 2, 0),
        (v_footer_menu_id, 'Cookie Policy', '/cookies', '_self', 3, 0),
        (v_footer_menu_id, 'Sitemap', '/sitemap.xml', '_self', 4, 0);
    END IF;
END $$;

-- ============================================================
-- SECTION CONTENT - HERO
-- ============================================================
DO $$
DECLARE
    v_site_id UUID;
    v_hero_id UUID;
BEGIN
    SELECT id INTO v_site_id FROM sites WHERE slug = 'default' LIMIT 1;
    SELECT ps.id INTO v_hero_id 
    FROM page_sections ps 
    JOIN pages p ON p.id = ps.page_id 
    WHERE p.slug = 'home' AND p.site_id = v_site_id AND ps.identifier = 'hero' LIMIT 1;

    IF v_hero_id IS NOT NULL THEN
        INSERT INTO section_contents (section_id, key, value, type, label, sort_order) VALUES
        (v_hero_id, 'badge_text', '🚀 Now with AI-powered content suggestions', 'text', 'Badge Text', 1),
        (v_hero_id, 'title', 'Build Landing Pages That Convert', 'text', 'Main Title', 2),
        (v_hero_id, 'title_highlight', 'Convert', 'text', 'Highlighted Word', 3),
        (v_hero_id, 'subtitle', 'Create, manage, and optimize your landing pages with our powerful CMS. No coding required. Launch in minutes, not months.', 'text', 'Subtitle', 4),
        (v_hero_id, 'cta_primary_text', 'Start Free Trial', 'text', 'Primary CTA Text', 5),
        (v_hero_id, 'cta_primary_link', '#pricing', 'link', 'Primary CTA Link', 6),
        (v_hero_id, 'cta_secondary_text', 'See Live Demo', 'text', 'Secondary CTA Text', 7),
        (v_hero_id, 'cta_secondary_link', '#demo', 'link', 'Secondary CTA Link', 8),
        (v_hero_id, 'hero_image', '', 'image', 'Hero Image URL', 9),
        (v_hero_id, 'hero_image_alt', 'LandingCMS Dashboard Preview', 'text', 'Hero Image Alt', 10),
        (v_hero_id, 'social_proof_text', 'Trusted by 10,000+ marketers and developers worldwide', 'text', 'Social Proof Text', 11),
        (v_hero_id, 'video_url', '', 'video', 'Demo Video URL', 12)
        ON CONFLICT (section_id, key) DO UPDATE SET value = EXCLUDED.value;
    END IF;
END $$;

-- ============================================================
-- SECTION CONTENT - FEATURES
-- ============================================================
DO $$
DECLARE
    v_site_id UUID;
    v_features_id UUID;
BEGIN
    SELECT id INTO v_site_id FROM sites WHERE slug = 'default' LIMIT 1;
    SELECT ps.id INTO v_features_id 
    FROM page_sections ps 
    JOIN pages p ON p.id = ps.page_id 
    WHERE p.slug = 'home' AND p.site_id = v_site_id AND ps.identifier = 'features' LIMIT 1;

    IF v_features_id IS NOT NULL THEN
        INSERT INTO section_contents (section_id, key, value, type, label, sort_order) VALUES
        (v_features_id, 'badge_text', 'Features', 'text', 'Badge Text', 1),
        (v_features_id, 'title', 'Everything you need to succeed online', 'text', 'Section Title', 2),
        (v_features_id, 'subtitle', 'Powerful tools designed for marketers, developers, and growth teams. Build faster, convert better.', 'text', 'Section Subtitle', 3)
        ON CONFLICT (section_id, key) DO UPDATE SET value = EXCLUDED.value;
    END IF;
END $$;

-- ============================================================
-- SECTION CONTENT - STATS
-- ============================================================
DO $$
DECLARE
    v_site_id UUID;
    v_stats_id UUID;
BEGIN
    SELECT id INTO v_site_id FROM sites WHERE slug = 'default' LIMIT 1;
    SELECT ps.id INTO v_stats_id 
    FROM page_sections ps 
    JOIN pages p ON p.id = ps.page_id 
    WHERE p.slug = 'home' AND p.site_id = v_site_id AND ps.identifier = 'stats' LIMIT 1;

    IF v_stats_id IS NOT NULL THEN
        INSERT INTO section_contents (section_id, key, value, type, label, sort_order) VALUES
        (v_stats_id, 'stat_1_value', '10,000+', 'text', 'Stat 1 Value', 1),
        (v_stats_id, 'stat_1_label', 'Active Users', 'text', 'Stat 1 Label', 2),
        (v_stats_id, 'stat_2_value', '99.9%', 'text', 'Stat 2 Value', 3),
        (v_stats_id, 'stat_2_label', 'Uptime SLA', 'text', 'Stat 2 Label', 4),
        (v_stats_id, 'stat_3_value', '2.3x', 'text', 'Stat 3 Value', 5),
        (v_stats_id, 'stat_3_label', 'Avg. Conversion Lift', 'text', 'Stat 3 Label', 6),
        (v_stats_id, 'stat_4_value', '< 5min', 'text', 'Stat 4 Value', 7),
        (v_stats_id, 'stat_4_label', 'Time to Launch', 'text', 'Stat 4 Label', 8)
        ON CONFLICT (section_id, key) DO UPDATE SET value = EXCLUDED.value;
    END IF;
END $$;

-- ============================================================
-- SECTION CONTENT - TESTIMONIALS
-- ============================================================
DO $$
DECLARE
    v_site_id UUID;
    v_testimonials_id UUID;
BEGIN
    SELECT id INTO v_site_id FROM sites WHERE slug = 'default' LIMIT 1;
    SELECT ps.id INTO v_testimonials_id 
    FROM page_sections ps 
    JOIN pages p ON p.id = ps.page_id 
    WHERE p.slug = 'home' AND p.site_id = v_site_id AND ps.identifier = 'testimonials' LIMIT 1;

    IF v_testimonials_id IS NOT NULL THEN
        INSERT INTO section_contents (section_id, key, value, type, label, sort_order) VALUES
        (v_testimonials_id, 'badge_text', 'Testimonials', 'text', 'Badge Text', 1),
        (v_testimonials_id, 'title', 'Loved by 10,000+ teams worldwide', 'text', 'Section Title', 2),
        (v_testimonials_id, 'subtitle', 'See what our customers are saying about LandingCMS.', 'text', 'Section Subtitle', 3)
        ON CONFLICT (section_id, key) DO UPDATE SET value = EXCLUDED.value;
    END IF;
END $$;

-- ============================================================
-- SECTION CONTENT - PRICING
-- ============================================================
DO $$
DECLARE
    v_site_id UUID;
    v_pricing_id UUID;
BEGIN
    SELECT id INTO v_site_id FROM sites WHERE slug = 'default' LIMIT 1;
    SELECT ps.id INTO v_pricing_id 
    FROM page_sections ps 
    JOIN pages p ON p.id = ps.page_id 
    WHERE p.slug = 'home' AND p.site_id = v_site_id AND ps.identifier = 'pricing' LIMIT 1;

    IF v_pricing_id IS NOT NULL THEN
        INSERT INTO section_contents (section_id, key, value, type, label, sort_order) VALUES
        (v_pricing_id, 'badge_text', 'Pricing', 'text', 'Badge Text', 1),
        (v_pricing_id, 'title', 'Simple, transparent pricing', 'text', 'Section Title', 2),
        (v_pricing_id, 'subtitle', 'Start free, scale as you grow. No hidden fees, no surprises.', 'text', 'Section Subtitle', 3)
        ON CONFLICT (section_id, key) DO UPDATE SET value = EXCLUDED.value;
    END IF;
END $$;

-- ============================================================
-- SECTION CONTENT - FAQ
-- ============================================================
DO $$
DECLARE
    v_site_id UUID;
    v_faq_id UUID;
BEGIN
    SELECT id INTO v_site_id FROM sites WHERE slug = 'default' LIMIT 1;
    SELECT ps.id INTO v_faq_id 
    FROM page_sections ps 
    JOIN pages p ON p.id = ps.page_id 
    WHERE p.slug = 'home' AND p.site_id = v_site_id AND ps.identifier = 'faq' LIMIT 1;

    IF v_faq_id IS NOT NULL THEN
        INSERT INTO section_contents (section_id, key, value, type, label, sort_order) VALUES
        (v_faq_id, 'badge_text', 'FAQ', 'text', 'Badge Text', 1),
        (v_faq_id, 'title', 'Frequently asked questions', 'text', 'Section Title', 2),
        (v_faq_id, 'subtitle', 'Everything you need to know about LandingCMS.', 'text', 'Section Subtitle', 3)
        ON CONFLICT (section_id, key) DO UPDATE SET value = EXCLUDED.value;
    END IF;
END $$;

-- ============================================================
-- SECTION CONTENT - CTA
-- ============================================================
DO $$
DECLARE
    v_site_id UUID;
    v_cta_id UUID;
BEGIN
    SELECT id INTO v_site_id FROM sites WHERE slug = 'default' LIMIT 1;
    SELECT ps.id INTO v_cta_id 
    FROM page_sections ps 
    JOIN pages p ON p.id = ps.page_id 
    WHERE p.slug = 'home' AND p.site_id = v_site_id AND ps.identifier = 'cta' LIMIT 1;

    IF v_cta_id IS NOT NULL THEN
        INSERT INTO section_contents (section_id, key, value, type, label, sort_order) VALUES
        (v_cta_id, 'title', 'Ready to build your perfect landing page?', 'text', 'CTA Title', 1),
        (v_cta_id, 'subtitle', 'Join 10,000+ teams already using LandingCMS. Start your free trial today — no credit card required.', 'text', 'CTA Subtitle', 2),
        (v_cta_id, 'cta_primary_text', 'Start Free Trial', 'text', 'Primary CTA Text', 3),
        (v_cta_id, 'cta_primary_link', '/signup', 'link', 'Primary CTA Link', 4),
        (v_cta_id, 'cta_secondary_text', 'Schedule a Demo', 'text', 'Secondary CTA Text', 5),
        (v_cta_id, 'cta_secondary_link', '/demo', 'link', 'Secondary CTA Link', 6)
        ON CONFLICT (section_id, key) DO UPDATE SET value = EXCLUDED.value;
    END IF;
END $$;

-- ============================================================
-- FEATURES (Component Table)
-- ============================================================
DO $$
DECLARE
    v_site_id UUID;
    v_features_section_id UUID;
BEGIN
    SELECT id INTO v_site_id FROM sites WHERE slug = 'default' LIMIT 1;
    SELECT ps.id INTO v_features_section_id 
    FROM page_sections ps 
    JOIN pages p ON p.id = ps.page_id 
    WHERE p.slug = 'home' AND p.site_id = v_site_id AND ps.identifier = 'features' LIMIT 1;

    -- Delete existing features for this site
    DELETE FROM features WHERE site_id = v_site_id;

    IF v_site_id IS NOT NULL THEN
        INSERT INTO features (site_id, section_id, title, description, icon, icon_color, is_active, sort_order) VALUES
        (v_site_id, v_features_section_id, 'Drag & Drop Editor', 'Build beautiful pages with our intuitive visual editor. No coding skills required.', 'zap', '#6366f1', true, 1),
        (v_site_id, v_features_section_id, 'SEO Optimized', 'Every page is built with SEO best practices. Customize meta tags, OG images, and more.', 'bar-chart', '#10b981', true, 2),
        (v_site_id, v_features_section_id, 'A/B Testing', 'Test different versions of your pages to maximize conversions automatically.', 'star', '#f59e0b', true, 3),
        (v_site_id, v_features_section_id, 'Analytics Dashboard', 'Track visitors, conversions, and revenue with our built-in analytics.', 'chart', '#3b82f6', true, 4),
        (v_site_id, v_features_section_id, 'Custom Domains', 'Connect your own domain with free SSL certificates included.', 'globe', '#8b5cf6', true, 5),
        (v_site_id, v_features_section_id, 'Team Collaboration', 'Invite your team with role-based permissions. Work together seamlessly.', 'users', '#ec4899', true, 6),
        (v_site_id, v_features_section_id, 'API & Webhooks', 'Integrate with your existing tools via our REST API and webhooks.', 'plug', '#14b8a6', true, 7),
        (v_site_id, v_features_section_id, 'Enterprise Security', 'SOC 2 Type II certified. End-to-end encryption. GDPR compliant.', 'shield', '#ef4444', true, 8),
        (v_site_id, v_features_section_id, '24/7 Support', 'Get help when you need it with our dedicated support team.', 'headphones', '#f97316', true, 9);
    END IF;
END $$;

-- ============================================================
-- TESTIMONIALS (Component Table)
-- ============================================================
DO $$
DECLARE
    v_site_id UUID;
    v_testimonials_section_id UUID;
BEGIN
    SELECT id INTO v_site_id FROM sites WHERE slug = 'default' LIMIT 1;
    SELECT ps.id INTO v_testimonials_section_id 
    FROM page_sections ps 
    JOIN pages p ON p.id = ps.page_id 
    WHERE p.slug = 'home' AND p.site_id = v_site_id AND ps.identifier = 'testimonials' LIMIT 1;

    -- Delete existing testimonials for this site
    DELETE FROM testimonials WHERE site_id = v_site_id;

    IF v_site_id IS NOT NULL THEN
        INSERT INTO testimonials (site_id, section_id, author_name, author_title, author_company, content, rating, is_featured, is_active, sort_order) VALUES
        (v_site_id, v_testimonials_section_id, 'Sarah Johnson', 'Head of Marketing', 'TechCorp Inc.', 'LandingCMS transformed our marketing workflow. We went from spending weeks on landing pages to launching in hours. Our conversion rate increased by 47% in the first month.', 5, true, true, 1),
        (v_site_id, v_testimonials_section_id, 'Michael Chen', 'CTO', 'StartupXYZ', 'As a developer, I was skeptical about no-code tools. But LandingCMS''s API and customization options won me over. It''s the perfect balance of simplicity and power.', 5, true, true, 2),
        (v_site_id, v_testimonials_section_id, 'Emily Rodriguez', 'Growth Lead', 'ScaleUp Co.', 'The A/B testing feature alone paid for our subscription 10x over. We''ve been able to optimize our pages continuously without any developer involvement.', 5, true, true, 3),
        (v_site_id, v_testimonials_section_id, 'David Kim', 'Founder', 'LaunchFast', 'I''ve tried every landing page builder out there. LandingCMS is the only one that doesn''t feel like a compromise. Fast, flexible, and the support team is incredible.', 5, false, true, 4),
        (v_site_id, v_testimonials_section_id, 'Lisa Thompson', 'VP Marketing', 'Enterprise Corp', 'We manage 50+ landing pages across 12 markets. LandingCMS''s multi-site support and team collaboration features are exactly what we needed at enterprise scale.', 5, false, true, 5),
        (v_site_id, v_testimonials_section_id, 'James Wilson', 'Product Manager', 'SaaS Platform', 'The analytics integration is seamless. We can see exactly which pages are converting and why. It''s changed how we approach our entire go-to-market strategy.', 4, false, true, 6);
    END IF;
END $$;

-- ============================================================
-- PRICING PLANS (Component Table)
-- ============================================================
DO $$
DECLARE
    v_site_id UUID;
    v_pricing_section_id UUID;
BEGIN
    SELECT id INTO v_site_id FROM sites WHERE slug = 'default' LIMIT 1;
    SELECT ps.id INTO v_pricing_section_id 
    FROM page_sections ps 
    JOIN pages p ON p.id = ps.page_id 
    WHERE p.slug = 'home' AND p.site_id = v_site_id AND ps.identifier = 'pricing' LIMIT 1;

    -- Delete existing pricing plans for this site
    DELETE FROM pricing_plans WHERE site_id = v_site_id;

    IF v_site_id IS NOT NULL THEN
        INSERT INTO pricing_plans (site_id, section_id, name, description, price_monthly, price_yearly, currency, is_popular, is_custom, badge_text, cta_text, cta_link, features, features_excluded, is_active, sort_order) VALUES
        (v_site_id, v_pricing_section_id, 'Starter', 'Perfect for individuals and small projects', 9.00, 90.00, 'USD', false, false, NULL, 'Get Started Free', '/signup?plan=starter', 
         '["1 landing page", "5,000 monthly visitors", "Basic analytics", "SSL certificate", "Custom domain", "Email support"]'::jsonb,
         '["A/B testing", "Team members", "API access", "Priority support"]'::jsonb,
         true, 1),
        (v_site_id, v_pricing_section_id, 'Pro', 'For growing businesses and marketing teams', 29.00, 290.00, 'USD', true, false, 'Most Popular', 'Start Free Trial', '/signup?plan=pro',
         '["Unlimited landing pages", "50,000 monthly visitors", "Advanced analytics", "A/B testing", "3 team members", "API access", "Priority support", "Custom domain", "SSL certificate", "Webhooks"]'::jsonb,
         '["Dedicated account manager", "Custom SLA", "SSO/SAML"]'::jsonb,
         true, 2),
        (v_site_id, v_pricing_section_id, 'Business', 'For scaling companies with advanced needs', 79.00, 790.00, 'USD', false, false, 'Best Value', 'Start Free Trial', '/signup?plan=business',
         '["Unlimited landing pages", "500,000 monthly visitors", "Advanced analytics", "Unlimited A/B testing", "10 team members", "API access", "Priority support", "Custom domain", "SSL certificate", "Webhooks", "White-label option", "Dedicated account manager"]'::jsonb,
         '["Custom SLA", "SSO/SAML"]'::jsonb,
         true, 3),
        (v_site_id, v_pricing_section_id, 'Enterprise', 'For large organizations with custom requirements', NULL, NULL, 'USD', false, true, NULL, 'Contact Sales', '/contact',
         '["Unlimited everything", "Custom visitor limits", "Custom analytics", "Unlimited A/B testing", "Unlimited team members", "Full API access", "24/7 dedicated support", "Custom domain", "SSL certificate", "Webhooks", "White-label", "Dedicated account manager", "Custom SLA", "SSO/SAML", "On-premise option"]'::jsonb,
         '[]'::jsonb,
         true, 4);
    END IF;
END $$;

-- ============================================================
-- FAQs (Component Table)
-- ============================================================
DO $$
DECLARE
    v_site_id UUID;
    v_faq_section_id UUID;
BEGIN
    SELECT id INTO v_site_id FROM sites WHERE slug = 'default' LIMIT 1;
    SELECT ps.id INTO v_faq_section_id 
    FROM page_sections ps 
    JOIN pages p ON p.id = ps.page_id 
    WHERE p.slug = 'home' AND p.site_id = v_site_id AND ps.identifier = 'faq' LIMIT 1;

    -- Delete existing FAQs for this site
    DELETE FROM faqs WHERE site_id = v_site_id;

    IF v_site_id IS NOT NULL THEN
        INSERT INTO faqs (site_id, section_id, question, answer, category, is_active, sort_order) VALUES
        (v_site_id, v_faq_section_id, 'How do I get started with LandingCMS?', 'Getting started is easy! Sign up for a free account, choose a template or start from scratch, and you can have your first landing page live in under 5 minutes. No credit card required for the free trial.', 'Getting Started', true, 1),
        (v_site_id, v_faq_section_id, 'Do I need coding skills to use LandingCMS?', 'Not at all! LandingCMS is designed for marketers and non-technical users. Our drag-and-drop editor makes it easy to create professional landing pages without writing a single line of code. However, if you are a developer, you can also use our API for advanced customizations.', 'Getting Started', true, 2),
        (v_site_id, v_faq_section_id, 'Can I use my own domain name?', 'Yes! You can connect your custom domain on all paid plans. We provide free SSL certificates for all custom domains, so your pages will be secure and trusted by visitors. Setup takes less than 5 minutes.', 'Features', true, 3),
        (v_site_id, v_faq_section_id, 'How does A/B testing work?', 'Our A/B testing feature lets you create multiple versions of your landing page and automatically split traffic between them. We track conversions and show you which version performs better. You can then publish the winner with one click.', 'Features', true, 4),
        (v_site_id, v_faq_section_id, 'Can I cancel my subscription at any time?', 'Yes, you can cancel your subscription at any time with no questions asked. If you cancel, you will retain access until the end of your billing period. We do not charge cancellation fees.', 'Billing', true, 5),
        (v_site_id, v_faq_section_id, 'What payment methods do you accept?', 'We accept all major credit cards (Visa, MasterCard, American Express, Discover), PayPal, and bank transfers for annual Enterprise plans. All payments are processed securely through Stripe.', 'Billing', true, 6),
        (v_site_id, v_faq_section_id, 'Is my data secure?', 'Absolutely. We take security seriously. LandingCMS is SOC 2 Type II certified, GDPR compliant, and all data is encrypted at rest and in transit using AES-256 encryption. We perform regular security audits and penetration testing.', 'Security', true, 7),
        (v_site_id, v_faq_section_id, 'Do you offer a free trial?', 'Yes! All plans come with a 14-day free trial with full access to all features. No credit card required to start. After the trial, you can choose the plan that best fits your needs.', 'Billing', true, 8),
        (v_site_id, v_faq_section_id, 'Can I migrate my existing landing pages?', 'Yes, we offer a free migration service for Business and Enterprise plans. Our team will help you migrate your existing landing pages from any platform. For Starter and Pro plans, we provide detailed migration guides.', 'Getting Started', true, 9),
        (v_site_id, v_faq_section_id, 'What integrations do you support?', 'LandingCMS integrates with 100+ tools including Google Analytics, HubSpot, Salesforce, Mailchimp, Zapier, Slack, and many more. We also offer a REST API and webhooks for custom integrations.', 'Features', true, 10);
    END IF;
END $$;

-- ============================================================
-- SOCIAL LINKS
-- ============================================================
DO $$
DECLARE
    v_site_id UUID;
BEGIN
    SELECT id INTO v_site_id FROM sites WHERE slug = 'default' LIMIT 1;
    
    IF v_site_id IS NOT NULL THEN
        -- Delete existing social links
        DELETE FROM social_links WHERE site_id = v_site_id;
        
        INSERT INTO social_links (site_id, platform, url, icon, label, is_active, sort_order) VALUES
        (v_site_id, 'twitter', 'https://twitter.com/landingcms', 'twitter', 'Twitter / X', true, 1),
        (v_site_id, 'linkedin', 'https://linkedin.com/company/landingcms', 'linkedin', 'LinkedIn', true, 2),
        (v_site_id, 'github', 'https://github.com/landingcms', 'github', 'GitHub', true, 3),
        (v_site_id, 'youtube', 'https://youtube.com/@landingcms', 'youtube', 'YouTube', true, 4);
    END IF;
END $$;

-- ============================================================
-- UPDATE PAGE SEO
-- ============================================================
DO $$
DECLARE
    v_site_id UUID;
BEGIN
    SELECT id INTO v_site_id FROM sites WHERE slug = 'default' LIMIT 1;
    
    IF v_site_id IS NOT NULL THEN
        UPDATE pages SET
            seo_title = 'LandingCMS - Build Dynamic Landing Pages That Convert',
            seo_description = 'Create, manage, and optimize your landing pages with our powerful CMS. No coding required. A/B testing, analytics, and 100+ integrations included.',
            seo_keywords = 'landing page builder, cms, no-code, website builder, marketing, conversion optimization, A/B testing',
            og_title = 'LandingCMS - Build Landing Pages That Convert',
            og_description = 'Create stunning landing pages without code. Manage all content from one powerful dashboard. Start free today.',
            twitter_title = 'LandingCMS - Build Landing Pages That Convert',
            twitter_description = 'Create stunning landing pages without code. Start free today.',
            twitter_card = 'summary_large_image',
            robots_meta = 'index, follow'
        WHERE slug = 'home' AND site_id = v_site_id;
    END IF;
END $$;

-- Record migration
INSERT INTO schema_migrations (version, description) VALUES
('009', 'Seed landing page content')
ON CONFLICT DO NOTHING;
