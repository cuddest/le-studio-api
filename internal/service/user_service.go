package service

import (
	"context"
	"errors"
	"strings"

	"le-studio-api/internal/domain"
	"le-studio-api/internal/dto"
	"le-studio-api/internal/repository"
	"le-studio-api/pkg/pagination"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserService handles user workflows.
type UserService struct {
	repos repository.Repositories
}

// NewUserService creates user service.
func NewUserService(repos repository.Repositories) *UserService { return &UserService{repos: repos} }

// GetMe returns the authenticated user.
func (s *UserService) GetMe(ctx context.Context, userID uint) (*domain.User, error) {
	return s.repos.Users.GetByID(ctx, userID)
}

// UpdateMe updates profile fields.
func (s *UserService) UpdateMe(ctx context.Context, userID uint, payload dto.UpdateMeDTO) (*domain.User, error) {
	user, err := s.repos.Users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if payload.FirstName != "" {
		user.FirstName = payload.FirstName
	}
	if payload.LastName != "" {
		user.LastName = payload.LastName
	}
	if payload.Phone != "" {
		user.Phone = payload.Phone
	}
	if payload.PhotoURL != "" {
		user.PhotoURL = payload.PhotoURL
	}
	if err := s.repos.Users.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

// ChangePassword updates the user's password.
func (s *UserService) ChangePassword(ctx context.Context, userID uint, payload dto.ChangePasswordDTO) error {
	user, err := s.repos.Users.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(payload.CurrentPassword)); err != nil {
		return errors.New("current password is invalid")
	}
	newHash, err := bcrypt.GenerateFromPassword([]byte(payload.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(newHash)
	return s.repos.Users.Update(ctx, user)
}

// AdminListUsers returns paginated users.
func (s *UserService) AdminListUsers(ctx context.Context, params pagination.Params) ([]domain.User, int64, error) {
	return s.repos.Users.List(ctx, params)
}

// AdminCreateUser creates a new user for admin.
func (s *UserService) AdminCreateUser(ctx context.Context, payload dto.RegisterDTO) (*domain.User, error) {
	email := strings.ToLower(strings.TrimSpace(payload.Email))
	if _, err := s.repos.Users.FindByEmail(ctx, email); err == nil {
		return nil, gorm.ErrDuplicatedKey
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &domain.User{
		FirstName:    payload.FirstName,
		LastName:     payload.LastName,
		Email:        email,
		Phone:        payload.Phone,
		PasswordHash: string(passwordHash),
		IsGuest:      false,
		IsActive:     true,
	}
	if err := s.repos.Users.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

// AdminDeleteUser removes a user by id.
func (s *UserService) AdminDeleteUser(ctx context.Context, id uint) error {
	return s.repos.Users.Delete(ctx, id)
}

// AdminToggleActive toggles user active status.
func (s *UserService) AdminToggleActive(ctx context.Context, id uint) (*domain.User, error) {
	user, err := s.repos.Users.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	user.IsActive = !user.IsActive
	if err := s.repos.Users.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

// AdminPromoteGuest sets a password and marks user non-guest.
func (s *UserService) AdminPromoteGuest(ctx context.Context, id uint, payload dto.PromoteGuestDTO) (*domain.User, error) {
	user, err := s.repos.Users.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	newHash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = string(newHash)
	user.IsGuest = false
	if err := s.repos.Users.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

// AdminGetUser returns user by id for admin.
func (s *UserService) AdminGetUser(ctx context.Context, id uint) (*domain.User, error) {
	return s.repos.Users.GetByID(ctx, id)
}

// AdminUpdateUser updates a user's profile fields.
func (s *UserService) AdminUpdateUser(ctx context.Context, id uint, payload dto.AdminUpdateUserDTO) (*domain.User, error) {
	user, err := s.repos.Users.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if payload.FirstName != "" {
		user.FirstName = payload.FirstName
	}
	if payload.LastName != "" {
		user.LastName = payload.LastName
	}
	if payload.Phone != "" {
		user.Phone = payload.Phone
	}
	if payload.PhotoURL != "" {
		user.PhotoURL = payload.PhotoURL
	}
	if payload.IsActive != nil {
		user.IsActive = *payload.IsActive
	}
	if err := s.repos.Users.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}
