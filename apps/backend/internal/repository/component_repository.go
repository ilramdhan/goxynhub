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

// ComponentRepository defines the interface for component data access
type ComponentRepository interface {
	// Features
	FindFeaturesByFilter(ctx context.Context, filter domain.ComponentFilter) ([]*domain.Feature, int, error)
	FindFeatureByID(ctx context.Context, id uuid.UUID) (*domain.Feature, error)
	CreateFeature(ctx context.Context, feature *domain.Feature) error
	UpdateFeature(ctx context.Context, feature *domain.Feature) error
	DeleteFeature(ctx context.Context, id uuid.UUID) error

	// Testimonials
	FindTestimonialsByFilter(ctx context.Context, filter domain.ComponentFilter) ([]*domain.Testimonial, int, error)
	FindTestimonialByID(ctx context.Context, id uuid.UUID) (*domain.Testimonial, error)
	CreateTestimonial(ctx context.Context, testimonial *domain.Testimonial) error
	UpdateTestimonial(ctx context.Context, testimonial *domain.Testimonial) error
	DeleteTestimonial(ctx context.Context, id uuid.UUID) error

	// Pricing Plans
	FindPricingPlansByFilter(ctx context.Context, filter domain.ComponentFilter) ([]*domain.PricingPlan, int, error)
	FindPricingPlanByID(ctx context.Context, id uuid.UUID) (*domain.PricingPlan, error)
	CreatePricingPlan(ctx context.Context, plan *domain.PricingPlan) error
	UpdatePricingPlan(ctx context.Context, plan *domain.PricingPlan) error
	DeletePricingPlan(ctx context.Context, id uuid.UUID) error

	// FAQs
	FindFAQsByFilter(ctx context.Context, filter domain.ComponentFilter) ([]*domain.FAQ, int, error)
	FindFAQByID(ctx context.Context, id uuid.UUID) (*domain.FAQ, error)
	CreateFAQ(ctx context.Context, faq *domain.FAQ) error
	UpdateFAQ(ctx context.Context, faq *domain.FAQ) error
	DeleteFAQ(ctx context.Context, id uuid.UUID) error

	// Navigation
	FindMenusBySiteID(ctx context.Context, siteID uuid.UUID) ([]*domain.NavigationMenu, error)
	FindMenuByIdentifier(ctx context.Context, siteID uuid.UUID, identifier string) (*domain.NavigationMenu, error)
	FindMenuByID(ctx context.Context, id uuid.UUID) (*domain.NavigationMenu, error)
	CreateNavigationMenu(ctx context.Context, menu *domain.NavigationMenu) error
	UpdateNavigationMenu(ctx context.Context, menu *domain.NavigationMenu) error
	DeleteNavigationMenu(ctx context.Context, id uuid.UUID) error
	FindItemsByMenuID(ctx context.Context, menuID uuid.UUID) ([]*domain.NavigationItem, error)
	CreateNavigationItem(ctx context.Context, item *domain.NavigationItem) error
	UpdateNavigationItem(ctx context.Context, item *domain.NavigationItem) error
	DeleteNavigationItem(ctx context.Context, id uuid.UUID) error

	// Media
	FindMediaByFilter(ctx context.Context, siteID uuid.UUID, filter domain.Pagination) ([]*domain.Media, int, error)
	FindMediaByID(ctx context.Context, id uuid.UUID) (*domain.Media, error)
	CreateMedia(ctx context.Context, media *domain.Media) error
	UpdateMedia(ctx context.Context, media *domain.Media) error
	DeleteMedia(ctx context.Context, id uuid.UUID) error

	// Audit Logs
	FindAuditLogs(ctx context.Context, filter domain.AuditLogFilter) ([]*domain.AuditLog, int, error)
	CreateAuditLog(ctx context.Context, log *domain.AuditLog) error
}

// componentRepository implements ComponentRepository
type componentRepository struct {
	db *sqlx.DB
}

// NewComponentRepository creates a new componentRepository
func NewComponentRepository(db *sqlx.DB) ComponentRepository {
	return &componentRepository{db: db}
}

// ─── Features ─────────────────────────────────────────────────────────────────

