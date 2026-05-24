package service

import (
	"context"

	"le-studio-api/internal/domain"
	"le-studio-api/internal/repository"
)

// TrainingTypeService handles training type workflows.
type TrainingTypeService struct {
	repos repository.Repositories
}

// NewTrainingTypeService creates training type service.
func NewTrainingTypeService(repos repository.Repositories) *TrainingTypeService { return &TrainingTypeService{repos: repos} }

// List returns training types.
func (s *TrainingTypeService) List(ctx context.Context, includeInactive bool) ([]domain.TrainingType, error) {
	return s.repos.Training.List(ctx, includeInactive)
}

// Get returns training type by id.
func (s *TrainingTypeService) Get(ctx context.Context, id uint) (*domain.TrainingType, error) {
	return s.repos.Training.GetByID(ctx, id)
}

// Create creates training type.
func (s *TrainingTypeService) Create(ctx context.Context, training *domain.TrainingType) (*domain.TrainingType, error) {
	if err := s.repos.Training.Create(ctx, training); err != nil {
		return nil, err
	}
	return training, nil
}

// Update updates training type.
func (s *TrainingTypeService) Update(ctx context.Context, training *domain.TrainingType) (*domain.TrainingType, error) {
	if err := s.repos.Training.Update(ctx, training); err != nil {
		return nil, err
	}
	return training, nil
}

// Delete removes training type.
func (s *TrainingTypeService) Delete(ctx context.Context, id uint) error {
	return s.repos.Training.Delete(ctx, id)
}