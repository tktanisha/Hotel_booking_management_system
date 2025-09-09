package hotel_repo

import (
	"github.com/google/uuid"
	"github.com/tktanisha/booking_system/internal/models"
)

//go:generate mockgen -source=hotel_interface.go -destination=../../mocks/mock_hotel_repo.go -package=mocks

type HotelRepositoryInterface interface {
	GetHotelByID(uuid.UUID) (*models.Hotels, error)
	CreateHotel(*models.Hotels) (*models.Hotels, error)
}
