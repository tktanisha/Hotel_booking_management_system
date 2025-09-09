package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/tktanisha/booking_system/internal/enums/room"
)

type Rooms struct {
	Id                uuid.UUID     `json:"id"`
	HotelId           uuid.UUID     `json:"hotel_id"`
	AvailableQuantity int           `json:"available_quantity"`
	RoomCategory      room.RoomType `json:"room_category"`
	CreatedAt         time.Time     `json:"created_at"`
}
