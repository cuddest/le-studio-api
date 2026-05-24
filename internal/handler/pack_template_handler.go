package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"le-studio-api/internal/dto"
	"le-studio-api/internal/service"
	"le-studio-api/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// PackTemplateHandler serves pack-template endpoints.
type PackTemplateHandler struct {
	svc *service.PackTemplateService
	v   *validator.Validate
}

// NewPackTemplateHandler creates PackTemplateHandler.
func NewPackTemplateHandler(svc *service.PackTemplateService, v *validator.Validate) *PackTemplateHandler {
	return &PackTemplateHandler{svc: svc, v: v}
}

func (h *PackTemplateHandler) List(c *gin.Context) {
	includeInactive := strings.EqualFold(c.Query("include_inactive"), "true")
	templates, err := h.svc.List(c.Request.Context(), includeInactive)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "TEMPLATE_LIST_FAILED", "Unable to load templates.", nil)
		return
	}
	response.OK(c, templates)
}

func (h *PackTemplateHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid id.", nil)
		return
	}
	tpl, err := h.svc.Get(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Template not found.", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "TEMPLATE_FETCH_FAILED", "Unable to load template.", nil)
		return
	}
	response.OK(c, tpl)
}

func (h *PackTemplateHandler) AdminCreate(c *gin.Context) {
	var payload dto.CreatePackTemplateDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid payload.", nil)
		return
	}
	if err := h.v.Struct(payload); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_FAILED", "Validation failed.", err)
		return
	}
	tpl, err := h.svc.Create(c.Request.Context(), payload)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "TEMPLATE_CREATE_FAILED", "Unable to create template.", nil)
		return
	}
	response.Created(c, tpl)
}

func (h *PackTemplateHandler) AdminUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid id.", nil)
		return
	}
	var payload dto.CreatePackTemplateDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid payload.", nil)
		return
	}
	if err := h.v.Struct(payload); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_FAILED", "Validation failed.", err)
		return
	}
	tpl, err := h.svc.Update(c.Request.Context(), uint(id), payload)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "TEMPLATE_UPDATE_FAILED", "Unable to update template.", nil)
		return
	}
	response.OK(c, tpl)
}

func (h *PackTemplateHandler) AdminDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid id.", nil)
		return
	}
	if err := h.svc.Delete(c.Request.Context(), uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "TEMPLATE_DELETE_FAILED", "Unable to delete template.", nil)
		return
	}
	response.OK(c, gin.H{"message": "deleted"})
}

func (h *PackTemplateHandler) AdminReorder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid id.", nil)
		return
	}
	var payload dto.ReorderPackTemplateDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid payload.", nil)
		return
	}
	if err := h.v.Struct(payload); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_FAILED", "Validation failed.", err)
		return
	}
	tpl, err := h.svc.Reorder(c.Request.Context(), uint(id), payload)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "TEMPLATE_UPDATE_FAILED", "Unable to reorder template.", nil)
		return
	}
	response.OK(c, tpl)
}
