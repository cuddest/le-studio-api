package postgres

import (
	"context"

	"le-studio-api/internal/domain"

	"gorm.io/gorm"
)

// TrainingTypeRepo is training-type repository.
type TrainingTypeRepo struct{ db *gorm.DB }

// NewTrainingTypeRepo creates training-type repository.
func NewTrainingTypeRepo(db *gorm.DB) *TrainingTypeRepo { return &TrainingTypeRepo{db: db} }

// Create inserts training type.
func (r *TrainingTypeRepo) Create(ctx context.Context, training *domain.TrainingType) error {
	return r.db.WithContext(ctx).Create(training).Error
}

// GetByID returns training type by id.
func (r *TrainingTypeRepo) GetByID(ctx context.Context, id uint) (*domain.TrainingType, error) {
	var training domain.TrainingType
	if err := r.db.WithContext(ctx).First(&training, id).Error; err != nil {
		return nil, err
	}
	return &training, nil
}

// List returns training types.
func (r *TrainingTypeRepo) List(ctx context.Context, includeInactive bool) ([]domain.TrainingType, error) {
	query := r.db.WithContext(ctx).Order("name asc")
	if !includeInactive {
		query = query.Where("is_active = ?", true)
	}
	var types []domain.TrainingType
	if err := query.Find(&types).Error; err != nil {
		return nil, err
	}
	return types, nil
}

// Update saves training type.
func (r *TrainingTypeRepo) Update(ctx context.Context, training *domain.TrainingType) error {
	return r.db.WithContext(ctx).Save(training).Error
}

// Delete removes training type.
func (r *TrainingTypeRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.TrainingType{}, id).Error
}
