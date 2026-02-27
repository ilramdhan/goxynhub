-- Migration: 006_create_navigation.sql
-- Description: Create navigation menus and items tables
-- Created: 2024-01-01

-- Navigation menus table
CREATE TABLE IF NOT EXISTS navigation_menus (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    site_id     UUID NOT NULL REFERENCES sites(id) ON DELETE CASCADE,
    name        VARCHAR(255) NOT NULL,
    identifier  VARCHAR(100) NOT NULL,  -- e.g., 'header', 'footer', 'sidebar'
    description TEXT,
    is_active   BOOLEAN NOT NULL DEFAULT true,
    metadata    JSONB DEFAULT '{}',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_nav_menus_site_identifier ON navigation_menus(site_id, identifier);
CREATE INDEX idx_nav_menus_site_id ON navigation_menus(site_id);

-- Navigation items table (supports nested/hierarchical menus)
CREATE TABLE IF NOT EXISTS navigation_items (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    menu_id         UUID NOT NULL REFERENCES navigation_menus(id) ON DELETE CASCADE,
    parent_id       UUID REFERENCES navigation_items(id) ON DELETE CASCADE,
    page_id         UUID REFERENCES pages(id) ON DELETE SET NULL,  -- link to internal page
    label           VARCHAR(255) NOT NULL,
    url             TEXT,                   -- external URL or path
    target          VARCHAR(20) DEFAULT '_self',  -- _self, _blank
    icon            VARCHAR(100),           -- icon name/class
    css_class       VARCHAR(255),
    is_active       BOOLEAN NOT NULL DEFAULT true,
    is_mega_menu    BOOLEAN NOT NULL DEFAULT false,
    sort_order      INTEGER NOT NULL DEFAULT 0,
    depth           INTEGER NOT NULL DEFAULT 0,  -- 0 = top level
    metadata        JSONB DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_nav_items_menu_id ON navigation_items(menu_id);
CREATE INDEX idx_nav_items_parent_id ON navigation_items(parent_id);
CREATE INDEX idx_nav_items_page_id ON navigation_items(page_id);
CREATE INDEX idx_nav_items_sort ON navigation_items(menu_id, sort_order);

-- Apply triggers
CREATE TRIGGER update_navigation_menus_updated_at
    BEFORE UPDATE ON navigation_menus
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_navigation_items_updated_at
    BEFORE UPDATE ON navigation_items
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Insert default navigation
DO $$
DECLARE
    v_site_id UUID;
    v_header_menu_id UUID;
    v_footer_menu_id UUID;
    v_home_page_id UUID;
BEGIN
    SELECT id INTO v_site_id FROM sites WHERE slug = 'default' LIMIT 1;
    SELECT id INTO v_home_page_id FROM pages WHERE slug = 'home' LIMIT 1;
    
    IF v_site_id IS NOT NULL THEN
        -- Header menu
        INSERT INTO navigation_menus (site_id, name, identifier, description)
        VALUES (v_site_id, 'Header Navigation', 'header', 'Main header navigation menu')
        RETURNING id INTO v_header_menu_id;
        
        -- Footer menu
        INSERT INTO navigation_menus (site_id, name, identifier, description)
        VALUES (v_site_id, 'Footer Navigation', 'footer', 'Footer navigation menu')
        RETURNING id INTO v_footer_menu_id;
        
        -- Header items
        INSERT INTO navigation_items (menu_id, page_id, label, url, sort_order) VALUES
        (v_header_menu_id, v_home_page_id, 'Home', '/', 1),
        (v_header_menu_id, NULL, 'Features', '/#features', 2),
        (v_header_menu_id, NULL, 'Pricing', '/#pricing', 3),
        (v_header_menu_id, NULL, 'FAQ', '/#faq', 4),
        (v_header_menu_id, NULL, 'Contact', '/contact', 5);
        
        -- Footer items
        INSERT INTO navigation_items (menu_id, label, url, sort_order) VALUES
        (v_footer_menu_id, 'Privacy Policy', '/privacy', 1),
        (v_footer_menu_id, 'Terms of Service', '/terms', 2),
        (v_footer_menu_id, 'Cookie Policy', '/cookies', 3);
    END IF;
END $$;

-- ============================================================
-- ROLLBACK SCRIPT
-- ============================================================
-- DROP TABLE IF EXISTS navigation_items CASCADE;
-- DROP TABLE IF EXISTS navigation_menus CASCADE;
