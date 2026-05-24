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

// UserHandler serves user endpoints.
type UserHandler struct {
	svc *service.UserService
	v   *validator.Validate
}

// NewUserHandler creates UserHandler.
func NewUserHandler(svc *service.UserService, v *validator.Validate) *UserHandler {
	return &UserHandler{svc: svc, v: v}
}

func (h *UserHandler) GetMe(c *gin.Context) {
	userID, err := parseUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token.", nil)
		return
	}
	user, err := h.svc.GetMe(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "USER_FETCH_FAILED", "Unable to load user.", nil)
		return
	}
	response.OK(c, user)
}

func (h *UserHandler) PatchMe(c *gin.Context) {
	userID, err := parseUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token.", nil)
		return
	}
	var payload dto.UpdateMeDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid payload.", nil)
		return
	}
	if err := h.v.Struct(payload); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_FAILED", "Validation failed.", err)
		return
	}
	user, err := h.svc.UpdateMe(c.Request.Context(), userID, payload)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "USER_UPDATE_FAILED", "Unable to update profile.", nil)
		return
	}
	response.OK(c, user)
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID, err := parseUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token.", nil)
		return
	}
	var payload dto.ChangePasswordDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid payload.", nil)
		return
	}
	if err := h.v.Struct(payload); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_FAILED", "Validation failed.", err)
		return
	}
	if err := h.svc.ChangePassword(c.Request.Context(), userID, payload); err != nil {
		response.Error(c, http.StatusBadRequest, "PASSWORD_CHANGE_FAILED", err.Error(), nil)
		return
	}
	response.OK(c, gin.H{"message": "password updated"})
}

func (h *UserHandler) AdminListUsers(c *gin.Context) {
	params := pagination.Parse(c)
	users, total, err := h.svc.AdminListUsers(c.Request.Context(), params)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "USER_LIST_FAILED", "Unable to load users.", nil)
		return
	}
	meta := response.Meta{Page: params.Page, Limit: params.Limit, Total: int(total), TotalPages: calcTotalPages(total, params.Limit)}
	response.Paginated(c, users, meta)
}

func (h *UserHandler) AdminCreateUser(c *gin.Context) {
	var payload dto.RegisterDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid payload.", nil)
		return
	}
	if err := h.v.Struct(payload); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_FAILED", "Validation failed.", err)
		return
	}
	user, err := h.svc.AdminCreateUser(c.Request.Context(), payload)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			response.Error(c, http.StatusConflict, "EMAIL_EXISTS", "Email already registered.", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "USER_CREATE_FAILED", "Unable to create user.", nil)
		return
	}
	response.Created(c, user)
}

func (h *UserHandler) AdminDeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid id.", nil)
		return
	}
	if err := h.svc.AdminDeleteUser(c.Request.Context(), uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "USER_DELETE_FAILED", "Unable to delete user.", nil)
		return
	}
	response.OK(c, gin.H{"message": "deleted"})
}

func (h *UserHandler) AdminToggleActive(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid id.", nil)
		return
	}
	user, err := h.svc.AdminToggleActive(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "USER_UPDATE_FAILED", "Unable to toggle user.", nil)
		return
	}
	response.OK(c, user)
}

func (h *UserHandler) AdminPromoteGuest(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid id.", nil)
		return
	}
	var payload dto.PromoteGuestDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid payload.", nil)
		return
	}
	if err := h.v.Struct(payload); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_FAILED", "Validation failed.", err)
		return
	}
	user, err := h.svc.AdminPromoteGuest(c.Request.Context(), uint(id), payload)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "USER_UPDATE_FAILED", "Unable to promote guest.", nil)
		return
	}
	response.OK(c, user)
}

func (h *UserHandler) AdminGetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid id.", nil)
		return
	}
	user, err := h.svc.AdminGetUser(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "User not found.", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "USER_FETCH_FAILED", "Unable to load user.", nil)
		return
	}
	response.OK(c, user)
}

func (h *UserHandler) AdminUpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid id.", nil)
		return
	}
	var payload dto.AdminUpdateUserDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid payload.", nil)
		return
	}
	if err := h.v.Struct(payload); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_FAILED", "Validation failed.", err)
		return
	}
	user, err := h.svc.AdminUpdateUser(c.Request.Context(), uint(id), payload)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "USER_UPDATE_FAILED", "Unable to update user.", nil)
		return
	}
	response.OK(c, user)
}

func parseUserID(c *gin.Context) (uint, error) {
	idStr := c.GetString("userID")
	if idStr == "" {
		return 0, errors.New("missing user id")
	}
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

func calcTotalPages(total int64, limit int) int {
	if limit <= 0 {
		return 0
	}
	pages := int(total) / limit
	if int(total)%limit != 0 {
		pages += 1
	}
	return pages
}
