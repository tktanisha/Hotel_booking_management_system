package auth_validators

import (
	"encoding/json"
	"net/http"

	"github.com/tktanisha/booking_system/internal/api/validators/payloads"
)

func RegisterValidate(r *http.Request) (*payloads.RegisterRequest, error) {
	var payload payloads.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return nil, err
	}
	if err := ValidateFullName(payload.Fullname); err != nil {
		return nil, err
	}
	if err := ValidateEmail(payload.Email); err != nil {
		return nil, err
	}
	if err := ValidatePassword(payload.Password); err != nil {
		return nil, err
	}
	return &payload, nil
}
