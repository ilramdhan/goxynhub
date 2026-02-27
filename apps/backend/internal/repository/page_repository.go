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

// PageRepository defines the interface for page data access
type PageRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Page, error)
	FindBySlug(ctx context.Context, siteID uuid.UUID, slug string) (*domain.Page, error)
	FindHomepage(ctx context.Context, siteID uuid.UUID) (*domain.Page, error)
	FindAll(ctx context.Context, filter domain.PageFilter) ([]*domain.Page, int, error)
	Create(ctx context.Context, page *domain.Page) error
	Update(ctx context.Context, page *domain.Page) error
	Delete(ctx context.Context, id uuid.UUID) error
	SetHomepage(ctx context.Context, siteID, pageID uuid.UUID) error

	// Section operations
	FindSectionsByPageID(ctx context.Context, pageID uuid.UUID) ([]*domain.PageSection, error)
	FindSectionByID(ctx context.Context, id uuid.UUID) (*domain.PageSection, error)
	CreateSection(ctx context.Context, section *domain.PageSection) error
	UpdateSection(ctx context.Context, section *domain.PageSection) error
	DeleteSection(ctx context.Context, id uuid.UUID) error
	ReorderSections(ctx context.Context, orders []domain.SectionOrder) error

	// Content operations
	FindContentsBySectionID(ctx context.Context, sectionID uuid.UUID) ([]*domain.SectionContent, error)
	FindContentByKey(ctx context.Context, sectionID uuid.UUID, key string) (*domain.SectionContent, error)
	UpsertContent(ctx context.Context, content *domain.SectionContent) error
	DeleteContent(ctx context.Context, id uuid.UUID) error
	DeleteContentByKey(ctx context.Context, sectionID uuid.UUID, key string) error
}

// pageRepository implements PageRepository
type pageRepository struct {
	db *sqlx.DB
}

// NewPageRepository creates a new pageRepository
func NewPageRepository(db *sqlx.DB) PageRepository {
	return &pageRepository{db: db}
}

// FindByID retrieves a page by ID
func (r *pageRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Page, error) {
	query := `
		SELECT id, site_id, title, slug, description, status, is_homepage,
		       seo_title, seo_description, seo_keywords,
		       og_title, og_description, og_image, og_type,
		       twitter_title, twitter_description, twitter_image, twitter_card,
		       schema_markup, custom_head, canonical_url, robots_meta, template,
		       sort_order, published_at, metadata, created_by, updated_by, created_at, updated_at
		FROM pages
		WHERE id = $1 AND deleted_at IS NULL
	`
	var page domain.Page
	if err := r.db.GetContext(ctx, &page, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("pageRepository.FindByID: %w", err)
	}
	return &page, nil
}

// FindBySlug retrieves a page by slug and site ID
func (r *pageRepository) FindBySlug(ctx context.Context, siteID uuid.UUID, slug string) (*domain.Page, error) {
	query := `
		SELECT id, site_id, title, slug, description, status, is_homepage,
		       seo_title, seo_description, seo_keywords,
		       og_title, og_description, og_image, og_type,
		       twitter_title, twitter_description, twitter_image, twitter_card,
		       schema_markup, custom_head, canonical_url, robots_meta, template,
		       sort_order, published_at, metadata, created_by, updated_by, created_at, updated_at
		FROM pages
		WHERE site_id = $1 AND slug = $2 AND deleted_at IS NULL
	`
	var page domain.Page
	if err := r.db.GetContext(ctx, &page, query, siteID, slug); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("pageRepository.FindBySlug: %w", err)
	}
	return &page, nil
}

// FindHomepage retrieves the homepage for a site
func (r *pageRepository) FindHomepage(ctx context.Context, siteID uuid.UUID) (*domain.Page, error) {
	query := `
		SELECT id, site_id, title, slug, description, status, is_homepage,
		       seo_title, seo_description, seo_keywords,
		       og_title, og_description, og_image, og_type,
		       twitter_title, twitter_description, twitter_image, twitter_card,
		       schema_markup, custom_head, canonical_url, robots_meta, template,
		       sort_order, published_at, metadata, created_by, updated_by, created_at, updated_at
		FROM pages
		WHERE site_id = $1 AND is_homepage = true AND status = 'published' AND deleted_at IS NULL
		LIMIT 1
	`
	var page domain.Page
	if err := r.db.GetContext(ctx, &page, query, siteID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("pageRepository.FindHomepage: %w", err)
	}
	return &page, nil
}

