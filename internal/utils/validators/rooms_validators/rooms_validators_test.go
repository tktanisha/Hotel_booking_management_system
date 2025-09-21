package room_validators_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/tktanisha/booking_system/internal/enums/room"
	"github.com/tktanisha/booking_system/internal/utils/validators/payloads"
	room_validators "github.com/tktanisha/booking_system/internal/utils/validators/rooms_validators"
)

func TestValidateCreateRoomPayload(t *testing.T) {
	validPayload := payloads.CreateRoomPayload{
		HotelID:  uuid.New(),
		RoomType: room.Single,
		Price:    100,
		Quantity: 2,
	}

	tests := []struct {
		name        string
		body        interface{}
		expectError bool
		errorMsg    string
	}{
		{"valid payload", validPayload, false, ""},
		{"invalid JSON", "{invalid json", true, "invalid request payload"},
		{"missing hotel_id", payloads.CreateRoomPayload{
			HotelID:  uuid.Nil,
			RoomType: room.Single,
			Price:    100,
			Quantity: 2,
		}, true, "hotel_id is required"},
		{"invalid room_type", payloads.CreateRoomPayload{
			HotelID:  uuid.New(),
			RoomType: "invalid",
			Price:    100,
			Quantity: 2,
		}, true, "invalid room_type"},
		{"non-positive price", payloads.CreateRoomPayload{
			HotelID:  uuid.New(),
			RoomType: room.Double,
			Price:    0,
			Quantity: 2,
		}, true, "price must be positive"},
		{"non-positive quantity", payloads.CreateRoomPayload{
			HotelID:  uuid.New(),
			RoomType: room.Suite,
			Price:    100,
			Quantity: 0,
		}, true, "quantity must be positive"},
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
			_, err := room_validators.ValidateCreateRoomPayload(req)

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

func TestValidateRoomPayload(t *testing.T) {
	validPayload := []*payloads.RoomPayload{
		{
			RoomType: room.Double,
			Quantity: 2,
		},
	}

	tests := []struct {
		name        string
		body        interface{}
		expectError bool
		errorMsg    string
	}{
		{"valid rooms payload", validPayload, false, ""},
		{"invalid JSON", "{invalid json", true, "invalid character"},
		{"negative quantity", []*payloads.RoomPayload{
			{
				RoomType: room.Single,
				Quantity: -1,
			},
		}, true, "quantity must be positive"},
		{"invalid room type", []*payloads.RoomPayload{
			{
				RoomType: "invalid",
				Quantity: 2,
			},
		}, true, "invalid room type"},
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
			_, err := room_validators.ValidateRoomPayload(req)

			if tt.expectError {
				if err == nil || !contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error containing %q, got %v", tt.errorMsg, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || bytes.Contains([]byte(s), []byte(substr)))
}
