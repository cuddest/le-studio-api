package dto

// UpdateMeDTO defines user profile update payload.
type UpdateMeDTO struct {
	FirstName string `json:"first_name" validate:"omitempty,min=2,max=50"`
	LastName  string `json:"last_name" validate:"omitempty,min=2,max=50"`
	Phone     string `json:"phone"`
	Gender    string `json:"gender" validate:"omitempty,oneof=male female"`
	PhotoURL  string `json:"photo_url"`
}

// ChangePasswordDTO defines own password-change payload.
type ChangePasswordDTO struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

// PromoteGuestDTO defines guest promotion payload.
type PromoteGuestDTO struct {
	Password string `json:"password" validate:"required,min=8"`
}

// AdminUpdateUserDTO defines admin user update payload.
type AdminUpdateUserDTO struct {
	FirstName string `json:"first_name" validate:"omitempty,min=2,max=50"`
	LastName  string `json:"last_name" validate:"omitempty,min=2,max=50"`
	Phone     string `json:"phone"`
	Gender    string `json:"gender" validate:"omitempty,oneof=male female"`
	PhotoURL  string `json:"photo_url"`
	IsActive  *bool  `json:"is_active"`
}
