package domain

import "gorm.io/gorm"

// Attendance stores attendance marking for bookings.
type Attendance struct {
	gorm.Model
	BookingID  uint    `gorm:"uniqueIndex;not null"`
	Booking    Booking `gorm:"foreignKey:BookingID"`
	MarkedByID uint
	Notes      string
}
