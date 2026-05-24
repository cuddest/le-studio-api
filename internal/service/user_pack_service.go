package service

import (
	"context"
	"time"

	"le-studio-api/internal/domain"
	"le-studio-api/internal/dto"
	"le-studio-api/internal/repository"
	"le-studio-api/pkg/pagination"
)

// UserPackService handles user pack workflows.
type UserPackService struct {
	repos repository.Repositories
}

// NewUserPackService creates user pack service.
func NewUserPackService(repos repository.Repositories) *UserPackService { return &UserPackService{repos: repos} }

// Purchase creates a pack for user.
func (s *UserPackService) Purchase(ctx context.Context, userID uint, payload dto.CreateUserPackDTO) (*domain.UserPack, error) {
	packTemplate, err := s.repos.Templates.GetByID(ctx, payload.PackTemplateID)
	if err != nil {
		return nil, err
	}
	pack := &domain.UserPack{
		UserID:         userID,
		PackTemplateID: packTemplate.ID,
		TotalSessions:  packTemplate.NumberOfSessions,
		UsedSessions:   0,
		IsPaid:         true,
		Status:         "active",
	}
	if err := s.repos.UserPacks.Create(ctx, pack); err != nil {
		return nil, err
	}
	return pack, nil
}

// ListByUser returns user packs.
func (s *UserPackService) ListByUser(ctx context.Context, userID uint) ([]domain.UserPack, error) {
	return s.repos.UserPacks.ListByUser(ctx, userID)
}

// AdminList returns paginated packs.
func (s *UserPackService) AdminList(ctx context.Context, params pagination.Params) ([]domain.UserPack, int64, error) {
	return s.repos.UserPacks.List(ctx, params)
}

// AdminGet returns pack by id.
func (s *UserPackService) AdminGet(ctx context.Context, id uint) (*domain.UserPack, error) {
	return s.repos.UserPacks.GetByID(ctx, id)
}

// AdminUpdate updates pack notes and expiry.
func (s *UserPackService) AdminUpdate(ctx context.Context, id uint, payload dto.UpdateUserPackDTO) (*domain.UserPack, error) {
	pack, err := s.repos.UserPacks.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	pack.Notes = payload.Notes
	pack.ExpiresAt = payload.ExpiresAt
	if err := s.repos.UserPacks.Update(ctx, pack); err != nil {
		return nil, err
	}
	return pack, nil
}

// AdminMarkPaid marks pack as paid.
func (s *UserPackService) AdminMarkPaid(ctx context.Context, id uint) (*domain.UserPack, error) {
	pack, err := s.repos.UserPacks.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !pack.IsPaid {
		now := time.Now()
		pack.IsPaid = true
		pack.PaidAt = &now
	}
	if err := s.repos.UserPacks.Update(ctx, pack); err != nil {
		return nil, err
	}
	return pack, nil
}

// AdminAdjust adjusts used sessions.
func (s *UserPackService) AdminAdjust(ctx context.Context, id uint, payload dto.AdjustUserPackDTO) (*domain.UserPack, error) {
	pack, err := s.repos.UserPacks.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	pack.UsedSessions = payload.UsedSessions
	if pack.UsedSessions > pack.TotalSessions {
		pack.Status = "exhausted"
	} else if pack.Status == "exhausted" {
		pack.Status = "active"
	}
	if err := s.repos.UserPacks.Update(ctx, pack); err != nil {
		return nil, err
	}
	return pack, nil
}

// AdminDelete removes pack.
func (s *UserPackService) AdminDelete(ctx context.Context, id uint) error {
	return s.repos.UserPacks.Delete(ctx, id)
}
