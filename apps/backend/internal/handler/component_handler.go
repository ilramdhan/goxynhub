package handler

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/domain"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/middleware"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/pkg/response"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/repository"
)

// ComponentHandler handles component-related endpoints
type ComponentHandler struct {
	compRepo        repository.ComponentRepository
	supabaseURL     string
	supabaseBucket  string
	supabaseKey     string
	maxUploadSize   int64
	allowedMimeTypes []string
	logger          zerolog.Logger
}

// NewComponentHandler creates a new ComponentHandler
func NewComponentHandler(
	compRepo repository.ComponentRepository,
	supabaseURL, supabaseBucket, supabaseKey string,
	maxUploadSize int64,
	allowedMimeTypes []string,
	logger zerolog.Logger,
) *ComponentHandler {
	return &ComponentHandler{
		compRepo:        compRepo,
		supabaseURL:     supabaseURL,
		supabaseBucket:  supabaseBucket,
		supabaseKey:     supabaseKey,
		maxUploadSize:   maxUploadSize,
		allowedMimeTypes: allowedMimeTypes,
		logger:          logger,
	}
}

// ─── Features ─────────────────────────────────────────────────────────────────

// ListFeatures handles GET /api/v1/admin/features
func (h *ComponentHandler) ListFeatures(c *gin.Context) {
	var filter domain.ComponentFilter
	siteIDStr := c.Query("site_id")
	if siteIDStr != "" {
		siteID, err := uuid.Parse(siteIDStr)
		if err != nil {
			response.BadRequest(c, "invalid site_id")
			return
		}
		filter.SiteID = &siteID
	}

	features, total, err := h.compRepo.FindFeaturesByFilter(c.Request.Context(), filter)
	if err != nil {
		h.logger.Error().Err(err).Msg("list features error")
		response.InternalError(c, err)
		return
	}

	filter.Normalize()
	response.OKPaginated(c, features, gin.H{
		"total": total, "page": filter.Page, "per_page": filter.PerPage,
	})
}

// CreateFeature handles POST /api/v1/admin/features
func (h *ComponentHandler) CreateFeature(c *gin.Context) {
	var input domain.CreateFeatureInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	feature := &domain.Feature{
		ID:          uuid.New(),
		SiteID:      input.SiteID,
		SectionID:   input.SectionID,
		Title:       input.Title,
		Description: input.Description,
		Icon:        input.Icon,
		IconColor:   input.IconColor,
		ImageURL:    input.ImageURL,
		ImageAlt:    input.ImageAlt,
		LinkURL:     input.LinkURL,
		LinkText:    input.LinkText,
		IsActive:    input.IsActive,
		SortOrder:   input.SortOrder,
	}

	if err := h.compRepo.CreateFeature(c.Request.Context(), feature); err != nil {
		h.logger.Error().Err(err).Msg("create feature error")
		response.InternalError(c, err)
		return
	}

	response.Created(c, feature)
}

// UpdateFeature handles PUT /api/v1/admin/features/:id
func (h *ComponentHandler) UpdateFeature(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid feature ID")
		return
	}

	feature, err := h.compRepo.FindFeatureByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "feature not found")
			return
		}
		response.InternalError(c, err)
		return
	}

	if err := c.ShouldBindJSON(feature); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	if err := h.compRepo.UpdateFeature(c.Request.Context(), feature); err != nil {
		h.logger.Error().Err(err).Msg("update feature error")
		response.InternalError(c, err)
		return
	}

	response.OK(c, feature)
}

// DeleteFeature handles DELETE /api/v1/admin/features/:id
func (h *ComponentHandler) DeleteFeature(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid feature ID")
		return
	}

	if err := h.compRepo.DeleteFeature(c.Request.Context(), id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "feature not found")
			return
		}
		response.InternalError(c, err)
		return
	}

	response.NoContent(c)
}

// ─── Testimonials ─────────────────────────────────────────────────────────────

// ListTestimonials handles GET /api/v1/admin/testimonials
func (h *ComponentHandler) ListTestimonials(c *gin.Context) {
	var filter domain.ComponentFilter
	siteIDStr := c.Query("site_id")
	if siteIDStr != "" {
		siteID, _ := uuid.Parse(siteIDStr)
		filter.SiteID = &siteID
	}

	testimonials, total, err := h.compRepo.FindTestimonialsByFilter(c.Request.Context(), filter)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	filter.Normalize()
	response.OKPaginated(c, testimonials, gin.H{"total": total, "page": filter.Page, "per_page": filter.PerPage})
}

