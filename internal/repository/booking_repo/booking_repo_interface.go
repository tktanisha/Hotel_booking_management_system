package booking_repo

import (
	"github.com/google/uuid"
	"github.com/tktanisha/booking_system/internal/models"
)

//go:generate mockgen -source=booking_repo_interface.go -destination=../../mocks/mock_booking_repo.go -package=mocks

type BookingRepoInterface interface {
	CreateBookingWithRooms(*models.Bookings, []*models.BookedRooms) (*models.Bookings, error)
	GetBookingById(uuid.UUID) (*models.Bookings, error)
	GetBookedRoomsByBookingId(uuid.UUID) ([]*models.BookedRooms, error)
	Save(*models.Bookings) error
}
