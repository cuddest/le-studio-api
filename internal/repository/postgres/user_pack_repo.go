package postgres

import (
	"context"

	"le-studio-api/internal/domain"
	"le-studio-api/pkg/pagination"

	"gorm.io/gorm"
)

// UserPackRepo is user-pack repository.
type UserPackRepo struct{ db *gorm.DB }

// NewUserPackRepo creates user-pack repository.
func NewUserPackRepo(db *gorm.DB) *UserPackRepo { return &UserPackRepo{db: db} }

// Create inserts user pack.
func (r *UserPackRepo) Create(ctx context.Context, pack *domain.UserPack) error {
	return r.db.WithContext(ctx).Create(pack).Error
}

// GetByID returns user pack by id.
func (r *UserPackRepo) GetByID(ctx context.Context, id uint) (*domain.UserPack, error) {
	var pack domain.UserPack
	if err := r.db.WithContext(ctx).Preload("User").Preload("PackTemplate").Preload("PackTemplate.TrainingTypes").First(&pack, id).Error; err != nil {
		return nil, err
	}
	return &pack, nil
}

// List returns paginated packs.
func (r *UserPackRepo) List(ctx context.Context, params pagination.Params) ([]domain.UserPack, int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&domain.UserPack{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var packs []domain.UserPack
	if err := r.db.WithContext(ctx).Preload("User").Preload("PackTemplate").Preload("PackTemplate.TrainingTypes").Order("created_at desc").Limit(params.Limit).Offset(params.Offset).Find(&packs).Error; err != nil {
		return nil, 0, err
	}
	return packs, total, nil
}

// ListByUser returns packs for a user.
func (r *UserPackRepo) ListByUser(ctx context.Context, userID uint) ([]domain.UserPack, error) {
	var packs []domain.UserPack
	if err := r.db.WithContext(ctx).Preload("PackTemplate").Preload("PackTemplate.TrainingTypes").Where("user_id = ?", userID).Order("created_at desc").Find(&packs).Error; err != nil {
		return nil, err
	}
	return packs, nil
}

// Update saves pack changes.
func (r *UserPackRepo) Update(ctx context.Context, pack *domain.UserPack) error {
	return r.db.WithContext(ctx).Save(pack).Error
}

// Delete removes pack by id.
func (r *UserPackRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.UserPack{}, id).Error
}
