package domain

import (
	"time"

	"gorm.io/gorm"
)

// WeeklySchedule represents weekly slot collection.
type WeeklySchedule struct {
	gorm.Model
	Title       string    `gorm:"not null"`
	WeekStart   time.Time `gorm:"not null"`
	WeekEnd     time.Time `gorm:"not null"`
	IsPublished bool      `gorm:"default:false"`
	PublishedAt *time.Time
	Slots       []Slot
}
