package service

import (
	"context"

	"le-studio-api/internal/domain"
	"le-studio-api/internal/dto"
	"le-studio-api/internal/repository"
	"le-studio-api/pkg/pagination"
)

// AttendanceService handles attendance workflows.
type AttendanceService struct {
	repos repository.Repositories
}

// NewAttendanceService creates attendance service.
func NewAttendanceService(repos repository.Repositories) *AttendanceService { return &AttendanceService{repos: repos} }

// Mark creates attendance record.
func (s *AttendanceService) Mark(ctx context.Context, adminID uint, payload dto.MarkAttendanceDTO) (*domain.Attendance, error) {
	attendance := &domain.Attendance{
		BookingID:  payload.BookingID,
		MarkedByID: adminID,
		Notes:      payload.Notes,
	}
	if err := s.repos.Attendance.Create(ctx, attendance); err != nil {
		return nil, err
	}
	return attendance, nil
}

// List returns attendance records.
func (s *AttendanceService) List(ctx context.Context, params pagination.Params) ([]domain.Attendance, int64, error) {
	return s.repos.Attendance.List(ctx, params)
}

// Update updates attendance notes.
func (s *AttendanceService) Update(ctx context.Context, id uint, payload dto.UpdateAttendanceDTO) (*domain.Attendance, error) {
	attendance, err := s.repos.Attendance.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	attendance.Notes = payload.Notes
	if err := s.repos.Attendance.Update(ctx, attendance); err != nil {
		return nil, err
	}
	return attendance, nil
}

// Delete removes attendance.
func (s *AttendanceService) Delete(ctx context.Context, id uint) error {
	return s.repos.Attendance.Delete(ctx, id)
}
