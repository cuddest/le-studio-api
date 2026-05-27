package service

import (
	"context"
	"log"
	"time"

	"le-studio-api/internal/domain"
	"le-studio-api/internal/dto"
	"le-studio-api/internal/repository"
)

// ScheduleService handles schedule workflows.
type ScheduleService struct {
	repos repository.Repositories
}

// NewScheduleService creates schedule service.
func NewScheduleService(repos repository.Repositories) *ScheduleService {
	return &ScheduleService{repos: repos}
}

// List returns schedules.
func (s *ScheduleService) List(ctx context.Context, includeUnpublished bool) ([]domain.WeeklySchedule, error) {
	return s.repos.Schedules.List(ctx, includeUnpublished)
}

// Get returns schedule by id.
func (s *ScheduleService) Get(ctx context.Context, id uint) (*domain.WeeklySchedule, error) {
	return s.repos.Schedules.GetByID(ctx, id)
}

// AdminCreate creates schedule.
func (s *ScheduleService) AdminCreate(ctx context.Context, payload dto.CreateScheduleDTO) (*domain.WeeklySchedule, error) {
	schedule := &domain.WeeklySchedule{
		Title:     payload.Title,
		WeekStart: payload.WeekStart,
		WeekEnd:   payload.WeekEnd,
	}
	if err := s.repos.Schedules.Create(ctx, schedule); err != nil {
		return nil, err
	}
	return schedule, nil
}

// AdminUpdate updates schedule.
func (s *ScheduleService) AdminUpdate(ctx context.Context, id uint, payload dto.CreateScheduleDTO) (*domain.WeeklySchedule, error) {
	schedule, err := s.repos.Schedules.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	schedule.Title = payload.Title
	schedule.WeekStart = payload.WeekStart
	schedule.WeekEnd = payload.WeekEnd
	if err := s.repos.Schedules.Update(ctx, schedule); err != nil {
		return nil, err
	}
	return schedule, nil
}

// AdminPublish marks schedule published.
func (s *ScheduleService) AdminPublish(ctx context.Context, id uint) (*domain.WeeklySchedule, error) {
	schedule, err := s.repos.Schedules.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	schedule.IsPublished = true
	schedule.PublishedAt = &now
	if err := s.repos.Schedules.Update(ctx, schedule); err != nil {
		return nil, err
	}
	return schedule, nil
}

// AdminDelete removes schedule.
func (s *ScheduleService) AdminDelete(ctx context.Context, id uint) error {
	return s.repos.Schedules.Delete(ctx, id)
}

// CreateSlot creates slot within schedule.
func (s *ScheduleService) CreateSlot(ctx context.Context, scheduleID uint, payload dto.CreateSlotDTO) (*domain.Slot, error) {
	slot := &domain.Slot{
		WeeklyScheduleID: scheduleID,
		TrainingTypeID:   payload.TrainingTypeID,
		CoachID:          payload.CoachID,
		SlotType:         payload.SlotType,
		Name:             payload.Name,
		Date:             payload.Date,
		StartTime:        payload.StartTime,
		EndTime:          payload.EndTime,
		Level:            payload.Level,
		Capacity:         payload.Capacity,
		DayOfWeek:        int(payload.Date.Weekday()),
	}
	if err := s.repos.Slots.Create(ctx, slot); err != nil {
		return nil, err
	}
	return slot, nil
}

// UpdateSlot updates slot fields.
func (s *ScheduleService) UpdateSlot(ctx context.Context, id uint, payload dto.CreateSlotDTO) (*domain.Slot, error) {
	slot, err := s.repos.Slots.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	log.Printf("UpdateSlot: before update id=%d training_type_id=%d", id, slot.TrainingTypeID)
	slot.TrainingTypeID = payload.TrainingTypeID
	slot.CoachID = payload.CoachID
	slot.SlotType = payload.SlotType
	slot.Name = payload.Name
	slot.Date = payload.Date
	slot.StartTime = payload.StartTime
	slot.EndTime = payload.EndTime
	slot.Level = payload.Level
	slot.Capacity = payload.Capacity
	slot.DayOfWeek = int(payload.Date.Weekday())
	if err := s.repos.Slots.Update(ctx, slot); err != nil {
		return nil, err
	}
	// Re-fetch to get fresh preloaded relationships
	updated, err := s.repos.Slots.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	log.Printf("UpdateSlot: after update id=%d training_type_id=%d", id, updated.TrainingTypeID)
	return updated, nil
}

// CancelSlot marks slot cancelled.
func (s *ScheduleService) CancelSlot(ctx context.Context, id uint) (*domain.Slot, error) {
	slot, err := s.repos.Slots.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	slot.IsCancelled = true
	if err := s.repos.Slots.Update(ctx, slot); err != nil {
		return nil, err
	}
	return slot, nil
}

// ListSlots returns slots for schedule.
func (s *ScheduleService) ListSlots(ctx context.Context, scheduleID uint, includeCancelled bool) ([]domain.Slot, error) {
	return s.repos.Slots.ListBySchedule(ctx, scheduleID, includeCancelled)
}