// CreateTestimonial handles POST /api/v1/admin/testimonials
func (h *ComponentHandler) CreateTestimonial(c *gin.Context) {
	var input domain.CreateTestimonialInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	t := &domain.Testimonial{
		ID:            uuid.New(),
		SiteID:        input.SiteID,
		SectionID:     input.SectionID,
		AuthorName:    input.AuthorName,
		AuthorTitle:   input.AuthorTitle,
		AuthorCompany: input.AuthorCompany,
		AuthorAvatar:  input.AuthorAvatar,
		Content:       input.Content,
		Rating:        input.Rating,
		Source:        input.Source,
		SourceURL:     input.SourceURL,
		IsFeatured:    input.IsFeatured,
		IsActive:      input.IsActive,
		SortOrder:     input.SortOrder,
	}

	if err := h.compRepo.CreateTestimonial(c.Request.Context(), t); err != nil {
		response.InternalError(c, err)
		return
	}

	response.Created(c, t)
}

// UpdateTestimonial handles PUT /api/v1/admin/testimonials/:id
func (h *ComponentHandler) UpdateTestimonial(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid testimonial ID")
		return
	}

	t, err := h.compRepo.FindTestimonialByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "testimonial not found")
			return
		}
		response.InternalError(c, err)
		return
	}

	if err := c.ShouldBindJSON(t); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	if err := h.compRepo.UpdateTestimonial(c.Request.Context(), t); err != nil {
		response.InternalError(c, err)
		return
	}

	response.OK(c, t)
}

// DeleteTestimonial handles DELETE /api/v1/admin/testimonials/:id
func (h *ComponentHandler) DeleteTestimonial(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid testimonial ID")
		return
	}

	if err := h.compRepo.DeleteTestimonial(c.Request.Context(), id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "testimonial not found")
			return
		}
		response.InternalError(c, err)
		return
	}

	response.NoContent(c)
}

// ─── Pricing Plans ────────────────────────────────────────────────────────────

// ListPricingPlans handles GET /api/v1/admin/pricing
func (h *ComponentHandler) ListPricingPlans(c *gin.Context) {
	var filter domain.ComponentFilter
	siteIDStr := c.Query("site_id")
	if siteIDStr != "" {
		siteID, _ := uuid.Parse(siteIDStr)
		filter.SiteID = &siteID
	}

	plans, total, err := h.compRepo.FindPricingPlansByFilter(c.Request.Context(), filter)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	filter.Normalize()
	response.OKPaginated(c, plans, gin.H{"total": total, "page": filter.Page, "per_page": filter.PerPage})
}

// CreatePricingPlan handles POST /api/v1/admin/pricing
func (h *ComponentHandler) CreatePricingPlan(c *gin.Context) {
	var input domain.CreatePricingPlanInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	features := make(domain.JSONArray, len(input.Features))
	for i, f := range input.Features {
		features[i] = f
	}
	featuresExcluded := make(domain.JSONArray, len(input.FeaturesExcluded))
	for i, f := range input.FeaturesExcluded {
		featuresExcluded[i] = f
	}

	plan := &domain.PricingPlan{
		ID:               uuid.New(),
		SiteID:           input.SiteID,
		SectionID:        input.SectionID,
		Name:             input.Name,
		Description:      input.Description,
		PriceMonthly:     input.PriceMonthly,
		PriceYearly:      input.PriceYearly,
		Currency:         input.Currency,
		PriceLabel:       input.PriceLabel,
		IsPopular:        input.IsPopular,
		IsCustom:         input.IsCustom,
		BadgeText:        input.BadgeText,
		CTAText:          input.CTAText,
		CTALink:          input.CTALink,
		Features:         features,
		FeaturesExcluded: featuresExcluded,
		IsActive:         input.IsActive,
		SortOrder:        input.SortOrder,
	}

	if err := h.compRepo.CreatePricingPlan(c.Request.Context(), plan); err != nil {
		response.InternalError(c, err)
		return
	}

	response.Created(c, plan)
}

// UpdatePricingPlan handles PUT /api/v1/admin/pricing/:id
func (h *ComponentHandler) UpdatePricingPlan(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid pricing plan ID")
		return
	}

	plan, err := h.compRepo.FindPricingPlanByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "pricing plan not found")
			return
		}
		response.InternalError(c, err)
		return
	}

	if err := c.ShouldBindJSON(plan); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	if err := h.compRepo.UpdatePricingPlan(c.Request.Context(), plan); err != nil {
		response.InternalError(c, err)
		return
	}

	response.OK(c, plan)
}

