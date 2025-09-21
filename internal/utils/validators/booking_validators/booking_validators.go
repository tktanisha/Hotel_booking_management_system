package booking_validators

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/tktanisha/booking_system/internal/utils/validators/payloads"
)

func CreateBookingValidator(r *http.Request) (*payloads.BookingPayload, error) {
	var payload payloads.BookingPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return nil, err
	}

	if payload.HotelId == uuid.Nil {
		return nil, errors.New("hotel_id is required")
	}
	if payload.CheckIn.IsZero() {
		return nil, errors.New("checkin date is required")
	}
	if payload.CheckOut.IsZero() {
		return nil, errors.New("checkout date is required")
	}
	if payload.CheckIn.After(payload.CheckOut) {
		return nil, errors.New("checkin date must be before checkout date")
	}
	if len(payload.Rooms) == 0 {
		return nil, errors.New("at least one room is required")
	}
	for _, room := range payload.Rooms {
		if room.RoomType == "" {
			return nil, errors.New("room_type is required")
		}
		if room.Quantity < 0 {
			return nil, errors.New("quantity must be positive")
		}
	}
	return &payload, nil
}
