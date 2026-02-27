package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/domain"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/service"
)

// ─── Mock SiteRepository ──────────────────────────────────────────────────────

type mockSiteRepository struct {
	sites    map[uuid.UUID]*domain.Site
	settings map[string]*domain.SiteSetting // key: siteID+":"+key
}

func newMockSiteRepository() *mockSiteRepository {
	return &mockSiteRepository{
		sites:    make(map[uuid.UUID]*domain.Site),
		settings: make(map[string]*domain.SiteSetting),
	}
}

func (m *mockSiteRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Site, error) {
	if s, ok := m.sites[id]; ok {
		return s, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockSiteRepository) FindBySlug(ctx context.Context, slug string) (*domain.Site, error) {
	for _, s := range m.sites {
		if s.Slug == slug {
			return s, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockSiteRepository) FindByDomain(ctx context.Context, domainName string) (*domain.Site, error) {
	for _, s := range m.sites {
		if s.Domain != nil && *s.Domain == domainName {
			return s, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockSiteRepository) FindAll(ctx context.Context, filter domain.SiteFilter) ([]*domain.Site, int, error) {
	var sites []*domain.Site
	for _, s := range m.sites {
		sites = append(sites, s)
	}
	return sites, len(sites), nil
}

func (m *mockSiteRepository) Create(ctx context.Context, site *domain.Site) error {
	site.CreatedAt = time.Now()
	site.UpdatedAt = time.Now()
	m.sites[site.ID] = site
	return nil
}

func (m *mockSiteRepository) Update(ctx context.Context, site *domain.Site) error {
	if _, ok := m.sites[site.ID]; !ok {
		return domain.ErrNotFound
	}
	site.UpdatedAt = time.Now()
	m.sites[site.ID] = site
	return nil
}

func (m *mockSiteRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if _, ok := m.sites[id]; !ok {
		return domain.ErrNotFound
	}
	delete(m.sites, id)
	return nil
}

func (m *mockSiteRepository) FindSettingsBySiteID(ctx context.Context, siteID uuid.UUID, publicOnly bool) ([]*domain.SiteSetting, error) {
	var settings []*domain.SiteSetting
	prefix := siteID.String() + ":"
	for k, s := range m.settings {
		if len(k) > len(prefix) && k[:len(prefix)] == prefix {
			if !publicOnly || s.IsPublic {
				settings = append(settings, s)
			}
		}
	}
	return settings, nil
}

func (m *mockSiteRepository) FindSettingByKey(ctx context.Context, siteID uuid.UUID, key string) (*domain.SiteSetting, error) {
	k := siteID.String() + ":" + key
	if s, ok := m.settings[k]; ok {
		return s, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockSiteRepository) UpsertSetting(ctx context.Context, siteID uuid.UUID, key, value string) error {
	k := siteID.String() + ":" + key
	if existing, ok := m.settings[k]; ok {
		existing.Value = &value
	} else {
		m.settings[k] = &domain.SiteSetting{
			ID:     uuid.New(),
			SiteID: siteID,
			Key:    key,
			Value:  &value,
		}
	}
	return nil
}

func (m *mockSiteRepository) BulkUpsertSettings(ctx context.Context, siteID uuid.UUID, settings map[string]string) error {
	for key, value := range settings {
		m.UpsertSetting(ctx, siteID, key, value)
	}
	return nil
}

// ─── Tests ────────────────────────────────────────────────────────────────────

func createTestSiteService(repo *mockSiteRepository) service.SiteService {
	logger := zerolog.Nop()
	return service.NewSiteService(repo, logger)
}

func TestSiteService_CreateSite_Success(t *testing.T) {
	repo := newMockSiteRepository()
	svc := createTestSiteService(repo)

	userID := uuid.New()
	input := domain.CreateSiteInput{
		Name: "Test Site",
		Slug: "test-site",
	}

	site, err := svc.CreateSite(context.Background(), input, userID)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if site == nil {
		t.Fatal("expected site, got nil")
	}
	if site.Name != "Test Site" {
		t.Errorf("expected name 'Test Site', got '%s'", site.Name)
	}
	if site.Slug != "test-site" {
		t.Errorf("expected slug 'test-site', got '%s'", site.Slug)
	}
	if !site.IsActive {
		t.Error("expected site to be active by default")
	}
}

func TestSiteService_GetSite_NotFound(t *testing.T) {
	repo := newMockSiteRepository()
	svc := createTestSiteService(repo)

	_, err := svc.GetSite(context.Background(), uuid.New())
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}
}

func TestSiteService_GetSiteBySlug_Success(t *testing.T) {
	repo := newMockSiteRepository()
	svc := createTestSiteService(repo)

	userID := uuid.New()
	input := domain.CreateSiteInput{Name: "My Site", Slug: "my-site"}
	created, _ := svc.CreateSite(context.Background(), input, userID)

	found, err := svc.GetSiteBySlug(context.Background(), "my-site")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if found.ID != created.ID {
		t.Errorf("expected same site ID")
	}
}

func TestSiteService_UpdateSite_Success(t *testing.T) {
	repo := newMockSiteRepository()
	svc := createTestSiteService(repo)

	userID := uuid.New()
	input := domain.CreateSiteInput{Name: "Original Name", Slug: "original-slug"}
	site, _ := svc.CreateSite(context.Background(), input, userID)

	newName := "Updated Name"
	isActive := false
	updateInput := domain.UpdateSiteInput{
		Name:     &newName,
		IsActive: &isActive,
	}

	updated, err := svc.UpdateSite(context.Background(), site.ID, updateInput)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if updated.Name != "Updated Name" {
		t.Errorf("expected name 'Updated Name', got '%s'", updated.Name)
	}
	if updated.IsActive {
		t.Error("expected site to be inactive after update")
	}
	// Slug should remain unchanged
	if updated.Slug != "original-slug" {
		t.Errorf("expected slug 'original-slug', got '%s'", updated.Slug)
	}
}

func TestSiteService_DeleteSite_Success(t *testing.T) {
	repo := newMockSiteRepository()
	svc := createTestSiteService(repo)

	userID := uuid.New()
	input := domain.CreateSiteInput{Name: "Site to Delete", Slug: "delete-me"}
	site, _ := svc.CreateSite(context.Background(), input, userID)

	err := svc.DeleteSite(context.Background(), site.ID)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	_, err = svc.GetSite(context.Background(), site.ID)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("expected ErrNotFound after delete, got: %v", err)
	}
}

func TestSiteService_UpdateSetting_Success(t *testing.T) {
	repo := newMockSiteRepository()
	svc := createTestSiteService(repo)

	userID := uuid.New()
	input := domain.CreateSiteInput{Name: "Test Site", Slug: "test-site"}
	site, _ := svc.CreateSite(context.Background(), input, userID)

	value := "My Awesome Product"
	err := svc.UpdateSetting(context.Background(), site.ID, "site_title", domain.UpdateSettingInput{
		Value: &value,
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Verify setting was saved
	settings, err := svc.GetSettings(context.Background(), site.ID, false)
	if err != nil {
		t.Fatalf("expected no error getting settings, got: %v", err)
	}

	found := false
	for _, s := range settings {
		if s.Key == "site_title" && s.Value != nil && *s.Value == "My Awesome Product" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected to find updated setting")
	}
}

func TestSiteService_BulkUpdateSettings_Success(t *testing.T) {
	repo := newMockSiteRepository()
	svc := createTestSiteService(repo)

	userID := uuid.New()
	input := domain.CreateSiteInput{Name: "Test Site", Slug: "test-site"}
	site, _ := svc.CreateSite(context.Background(), input, userID)

	err := svc.BulkUpdateSettings(context.Background(), site.ID, domain.BulkUpdateSettingsInput{
		Settings: map[string]string{
			"site_title":       "My Product",
			"seo_title":        "My Product - Best Solution",
			"seo_description":  "The best solution for your needs",
		},
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	settingsMap, err := svc.GetSettingsMap(context.Background(), site.ID, false)
	if err != nil {
		t.Fatalf("expected no error getting settings map, got: %v", err)
	}

	if settingsMap["site_title"] != "My Product" {
		t.Errorf("expected site_title 'My Product', got '%s'", settingsMap["site_title"])
	}
	if settingsMap["seo_title"] != "My Product - Best Solution" {
		t.Errorf("expected seo_title 'My Product - Best Solution', got '%s'", settingsMap["seo_title"])
	}
}

func TestSiteService_ListSites_Success(t *testing.T) {
	repo := newMockSiteRepository()
	svc := createTestSiteService(repo)

	userID := uuid.New()

	// Create multiple sites
	svc.CreateSite(context.Background(), domain.CreateSiteInput{Name: "Site 1", Slug: "site-1"}, userID)
	svc.CreateSite(context.Background(), domain.CreateSiteInput{Name: "Site 2", Slug: "site-2"}, userID)
	svc.CreateSite(context.Background(), domain.CreateSiteInput{Name: "Site 3", Slug: "site-3"}, userID)

	result, err := svc.ListSites(context.Background(), domain.SiteFilter{})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if result.Total != 3 {
		t.Errorf("expected 3 sites, got %d", result.Total)
	}
}
