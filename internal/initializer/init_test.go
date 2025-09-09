package initializer_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/tktanisha/booking_system/internal/initializer"
	"github.com/tktanisha/booking_system/internal/mocks"
)

func TestInitialize(t *testing.T) {
	// Create a mock DB
	sqlDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer sqlDB.Close()

	ctrl := gomock.NewController(t)
	mockDB := mocks.NewMockDB(ctrl)

	// Call the Initialize function
	initializer.Initialize(mockDB)

	// Validate that all global variables are initialized
	if initializer.AuthService == nil {
		t.Errorf("AuthService is nil")
	}
	if initializer.BookingService == nil {
		t.Errorf("BookingService is nil")
	}
	if initializer.HotelService == nil {
		t.Errorf("HotelService is nil")
	}
	if initializer.RoomService == nil {
		t.Errorf("RoomService is nil")
	}
}
