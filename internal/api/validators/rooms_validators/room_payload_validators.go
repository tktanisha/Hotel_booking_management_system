package room_validators

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/tktanisha/booking_system/internal/api/validators/payloads"
	"github.com/tktanisha/booking_system/internal/enums/room"
)

func ValidateRoomPayload(r *http.Request) ([]*payloads.RoomPayload, error) {
	var payload []*payloads.RoomPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return nil, err
	}

	validRoomTypes := map[room.RoomType]bool{
		room.Single: true,
		room.Double: true,
		room.Suite:  true,
	}

	for _, payload := range payload {
		if payload.Quantity <= 0 {
			return nil, errors.New("quantity must be positive")
		}

		if !validRoomTypes[payload.RoomType] {
			return nil, errors.New("invalid room type")
		}
	}

	return payload, nil
}
