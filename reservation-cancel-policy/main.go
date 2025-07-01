package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/pkg/errors"
)

// In-memory repository implementations
type InMemoryReservationRepository struct {
	mu           sync.RWMutex
	reservations map[ReservationID]*Reservation
}

func NewInMemoryReservationRepository() *InMemoryReservationRepository {
	return &InMemoryReservationRepository{
		reservations: make(map[ReservationID]*Reservation),
	}
}

func (r *InMemoryReservationRepository) GetByID(id ReservationID) (*Reservation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	reservation, exists := r.reservations[id]
	if !exists {
		return nil, errors.New("reservation not found")
	}
	// Return a copy to avoid external modifications
	resCopy := *reservation
	return &resCopy, nil
}

func (r *InMemoryReservationRepository) Save(reservation *Reservation) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.reservations[reservation.id] = reservation
	return nil
}

func (r *InMemoryReservationRepository) GetByUserID(userID UserID) ([]*Reservation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var userReservations []*Reservation
	for _, res := range r.reservations {
		if res.UserID == userID {
			resCopy := *res
			userReservations = append(userReservations, &resCopy)
		}
	}
	return userReservations, nil
}

type InMemoryUserRepository struct {
	mu    sync.RWMutex
	users map[UserID]*User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[UserID]*User),
	}
}

func (r *InMemoryUserRepository) GetByID(id UserID) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	userCopy := *user
	return &userCopy, nil
}

func (r *InMemoryUserRepository) Save(user *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.users[user.id] = user
	return nil
}

// Mock implementations of services
type MockPaymentService struct{}

func (s *MockPaymentService) ProcessRefund(userID UserID, amount Money) error {
	fmt.Printf("💰 Processing refund of %.2f %s to user %s\n", 
		float64(amount.Amount)/100, amount.Currency, userID)
	// Simulate processing time
	time.Sleep(50 * time.Millisecond)
	return nil
}

type MockNotificationService struct{}

func (s *MockNotificationService) NotifyCancellation(userID UserID, result *CancellationResult) error {
	refundText := "without refund"
	if result.RefundAmount.Amount > 0 {
		refundText = fmt.Sprintf("with refund of %.2f %s", 
			float64(result.RefundAmount.Amount)/100, result.RefundAmount.Currency)
	}
	fmt.Printf("📧 Notifying user %s: Reservation %s cancelled by %s %s\n", 
		userID, result.ReservationID, result.CancelledBy.GetID(), refundText)
	return nil
}

// Helper function to create test data
func setupTestData(userRepo *InMemoryUserRepository) {
	// Create users
	endUser1 := &User{
		id:   "user-123",
		role: RoleEndUser,
	}
	endUser2 := &User{
		id:   "user-456",
		role: RoleEndUser,
	}
	adminUser := &User{
		id:   "admin-001",
		role: RoleAdmin,
	}

	userRepo.Save(endUser1)
	userRepo.Save(endUser2)
	userRepo.Save(adminUser)
}

func printSeparator() {
	fmt.Println("\n" + strings.Repeat("-", 60) + "\n")
}

