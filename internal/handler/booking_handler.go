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

// BookingHandler serves booking endpoints.
type BookingHandler struct {
	svc *service.BookingService
	v   *validator.Validate
}

// NewBookingHandler creates BookingHandler.
func NewBookingHandler(svc *service.BookingService, v *validator.Validate) *BookingHandler {
	return &BookingHandler{svc: svc, v: v}
}

func (h *BookingHandler) Create(c *gin.Context) {
	userID, err := parseUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token.", nil)
		return
	}
	var payload dto.CreateBookingDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid payload.", nil)
		return
	}
	if err := h.v.Struct(payload); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_FAILED", "Validation failed.", err)
		return
	}
	booking, err := h.svc.Create(c.Request.Context(), userID, payload)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Slot or pack not found.", nil)
			return
		}
		response.Error(c, http.StatusBadRequest, "BOOKING_FAILED", err.Error(), nil)
		return
	}
	response.Created(c, booking)
}

func (h *BookingHandler) AdminCreate(c *gin.Context) {
	var payload dto.AdminCreateBookingDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid payload.", nil)
		return
	}
	if err := h.v.Struct(payload); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_FAILED", "Validation failed.", err)
		return
	}
	booking, err := h.svc.AdminCreate(c.Request.Context(), payload)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Slot or pack not found.", nil)
			return
		}
		response.Error(c, http.StatusBadRequest, "BOOKING_FAILED", err.Error(), nil)
		return
	}
	response.Created(c, booking)
}

func (h *BookingHandler) Cancel(c *gin.Context) {
	userID, err := parseUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token.", nil)
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid id.", nil)
		return
	}
	booking, err := h.svc.Cancel(c.Request.Context(), userID, uint(id))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "BOOKING_CANCEL_FAILED", err.Error(), nil)
		return
	}
	response.OK(c, booking)
}

func (h *BookingHandler) AdminCancel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid id.", nil)
		return
	}
	booking, err := h.svc.AdminCancel(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "BOOKING_CANCEL_FAILED", err.Error(), nil)
		return
	}
	response.OK(c, booking)
}

func (h *BookingHandler) Get(c *gin.Context) {
	userID, err := parseUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token.", nil)
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid id.", nil)
		return
	}
	booking, err := h.svc.Get(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Booking not found.", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "BOOKING_FETCH_FAILED", "Unable to load booking.", nil)
		return
	}
	if booking.UserID != userID {
		response.Error(c, http.StatusForbidden, "FORBIDDEN", "Booking does not belong to user.", nil)
		return
	}
	response.OK(c, booking)
}

func (h *BookingHandler) ListByUser(c *gin.Context) {
	userID, err := parseUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token.", nil)
		return
	}
	bookings, err := h.svc.ListByUser(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "BOOKING_LIST_FAILED", "Unable to load bookings.", nil)
		return
	}
	response.OK(c, bookings)
}

func (h *BookingHandler) AdminList(c *gin.Context) {
	params := pagination.Parse(c)
	bookings, total, err := h.svc.AdminList(c.Request.Context(), params)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "BOOKING_LIST_FAILED", "Unable to load bookings.", nil)
		return
	}
	meta := response.Meta{Page: params.Page, Limit: params.Limit, Total: int(total), TotalPages: calcTotalPages(total, params.Limit)}
	response.Paginated(c, bookings, meta)
}
