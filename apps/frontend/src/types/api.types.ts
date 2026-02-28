// API Response types

export interface ApiResponse<T = unknown> {
  success: boolean;
  message: string;
  data?: T;
  errors?: ValidationError[];
  meta?: PaginationMeta;
}

export interface ValidationError {
  field: string;
  message: string;
}

export interface PaginationMeta {
  page: number;
  per_page: number;
  total: number;
  total_pages: number;
}

export interface PaginatedResponse<T> {
  data: T[];
  meta: PaginationMeta;
}

// Auth types
export type UserRole = "super_admin" | "admin" | "editor";
export type UserStatus = "active" | "inactive" | "suspended";

export interface User {
  id: string;
  email: string;
  full_name: string;
  avatar_url: string | null;
  role: UserRole;
  status: UserStatus;
  email_verified: boolean;
  last_login_at: string | null;
  created_at: string;
  updated_at: string;
}

export interface AuthTokens {
  access_token: string;
  expires_at: string;
  token_type: string;
}

export interface LoginResponse {
  access_token: string;
  expires_at: string;
  token_type: string;
  user: Pick<User, "id" | "email" | "full_name" | "role" | "avatar_url">;
}

export interface LoginInput {
  email: string;
  password: string;
}

export interface ChangePasswordInput {
  current_password: string;
  new_password: string;
}

// Site types
export interface Site {
  id: string;
  name: string;
  slug: string;
  domain: string | null;
  description: string | null;
  logo_url: string | null;
  favicon_url: string | null;
  is_active: boolean;
  created_at: string;
  updated_at: string;
  settings?: SiteSetting[];
}

export interface SiteSetting {
  id: string;
  site_id: string;
  key: string;
  value: string | null;
  value_json: Record<string, unknown> | null;
  type: string;
  group_name: string;
  label: string | null;
  description: string | null;
  is_public: boolean;
  sort_order: number;
}

export type SiteSettingsMap = Record<string, string>;

// Page types
export type PageStatus = "draft" | "published" | "archived";
export type SectionType =
  | "hero"
  | "features"
  | "pricing"
  | "testimonials"
  | "faq"
  | "cta"
  | "about"
  | "team"
  | "gallery"
  | "stats"
  | "logos"
  | "newsletter"
  | "contact"
  | "video"
  | "custom"
  | "html";

export type ContentType =
  | "text"
  | "html"
  | "markdown"
  | "image"
  | "video"
  | "link"
  | "button"
  | "color"
  | "number"
  | "boolean"
  | "json"
  | "file";

export interface Page {
  id: string;
  site_id: string;
  title: string;
  slug: string;
  description: string | null;
  status: PageStatus;
  is_homepage: boolean;
  // SEO
  seo_title: string | null;
  seo_description: string | null;
  seo_keywords: string | null;
  // OG
  og_title: string | null;
  og_description: string | null;
  og_image: string | null;
  og_type: string | null;
  // Twitter
  twitter_title: string | null;
  twitter_description: string | null;
  twitter_image: string | null;
  twitter_card: string | null;
  // Other
  custom_head: string | null;
  canonical_url: string | null;
  robots_meta: string | null;
  template: string | null;
  sort_order: number;
  published_at: string | null;
  created_at: string;
  updated_at: string;
  sections?: PageSection[];
}

export interface PageSection {
  id: string;
  page_id: string;
  name: string;
  type: SectionType;
  identifier: string | null;
  is_visible: boolean;
  sort_order: number;
  // Background
  bg_color: string | null;
  bg_image: string | null;
  bg_video: string | null;
  bg_overlay: boolean;
  bg_overlay_color: string | null;
  bg_overlay_opacity: number | null;
  // Layout
  layout: string | null;
  padding_top: string | null;
  padding_bottom: string | null;
  // Style
  animation: string | null;
  css_class: string | null;
  custom_css: string | null;
  created_at: string;
  updated_at: string;
  contents?: SectionContent[];
}

