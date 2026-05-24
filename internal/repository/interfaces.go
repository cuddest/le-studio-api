package repository

import (
	"context"

	"le-studio-api/internal/domain"
	"le-studio-api/pkg/pagination"
)

// AdminRepository defines admin data access.
type AdminRepository interface {
	Create(ctx context.Context, admin *domain.Admin) error
	FindByEmail(ctx context.Context, email string) (*domain.Admin, error)
	GetByID(ctx context.Context, id uint) (*domain.Admin, error)
	List(ctx context.Context) ([]domain.Admin, error)
}

// UserRepository defines user data access.
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id uint) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	List(ctx context.Context, params pagination.Params) ([]domain.User, int64, error)
	Delete(ctx context.Context, id uint) error
}

// CoachRepository defines coach data access.
type CoachRepository interface {
	Create(ctx context.Context, coach *domain.Coach) error
	GetByID(ctx context.Context, id uint) (*domain.Coach, error)
	List(ctx context.Context, includeInactive bool) ([]domain.Coach, error)
	Update(ctx context.Context, coach *domain.Coach) error
	Delete(ctx context.Context, id uint) error
}

// TrainingTypeRepository defines training type access.
type TrainingTypeRepository interface {
	Create(ctx context.Context, training *domain.TrainingType) error
	GetByID(ctx context.Context, id uint) (*domain.TrainingType, error)
	List(ctx context.Context, includeInactive bool) ([]domain.TrainingType, error)
	Update(ctx context.Context, training *domain.TrainingType) error
	Delete(ctx context.Context, id uint) error
}

// PackTemplateRepository defines pack template access.
type PackTemplateRepository interface {
	Create(ctx context.Context, tpl *domain.PackTemplate) error
	GetByID(ctx context.Context, id uint) (*domain.PackTemplate, error)
	List(ctx context.Context, includeInactive bool) ([]domain.PackTemplate, error)
	Update(ctx context.Context, tpl *domain.PackTemplate) error
	Delete(ctx context.Context, id uint) error
}

// UserPackRepository defines user pack access.
type UserPackRepository interface {
	Create(ctx context.Context, pack *domain.UserPack) error
	GetByID(ctx context.Context, id uint) (*domain.UserPack, error)
	List(ctx context.Context, params pagination.Params) ([]domain.UserPack, int64, error)
	ListByUser(ctx context.Context, userID uint) ([]domain.UserPack, error)
	Update(ctx context.Context, pack *domain.UserPack) error
	Delete(ctx context.Context, id uint) error
}

// WeeklyScheduleRepository defines schedule access.
type WeeklyScheduleRepository interface {
	Create(ctx context.Context, schedule *domain.WeeklySchedule) error
	GetByID(ctx context.Context, id uint) (*domain.WeeklySchedule, error)
	List(ctx context.Context, includeUnpublished bool) ([]domain.WeeklySchedule, error)
	Update(ctx context.Context, schedule *domain.WeeklySchedule) error
	Delete(ctx context.Context, id uint) error
}

// SlotRepository defines slot access.
type SlotRepository interface {
	Create(ctx context.Context, slot *domain.Slot) error
	GetByID(ctx context.Context, id uint) (*domain.Slot, error)
	ListBySchedule(ctx context.Context, scheduleID uint, includeCancelled bool) ([]domain.Slot, error)
	Update(ctx context.Context, slot *domain.Slot) error
	Delete(ctx context.Context, id uint) error
}

// BookingRepository defines booking access.
type BookingRepository interface {
	Create(ctx context.Context, booking *domain.Booking) error
	GetByID(ctx context.Context, id uint) (*domain.Booking, error)
	ListByUser(ctx context.Context, userID uint) ([]domain.Booking, error)
	List(ctx context.Context, params pagination.Params) ([]domain.Booking, int64, error)
	Update(ctx context.Context, booking *domain.Booking) error
}

// AttendanceRepository defines attendance access.
type AttendanceRepository interface {
	Create(ctx context.Context, attendance *domain.Attendance) error
	GetByID(ctx context.Context, id uint) (*domain.Attendance, error)
	List(ctx context.Context, params pagination.Params) ([]domain.Attendance, int64, error)
	Update(ctx context.Context, attendance *domain.Attendance) error
	Delete(ctx context.Context, id uint) error
}

// RefreshTokenRepository defines refresh token access.
type RefreshTokenRepository interface {
	Create(ctx context.Context, token *domain.RefreshToken) error
	FindByHash(ctx context.Context, hash string) (*domain.RefreshToken, error)
	Revoke(ctx context.Context, token *domain.RefreshToken) error
}

// Repositories groups repository implementations.
type Repositories struct {
	Admins       AdminRepository
	Users        UserRepository
	Coaches      CoachRepository
	Training     TrainingTypeRepository
	Templates    PackTemplateRepository
	UserPacks    UserPackRepository
	Schedules    WeeklyScheduleRepository
	Slots        SlotRepository
	Bookings     BookingRepository
	Attendance   AttendanceRepository
	RefreshToken RefreshTokenRepository
}
