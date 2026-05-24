package service

import (
	"context"

	"le-studio-api/internal/domain"
	"le-studio-api/internal/dto"
	"le-studio-api/internal/repository"
)

// CoachService handles coach workflows.
type CoachService struct {
	repos repository.Repositories
}

// NewCoachService creates coach service.
func NewCoachService(repos repository.Repositories) *CoachService { return &CoachService{repos: repos} }

// List returns coaches.
func (s *CoachService) List(ctx context.Context, includeInactive bool) ([]domain.Coach, error) {
	return s.repos.Coaches.List(ctx, includeInactive)
}

// Get returns coach by id.
func (s *CoachService) Get(ctx context.Context, id uint) (*domain.Coach, error) {
	return s.repos.Coaches.GetByID(ctx, id)
}

// Create creates coach.
func (s *CoachService) Create(ctx context.Context, payload dto.CreateCoachDTO) (*domain.Coach, error) {
	coach := &domain.Coach{
		FirstName:   payload.FirstName,
		LastName:    payload.LastName,
		Bio:         payload.Bio,
		PhotoURL:    payload.PhotoURL,
		Specialties: payload.Specialties,
		IsActive:    true,
	}
	if payload.IsActive != nil {
		coach.IsActive = *payload.IsActive
	}
	if err := s.repos.Coaches.Create(ctx, coach); err != nil {
		return nil, err
	}
	return coach, nil
}

// Update updates coach.
func (s *CoachService) Update(ctx context.Context, id uint, payload dto.CreateCoachDTO) (*domain.Coach, error) {
	coach, err := s.repos.Coaches.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	coach.FirstName = payload.FirstName
	coach.LastName = payload.LastName
	coach.Bio = payload.Bio
	coach.PhotoURL = payload.PhotoURL
	coach.Specialties = payload.Specialties
	if payload.IsActive != nil {
		coach.IsActive = *payload.IsActive
	}
	if err := s.repos.Coaches.Update(ctx, coach); err != nil {
		return nil, err
	}
	return coach, nil
}

// Delete removes coach.
func (s *CoachService) Delete(ctx context.Context, id uint) error {
	return s.repos.Coaches.Delete(ctx, id)
}

// ToggleActive flips coach active status.
func (s *CoachService) ToggleActive(ctx context.Context, id uint) (*domain.Coach, error) {
	coach, err := s.repos.Coaches.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	coach.IsActive = !coach.IsActive
	if err := s.repos.Coaches.Update(ctx, coach); err != nil {
		return nil, err
	}
	return coach, nil
}
