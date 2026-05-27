package domain

import (
	"time"

	"gorm.io/gorm"
)

const (
	SlotTypeMixte     = "mixte"
	SlotTypeWomenOnly = "women_only"
	SlotTypeMenOnly   = "men_only"
)

// Slot represents one bookable session.
type Slot struct {
	gorm.Model
	WeeklyScheduleID uint           `gorm:"not null;index:idx_slot_schedule_date,priority:1" json:"weekly_schedule_id"`
	WeeklySchedule   WeeklySchedule `gorm:"foreignKey:WeeklyScheduleID" json:"-"`
	TrainingTypeID   uint           `gorm:"not null;index" json:"training_type_id"`
	TrainingType     TrainingType   `gorm:"foreignKey:TrainingTypeID" json:"training_type"`
	CoachID          uint           `gorm:"not null;index" json:"coach_id"`
	Coach            Coach          `gorm:"foreignKey:CoachID" json:"coach"`
	SlotType         string         `gorm:"not null;default:mixte;index" json:"slot_type"`
	Name             string         `gorm:"not null" json:"name"`
	DayOfWeek        int            `gorm:"not null" json:"day_of_week"`
	Date             time.Time      `gorm:"not null;index:idx_slot_schedule_date,priority:2" json:"date"`
	StartTime        time.Time      `gorm:"not null" json:"start_time"`
	EndTime          time.Time      `gorm:"not null" json:"end_time"`
	Level            string         `json:"level"`
	Capacity         int            `gorm:"not null" json:"capacity"`
	BookedCount      int            `gorm:"default:0" json:"booked_count"`
	IsCancelled      bool           `gorm:"default:false" json:"is_cancelled"`
	Bookings         []Booking      `json:"-"`
}
