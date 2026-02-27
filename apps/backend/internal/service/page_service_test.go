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

// ─── Mock PageRepository ──────────────────────────────────────────────────────

type mockPageRepository struct {
	pages    map[uuid.UUID]*domain.Page
	sections map[uuid.UUID]*domain.PageSection
	contents map[uuid.UUID]*domain.SectionContent
}

func newMockPageRepository() *mockPageRepository {
	return &mockPageRepository{
		pages:    make(map[uuid.UUID]*domain.Page),
		sections: make(map[uuid.UUID]*domain.PageSection),
		contents: make(map[uuid.UUID]*domain.SectionContent),
	}
}

func (m *mockPageRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Page, error) {
	if p, ok := m.pages[id]; ok {
		return p, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockPageRepository) FindBySlug(ctx context.Context, siteID uuid.UUID, slug string) (*domain.Page, error) {
	for _, p := range m.pages {
		if p.SiteID == siteID && p.Slug == slug {
			return p, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockPageRepository) FindHomepage(ctx context.Context, siteID uuid.UUID) (*domain.Page, error) {
	for _, p := range m.pages {
		if p.SiteID == siteID && p.IsHomepage && p.Status == domain.PageStatusPublished {
			return p, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockPageRepository) FindAll(ctx context.Context, filter domain.PageFilter) ([]*domain.Page, int, error) {
	var pages []*domain.Page
	for _, p := range m.pages {
		if filter.SiteID != nil && p.SiteID != *filter.SiteID {
			continue
		}
		if filter.Status != nil && p.Status != *filter.Status {
			continue
		}
		pages = append(pages, p)
	}
	return pages, len(pages), nil
}

func (m *mockPageRepository) Create(ctx context.Context, page *domain.Page) error {
	page.CreatedAt = time.Now()
	page.UpdatedAt = time.Now()
	m.pages[page.ID] = page
	return nil
}

func (m *mockPageRepository) Update(ctx context.Context, page *domain.Page) error {
	if _, ok := m.pages[page.ID]; !ok {
		return domain.ErrNotFound
	}
	page.UpdatedAt = time.Now()
	m.pages[page.ID] = page
	return nil
}

func (m *mockPageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if _, ok := m.pages[id]; !ok {
		return domain.ErrNotFound
	}
	delete(m.pages, id)
	return nil
}

func (m *mockPageRepository) SetHomepage(ctx context.Context, siteID, pageID uuid.UUID) error {
	for _, p := range m.pages {
		if p.SiteID == siteID {
			p.IsHomepage = p.ID == pageID
		}
	}
	return nil
}

func (m *mockPageRepository) FindSectionsByPageID(ctx context.Context, pageID uuid.UUID) ([]*domain.PageSection, error) {
	var sections []*domain.PageSection
	for _, s := range m.sections {
		if s.PageID == pageID {
			sections = append(sections, s)
		}
	}
	return sections, nil
}

func (m *mockPageRepository) FindSectionByID(ctx context.Context, id uuid.UUID) (*domain.PageSection, error) {
	if s, ok := m.sections[id]; ok {
		return s, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockPageRepository) CreateSection(ctx context.Context, section *domain.PageSection) error {
	section.CreatedAt = time.Now()
	section.UpdatedAt = time.Now()
	m.sections[section.ID] = section
	return nil
}

func (m *mockPageRepository) UpdateSection(ctx context.Context, section *domain.PageSection) error {
	if _, ok := m.sections[section.ID]; !ok {
		return domain.ErrNotFound
	}
	section.UpdatedAt = time.Now()
	m.sections[section.ID] = section
	return nil
}

func (m *mockPageRepository) DeleteSection(ctx context.Context, id uuid.UUID) error {
	if _, ok := m.sections[id]; !ok {
		return domain.ErrNotFound
	}
	delete(m.sections, id)
	return nil
}

func (m *mockPageRepository) ReorderSections(ctx context.Context, orders []domain.SectionOrder) error {
	for _, order := range orders {
		if s, ok := m.sections[order.ID]; ok {
			s.SortOrder = order.SortOrder
		}
	}
	return nil
}

func (m *mockPageRepository) FindContentsBySectionID(ctx context.Context, sectionID uuid.UUID) ([]*domain.SectionContent, error) {
	var contents []*domain.SectionContent
	for _, c := range m.contents {
		if c.SectionID == sectionID {
			contents = append(contents, c)
		}
	}
	return contents, nil
}

func (m *mockPageRepository) FindContentByKey(ctx context.Context, sectionID uuid.UUID, key string) (*domain.SectionContent, error) {
	for _, c := range m.contents {
		if c.SectionID == sectionID && c.Key == key {
			return c, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockPageRepository) UpsertContent(ctx context.Context, content *domain.SectionContent) error {
	// Check if exists
	for _, c := range m.contents {
		if c.SectionID == content.SectionID && c.Key == content.Key {
			content.ID = c.ID
			content.CreatedAt = c.CreatedAt
			content.UpdatedAt = time.Now()
			m.contents[c.ID] = content
			return nil
		}
	}
	// Create new
	content.CreatedAt = time.Now()
	content.UpdatedAt = time.Now()
	m.contents[content.ID] = content
	return nil
}

func (m *mockPageRepository) DeleteContent(ctx context.Context, id uuid.UUID) error {
	if _, ok := m.contents[id]; !ok {
		return domain.ErrNotFound
	}
	delete(m.contents, id)
	return nil
}

func (m *mockPageRepository) DeleteContentByKey(ctx context.Context, sectionID uuid.UUID, key string) error {
	for id, c := range m.contents {
		if c.SectionID == sectionID && c.Key == key {
			delete(m.contents, id)
			return nil
		}
	}
	return nil
}

// ─── Tests ────────────────────────────────────────────────────────────────────

func createTestPageService(repo *mockPageRepository) service.PageService {
	logger := zerolog.Nop()
	return service.NewPageService(repo, logger)
}

func TestPageService_CreatePage_Success(t *testing.T) {
	repo := newMockPageRepository()
	svc := createTestPageService(repo)

	siteID := uuid.New()
	userID := uuid.New()
	input := domain.CreatePageInput{
		SiteID: siteID,
		Title:  "Test Page",
		Slug:   "test-page",
		Status: domain.PageStatusDraft,
	}

	page, err := svc.CreatePage(context.Background(), input, userID)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if page == nil {
		t.Fatal("expected page, got nil")
	}
	if page.Title != "Test Page" {
		t.Errorf("expected title 'Test Page', got '%s'", page.Title)
	}
	if page.Slug != "test-page" {
		t.Errorf("expected slug 'test-page', got '%s'", page.Slug)
	}
	if page.Status != domain.PageStatusDraft {
		t.Errorf("expected status 'draft', got '%s'", page.Status)
	}
}

func TestPageService_CreatePage_SlugGeneration(t *testing.T) {
	repo := newMockPageRepository()
	svc := createTestPageService(repo)

	siteID := uuid.New()
	userID := uuid.New()
	input := domain.CreatePageInput{
		SiteID: siteID,
		Title:  "My Awesome Page!",
		Slug:   "", // Empty slug - should be auto-generated
		Status: domain.PageStatusDraft,
	}

	page, err := svc.CreatePage(context.Background(), input, userID)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if page.Slug != "my-awesome-page" {
		t.Errorf("expected slug 'my-awesome-page', got '%s'", page.Slug)
	}
}

func TestPageService_GetPage_NotFound(t *testing.T) {
	repo := newMockPageRepository()
	svc := createTestPageService(repo)

	_, err := svc.GetPage(context.Background(), uuid.New())
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}
}

func TestPageService_UpdatePage_Success(t *testing.T) {
	repo := newMockPageRepository()
	svc := createTestPageService(repo)

	siteID := uuid.New()
	userID := uuid.New()

	// Create page first
	createInput := domain.CreatePageInput{
		SiteID: siteID,
		Title:  "Original Title",
		Slug:   "original-slug",
		Status: domain.PageStatusDraft,
	}
	page, _ := svc.CreatePage(context.Background(), createInput, userID)

	// Update page
	newTitle := "Updated Title"
	newStatus := domain.PageStatusPublished
	updateInput := domain.UpdatePageInput{
		Title:  &newTitle,
		Status: &newStatus,
	}

	updatedPage, err := svc.UpdatePage(context.Background(), page.ID, updateInput, userID)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if updatedPage.Title != "Updated Title" {
		t.Errorf("expected title 'Updated Title', got '%s'", updatedPage.Title)
	}
	if updatedPage.Status != domain.PageStatusPublished {
		t.Errorf("expected status 'published', got '%s'", updatedPage.Status)
	}
	// Slug should remain unchanged
	if updatedPage.Slug != "original-slug" {
		t.Errorf("expected slug 'original-slug', got '%s'", updatedPage.Slug)
	}
}

func TestPageService_DeletePage_Success(t *testing.T) {
	repo := newMockPageRepository()
	svc := createTestPageService(repo)

	siteID := uuid.New()
	userID := uuid.New()

	// Create page
	createInput := domain.CreatePageInput{
		SiteID: siteID,
		Title:  "Page to Delete",
		Slug:   "page-to-delete",
		Status: domain.PageStatusDraft,
	}
	page, _ := svc.CreatePage(context.Background(), createInput, userID)

	// Delete page
	err := svc.DeletePage(context.Background(), page.ID)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Verify page is gone
	_, err = svc.GetPage(context.Background(), page.ID)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("expected ErrNotFound after delete, got: %v", err)
	}
}

func TestPageService_PublishPage_Success(t *testing.T) {
	repo := newMockPageRepository()
	svc := createTestPageService(repo)

	siteID := uuid.New()
	userID := uuid.New()

	// Create draft page
	createInput := domain.CreatePageInput{
		SiteID: siteID,
		Title:  "Draft Page",
		Slug:   "draft-page",
		Status: domain.PageStatusDraft,
	}
	page, _ := svc.CreatePage(context.Background(), createInput, userID)

	// Publish
	publishedPage, err := svc.PublishPage(context.Background(), page.ID, userID)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if publishedPage.Status != domain.PageStatusPublished {
		t.Errorf("expected status 'published', got '%s'", publishedPage.Status)
	}
}

func TestPageService_CreateSection_Success(t *testing.T) {
	repo := newMockPageRepository()
	svc := createTestPageService(repo)

	pageID := uuid.New()
	input := domain.CreateSectionInput{
		PageID:    pageID,
		Name:      "Hero Section",
		Type:      domain.SectionTypeHero,
		IsVisible: true,
		SortOrder: 1,
	}

	section, err := svc.CreateSection(context.Background(), input)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if section.Name != "Hero Section" {
		t.Errorf("expected name 'Hero Section', got '%s'", section.Name)
	}
	if section.Type != domain.SectionTypeHero {
		t.Errorf("expected type 'hero', got '%s'", section.Type)
	}
}

func TestPageService_UpsertContent_CreateNew(t *testing.T) {
	repo := newMockPageRepository()
	svc := createTestPageService(repo)

	sectionID := uuid.New()
	value := "Hello World"
	input := domain.UpsertContentInput{
		Key:   "title",
		Value: &value,
		Type:  domain.ContentTypeText,
	}

	content, err := svc.UpsertContent(context.Background(), sectionID, input)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if content.Key != "title" {
		t.Errorf("expected key 'title', got '%s'", content.Key)
	}
	if *content.Value != "Hello World" {
		t.Errorf("expected value 'Hello World', got '%s'", *content.Value)
	}
}

func TestPageService_UpsertContent_UpdateExisting(t *testing.T) {
	repo := newMockPageRepository()
	svc := createTestPageService(repo)

	sectionID := uuid.New()

	// Create initial content
	value1 := "Original Value"
	input1 := domain.UpsertContentInput{
		Key:   "title",
		Value: &value1,
		Type:  domain.ContentTypeText,
	}
	svc.UpsertContent(context.Background(), sectionID, input1)

	// Update content
	value2 := "Updated Value"
	input2 := domain.UpsertContentInput{
		Key:   "title",
		Value: &value2,
		Type:  domain.ContentTypeText,
	}
	content, err := svc.UpsertContent(context.Background(), sectionID, input2)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if *content.Value != "Updated Value" {
		t.Errorf("expected value 'Updated Value', got '%s'", *content.Value)
	}

	// Verify only one content item exists
	contents, _ := svc.GetSectionContents(context.Background(), sectionID)
	if len(contents) != 1 {
		t.Errorf("expected 1 content item, got %d", len(contents))
	}
}

func TestPageService_ReorderSections(t *testing.T) {
	repo := newMockPageRepository()
	svc := createTestPageService(repo)

	pageID := uuid.New()

	// Create sections
	s1, _ := svc.CreateSection(context.Background(), domain.CreateSectionInput{
		PageID: pageID, Name: "Section 1", Type: domain.SectionTypeHero, SortOrder: 1,
	})
	s2, _ := svc.CreateSection(context.Background(), domain.CreateSectionInput{
		PageID: pageID, Name: "Section 2", Type: domain.SectionTypeFeatures, SortOrder: 2,
	})

	// Reorder
	err := svc.ReorderSections(context.Background(), domain.ReorderSectionsInput{
		Sections: []domain.SectionOrder{
			{ID: s1.ID, SortOrder: 2},
			{ID: s2.ID, SortOrder: 1},
		},
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Verify new order
	updatedS1, _ := svc.GetSection(context.Background(), s1.ID)
	updatedS2, _ := svc.GetSection(context.Background(), s2.ID)

	if updatedS1.SortOrder != 2 {
		t.Errorf("expected s1 sort_order 2, got %d", updatedS1.SortOrder)
	}
	if updatedS2.SortOrder != 1 {
		t.Errorf("expected s2 sort_order 1, got %d", updatedS2.SortOrder)
	}
}
