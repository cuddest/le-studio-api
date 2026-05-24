package domain

import (
	"time"

	"gorm.io/gorm"
)

// UserPack represents a purchased pack instance.
type UserPack struct {
	gorm.Model
	UserID         uint         `gorm:"not null;index:idx_user_pack_user_status,priority:1"`
	User           User         `gorm:"foreignKey:UserID"`
	PackTemplateID uint         `gorm:"not null"`
	PackTemplate   PackTemplate `gorm:"foreignKey:PackTemplateID"`
	TotalSessions  int          `gorm:"not null"`
	UsedSessions   int          `gorm:"default:0"`
	IsPaid         bool         `gorm:"default:false;index"`
	PaidAt         *time.Time
	Notes          string
	Status         string `gorm:"default:'active';index:idx_user_pack_user_status,priority:2"`
	ExpiresAt      *time.Time
}

// RemainingSessions returns remaining sessions.
func (p *UserPack) RemainingSessions() int { return p.TotalSessions - p.UsedSessions }

// IsExhausted returns true if pack has no sessions left.
func (p *UserPack) IsExhausted() bool { return p.RemainingSessions() <= 0 }
