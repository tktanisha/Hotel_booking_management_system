package factory

import (
	"github.com/tktanisha/booking_system/internal/api/validators/payloads"
	"github.com/tktanisha/booking_system/internal/models"
)

type RoomFactory interface {
	Create(payload *payloads.CreateRoomPayload) *models.Rooms
}
