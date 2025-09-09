package hotel_service_test

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/tktanisha/booking_system/internal/api/validators/payloads"
	"github.com/tktanisha/booking_system/internal/mocks"
	"github.com/tktanisha/booking_system/internal/models"
	"github.com/tktanisha/booking_system/internal/services/hotel_service"
)

func TestHotelService_GetHotelByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockHotelRepositoryInterface(ctrl)
	service := hotel_service.NewHotelService(mockRepo)

	hotelID := uuid.New()
	expectedHotel := &models.Hotels{Id: hotelID, Name: "Test Hotel"}

	mockRepo.EXPECT().
		GetHotelByID(hotelID).
		Return(expectedHotel, nil)

	result, err := service.GetHotelByID(hotelID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name != "Test Hotel" {
		t.Errorf("expected Test Hotel, got %v", result.Name)
	}
}

func TestHotelService_CreateHotel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockHotelRepositoryInterface(ctrl)
	service := hotel_service.NewHotelService(mockRepo)

	ctx := &models.UserContext{Id: uuid.New()}
	payload := &payloads.CreateHotelPayload{Name: "New Hotel", Address: "123 Street"}

	mockRepo.EXPECT().
		CreateHotel(gomock.Any()).
		DoAndReturn(func(hotel *models.Hotels) (*models.Hotels, error) {
			hotel.CreatedAt = time.Now()
			return hotel, nil
		})

	result, err := service.CreateHotel(ctx, payload)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name != "New Hotel" {
		t.Errorf("expected New Hotel, got %v", result.Name)
	}
}
