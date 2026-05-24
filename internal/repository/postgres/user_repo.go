package postgres

import (
	"context"

	"le-studio-api/internal/domain"
	"le-studio-api/pkg/pagination"

	"gorm.io/gorm"
)

// UserRepo is user repository.
type UserRepo struct{ db *gorm.DB }

// NewUserRepo creates user repository.
func NewUserRepo(db *gorm.DB) *UserRepo { return &UserRepo{db: db} }

// Create inserts a new user.
func (r *UserRepo) Create(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// FindByEmail returns user by email.
func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByID returns user by id.
func (r *UserRepo) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Update saves user changes.
func (r *UserRepo) Update(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// List returns paginated users and total count.
func (r *UserRepo) List(ctx context.Context, params pagination.Params) ([]domain.User, int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&domain.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var users []domain.User
	if err := r.db.WithContext(ctx).Order("created_at desc").Limit(params.Limit).Offset(params.Offset).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

// Delete removes user by id.
func (r *UserRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.User{}, id).Error
}
