package handler

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/domain"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/pkg/response"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/repository"
)

// UserHandler handles user management endpoints
type UserHandler struct {
	userRepo repository.UserRepository
	logger   zerolog.Logger
	bcryptCost int
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userRepo repository.UserRepository, logger zerolog.Logger, bcryptCost int) *UserHandler {
	return &UserHandler{
		userRepo:   userRepo,
		logger:     logger,
		bcryptCost: bcryptCost,
	}
}

// ListUsers handles GET /api/v1/admin/users
func (h *UserHandler) ListUsers(c *gin.Context) {
	var filter domain.UserFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.BadRequest(c, "invalid query parameters")
		return
	}

	search := c.Query("search")
	if search != "" {
		filter.Search = &search
	}
	roleStr := c.Query("role")
	if roleStr != "" {
		role := domain.UserRole(roleStr)
		filter.Role = &role
	}
	statusStr := c.Query("status")
	if statusStr != "" {
		status := domain.UserStatus(statusStr)
		filter.Status = &status
	}

	users, total, err := h.userRepo.FindAll(c.Request.Context(), filter)
	if err != nil {
		h.logger.Error().Err(err).Msg("list users error")
		response.InternalError(c, err)
		return
	}

	filter.Normalize()
	totalPages := total / filter.PerPage
	if total%filter.PerPage > 0 {
		totalPages++
	}

	response.OKPaginated(c, users, gin.H{
		"page":        filter.Page,
		"per_page":    filter.PerPage,
		"total":       total,
		"total_pages": totalPages,
	})
}

// GetUser handles GET /api/v1/admin/users/:id
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid user ID")
		return
	}

	user, err := h.userRepo.FindByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "user not found")
			return
		}
		h.logger.Error().Err(err).Msg("get user error")
		response.InternalError(c, err)
		return
	}

	response.OK(c, user)
}

// CreateUser handles POST /api/v1/admin/users
func (h *UserHandler) CreateUser(c *gin.Context) {
	var input domain.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	// Check if email already exists
	_, err := h.userRepo.FindByEmail(c.Request.Context(), input.Email)
	if err == nil {
		response.Conflict(c, "a user with this email already exists")
		return
	}
	if !errors.Is(err, domain.ErrNotFound) {
		h.logger.Error().Err(err).Msg("check email error")
		response.InternalError(c, err)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), h.bcryptCost)
	if err != nil {
		h.logger.Error().Err(err).Msg("hash password error")
		response.InternalError(c, fmt.Errorf("failed to hash password"))
		return
	}

	user := &domain.User{
		ID:            uuid.New(),
		Email:         input.Email,
		PasswordHash:  string(hashedPassword),
		FullName:      input.FullName,
		Role:          input.Role,
		Status:        domain.StatusActive,
		EmailVerified: true,
	}

	if err := h.userRepo.Create(c.Request.Context(), user); err != nil {
		h.logger.Error().Err(err).Msg("create user error")
		response.InternalError(c, err)
		return
	}

	response.Created(c, user)
}

// UpdateUser handles PUT /api/v1/admin/users/:id
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid user ID")
		return
	}

	var input domain.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	user, err := h.userRepo.FindByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "user not found")
			return
		}
		h.logger.Error().Err(err).Msg("get user error")
		response.InternalError(c, err)
		return
	}

	if input.FullName != nil {
		user.FullName = *input.FullName
	}
	if input.AvatarURL != nil {
		user.AvatarURL = input.AvatarURL
	}
	if input.Role != nil {
		user.Role = *input.Role
	}
	if input.Status != nil {
		user.Status = *input.Status
	}

	if err := h.userRepo.Update(c.Request.Context(), user); err != nil {
		h.logger.Error().Err(err).Msg("update user error")
		response.InternalError(c, err)
		return
	}

	response.OK(c, user)
}

// DeleteUser handles DELETE /api/v1/admin/users/:id
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid user ID")
		return
	}

	if err := h.userRepo.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "user not found")
			return
		}
		h.logger.Error().Err(err).Msg("delete user error")
		response.InternalError(c, err)
		return
	}

	response.NoContent(c)
}
