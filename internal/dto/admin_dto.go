package dto

// UpdateAdminDTO defines admin profile update payload.
type UpdateAdminDTO struct {
	Name     string `json:"name" validate:"omitempty,min=2,max=100"`
	PhotoURL string `json:"photo_url"`
}