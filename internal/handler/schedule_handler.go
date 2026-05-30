package handler

import (
	"errors"
	"log"
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

// ScheduleHandler serves schedule endpoints.
type ScheduleHandler struct {
	svc *service.ScheduleService
	v   *validator.Validate
}

// NewScheduleHandler creates ScheduleHandler.
func NewScheduleHandler(svc *service.ScheduleService, v *validator.Validate) *ScheduleHandler {
	return &ScheduleHandler{svc: svc, v: v}
}

func (h *ScheduleHandler) List(c *gin.Context) {
	includeUnpublished := strings.EqualFold(c.Query("include_unpublished"), "true")
	schedules, err := h.svc.List(c.Request.Context(), includeUnpublished)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "SCHEDULE_LIST_FAILED", "Unable to load schedules.", nil)
		return
	}
	response.OK(c, schedules)
}

func (h *ScheduleHandler) AdminList(c *gin.Context) {
	schedules, err := h.svc.List(c.Request.Context(), true)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "SCHEDULE_LIST_FAILED", "Unable to load schedules.", nil)
		return
	}
	response.OK(c, schedules)
}

func (h *ScheduleHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid id.", nil)
		return
	}
	schedule, err := h.svc.Get(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "Schedule not found.", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "SCHEDULE_FETCH_FAILED", "Unable to load schedule.", nil)
		return
	}
	response.OK(c, schedule)
}

func (h *ScheduleHandler) ListSlots(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid id.", nil)
		return
	}
	includeCancelled := strings.EqualFold(c.Query("include_cancelled"), "true")
	slots, err := h.svc.ListSlots(c.Request.Context(), uint(id), includeCancelled)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "SLOT_LIST_FAILED", "Unable to load slots.", nil)
		return
	}
	response.OK(c, slots)
}

func (h *ScheduleHandler) AdminCreate(c *gin.Context) {
	var payload dto.CreateScheduleDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid payload.", nil)
		return
	}
	if err := h.v.Struct(payload); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_FAILED", "Validation failed.", err)
		return
	}
	schedule, err := h.svc.AdminCreate(c.Request.Context(), payload)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "SCHEDULE_CREATE_FAILED", "Unable to create schedule.", nil)
		return
	}
	response.Created(c, schedule)
}

func (h *ScheduleHandler) AdminPublish(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid id.", nil)
		return
	}
	schedule, err := h.svc.AdminPublish(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "SCHEDULE_PUBLISH_FAILED", "Unable to publish schedule.", nil)
		return
	}
	response.OK(c, schedule)
}

func (h *ScheduleHandler) AdminUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid id.", nil)
		return
	}
	var payload dto.CreateScheduleDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid payload.", nil)
		return
	}
	if err := h.v.Struct(payload); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_FAILED", "Validation failed.", err)
		return
	}
	schedule, err := h.svc.AdminUpdate(c.Request.Context(), uint(id), payload)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "SCHEDULE_UPDATE_FAILED", "Unable to update schedule.", nil)
		return
	}
	response.OK(c, schedule)
}

func (h *ScheduleHandler) AdminDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid id.", nil)
		return
	}
	if err := h.svc.AdminDelete(c.Request.Context(), uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "SCHEDULE_DELETE_FAILED", "Unable to delete schedule.", nil)
		return
	}
	response.OK(c, gin.H{"message": "deleted"})
}

func (h *ScheduleHandler) AdminCreateSlot(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid id.", nil)
		return
	}
	var payload dto.CreateSlotDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid payload.", nil)
		return
	}
	if err := h.v.Struct(payload); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_FAILED", "Validation failed.", err)
		return
	}
	slot, err := h.svc.CreateSlot(c.Request.Context(), uint(id), payload)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "SLOT_CREATE_FAILED", "Unable to create slot.", nil)
		return
	}
	response.Created(c, slot)
}

func (h *ScheduleHandler) AdminUpdateSlot(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid id.", nil)
		return
	}
	var payload dto.CreateSlotDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid payload.", nil)
		return
	}
	if err := h.v.Struct(payload); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_FAILED", "Validation failed.", err)
		return
	}
	// Log received payload to help debug slot updates
	log.Printf("AdminUpdateSlot: received payload for slot %s: training_type_id=%d, coach_id=%d, name=%s", c.Param("id"), payload.TrainingTypeID, payload.CoachID, payload.Name)
	slot, err := h.svc.UpdateSlot(c.Request.Context(), uint(id), payload)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "SLOT_UPDATE_FAILED", "Unable to update slot.", nil)
		return
	}
	log.Printf("AdminUpdateSlot: update result for slot %d: training_type_id=%d", slot.ID, slot.TrainingTypeID)
	response.OK(c, slot)
}

func (h *ScheduleHandler) AdminDeleteSlot(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid id.", nil)
		return
	}
	if err := h.svc.DeleteSlot(c.Request.Context(), uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "SLOT_DELETE_FAILED", "Unable to delete slot.", nil)
		return
	}
	response.OK(c, gin.H{"message": "deleted"})
}

func (h *ScheduleHandler) AdminCancelSlot(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid id.", nil)
		return
	}
	slot, err := h.svc.CancelSlot(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "SLOT_CANCEL_FAILED", "Unable to cancel slot.", nil)
		return
	}
	response.OK(c, slot)
}
