package handler

import (
	"errors"
	"net/http"
	"strconv"

	"le-studio-api/internal/dto"
	"le-studio-api/internal/service"
	"le-studio-api/pkg/pagination"
	"le-studio-api/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// UserPackHandler serves user-pack endpoints.
type UserPackHandler struct {
	svc *service.UserPackService
	v   *validator.Validate
}

// NewUserPackHandler creates UserPackHandler.
func NewUserPackHandler(svc *service.UserPackService, v *validator.Validate) *UserPackHandler {
	return &UserPackHandler{svc: svc, v: v}
}

func (h *UserPackHandler) Purchase(c *gin.Context) {
	userID, err := parseUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token.", nil)
		return
	}
	var payload dto.CreateUserPackDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid payload.", nil)
		return
	}
	if err := h.v.Struct(payload); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_FAILED", "Validation failed.", formatValidationErrors(err))
		return
	}
	pack, err := h.svc.Purchase(c.Request.Context(), userID, payload)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Pack template not found.", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "PURCHASE_FAILED", "Unable to purchase pack.", nil)
		return
	}
	response.Created(c, pack)
}

func (h *UserPackHandler) AdminCreate(c *gin.Context) {
	var payload dto.CreateUserPackDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid payload.", nil)
		return
	}
	if err := h.v.Struct(payload); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_FAILED", "Validation failed.", formatValidationErrors(err))
		return
	}
	pack, err := h.svc.Purchase(c.Request.Context(), payload.UserID, payload)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Pack template not found.", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "PURCHASE_FAILED", "Unable to create user pack.", nil)
		return
	}
	response.Created(c, pack)
}

func (h *UserPackHandler) ListByUser(c *gin.Context) {
	userID, err := parseUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token.", nil)
		return
	}
	packs, err := h.svc.ListByUser(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "PACK_LIST_FAILED", "Unable to load packs.", nil)
		return
	}
	response.OK(c, packs)
}

func (h *UserPackHandler) AdminList(c *gin.Context) {
	params := pagination.Parse(c)
	packs, total, err := h.svc.AdminList(c.Request.Context(), params)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "PACK_LIST_FAILED", "Unable to load packs.", nil)
		return
	}
	meta := response.Meta{Page: params.Page, Limit: params.Limit, Total: int(total), TotalPages: calcTotalPages(total, params.Limit)}
	response.Paginated(c, packs, meta)
}

func (h *UserPackHandler) AdminGet(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid id.", nil)
		return
	}
	pack, err := h.svc.AdminGet(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Pack not found.", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "PACK_FETCH_FAILED", "Unable to load pack.", nil)
		return
	}
	response.OK(c, pack)
}

func (h *UserPackHandler) AdminUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid id.", nil)
		return
	}
	var payload dto.UpdateUserPackDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid payload.", nil)
		return
	}
	pack, err := h.svc.AdminUpdate(c.Request.Context(), uint(id), payload)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "PACK_UPDATE_FAILED", "Unable to update pack.", nil)
		return
	}
	response.OK(c, pack)
}

func (h *UserPackHandler) AdminMarkPaid(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid id.", nil)
		return
	}
	pack, err := h.svc.AdminMarkPaid(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "PACK_UPDATE_FAILED", "Unable to update pack.", nil)
		return
	}
	response.OK(c, pack)
}

func (h *UserPackHandler) AdminAdjust(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid id.", nil)
		return
	}
	var payload dto.AdjustUserPackDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid payload.", nil)
		return
	}
	if err := h.v.Struct(payload); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_FAILED", "Validation failed.", formatValidationErrors(err))
		return
	}
	pack, err := h.svc.AdminAdjust(c.Request.Context(), uint(id), payload)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "PACK_UPDATE_FAILED", "Unable to update pack.", nil)
		return
	}
	response.OK(c, pack)
}

func (h *UserPackHandler) AdminDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid id.", nil)
		return
	}
	if err := h.svc.AdminDelete(c.Request.Context(), uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "PACK_DELETE_FAILED", "Unable to delete pack.", nil)
		return
	}
	response.OK(c, gin.H{"message": "deleted"})
}
