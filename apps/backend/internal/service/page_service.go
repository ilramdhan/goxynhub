package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/domain"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/repository"
)

// PageService defines the interface for page operations
type PageService interface {
	GetPage(ctx context.Context, id uuid.UUID) (*domain.Page, error)
	GetPageBySlug(ctx context.Context, siteID uuid.UUID, slug string) (*domain.Page, error)
	GetHomepage(ctx context.Context, siteID uuid.UUID) (*domain.Page, error)
	GetPageWithContent(ctx context.Context, siteID uuid.UUID, slug string) (*domain.Page, error)
	ListPages(ctx context.Context, filter domain.PageFilter) (*domain.PaginatedResult[*domain.Page], error)
	CreatePage(ctx context.Context, input domain.CreatePageInput, userID uuid.UUID) (*domain.Page, error)
	UpdatePage(ctx context.Context, id uuid.UUID, input domain.UpdatePageInput, userID uuid.UUID) (*domain.Page, error)
	DeletePage(ctx context.Context, id uuid.UUID) error
	PublishPage(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Page, error)
	UnpublishPage(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Page, error)

	// Section operations
	GetSection(ctx context.Context, id uuid.UUID) (*domain.PageSection, error)
	ListSections(ctx context.Context, pageID uuid.UUID) ([]*domain.PageSection, error)
	CreateSection(ctx context.Context, input domain.CreateSectionInput) (*domain.PageSection, error)
	UpdateSection(ctx context.Context, id uuid.UUID, input domain.UpdateSectionInput) (*domain.PageSection, error)
	DeleteSection(ctx context.Context, id uuid.UUID) error
	ReorderSections(ctx context.Context, input domain.ReorderSectionsInput) error

	// Content operations
	GetSectionContents(ctx context.Context, sectionID uuid.UUID) ([]*domain.SectionContent, error)
	UpsertContent(ctx context.Context, sectionID uuid.UUID, input domain.UpsertContentInput) (*domain.SectionContent, error)
	DeleteContent(ctx context.Context, id uuid.UUID) error
	BulkUpsertContents(ctx context.Context, sectionID uuid.UUID, inputs []domain.UpsertContentInput) ([]*domain.SectionContent, error)
}

// pageService implements PageService
type pageService struct {
	pageRepo repository.PageRepository
	logger   zerolog.Logger
}

// NewPageService creates a new pageService
func NewPageService(pageRepo repository.PageRepository, logger zerolog.Logger) PageService {
	return &pageService{
		pageRepo: pageRepo,
		logger:   logger,
	}
}

// GetPage retrieves a page by ID
func (s *pageService) GetPage(ctx context.Context, id uuid.UUID) (*domain.Page, error) {
	page, err := s.pageRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("pageService.GetPage: %w", err)
	}
	return page, nil
}

// GetPageBySlug retrieves a page by slug
func (s *pageService) GetPageBySlug(ctx context.Context, siteID uuid.UUID, slug string) (*domain.Page, error) {
	page, err := s.pageRepo.FindBySlug(ctx, siteID, slug)
	if err != nil {
		return nil, fmt.Errorf("pageService.GetPageBySlug: %w", err)
	}
	return page, nil
}

// GetHomepage retrieves the homepage for a site
func (s *pageService) GetHomepage(ctx context.Context, siteID uuid.UUID) (*domain.Page, error) {
	page, err := s.pageRepo.FindHomepage(ctx, siteID)
	if err != nil {
		return nil, fmt.Errorf("pageService.GetHomepage: %w", err)
	}
	return page, nil
}

// GetPageWithContent retrieves a page with all its sections and content
func (s *pageService) GetPageWithContent(ctx context.Context, siteID uuid.UUID, slug string) (*domain.Page, error) {
	var page *domain.Page
	var err error

	if slug == "" || slug == "home" {
		page, err = s.pageRepo.FindHomepage(ctx, siteID)
	} else {
		page, err = s.pageRepo.FindBySlug(ctx, siteID, slug)
	}

	if err != nil {
		return nil, fmt.Errorf("pageService.GetPageWithContent: %w", err)
	}

	// Load sections
	sections, err := s.pageRepo.FindSectionsByPageID(ctx, page.ID)
	if err != nil {
		return nil, fmt.Errorf("pageService.GetPageWithContent sections: %w", err)
	}

	// Load content for each section
	for _, section := range sections {
		contents, err := s.pageRepo.FindContentsBySectionID(ctx, section.ID)
		if err != nil {
			return nil, fmt.Errorf("pageService.GetPageWithContent contents: %w", err)
		}
		section.Contents = contents
	}

	page.Sections = sections
	return page, nil
}

