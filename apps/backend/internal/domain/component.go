package domain

import (
	"time"

	"github.com/google/uuid"
)

// Feature represents a feature item
type Feature struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	SiteID      uuid.UUID  `db:"site_id" json:"site_id"`
	SectionID   *uuid.UUID `db:"section_id" json:"section_id"`
	Title       string     `db:"title" json:"title"`
	Description *string    `db:"description" json:"description"`
	Icon        *string    `db:"icon" json:"icon"`
	IconColor   *string    `db:"icon_color" json:"icon_color"`
	ImageURL    *string    `db:"image_url" json:"image_url"`
	ImageAlt    *string    `db:"image_alt" json:"image_alt"`
	LinkURL     *string    `db:"link_url" json:"link_url"`
	LinkText    *string    `db:"link_text" json:"link_text"`
	IsActive    bool       `db:"is_active" json:"is_active"`
	SortOrder   int        `db:"sort_order" json:"sort_order"`
	Metadata    JSONMap    `db:"metadata" json:"metadata,omitempty"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at" json:"-"`
}

// Testimonial represents a testimonial/review
type Testimonial struct {
	ID            uuid.UUID  `db:"id" json:"id"`
	SiteID        uuid.UUID  `db:"site_id" json:"site_id"`
	SectionID     *uuid.UUID `db:"section_id" json:"section_id"`
	AuthorName    string     `db:"author_name" json:"author_name"`
	AuthorTitle   *string    `db:"author_title" json:"author_title"`
	AuthorCompany *string    `db:"author_company" json:"author_company"`
	AuthorAvatar  *string    `db:"author_avatar" json:"author_avatar"`
	Content       string     `db:"content" json:"content"`
	Rating        *int       `db:"rating" json:"rating"`
	Source        *string    `db:"source" json:"source"`
	SourceURL     *string    `db:"source_url" json:"source_url"`
	IsFeatured    bool       `db:"is_featured" json:"is_featured"`
	IsActive      bool       `db:"is_active" json:"is_active"`
	SortOrder     int        `db:"sort_order" json:"sort_order"`
	Metadata      JSONMap    `db:"metadata" json:"metadata,omitempty"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt     *time.Time `db:"deleted_at" json:"-"`
}

// PricingPlan represents a pricing plan
type PricingPlan struct {
	ID               uuid.UUID  `db:"id" json:"id"`
	SiteID           uuid.UUID  `db:"site_id" json:"site_id"`
	SectionID        *uuid.UUID `db:"section_id" json:"section_id"`
	Name             string     `db:"name" json:"name"`
	Description      *string    `db:"description" json:"description"`
	PriceMonthly     *float64   `db:"price_monthly" json:"price_monthly"`
	PriceYearly      *float64   `db:"price_yearly" json:"price_yearly"`
	Currency         string     `db:"currency" json:"currency"`
	PriceLabel       *string    `db:"price_label" json:"price_label"`
	IsPopular        bool       `db:"is_popular" json:"is_popular"`
	IsCustom         bool       `db:"is_custom" json:"is_custom"`
	BadgeText        *string    `db:"badge_text" json:"badge_text"`
	CTAText          string     `db:"cta_text" json:"cta_text"`
	CTALink          *string    `db:"cta_link" json:"cta_link"`
	Features         JSONArray  `db:"features" json:"features"`
	FeaturesExcluded JSONArray  `db:"features_excluded" json:"features_excluded"`
	IsActive         bool       `db:"is_active" json:"is_active"`
	SortOrder        int        `db:"sort_order" json:"sort_order"`
	Metadata         JSONMap    `db:"metadata" json:"metadata,omitempty"`
	CreatedAt        time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt        *time.Time `db:"deleted_at" json:"-"`
}

// FAQ represents a frequently asked question
type FAQ struct {
	ID        uuid.UUID  `db:"id" json:"id"`
	SiteID    uuid.UUID  `db:"site_id" json:"site_id"`
	SectionID *uuid.UUID `db:"section_id" json:"section_id"`
	Question  string     `db:"question" json:"question"`
	Answer    string     `db:"answer" json:"answer"`
	Category  *string    `db:"category" json:"category"`
	IsActive  bool       `db:"is_active" json:"is_active"`
	SortOrder int        `db:"sort_order" json:"sort_order"`
	Metadata  JSONMap    `db:"metadata" json:"metadata,omitempty"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"-"`
}

// NavigationMenu represents a navigation menu
type NavigationMenu struct {
	ID          uuid.UUID         `db:"id" json:"id"`
	SiteID      uuid.UUID         `db:"site_id" json:"site_id"`
	Name        string            `db:"name" json:"name"`
	Identifier  string            `db:"identifier" json:"identifier"`
	Description *string           `db:"description" json:"description"`
	IsActive    bool              `db:"is_active" json:"is_active"`
	Metadata    JSONMap           `db:"metadata" json:"metadata,omitempty"`
	CreatedAt   time.Time         `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time         `db:"updated_at" json:"updated_at"`
	Items       []*NavigationItem `db:"-" json:"items,omitempty"`
}

