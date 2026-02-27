package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/domain"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/repository"
)

// SiteService defines the interface for site operations
type SiteService interface {
	GetSite(ctx context.Context, id uuid.UUID) (*domain.Site, error)
	GetSiteBySlug(ctx context.Context, slug string) (*domain.Site, error)
	GetSiteByDomain(ctx context.Context, domain string) (*domain.Site, error)
	GetSiteWithSettings(ctx context.Context, id uuid.UUID, publicOnly bool) (*domain.Site, error)
	ListSites(ctx context.Context, filter domain.SiteFilter) (*domain.PaginatedResult[*domain.Site], error)
	CreateSite(ctx context.Context, input domain.CreateSiteInput, userID uuid.UUID) (*domain.Site, error)
	UpdateSite(ctx context.Context, id uuid.UUID, input domain.UpdateSiteInput) (*domain.Site, error)
	DeleteSite(ctx context.Context, id uuid.UUID) error

	// Settings
	GetSettings(ctx context.Context, siteID uuid.UUID, publicOnly bool) ([]*domain.SiteSetting, error)
	GetSettingsMap(ctx context.Context, siteID uuid.UUID, publicOnly bool) (domain.SiteSettingsMap, error)
	UpdateSetting(ctx context.Context, siteID uuid.UUID, key string, input domain.UpdateSettingInput) error
	BulkUpdateSettings(ctx context.Context, siteID uuid.UUID, input domain.BulkUpdateSettingsInput) error
}

// siteService implements SiteService
type siteService struct {
	siteRepo repository.SiteRepository
	logger   zerolog.Logger
}

// NewSiteService creates a new siteService
func NewSiteService(siteRepo repository.SiteRepository, logger zerolog.Logger) SiteService {
	return &siteService{
		siteRepo: siteRepo,
		logger:   logger,
	}
}

// GetSite retrieves a site by ID
func (s *siteService) GetSite(ctx context.Context, id uuid.UUID) (*domain.Site, error) {
	site, err := s.siteRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("siteService.GetSite: %w", err)
	}
	return site, nil
}

// GetSiteBySlug retrieves a site by slug
func (s *siteService) GetSiteBySlug(ctx context.Context, slug string) (*domain.Site, error) {
	site, err := s.siteRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("siteService.GetSiteBySlug: %w", err)
	}
	return site, nil
}

// GetSiteByDomain retrieves a site by domain
func (s *siteService) GetSiteByDomain(ctx context.Context, domainName string) (*domain.Site, error) {
	site, err := s.siteRepo.FindByDomain(ctx, domainName)
	if err != nil {
		return nil, fmt.Errorf("siteService.GetSiteByDomain: %w", err)
	}
	return site, nil
}

// GetSiteWithSettings retrieves a site with its settings
func (s *siteService) GetSiteWithSettings(ctx context.Context, id uuid.UUID, publicOnly bool) (*domain.Site, error) {
	site, err := s.siteRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("siteService.GetSiteWithSettings: %w", err)
	}

	settings, err := s.siteRepo.FindSettingsBySiteID(ctx, id, publicOnly)
	if err != nil {
		return nil, fmt.Errorf("siteService.GetSiteWithSettings settings: %w", err)
	}

	site.Settings = settings
	return site, nil
}

// ListSites retrieves all sites with optional filtering
func (s *siteService) ListSites(ctx context.Context, filter domain.SiteFilter) (*domain.PaginatedResult[*domain.Site], error) {
	sites, total, err := s.siteRepo.FindAll(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("siteService.ListSites: %w", err)
	}

	result := domain.NewPaginatedResult(sites, total, filter.Pagination)
	return &result, nil
}

// CreateSite creates a new site
func (s *siteService) CreateSite(ctx context.Context, input domain.CreateSiteInput, userID uuid.UUID) (*domain.Site, error) {
	site := &domain.Site{
		ID:          uuid.New(),
		Name:        input.Name,
		Slug:        input.Slug,
		Domain:      input.Domain,
		Description: input.Description,
		LogoURL:     input.LogoURL,
		FaviconURL:  input.FaviconURL,
		IsActive:    true,
		CreatedBy:   &userID,
	}

	if err := s.siteRepo.Create(ctx, site); err != nil {
		return nil, fmt.Errorf("siteService.CreateSite: %w", err)
	}

	s.logger.Info().
		Str("site_id", site.ID.String()).
		Str("slug", site.Slug).
		Msg("site created")

	return site, nil
}

// UpdateSite updates an existing site
func (s *siteService) UpdateSite(ctx context.Context, id uuid.UUID, input domain.UpdateSiteInput) (*domain.Site, error) {
	site, err := s.siteRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("siteService.UpdateSite find: %w", err)
	}

	if input.Name != nil {
		site.Name = *input.Name
	}
	if input.Slug != nil {
		site.Slug = *input.Slug
	}
	if input.Domain != nil {
		site.Domain = input.Domain
	}
	if input.Description != nil {
		site.Description = input.Description
	}
	if input.LogoURL != nil {
		site.LogoURL = input.LogoURL
	}
	if input.FaviconURL != nil {
		site.FaviconURL = input.FaviconURL
	}
	if input.IsActive != nil {
		site.IsActive = *input.IsActive
	}

	if err := s.siteRepo.Update(ctx, site); err != nil {
		return nil, fmt.Errorf("siteService.UpdateSite: %w", err)
	}

	return site, nil
}

// DeleteSite soft-deletes a site
func (s *siteService) DeleteSite(ctx context.Context, id uuid.UUID) error {
	if err := s.siteRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("siteService.DeleteSite: %w", err)
	}
	return nil
}

// GetSettings retrieves all settings for a site
func (s *siteService) GetSettings(ctx context.Context, siteID uuid.UUID, publicOnly bool) ([]*domain.SiteSetting, error) {
	settings, err := s.siteRepo.FindSettingsBySiteID(ctx, siteID, publicOnly)
	if err != nil {
		return nil, fmt.Errorf("siteService.GetSettings: %w", err)
	}
	return settings, nil
}

// GetSettingsMap retrieves settings as a key-value map
func (s *siteService) GetSettingsMap(ctx context.Context, siteID uuid.UUID, publicOnly bool) (domain.SiteSettingsMap, error) {
	settings, err := s.siteRepo.FindSettingsBySiteID(ctx, siteID, publicOnly)
	if err != nil {
		return nil, fmt.Errorf("siteService.GetSettingsMap: %w", err)
	}

	settingsMap := make(domain.SiteSettingsMap)
	for _, s := range settings {
		if s.Value != nil {
			settingsMap[s.Key] = *s.Value
		}
	}
	return settingsMap, nil
}

// UpdateSetting updates a single site setting
func (s *siteService) UpdateSetting(ctx context.Context, siteID uuid.UUID, key string, input domain.UpdateSettingInput) error {
	value := ""
	if input.Value != nil {
		value = *input.Value
	}
	if err := s.siteRepo.UpsertSetting(ctx, siteID, key, value); err != nil {
		return fmt.Errorf("siteService.UpdateSetting: %w", err)
	}
	return nil
}

// BulkUpdateSettings updates multiple site settings at once
func (s *siteService) BulkUpdateSettings(ctx context.Context, siteID uuid.UUID, input domain.BulkUpdateSettingsInput) error {
	if err := s.siteRepo.BulkUpsertSettings(ctx, siteID, input.Settings); err != nil {
		return fmt.Errorf("siteService.BulkUpdateSettings: %w", err)
	}
	return nil
}
