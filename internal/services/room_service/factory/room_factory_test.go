package factory_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/tktanisha/booking_system/internal/enums/room"
	"github.com/tktanisha/booking_system/internal/services/room_service/factory"
	"github.com/tktanisha/booking_system/internal/utils/validators/payloads"
)

func TestGetRoomFactory(t *testing.T) {
	tests := []struct {
		name        string
		roomType    room.RoomType
		expectError bool
	}{
		{"single room factory", room.Single, false},
		{"double room factory", room.Double, false},
		{"suite room factory", room.Suite, false},
		{"invalid room type", "dhfkjs", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := factory.GetRoomFactory(tt.roomType)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				if f != nil {
					t.Errorf("expected nil factory, got %+v", f)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if f == nil {
					t.Errorf("expected non-nil factory, got nil")
				}
			}
		})
	}
}

func TestRoomFactories_Create(t *testing.T) {
	hotelID := uuid.New()
	quantity := 5
	payload := &payloads.CreateRoomPayload{
		HotelID:  hotelID,
		Quantity: quantity,
	}

	factories := []struct {
		name     string
		factory  factory.RoomFactory
		roomType room.RoomType
	}{
		{"single room", &factory.SingleRoomFactory{}, room.Single},
		{"double room", &factory.DoubleRoomFactory{}, room.Double},
		{"suite room", &factory.SuiteRoomFactory{}, room.Suite},
	}

	for _, tt := range factories {
		t.Run(tt.name, func(t *testing.T) {
			roomObj := tt.factory.Create(payload)

			if roomObj.HotelId != hotelID {
				t.Errorf("expected HotelId=%v, got %v", hotelID, roomObj.HotelId)
			}
			if roomObj.AvailableQuantity != quantity {
				t.Errorf("expected Quantity=%d, got %d", quantity, roomObj.AvailableQuantity)
			}
			if roomObj.RoomCategory != tt.roomType {
				t.Errorf("expected RoomCategory=%v, got %v", tt.roomType, roomObj.RoomCategory)
			}
			if roomObj.Id == uuid.Nil {
				t.Errorf("expected non-nil Id, got %v", roomObj.Id)
			}
			if time.Since(roomObj.CreatedAt) > time.Second {
				t.Errorf("CreatedAt timestamp is too old: %v", roomObj.CreatedAt)
			}
		})
	}
}
