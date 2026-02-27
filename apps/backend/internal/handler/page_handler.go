package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/domain"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/middleware"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/pkg/response"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/service"
)

// PageHandler handles page-related endpoints
type PageHandler struct {
	pageService service.PageService
	logger      zerolog.Logger
}

// NewPageHandler creates a new PageHandler
func NewPageHandler(pageService service.PageService, logger zerolog.Logger) *PageHandler {
	return &PageHandler{
		pageService: pageService,
		logger:      logger,
	}
}

// GetPublicPage handles GET /api/v1/public/pages/:slug
func (h *PageHandler) GetPublicPage(c *gin.Context) {
	siteIDStr := c.Query("site_id")
	slug := c.Param("slug")

	var siteID uuid.UUID
	var err error
	if siteIDStr != "" {
		siteID, err = uuid.Parse(siteIDStr)
		if err != nil {
			response.BadRequest(c, "invalid site_id")
			return
		}
	} else {
		// Use default site - in production, resolve from domain
		// For now, we'll require site_id or use a default
		response.BadRequest(c, "site_id is required")
		return
	}

	page, err := h.pageService.GetPageWithContent(c.Request.Context(), siteID, slug)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "page not found")
			return
		}
		h.logger.Error().Err(err).Str("slug", slug).Msg("get public page error")
		response.InternalError(c, err)
		return
	}

	// Only return published pages for public API
	if page.Status != domain.PageStatusPublished {
		response.NotFound(c, "page not found")
		return
	}

	response.OK(c, page)
}

// ListPages handles GET /api/v1/admin/pages
func (h *PageHandler) ListPages(c *gin.Context) {
	var filter domain.PageFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.BadRequest(c, "invalid query parameters")
		return
	}

	siteIDStr := c.Query("site_id")
	if siteIDStr != "" {
		siteID, err := uuid.Parse(siteIDStr)
		if err != nil {
			response.BadRequest(c, "invalid site_id")
			return
		}
		filter.SiteID = &siteID
	}

	statusStr := c.Query("status")
	if statusStr != "" {
		status := domain.PageStatus(statusStr)
		filter.Status = &status
	}

	search := c.Query("search")
	if search != "" {
		filter.Search = &search
	}

	result, err := h.pageService.ListPages(c.Request.Context(), filter)
	if err != nil {
		h.logger.Error().Err(err).Msg("list pages error")
		response.InternalError(c, err)
		return
	}

	response.OKPaginated(c, result.Data, gin.H{
		"page":        result.Page,
		"per_page":    result.PerPage,
		"total":       result.Total,
		"total_pages": result.TotalPages,
	})
}

// GetPage handles GET /api/v1/admin/pages/:id
func (h *PageHandler) GetPage(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid page ID")
		return
	}

	page, err := h.pageService.GetPage(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "page not found")
			return
		}
		h.logger.Error().Err(err).Str("id", id.String()).Msg("get page error")
		response.InternalError(c, err)
		return
	}

	response.OK(c, page)
}

// CreatePage handles POST /api/v1/admin/pages
func (h *PageHandler) CreatePage(c *gin.Context) {
	userIDVal, _ := c.Get(middleware.ContextKeyUserID)
	userID, _ := userIDVal.(uuid.UUID)

	var input domain.CreatePageInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	page, err := h.pageService.CreatePage(c.Request.Context(), input, userID)
	if err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			response.Conflict(c, "a page with this slug already exists")
			return
		}
		h.logger.Error().Err(err).Msg("create page error")
		response.InternalError(c, err)
		return
	}

	response.Created(c, page)
}

// UpdatePage handles PUT /api/v1/admin/pages/:id
func (h *PageHandler) UpdatePage(c *gin.Context) {
	userIDVal, _ := c.Get(middleware.ContextKeyUserID)
	userID, _ := userIDVal.(uuid.UUID)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid page ID")
		return
	}

	var input domain.UpdatePageInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	page, err := h.pageService.UpdatePage(c.Request.Context(), id, input, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "page not found")
			return
		}
		if errors.Is(err, domain.ErrAlreadyExists) {
			response.Conflict(c, "a page with this slug already exists")
			return
		}
		h.logger.Error().Err(err).Str("id", id.String()).Msg("update page error")
		response.InternalError(c, err)
		return
	}

	response.OK(c, page)
}

// DeletePage handles DELETE /api/v1/admin/pages/:id
func (h *PageHandler) DeletePage(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid page ID")
		return
	}

	if err := h.pageService.DeletePage(c.Request.Context(), id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "page not found")
			return
		}
		h.logger.Error().Err(err).Str("id", id.String()).Msg("delete page error")
		response.InternalError(c, err)
		return
	}

	response.NoContent(c)
}

// PublishPage handles PATCH /api/v1/admin/pages/:id/publish
func (h *PageHandler) PublishPage(c *gin.Context) {
	userIDVal, _ := c.Get(middleware.ContextKeyUserID)
	userID, _ := userIDVal.(uuid.UUID)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid page ID")
		return
	}

	page, err := h.pageService.PublishPage(c.Request.Context(), id, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "page not found")
			return
		}
		h.logger.Error().Err(err).Str("id", id.String()).Msg("publish page error")
		response.InternalError(c, err)
		return
	}

	response.OKWithMessage(c, "page published successfully", page)
}

