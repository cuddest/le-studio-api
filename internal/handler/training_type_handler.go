package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"le-studio-api/internal/domain"
	"le-studio-api/internal/dto"
	"le-studio-api/internal/service"
	"le-studio-api/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// TrainingTypeHandler serves training type endpoints.
type TrainingTypeHandler struct {
	svc *service.TrainingTypeService
	v   *validator.Validate
}

// NewTrainingTypeHandler creates TrainingTypeHandler.
func NewTrainingTypeHandler(svc *service.TrainingTypeService, v *validator.Validate) *TrainingTypeHandler {
	return &TrainingTypeHandler{svc: svc, v: v}
}

// List returns training types.
func (h *TrainingTypeHandler) List(c *gin.Context) {
	includeInactive := strings.EqualFold(c.Query("include_inactive"), "true")
	types, err := h.svc.List(c.Request.Context(), includeInactive)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "TRAINING_LIST_FAILED", "Unable to load training types.", nil)
		return
	}
	response.OK(c, types)
}

// AdminCreate creates a training type.
func (h *TrainingTypeHandler) AdminCreate(c *gin.Context) {
	var payload dto.TrainingTypeDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid payload.", nil)
		return
	}
	if err := h.v.Struct(payload); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_FAILED", "Validation failed.", formatValidationErrors(err))
		return
	}
	training := &domain.TrainingType{
		Code:        payload.Code,
		Name:        payload.Name,
		Description: payload.Description,
		Color:       payload.Color,
		IsActive:    true,
	}
	if payload.ParentID != nil {
		training.ParentID = payload.ParentID
	}
	if payload.IsActive != nil {
		training.IsActive = *payload.IsActive
	}
	created, err := h.svc.Create(c.Request.Context(), training)
	if err != nil {
		if errors.Is(err, service.ErrTrainingTypeInvalidParent) {
			response.Error(c, http.StatusBadRequest, "INVALID_PARENT", "Parent must be a top-level training type (one without its own parent).", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "TRAINING_CREATE_FAILED", "Unable to create training type.", nil)
		return
	}
	response.Created(c, created)
}

// AdminUpdate updates a training type.
func (h *TrainingTypeHandler) AdminUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid id.", nil)
		return
	}
	var payload dto.TrainingTypeDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid payload.", nil)
		return
	}
	if err := h.v.Struct(payload); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_FAILED", "Validation failed.", formatValidationErrors(err))
		return
	}
	training, err := h.svc.Get(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Training type not found.", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "TRAINING_FETCH_FAILED", "Unable to load training type.", nil)
		return
	}
	training.Code = payload.Code
	training.Name = payload.Name
	training.Description = payload.Description
	training.Color = payload.Color
	training.ParentID = payload.ParentID
	if payload.IsActive != nil {
		training.IsActive = *payload.IsActive
	}
	updated, err := h.svc.Update(c.Request.Context(), training)
	if err != nil {
		if errors.Is(err, service.ErrTrainingTypeInvalidParent) {
			response.Error(c, http.StatusBadRequest, "INVALID_PARENT", "Parent must be a top-level training type (one without its own parent).", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "TRAINING_UPDATE_FAILED", "Unable to update training type.", nil)
		return
	}
	response.OK(c, updated)
}

// AdminDelete deletes a training type.
func (h *TrainingTypeHandler) AdminDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid id.", nil)
		return
	}
	if err := h.svc.Delete(c.Request.Context(), uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "TRAINING_DELETE_FAILED", "Unable to delete training type.", nil)
		return
	}
	response.OK(c, gin.H{"message": "deleted"})
}
