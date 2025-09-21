package room_validators

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/tktanisha/booking_system/internal/enums/room"
	"github.com/tktanisha/booking_system/internal/utils/validators/payloads"
)

func ValidateCreateRoomPayload(r *http.Request) (*payloads.CreateRoomPayload, error) {
	var payload payloads.CreateRoomPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return nil, errors.New("invalid request payload")
	}

	if payload.HotelID == uuid.Nil {
		return nil, errors.New("hotel_id is required")
	}

	validRoomTypes := map[room.RoomType]bool{
		room.Single: true,
		room.Double: true,
		room.Suite:  true,
	}
	if !validRoomTypes[payload.RoomType] {
		return nil, errors.New("invalid room_type")
	}

	if payload.Price <= 0 {
		return nil, errors.New("price must be positive")
	}
	if payload.Quantity <= 0 {
		return nil, errors.New("quantity must be positive")
	}
	return &payload, nil
}
