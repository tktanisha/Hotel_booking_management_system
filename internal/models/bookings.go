package models

import (
	"time"

	"github.com/google/uuid"
	booking_status "github.com/tktanisha/booking_system/internal/enums/booking"
)

type Bookings struct {
	Id        uuid.UUID                    `json:"id"`
	UserId    uuid.UUID                    `json:"user_id"`
	HotelId   uuid.UUID                    `json:"hotel_id"`
	CheckIn   time.Time                    `json:"checkin"`
	CheckOut  time.Time                    `json:"checkout"`
	Status    booking_status.BookingStatus `json:"status"`
	CreatedAt time.Time                    `json:"created_at"`
}
