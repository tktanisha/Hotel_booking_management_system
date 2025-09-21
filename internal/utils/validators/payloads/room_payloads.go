package payloads

import (
	"github.com/google/uuid"
	"github.com/tktanisha/booking_system/internal/enums/room"
)

type CreateRoomPayload struct {
	HotelID  uuid.UUID     `json:"hotel_id"`
	RoomType room.RoomType `json:"room_type"`
	Price    float64       `json:"price"`
	Quantity int           `json:"quantity"`
}
type RoomPayload struct {
	RoomType room.RoomType `json:"room_type"`
	Quantity int           `json:"quantity"`
}
