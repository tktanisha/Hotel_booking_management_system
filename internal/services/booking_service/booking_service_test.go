package booking_service_test

import (
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	booking_status "github.com/tktanisha/booking_system/internal/enums/booking"
	"github.com/tktanisha/booking_system/internal/enums/room"
	"github.com/tktanisha/booking_system/internal/mocks"
	"github.com/tktanisha/booking_system/internal/models"
	"github.com/tktanisha/booking_system/internal/services/booking_service"
	"github.com/tktanisha/booking_system/internal/utils/validators/payloads"
)

func TestBookingService_CancelBooking(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mocks.NewMockBookingRepoInterface(ctrl)
	mockRoomService := mocks.NewMockRoomServiceInterface(ctrl)
	service := booking_service.NewBookingService(mockBookingRepo, mockRoomService)

	bookingID := uuid.New()
	hotelID := uuid.New()

	tests := []struct {
		name        string
		mockSetup   func()
		expectError bool
	}{
		{
			name: "successfully cancelled",
			mockSetup: func() {
				mockBookingRepo.EXPECT().GetBookingById(bookingID).Return(
					&models.Bookings{
						Id:      bookingID,
						UserId:  uuid.New(),
						HotelId: hotelID,
						CheckIn: time.Now().Add(24 * time.Hour),
						Status:  booking_status.StatusConfirmed,
					}, nil)
				mockBookingRepo.EXPECT().GetBookedRoomsByBookingId(bookingID).Return(
					[]*models.BookedRooms{{RoomType: room.Single, RoomQuantity: 2}}, nil)
				mockRoomService.EXPECT().IncreaseRoomQuantity(gomock.Any(), hotelID).Return(
					&models.Rooms{}, nil)
				mockBookingRepo.EXPECT().Save(gomock.Any()).Return(nil)
			},
			expectError: false,
		},
		{
			name: "error fetching booking",
			mockSetup: func() {
				mockBookingRepo.EXPECT().GetBookingById(bookingID).Return(nil, errors.New("booking not found"))
			},
			expectError: true,
		},
		{
			name: "cannot cancel after check-in",
			mockSetup: func() {
				mockBookingRepo.EXPECT().GetBookingById(bookingID).Return(&models.Bookings{CheckIn: time.Now().Add(-2 * time.Hour)}, nil)
			},
			expectError: true,
		},
		{
			name: "booking already cancelled",
			mockSetup: func() {
				mockBookingRepo.EXPECT().GetBookingById(bookingID).Return(
					&models.Bookings{CheckIn: time.Now().Add(-2 * time.Hour), Status: booking_status.StatusCancelled}, nil)
			},
			expectError: true,
		},
		{
			name: "error fetching booked rooms",
			mockSetup: func() {
				mockBookingRepo.EXPECT().GetBookingById(bookingID).Return(
					&models.Bookings{CheckIn: time.Now().Add(24 * time.Hour), Status: booking_status.StatusConfirmed}, nil)

				mockBookingRepo.EXPECT().GetBookedRoomsByBookingId(bookingID).Return(
					nil, errors.New("db error"))
			},
			expectError: true,
		},
		{
			name: "error increasing room quantity",
			mockSetup: func() {
				mockBookingRepo.EXPECT().GetBookingById(bookingID).Return(
					&models.Bookings{
						Id:      bookingID,
						HotelId: hotelID,
						CheckIn: time.Now().Add(24 * time.Hour),
						Status:  booking_status.StatusConfirmed,
					}, nil)
				mockBookingRepo.EXPECT().GetBookedRoomsByBookingId(bookingID).Return(
					[]*models.BookedRooms{{RoomType: room.Single, RoomQuantity: 2}}, nil)
				mockRoomService.EXPECT().IncreaseRoomQuantity(gomock.Any(), hotelID).Return(
					nil, errors.New("unable to increase"))
			},
			expectError: true,
		},
		{
			name: "error saving booking",
			mockSetup: func() {
				mockBookingRepo.EXPECT().GetBookingById(bookingID).Return(
					&models.Bookings{
						Id:      bookingID,
						HotelId: hotelID,
						CheckIn: time.Now().Add(24 * time.Hour),
						Status:  booking_status.StatusConfirmed,
					}, nil)
				mockBookingRepo.EXPECT().GetBookedRoomsByBookingId(bookingID).Return(
					[]*models.BookedRooms{{RoomType: room.Single, RoomQuantity: 2}}, nil)
				mockRoomService.EXPECT().IncreaseRoomQuantity(gomock.Any(), hotelID).Return(
					&models.Rooms{}, nil)
				mockBookingRepo.EXPECT().Save(gomock.Any()).Return(errors.New("interanl server error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			_, err := service.CancelBooking(bookingID)
			if (err != nil) != tt.expectError {
				t.Errorf("expected error=%v, got=%v", tt.expectError, err)
			}
		})
	}
}

func TestBookingService_CreateBooking(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mocks.NewMockBookingRepoInterface(ctrl)
	mockRoomService := mocks.NewMockRoomServiceInterface(ctrl)
	service := booking_service.NewBookingService(mockBookingRepo, mockRoomService)

	userCtx := &models.UserContext{Id: uuid.New()}
	hotelID := uuid.New()
	roomPayload := &payloads.RoomPayload{RoomType: "Deluxe", Quantity: 2}
	payload := &payloads.BookingPayload{
		HotelId:  hotelID,
		CheckIn:  time.Now().Add(24 * time.Hour),
		CheckOut: time.Now().Add(48 * time.Hour),
		Rooms:    []*payloads.RoomPayload{roomPayload},
	}

	t.Run("success", func(t *testing.T) {
		mockRoomService.EXPECT().IsAvailable(roomPayload, hotelID).Return(true)
		mockRoomService.EXPECT().ReduceRoomQuantity(roomPayload, hotelID).Return(nil)
		mockBookingRepo.EXPECT().CreateBookingWithRooms(gomock.Any(), gomock.Any()).Return(&models.Bookings{Id: uuid.New()}, nil)

		_, err := service.CreateBooking(userCtx, payload)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("room not available", func(t *testing.T) {
		mockRoomService.EXPECT().IsAvailable(roomPayload, hotelID).Return(false)

		_, err := service.CreateBooking(userCtx, payload)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("reduce room quantity failure", func(t *testing.T) {
		mockRoomService.EXPECT().IsAvailable(roomPayload, hotelID).Return(true)
		mockRoomService.EXPECT().ReduceRoomQuantity(roomPayload, hotelID).Return(errors.New("reduce error"))

		_, err := service.CreateBooking(userCtx, payload)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("create booking failure", func(t *testing.T) {
		mockRoomService.EXPECT().IsAvailable(roomPayload, hotelID).Return(true)
		mockRoomService.EXPECT().ReduceRoomQuantity(roomPayload, hotelID).Return(nil)
		mockBookingRepo.EXPECT().CreateBookingWithRooms(gomock.Any(), gomock.Any()).Return(nil, errors.New("db error"))

		_, err := service.CreateBooking(userCtx, payload)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})
}

func TestBookingService_CheckoutBooking(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mocks.NewMockBookingRepoInterface(ctrl)
	mockRoomService := mocks.NewMockRoomServiceInterface(ctrl)
	service := booking_service.NewBookingService(mockBookingRepo, mockRoomService)

	bookingID := uuid.New()
	hotelID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockBookingRepo.EXPECT().GetBookingById(bookingID).Return(&models.Bookings{
			Id:      bookingID,
			HotelId: hotelID,
			Status:  booking_status.StatusConfirmed,
		}, nil)
		mockBookingRepo.EXPECT().GetBookedRoomsByBookingId(bookingID).Return([]*models.BookedRooms{
			{RoomType: "Suite", RoomQuantity: 1},
		}, nil)
		mockRoomService.EXPECT().IncreaseRoomQuantity(gomock.Any(), hotelID).Return(&models.Rooms{}, nil)
		mockBookingRepo.EXPECT().Save(gomock.Any()).Return(nil)

		_, err := service.CheckoutBooking(bookingID)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("error fetching booking", func(t *testing.T) {
		mockBookingRepo.EXPECT().GetBookingById(bookingID).Return(nil, errors.New("db error"))

		_, err := service.CheckoutBooking(bookingID)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("booking not confirmed", func(t *testing.T) {
		mockBookingRepo.EXPECT().GetBookingById(bookingID).Return(&models.Bookings{
			Id:     bookingID,
			Status: booking_status.StatusCancelled,
		}, nil)

		_, err := service.CheckoutBooking(bookingID)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("error fetching booked rooms", func(t *testing.T) {
		mockBookingRepo.EXPECT().GetBookingById(bookingID).Return(&models.Bookings{
			Id:      bookingID,
			HotelId: hotelID,
			Status:  booking_status.StatusConfirmed,
		}, nil)
		mockBookingRepo.EXPECT().GetBookedRoomsByBookingId(bookingID).Return(nil, errors.New("fetch error"))

		_, err := service.CheckoutBooking(bookingID)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("error increasing room quantity", func(t *testing.T) {
		mockBookingRepo.EXPECT().GetBookingById(bookingID).Return(&models.Bookings{
			Id:      bookingID,
			HotelId: hotelID,
			Status:  booking_status.StatusConfirmed,
		}, nil)
		mockBookingRepo.EXPECT().GetBookedRoomsByBookingId(bookingID).Return([]*models.BookedRooms{
			{RoomType: "Suite", RoomQuantity: 1},
		}, nil)
		mockRoomService.EXPECT().IncreaseRoomQuantity(gomock.Any(), hotelID).Return(nil, errors.New("increase error"))

		_, err := service.CheckoutBooking(bookingID)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("error saving booking", func(t *testing.T) {
		mockBookingRepo.EXPECT().GetBookingById(bookingID).Return(&models.Bookings{
			Id:      bookingID,
			HotelId: hotelID,
			Status:  booking_status.StatusConfirmed,
		}, nil)
		mockBookingRepo.EXPECT().GetBookedRoomsByBookingId(bookingID).Return([]*models.BookedRooms{
			{RoomType: "Suite", RoomQuantity: 1},
		}, nil)
		mockRoomService.EXPECT().IncreaseRoomQuantity(gomock.Any(), hotelID).Return(&models.Rooms{}, nil)
		mockBookingRepo.EXPECT().Save(gomock.Any()).Return(errors.New("save error"))

		_, err := service.CheckoutBooking(bookingID)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})
}
