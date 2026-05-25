package dto

import "time"

// CreateScheduleDTO defines schedule create payload.
type CreateScheduleDTO struct {
	Title     string    `json:"title" validate:"required,min=2,max=100"`
	WeekStart time.Time `json:"week_start" validate:"required"`
	WeekEnd   time.Time `json:"week_end" validate:"required,gtfield=WeekStart"`
}

// CreateSlotDTO defines slot create/update payload.
type CreateSlotDTO struct {
	TrainingTypeID uint      `json:"training_type_id" validate:"required"`
	CoachID        uint      `json:"coach_id" validate:"required"`
	SlotType       string    `json:"slot_type" validate:"required,oneof=mixte women_only men_only"`
	Name           string    `json:"name" validate:"required,min=2,max=100"`
	Date           time.Time `json:"date" validate:"required"`
	StartTime      time.Time `json:"start_time" validate:"required"`
	EndTime        time.Time `json:"end_time" validate:"required,gtfield=StartTime"`
	Level          string    `json:"level" validate:"required"`
	Capacity       int       `json:"capacity" validate:"required,min=1,max=500"`
}
