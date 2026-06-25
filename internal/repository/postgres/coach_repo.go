package postgres

import (
	"context"

	"le-studio-api/internal/domain"

	"gorm.io/gorm"
)

// CoachRepo is coach repository.
type CoachRepo struct{ db *gorm.DB }

// NewCoachRepo creates CoachRepo.
func NewCoachRepo(db *gorm.DB) *CoachRepo { return &CoachRepo{db: db} }

// Create inserts coach.
func (r *CoachRepo) Create(ctx context.Context, coach *domain.Coach) error {
	return r.db.WithContext(ctx).Create(coach).Error
}

// GetByID returns coach by id.
func (r *CoachRepo) GetByID(ctx context.Context, id uint) (*domain.Coach, error) {
	var coach domain.Coach
	if err := r.db.WithContext(ctx).First(&coach, id).Error; err != nil {
		return nil, err
	}
	return &coach, nil
}

// GetByIDUnscoped returns coach by id, including soft-deleted rows.
func (r *CoachRepo) GetByIDUnscoped(ctx context.Context, id uint) (*domain.Coach, error) {
	var coach domain.Coach
	if err := r.db.WithContext(ctx).Unscoped().First(&coach, id).Error; err != nil {
		return nil, err
	}
	return &coach, nil
}

// List returns coaches, optionally including inactive and/or soft-deleted rows.
func (r *CoachRepo) List(ctx context.Context, includeInactive, includeDeleted bool) ([]domain.Coach, error) {
	query := r.db.WithContext(ctx).Order("last_name asc")
	if includeDeleted {
		query = query.Unscoped()
	}
	if !includeInactive {
		query = query.Where("is_active = ?", true)
	}
	var coaches []domain.Coach
	if err := query.Find(&coaches).Error; err != nil {
		return nil, err
	}
	return coaches, nil
}

// Update saves coach changes.
func (r *CoachRepo) Update(ctx context.Context, coach *domain.Coach) error {
	return r.db.WithContext(ctx).Save(coach).Error
}

// Delete permanently removes a coach row (hard delete, bypasses soft delete).
func (r *CoachRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Unscoped().Delete(&domain.Coach{}, id).Error
}

// Restore clears the soft-delete flag and re-activates the coach. Uses
// Unscoped so that already-soft-deleted rows can be revived.
func (r *CoachRepo) Restore(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Unscoped().Model(&domain.Coach{}).
		Where("id = ?", id).
		Updates(map[string]any{"deleted_at": nil, "is_active": true}).Error
}
