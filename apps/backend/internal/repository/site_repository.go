package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/domain"
)

// SiteRepository defines the interface for site data access
type SiteRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Site, error)
	FindBySlug(ctx context.Context, slug string) (*domain.Site, error)
	FindByDomain(ctx context.Context, domain string) (*domain.Site, error)
	FindAll(ctx context.Context, filter domain.SiteFilter) ([]*domain.Site, int, error)
	Create(ctx context.Context, site *domain.Site) error
	Update(ctx context.Context, site *domain.Site) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Settings
	FindSettingsBySiteID(ctx context.Context, siteID uuid.UUID, publicOnly bool) ([]*domain.SiteSetting, error)
	FindSettingByKey(ctx context.Context, siteID uuid.UUID, key string) (*domain.SiteSetting, error)
	UpsertSetting(ctx context.Context, siteID uuid.UUID, key, value string) error
	BulkUpsertSettings(ctx context.Context, siteID uuid.UUID, settings map[string]string) error
}

// siteRepository implements SiteRepository
type siteRepository struct {
	db *sqlx.DB
}

// NewSiteRepository creates a new siteRepository
func NewSiteRepository(db *sqlx.DB) SiteRepository {
	return &siteRepository{db: db}
}

// FindByID retrieves a site by ID
func (r *siteRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Site, error) {
	query := `
		SELECT id, name, slug, domain, description, logo_url, favicon_url, is_active, metadata, created_by, created_at, updated_at
		FROM sites
		WHERE id = $1 AND deleted_at IS NULL
	`
	var site domain.Site
	if err := r.db.GetContext(ctx, &site, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("siteRepository.FindByID: %w", err)
	}
	return &site, nil
}

// FindBySlug retrieves a site by slug
func (r *siteRepository) FindBySlug(ctx context.Context, slug string) (*domain.Site, error) {
	query := `
		SELECT id, name, slug, domain, description, logo_url, favicon_url, is_active, metadata, created_by, created_at, updated_at
		FROM sites
		WHERE slug = $1 AND deleted_at IS NULL
	`
	var site domain.Site
	if err := r.db.GetContext(ctx, &site, query, slug); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("siteRepository.FindBySlug: %w", err)
	}
	return &site, nil
}

// FindByDomain retrieves a site by domain
func (r *siteRepository) FindByDomain(ctx context.Context, domainName string) (*domain.Site, error) {
	query := `
		SELECT id, name, slug, domain, description, logo_url, favicon_url, is_active, metadata, created_by, created_at, updated_at
		FROM sites
		WHERE domain = $1 AND deleted_at IS NULL AND is_active = true
	`
	var site domain.Site
	if err := r.db.GetContext(ctx, &site, query, domainName); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("siteRepository.FindByDomain: %w", err)
	}
	return &site, nil
}