// NavigationItem represents a navigation menu item
type NavigationItem struct {
	ID          uuid.UUID         `db:"id" json:"id"`
	MenuID      uuid.UUID         `db:"menu_id" json:"menu_id"`
	ParentID    *uuid.UUID        `db:"parent_id" json:"parent_id"`
	PageID      *uuid.UUID        `db:"page_id" json:"page_id"`
	Label       string            `db:"label" json:"label"`
	URL         *string           `db:"url" json:"url"`
	Target      string            `db:"target" json:"target"`
	Icon        *string           `db:"icon" json:"icon"`
	CSSClass    *string           `db:"css_class" json:"css_class"`
	IsActive    bool              `db:"is_active" json:"is_active"`
	IsMegaMenu  bool              `db:"is_mega_menu" json:"is_mega_menu"`
	SortOrder   int               `db:"sort_order" json:"sort_order"`
	Depth       int               `db:"depth" json:"depth"`
	Metadata    JSONMap           `db:"metadata" json:"metadata,omitempty"`
	CreatedAt   time.Time         `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time         `db:"updated_at" json:"updated_at"`
	Children    []*NavigationItem `db:"-" json:"children,omitempty"`
}

// Media represents an uploaded media file
type Media struct {
	ID           uuid.UUID  `db:"id" json:"id"`
	SiteID       uuid.UUID  `db:"site_id" json:"site_id"`
	Name         string     `db:"name" json:"name"`
	OriginalName string     `db:"original_name" json:"original_name"`
	FilePath     string     `db:"file_path" json:"file_path"`
	PublicURL    string     `db:"public_url" json:"public_url"`
	ThumbnailURL *string    `db:"thumbnail_url" json:"thumbnail_url"`
	Type         string     `db:"type" json:"type"`
	MimeType     string     `db:"mime_type" json:"mime_type"`
	FileSize     int64      `db:"file_size" json:"file_size"`
	Width        *int       `db:"width" json:"width"`
	Height       *int       `db:"height" json:"height"`
	Duration     *int       `db:"duration" json:"duration"`
	AltText      *string    `db:"alt_text" json:"alt_text"`
	Caption      *string    `db:"caption" json:"caption"`
	Tags         StringArray `db:"tags" json:"tags"`
	Folder       string     `db:"folder" json:"folder"`
	IsUsed       bool       `db:"is_used" json:"is_used"`
	Metadata     JSONMap    `db:"metadata" json:"metadata,omitempty"`
	UploadedBy   *uuid.UUID `db:"uploaded_by" json:"uploaded_by"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at" json:"-"`
}

