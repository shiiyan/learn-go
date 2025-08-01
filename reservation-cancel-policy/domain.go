package main

import (
	"time"

	"github.com/pkg/errors"
)

type Clock interface {
	Now() time.Time
}
type RealClock struct{}

func (c RealClock) Now() time.Time {
	return time.Now()
}

type ReservationID string
type UserID string

type Status string

const (
	StatusActive    Status = "active"
	StatusCancelled Status = "cancelled"
)

type Money struct {
	Amount   int64
	Currency string
}

type CancellationPolicy interface {
	CanCancel(reservation *Reservation, canceller Canceller) error
	CancelWithoutRefund() bool
}

type EndUserCancellationPolicy struct{}

func (p EndUserCancellationPolicy) CanCancel(reservation *Reservation, canceller Canceller) error {
	// end user can only cancel their own reservations
	if reservation.UserID != canceller.GetID() {
		return errors.New("canceller is not the owner of the reservation")
	}

	return nil
}

func (p EndUserCancellationPolicy) CancelWithoutRefund() bool {
	// end user cannot cancel without refund
	return false
}

type AdminCancellationWithoutRefundPolicy struct{}

func (p AdminCancellationWithoutRefundPolicy) CanCancel(reservation *Reservation, canceller Canceller) error {
	// admin can cancel any reservation
	return nil
}

func (p AdminCancellationWithoutRefundPolicy) CancelWithoutRefund() bool {
	// admin can cancel without refund
	return true
}

type AdminCancellationWithRefundPolicy struct{}

func (p AdminCancellationWithRefundPolicy) CanCancel(reservation *Reservation, canceller Canceller) error {
	// admin can cancel any reservation
	return nil
}

func (p AdminCancellationWithRefundPolicy) CancelWithoutRefund() bool {
	// admin cancels with refund (not without)
	return false
}

func NewCancellationPolicy(userRole Role, shouldRefund bool) CancellationPolicy {
	switch userRole {
	case RoleAdmin:
		if shouldRefund {
			return AdminCancellationWithRefundPolicy{}
		}
		return AdminCancellationWithoutRefundPolicy{}
	case RoleEndUser:
		return EndUserCancellationPolicy{}
	}
	return nil
}

type Canceller interface {
	GetID() UserID
	CanCancelWithoutRefund() bool
}

type User struct {
	id   UserID
	role Role
}

type Role string

const (
	RoleEndUser Role = "end_user"
	RoleAdmin   Role = "admin"
)

func (u User) GetID() UserID {
	return u.id
}

func (u User) CanCancelWithoutRefund() bool {
	return u.role == RoleAdmin
}

type Reservation struct {
	id          ReservationID
	UserID      UserID
	status      Status
	amount      Money
	createdAt   time.Time
	cancelledAt *time.Time // nil if not cancelled
	canceller   Canceller // nil if not cancelled
}

func (r *Reservation) Cancel(canceller Canceller, policy CancellationPolicy, clock Clock) (*CancellationResult, error) {
	if r.status == StatusCancelled {
		return nil, errors.New("reservation already cancelled")
	}

	if err := policy.CanCancel(r, canceller); err != nil {
		return nil, errors.Wrap(err, "cannot cancel reservation")
	}

	if policy.CancelWithoutRefund() && !canceller.CanCancelWithoutRefund() {
		return nil, errors.New("canceller cannot cancel without refund")
	}

	money := r.calculateRefund(policy)

	r.status = StatusCancelled
	now := clock.Now()
	r.cancelledAt = &now
	r.canceller = canceller

	return &CancellationResult{
		ReservationID: r.id,
		RefundAmount:  money,
		CancelledAt:   *r.cancelledAt,
		CancelledBy:   canceller,
	}, nil
}

type CancellationResult struct {
	ReservationID ReservationID
	RefundAmount  Money
	CancelledAt   time.Time
	CancelledBy   Canceller
}

func (r *Reservation) calculateRefund(policy CancellationPolicy) Money {
	if policy.CancelWithoutRefund() {
		return Money{Amount: 0, Currency: ""}
	}

	return r.amount
}
