package domain

import "gorm.io/gorm"

// TrainingType represents a class discipline category.
type TrainingType struct {
	gorm.Model
	Code        string         `gorm:"uniqueIndex;not null" json:"code"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `json:"description"`
	Color       string         `json:"color"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	ParentID    *uint          `json:"parent_id"`
	Parent      *TrainingType  `gorm:"foreignKey:ParentID" json:"parent"`
}
