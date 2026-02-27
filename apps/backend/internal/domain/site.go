package domain

import (
	"time"

	"github.com/google/uuid"
)

// Site represents a website managed by the CMS
type Site struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	Name        string     `db:"name" json:"name"`
	Slug        string     `db:"slug" json:"slug"`
	Domain      *string    `db:"domain" json:"domain"`
	Description *string    `db:"description" json:"description"`
	LogoURL     *string    `db:"logo_url" json:"logo_url"`
	FaviconURL  *string    `db:"favicon_url" json:"favicon_url"`
	IsActive    bool       `db:"is_active" json:"is_active"`
	Metadata    JSONMap    `db:"metadata" json:"metadata,omitempty"`
	CreatedBy   *uuid.UUID `db:"created_by" json:"created_by"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at" json:"-"`
	// Relations
	Settings []*SiteSetting `db:"-" json:"settings,omitempty"`
}

// SiteSetting represents a key-value setting for a site
type SiteSetting struct {
	ID          uuid.UUID `db:"id" json:"id"`
	SiteID      uuid.UUID `db:"site_id" json:"site_id"`
	Key         string    `db:"key" json:"key"`
	Value       *string   `db:"value" json:"value"`
	ValueJSON   JSONMap   `db:"value_json" json:"value_json,omitempty"`
	Type        string    `db:"type" json:"type"`
	GroupName   string    `db:"group_name" json:"group_name"`
	Label       *string   `db:"label" json:"label"`
	Description *string   `db:"description" json:"description"`
	IsPublic    bool      `db:"is_public" json:"is_public"`
	SortOrder   int       `db:"sort_order" json:"sort_order"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// SiteSettingsMap is a convenience type for settings as a map
type SiteSettingsMap map[string]string

// SiteFilter holds filter parameters for site queries
type SiteFilter struct {
	IsActive *bool
	Search   *string
	Pagination
}

// CreateSiteInput holds data for creating a new site
type CreateSiteInput struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Slug        string  `json:"slug" validate:"required,min=1,max=255"`
	Domain      *string `json:"domain" validate:"omitempty,hostname"`
	Description *string `json:"description" validate:"omitempty,max=500"`
	LogoURL     *string `json:"logo_url" validate:"omitempty,url"`
	FaviconURL  *string `json:"favicon_url" validate:"omitempty,url"`
}

// UpdateSiteInput holds data for updating a site
type UpdateSiteInput struct {
	Name        *string `json:"name" validate:"omitempty,min=1,max=255"`
	Slug        *string `json:"slug" validate:"omitempty,min=1,max=255"`
	Domain      *string `json:"domain" validate:"omitempty,hostname"`
	Description *string `json:"description" validate:"omitempty,max=500"`
	LogoURL     *string `json:"logo_url" validate:"omitempty,url"`
	FaviconURL  *string `json:"favicon_url" validate:"omitempty,url"`
	IsActive    *bool   `json:"is_active"`
}

// UpdateSettingInput holds data for updating a site setting
type UpdateSettingInput struct {
	Value     *string `json:"value"`
	ValueJSON JSONMap `json:"value_json"`
}

// BulkUpdateSettingsInput holds data for bulk updating site settings
type BulkUpdateSettingsInput struct {
	Settings map[string]string `json:"settings" validate:"required"`
}
