package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"le-studio-api/internal/domain"
	"le-studio-api/internal/dto"
	"le-studio-api/internal/repository"
	"le-studio-api/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthService handles auth workflows.
type AuthService struct {
	repos              repository.Repositories
	db                 *gorm.DB
	jwtSecret          string
	accessTokenTTL     time.Duration
	refreshTokenTTL    time.Duration
}

// NewAuthService creates auth service.
func NewAuthService(repos repository.Repositories, db *gorm.DB, jwtSecret string, accessTTL, refreshTTL time.Duration) *AuthService {
	return &AuthService{repos: repos, db: db, jwtSecret: jwtSecret, accessTokenTTL: accessTTL, refreshTokenTTL: refreshTTL}
}

// Register creates a user account and tokens.
func (s *AuthService) Register(ctx context.Context, payload dto.RegisterDTO, userAgent, ip string) (string, string, *domain.User, error) {
	email := strings.ToLower(strings.TrimSpace(payload.Email))
	if _, err := s.repos.Users.FindByEmail(ctx, email); err == nil {
		return "", "", nil, gorm.ErrDuplicatedKey
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", "", nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", nil, err
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
		return "", "", nil, err
	}

	accessToken, refreshToken, err := s.issueTokens(ctx, user.ID, "user", user.Email, userAgent, ip)
	if err != nil {
		return "", "", nil, err
	}
	return accessToken, refreshToken, user, nil
}

// Login validates credentials and issues tokens.
func (s *AuthService) Login(ctx context.Context, payload dto.LoginDTO, userAgent, ip string) (string, string, *domain.User, error) {
	email := strings.ToLower(strings.TrimSpace(payload.Email))
	user, err := s.repos.Users.FindByEmail(ctx, email)
	if err != nil {
		return "", "", nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(payload.Password)); err != nil {
		return "", "", nil, errors.New("invalid credentials")
	}
	accessToken, refreshToken, err := s.issueTokens(ctx, user.ID, "user", user.Email, userAgent, ip)
	if err != nil {
		return "", "", nil, err
	}
	return accessToken, refreshToken, user, nil
}

// AdminLogin validates admin credentials and issues tokens.
func (s *AuthService) AdminLogin(ctx context.Context, payload dto.LoginDTO, userAgent, ip string) (string, string, *domain.Admin, error) {
	email := strings.ToLower(strings.TrimSpace(payload.Email))
	admin, err := s.repos.Admins.FindByEmail(ctx, email)
	if err != nil {
		return "", "", nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(payload.Password)); err != nil {
		return "", "", nil, errors.New("invalid credentials")
	}
	accessToken, refreshToken, err := s.issueTokens(ctx, admin.ID, "admin", admin.Email, userAgent, ip)
	if err != nil {
		return "", "", nil, err
	}
	return accessToken, refreshToken, admin, nil
}

// Refresh rotates refresh token and returns new access/refresh tokens.
func (s *AuthService) Refresh(ctx context.Context, refreshToken, userAgent, ip string) (string, string, uint, string, error) {
	hash := jwt.HashToken(refreshToken)
	stored, err := s.repos.RefreshToken.FindByHash(ctx, hash)
	if err != nil {
		return "", "", 0, "", err
	}
	if stored.RevokedAt != nil || time.Now().After(stored.ExpiresAt) {
		return "", "", 0, "", errors.New("refresh token expired")
	}
	if err := s.repos.RefreshToken.Revoke(ctx, stored); err != nil {
		return "", "", 0, "", err
	}

	email := ""
	if stored.Role == "admin" {
		admin, err := s.repos.Admins.GetByID(ctx, stored.SubjectID)
		if err != nil {
			return "", "", 0, "", err
		}
		email = admin.Email
	} else {
		user, err := s.repos.Users.GetByID(ctx, stored.SubjectID)
		if err != nil {
			return "", "", 0, "", err
		}
		email = user.Email
	}
	accessToken, newRefreshToken, err := s.issueTokens(ctx, stored.SubjectID, stored.Role, email, userAgent, ip)
	if err != nil {
		return "", "", 0, "", err
	}
	return accessToken, newRefreshToken, stored.SubjectID, stored.Role, nil
}

// Logout revokes the refresh token.
func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	hash := jwt.HashToken(refreshToken)
	stored, err := s.repos.RefreshToken.FindByHash(ctx, hash)
	if err != nil {
		return err
	}
	return s.repos.RefreshToken.Revoke(ctx, stored)
}

// GuestPurchase creates a guest user and initial pack purchase.
func (s *AuthService) GuestPurchase(ctx context.Context, payload dto.GuestPurchaseDTO, userAgent, ip string) (string, string, *domain.User, *domain.UserPack, error) {
	email := strings.ToLower(strings.TrimSpace(payload.Email))
	if _, err := s.repos.Users.FindByEmail(ctx, email); err == nil {
		return "", "", nil, nil, gorm.ErrDuplicatedKey
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", "", nil, nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(randomFallbackPassword()), bcrypt.DefaultCost)
	if err != nil {
		return "", "", nil, nil, err
	}

	user := &domain.User{
		FirstName:    payload.FirstName,
		LastName:     payload.LastName,
		Email:        email,
		Phone:        payload.Phone,
		PasswordHash: string(passwordHash),
		IsGuest:      true,
		IsActive:     true,
	}

	var createdPack *domain.UserPack
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		var template domain.PackTemplate
		if err := tx.First(&template, payload.PackTemplateID).Error; err != nil {
			return err
		}
		pack := &domain.UserPack{
			UserID:         user.ID,
			PackTemplateID: template.ID,
			TotalSessions:  template.NumberOfSessions,
			UsedSessions:   0,
			IsPaid:         true,
			Status:         "active",
		}
		if err := tx.Create(pack).Error; err != nil {
			return err
		}
		createdPack = pack
		return nil
	}); err != nil {
		return "", "", nil, nil, err
	}

	accessToken, refreshToken, err := s.issueTokens(ctx, user.ID, "user", user.Email, userAgent, ip)
	if err != nil {
		return "", "", nil, nil, err
	}
	return accessToken, refreshToken, user, createdPack, nil
}

func (s *AuthService) issueTokens(ctx context.Context, subjectID uint, role, email, userAgent, ip string) (string, string, error) {
	accessToken, err := jwt.GenerateAccessToken(s.jwtSecret, subjectID, role, email, s.accessTokenTTL)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := randomToken(32)
	if err != nil {
		return "", "", err
	}
	hash := jwt.HashToken(refreshToken)
	refresh := &domain.RefreshToken{
		SubjectID: subjectID,
		Role:      role,
		TokenHash: hash,
		ExpiresAt: time.Now().Add(s.refreshTokenTTL),
		UserAgent: userAgent,
		IPAddress: ip,
	}
	if err := s.repos.RefreshToken.Create(ctx, refresh); err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func randomToken(length int) (string, error) {
	buf := make([]byte, length)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func randomFallbackPassword() string {
	return "guest-" + time.Now().Format("20060102150405")
}