// UnpublishPage handles PATCH /api/v1/admin/pages/:id/unpublish
func (h *PageHandler) UnpublishPage(c *gin.Context) {
	userIDVal, _ := c.Get(middleware.ContextKeyUserID)
	userID, _ := userIDVal.(uuid.UUID)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid page ID")
		return
	}

	page, err := h.pageService.UnpublishPage(c.Request.Context(), id, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "page not found")
			return
		}
		h.logger.Error().Err(err).Str("id", id.String()).Msg("unpublish page error")
		response.InternalError(c, err)
		return
	}

	response.OKWithMessage(c, "page unpublished successfully", page)
}

// ListSections handles GET /api/v1/admin/pages/:id/sections
func (h *PageHandler) ListSections(c *gin.Context) {
	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid page ID")
		return
	}

	sections, err := h.pageService.ListSections(c.Request.Context(), pageID)
	if err != nil {
		h.logger.Error().Err(err).Str("page_id", pageID.String()).Msg("list sections error")
		response.InternalError(c, err)
		return
	}

	response.OK(c, sections)
}

// CreateSection handles POST /api/v1/admin/pages/:id/sections
func (h *PageHandler) CreateSection(c *gin.Context) {
	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid page ID")
		return
	}

	var input domain.CreateSectionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}
	input.PageID = pageID

	section, err := h.pageService.CreateSection(c.Request.Context(), input)
	if err != nil {
		h.logger.Error().Err(err).Msg("create section error")
		response.InternalError(c, err)
		return
	}

	response.Created(c, section)
}

// UpdateSection handles PUT /api/v1/admin/sections/:id
func (h *PageHandler) UpdateSection(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid section ID")
		return
	}

	var input domain.UpdateSectionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	section, err := h.pageService.UpdateSection(c.Request.Context(), id, input)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "section not found")
			return
		}
		h.logger.Error().Err(err).Str("id", id.String()).Msg("update section error")
		response.InternalError(c, err)
		return
	}

	response.OK(c, section)
}

// DeleteSection handles DELETE /api/v1/admin/sections/:id
func (h *PageHandler) DeleteSection(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid section ID")
		return
	}

	if err := h.pageService.DeleteSection(c.Request.Context(), id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "section not found")
			return
		}
		h.logger.Error().Err(err).Str("id", id.String()).Msg("delete section error")
		response.InternalError(c, err)
		return
	}

	response.NoContent(c)
}

// ReorderSections handles PATCH /api/v1/admin/sections/reorder
func (h *PageHandler) ReorderSections(c *gin.Context) {
	var input domain.ReorderSectionsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	if err := h.pageService.ReorderSections(c.Request.Context(), input); err != nil {
		h.logger.Error().Err(err).Msg("reorder sections error")
		response.InternalError(c, err)
		return
	}

	response.OKWithMessage(c, "sections reordered successfully", nil)
}

// ListContents handles GET /api/v1/admin/sections/:id/contents
func (h *PageHandler) ListContents(c *gin.Context) {
	sectionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid section ID")
		return
	}

	contents, err := h.pageService.GetSectionContents(c.Request.Context(), sectionID)
	if err != nil {
		h.logger.Error().Err(err).Str("section_id", sectionID.String()).Msg("list contents error")
		response.InternalError(c, err)
		return
	}

	response.OK(c, contents)
}

// UpsertContent handles POST /api/v1/admin/sections/:id/contents
func (h *PageHandler) UpsertContent(c *gin.Context) {
	sectionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid section ID")
		return
	}

	var input domain.UpsertContentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	content, err := h.pageService.UpsertContent(c.Request.Context(), sectionID, input)
	if err != nil {
		h.logger.Error().Err(err).Msg("upsert content error")
		response.InternalError(c, err)
		return
	}

	response.OK(c, content)
}

// BulkUpsertContents handles POST /api/v1/admin/sections/:id/contents/bulk
func (h *PageHandler) BulkUpsertContents(c *gin.Context) {
	sectionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid section ID")
		return
	}

	var inputs []domain.UpsertContentInput
	if err := c.ShouldBindJSON(&inputs); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	contents, err := h.pageService.BulkUpsertContents(c.Request.Context(), sectionID, inputs)
	if err != nil {
		h.logger.Error().Err(err).Msg("bulk upsert contents error")
		response.InternalError(c, err)
		return
	}

	response.OK(c, contents)
}

// DeleteContent handles DELETE /api/v1/admin/contents/:id
func (h *PageHandler) DeleteContent(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid content ID")
		return
	}

	if err := h.pageService.DeleteContent(c.Request.Context(), id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "content not found")
			return
		}
		h.logger.Error().Err(err).Str("id", id.String()).Msg("delete content error")
		response.InternalError(c, err)
		return
	}

	response.NoContent(c)
}
