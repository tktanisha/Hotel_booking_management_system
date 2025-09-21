package room_service

import (
	"github.com/google/uuid"
	"github.com/tktanisha/booking_system/internal/models"
	"github.com/tktanisha/booking_system/internal/utils/validators/payloads"
)

//go:generate mockgen -source=room_interface.go -destination=../../mocks/mock_room_service.go -package=mocks

type RoomServiceInterface interface {
	CreateRoom(*payloads.CreateRoomPayload) (*models.Rooms, error)
	IsAvailable(*payloads.RoomPayload, uuid.UUID) bool
	IncreaseRoomQuantity(*payloads.RoomPayload, uuid.UUID) (*models.Rooms, error)
	ReduceRoomQuantity(*payloads.RoomPayload, uuid.UUID) error
	GetAllRoomByHotelID(hotelID uuid.UUID) ([]*models.Rooms, error) //it also show how much room are available
}
