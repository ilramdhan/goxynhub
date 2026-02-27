package domain

import (
	"time"

	"github.com/google/uuid"
)

// PageStatus represents the publication status of a page
type PageStatus string

const (
	PageStatusDraft     PageStatus = "draft"
	PageStatusPublished PageStatus = "published"
	PageStatusArchived  PageStatus = "archived"
)

// SectionType represents the type of a page section
type SectionType string

const (
	SectionTypeHero         SectionType = "hero"
	SectionTypeFeatures     SectionType = "features"
	SectionTypePricing      SectionType = "pricing"
	SectionTypeTestimonials SectionType = "testimonials"
	SectionTypeFAQ          SectionType = "faq"
	SectionTypeCTA          SectionType = "cta"
	SectionTypeAbout        SectionType = "about"
	SectionTypeTeam         SectionType = "team"
	SectionTypeGallery      SectionType = "gallery"
	SectionTypeStats        SectionType = "stats"
	SectionTypeLogos        SectionType = "logos"
	SectionTypeNewsletter   SectionType = "newsletter"
	SectionTypeContact      SectionType = "contact"
	SectionTypeVideo        SectionType = "video"
	SectionTypeCustom       SectionType = "custom"
	SectionTypeHTML         SectionType = "html"
)

// ContentType represents the type of section content
type ContentType string

const (
	ContentTypeText     ContentType = "text"
	ContentTypeHTML     ContentType = "html"
	ContentTypeMarkdown ContentType = "markdown"
	ContentTypeImage    ContentType = "image"
	ContentTypeVideo    ContentType = "video"
	ContentTypeLink     ContentType = "link"
	ContentTypeButton   ContentType = "button"
	ContentTypeColor    ContentType = "color"
	ContentTypeNumber   ContentType = "number"
	ContentTypeBoolean  ContentType = "boolean"
	ContentTypeJSON     ContentType = "json"
	ContentTypeFile     ContentType = "file"
)

// Page represents a landing page
type Page struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	SiteID      uuid.UUID  `db:"site_id" json:"site_id"`
	Title       string     `db:"title" json:"title"`
	Slug        string     `db:"slug" json:"slug"`
	Description *string    `db:"description" json:"description"`
	Status      PageStatus `db:"status" json:"status"`
	IsHomepage  bool       `db:"is_homepage" json:"is_homepage"`
	// SEO
	SEOTitle       *string `db:"seo_title" json:"seo_title"`
	SEODescription *string `db:"seo_description" json:"seo_description"`
	SEOKeywords    *string `db:"seo_keywords" json:"seo_keywords"`
	// Open Graph
	OGTitle       *string `db:"og_title" json:"og_title"`
	OGDescription *string `db:"og_description" json:"og_description"`
	OGImage       *string `db:"og_image" json:"og_image"`
	OGType        *string `db:"og_type" json:"og_type"`
	// Twitter Card
	TwitterTitle       *string `db:"twitter_title" json:"twitter_title"`
	TwitterDescription *string `db:"twitter_description" json:"twitter_description"`
	TwitterImage       *string `db:"twitter_image" json:"twitter_image"`
	TwitterCard        *string `db:"twitter_card" json:"twitter_card"`
	// Schema
	SchemaMarkup JSONMap `db:"schema_markup" json:"schema_markup,omitempty"`
	// Other
	CustomHead   *string    `db:"custom_head" json:"custom_head"`
	CanonicalURL *string    `db:"canonical_url" json:"canonical_url"`
	RobotsMeta   *string    `db:"robots_meta" json:"robots_meta"`
	Template     *string    `db:"template" json:"template"`
	SortOrder    int        `db:"sort_order" json:"sort_order"`
	PublishedAt  *time.Time `db:"published_at" json:"published_at"`
	Metadata     JSONMap    `db:"metadata" json:"metadata,omitempty"`
	CreatedBy    *uuid.UUID `db:"created_by" json:"created_by"`
	UpdatedBy    *uuid.UUID `db:"updated_by" json:"updated_by"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at" json:"-"`
	// Relations (populated on demand)
	Sections []*PageSection `db:"-" json:"sections,omitempty"`
}

