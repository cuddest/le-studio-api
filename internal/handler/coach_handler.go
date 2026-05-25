package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"le-studio-api/internal/dto"
	"le-studio-api/internal/service"
	"le-studio-api/pkg/cloudinary"
	"le-studio-api/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// CoachHandler serves coach endpoints.
type CoachHandler struct {
	svc *service.CoachService
	v   *validator.Validate
	cld *cloudinary.Client
}

// NewCoachHandler creates CoachHandler.
func NewCoachHandler(svc *service.CoachService, v *validator.Validate) *CoachHandler {
	return &CoachHandler{svc: svc, v: v}
}

// SetCloudinaryClient sets the cloudinary client for photo uploads.
func (h *CoachHandler) SetCloudinaryClient(cld *cloudinary.Client) {
	h.cld = cld
}

func (h *CoachHandler) List(c *gin.Context) {
	includeInactive := strings.EqualFold(c.Query("include_inactive"), "true")
	coaches, err := h.svc.List(c.Request.Context(), includeInactive)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "COACH_LIST_FAILED", "Unable to load coaches.", nil)
		return
	}
	response.OK(c, coaches)
}

func (h *CoachHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid id.", nil)
		return
	}
	coach, err := h.svc.Get(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Coach not found.", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "COACH_FETCH_FAILED", "Unable to load coach.", nil)
		return
	}
	response.OK(c, coach)
}

func (h *CoachHandler) AdminCreate(c *gin.Context) {
	var payload dto.CreateCoachDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid payload.", nil)
		return
	}
	if err := h.v.Struct(payload); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_FAILED", "Validation failed.", err)
		return
	}
	coach, err := h.svc.Create(c.Request.Context(), payload)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "COACH_CREATE_FAILED", "Unable to create coach.", nil)
		return
	}
	response.Created(c, coach)
}

func (h *CoachHandler) AdminUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid id.", nil)
		return
	}
	var payload dto.CreateCoachDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid payload.", nil)
		return
	}
	if err := h.v.Struct(payload); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_FAILED", "Validation failed.", err)
		return
	}
	coach, err := h.svc.Update(c.Request.Context(), uint(id), payload)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "COACH_UPDATE_FAILED", "Unable to update coach.", nil)
		return
	}
	response.OK(c, coach)
}

func (h *CoachHandler) AdminDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid id.", nil)
		return
	}
	if err := h.svc.Delete(c.Request.Context(), uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "COACH_DELETE_FAILED", "Unable to delete coach.", nil)
		return
	}
	response.OK(c, gin.H{"message": "deleted"})
}

func (h *CoachHandler) AdminToggleActive(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid id.", nil)
		return
	}
	coach, err := h.svc.ToggleActive(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "COACH_UPDATE_FAILED", "Unable to update coach.", nil)
		return
	}
	response.OK(c, coach)
}

func (h *CoachHandler) AdminUploadPhoto(c *gin.Context) {
	if h.cld == nil {
		response.Error(c, http.StatusServiceUnavailable, "CLOUDINARY_ERROR", "File upload service unavailable.", nil)
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid id.", nil)
		return
	}

	file, err := c.FormFile("photo")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_FILE", "No file provided.", nil)
		return
	}

	// Upload to Cloudinary
	photoURL, err := h.cld.UploadFile(c.Request.Context(), file)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "UPLOAD_FAILED", err.Error(), nil)
		return
	}

	// Update coach with photo URL
	coach, err := h.svc.Update(c.Request.Context(), uint(id), dto.CreateCoachDTO{PhotoURL: photoURL})
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "COACH_UPDATE_FAILED", "Unable to update coach.", nil)
		return
	}

	response.OK(c, gin.H{"photo_url": photoURL, "coach": coach})
}