// FindAll retrieves all pages with optional filtering
func (r *pageRepository) FindAll(ctx context.Context, filter domain.PageFilter) ([]*domain.Page, int, error) {
	args := []interface{}{}
	argIdx := 1
	where := "WHERE deleted_at IS NULL"

	if filter.SiteID != nil {
		where += fmt.Sprintf(" AND site_id = $%d", argIdx)
		args = append(args, *filter.SiteID)
		argIdx++
	}
	if filter.Status != nil {
		where += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, *filter.Status)
		argIdx++
	}
	if filter.IsHomepage != nil {
		where += fmt.Sprintf(" AND is_homepage = $%d", argIdx)
		args = append(args, *filter.IsHomepage)
		argIdx++
	}
	if filter.Search != nil && *filter.Search != "" {
		where += fmt.Sprintf(" AND (title ILIKE $%d OR slug ILIKE $%d)", argIdx, argIdx+1)
		searchTerm := "%" + *filter.Search + "%"
		args = append(args, searchTerm, searchTerm)
		argIdx += 2
	}

	// Count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM pages %s", where)
	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("pageRepository.FindAll count: %w", err)
	}

	// Data
	filter.Normalize()
	dataQuery := fmt.Sprintf(`
		SELECT id, site_id, title, slug, description, status, is_homepage,
		       seo_title, seo_description, seo_keywords,
		       og_title, og_description, og_image, og_type,
		       twitter_title, twitter_description, twitter_image, twitter_card,
		       schema_markup, custom_head, canonical_url, robots_meta, template,
		       sort_order, published_at, metadata, created_by, updated_by, created_at, updated_at
		FROM pages %s
		ORDER BY sort_order ASC, created_at DESC
		LIMIT $%d OFFSET $%d
	`, where, argIdx, argIdx+1)
	args = append(args, filter.PerPage, filter.Offset())

	var pages []*domain.Page
	if err := r.db.SelectContext(ctx, &pages, dataQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("pageRepository.FindAll: %w", err)
	}

	return pages, total, nil
}

// Create inserts a new page
func (r *pageRepository) Create(ctx context.Context, page *domain.Page) error {
	query := `
		INSERT INTO pages (
			id, site_id, title, slug, description, status, is_homepage,
			seo_title, seo_description, seo_keywords,
			og_title, og_description, og_image, og_type,
			twitter_title, twitter_description, twitter_image, twitter_card,
			canonical_url, robots_meta, template, sort_order, metadata, created_by, updated_by
		) VALUES (
			:id, :site_id, :title, :slug, :description, :status, :is_homepage,
			:seo_title, :seo_description, :seo_keywords,
			:og_title, :og_description, :og_image, :og_type,
			:twitter_title, :twitter_description, :twitter_image, :twitter_card,
			:canonical_url, :robots_meta, :template, :sort_order, :metadata, :created_by, :updated_by
		)
		RETURNING created_at, updated_at
	`
	rows, err := r.db.NamedQueryContext(ctx, query, page)
	if err != nil {
		return fmt.Errorf("pageRepository.Create: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&page.CreatedAt, &page.UpdatedAt); err != nil {
			return fmt.Errorf("pageRepository.Create scan: %w", err)
		}
	}
	return nil
}

// Update updates an existing page
func (r *pageRepository) Update(ctx context.Context, page *domain.Page) error {
	query := `
		UPDATE pages SET
			title = :title, slug = :slug, description = :description,
			status = :status, is_homepage = :is_homepage,
			seo_title = :seo_title, seo_description = :seo_description, seo_keywords = :seo_keywords,
			og_title = :og_title, og_description = :og_description, og_image = :og_image, og_type = :og_type,
			twitter_title = :twitter_title, twitter_description = :twitter_description,
			twitter_image = :twitter_image, twitter_card = :twitter_card,
			custom_head = :custom_head, canonical_url = :canonical_url,
			robots_meta = :robots_meta, template = :template,
			metadata = :metadata, updated_by = :updated_by, updated_at = NOW()
		WHERE id = :id AND deleted_at IS NULL
		RETURNING updated_at
	`
	rows, err := r.db.NamedQueryContext(ctx, query, page)
	if err != nil {
		return fmt.Errorf("pageRepository.Update: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&page.UpdatedAt); err != nil {
			return fmt.Errorf("pageRepository.Update scan: %w", err)
		}
	}
	return nil
}

