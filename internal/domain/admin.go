package domain

import "gorm.io/gorm"

// Admin represents an administrator account.
type Admin struct {
	gorm.Model
	Name         string `gorm:"not null"`
	Email        string `gorm:"uniqueIndex;not null;index"`
	PasswordHash string `gorm:"not null" json:"-"`
	PhotoURL     string
}
