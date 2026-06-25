package postgres

import (
	"context"
	"time"

	"le-studio-api/internal/domain"

	"gorm.io/gorm"
)

// SlotRepo is slot repository.
type SlotRepo struct{ db *gorm.DB }

// NewSlotRepo creates slot repository.
func NewSlotRepo(db *gorm.DB) *SlotRepo { return &SlotRepo{db: db} }

// Create inserts slot.
func (r *SlotRepo) Create(ctx context.Context, slot *domain.Slot) error {
	return r.db.WithContext(ctx).Create(slot).Error
}

// GetByID returns slot by id.
func (r *SlotRepo) GetByID(ctx context.Context, id uint) (*domain.Slot, error) {
	var slot domain.Slot
	if err := r.db.WithContext(ctx).Preload("Coach").Preload("TrainingType").First(&slot, id).Error; err != nil {
		return nil, err
	}
	return &slot, nil
}

// ListBySchedule returns slots for schedule.
func (r *SlotRepo) ListBySchedule(ctx context.Context, scheduleID uint, includeCancelled bool) ([]domain.Slot, error) {
	query := r.db.WithContext(ctx).Preload("Coach").Preload("TrainingType").Where("weekly_schedule_id = ?", scheduleID).Order("date asc, start_time asc")
	if !includeCancelled {
		query = query.Where("is_cancelled = ?", false)
	}
	var slots []domain.Slot
	if err := query.Find(&slots).Error; err != nil {
		return nil, err
	}
	return slots, nil
}

// ExistsOverlap reports whether any non-cancelled slot in the same schedule
// overlaps the half-open interval [startTime, endTime). Pass excludeSlotID > 0
// when updating to ignore the slot being modified.
func (r *SlotRepo) ExistsOverlap(ctx context.Context, scheduleID uint, startTime, endTime time.Time, excludeSlotID uint) (bool, error) {
	var count int64
	q := r.db.WithContext(ctx).Model(&domain.Slot{}).
		Where("weekly_schedule_id = ?", scheduleID).
		Where("is_cancelled = ?", false).
		Where("start_time < ? AND end_time > ?", endTime, startTime)
	if excludeSlotID > 0 {
		q = q.Where("id <> ?", excludeSlotID)
	}
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// Update saves slot changes.
func (r *SlotRepo) Update(ctx context.Context, slot *domain.Slot) error {
	return r.db.WithContext(ctx).Save(slot).Error
}

// Delete removes slot by id.
func (r *SlotRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Slot{}, id).Error
}