// Delete soft-deletes a page
func (r *pageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE pages SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("pageRepository.Delete: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// SetHomepage sets a page as the homepage for a site (unsets others)
func (r *pageRepository) SetHomepage(ctx context.Context, siteID, pageID uuid.UUID) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("pageRepository.SetHomepage begin tx: %w", err)
	}
	defer tx.Rollback()

	// Unset current homepage
	_, err = tx.ExecContext(ctx, `UPDATE pages SET is_homepage = false WHERE site_id = $1 AND is_homepage = true`, siteID)
	if err != nil {
		return fmt.Errorf("pageRepository.SetHomepage unset: %w", err)
	}

	// Set new homepage
	_, err = tx.ExecContext(ctx, `UPDATE pages SET is_homepage = true WHERE id = $1`, pageID)
	if err != nil {
		return fmt.Errorf("pageRepository.SetHomepage set: %w", err)
	}

	return tx.Commit()
}

// FindSectionsByPageID retrieves all sections for a page
func (r *pageRepository) FindSectionsByPageID(ctx context.Context, pageID uuid.UUID) ([]*domain.PageSection, error) {
	query := `
		SELECT id, page_id, name, type, identifier, is_visible, sort_order,
		       bg_color, bg_image, bg_video, bg_overlay, bg_overlay_color, bg_overlay_opacity,
		       layout, padding_top, padding_bottom, animation, css_class, custom_css,
		       metadata, created_at, updated_at
		FROM page_sections
		WHERE page_id = $1
		ORDER BY sort_order ASC
	`
	var sections []*domain.PageSection
	if err := r.db.SelectContext(ctx, &sections, query, pageID); err != nil {
		return nil, fmt.Errorf("pageRepository.FindSectionsByPageID: %w", err)
	}
	return sections, nil
}

// FindSectionByID retrieves a section by ID
func (r *pageRepository) FindSectionByID(ctx context.Context, id uuid.UUID) (*domain.PageSection, error) {
	query := `
		SELECT id, page_id, name, type, identifier, is_visible, sort_order,
		       bg_color, bg_image, bg_video, bg_overlay, bg_overlay_color, bg_overlay_opacity,
		       layout, padding_top, padding_bottom, animation, css_class, custom_css,
		       metadata, created_at, updated_at
		FROM page_sections
		WHERE id = $1
	`
	var section domain.PageSection
	if err := r.db.GetContext(ctx, &section, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("pageRepository.FindSectionByID: %w", err)
	}
	return &section, nil
}

// CreateSection inserts a new page section
func (r *pageRepository) CreateSection(ctx context.Context, section *domain.PageSection) error {
	query := `
		INSERT INTO page_sections (id, page_id, name, type, identifier, is_visible, sort_order, metadata)
		VALUES (:id, :page_id, :name, :type, :identifier, :is_visible, :sort_order, :metadata)
		RETURNING created_at, updated_at
	`
	rows, err := r.db.NamedQueryContext(ctx, query, section)
	if err != nil {
		return fmt.Errorf("pageRepository.CreateSection: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&section.CreatedAt, &section.UpdatedAt); err != nil {
			return fmt.Errorf("pageRepository.CreateSection scan: %w", err)
		}
	}
	return nil
}

// UpdateSection updates an existing page section
func (r *pageRepository) UpdateSection(ctx context.Context, section *domain.PageSection) error {
	query := `
		UPDATE page_sections SET
			name = :name, type = :type, is_visible = :is_visible, sort_order = :sort_order,
			bg_color = :bg_color, bg_image = :bg_image, bg_video = :bg_video,
			bg_overlay = :bg_overlay, bg_overlay_color = :bg_overlay_color, bg_overlay_opacity = :bg_overlay_opacity,
			layout = :layout, padding_top = :padding_top, padding_bottom = :padding_bottom,
			animation = :animation, css_class = :css_class, custom_css = :custom_css,
			metadata = :metadata, updated_at = NOW()
		WHERE id = :id
		RETURNING updated_at
	`
	rows, err := r.db.NamedQueryContext(ctx, query, section)
	if err != nil {
		return fmt.Errorf("pageRepository.UpdateSection: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&section.UpdatedAt); err != nil {
			return fmt.Errorf("pageRepository.UpdateSection scan: %w", err)
		}
	}
	return nil
}

