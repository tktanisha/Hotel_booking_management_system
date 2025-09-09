package factory

import (
	"time"

	"github.com/google/uuid"
	"github.com/tktanisha/booking_system/internal/api/validators/payloads"
	"github.com/tktanisha/booking_system/internal/enums/room"
	"github.com/tktanisha/booking_system/internal/models"
)

type SingleRoomFactory struct{}

func (s *SingleRoomFactory) Create(payload *payloads.CreateRoomPayload) *models.Rooms {
	return &models.Rooms{
		Id:                uuid.New(),
		HotelId:           payload.HotelID,
		AvailableQuantity: payload.Quantity,
		RoomCategory:      room.Single,
		CreatedAt:         time.Now(),
	}
}