// ListPages retrieves all pages with optional filtering
func (s *pageService) ListPages(ctx context.Context, filter domain.PageFilter) (*domain.PaginatedResult[*domain.Page], error) {
	pages, total, err := s.pageRepo.FindAll(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("pageService.ListPages: %w", err)
	}

	result := domain.NewPaginatedResult(pages, total, filter.Pagination)
	return &result, nil
}

// CreatePage creates a new page
func (s *pageService) CreatePage(ctx context.Context, input domain.CreatePageInput, userID uuid.UUID) (*domain.Page, error) {
	// Generate slug if not provided
	slug := input.Slug
	if slug == "" {
		slug = generateSlug(input.Title)
	} else {
		slug = normalizeSlug(slug)
	}

	page := &domain.Page{
		ID:             uuid.New(),
		SiteID:         input.SiteID,
		Title:          input.Title,
		Slug:           slug,
		Description:    input.Description,
		Status:         input.Status,
		IsHomepage:     input.IsHomepage,
		SEOTitle:       input.SEOTitle,
		SEODescription: input.SEODescription,
		SEOKeywords:    input.SEOKeywords,
		OGTitle:        input.OGTitle,
		OGDescription:  input.OGDescription,
		OGImage:        input.OGImage,
		OGType:         input.OGType,
		TwitterTitle:       input.TwitterTitle,
		TwitterDescription: input.TwitterDescription,
		TwitterImage:       input.TwitterImage,
		TwitterCard:        input.TwitterCard,
		CanonicalURL:   input.CanonicalURL,
		RobotsMeta:     input.RobotsMeta,
		Template:       input.Template,
		CreatedBy:      &userID,
		UpdatedBy:      &userID,
	}

	if err := s.pageRepo.Create(ctx, page); err != nil {
		return nil, fmt.Errorf("pageService.CreatePage: %w", err)
	}

	// If this is set as homepage, update other pages
	if input.IsHomepage {
		if err := s.pageRepo.SetHomepage(ctx, input.SiteID, page.ID); err != nil {
			s.logger.Error().Err(err).Msg("failed to set homepage")
		}
	}

	s.logger.Info().
		Str("page_id", page.ID.String()).
		Str("slug", page.Slug).
		Str("user_id", userID.String()).
		Msg("page created")

	return page, nil
}

// UpdatePage updates an existing page
func (s *pageService) UpdatePage(ctx context.Context, id uuid.UUID, input domain.UpdatePageInput, userID uuid.UUID) (*domain.Page, error) {
	page, err := s.pageRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("pageService.UpdatePage find: %w", err)
	}

	// Apply updates
	if input.Title != nil {
		page.Title = *input.Title
	}
	if input.Slug != nil {
		page.Slug = normalizeSlug(*input.Slug)
	}
	if input.Description != nil {
		page.Description = input.Description
	}
	if input.Status != nil {
		page.Status = *input.Status
	}
	if input.IsHomepage != nil {
		page.IsHomepage = *input.IsHomepage
	}
	if input.SEOTitle != nil {
		page.SEOTitle = input.SEOTitle
	}
	if input.SEODescription != nil {
		page.SEODescription = input.SEODescription
	}
	if input.SEOKeywords != nil {
		page.SEOKeywords = input.SEOKeywords
	}
	if input.OGTitle != nil {
		page.OGTitle = input.OGTitle
	}
	if input.OGDescription != nil {
		page.OGDescription = input.OGDescription
	}
	if input.OGImage != nil {
		page.OGImage = input.OGImage
	}
	if input.OGType != nil {
		page.OGType = input.OGType
	}
	if input.TwitterTitle != nil {
		page.TwitterTitle = input.TwitterTitle
	}
	if input.TwitterDescription != nil {
		page.TwitterDescription = input.TwitterDescription
	}
	if input.TwitterImage != nil {
		page.TwitterImage = input.TwitterImage
	}
	if input.TwitterCard != nil {
		page.TwitterCard = input.TwitterCard
	}
	if input.CanonicalURL != nil {
		page.CanonicalURL = input.CanonicalURL
	}
	if input.RobotsMeta != nil {
		page.RobotsMeta = input.RobotsMeta
	}
	if input.Template != nil {
		page.Template = input.Template
	}
	if input.CustomHead != nil {
		page.CustomHead = input.CustomHead
	}
	page.UpdatedBy = &userID

	if err := s.pageRepo.Update(ctx, page); err != nil {
		return nil, fmt.Errorf("pageService.UpdatePage: %w", err)
	}

	// Handle homepage change
	if input.IsHomepage != nil && *input.IsHomepage {
		if err := s.pageRepo.SetHomepage(ctx, page.SiteID, page.ID); err != nil {
			s.logger.Error().Err(err).Msg("failed to set homepage")
		}
	}

	return page, nil
}

