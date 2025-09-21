package payloads

import (
	"time"

	"github.com/google/uuid"
)

type BookingPayload struct {
	HotelId  uuid.UUID      `json:"hotel_id"`
	CheckIn  time.Time      `json:"checkin"`
	CheckOut time.Time      `json:"checkout"`
	Rooms    []*RoomPayload `json:"rooms"`
}
