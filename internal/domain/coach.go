package domain

import "gorm.io/gorm"

// Coach represents an instructor profile.
type Coach struct {
	gorm.Model
	FirstName     string `gorm:"not null" json:"first_name"`
	LastName      string `gorm:"not null" json:"last_name"`
	Bio           string `json:"bio"`
	PhotoURL      string `json:"photo_url"`
	PhotoPublicID string `json:"photo_public_id"`
	Specialties   string `json:"specialties"`
	IsActive      bool   `gorm:"default:true" json:"is_active"`
}