// DeletePage soft-deletes a page
func (s *pageService) DeletePage(ctx context.Context, id uuid.UUID) error {
	if err := s.pageRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("pageService.DeletePage: %w", err)
	}
	return nil
}

// PublishPage publishes a page
func (s *pageService) PublishPage(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Page, error) {
	status := domain.PageStatusPublished
	return s.UpdatePage(ctx, id, domain.UpdatePageInput{Status: &status}, userID)
}

// UnpublishPage unpublishes a page
func (s *pageService) UnpublishPage(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Page, error) {
	status := domain.PageStatusDraft
	return s.UpdatePage(ctx, id, domain.UpdatePageInput{Status: &status}, userID)
}

// GetSection retrieves a section by ID
func (s *pageService) GetSection(ctx context.Context, id uuid.UUID) (*domain.PageSection, error) {
	section, err := s.pageRepo.FindSectionByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("pageService.GetSection: %w", err)
	}
	return section, nil
}

// ListSections retrieves all sections for a page
func (s *pageService) ListSections(ctx context.Context, pageID uuid.UUID) ([]*domain.PageSection, error) {
	sections, err := s.pageRepo.FindSectionsByPageID(ctx, pageID)
	if err != nil {
		return nil, fmt.Errorf("pageService.ListSections: %w", err)
	}
	return sections, nil
}

// CreateSection creates a new page section
func (s *pageService) CreateSection(ctx context.Context, input domain.CreateSectionInput) (*domain.PageSection, error) {
	section := &domain.PageSection{
		ID:         uuid.New(),
		PageID:     input.PageID,
		Name:       input.Name,
		Type:       input.Type,
		Identifier: input.Identifier,
		IsVisible:  input.IsVisible,
		SortOrder:  input.SortOrder,
	}

	if err := s.pageRepo.CreateSection(ctx, section); err != nil {
		return nil, fmt.Errorf("pageService.CreateSection: %w", err)
	}

	return section, nil
}

// UpdateSection updates an existing section
func (s *pageService) UpdateSection(ctx context.Context, id uuid.UUID, input domain.UpdateSectionInput) (*domain.PageSection, error) {
	section, err := s.pageRepo.FindSectionByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("pageService.UpdateSection find: %w", err)
	}

	if input.Name != nil {
		section.Name = *input.Name
	}
	if input.Type != nil {
		section.Type = *input.Type
	}
	if input.IsVisible != nil {
		section.IsVisible = *input.IsVisible
	}
	if input.SortOrder != nil {
		section.SortOrder = *input.SortOrder
	}
	if input.BGColor != nil {
		section.BGColor = input.BGColor
	}
	if input.BGImage != nil {
		section.BGImage = input.BGImage
	}
	if input.BGVideo != nil {
		section.BGVideo = input.BGVideo
	}
	if input.BGOverlay != nil {
		section.BGOverlay = *input.BGOverlay
	}
	if input.BGOverlayColor != nil {
		section.BGOverlayColor = input.BGOverlayColor
	}
	if input.BGOverlayOpacity != nil {
		section.BGOverlayOpacity = input.BGOverlayOpacity
	}
	if input.Layout != nil {
		section.Layout = input.Layout
	}
	if input.PaddingTop != nil {
		section.PaddingTop = input.PaddingTop
	}
	if input.PaddingBottom != nil {
		section.PaddingBottom = input.PaddingBottom
	}
	if input.Animation != nil {
		section.Animation = input.Animation
	}
	if input.CSSClass != nil {
		section.CSSClass = input.CSSClass
	}
	if input.CustomCSS != nil {
		section.CustomCSS = input.CustomCSS
	}

	if err := s.pageRepo.UpdateSection(ctx, section); err != nil {
		return nil, fmt.Errorf("pageService.UpdateSection: %w", err)
	}

	return section, nil
}

