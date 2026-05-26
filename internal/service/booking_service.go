package service

import (
	"context"
	"errors"
	"time"

	"le-studio-api/internal/domain"
	"le-studio-api/internal/dto"
	"le-studio-api/internal/repository"
	"le-studio-api/pkg/pagination"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// BookingService handles booking workflows.
type BookingService struct {
	repos repository.Repositories
	db    *gorm.DB
}

// NewBookingService creates booking service.
func NewBookingService(repos repository.Repositories, db *gorm.DB) *BookingService {
	return &BookingService{repos: repos, db: db}
}

// Create creates booking for user.
func (s *BookingService) Create(ctx context.Context, userID uint, payload dto.CreateBookingDTO) (*domain.Booking, error) {
	return s.createForUser(ctx, userID, payload.SlotID, payload.UserPackID)
}

// AdminCreate creates booking on behalf of a user.
func (s *BookingService) AdminCreate(ctx context.Context, payload dto.AdminCreateBookingDTO) (*domain.Booking, error) {
	return s.createForUser(ctx, payload.UserID, payload.SlotID, payload.UserPackID)
}

func (s *BookingService) createForUser(ctx context.Context, userID, slotID, userPackID uint) (*domain.Booking, error) {
	var created *domain.Booking
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var slot domain.Slot
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Preload("TrainingType").Preload("TrainingType.Parent").First(&slot, slotID).Error; err != nil {
			return err
		}
		if slot.IsCancelled {
			return errors.New("slot is cancelled")
		}
		if slot.BookedCount >= slot.Capacity {
			return errors.New("slot is full")
		}

		var pack domain.UserPack
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Preload("PackTemplate").Preload("PackTemplate.TrainingTypes").First(&pack, userPackID).Error; err != nil {
			return err
		}
		if pack.UserID != userID {
			return errors.New("pack does not belong to user")
		}
		if !pack.IsPaid {
			return errors.New("pack is unpaid")
		}
		if pack.IsExhausted() || pack.Status != "active" {
			return errors.New("pack has no remaining sessions")
		}

		// Validate pack includes the slot's training type (or one of its parent types)
		allowed := false
		allowedMap := map[uint]bool{}
		if pack.PackTemplate.ID != 0 {
			for _, t := range pack.PackTemplate.TrainingTypes {
				allowedMap[t.ID] = true
			}
			// if pack template has a single TrainingTypeID (legacy), include it too
			if pack.PackTemplate.TrainingTypeID != 0 {
				allowedMap[pack.PackTemplate.TrainingTypeID] = true
			}
		}
		// walk up slot training type parents
		p := &slot.TrainingType
		for p != nil {
			if allowedMap[p.ID] {
				allowed = true
				break
			}
			p = p.Parent
		}
		if !allowed {
			return errors.New("pack does not include this training type")
		}

		booking := &domain.Booking{
			UserID:     userID,
			SlotID:     slot.ID,
			UserPackID: pack.ID,
			Status:     "confirmed",
		}
		if err := tx.Create(booking).Error; err != nil {
			return err
		}
		slot.BookedCount += 1
		if err := tx.Save(&slot).Error; err != nil {
			return err
		}
		pack.UsedSessions += 1
		if pack.IsExhausted() {
			pack.Status = "exhausted"
		}
		if err := tx.Save(&pack).Error; err != nil {
			return err
		}
		created = booking
		return nil
	}); err != nil {
		return nil, err
	}
	return s.repos.Bookings.GetByID(ctx, created.ID)
}

// Cancel cancels a user booking.
func (s *BookingService) Cancel(ctx context.Context, userID, bookingID uint) (*domain.Booking, error) {
	return s.cancelInternal(ctx, bookingID, &userID)
}

// AdminCancel cancels booking without ownership check.
func (s *BookingService) AdminCancel(ctx context.Context, bookingID uint) (*domain.Booking, error) {
	return s.cancelInternal(ctx, bookingID, nil)
}

func (s *BookingService) cancelInternal(ctx context.Context, bookingID uint, userID *uint) (*domain.Booking, error) {
	var updated *domain.Booking
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var booking domain.Booking
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&booking, bookingID).Error; err != nil {
			return err
		}
		if userID != nil && booking.UserID != *userID {
			return errors.New("booking does not belong to user")
		}
		if booking.Status == "cancelled" {
			updated = &booking
			return nil
		}
		now := time.Now()
		booking.Status = "cancelled"
		booking.CancelledAt = &now
		if err := tx.Save(&booking).Error; err != nil {
			return err
		}
		var slot domain.Slot
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&slot, booking.SlotID).Error; err != nil {
			return err
		}
		if slot.BookedCount > 0 {
			slot.BookedCount -= 1
			if err := tx.Save(&slot).Error; err != nil {
				return err
			}
		}
		var pack domain.UserPack
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&pack, booking.UserPackID).Error; err != nil {
			return err
		}
		if pack.UsedSessions > 0 {
			pack.UsedSessions -= 1
			if pack.Status == "exhausted" && pack.RemainingSessions() > 0 {
				pack.Status = "active"
			}
			if err := tx.Save(&pack).Error; err != nil {
				return err
			}
		}
		updated = &booking
		return nil
	}); err != nil {
		return nil, err
	}
	return s.repos.Bookings.GetByID(ctx, updated.ID)
}

// Get returns booking by id.
func (s *BookingService) Get(ctx context.Context, bookingID uint) (*domain.Booking, error) {
	return s.repos.Bookings.GetByID(ctx, bookingID)
}

// ListByUser returns user bookings.
func (s *BookingService) ListByUser(ctx context.Context, userID uint) ([]domain.Booking, error) {
	return s.repos.Bookings.ListByUser(ctx, userID)
}

// AdminList returns paginated bookings.
func (s *BookingService) AdminList(ctx context.Context, params pagination.Params) ([]domain.Booking, int64, error) {
	return s.repos.Bookings.List(ctx, params)
}
