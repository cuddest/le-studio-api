package domain

import "gorm.io/gorm"

// PackTemplate represents an admin-configured purchasable pack.
type PackTemplate struct {
	gorm.Model
	Name             string         `gorm:"not null"`
	NumberOfSessions int            `gorm:"not null"`
	Price            float64        `gorm:"not null"`
	TrainingTypeID   uint           `gorm:"not null"`
	TrainingType     TrainingType   `gorm:"foreignKey:TrainingTypeID"`
	TrainingTypes    []TrainingType `gorm:"many2many:pack_template_training_types;"`
	IsActive         bool           `gorm:"default:true"`
	DisplayOrder     int
}