// PageSection represents a section within a page
type PageSection struct {
	ID         uuid.UUID   `db:"id" json:"id"`
	PageID     uuid.UUID   `db:"page_id" json:"page_id"`
	Name       string      `db:"name" json:"name"`
	Type       SectionType `db:"type" json:"type"`
	Identifier *string     `db:"identifier" json:"identifier"`
	IsVisible  bool        `db:"is_visible" json:"is_visible"`
	SortOrder  int         `db:"sort_order" json:"sort_order"`
	// Background
	BGColor          *string  `db:"bg_color" json:"bg_color"`
	BGImage          *string  `db:"bg_image" json:"bg_image"`
	BGVideo          *string  `db:"bg_video" json:"bg_video"`
	BGOverlay        bool     `db:"bg_overlay" json:"bg_overlay"`
	BGOverlayColor   *string  `db:"bg_overlay_color" json:"bg_overlay_color"`
	BGOverlayOpacity *float64 `db:"bg_overlay_opacity" json:"bg_overlay_opacity"`
	// Layout
	Layout        *string `db:"layout" json:"layout"`
	PaddingTop    *string `db:"padding_top" json:"padding_top"`
	PaddingBottom *string `db:"padding_bottom" json:"padding_bottom"`
	// Style
	Animation *string `db:"animation" json:"animation"`
	CSSClass  *string `db:"css_class" json:"css_class"`
	CustomCSS *string `db:"custom_css" json:"custom_css"`
	Metadata  JSONMap `db:"metadata" json:"metadata,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	// Relations
	Contents []*SectionContent `db:"-" json:"contents,omitempty"`
}

// SectionContent represents a key-value content item within a section
type SectionContent struct {
	ID          uuid.UUID   `db:"id" json:"id"`
	SectionID   uuid.UUID   `db:"section_id" json:"section_id"`
	Key         string      `db:"key" json:"key"`
	Value       *string     `db:"value" json:"value"`
	ValueJSON   JSONMap     `db:"value_json" json:"value_json,omitempty"`
	Type        ContentType `db:"type" json:"type"`
	Label       *string     `db:"label" json:"label"`
	Description *string     `db:"description" json:"description"`
	Placeholder *string     `db:"placeholder" json:"placeholder"`
	IsRequired  bool        `db:"is_required" json:"is_required"`
	SortOrder   int         `db:"sort_order" json:"sort_order"`
	// Image/file specific
	AltText *string `db:"alt_text" json:"alt_text"`
	Width   *int    `db:"width" json:"width"`
	Height  *int    `db:"height" json:"height"`
	// Link/button specific
	LinkURL    *string `db:"link_url" json:"link_url"`
	LinkTarget *string `db:"link_target" json:"link_target"`
	Metadata   JSONMap `db:"metadata" json:"metadata,omitempty"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

// PageFilter holds filter parameters for page queries
type PageFilter struct {
	SiteID     *uuid.UUID
	Status     *PageStatus
	IsHomepage *bool
	Search     *string
	Pagination
}

// CreatePageInput holds data for creating a new page
type CreatePageInput struct {
	SiteID      uuid.UUID  `json:"site_id" validate:"required"`
	Title       string     `json:"title" validate:"required,min=1,max=255"`
	Slug        string     `json:"slug" validate:"required,min=1,max=255"`
	Description *string    `json:"description" validate:"omitempty,max=500"`
	Status      PageStatus `json:"status" validate:"required,oneof=draft published archived"`
	IsHomepage  bool       `json:"is_homepage"`
	// SEO
	SEOTitle       *string `json:"seo_title" validate:"omitempty,max=255"`
	SEODescription *string `json:"seo_description" validate:"omitempty,max=500"`
	SEOKeywords    *string `json:"seo_keywords"`
	// OG
	OGTitle       *string `json:"og_title" validate:"omitempty,max=255"`
	OGDescription *string `json:"og_description" validate:"omitempty,max=500"`
	OGImage       *string `json:"og_image" validate:"omitempty,url"`
	OGType        *string `json:"og_type"`
	// Twitter
	TwitterTitle       *string `json:"twitter_title" validate:"omitempty,max=255"`
	TwitterDescription *string `json:"twitter_description" validate:"omitempty,max=500"`
	TwitterImage       *string `json:"twitter_image" validate:"omitempty,url"`
	TwitterCard        *string `json:"twitter_card"`
	// Other
	CanonicalURL *string `json:"canonical_url" validate:"omitempty,url"`
	RobotsMeta   *string `json:"robots_meta"`
	Template     *string `json:"template"`
}

// UpdatePageInput holds data for updating a page
type UpdatePageInput struct {
	Title       *string     `json:"title" validate:"omitempty,min=1,max=255"`
	Slug        *string     `json:"slug" validate:"omitempty,min=1,max=255"`
	Description *string     `json:"description" validate:"omitempty,max=500"`
	Status      *PageStatus `json:"status" validate:"omitempty,oneof=draft published archived"`
	IsHomepage  *bool       `json:"is_homepage"`
	// SEO
	SEOTitle       *string `json:"seo_title" validate:"omitempty,max=255"`
	SEODescription *string `json:"seo_description" validate:"omitempty,max=500"`
	SEOKeywords    *string `json:"seo_keywords"`
	// OG
	OGTitle       *string `json:"og_title" validate:"omitempty,max=255"`
	OGDescription *string `json:"og_description" validate:"omitempty,max=500"`
	OGImage       *string `json:"og_image" validate:"omitempty,url"`
	OGType        *string `json:"og_type"`
	// Twitter
	TwitterTitle       *string `json:"twitter_title" validate:"omitempty,max=255"`
	TwitterDescription *string `json:"twitter_description" validate:"omitempty,max=500"`
	TwitterImage       *string `json:"twitter_image" validate:"omitempty,url"`
	TwitterCard        *string `json:"twitter_card"`
	// Other
	CanonicalURL *string `json:"canonical_url" validate:"omitempty,url"`
	RobotsMeta   *string `json:"robots_meta"`
	Template     *string `json:"template"`
	CustomHead   *string `json:"custom_head"`
}

// CreateSectionInput holds data for creating a page section
type CreateSectionInput struct {
	PageID     uuid.UUID   `json:"page_id" validate:"required"`
	Name       string      `json:"name" validate:"required,min=1,max=255"`
	Type       SectionType `json:"type" validate:"required"`
	Identifier *string     `json:"identifier" validate:"omitempty,max=100"`
	IsVisible  bool        `json:"is_visible"`
	SortOrder  int         `json:"sort_order"`
}

// UpdateSectionInput holds data for updating a page section
type UpdateSectionInput struct {
	Name             *string     `json:"name" validate:"omitempty,min=1,max=255"`
	Type             *SectionType `json:"type"`
	IsVisible        *bool       `json:"is_visible"`
	SortOrder        *int        `json:"sort_order"`
	BGColor          *string     `json:"bg_color"`
	BGImage          *string     `json:"bg_image"`
	BGVideo          *string     `json:"bg_video"`
	BGOverlay        *bool       `json:"bg_overlay"`
	BGOverlayColor   *string     `json:"bg_overlay_color"`
	BGOverlayOpacity *float64    `json:"bg_overlay_opacity"`
	Layout           *string     `json:"layout"`
	PaddingTop       *string     `json:"padding_top"`
	PaddingBottom    *string     `json:"padding_bottom"`
	Animation        *string     `json:"animation"`
	CSSClass         *string     `json:"css_class"`
	CustomCSS        *string     `json:"custom_css"`
}

// UpsertContentInput holds data for creating/updating section content
type UpsertContentInput struct {
	Key         string      `json:"key" validate:"required,min=1,max=255"`
	Value       *string     `json:"value"`
	ValueJSON   JSONMap     `json:"value_json"`
	Type        ContentType `json:"type" validate:"required"`
	Label       *string     `json:"label"`
	Description *string     `json:"description"`
	Placeholder *string     `json:"placeholder"`
	IsRequired  bool        `json:"is_required"`
	SortOrder   int         `json:"sort_order"`
	AltText     *string     `json:"alt_text"`
	Width       *int        `json:"width"`
	Height      *int        `json:"height"`
	LinkURL     *string     `json:"link_url"`
	LinkTarget  *string     `json:"link_target"`
}

// ReorderSectionsInput holds data for reordering sections
type ReorderSectionsInput struct {
	Sections []SectionOrder `json:"sections" validate:"required,min=1"`
}

// SectionOrder holds the ID and new sort order for a section
type SectionOrder struct {
	ID        uuid.UUID `json:"id" validate:"required"`
	SortOrder int       `json:"sort_order"`
}
