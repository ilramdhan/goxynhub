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

// SiteHandler handles site-related endpoints
type SiteHandler struct {
	siteService service.SiteService
	logger      zerolog.Logger
}

// NewSiteHandler creates a new SiteHandler
func NewSiteHandler(siteService service.SiteService, logger zerolog.Logger) *SiteHandler {
	return &SiteHandler{
		siteService: siteService,
		logger:      logger,
	}
}

// GetPublicSite handles GET /api/v1/public/site/:slug
func (h *SiteHandler) GetPublicSite(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		slug = "default"
	}

	site, err := h.siteService.GetSiteWithSettings(c.Request.Context(), uuid.Nil, true)
	if err != nil {
		// Try by slug
		siteBySlug, slugErr := h.siteService.GetSiteBySlug(c.Request.Context(), slug)
		if slugErr != nil {
			if errors.Is(slugErr, domain.ErrNotFound) {
				response.NotFound(c, "site not found")
				return
			}
			h.logger.Error().Err(slugErr).Msg("get public site error")
			response.InternalError(c, slugErr)
			return
		}

		// Get with public settings
		siteWithSettings, settingsErr := h.siteService.GetSiteWithSettings(c.Request.Context(), siteBySlug.ID, true)
		if settingsErr != nil {
			response.OK(c, siteBySlug)
			return
		}
		response.OK(c, siteWithSettings)
		return
	}

	response.OK(c, site)
}

// GetPublicSiteByID handles GET /api/v1/public/sites/:id
func (h *SiteHandler) GetPublicSiteByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid site ID")
		return
	}

	site, err := h.siteService.GetSiteWithSettings(c.Request.Context(), id, true)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "site not found")
			return
		}
		h.logger.Error().Err(err).Msg("get public site error")
		response.InternalError(c, err)
		return
	}

	response.OK(c, site)
}

// ListSites handles GET /api/v1/admin/sites
func (h *SiteHandler) ListSites(c *gin.Context) {
	var filter domain.SiteFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.BadRequest(c, "invalid query parameters")
		return
	}

	search := c.Query("search")
	if search != "" {
		filter.Search = &search
	}

	result, err := h.siteService.ListSites(c.Request.Context(), filter)
	if err != nil {
		h.logger.Error().Err(err).Msg("list sites error")
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

// GetSite handles GET /api/v1/admin/sites/:id
func (h *SiteHandler) GetSite(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid site ID")
		return
	}

	site, err := h.siteService.GetSiteWithSettings(c.Request.Context(), id, false)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "site not found")
			return
		}
		h.logger.Error().Err(err).Msg("get site error")
		response.InternalError(c, err)
		return
	}

	response.OK(c, site)
}

// CreateSite handles POST /api/v1/admin/sites
func (h *SiteHandler) CreateSite(c *gin.Context) {
	userIDVal, _ := c.Get(middleware.ContextKeyUserID)
	userID, _ := userIDVal.(uuid.UUID)

	var input domain.CreateSiteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	site, err := h.siteService.CreateSite(c.Request.Context(), input, userID)
	if err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			response.Conflict(c, "a site with this slug already exists")
			return
		}
		h.logger.Error().Err(err).Msg("create site error")
		response.InternalError(c, err)
		return
	}

	response.Created(c, site)
}

// UpdateSite handles PUT /api/v1/admin/sites/:id
func (h *SiteHandler) UpdateSite(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid site ID")
		return
	}

	var input domain.UpdateSiteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	site, err := h.siteService.UpdateSite(c.Request.Context(), id, input)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "site not found")
			return
		}
		h.logger.Error().Err(err).Msg("update site error")
		response.InternalError(c, err)
		return
	}

	response.OK(c, site)
}

// DeleteSite handles DELETE /api/v1/admin/sites/:id
func (h *SiteHandler) DeleteSite(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid site ID")
		return
	}

	if err := h.siteService.DeleteSite(c.Request.Context(), id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "site not found")
			return
		}
		h.logger.Error().Err(err).Msg("delete site error")
		response.InternalError(c, err)
		return
	}

	response.NoContent(c)
}

// GetSettings handles GET /api/v1/admin/sites/:id/settings
func (h *SiteHandler) GetSettings(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid site ID")
		return
	}

	settings, err := h.siteService.GetSettings(c.Request.Context(), id, false)
	if err != nil {
		h.logger.Error().Err(err).Msg("get settings error")
		response.InternalError(c, err)
		return
	}

	response.OK(c, settings)
}

// UpdateSetting handles PUT /api/v1/admin/sites/:id/settings/:key
func (h *SiteHandler) UpdateSetting(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid site ID")
		return
	}

	key := c.Param("key")
	if key == "" {
		response.BadRequest(c, "setting key is required")
		return
	}

	var input domain.UpdateSettingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	if err := h.siteService.UpdateSetting(c.Request.Context(), id, key, input); err != nil {
		h.logger.Error().Err(err).Msg("update setting error")
		response.InternalError(c, err)
		return
	}

	response.OKWithMessage(c, "setting updated successfully", nil)
}

// BulkUpdateSettings handles PUT /api/v1/admin/sites/:id/settings
func (h *SiteHandler) BulkUpdateSettings(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid site ID")
		return
	}

	var input domain.BulkUpdateSettingsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	if err := h.siteService.BulkUpdateSettings(c.Request.Context(), id, input); err != nil {
		h.logger.Error().Err(err).Msg("bulk update settings error")
		response.InternalError(c, err)
		return
	}

	// Return updated settings
	settings, err := h.siteService.GetSettings(c.Request.Context(), id, false)
	if err != nil {
		response.OKWithMessage(c, "settings updated successfully", nil)
		return
	}

	response.OKWithMessage(c, "settings updated successfully", settings)
}
