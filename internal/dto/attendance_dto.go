package dto

// MarkAttendanceDTO defines attendance create payload.
type MarkAttendanceDTO struct {
	BookingID uint   `json:"booking_id" validate:"required"`
	Notes     string `json:"notes"`
}

// UpdateAttendanceDTO defines attendance patch payload.
type UpdateAttendanceDTO struct {
	Notes string `json:"notes"`
}