// DeleteSection deletes a page section
func (r *pageRepository) DeleteSection(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM page_sections WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("pageRepository.DeleteSection: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// ReorderSections updates the sort order of multiple sections
func (r *pageRepository) ReorderSections(ctx context.Context, orders []domain.SectionOrder) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("pageRepository.ReorderSections begin tx: %w", err)
	}
	defer tx.Rollback()

	for _, order := range orders {
		_, err := tx.ExecContext(ctx, `UPDATE page_sections SET sort_order = $1 WHERE id = $2`, order.SortOrder, order.ID)
		if err != nil {
			return fmt.Errorf("pageRepository.ReorderSections update: %w", err)
		}
	}

	return tx.Commit()
}

// FindContentsBySectionID retrieves all content items for a section
func (r *pageRepository) FindContentsBySectionID(ctx context.Context, sectionID uuid.UUID) ([]*domain.SectionContent, error) {
	query := `
		SELECT id, section_id, key, value, value_json, type, label, description, placeholder,
		       is_required, sort_order, alt_text, width, height, link_url, link_target,
		       metadata, created_at, updated_at
		FROM section_contents
		WHERE section_id = $1
		ORDER BY sort_order ASC
	`
	var contents []*domain.SectionContent
	if err := r.db.SelectContext(ctx, &contents, query, sectionID); err != nil {
		return nil, fmt.Errorf("pageRepository.FindContentsBySectionID: %w", err)
	}
	return contents, nil
}

// FindContentByKey retrieves a content item by section ID and key
func (r *pageRepository) FindContentByKey(ctx context.Context, sectionID uuid.UUID, key string) (*domain.SectionContent, error) {
	query := `
		SELECT id, section_id, key, value, value_json, type, label, description, placeholder,
		       is_required, sort_order, alt_text, width, height, link_url, link_target,
		       metadata, created_at, updated_at
		FROM section_contents
		WHERE section_id = $1 AND key = $2
	`
	var content domain.SectionContent
	if err := r.db.GetContext(ctx, &content, query, sectionID, key); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("pageRepository.FindContentByKey: %w", err)
	}
	return &content, nil
}

// UpsertContent creates or updates a content item
func (r *pageRepository) UpsertContent(ctx context.Context, content *domain.SectionContent) error {
	query := `
		INSERT INTO section_contents (
			id, section_id, key, value, value_json, type, label, description, placeholder,
			is_required, sort_order, alt_text, width, height, link_url, link_target, metadata
		) VALUES (
			:id, :section_id, :key, :value, :value_json, :type, :label, :description, :placeholder,
			:is_required, :sort_order, :alt_text, :width, :height, :link_url, :link_target, :metadata
		)
		ON CONFLICT (section_id, key) DO UPDATE SET
			value = EXCLUDED.value,
			value_json = EXCLUDED.value_json,
			type = EXCLUDED.type,
			label = EXCLUDED.label,
			description = EXCLUDED.description,
			placeholder = EXCLUDED.placeholder,
			is_required = EXCLUDED.is_required,
			sort_order = EXCLUDED.sort_order,
			alt_text = EXCLUDED.alt_text,
			width = EXCLUDED.width,
			height = EXCLUDED.height,
			link_url = EXCLUDED.link_url,
			link_target = EXCLUDED.link_target,
			metadata = EXCLUDED.metadata,
			updated_at = NOW()
		RETURNING id, created_at, updated_at
	`
	rows, err := r.db.NamedQueryContext(ctx, query, content)
	if err != nil {
		return fmt.Errorf("pageRepository.UpsertContent: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&content.ID, &content.CreatedAt, &content.UpdatedAt); err != nil {
			return fmt.Errorf("pageRepository.UpsertContent scan: %w", err)
		}
	}
	return nil
}

// DeleteContent deletes a content item by ID
func (r *pageRepository) DeleteContent(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM section_contents WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("pageRepository.DeleteContent: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// DeleteContentByKey deletes a content item by section ID and key
func (r *pageRepository) DeleteContentByKey(ctx context.Context, sectionID uuid.UUID, key string) error {
	query := `DELETE FROM section_contents WHERE section_id = $1 AND key = $2`
	_, err := r.db.ExecContext(ctx, query, sectionID, key)
	if err != nil {
		return fmt.Errorf("pageRepository.DeleteContentByKey: %w", err)
	}
	return nil
}
