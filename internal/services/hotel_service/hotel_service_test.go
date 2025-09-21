package hotel_service_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/tktanisha/booking_system/internal/mocks"
	"github.com/tktanisha/booking_system/internal/models"
	"github.com/tktanisha/booking_system/internal/services/hotel_service"
	"github.com/tktanisha/booking_system/internal/utils/validators/payloads"
)

func TestHotelService_GetHotelByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockHotelRepositoryInterface(ctrl)
	service := hotel_service.NewHotelService(mockRepo)

	HotelID := uuid.New()

	test := []struct {
		name     string
		hotelId  uuid.UUID
		mockFunc func()
		wantErr  bool
	}{
		{
			name:    "unable to get the hotel",
			hotelId: HotelID,
			mockFunc: func() {
				mockRepo.EXPECT().GetHotelByID(HotelID).Return(nil, errors.New("not found"))
			},
			wantErr: true,
		},
		{
			name:    "succesfully fetched the hotel",
			hotelId: HotelID,
			mockFunc: func() {
				mockRepo.EXPECT().GetHotelByID(HotelID).Return(&models.Hotels{}, nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			_, err := service.GetHotelByID(tt.hotelId)

			if (err != nil) != tt.wantErr {
				t.Errorf("want = %v vand got =%v", tt.wantErr, err)
			}
		})
	}

}

func TestHotelService_CreateHotel(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mocks.NewMockHotelRepositoryInterface(ctrl)
	service := hotel_service.NewHotelService(mockRepo)

	managerCtx := &models.UserContext{
		Id:   uuid.New(),
		Role: "manager",
	}

	tests := []struct {
		name     string
		userCtx  *models.UserContext
		payload  *payloads.CreateHotelPayload
		mockFunc func()
		wantErr  bool
	}{
		{
			name:    "Successful hotel creation",
			userCtx: managerCtx,
			payload: &payloads.CreateHotelPayload{
				Name:    "Test Hotel",
				Address: "123 Test St",
			},
			mockFunc: func() {
				mockRepo.EXPECT().
					CreateHotel(gomock.Any()).
					Return(&models.Hotels{Id: uuid.New(), Name: "Test Hotel", Address: "123 Test St", ManagerId: managerCtx.Id}, nil)
			},
			wantErr: false,
		},
		{
			name:    "Repository error during creation",
			userCtx: managerCtx,
			payload: &payloads.CreateHotelPayload{
				Name:    "Error Hotel",
				Address: "456 Error Rd",
			},
			mockFunc: func() {
				mockRepo.EXPECT().
					CreateHotel(gomock.Any()).
					Return(nil, errors.New("repository error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			_, err := service.CreateHotel(tt.userCtx, tt.payload)

			if (err != nil) != tt.wantErr {
				t.Errorf("want error: %v, but got: %v", tt.wantErr, err)

			}
		})

	}
}
