package hotel_service

import (
	"time"

	"github.com/google/uuid"
	"github.com/tktanisha/booking_system/internal/models"
	"github.com/tktanisha/booking_system/internal/repository/hotel_repo"
	"github.com/tktanisha/booking_system/internal/utils/validators/payloads"
)

type HotelService struct {
	hotelRepo hotel_repo.HotelRepositoryInterface
}

func NewHotelService(hotelRepo hotel_repo.HotelRepositoryInterface) *HotelService {
	return &HotelService{
		hotelRepo: hotelRepo,
	}
}

func (h *HotelService) GetHotelByID(hotelID uuid.UUID) (*models.Hotels, error) {
	return h.hotelRepo.GetHotelByID(hotelID)
}

func (h *HotelService) CreateHotel(ctx *models.UserContext, payload *payloads.CreateHotelPayload) (*models.Hotels, error) {
	hotel := &models.Hotels{
		Id:        uuid.New(),
		Name:      payload.Name,
		Address:   payload.Address,
		CreatedAt: time.Now(),
		ManagerId: ctx.Id,
	}
	return h.hotelRepo.CreateHotel(hotel)
}
