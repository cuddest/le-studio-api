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

// AttendanceHandler serves attendance endpoints.
type AttendanceHandler struct {
	svc *service.AttendanceService
	v   *validator.Validate
}

// NewAttendanceHandler creates AttendanceHandler.
func NewAttendanceHandler(svc *service.AttendanceService, v *validator.Validate) *AttendanceHandler {
	return &AttendanceHandler{svc: svc, v: v}
}

func (h *AttendanceHandler) Mark(c *gin.Context) {
	adminID, err := parseUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token.", nil)
		return
	}
	var payload dto.MarkAttendanceDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid payload.", nil)
		return
	}
	if err := h.v.Struct(payload); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_FAILED", "Validation failed.", err)
		return
	}
	attendance, err := h.svc.Mark(c.Request.Context(), adminID, payload)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "ATTENDANCE_MARK_FAILED", "Unable to mark attendance.", nil)
		return
	}
	response.Created(c, attendance)
}

func (h *AttendanceHandler) List(c *gin.Context) {
	params := pagination.Parse(c)
	attendance, total, err := h.svc.List(c.Request.Context(), params)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "ATTENDANCE_LIST_FAILED", "Unable to load attendance.", nil)
		return
	}
	meta := response.Meta{Page: params.Page, Limit: params.Limit, Total: int(total), TotalPages: calcTotalPages(total, params.Limit)}
	response.Paginated(c, attendance, meta)
}

func (h *AttendanceHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid id.", nil)
		return
	}
	var payload dto.UpdateAttendanceDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid payload.", nil)
		return
	}
	attendance, err := h.svc.Update(c.Request.Context(), uint(id), payload)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Attendance not found.", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "ATTENDANCE_UPDATE_FAILED", "Unable to update attendance.", nil)
		return
	}
	response.OK(c, attendance)
}

func (h *AttendanceHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid id.", nil)
		return
	}
	if err := h.svc.Delete(c.Request.Context(), uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "ATTENDANCE_DELETE_FAILED", "Unable to delete attendance.", nil)
		return
	}
	response.OK(c, gin.H{"message": "deleted"})
}