func (r *componentRepository) FindFeaturesByFilter(ctx context.Context, filter domain.ComponentFilter) ([]*domain.Feature, int, error) {
	args := []interface{}{}
	argIdx := 1
	where := "WHERE deleted_at IS NULL"

	if filter.SiteID != nil {
		where += fmt.Sprintf(" AND site_id = $%d", argIdx)
		args = append(args, *filter.SiteID)
		argIdx++
	}
	if filter.SectionID != nil {
		where += fmt.Sprintf(" AND section_id = $%d", argIdx)
		args = append(args, *filter.SectionID)
		argIdx++
	}
	if filter.IsActive != nil {
		where += fmt.Sprintf(" AND is_active = $%d", argIdx)
		args = append(args, *filter.IsActive)
		argIdx++
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM features %s", where)
	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("componentRepository.FindFeaturesByFilter count: %w", err)
	}

	filter.Normalize()
	dataQuery := fmt.Sprintf(`
		SELECT id, site_id, section_id, title, description, icon, icon_color, image_url, image_alt,
		       link_url, link_text, is_active, sort_order, metadata, created_at, updated_at
		FROM features %s ORDER BY sort_order ASC LIMIT $%d OFFSET $%d
	`, where, argIdx, argIdx+1)
	args = append(args, filter.PerPage, filter.Offset())

	var features []*domain.Feature
	if err := r.db.SelectContext(ctx, &features, dataQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("componentRepository.FindFeaturesByFilter: %w", err)
	}
	return features, total, nil
}

func (r *componentRepository) FindFeatureByID(ctx context.Context, id uuid.UUID) (*domain.Feature, error) {
	query := `SELECT id, site_id, section_id, title, description, icon, icon_color, image_url, image_alt,
		       link_url, link_text, is_active, sort_order, metadata, created_at, updated_at
		FROM features WHERE id = $1 AND deleted_at IS NULL`
	var feature domain.Feature
	if err := r.db.GetContext(ctx, &feature, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("componentRepository.FindFeatureByID: %w", err)
	}
	return &feature, nil
}

