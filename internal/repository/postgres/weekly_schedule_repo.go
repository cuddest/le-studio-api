package postgres

import (
	"context"

	"le-studio-api/internal/domain"

	"gorm.io/gorm"
)

// WeeklyScheduleRepo is schedule repository.
type WeeklyScheduleRepo struct{ db *gorm.DB }

// NewWeeklyScheduleRepo creates schedule repository.
func NewWeeklyScheduleRepo(db *gorm.DB) *WeeklyScheduleRepo { return &WeeklyScheduleRepo{db: db} }

// Create inserts schedule.
func (r *WeeklyScheduleRepo) Create(ctx context.Context, schedule *domain.WeeklySchedule) error {
	return r.db.WithContext(ctx).Create(schedule).Error
}

// GetByID returns schedule by id.
func (r *WeeklyScheduleRepo) GetByID(ctx context.Context, id uint) (*domain.WeeklySchedule, error) {
	var schedule domain.WeeklySchedule
	if err := r.db.WithContext(ctx).Preload("Slots").First(&schedule, id).Error; err != nil {
		return nil, err
	}
	return &schedule, nil
}

// List returns schedules.
func (r *WeeklyScheduleRepo) List(ctx context.Context, includeUnpublished bool) ([]domain.WeeklySchedule, error) {
	query := r.db.WithContext(ctx).Order("week_start desc")
	if !includeUnpublished {
		query = query.Where("is_published = ?", true)
	}
	var schedules []domain.WeeklySchedule
	if err := query.Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

// Update saves schedule changes.
func (r *WeeklyScheduleRepo) Update(ctx context.Context, schedule *domain.WeeklySchedule) error {
	return r.db.WithContext(ctx).Save(schedule).Error
}

// Delete removes schedule.
func (r *WeeklyScheduleRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.WeeklySchedule{}, id).Error
}