// DeletePricingPlan handles DELETE /api/v1/admin/pricing/:id
func (h *ComponentHandler) DeletePricingPlan(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid pricing plan ID")
		return
	}

	if err := h.compRepo.DeletePricingPlan(c.Request.Context(), id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "pricing plan not found")
			return
		}
		response.InternalError(c, err)
		return
	}

	response.NoContent(c)
}

// ─── FAQs ─────────────────────────────────────────────────────────────────────

// ListFAQs handles GET /api/v1/admin/faqs
func (h *ComponentHandler) ListFAQs(c *gin.Context) {
	var filter domain.ComponentFilter
	siteIDStr := c.Query("site_id")
	if siteIDStr != "" {
		siteID, _ := uuid.Parse(siteIDStr)
		filter.SiteID = &siteID
	}

	faqs, total, err := h.compRepo.FindFAQsByFilter(c.Request.Context(), filter)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	filter.Normalize()
	response.OKPaginated(c, faqs, gin.H{"total": total, "page": filter.Page, "per_page": filter.PerPage})
}

// CreateFAQ handles POST /api/v1/admin/faqs
func (h *ComponentHandler) CreateFAQ(c *gin.Context) {
	var input domain.CreateFAQInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	faq := &domain.FAQ{
		ID:        uuid.New(),
		SiteID:    input.SiteID,
		SectionID: input.SectionID,
		Question:  input.Question,
		Answer:    input.Answer,
		Category:  input.Category,
		IsActive:  input.IsActive,
		SortOrder: input.SortOrder,
	}

	if err := h.compRepo.CreateFAQ(c.Request.Context(), faq); err != nil {
		response.InternalError(c, err)
		return
	}

	response.Created(c, faq)
}

// UpdateFAQ handles PUT /api/v1/admin/faqs/:id
func (h *ComponentHandler) UpdateFAQ(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid FAQ ID")
		return
	}

	faq, err := h.compRepo.FindFAQByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "FAQ not found")
			return
		}
		response.InternalError(c, err)
		return
	}

	if err := c.ShouldBindJSON(faq); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	if err := h.compRepo.UpdateFAQ(c.Request.Context(), faq); err != nil {
		response.InternalError(c, err)
		return
	}

	response.OK(c, faq)
}

// DeleteFAQ handles DELETE /api/v1/admin/faqs/:id
func (h *ComponentHandler) DeleteFAQ(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid FAQ ID")
		return
	}

	if err := h.compRepo.DeleteFAQ(c.Request.Context(), id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "FAQ not found")
			return
		}
		response.InternalError(c, err)
		return
	}

	response.NoContent(c)
}

// ─── Navigation ───────────────────────────────────────────────────────────────

// GetNavigation handles GET /api/v1/public/navigation/:siteId/:identifier
func (h *ComponentHandler) GetNavigation(c *gin.Context) {
	siteID, err := uuid.Parse(c.Param("siteId"))
	if err != nil {
		response.BadRequest(c, "invalid site ID")
		return
	}

	identifier := c.Param("identifier")
	if identifier == "" {
		identifier = "header"
	}

	menu, err := h.compRepo.FindMenuByIdentifier(c.Request.Context(), siteID, identifier)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "navigation menu not found")
			return
		}
		response.InternalError(c, err)
		return
	}

	items, err := h.compRepo.FindItemsByMenuID(c.Request.Context(), menu.ID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	// Build tree structure
	menu.Items = buildNavigationTree(items)
	response.OK(c, menu)
}