// FindAll retrieves all sites with optional filtering
func (r *siteRepository) FindAll(ctx context.Context, filter domain.SiteFilter) ([]*domain.Site, int, error) {
	args := []interface{}{}
	argIdx := 1
	where := "WHERE deleted_at IS NULL"

	if filter.IsActive != nil {
		where += fmt.Sprintf(" AND is_active = $%d", argIdx)
		args = append(args, *filter.IsActive)
		argIdx++
	}
	if filter.Search != nil && *filter.Search != "" {
		where += fmt.Sprintf(" AND (name ILIKE $%d OR slug ILIKE $%d)", argIdx, argIdx+1)
		searchTerm := "%" + *filter.Search + "%"
		args = append(args, searchTerm, searchTerm)
		argIdx += 2
	}

	// Count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM sites %s", where)
	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("siteRepository.FindAll count: %w", err)
	}

	// Data
	filter.Normalize()
	dataQuery := fmt.Sprintf(`
		SELECT id, name, slug, domain, description, logo_url, favicon_url, is_active, metadata, created_by, created_at, updated_at
		FROM sites %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, where, argIdx, argIdx+1)
	args = append(args, filter.PerPage, filter.Offset())

	var sites []*domain.Site
	if err := r.db.SelectContext(ctx, &sites, dataQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("siteRepository.FindAll: %w", err)
	}

	return sites, total, nil
}

// Create inserts a new site
func (r *siteRepository) Create(ctx context.Context, site *domain.Site) error {
	query := `
		INSERT INTO sites (id, name, slug, domain, description, logo_url, favicon_url, is_active, metadata, created_by)
		VALUES (:id, :name, :slug, :domain, :description, :logo_url, :favicon_url, :is_active, :metadata, :created_by)
		RETURNING created_at, updated_at
	`
	rows, err := r.db.NamedQueryContext(ctx, query, site)
	if err != nil {
		return fmt.Errorf("siteRepository.Create: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&site.CreatedAt, &site.UpdatedAt); err != nil {
			return fmt.Errorf("siteRepository.Create scan: %w", err)
		}
	}
	return nil
}

// Update updates an existing site
func (r *siteRepository) Update(ctx context.Context, site *domain.Site) error {
	query := `
		UPDATE sites SET
			name = :name, slug = :slug, domain = :domain, description = :description,
			logo_url = :logo_url, favicon_url = :favicon_url, is_active = :is_active,
			metadata = :metadata, updated_at = NOW()
		WHERE id = :id AND deleted_at IS NULL
		RETURNING updated_at
	`
	rows, err := r.db.NamedQueryContext(ctx, query, site)
	if err != nil {
		return fmt.Errorf("siteRepository.Update: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&site.UpdatedAt); err != nil {
			return fmt.Errorf("siteRepository.Update scan: %w", err)
		}
	}
	return nil
}

// Delete soft-deletes a site
func (r *siteRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE sites SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("siteRepository.Delete: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// FindSettingsBySiteID retrieves all settings for a site
func (r *siteRepository) FindSettingsBySiteID(ctx context.Context, siteID uuid.UUID, publicOnly bool) ([]*domain.SiteSetting, error) {
	query := `
		SELECT id, site_id, key, value, value_json, type, group_name, label, description, is_public, sort_order, created_at, updated_at
		FROM site_settings
		WHERE site_id = $1
	`
	args := []interface{}{siteID}
	if publicOnly {
		query += " AND is_public = true"
	}
	query += " ORDER BY group_name, sort_order"

	var settings []*domain.SiteSetting
	if err := r.db.SelectContext(ctx, &settings, query, args...); err != nil {
		return nil, fmt.Errorf("siteRepository.FindSettingsBySiteID: %w", err)
	}
	return settings, nil
}

// FindSettingByKey retrieves a specific setting by key
func (r *siteRepository) FindSettingByKey(ctx context.Context, siteID uuid.UUID, key string) (*domain.SiteSetting, error) {
	query := `
		SELECT id, site_id, key, value, value_json, type, group_name, label, description, is_public, sort_order, created_at, updated_at
		FROM site_settings
		WHERE site_id = $1 AND key = $2
	`
	var setting domain.SiteSetting
	if err := r.db.GetContext(ctx, &setting, query, siteID, key); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("siteRepository.FindSettingByKey: %w", err)
	}
	return &setting, nil
}

// UpsertSetting creates or updates a site setting
func (r *siteRepository) UpsertSetting(ctx context.Context, siteID uuid.UUID, key, value string) error {
	query := `
		INSERT INTO site_settings (id, site_id, key, value)
		VALUES (gen_random_uuid(), $1, $2, $3)
		ON CONFLICT (site_id, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW()
	`
	_, err := r.db.ExecContext(ctx, query, siteID, key, value)
	if err != nil {
		return fmt.Errorf("siteRepository.UpsertSetting: %w", err)
	}
	return nil
}

// BulkUpsertSettings creates or updates multiple site settings
func (r *siteRepository) BulkUpsertSettings(ctx context.Context, siteID uuid.UUID, settings map[string]string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("siteRepository.BulkUpsertSettings begin tx: %w", err)
	}
	defer tx.Rollback()

	for key, value := range settings {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO site_settings (id, site_id, key, value)
			VALUES (gen_random_uuid(), $1, $2, $3)
			ON CONFLICT (site_id, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW()
		`, siteID, key, value)
		if err != nil {
			return fmt.Errorf("siteRepository.BulkUpsertSettings upsert: %w", err)
		}
	}

	return tx.Commit()
}
