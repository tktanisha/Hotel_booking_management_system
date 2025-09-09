package hotel_service

import (
	"github.com/google/uuid"
	"github.com/tktanisha/booking_system/internal/api/validators/payloads"
	"github.com/tktanisha/booking_system/internal/models"
)

//go:generate mockgen -source=hotel_interface.go -destination=../../mocks/mock_hotel_service.go -package=mocks

type HotelServiceInterface interface {
	GetHotelByID(uuid.UUID) (*models.Hotels, error)
	CreateHotel(*models.UserContext, *payloads.CreateHotelPayload) (*models.Hotels, error)
}
