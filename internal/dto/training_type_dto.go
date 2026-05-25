package dto

// TrainingTypeDTO defines training type create/update payload.
type TrainingTypeDTO struct {
	Code        string `json:"code" validate:"required,min=1,max=10"`
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Description string `json:"description"`
	Color       string `json:"color"`
	IsActive    *bool  `json:"is_active"`
	ParentID    *uint  `json:"parent_id"`
}