// ListNavigation handles GET /api/v1/admin/navigation
func (h *ComponentHandler) ListNavigation(c *gin.Context) {
	siteIDStr := c.Query("site_id")
	if siteIDStr == "" {
		response.BadRequest(c, "site_id is required")
		return
	}

	siteID, err := uuid.Parse(siteIDStr)
	if err != nil {
		response.BadRequest(c, "invalid site_id")
		return
	}

	menus, err := h.compRepo.FindMenusBySiteID(c.Request.Context(), siteID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	// Load items for each menu
	for _, menu := range menus {
		items, err := h.compRepo.FindItemsByMenuID(c.Request.Context(), menu.ID)
		if err == nil {
			menu.Items = buildNavigationTree(items)
		}
	}

	response.OK(c, menus)
}

// CreateNavigationItem handles POST /api/v1/admin/navigation/:menuId/items
func (h *ComponentHandler) CreateNavigationItem(c *gin.Context) {
	menuID, err := uuid.Parse(c.Param("menuId"))
	if err != nil {
		response.BadRequest(c, "invalid menu ID")
		return
	}

	var item domain.NavigationItem
	if err := c.ShouldBindJSON(&item); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	item.ID = uuid.New()
	item.MenuID = menuID

	if err := h.compRepo.CreateNavigationItem(c.Request.Context(), &item); err != nil {
		response.InternalError(c, err)
		return
	}

	response.Created(c, item)
}

// UpdateNavigationItem handles PUT /api/v1/admin/navigation/items/:id
func (h *ComponentHandler) UpdateNavigationItem(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid item ID")
		return
	}

	var item domain.NavigationItem
	if err := c.ShouldBindJSON(&item); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}
	item.ID = id

	if err := h.compRepo.UpdateNavigationItem(c.Request.Context(), &item); err != nil {
		response.InternalError(c, err)
		return
	}

	response.OK(c, item)
}

// DeleteNavigationItem handles DELETE /api/v1/admin/navigation/items/:id
func (h *ComponentHandler) DeleteNavigationItem(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid item ID")
		return
	}

	if err := h.compRepo.DeleteNavigationItem(c.Request.Context(), id); err != nil {
		response.InternalError(c, err)
		return
	}

	response.NoContent(c)
}

// ─── Media ────────────────────────────────────────────────────────────────────

// ListMedia handles GET /api/v1/admin/media
func (h *ComponentHandler) ListMedia(c *gin.Context) {
	siteIDStr := c.Query("site_id")
	if siteIDStr == "" {
		response.BadRequest(c, "site_id is required")
		return
	}

	siteID, err := uuid.Parse(siteIDStr)
	if err != nil {
		response.BadRequest(c, "invalid site_id")
		return
	}

	var pagination domain.Pagination
	if err := c.ShouldBindQuery(&pagination); err != nil {
		pagination = domain.Pagination{Page: 1, PerPage: 20}
	}

	media, total, err := h.compRepo.FindMediaByFilter(c.Request.Context(), siteID, pagination)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	pagination.Normalize()
	response.OKPaginated(c, media, gin.H{"total": total, "page": pagination.Page, "per_page": pagination.PerPage})
}

// UploadMedia handles POST /api/v1/admin/media/upload
func (h *ComponentHandler) UploadMedia(c *gin.Context) {
	userIDVal, _ := c.Get(middleware.ContextKeyUserID)
	userID, _ := userIDVal.(uuid.UUID)

	siteIDStr := c.PostForm("site_id")
	if siteIDStr == "" {
		response.BadRequest(c, "site_id is required")
		return
	}

	siteID, err := uuid.Parse(siteIDStr)
	if err != nil {
		response.BadRequest(c, "invalid site_id")
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.BadRequest(c, "file is required")
		return
	}
	defer file.Close()

	// Validate file size
	if header.Size > h.maxUploadSize {
		response.BadRequest(c, fmt.Sprintf("file size exceeds maximum allowed size of %d bytes", h.maxUploadSize))
		return
	}

	// Validate MIME type
	mimeType := header.Header.Get("Content-Type")
	if !h.isAllowedMimeType(mimeType) {
		response.BadRequest(c, "file type not allowed")
		return
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	fileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	folder := c.PostForm("folder")
	if folder == "" {
		folder = "/"
	}
	filePath := fmt.Sprintf("%s/%s", strings.TrimPrefix(folder, "/"), fileName)

	// Upload to Supabase Storage
	publicURL, err := h.uploadToSupabase(file, filePath, mimeType)
	if err != nil {
		h.logger.Error().Err(err).Msg("upload to supabase error")
		response.InternalError(c, fmt.Errorf("failed to upload file"))
		return
	}

	// Determine media type
	mediaType := getMediaType(mimeType)

	// Save to database
	media := &domain.Media{
		ID:           uuid.New(),
		SiteID:       siteID,
		Name:         strings.TrimSuffix(header.Filename, ext),
		OriginalName: header.Filename,
		FilePath:     filePath,
		PublicURL:    publicURL,
		Type:         mediaType,
		MimeType:     mimeType,
		FileSize:     header.Size,
		Folder:       folder,
		UploadedBy:   &userID,
	}

	altText := c.PostForm("alt_text")
	if altText != "" {
		media.AltText = &altText
	}

	if err := h.compRepo.CreateMedia(c.Request.Context(), media); err != nil {
		h.logger.Error().Err(err).Msg("save media error")
		response.InternalError(c, err)
		return
	}

	response.Created(c, media)
}

// UpdateMedia handles PUT /api/v1/admin/media/:id
func (h *ComponentHandler) UpdateMedia(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid media ID")
		return
	}

	media, err := h.compRepo.FindMediaByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "media not found")
			return
		}
		response.InternalError(c, err)
		return
	}

	var updateData struct {
		Name    *string `json:"name"`
		AltText *string `json:"alt_text"`
		Caption *string `json:"caption"`
		Folder  *string `json:"folder"`
	}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	if updateData.Name != nil {
		media.Name = *updateData.Name
	}
	if updateData.AltText != nil {
		media.AltText = updateData.AltText
	}
	if updateData.Caption != nil {
		media.Caption = updateData.Caption
	}

	if err := h.compRepo.UpdateMedia(c.Request.Context(), media); err != nil {
		response.InternalError(c, err)
		return
	}

	response.OK(c, media)
}

