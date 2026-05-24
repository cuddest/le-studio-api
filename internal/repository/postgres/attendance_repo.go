package postgres

import (
	"context"

	"le-studio-api/internal/domain"
	"le-studio-api/pkg/pagination"

	"gorm.io/gorm"
)

// AttendanceRepo is attendance repository.
type AttendanceRepo struct{ db *gorm.DB }

// NewAttendanceRepo creates attendance repository.
func NewAttendanceRepo(db *gorm.DB) *AttendanceRepo { return &AttendanceRepo{db: db} }

// Create inserts attendance record.
func (r *AttendanceRepo) Create(ctx context.Context, attendance *domain.Attendance) error {
	return r.db.WithContext(ctx).Create(attendance).Error
}

// GetByID returns attendance by id.
func (r *AttendanceRepo) GetByID(ctx context.Context, id uint) (*domain.Attendance, error) {
	var attendance domain.Attendance
	if err := r.db.WithContext(ctx).Preload("Booking").First(&attendance, id).Error; err != nil {
		return nil, err
	}
	return &attendance, nil
}

// List returns paginated attendance.
func (r *AttendanceRepo) List(ctx context.Context, params pagination.Params) ([]domain.Attendance, int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&domain.Attendance{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var attendance []domain.Attendance
	if err := r.db.WithContext(ctx).
		Preload("Booking").
		Order("created_at desc").
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&attendance).Error; err != nil {
		return nil, 0, err
	}
	return attendance, total, nil
}

// Update saves attendance changes.
func (r *AttendanceRepo) Update(ctx context.Context, attendance *domain.Attendance) error {
	return r.db.WithContext(ctx).Save(attendance).Error
}

// Delete removes attendance by id.
func (r *AttendanceRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Attendance{}, id).Error
}
