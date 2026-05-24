package postgres

import (
	"context"

	"le-studio-api/internal/domain"

	"gorm.io/gorm"
)

// PackTemplateRepo is pack-template repository.
type PackTemplateRepo struct{ db *gorm.DB }

// NewPackTemplateRepo creates pack-template repository.
func NewPackTemplateRepo(db *gorm.DB) *PackTemplateRepo { return &PackTemplateRepo{db: db} }

// Create inserts pack template.
func (r *PackTemplateRepo) Create(ctx context.Context, tpl *domain.PackTemplate) error {
	return r.db.WithContext(ctx).Create(tpl).Error
}

// GetByID returns pack template by id.
func (r *PackTemplateRepo) GetByID(ctx context.Context, id uint) (*domain.PackTemplate, error) {
	var tpl domain.PackTemplate
	if err := r.db.WithContext(ctx).Preload("TrainingType").First(&tpl, id).Error; err != nil {
		return nil, err
	}
	return &tpl, nil
}

// List returns pack templates.
func (r *PackTemplateRepo) List(ctx context.Context, includeInactive bool) ([]domain.PackTemplate, error) {
	query := r.db.WithContext(ctx).Preload("TrainingType").Order("display_order asc")
	if !includeInactive {
		query = query.Where("is_active = ?", true)
	}
	var templates []domain.PackTemplate
	if err := query.Find(&templates).Error; err != nil {
		return nil, err
	}
	return templates, nil
}

// Update saves pack template.
func (r *PackTemplateRepo) Update(ctx context.Context, tpl *domain.PackTemplate) error {
	return r.db.WithContext(ctx).Save(tpl).Error
}

// Delete removes pack template.
func (r *PackTemplateRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.PackTemplate{}, id).Error
}
