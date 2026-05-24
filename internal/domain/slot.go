package domain

import (
	"time"

	"gorm.io/gorm"
)

// Slot represents one bookable session.
type Slot struct {
	gorm.Model
	WeeklyScheduleID uint           `gorm:"not null;index:idx_slot_schedule_date,priority:1"`
	WeeklySchedule   WeeklySchedule `gorm:"foreignKey:WeeklyScheduleID"`
	TrainingTypeID   uint           `gorm:"not null;index"`
	TrainingType     TrainingType   `gorm:"foreignKey:TrainingTypeID"`
	CoachID          uint           `gorm:"not null;index"`
	Coach            Coach          `gorm:"foreignKey:CoachID"`
	Name             string         `gorm:"not null"`
	DayOfWeek        int            `gorm:"not null"`
	Date             time.Time      `gorm:"not null;index:idx_slot_schedule_date,priority:2"`
	StartTime        time.Time      `gorm:"not null"`
	EndTime          time.Time      `gorm:"not null"`
	Level            string
	Capacity         int  `gorm:"not null"`
	BookedCount      int  `gorm:"default:0"`
	IsCancelled      bool `gorm:"default:false"`
	Bookings         []Booking
}
