package factory

import (
	"github.com/tktanisha/booking_system/internal/models"
	"github.com/tktanisha/booking_system/internal/utils/validators/payloads"
)

type RoomFactory interface {
	Create(payload *payloads.CreateRoomPayload) *models.Rooms
}