// DeleteSection deletes a section
func (s *pageService) DeleteSection(ctx context.Context, id uuid.UUID) error {
	if err := s.pageRepo.DeleteSection(ctx, id); err != nil {
		return fmt.Errorf("pageService.DeleteSection: %w", err)
	}
	return nil
}

// ReorderSections reorders sections
func (s *pageService) ReorderSections(ctx context.Context, input domain.ReorderSectionsInput) error {
	if err := s.pageRepo.ReorderSections(ctx, input.Sections); err != nil {
		return fmt.Errorf("pageService.ReorderSections: %w", err)
	}
	return nil
}

// GetSectionContents retrieves all content for a section
func (s *pageService) GetSectionContents(ctx context.Context, sectionID uuid.UUID) ([]*domain.SectionContent, error) {
	contents, err := s.pageRepo.FindContentsBySectionID(ctx, sectionID)
	if err != nil {
		return nil, fmt.Errorf("pageService.GetSectionContents: %w", err)
	}
	return contents, nil
}

// UpsertContent creates or updates a content item
func (s *pageService) UpsertContent(ctx context.Context, sectionID uuid.UUID, input domain.UpsertContentInput) (*domain.SectionContent, error) {
	content := &domain.SectionContent{
		ID:          uuid.New(),
		SectionID:   sectionID,
		Key:         input.Key,
		Value:       input.Value,
		ValueJSON:   input.ValueJSON,
		Type:        input.Type,
		Label:       input.Label,
		Description: input.Description,
		Placeholder: input.Placeholder,
		IsRequired:  input.IsRequired,
		SortOrder:   input.SortOrder,
		AltText:     input.AltText,
		Width:       input.Width,
		Height:      input.Height,
		LinkURL:     input.LinkURL,
		LinkTarget:  input.LinkTarget,
	}

	if err := s.pageRepo.UpsertContent(ctx, content); err != nil {
		return nil, fmt.Errorf("pageService.UpsertContent: %w", err)
	}

	return content, nil
}

// DeleteContent deletes a content item
func (s *pageService) DeleteContent(ctx context.Context, id uuid.UUID) error {
	if err := s.pageRepo.DeleteContent(ctx, id); err != nil {
		return fmt.Errorf("pageService.DeleteContent: %w", err)
	}
	return nil
}

// BulkUpsertContents creates or updates multiple content items
func (s *pageService) BulkUpsertContents(ctx context.Context, sectionID uuid.UUID, inputs []domain.UpsertContentInput) ([]*domain.SectionContent, error) {
	var results []*domain.SectionContent
	for _, input := range inputs {
		content, err := s.UpsertContent(ctx, sectionID, input)
		if err != nil {
			return nil, fmt.Errorf("pageService.BulkUpsertContents: %w", err)
		}
		results = append(results, content)
	}
	return results, nil
}

// generateSlug creates a URL-friendly slug from a title
func generateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = regexp.MustCompile(`[^a-z0-9\s-]`).ReplaceAllString(slug, "")
	slug = regexp.MustCompile(`[\s-]+`).ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	return slug
}

// normalizeSlug ensures a slug is URL-friendly
func normalizeSlug(slug string) string {
	return generateSlug(slug)
}