func main() {
	// Initialize repositories
	userRepo := NewInMemoryUserRepository()
	reservationRepo := NewInMemoryReservationRepository()

	// Initialize services
	paymentService := &MockPaymentService{}
	notificationService := &MockNotificationService{}
	clock := RealClock{}

	// Setup test data
	setupTestData(userRepo)

	// Create application service
	service := NewReservationService(
		reservationRepo,
		userRepo,
		paymentService,
		notificationService,
		clock,
	)

	fmt.Println("=== Reservation Cancellation System Demo ===\n")

	// Create some reservations first
	fmt.Println("📝 Creating reservations...")
	
	res1, err := service.CreateReservation(CreateReservationCommand{
		UserID:   "user-123",
		Amount:   15000, // $150.00
		Currency: "USD",
	})
	if err != nil {
		log.Fatalf("Failed to create reservation: %v", err)
	}
	fmt.Printf("Created reservation %s for user-123 ($150.00)\n", res1.id)

	res2, err := service.CreateReservation(CreateReservationCommand{
		UserID:   "user-456",
		Amount:   25000, // $250.00
		Currency: "USD",
	})
	if err != nil {
		log.Fatalf("Failed to create reservation: %v", err)
	}
	fmt.Printf("Created reservation %s for user-456 ($250.00)\n", res2.id)

	res3, err := service.CreateReservation(CreateReservationCommand{
		UserID:   "user-123",
		Amount:   10000, // $100.00
		Currency: "USD",
	})
	if err != nil {
		log.Fatalf("Failed to create reservation: %v", err)
	}
	fmt.Printf("Created reservation %s for user-123 ($100.00)\n", res3.id)

	printSeparator()

	// Example 1: End user cancels their own reservation (always with refund)
	fmt.Println("🟢 Example 1: End user cancels their own reservation")
	fmt.Println("Expected: Success with refund\n")
	
	result1, err := service.CancelReservation(CancelReservationCommand{
		ReservationID: string(res1.id),
		CancellerID:   "user-123",
		shouldRefund: true, // This is ignored for end users
	})
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Printf("✅ Success! Refund amount: $%.2f\n", float64(result1.RefundAmount.Amount)/100)
	}

	printSeparator()

	// // Example 2: End user tries to cancel someone else's reservation
	// fmt.Println("🔴 Example 2: End user tries to cancel another user's reservation")
	// fmt.Println("Expected: Error - not the owner")
	
	// _, err2 := service.CancelReservation(CancelReservationCommand{
	// 	ReservationID: string(res2.id),
	// 	CancellerID:   "user-123",
	// 	shouldRefund: true,
	// })
	// if err2 != nil {
	// 	fmt.Printf("❌ Error (as expected): %v\n", err2)
	// } else {
	// 	fmt.Println("✅ Unexpected success!")
	// }

	// printSeparator()

	// // Example 3: Admin cancels with refund
	// fmt.Println("🟢 Example 3: Admin cancels user's reservation WITH refund")
	// fmt.Println("Expected: Success with refund")
	
	// result3, err := service.CancelReservation(CancelReservationCommand{
	// 	ReservationID: string(res2.id),
	// 	CancellerID:   "admin-001",
	// 	shouldRefund: true, // Admin wants to give refund
	// })
	// if err != nil {
	// 	fmt.Printf("❌ Error: %v\n", err)
	// } else {
	// 	fmt.Printf("✅ Success! Refund amount: $%.2f\n", float64(result3.RefundAmount.Amount)/100)
	// }

	// printSeparator()

	// // Example 4: Admin cancels WITHOUT refund
	// fmt.Println("🟠 Example 4: Admin cancels user's reservation WITHOUT refund")
	// fmt.Println("Expected: Success with no refund")
	
	// result4, err := service.CancelReservation(CancelReservationCommand{
	// 	ReservationID: string(res3.id),
	// 	CancellerID:   "admin-001",
	// 	shouldRefund: false, // Admin forces no refund
	// })
	// if err != nil {
	// 	fmt.Printf("❌ Error: %v\n", err)
	// } else {
	// 	fmt.Printf("✅ Success! Refund amount: $%.2f (no refund as requested)\n", 
	// 		float64(result4.RefundAmount.Amount)/100)
	// }

	// printSeparator()

	// // Example 5: Try to cancel already cancelled reservation
	// fmt.Println("🔴 Example 5: Try to cancel an already cancelled reservation")
	// fmt.Println("Expected: Error - already cancelled\n")
	
	// _, err5 := service.CancelReservation(CancelReservationCommand{
	// 	ReservationID: string(res1.id),
	// 	CancellerID:   "user-123",
	// 	shouldRefund: true,
	// })
	// if err5 != nil {
	// 	fmt.Printf("❌ Error (as expected): %v\n", err5)
	// } else {
	// 	fmt.Println("✅ Unexpected success!")
	// }

	// printSeparator()

	// // Show final state of all reservations
	// fmt.Println("📊 Final State of All Reservations:\n")
	
	// for _, resID := range []string{string(res1.id), string(res2.id), string(res3.id)} {
	// 	details, err := service.GetReservationDetails(resID)
	// 	if err != nil {
	// 		fmt.Printf("Error getting details for %s: %v\n", resID, err)
	// 		continue
	// 	}
		
	// 	fmt.Printf("Reservation %s:\n", details.ID)
	// 	fmt.Printf("  User: %s\n", details.UserID)
	// 	fmt.Printf("  Status: %s\n", details.Status)
	// 	fmt.Printf("  Amount: $%.2f %s\n", float64(details.Amount)/100, details.Currency)
	// 	fmt.Printf("  Created: %s\n", details.CreatedAt)
	// 	if details.CancelledAt != "" {
	// 		fmt.Printf("  Cancelled: %s by %s\n", details.CancelledAt, details.CancelledBy)
	// 	}
	// 	fmt.Println()
	// }
}

// Import for string operations
var strings = struct {
	Repeat func(string, int) string
}{
	Repeat: func(s string, count int) string {
		result := ""
		for range count {
			result += s
		}
		return result
	},
}