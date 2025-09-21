package hotel_validators

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/tktanisha/booking_system/internal/utils/validators/payloads"
)

func ValidateHotelPayload(r *http.Request) (*payloads.CreateHotelPayload, error) {
	var payload payloads.CreateHotelPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return nil, errors.New("invalid request payload")
	}

	if payload.Name == "" {
		return nil, errors.New("name is required")
	}

	if len(payload.Name) < 3 || len(payload.Name) > 100 {
		return nil, errors.New("name must be between 3 and 100 characters")
	}

	if payload.Address == "" {
		return nil, errors.New("address is required")
	}

	if len(payload.Address) < 10 || len(payload.Address) > 200 {
		return nil, errors.New("address must be between 10 and 200 characters")
	}

	return &payload, nil
}
