package booking_service

import (
	"github.com/google/uuid"

	"github.com/tktanisha/booking_system/internal/api/validators/payloads"
	"github.com/tktanisha/booking_system/internal/models"
)

//go:generate mockgen -source=booking_service_interface.go -destination=../../mocks/mock_booking_service.go -package=mocks

type BookingServiceInterface interface {
	CreateBooking(*models.UserContext, *payloads.BookingPayload) (*models.Bookings, error)
	CancelBooking(uuid.UUID) (*models.Bookings, error)
	CheckoutBooking(uuid.UUID) (*models.Bookings, error)
}
