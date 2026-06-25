package service

import (
	"context"
	"errors"

	"le-studio-api/internal/domain"
	"le-studio-api/internal/repository"
)

// ErrTrainingTypeInvalidParent is returned when a training type's parent is
// invalid: the parent must exist, must not be the type itself, and must itself
// be a top-level type (no parent of its own). This enforces a flat hierarchy.
var ErrTrainingTypeInvalidParent = errors.New("parent must be a top-level training type")

// TrainingTypeService handles training type workflows.
type TrainingTypeService struct {
	repos repository.Repositories
}

// NewTrainingTypeService creates TrainingTypeService.
func NewTrainingTypeService(repos repository.Repositories) *TrainingTypeService {
	return &TrainingTypeService{repos: repos}
}

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
	if err := s.validateParent(ctx, training.ID, training.ParentID); err != nil {
		return nil, err
	}
	if err := s.repos.Training.Create(ctx, training); err != nil {
		return nil, err
	}
	return training, nil
}

// Update updates training type.
func (s *TrainingTypeService) Update(ctx context.Context, training *domain.TrainingType) (*domain.TrainingType, error) {
	if err := s.validateParent(ctx, training.ID, training.ParentID); err != nil {
		return nil, err
	}
	if err := s.repos.Training.Update(ctx, training); err != nil {
		return nil, err
	}
	return training, nil
}

// Delete removes training type.
func (s *TrainingTypeService) Delete(ctx context.Context, id uint) error {
	return s.repos.Training.Delete(ctx, id)
}

// validateParent enforces: parent must exist, must not be the type itself, and
// must itself be top-level (no parent of its own).
func (s *TrainingTypeService) validateParent(ctx context.Context, selfID uint, parentID *uint) error {
	if parentID == nil {
		return nil
	}
	if selfID != 0 && *parentID == selfID {
		return ErrTrainingTypeInvalidParent
	}
	parent, err := s.repos.Training.GetByID(ctx, *parentID)
	if err != nil {
		return err
	}
	if parent.ParentID != nil {
		return ErrTrainingTypeInvalidParent
	}
	return nil
}