// DeleteMedia handles DELETE /api/v1/admin/media/:id
func (h *ComponentHandler) DeleteMedia(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid media ID")
		return
	}

	if err := h.compRepo.DeleteMedia(c.Request.Context(), id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "media not found")
			return
		}
		response.InternalError(c, err)
		return
	}

	response.NoContent(c)
}

// ─── Audit Logs ───────────────────────────────────────────────────────────────

// ListAuditLogs handles GET /api/v1/admin/audit-logs
func (h *ComponentHandler) ListAuditLogs(c *gin.Context) {
	var filter domain.AuditLogFilter

	siteIDStr := c.Query("site_id")
	if siteIDStr != "" {
		siteID, _ := uuid.Parse(siteIDStr)
		filter.SiteID = &siteID
	}

	action := c.Query("action")
	if action != "" {
		filter.Action = &action
	}

	resourceType := c.Query("resource_type")
	if resourceType != "" {
		filter.ResourceType = &resourceType
	}

	logs, total, err := h.compRepo.FindAuditLogs(c.Request.Context(), filter)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	filter.Normalize()
	response.OKPaginated(c, logs, gin.H{"total": total, "page": filter.Page, "per_page": filter.PerPage})
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func (h *ComponentHandler) isAllowedMimeType(mimeType string) bool {
	for _, allowed := range h.allowedMimeTypes {
		if strings.EqualFold(allowed, mimeType) {
			return true
		}
	}
	return false
}

func (h *ComponentHandler) uploadToSupabase(file multipart.File, filePath, mimeType string) (string, error) {
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("read file: %w", err)
	}

	uploadURL := fmt.Sprintf("%s/storage/v1/object/%s/%s", h.supabaseURL, h.supabaseBucket, filePath)

	req, err := http.NewRequest(http.MethodPost, uploadURL, strings.NewReader(string(fileBytes)))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+h.supabaseKey)
	req.Header.Set("Content-Type", mimeType)
	req.Header.Set("x-upsert", "true")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("upload request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("upload failed with status: %d", resp.StatusCode)
	}

	publicURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", h.supabaseURL, h.supabaseBucket, filePath)
	return publicURL, nil
}

func getMediaType(mimeType string) string {
	switch {
	case strings.HasPrefix(mimeType, "image/"):
		return "image"
	case strings.HasPrefix(mimeType, "video/"):
		return "video"
	case strings.HasPrefix(mimeType, "audio/"):
		return "audio"
	case mimeType == "application/pdf":
		return "document"
	default:
		return "other"
	}
}

func buildNavigationTree(items []*domain.NavigationItem) []*domain.NavigationItem {
	itemMap := make(map[uuid.UUID]*domain.NavigationItem)
	for _, item := range items {
		itemMap[item.ID] = item
		item.Children = []*domain.NavigationItem{}
	}

	var roots []*domain.NavigationItem
	for _, item := range items {
		if item.ParentID == nil {
			roots = append(roots, item)
		} else {
			if parent, ok := itemMap[*item.ParentID]; ok {
				parent.Children = append(parent.Children, item)
			}
		}
	}
	return roots
}
