package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/tktanisha/booking_system/internal/enums/room"
)

type BookedRooms struct {
	Id           uuid.UUID     `json:"id"`
	BookingId    uuid.UUID     `json:"booking_id"`
	RoomType     room.RoomType `json:"room_type"`
	RoomQuantity int           `json:"room_quantity"`
	CreatedAt    time.Time     `json:"created_at"`
}
