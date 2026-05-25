package postgres

import (
	"context"

	"le-studio-api/internal/domain"
	"le-studio-api/pkg/pagination"

	"gorm.io/gorm"
)

// BookingRepo is booking repository.
type BookingRepo struct{ db *gorm.DB }

// NewBookingRepo creates booking repository.
func NewBookingRepo(db *gorm.DB) *BookingRepo { return &BookingRepo{db: db} }

// Create inserts booking.
func (r *BookingRepo) Create(ctx context.Context, booking *domain.Booking) error {
	return r.db.WithContext(ctx).Create(booking).Error
}

// GetByID returns booking by id.
func (r *BookingRepo) GetByID(ctx context.Context, id uint) (*domain.Booking, error) {
	var booking domain.Booking
	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Slot").
		Preload("Slot.Coach").
		Preload("Slot.TrainingType").
		Preload("UserPack").
		Preload("UserPack.PackTemplate").
		Preload("UserPack.PackTemplate.TrainingTypes").
		First(&booking, id).Error; err != nil {
		return nil, err
	}
	return &booking, nil
}

// ListByUser returns bookings for user.
func (r *BookingRepo) ListByUser(ctx context.Context, userID uint) ([]domain.Booking, error) {
	var bookings []domain.Booking
	if err := r.db.WithContext(ctx).
		Preload("Slot").
		Preload("Slot.Coach").
		Preload("Slot.TrainingType").
		Preload("UserPack").
		Preload("UserPack.PackTemplate").
		Preload("UserPack.PackTemplate.TrainingTypes").
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

// List returns paginated bookings.
func (r *BookingRepo) List(ctx context.Context, params pagination.Params) ([]domain.Booking, int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&domain.Booking{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var bookings []domain.Booking
	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Slot").
		Preload("Slot.Coach").
		Preload("Slot.TrainingType").
		Preload("UserPack").
		Preload("UserPack.PackTemplate").
		Preload("UserPack.PackTemplate.TrainingTypes").
		Order("created_at desc").
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&bookings).Error; err != nil {
		return nil, 0, err
	}
	return bookings, total, nil
}

// Update saves booking changes.
func (r *BookingRepo) Update(ctx context.Context, booking *domain.Booking) error {
	return r.db.WithContext(ctx).Save(booking).Error
}