func (r *componentRepository) CreateFeature(ctx context.Context, feature *domain.Feature) error {
	query := `INSERT INTO features (id, site_id, section_id, title, description, icon, icon_color, image_url, image_alt, link_url, link_text, is_active, sort_order, metadata)
		VALUES (:id, :site_id, :section_id, :title, :description, :icon, :icon_color, :image_url, :image_alt, :link_url, :link_text, :is_active, :sort_order, :metadata)
		RETURNING created_at, updated_at`
	rows, err := r.db.NamedQueryContext(ctx, query, feature)
	if err != nil {
		return fmt.Errorf("componentRepository.CreateFeature: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&feature.CreatedAt, &feature.UpdatedAt)
	}
	return nil
}

func (r *componentRepository) UpdateFeature(ctx context.Context, feature *domain.Feature) error {
	query := `UPDATE features SET title=:title, description=:description, icon=:icon, icon_color=:icon_color,
		image_url=:image_url, image_alt=:image_alt, link_url=:link_url, link_text=:link_text,
		is_active=:is_active, sort_order=:sort_order, metadata=:metadata, updated_at=NOW()
		WHERE id=:id AND deleted_at IS NULL RETURNING updated_at`
	rows, err := r.db.NamedQueryContext(ctx, query, feature)
	if err != nil {
		return fmt.Errorf("componentRepository.UpdateFeature: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&feature.UpdatedAt)
	}
	return nil
}

func (r *componentRepository) DeleteFeature(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.ExecContext(ctx, `UPDATE features SET deleted_at=NOW() WHERE id=$1 AND deleted_at IS NULL`, id)
	if err != nil {
		return fmt.Errorf("componentRepository.DeleteFeature: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// ─── Testimonials ─────────────────────────────────────────────────────────────

func (r *componentRepository) FindTestimonialsByFilter(ctx context.Context, filter domain.ComponentFilter) ([]*domain.Testimonial, int, error) {
	args := []interface{}{}
	argIdx := 1
	where := "WHERE deleted_at IS NULL"

	if filter.SiteID != nil {
		where += fmt.Sprintf(" AND site_id = $%d", argIdx)
		args = append(args, *filter.SiteID)
		argIdx++
	}
	if filter.IsActive != nil {
		where += fmt.Sprintf(" AND is_active = $%d", argIdx)
		args = append(args, *filter.IsActive)
		argIdx++
	}

	var total int
	r.db.GetContext(ctx, &total, fmt.Sprintf("SELECT COUNT(*) FROM testimonials %s", where), args...)

	filter.Normalize()
	dataQuery := fmt.Sprintf(`SELECT id, site_id, section_id, author_name, author_title, author_company, author_avatar,
		content, rating, source, source_url, is_featured, is_active, sort_order, metadata, created_at, updated_at
		FROM testimonials %s ORDER BY sort_order ASC LIMIT $%d OFFSET $%d`, where, argIdx, argIdx+1)
	args = append(args, filter.PerPage, filter.Offset())

	var testimonials []*domain.Testimonial
	if err := r.db.SelectContext(ctx, &testimonials, dataQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("componentRepository.FindTestimonialsByFilter: %w", err)
	}
	return testimonials, total, nil
}

func (r *componentRepository) FindTestimonialByID(ctx context.Context, id uuid.UUID) (*domain.Testimonial, error) {
	query := `SELECT id, site_id, section_id, author_name, author_title, author_company, author_avatar,
		content, rating, source, source_url, is_featured, is_active, sort_order, metadata, created_at, updated_at
		FROM testimonials WHERE id = $1 AND deleted_at IS NULL`
	var t domain.Testimonial
	if err := r.db.GetContext(ctx, &t, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("componentRepository.FindTestimonialByID: %w", err)
	}
	return &t, nil
}

func (r *componentRepository) CreateTestimonial(ctx context.Context, t *domain.Testimonial) error {
	query := `INSERT INTO testimonials (id, site_id, section_id, author_name, author_title, author_company, author_avatar,
		content, rating, source, source_url, is_featured, is_active, sort_order, metadata)
		VALUES (:id, :site_id, :section_id, :author_name, :author_title, :author_company, :author_avatar,
		:content, :rating, :source, :source_url, :is_featured, :is_active, :sort_order, :metadata)
		RETURNING created_at, updated_at`
	rows, err := r.db.NamedQueryContext(ctx, query, t)
	if err != nil {
		return fmt.Errorf("componentRepository.CreateTestimonial: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&t.CreatedAt, &t.UpdatedAt)
	}
	return nil
}

func (r *componentRepository) UpdateTestimonial(ctx context.Context, t *domain.Testimonial) error {
	query := `UPDATE testimonials SET author_name=:author_name, author_title=:author_title, author_company=:author_company,
		author_avatar=:author_avatar, content=:content, rating=:rating, source=:source, source_url=:source_url,
		is_featured=:is_featured, is_active=:is_active, sort_order=:sort_order, metadata=:metadata, updated_at=NOW()
		WHERE id=:id AND deleted_at IS NULL RETURNING updated_at`
	rows, err := r.db.NamedQueryContext(ctx, query, t)
	if err != nil {
		return fmt.Errorf("componentRepository.UpdateTestimonial: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&t.UpdatedAt)
	}
	return nil
}

func (r *componentRepository) DeleteTestimonial(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.ExecContext(ctx, `UPDATE testimonials SET deleted_at=NOW() WHERE id=$1 AND deleted_at IS NULL`, id)
	if err != nil {
		return fmt.Errorf("componentRepository.DeleteTestimonial: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// ─── Pricing Plans ────────────────────────────────────────────────────────────

func (r *componentRepository) FindPricingPlansByFilter(ctx context.Context, filter domain.ComponentFilter) ([]*domain.PricingPlan, int, error) {
	args := []interface{}{}
	argIdx := 1
	where := "WHERE deleted_at IS NULL"

	if filter.SiteID != nil {
		where += fmt.Sprintf(" AND site_id = $%d", argIdx)
		args = append(args, *filter.SiteID)
		argIdx++
	}

	var total int
	r.db.GetContext(ctx, &total, fmt.Sprintf("SELECT COUNT(*) FROM pricing_plans %s", where), args...)

	filter.Normalize()
	dataQuery := fmt.Sprintf(`SELECT id, site_id, section_id, name, description, price_monthly, price_yearly, currency,
		price_label, is_popular, is_custom, badge_text, cta_text, cta_link, features, features_excluded,
		is_active, sort_order, metadata, created_at, updated_at
		FROM pricing_plans %s ORDER BY sort_order ASC LIMIT $%d OFFSET $%d`, where, argIdx, argIdx+1)
	args = append(args, filter.PerPage, filter.Offset())

	var plans []*domain.PricingPlan
	if err := r.db.SelectContext(ctx, &plans, dataQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("componentRepository.FindPricingPlansByFilter: %w", err)
	}
	return plans, total, nil
}

func (r *componentRepository) FindPricingPlanByID(ctx context.Context, id uuid.UUID) (*domain.PricingPlan, error) {
	query := `SELECT id, site_id, section_id, name, description, price_monthly, price_yearly, currency,
		price_label, is_popular, is_custom, badge_text, cta_text, cta_link, features, features_excluded,
		is_active, sort_order, metadata, created_at, updated_at
		FROM pricing_plans WHERE id = $1 AND deleted_at IS NULL`
	var p domain.PricingPlan
	if err := r.db.GetContext(ctx, &p, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("componentRepository.FindPricingPlanByID: %w", err)
	}
	return &p, nil
}

func (r *componentRepository) CreatePricingPlan(ctx context.Context, p *domain.PricingPlan) error {
	query := `INSERT INTO pricing_plans (id, site_id, section_id, name, description, price_monthly, price_yearly, currency,
		price_label, is_popular, is_custom, badge_text, cta_text, cta_link, features, features_excluded, is_active, sort_order, metadata)
		VALUES (:id, :site_id, :section_id, :name, :description, :price_monthly, :price_yearly, :currency,
		:price_label, :is_popular, :is_custom, :badge_text, :cta_text, :cta_link, :features, :features_excluded, :is_active, :sort_order, :metadata)
		RETURNING created_at, updated_at`
	rows, err := r.db.NamedQueryContext(ctx, query, p)
	if err != nil {
		return fmt.Errorf("componentRepository.CreatePricingPlan: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&p.CreatedAt, &p.UpdatedAt)
	}
	return nil
}

func (r *componentRepository) UpdatePricingPlan(ctx context.Context, p *domain.PricingPlan) error {
	query := `UPDATE pricing_plans SET name=:name, description=:description, price_monthly=:price_monthly,
		price_yearly=:price_yearly, currency=:currency, price_label=:price_label, is_popular=:is_popular,
		is_custom=:is_custom, badge_text=:badge_text, cta_text=:cta_text, cta_link=:cta_link,
		features=:features, features_excluded=:features_excluded, is_active=:is_active, sort_order=:sort_order,
		metadata=:metadata, updated_at=NOW() WHERE id=:id AND deleted_at IS NULL RETURNING updated_at`
	rows, err := r.db.NamedQueryContext(ctx, query, p)
	if err != nil {
		return fmt.Errorf("componentRepository.UpdatePricingPlan: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&p.UpdatedAt)
	}
	return nil
}

func (r *componentRepository) DeletePricingPlan(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.ExecContext(ctx, `UPDATE pricing_plans SET deleted_at=NOW() WHERE id=$1 AND deleted_at IS NULL`, id)
	if err != nil {
		return fmt.Errorf("componentRepository.DeletePricingPlan: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// ─── FAQs ─────────────────────────────────────────────────────────────────────

func (r *componentRepository) FindFAQsByFilter(ctx context.Context, filter domain.ComponentFilter) ([]*domain.FAQ, int, error) {
	args := []interface{}{}
	argIdx := 1
	where := "WHERE deleted_at IS NULL"

	if filter.SiteID != nil {
		where += fmt.Sprintf(" AND site_id = $%d", argIdx)
		args = append(args, *filter.SiteID)
		argIdx++
	}

	var total int
	r.db.GetContext(ctx, &total, fmt.Sprintf("SELECT COUNT(*) FROM faqs %s", where), args...)

	filter.Normalize()
	dataQuery := fmt.Sprintf(`SELECT id, site_id, section_id, question, answer, category, is_active, sort_order, metadata, created_at, updated_at
		FROM faqs %s ORDER BY sort_order ASC LIMIT $%d OFFSET $%d`, where, argIdx, argIdx+1)
	args = append(args, filter.PerPage, filter.Offset())

	var faqs []*domain.FAQ
	if err := r.db.SelectContext(ctx, &faqs, dataQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("componentRepository.FindFAQsByFilter: %w", err)
	}
	return faqs, total, nil
}

func (r *componentRepository) FindFAQByID(ctx context.Context, id uuid.UUID) (*domain.FAQ, error) {
	query := `SELECT id, site_id, section_id, question, answer, category, is_active, sort_order, metadata, created_at, updated_at
		FROM faqs WHERE id = $1 AND deleted_at IS NULL`
	var faq domain.FAQ
	if err := r.db.GetContext(ctx, &faq, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("componentRepository.FindFAQByID: %w", err)
	}
	return &faq, nil
}

func (r *componentRepository) CreateFAQ(ctx context.Context, faq *domain.FAQ) error {
	query := `INSERT INTO faqs (id, site_id, section_id, question, answer, category, is_active, sort_order, metadata)
		VALUES (:id, :site_id, :section_id, :question, :answer, :category, :is_active, :sort_order, :metadata)
		RETURNING created_at, updated_at`
	rows, err := r.db.NamedQueryContext(ctx, query, faq)
	if err != nil {
		return fmt.Errorf("componentRepository.CreateFAQ: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&faq.CreatedAt, &faq.UpdatedAt)
	}
	return nil
}

func (r *componentRepository) UpdateFAQ(ctx context.Context, faq *domain.FAQ) error {
	query := `UPDATE faqs SET question=:question, answer=:answer, category=:category, is_active=:is_active,
		sort_order=:sort_order, metadata=:metadata, updated_at=NOW() WHERE id=:id AND deleted_at IS NULL RETURNING updated_at`
	rows, err := r.db.NamedQueryContext(ctx, query, faq)
	if err != nil {
		return fmt.Errorf("componentRepository.UpdateFAQ: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&faq.UpdatedAt)
	}
	return nil
}

func (r *componentRepository) DeleteFAQ(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.ExecContext(ctx, `UPDATE faqs SET deleted_at=NOW() WHERE id=$1 AND deleted_at IS NULL`, id)
	if err != nil {
		return fmt.Errorf("componentRepository.DeleteFAQ: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// ─── Navigation ───────────────────────────────────────────────────────────────

func (r *componentRepository) FindMenusBySiteID(ctx context.Context, siteID uuid.UUID) ([]*domain.NavigationMenu, error) {
	query := `SELECT id, site_id, name, identifier, description, is_active, metadata, created_at, updated_at
		FROM navigation_menus WHERE site_id = $1 ORDER BY name`
	var menus []*domain.NavigationMenu
	if err := r.db.SelectContext(ctx, &menus, query, siteID); err != nil {
		return nil, fmt.Errorf("componentRepository.FindMenusBySiteID: %w", err)
	}
	return menus, nil
}

func (r *componentRepository) FindMenuByIdentifier(ctx context.Context, siteID uuid.UUID, identifier string) (*domain.NavigationMenu, error) {
	query := `SELECT id, site_id, name, identifier, description, is_active, metadata, created_at, updated_at
		FROM navigation_menus WHERE site_id = $1 AND identifier = $2`
	var menu domain.NavigationMenu
	if err := r.db.GetContext(ctx, &menu, query, siteID, identifier); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("componentRepository.FindMenuByIdentifier: %w", err)
	}
	return &menu, nil
}

func (r *componentRepository) FindMenuByID(ctx context.Context, id uuid.UUID) (*domain.NavigationMenu, error) {
	query := `SELECT id, site_id, name, identifier, description, is_active, metadata, created_at, updated_at
		FROM navigation_menus WHERE id = $1`
	var menu domain.NavigationMenu
	if err := r.db.GetContext(ctx, &menu, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("componentRepository.FindMenuByID: %w", err)
	}
	return &menu, nil
}

func (r *componentRepository) FindItemsByMenuID(ctx context.Context, menuID uuid.UUID) ([]*domain.NavigationItem, error) {
	query := `SELECT id, menu_id, parent_id, page_id, label, url, target, icon, css_class, is_active, is_mega_menu, sort_order, depth, metadata, created_at, updated_at
		FROM navigation_items WHERE menu_id = $1 ORDER BY depth, sort_order`
	var items []*domain.NavigationItem
	if err := r.db.SelectContext(ctx, &items, query, menuID); err != nil {
		return nil, fmt.Errorf("componentRepository.FindItemsByMenuID: %w", err)
	}
	return items, nil
}

func (r *componentRepository) CreateNavigationMenu(ctx context.Context, menu *domain.NavigationMenu) error {
	query := `INSERT INTO navigation_menus (id, site_id, name, identifier, description, is_active, metadata)
		VALUES (:id, :site_id, :name, :identifier, :description, :is_active, :metadata)
		RETURNING created_at, updated_at`
	rows, err := r.db.NamedQueryContext(ctx, query, menu)
	if err != nil {
		return fmt.Errorf("componentRepository.CreateNavigationMenu: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&menu.CreatedAt, &menu.UpdatedAt)
	}
	return nil
}

func (r *componentRepository) UpdateNavigationMenu(ctx context.Context, menu *domain.NavigationMenu) error {
	query := `UPDATE navigation_menus SET name=:name, identifier=:identifier, description=:description,
		is_active=:is_active, metadata=:metadata, updated_at=NOW()
		WHERE id=:id RETURNING updated_at`
	rows, err := r.db.NamedQueryContext(ctx, query, menu)
	if err != nil {
		return fmt.Errorf("componentRepository.UpdateNavigationMenu: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&menu.UpdatedAt)
	}
	return nil
}

func (r *componentRepository) DeleteNavigationMenu(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM navigation_menus WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("componentRepository.DeleteNavigationMenu: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *componentRepository) CreateNavigationItem(ctx context.Context, item *domain.NavigationItem) error {
	query := `INSERT INTO navigation_items (id, menu_id, parent_id, page_id, label, url, target, icon, css_class, is_active, is_mega_menu, sort_order, depth, metadata)
		VALUES (:id, :menu_id, :parent_id, :page_id, :label, :url, :target, :icon, :css_class, :is_active, :is_mega_menu, :sort_order, :depth, :metadata)
		RETURNING created_at, updated_at`
	rows, err := r.db.NamedQueryContext(ctx, query, item)
	if err != nil {
		return fmt.Errorf("componentRepository.CreateNavigationItem: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&item.CreatedAt, &item.UpdatedAt)
	}
	return nil
}

func (r *componentRepository) UpdateNavigationItem(ctx context.Context, item *domain.NavigationItem) error {
	query := `UPDATE navigation_items SET label=:label, url=:url, target=:target, icon=:icon, css_class=:css_class,
		is_active=:is_active, is_mega_menu=:is_mega_menu, sort_order=:sort_order, metadata=:metadata, updated_at=NOW()
		WHERE id=:id RETURNING updated_at`
	rows, err := r.db.NamedQueryContext(ctx, query, item)
	if err != nil {
		return fmt.Errorf("componentRepository.UpdateNavigationItem: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&item.UpdatedAt)
	}
	return nil
}

func (r *componentRepository) DeleteNavigationItem(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM navigation_items WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("componentRepository.DeleteNavigationItem: %w", err)
	}
	return nil
}

// ─── Media ────────────────────────────────────────────────────────────────────

func (r *componentRepository) FindMediaByFilter(ctx context.Context, siteID uuid.UUID, filter domain.Pagination) ([]*domain.Media, int, error) {
	filter.Normalize()
	var total int
	r.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM media WHERE site_id = $1 AND deleted_at IS NULL`, siteID)

	query := `SELECT id, site_id, name, original_name, file_path, public_url, thumbnail_url, type, mime_type,
		file_size, width, height, duration, alt_text, caption, tags, folder, is_used, metadata, uploaded_by, created_at, updated_at
		FROM media WHERE site_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	var media []*domain.Media
	if err := r.db.SelectContext(ctx, &media, query, siteID, filter.PerPage, filter.Offset()); err != nil {
		return nil, 0, fmt.Errorf("componentRepository.FindMediaByFilter: %w", err)
	}
	return media, total, nil
}

func (r *componentRepository) FindMediaByID(ctx context.Context, id uuid.UUID) (*domain.Media, error) {
	query := `SELECT id, site_id, name, original_name, file_path, public_url, thumbnail_url, type, mime_type,
		file_size, width, height, duration, alt_text, caption, tags, folder, is_used, metadata, uploaded_by, created_at, updated_at
		FROM media WHERE id = $1 AND deleted_at IS NULL`
	var m domain.Media
	if err := r.db.GetContext(ctx, &m, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("componentRepository.FindMediaByID: %w", err)
	}
	return &m, nil
}

func (r *componentRepository) CreateMedia(ctx context.Context, m *domain.Media) error {
	query := `INSERT INTO media (id, site_id, name, original_name, file_path, public_url, thumbnail_url, type, mime_type,
		file_size, width, height, alt_text, caption, tags, folder, uploaded_by, metadata)
		VALUES (:id, :site_id, :name, :original_name, :file_path, :public_url, :thumbnail_url, :type, :mime_type,
		:file_size, :width, :height, :alt_text, :caption, :tags, :folder, :uploaded_by, :metadata)
		RETURNING created_at, updated_at`
	rows, err := r.db.NamedQueryContext(ctx, query, m)
	if err != nil {
		return fmt.Errorf("componentRepository.CreateMedia: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&m.CreatedAt, &m.UpdatedAt)
	}
	return nil
}

func (r *componentRepository) UpdateMedia(ctx context.Context, m *domain.Media) error {
	query := `UPDATE media SET name=:name, alt_text=:alt_text, caption=:caption, tags=:tags, folder=:folder,
		metadata=:metadata, updated_at=NOW() WHERE id=:id AND deleted_at IS NULL RETURNING updated_at`
	rows, err := r.db.NamedQueryContext(ctx, query, m)
	if err != nil {
		return fmt.Errorf("componentRepository.UpdateMedia: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&m.UpdatedAt)
	}
	return nil
}

func (r *componentRepository) DeleteMedia(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.ExecContext(ctx, `UPDATE media SET deleted_at=NOW() WHERE id=$1 AND deleted_at IS NULL`, id)
	if err != nil {
		return fmt.Errorf("componentRepository.DeleteMedia: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// ─── Audit Logs ───────────────────────────────────────────────────────────────

func (r *componentRepository) FindAuditLogs(ctx context.Context, filter domain.AuditLogFilter) ([]*domain.AuditLog, int, error) {
	args := []interface{}{}
	argIdx := 1
	where := "WHERE 1=1"

	if filter.SiteID != nil {
		where += fmt.Sprintf(" AND site_id = $%d", argIdx)
		args = append(args, *filter.SiteID)
		argIdx++
	}
	if filter.UserID != nil {
		where += fmt.Sprintf(" AND user_id = $%d", argIdx)
		args = append(args, *filter.UserID)
		argIdx++
	}
	if filter.Action != nil {
		where += fmt.Sprintf(" AND action = $%d", argIdx)
		args = append(args, *filter.Action)
		argIdx++
	}
	if filter.ResourceType != nil {
		where += fmt.Sprintf(" AND resource_type = $%d", argIdx)
		args = append(args, *filter.ResourceType)
		argIdx++
	}

	var total int
	r.db.GetContext(ctx, &total, fmt.Sprintf("SELECT COUNT(*) FROM audit_logs %s", where), args...)

	filter.Normalize()
	dataQuery := fmt.Sprintf(`SELECT id, user_id, user_email, user_role, action, resource_type, resource_id, resource_name,
		old_values, new_values, ip_address, site_id, metadata, created_at
		FROM audit_logs %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`, where, argIdx, argIdx+1)
	args = append(args, filter.PerPage, filter.Offset())

	var logs []*domain.AuditLog
	if err := r.db.SelectContext(ctx, &logs, dataQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("componentRepository.FindAuditLogs: %w", err)
	}
	return logs, total, nil
}

func (r *componentRepository) CreateAuditLog(ctx context.Context, log *domain.AuditLog) error {
	query := `INSERT INTO audit_logs (id, user_id, user_email, user_role, action, resource_type, resource_id, resource_name,
		old_values, new_values, ip_address, site_id, metadata)
		VALUES (:id, :user_id, :user_email, :user_role, :action, :resource_type, :resource_id, :resource_name,
		:old_values, :new_values, :ip_address, :site_id, :metadata)`
	_, err := r.db.NamedExecContext(ctx, query, log)
	if err != nil {
		return fmt.Errorf("componentRepository.CreateAuditLog: %w", err)
	}
	return nil
}
