package domain

import "gorm.io/gorm"

// Coach represents an instructor profile.
type Coach struct {
	gorm.Model
	FirstName   string `gorm:"not null"`
	LastName    string `gorm:"not null"`
	Bio         string
	PhotoURL    string
	Specialties string
	IsActive    bool `gorm:"default:true"`
}
