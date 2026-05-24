package dto

import "time"

// CreateUserPackDTO defines payload for purchasing a pack.
type CreateUserPackDTO struct {
	PackTemplateID uint `json:"pack_template_id" validate:"required"`
	UserID         uint `json:"user_id"`
}

// UpdateUserPackDTO defines patch payload for admin updates.
type UpdateUserPackDTO struct {
	Notes     string     `json:"notes"`
	ExpiresAt *time.Time `json:"expires_at"`
}

// AdjustUserPackDTO defines manual used-sessions adjustment payload.
type AdjustUserPackDTO struct {
	UsedSessions int `json:"used_sessions" validate:"required,min=0"`
}