export interface SectionContent {
  id: string;
  section_id: string;
  key: string;
  value: string | null;
  value_json: Record<string, unknown> | null;
  type: ContentType;
  label: string | null;
  description: string | null;
  placeholder: string | null;
  is_required: boolean;
  sort_order: number;
  alt_text: string | null;
  width: number | null;
  height: number | null;
  link_url: string | null;
  link_target: string | null;
  created_at: string;
  updated_at: string;
}

// Content map for easy access by key
export type ContentMap = Record<string, SectionContent>;

// Helper to convert section contents array to map
export function contentsToMap(contents: SectionContent[]): ContentMap {
  return contents.reduce((acc, content) => {
    acc[content.key] = content;
    return acc;
  }, {} as ContentMap);
}

// Helper to get content value
export function getContentValue(
  contents: ContentMap,
  key: string,
  defaultValue = ""
): string {
  return contents[key]?.value ?? defaultValue;
}

// Component types
export interface Feature {
  id: string;
  site_id: string;
  section_id: string | null;
  title: string;
  description: string | null;
  icon: string | null;
  icon_color: string | null;
  image_url: string | null;
  image_alt: string | null;
  link_url: string | null;
  link_text: string | null;
  is_active: boolean;
  sort_order: number;
  created_at: string;
  updated_at: string;
}

export interface Testimonial {
  id: string;
  site_id: string;
  section_id: string | null;
  author_name: string;
  author_title: string | null;
  author_company: string | null;
  author_avatar: string | null;
  content: string;
  rating: number | null;
  source: string | null;
  source_url: string | null;
  is_featured: boolean;
  is_active: boolean;
  sort_order: number;
  created_at: string;
  updated_at: string;
}

export interface PricingPlan {
  id: string;
  site_id: string;
  section_id: string | null;
  name: string;
  description: string | null;
  price_monthly: number | null;
  price_yearly: number | null;
  currency: string;
  price_label: string | null;
  is_popular: boolean;
  is_custom: boolean;
  badge_text: string | null;
  cta_text: string;
  cta_link: string | null;
  features: string[];
  features_excluded: string[];
  is_active: boolean;
  sort_order: number;
  created_at: string;
  updated_at: string;
}

export interface FAQ {
  id: string;
  site_id: string;
  section_id: string | null;
  question: string;
  answer: string;
  category: string | null;
  is_active: boolean;
  sort_order: number;
  created_at: string;
  updated_at: string;
}

export interface NavigationMenu {
  id: string;
  site_id: string;
  name: string;
  identifier: string;
  description: string | null;
  is_active: boolean;
  items?: NavigationItem[];
}

export interface NavigationItem {
  id: string;
  menu_id: string;
  parent_id: string | null;
  page_id: string | null;
  label: string;
  url: string | null;
  target: string;
  icon: string | null;
  css_class: string | null;
  is_active: boolean;
  is_mega_menu: boolean;
  sort_order: number;
  depth: number;
  children?: NavigationItem[];
}

// Media types
export type MediaType = "image" | "video" | "document" | "audio" | "other";

export interface Media {
  id: string;
  site_id: string;
  name: string;
  original_name: string;
  file_path: string;
  public_url: string;
  thumbnail_url: string | null;
  type: MediaType;
  mime_type: string;
  file_size: number;
  width: number | null;
  height: number | null;
  duration: number | null;
  alt_text: string | null;
  caption: string | null;
  tags: string[];
  folder: string;
  created_at: string;
  updated_at: string;
}

// Audit log types
export interface AuditLog {
  id: string;
  user_id: string | null;
  user_email: string | null;
  user_role: string | null;
  action: string;
  resource_type: string;
  resource_id: string | null;
  resource_name: string | null;
  old_values: Record<string, unknown> | null;
  new_values: Record<string, unknown> | null;
  ip_address: string | null;
  site_id: string | null;
  created_at: string;
}

// Query parameter types
export interface PageQueryParams {
  page?: number;
  per_page?: number;
  site_id?: string;
  status?: PageStatus;
  search?: string;
}

export interface MediaQueryParams {
  page?: number;
  per_page?: number;
  type?: MediaType;
  folder?: string;
  search?: string;
}
