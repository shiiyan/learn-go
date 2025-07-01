package main

import (
	"fmt"
	"github.com/pkg/errors"
)

// Repository interfaces
type ReservationRepository interface {
	GetByID(id ReservationID) (*Reservation, error)
	Save(r *Reservation) error
	GetByUserID(userID UserID) ([]*Reservation, error)
}

type UserRepository interface {
	GetByID(id UserID) (*User, error)
}

// Payment service interface
type PaymentService interface {
	ProcessRefund(userID UserID, amount Money) error
}

// Notification service interface
type NotificationService interface {
	NotifyCancellation(userID UserID, result *CancellationResult) error
}

// Application Service
type ReservationService struct {
	reservationRepo     ReservationRepository
	userRepo            UserRepository
	paymentService      PaymentService
	notificationService NotificationService
	clock               Clock
}

func NewReservationService(
	reservationRepo ReservationRepository,
	userRepo UserRepository,
	paymentService PaymentService,
	notificationService NotificationService,
	clock Clock,
) *ReservationService {
	return &ReservationService{
		reservationRepo:     reservationRepo,
		userRepo:            userRepo,
		paymentService:      paymentService,
		notificationService: notificationService,
		clock:               clock,
	}
}

// Command for cancelling a reservation
type CancelReservationCommand struct {
	ReservationID string
	CancellerID   string
	shouldRefund bool // Only applicable for admins
}

// Main application service method for cancellation
func (s *ReservationService) CancelReservation(cmd CancelReservationCommand) (*CancellationResult, error) {
	// Load the reservation
	reservation, err := s.reservationRepo.GetByID(ReservationID(cmd.ReservationID))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get reservation")
	}

	// Load the canceller (user)
	canceller, err := s.userRepo.GetByID(UserID(cmd.CancellerID))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get canceller")
	}

	// Create appropriate policy based on user role and command
	// For admins, respect the ForceNoRefund flag
	policy := NewCancellationPolicy(canceller.role, cmd.shouldRefund)
	
	if policy == nil {
		return nil, errors.New("invalid cancellation policy")
	}

	// Execute domain logic
	result, err := reservation.Cancel(canceller, policy, s.clock)
	if err != nil {
		return nil, errors.Wrap(err, "failed to cancel reservation")
	}

	// Handle side effects

	// 1. Process refund if needed
	if result.RefundAmount.Amount > 0 {
		if err := s.paymentService.ProcessRefund(reservation.UserID, result.RefundAmount); err != nil {
			// In a real application, you might want to handle this with compensation
			return nil, errors.Wrap(err, "failed to process refund")
		}
	}

	// 2. Save the updated reservation
	if err := s.reservationRepo.Save(reservation); err != nil {
		return nil, errors.Wrap(err, "failed to save reservation")
	}

	// 3. Send notification
	if err := s.notificationService.NotifyCancellation(reservation.UserID, result); err != nil {
		// Notification failure shouldn't fail the whole operation
		fmt.Printf("Failed to send notification: %v\n", err)
	}

	return result, nil
}

// Command for creating a reservation
type CreateReservationCommand struct {
	UserID   string
	Amount   int64
	Currency string
}

// Create a new reservation
func (s *ReservationService) CreateReservation(cmd CreateReservationCommand) (*Reservation, error) {
	// Verify user exists
	_, err := s.userRepo.GetByID(UserID(cmd.UserID))
	if err != nil {
		return nil, errors.Wrap(err, "user not found")
	}

	// Create new reservation
	reservation := &Reservation{
		id:        ReservationID(fmt.Sprintf("res-%d", s.clock.Now().Unix())),
		UserID:    UserID(cmd.UserID),
		status:    StatusActive,
		amount:    Money{Amount: cmd.Amount, Currency: cmd.Currency},
		createdAt: s.clock.Now(),
	}

	// Save reservation
	if err := s.reservationRepo.Save(reservation); err != nil {
		return nil, errors.Wrap(err, "failed to save reservation")
	}

	return reservation, nil
}

// Query methods
func (s *ReservationService) GetReservation(id string) (*Reservation, error) {
	return s.reservationRepo.GetByID(ReservationID(id))
}

func (s *ReservationService) GetUserReservations(userID string) ([]*Reservation, error) {
	return s.reservationRepo.GetByUserID(UserID(userID))
}

// DTO for presenting reservation information
type ReservationDTO struct {
	ID          string
	UserID      string
	Status      string
	Amount      int64
	Currency    string
	CreatedAt   string
	CancelledAt string
	CancelledBy string
}

func (s *ReservationService) GetReservationDetails(id string) (*ReservationDTO, error) {
	reservation, err := s.reservationRepo.GetByID(ReservationID(id))
	if err != nil {
		return nil, err
	}

	dto := &ReservationDTO{
		ID:        string(reservation.id),
		UserID:    string(reservation.UserID),
		Status:    string(reservation.status),
		Amount:    reservation.amount.Amount,
		Currency:  reservation.amount.Currency,
		CreatedAt: reservation.createdAt.Format("2006-01-02 15:04:05"),
	}

	if reservation.cancelledAt != nil {
		dto.CancelledAt = reservation.cancelledAt.Format("2006-01-02 15:04:05")
	}

	if reservation.canceller != nil {
		dto.CancelledBy = string(reservation.canceller.GetID())
	}

	return dto, nil
}