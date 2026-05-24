package domain

import (
	"time"

	"gorm.io/gorm"
)

// RefreshToken stores hashed refresh tokens.
type RefreshToken struct {
	gorm.Model
	SubjectID uint      `gorm:"not null;index"`
	Role      string    `gorm:"not null;index"`
	TokenHash string    `gorm:"not null;uniqueIndex"`
	ExpiresAt time.Time `gorm:"not null"`
	RevokedAt *time.Time
	UserAgent string
	IPAddress string
}
