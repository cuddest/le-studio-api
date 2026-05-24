package service

import (
	"context"

	"le-studio-api/internal/domain"
	"le-studio-api/internal/dto"
	"le-studio-api/internal/repository"
)

// PackTemplateService handles pack template workflows.
type PackTemplateService struct {
	repos repository.Repositories
}

// NewPackTemplateService creates pack template service.
func NewPackTemplateService(repos repository.Repositories) *PackTemplateService {
	return &PackTemplateService{repos: repos}
}

// List returns templates.
func (s *PackTemplateService) List(ctx context.Context, includeInactive bool) ([]domain.PackTemplate, error) {
	return s.repos.Templates.List(ctx, includeInactive)
}

// Get returns template by id.
func (s *PackTemplateService) Get(ctx context.Context, id uint) (*domain.PackTemplate, error) {
	return s.repos.Templates.GetByID(ctx, id)
}

// Create creates a template.
func (s *PackTemplateService) Create(ctx context.Context, payload dto.CreatePackTemplateDTO) (*domain.PackTemplate, error) {
	tpl := &domain.PackTemplate{
		Name:             payload.Name,
		NumberOfSessions: payload.NumberOfSessions,
		Price:            payload.Price,
		TrainingTypeID:   payload.TrainingTypeID,
		IsActive:         true,
		DisplayOrder:     payload.DisplayOrder,
	}
	if payload.IsActive != nil {
		tpl.IsActive = *payload.IsActive
	}
	if err := s.repos.Templates.Create(ctx, tpl); err != nil {
		return nil, err
	}
	return tpl, nil
}

// Update updates template.
func (s *PackTemplateService) Update(ctx context.Context, id uint, payload dto.CreatePackTemplateDTO) (*domain.PackTemplate, error) {
	tpl, err := s.repos.Templates.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	tpl.Name = payload.Name
	tpl.NumberOfSessions = payload.NumberOfSessions
	tpl.Price = payload.Price
	tpl.TrainingTypeID = payload.TrainingTypeID
	tpl.DisplayOrder = payload.DisplayOrder
	if payload.IsActive != nil {
		tpl.IsActive = *payload.IsActive
	}
	if err := s.repos.Templates.Update(ctx, tpl); err != nil {
		return nil, err
	}
	return tpl, nil
}

// Delete removes template.
func (s *PackTemplateService) Delete(ctx context.Context, id uint) error {
	return s.repos.Templates.Delete(ctx, id)
}

// Reorder updates display order for a template.
func (s *PackTemplateService) Reorder(ctx context.Context, id uint, payload dto.ReorderPackTemplateDTO) (*domain.PackTemplate, error) {
	tpl, err := s.repos.Templates.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	tpl.DisplayOrder = payload.DisplayOrder
	if err := s.repos.Templates.Update(ctx, tpl); err != nil {
		return nil, err
	}
	return tpl, nil
}
