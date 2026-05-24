package dto

// CreatePackTemplateDTO defines pack template create/update payload.
type CreatePackTemplateDTO struct {
	Name             string  `json:"name" validate:"required,min=2,max=100"`
	NumberOfSessions int     `json:"number_of_sessions" validate:"required,min=1,max=500"`
	Price            float64 `json:"price" validate:"required,gte=0"`
	TrainingTypeID   uint    `json:"training_type_id" validate:"required"`
	IsActive         *bool   `json:"is_active"`
	DisplayOrder     int     `json:"display_order"`
}

// ReorderPackTemplateDTO defines display order update payload.
type ReorderPackTemplateDTO struct {
	DisplayOrder int `json:"display_order" validate:"required"`
}
