package handler

import (
	"errors"
	"net/http"

	"le-studio-api/internal/dto"
	"le-studio-api/internal/service"
	"le-studio-api/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// AuthHandler serves auth endpoints.
type AuthHandler struct {
	svc *service.AuthService
	v   *validator.Validate
}

// NewAuthHandler creates AuthHandler.
func NewAuthHandler(svc *service.AuthService, v *validator.Validate) *AuthHandler {
	return &AuthHandler{svc: svc, v: v}
}

// Register handles registration.
func (h *AuthHandler) Register(c *gin.Context) {
	var payload dto.RegisterDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid payload.", nil)
		return
	}
	if err := h.v.Struct(payload); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_FAILED", "Validation failed.", err)
		return
	}
	access, refresh, user, err := h.svc.Register(c.Request.Context(), payload, c.Request.UserAgent(), c.ClientIP())
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			response.Error(c, http.StatusConflict, "EMAIL_EXISTS", "Email already registered.", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "REGISTER_FAILED", "Unable to register.", nil)
		return
	}
	response.Created(c, gin.H{"access_token": access, "refresh_token": refresh, "user": user})
}

// Login handles login.
func (h *AuthHandler) Login(c *gin.Context) {
	var payload dto.LoginDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid payload.", nil)
		return
	}
	if err := h.v.Struct(payload); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_FAILED", "Validation failed.", err)
		return
	}
	access, refresh, user, err := h.svc.Login(c.Request.Context(), payload, c.Request.UserAgent(), c.ClientIP())
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password.", nil)
		return
	}
	response.OK(c, gin.H{"access_token": access, "refresh_token": refresh, "user": user})
}

// Refresh handles refresh.
func (h *AuthHandler) Refresh(c *gin.Context) {
	var payload dto.RefreshDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid payload.", nil)
		return
	}
	if err := h.v.Struct(payload); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_FAILED", "Validation failed.", err)
		return
	}
	access, refresh, _, _, err := h.svc.Refresh(c.Request.Context(), payload.RefreshToken, c.Request.UserAgent(), c.ClientIP())
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "INVALID_REFRESH", "Invalid refresh token.", nil)
		return
	}
	response.OK(c, gin.H{"access_token": access, "refresh_token": refresh})
}

// Logout handles logout.
func (h *AuthHandler) Logout(c *gin.Context) {
	var payload dto.RefreshDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid payload.", nil)
		return
	}
	if err := h.v.Struct(payload); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_FAILED", "Validation failed.", err)
		return
	}
	if err := h.svc.Logout(c.Request.Context(), payload.RefreshToken); err != nil {
		response.Error(c, http.StatusUnauthorized, "INVALID_REFRESH", "Invalid refresh token.", nil)
		return
	}
	response.OK(c, gin.H{"message": "logged out"})
}

// Guest handles guest auth.
func (h *AuthHandler) Guest(c *gin.Context) {
	var payload dto.GuestPurchaseDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid payload.", nil)
		return
	}
	if err := h.v.Struct(payload); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_FAILED", "Validation failed.", err)
		return
	}
	access, refresh, user, pack, err := h.svc.GuestPurchase(c.Request.Context(), payload, c.Request.UserAgent(), c.ClientIP())
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			response.Error(c, http.StatusConflict, "EMAIL_EXISTS", "Email already registered.", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "GUEST_PURCHASE_FAILED", "Unable to create guest purchase.", nil)
		return
	}
	response.Created(c, gin.H{"access_token": access, "refresh_token": refresh, "user": user, "user_pack": pack})
}

// AdminLogin handles admin login.
func (h *AuthHandler) AdminLogin(c *gin.Context) {
	var payload dto.LoginDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid payload.", nil)
		return
	}
	if err := h.v.Struct(payload); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_FAILED", "Validation failed.", err)
		return
	}
	access, refresh, admin, err := h.svc.AdminLogin(c.Request.Context(), payload, c.Request.UserAgent(), c.ClientIP())
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password.", nil)
		return
	}
	response.OK(c, gin.H{"access_token": access, "refresh_token": refresh, "admin": admin})
}
