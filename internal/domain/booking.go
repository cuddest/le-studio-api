package domain

import (
	"time"

	"gorm.io/gorm"
)

// Booking represents a user reservation.
type Booking struct {
	gorm.Model
	UserID      uint     `gorm:"not null;uniqueIndex:idx_booking_user_slot,priority:1" json:"user_id"`
	User        User     `gorm:"foreignKey:UserID" json:"-"`
	SlotID      uint     `gorm:"not null;uniqueIndex:idx_booking_user_slot,priority:2" json:"slot_id"`
	Slot        Slot     `gorm:"foreignKey:SlotID" json:"slot"`
	UserPackID  uint     `gorm:"not null" json:"user_pack_id"`
	UserPack    UserPack `gorm:"foreignKey:UserPackID" json:"user_pack"`
	Status      string   `gorm:"default:'confirmed';index" json:"status"`
	CancelledAt *time.Time `json:"cancelled_at"`
}
