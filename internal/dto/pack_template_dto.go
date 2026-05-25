package dto

// CreatePackTemplateDTO defines pack template create/update payload.
type CreatePackTemplateDTO struct {
	Name             string  `json:"name" validate:"required,min=2,max=100"`
	NumberOfSessions int     `json:"number_of_sessions" validate:"required,min=1,max=500"`
	Price            float64 `json:"price" validate:"required,gte=0"`
	// For backward compatibility you can send `training_type_id` for a single type
	TrainingTypeID uint `json:"training_type_id"`
	// Or provide multiple type ids using `training_type_ids`
	TrainingTypeIDs []uint `json:"training_type_ids"`
	IsActive        *bool  `json:"is_active"`
	DisplayOrder    int    `json:"display_order"`
}

// ReorderPackTemplateDTO defines display order update payload.
type ReorderPackTemplateDTO struct {
	DisplayOrder int `json:"display_order" validate:"required"`
}
