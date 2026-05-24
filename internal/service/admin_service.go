package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"le-studio-api/internal/domain"
	"le-studio-api/internal/dto"
	"le-studio-api/internal/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AdminService handles admin workflows.
type AdminService struct {
	db    *gorm.DB
	repos repository.Repositories
}

// NewAdminService creates admin service.
func NewAdminService(repos repository.Repositories, db *gorm.DB) *AdminService {
	return &AdminService{db: db, repos: repos}
}

// StatsOverview aggregates dashboard stats.
func (s *AdminService) StatsOverview(ctx context.Context) (map[string]any, error) {
	var totalUsers int64
	if err := s.db.WithContext(ctx).Model(&domain.User{}).Count(&totalUsers).Error; err != nil {
		return nil, err
	}
	var activePacks int64
	if err := s.db.WithContext(ctx).Model(&domain.UserPack{}).Where("status = ?", "active").Count(&activePacks).Error; err != nil {
		return nil, err
	}
	var revenue float64
	_ = s.db.WithContext(ctx).
		Table("user_packs").
		Select("coalesce(sum(pack_templates.price), 0)").
		Joins("join pack_templates on pack_templates.id = user_packs.pack_template_id").
		Where("user_packs.is_paid = ?", true).
		Scan(&revenue).Error

	weekStart := time.Now().AddDate(0, 0, -int(time.Now().Weekday()))
	weekEnd := weekStart.AddDate(0, 0, 7)
	var sessionsThisWeek int64
	_ = s.db.WithContext(ctx).
		Model(&domain.Booking{}).
		Where("created_at >= ? AND created_at < ?", weekStart, weekEnd).
		Count(&sessionsThisWeek).Error

	var attendanceCount int64
	_ = s.db.WithContext(ctx).
		Model(&domain.Attendance{}).
		Where("created_at >= ? AND created_at < ?", weekStart, weekEnd).
		Count(&attendanceCount).Error
	attendanceRate := 0.0
	if sessionsThisWeek > 0 {
		attendanceRate = float64(attendanceCount) / float64(sessionsThisWeek)
	}

	var upcomingSlots int64
	_ = s.db.WithContext(ctx).
		Model(&domain.Slot{}).
		Where("date >= ? AND is_cancelled = ?", time.Now(), false).
		Count(&upcomingSlots).Error

	var topTraining []map[string]any
	_ = s.db.WithContext(ctx).
		Table("bookings").
		Select("training_types.name as name, count(bookings.id) as count").
		Joins("join slots on slots.id = bookings.slot_id").
		Joins("join training_types on training_types.id = slots.training_type_id").
		Group("training_types.name").
		Order("count desc").
		Limit(5).
		Scan(&topTraining).Error

	return map[string]any{
		"total_users":        totalUsers,
		"active_packs":       activePacks,
		"revenue":            revenue,
		"sessions_this_week": sessionsThisWeek,
		"attendance_rate":    attendanceRate,
		"top_training_types": topTraining,
		"upcoming_slots":     upcomingSlots,
	}, nil
}

// GetProfile returns admin by id.
func (s *AdminService) GetProfile(ctx context.Context, id uint) (*domain.Admin, error) {
	return s.repos.Admins.GetByID(ctx, id)
}

// UpdateProfile updates admin profile fields.
func (s *AdminService) UpdateProfile(ctx context.Context, id uint, payload dto.UpdateAdminDTO) (*domain.Admin, error) {
	admin, err := s.repos.Admins.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if payload.Name != "" {
		admin.Name = payload.Name
	}
	if payload.PhotoURL != "" {
		admin.PhotoURL = payload.PhotoURL
	}
	if err := s.db.WithContext(ctx).Save(admin).Error; err != nil {
		return nil, err
	}
	return admin, nil
}

// ListAdmins returns all admin accounts.
func (s *AdminService) ListAdmins(ctx context.Context) ([]domain.Admin, error) {
	return s.repos.Admins.List(ctx)
}

// CreateAdmin creates a new admin account.
func (s *AdminService) CreateAdmin(ctx context.Context, payload dto.RegisterAdminDTO) (*domain.Admin, error) {
	email := strings.ToLower(strings.TrimSpace(payload.Email))
	if _, err := s.repos.Admins.FindByEmail(ctx, email); err == nil {
		return nil, gorm.ErrDuplicatedKey
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	admin := &domain.Admin{
		Name:         strings.TrimSpace(payload.Name),
		Email:        email,
		PasswordHash: string(passwordHash),
	}
	if err := s.repos.Admins.Create(ctx, admin); err != nil {
		return nil, err
	}
	return admin, nil
}