package handler

import (
	"errors"
	"net/http"
	"strconv"

	"le-studio-api/internal/dto"
	"le-studio-api/internal/service"
	"le-studio-api/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// AdminHandler serves admin dashboard endpoints.
type AdminHandler struct {
	svc *service.AdminService
	v   *validator.Validate
}

// NewAdminHandler creates AdminHandler.
func NewAdminHandler(svc *service.AdminService, v *validator.Validate) *AdminHandler {
	return &AdminHandler{svc: svc, v: v}
}

// StatsOverview returns overview stats.
func (h *AdminHandler) StatsOverview(c *gin.Context) {
	stats, err := h.svc.StatsOverview(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "STATS_FAILED", "Unable to load stats.", nil)
		return
	}
	response.OK(c, stats)
}

// GetMe returns current admin profile.
func (h *AdminHandler) GetMe(c *gin.Context) {
	adminID, err := strconv.ParseUint(c.GetString("userID"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token.", nil)
		return
	}
	admin, err := h.svc.GetProfile(c.Request.Context(), uint(adminID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Admin not found.", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "ADMIN_FETCH_FAILED", "Unable to load admin.", nil)
		return
	}
	response.OK(c, admin)
}

// UpdateMe updates current admin profile.
func (h *AdminHandler) UpdateMe(c *gin.Context) {
	adminID, err := strconv.ParseUint(c.GetString("userID"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token.", nil)
		return
	}
	var payload dto.UpdateAdminDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid payload.", nil)
		return
	}
	if err := h.v.Struct(payload); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_FAILED", "Validation failed.", err)
		return
	}
	admin, err := h.svc.UpdateProfile(c.Request.Context(), uint(adminID), payload)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "ADMIN_UPDATE_FAILED", "Unable to update admin.", nil)
		return
	}
	response.OK(c, admin)
}

// ListAdmins returns all admin accounts.
func (h *AdminHandler) ListAdmins(c *gin.Context) {
	admins, err := h.svc.ListAdmins(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "ADMIN_LIST_FAILED", "Unable to load admins.", nil)
		return
	}
	response.OK(c, admins)
}

// CreateAdmin creates a new admin account from dashboard.
func (h *AdminHandler) CreateAdmin(c *gin.Context) {
	var payload dto.RegisterAdminDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid payload.", nil)
		return
	}
	if err := h.v.Struct(payload); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_FAILED", "Validation failed.", err)
		return
	}

	admin, err := h.svc.CreateAdmin(c.Request.Context(), payload)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			response.Error(c, http.StatusConflict, "EMAIL_EXISTS", "Admin email already registered.", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "ADMIN_CREATE_FAILED", "Unable to create admin.", nil)
		return
	}
	response.Created(c, admin)
}
