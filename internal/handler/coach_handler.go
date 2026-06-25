package handler

import (
	"errors"
	"mime/multipart"
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

type coachRequest struct {
	Payload dto.CreateCoachDTO
	Photo   *multipart.FileHeader
}

func parseCoachRequest(c *gin.Context) (coachRequest, error) {
	if strings.HasPrefix(c.ContentType(), "multipart/form-data") {
		isActive := c.PostForm("is_active")
		var activePtr *bool
		if isActive != "" {
			parsed := strings.EqualFold(isActive, "true") || isActive == "1" || strings.EqualFold(isActive, "on")
			activePtr = &parsed
		}

		payload := dto.CreateCoachDTO{
			FirstName:     c.PostForm("first_name"),
			LastName:      c.PostForm("last_name"),
			Bio:           c.PostForm("bio"),
			PhotoURL:      c.PostForm("photo_url"),
			PhotoPublicID: c.PostForm("photo_public_id"),
			Specialties:   c.PostForm("specialties"),
			IsActive:      activePtr,
		}

		var file *multipart.FileHeader
		if uploaded, err := c.FormFile("photo"); err == nil {
			file = uploaded
		}

		return coachRequest{Payload: payload, Photo: file}, nil
	}

	var payload dto.CreateCoachDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		return coachRequest{}, err
	}
	return coachRequest{Payload: payload}, nil
}

func (h *CoachHandler) List(c *gin.Context) {
	includeInactive := strings.EqualFold(c.Query("include_inactive"), "true")
	includeDeleted := strings.EqualFold(c.Query("include_deleted"), "true")
	coaches, err := h.svc.List(c.Request.Context(), includeInactive, includeDeleted)
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
	request, err := parseCoachRequest(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid payload.", nil)
		return
	}
	if err := h.v.Struct(request.Payload); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_FAILED", "Validation failed.", formatValidationErrors(err))
		return
	}
	if request.Photo != nil {
		if h.cld == nil {
			response.Error(c, http.StatusServiceUnavailable, "CLOUDINARY_ERROR", "Photo uploads are unavailable.", nil)
			return
		}
		uploaded, err := h.cld.UploadFile(c.Request.Context(), request.Photo)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "UPLOAD_FAILED", "Unable to upload coach photo.", nil)
			return
		}
		request.Payload.PhotoURL = uploaded.URL
		request.Payload.PhotoPublicID = uploaded.PublicID
	}
	coach, err := h.svc.Create(c.Request.Context(), request.Payload)
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
	existing, err := h.svc.Get(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "NOT_FOUND", "Coach not found.", nil)
		return
	}

	request, err := parseCoachRequest(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid payload.", nil)
		return
	}
	if err := h.v.Struct(request.Payload); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_FAILED", "Validation failed.", formatValidationErrors(err))
		return
	}
	if request.Photo != nil {
		if h.cld == nil {
			response.Error(c, http.StatusServiceUnavailable, "CLOUDINARY_ERROR", "Photo uploads are unavailable.", nil)
			return
		}
		uploaded, err := h.cld.UploadFile(c.Request.Context(), request.Photo)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "UPLOAD_FAILED", "Unable to upload coach photo.", nil)
			return
		}
		request.Payload.PhotoURL = uploaded.URL
		request.Payload.PhotoPublicID = uploaded.PublicID
	}

	coach, err := h.svc.Update(c.Request.Context(), uint(id), request.Payload)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "COACH_UPDATE_FAILED", "Unable to update coach.", nil)
		return
	}
	if request.Photo != nil && h.cld != nil && existing.PhotoPublicID != "" && existing.PhotoPublicID != coach.PhotoPublicID {
		_ = h.cld.DeleteByPublicID(c.Request.Context(), existing.PhotoPublicID)
	}
	response.OK(c, coach)
}

func (h *CoachHandler) AdminDelete(c *gin.Context) {
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
	if h.cld != nil && coach.PhotoPublicID != "" {
		if err := h.cld.DeleteByPublicID(c.Request.Context(), coach.PhotoPublicID); err != nil {
			response.Error(c, http.StatusInternalServerError, "PHOTO_DELETE_FAILED", "Unable to delete coach photo.", nil)
			return
		}
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
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid id.", nil)
		return
	}

	existing, err := h.svc.Get(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Coach not found.", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "COACH_FETCH_FAILED", "Unable to load coach.", nil)
		return
	}

	request, err := parseCoachRequest(c)
	if err != nil || request.Photo == nil {
		response.Error(c, http.StatusBadRequest, "INVALID_FILE", "No photo provided.", nil)
		return
	}
	if h.cld == nil {
		response.Error(c, http.StatusServiceUnavailable, "CLOUDINARY_ERROR", "Photo uploads are unavailable.", nil)
		return
	}

	uploaded, err := h.cld.UploadFile(c.Request.Context(), request.Photo)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "UPLOAD_FAILED", "Unable to upload coach photo.", nil)
		return
	}

	request.Payload.FirstName = existing.FirstName
	request.Payload.LastName = existing.LastName
	request.Payload.Bio = existing.Bio
	request.Payload.Specialties = existing.Specialties
	request.Payload.IsActive = &existing.IsActive
	request.Payload.PhotoURL = uploaded.URL
	request.Payload.PhotoPublicID = uploaded.PublicID

	coach, err := h.svc.Update(c.Request.Context(), uint(id), request.Payload)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "COACH_UPDATE_FAILED", "Unable to update coach.", nil)
		return
	}
	if existing.PhotoPublicID != "" && existing.PhotoPublicID != coach.PhotoPublicID {
		_ = h.cld.DeleteByPublicID(c.Request.Context(), existing.PhotoPublicID)
	}
	response.OK(c, gin.H{"photo_url": coach.PhotoURL, "photo_public_id": coach.PhotoPublicID, "coach": coach})
}
