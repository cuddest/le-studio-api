package domain

import "gorm.io/gorm"

// TrainingType represents a class discipline category.
type TrainingType struct {
	gorm.Model
	Code        string `gorm:"uniqueIndex;not null"`
	Name        string `gorm:"not null"`
	Description string
	Color       string
	IsActive    bool `gorm:"default:true"`
	ParentID    *uint
	Parent      *TrainingType `gorm:"foreignKey:ParentID"`
}
