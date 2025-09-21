package factory

import (
	"time"

	"github.com/google/uuid"
	"github.com/tktanisha/booking_system/internal/enums/room"
	"github.com/tktanisha/booking_system/internal/models"
	"github.com/tktanisha/booking_system/internal/utils/validators/payloads"
)

type DoubleRoomFactory struct{}

func (d *DoubleRoomFactory) Create(payload *payloads.CreateRoomPayload) *models.Rooms {
	return &models.Rooms{
		Id:                uuid.New(),
		HotelId:           payload.HotelID,
		AvailableQuantity: payload.Quantity,
		RoomCategory:      room.Double,
		CreatedAt:         time.Now(),
	}
}
