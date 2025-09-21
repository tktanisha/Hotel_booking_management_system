package booking_validators_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/tktanisha/booking_system/internal/enums/room"
	"github.com/tktanisha/booking_system/internal/utils/validators/booking_validators"
	"github.com/tktanisha/booking_system/internal/utils/validators/payloads"
)

func TestCreateBookingValidator(t *testing.T) {
	validHotelID := uuid.New()
	now := time.Now()
	later := now.Add(24 * time.Hour)

	tests := []struct {
		name        string
		payload     payloads.BookingPayload
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid payload",
			payload: payloads.BookingPayload{
				HotelId:  validHotelID,
				CheckIn:  now,
				CheckOut: later,
				Rooms: []*payloads.RoomPayload{
					{RoomType: room.Single, Quantity: 2},
				},
			},
			expectError: false,
		},
		{
			name:        "invalid JSON",
			payload:     payloads.BookingPayload{}, // we'll send invalid JSON separately
			expectError: true,
			errorMsg:    "invalid character",
		},
		{
			name: "missing hotel_id",
			payload: payloads.BookingPayload{
				CheckIn:  now,
				CheckOut: later,
				Rooms:    []*payloads.RoomPayload{{RoomType: room.Single, Quantity: 1}},
			},
			expectError: true,
			errorMsg:    "hotel_id is required",
		},
		{
			name: "checkin after checkout",
			payload: payloads.BookingPayload{
				HotelId:  validHotelID,
				CheckIn:  later,
				CheckOut: now,
				Rooms:    []*payloads.RoomPayload{{RoomType: room.Single, Quantity: 1}},
			},
			expectError: true,
			errorMsg:    "checkin date must be before checkout date",
		},
		{
			name: "no rooms",
			payload: payloads.BookingPayload{
				HotelId:  validHotelID,
				CheckIn:  now,
				CheckOut: later,
				Rooms:    []*payloads.RoomPayload{},
			},
			expectError: true,
			errorMsg:    "at least one room is required",
		},
		{
			name: "room with empty type",
			payload: payloads.BookingPayload{
				HotelId:  validHotelID,
				CheckIn:  now,
				CheckOut: later,
				Rooms:    []*payloads.RoomPayload{{RoomType: "", Quantity: 1}},
			},
			expectError: true,
			errorMsg:    "room_type is required",
		},
		{
			name: "room with non-positive quantity",
			payload: payloads.BookingPayload{
				HotelId:  validHotelID,
				CheckIn:  now,
				CheckOut: later,
				Rooms:    []*payloads.RoomPayload{{RoomType: room.Single, Quantity: 0}},
			},
			expectError: true,
			errorMsg:    "quantity must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.name == "invalid JSON" {
				req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte("{invalid json")))
			} else {
				body, _ := json.Marshal(tt.payload)
				req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
			}

			got, err := booking_validators.CreateBookingValidator(req)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got nil")
				} else if !contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error containing %q, got %q", tt.errorMsg, err.Error())
				}
				if got != nil {
					t.Errorf("expected nil payload, got %v", got)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if got == nil {
					t.Errorf("expected payload, got nil")
				} else if got.HotelId != tt.payload.HotelId {
					t.Errorf("expected hotel_id %v, got %v", tt.payload.HotelId, got.HotelId)
				}
			}
		})
	}
}

func TestValidateCancelBooking(t *testing.T) {
	tests := []struct {
		name        string
		bookingID   uuid.UUID
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid booking ID",
			bookingID:   uuid.New(),
			expectError: false,
		},
		{
			name:        "missing booking ID (Nil)",
			bookingID:   uuid.Nil,
			expectError: true,
			errorMsg:    "booking_id is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := booking_validators.ValidateCancelBooking(tt.bookingID)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got nil")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("expected error message %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}
		})
	}
}

// helper function to match substrings
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || bytes.Contains([]byte(s), []byte(substr)))
}
