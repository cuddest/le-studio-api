package dto

// CreateCoachDTO defines coach create/update payload.
type CreateCoachDTO struct {
	FirstName     string `json:"first_name" validate:"required,min=2,max=50"`
	LastName      string `json:"last_name" validate:"required,min=2,max=50"`
	Bio           string `json:"bio"`
	PhotoURL      string `json:"photo_url"`
	PhotoPublicID string `json:"photo_public_id"`
	Specialties   string `json:"specialties"`
	IsActive      *bool  `json:"is_active"`
}
