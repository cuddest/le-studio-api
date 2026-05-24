package domain

import (
	"time"

	"gorm.io/gorm"
)

// Booking represents a user reservation.
type Booking struct {
	gorm.Model
	UserID      uint     `gorm:"not null;uniqueIndex:idx_booking_user_slot,priority:1"`
	User        User     `gorm:"foreignKey:UserID"`
	SlotID      uint     `gorm:"not null;uniqueIndex:idx_booking_user_slot,priority:2"`
	Slot        Slot     `gorm:"foreignKey:SlotID"`
	UserPackID  uint     `gorm:"not null"`
	UserPack    UserPack `gorm:"foreignKey:UserPackID"`
	Status      string   `gorm:"default:'confirmed';index"`
	CancelledAt *time.Time
}
