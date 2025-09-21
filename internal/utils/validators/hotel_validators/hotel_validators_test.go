package hotel_validators_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/tktanisha/booking_system/internal/utils/validators/hotel_validators"
	"github.com/tktanisha/booking_system/internal/utils/validators/payloads"
)

func TestValidateHotelPayload(t *testing.T) {
	validPayload := payloads.CreateHotelPayload{
		Name:    "Grand Hotel",
		Address: "123 Main Street, City Center",
	}

	tests := []struct {
		name        string
		body        interface{}
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid payload",
			body:        validPayload,
			expectError: false,
		},
		{
			name:        "invalid JSON",
			body:        "{invalid json",
			expectError: true,
			errorMsg:    "invalid request payload",
		},
		{
			name: "missing name",
			body: payloads.CreateHotelPayload{
				Name:    "",
				Address: "123 Main Street, City Center",
			},
			expectError: true,
			errorMsg:    "name is required",
		},
		{
			name: "short name",
			body: payloads.CreateHotelPayload{
				Name:    "AB",
				Address: "123 Main Street, City Center",
			},
			expectError: true,
			errorMsg:    "name must be between 3 and 100 characters",
		},
		{
			name: "long name",
			body: payloads.CreateHotelPayload{
				Name:    string(make([]byte, 101)),
				Address: "123 Main Street, City Center",
			},
			expectError: true,
			errorMsg:    "name must be between 3 and 100 characters",
		},
		{
			name: "missing address",
			body: payloads.CreateHotelPayload{
				Name:    "Grand Hotel",
				Address: "",
			},
			expectError: true,
			errorMsg:    "address is required",
		},
		{
			name: "short address",
			body: payloads.CreateHotelPayload{
				Name:    "Grand Hotel",
				Address: "short",
			},
			expectError: true,
			errorMsg:    "address must be between 10 and 200 characters",
		},
		{
			name: "long address",
			body: payloads.CreateHotelPayload{
				Name:    "Grand Hotel",
				Address: string(make([]byte, 201)),
			},
			expectError: true,
			errorMsg:    "address must be between 10 and 200 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			switch v := tt.body.(type) {
			case string:
				body = []byte(v)
			default:
				body, _ = json.Marshal(v)
			}

			req := httptest.NewRequest("POST", "/", bytes.NewBuffer(body))
			_, err := hotel_validators.ValidateHotelPayload(req)

			if tt.expectError {
				if err == nil || err.Error() != tt.errorMsg {
					t.Errorf("expected error %q, got %v", tt.errorMsg, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}
