package postgres

import (
	"context"

	"le-studio-api/internal/domain"

	"gorm.io/gorm"
)

// AdminRepo is admin repository.
type AdminRepo struct{ db *gorm.DB }

// NewAdminRepo creates admin repository.
func NewAdminRepo(db *gorm.DB) *AdminRepo { return &AdminRepo{db: db} }

// Create persists a new admin record.
func (r *AdminRepo) Create(ctx context.Context, admin *domain.Admin) error {
	return r.db.WithContext(ctx).Create(admin).Error
}

// FindByEmail returns admin by email.
func (r *AdminRepo) FindByEmail(ctx context.Context, email string) (*domain.Admin, error) {
	var admin domain.Admin
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&admin).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

// GetByID returns admin by id.
func (r *AdminRepo) GetByID(ctx context.Context, id uint) (*domain.Admin, error) {
	var admin domain.Admin
	if err := r.db.WithContext(ctx).First(&admin, id).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}
