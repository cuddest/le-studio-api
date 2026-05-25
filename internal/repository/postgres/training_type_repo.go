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
	if err := r.db.WithContext(ctx).Preload("Parent").First(&training, id).Error; err != nil {
		return nil, err
	}
	return &training, nil
}

// GetByIDs returns multiple training types by ids.
func (r *TrainingTypeRepo) GetByIDs(ctx context.Context, ids []uint) ([]domain.TrainingType, error) {
	var types []domain.TrainingType
	if len(ids) == 0 {
		return types, nil
	}
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&types).Error; err != nil {
		return nil, err
	}
	return types, nil
}

// List returns training types.
func (r *TrainingTypeRepo) List(ctx context.Context, includeInactive bool) ([]domain.TrainingType, error) {
	query := r.db.WithContext(ctx).Preload("Parent").Order("name asc")
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