// AuditLog represents an audit log entry
type AuditLog struct {
	ID           uuid.UUID  `db:"id" json:"id"`
	UserID       *uuid.UUID `db:"user_id" json:"user_id"`
	UserEmail    *string    `db:"user_email" json:"user_email"`
	UserRole     *string    `db:"user_role" json:"user_role"`
	Action       string     `db:"action" json:"action"`
	ResourceType string     `db:"resource_type" json:"resource_type"`
	ResourceID   *uuid.UUID `db:"resource_id" json:"resource_id"`
	ResourceName *string    `db:"resource_name" json:"resource_name"`
	OldValues    JSONMap    `db:"old_values" json:"old_values,omitempty"`
	NewValues    JSONMap    `db:"new_values" json:"new_values,omitempty"`
	IPAddress    *string    `db:"ip_address" json:"ip_address"`
	UserAgent    *string    `db:"user_agent" json:"user_agent"`
	SiteID       *uuid.UUID `db:"site_id" json:"site_id"`
	Metadata     JSONMap    `db:"metadata" json:"metadata,omitempty"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
}

// Component filter types
type ComponentFilter struct {
	SiteID    *uuid.UUID
	SectionID *uuid.UUID
	IsActive  *bool
	Search    *string
	Pagination
}

// Input types for components
type CreateFeatureInput struct {
	SiteID      uuid.UUID  `json:"site_id" validate:"required"`
	SectionID   *uuid.UUID `json:"section_id"`
	Title       string     `json:"title" validate:"required,min=1,max=255"`
	Description *string    `json:"description"`
	Icon        *string    `json:"icon"`
	IconColor   *string    `json:"icon_color"`
	ImageURL    *string    `json:"image_url"`
	ImageAlt    *string    `json:"image_alt"`
	LinkURL     *string    `json:"link_url"`
	LinkText    *string    `json:"link_text"`
	IsActive    bool       `json:"is_active"`
	SortOrder   int        `json:"sort_order"`
}

type CreateTestimonialInput struct {
	SiteID        uuid.UUID  `json:"site_id" validate:"required"`
	SectionID     *uuid.UUID `json:"section_id"`
	AuthorName    string     `json:"author_name" validate:"required,min=1,max=255"`
	AuthorTitle   *string    `json:"author_title"`
	AuthorCompany *string    `json:"author_company"`
	AuthorAvatar  *string    `json:"author_avatar"`
	Content       string     `json:"content" validate:"required"`
	Rating        *int       `json:"rating" validate:"omitempty,min=1,max=5"`
	Source        *string    `json:"source"`
	SourceURL     *string    `json:"source_url"`
	IsFeatured    bool       `json:"is_featured"`
	IsActive      bool       `json:"is_active"`
	SortOrder     int        `json:"sort_order"`
}

type CreatePricingPlanInput struct {
	SiteID           uuid.UUID  `json:"site_id" validate:"required"`
	SectionID        *uuid.UUID `json:"section_id"`
	Name             string     `json:"name" validate:"required,min=1,max=255"`
	Description      *string    `json:"description"`
	PriceMonthly     *float64   `json:"price_monthly"`
	PriceYearly      *float64   `json:"price_yearly"`
	Currency         string     `json:"currency"`
	PriceLabel       *string    `json:"price_label"`
	IsPopular        bool       `json:"is_popular"`
	IsCustom         bool       `json:"is_custom"`
	BadgeText        *string    `json:"badge_text"`
	CTAText          string     `json:"cta_text" validate:"required"`
	CTALink          *string    `json:"cta_link"`
	Features         []string   `json:"features"`
	FeaturesExcluded []string   `json:"features_excluded"`
	IsActive         bool       `json:"is_active"`
	SortOrder        int        `json:"sort_order"`
}

type CreateFAQInput struct {
	SiteID    uuid.UUID  `json:"site_id" validate:"required"`
	SectionID *uuid.UUID `json:"section_id"`
	Question  string     `json:"question" validate:"required"`
	Answer    string     `json:"answer" validate:"required"`
	Category  *string    `json:"category"`
	IsActive  bool       `json:"is_active"`
	SortOrder int        `json:"sort_order"`
}

type AuditLogFilter struct {
	UserID       *uuid.UUID
	Action       *string
	ResourceType *string
	SiteID       *uuid.UUID
	Pagination
}
