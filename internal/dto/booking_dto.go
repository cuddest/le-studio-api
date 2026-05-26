package dto

// CreateBookingDTO defines booking create payload.
type CreateBookingDTO struct {
	SlotID     uint `json:"slot_id" validate:"required"`
	UserPackID uint `json:"user_pack_id" validate:"required"`
}

// AdminCreateBookingDTO defines admin booking create payload.
type AdminCreateBookingDTO struct {
	UserID     uint `json:"user_id" validate:"required"`
	SlotID     uint `json:"slot_id" validate:"required"`
	UserPackID uint `json:"user_pack_id" validate:"required"`
}
