package domain

import "gorm.io/gorm"

// User represents a studio customer account.
type User struct {
	gorm.Model
	FirstName    string `gorm:"not null"`
	LastName     string `gorm:"not null"`
	Email        string `gorm:"uniqueIndex;not null;index"`
	Phone        string
	PasswordHash string `json:"-"`
	PhotoURL     string
	IsGuest      bool `gorm:"default:false;index"`
	IsActive     bool `gorm:"default:true;index"`

	UserPacks []UserPack
	Bookings  []Booking
}
