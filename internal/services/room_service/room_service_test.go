package room_service_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"

	"github.com/tktanisha/booking_system/internal/api/validators/payloads"
	"github.com/tktanisha/booking_system/internal/enums/room"
	"github.com/tktanisha/booking_system/internal/mocks"
	"github.com/tktanisha/booking_system/internal/models"
	"github.com/tktanisha/booking_system/internal/services/room_service"
)

// NOTE: adjust this import path if your mock package path differs.

func TestRoomService_CreateRoom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRoomRepoInterface(ctrl)
	svc := room_service.NewRoomService(mockRepo)

	t.Run("success create room", func(t *testing.T) {
		payload := &payloads.CreateRoomPayload{
			RoomType: room.Single,
			Price:    100,
			Quantity: 5,
		}

		// Expect CreateRoom to be called with any models.Rooms and return that same object
		mockRepo.
			EXPECT().
			CreateRoom(gomock.Any()).
			DoAndReturn(func(r *models.Rooms) (*models.Rooms, error) {
				// simulate DB returning same room with an ID
				r.Id = uuid.New()
				return r, nil
			})

		got, err := svc.CreateRoom(payload)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got == nil {
			t.Fatalf("expected a room, got nil")
		}
		if got.RoomCategory != payload.RoomType {
			t.Errorf("expected room type %s, got %s", payload.RoomType, got.RoomCategory)
		}
	})

	// factory error: use a room type that should make factory.GetRoomFactory return error
	t.Run("factory error", func(t *testing.T) {
		payload := &payloads.CreateRoomPayload{
			RoomType: "invalid-room-type", // factory should return error for this
			Price:    10,
			Quantity: 1,
		}

		// No repo expectations because service should fail before calling repo
		_, err := svc.CreateRoom(payload)
		if err == nil {
			t.Fatalf("expected error from factory, got nil")
		}
	})
}

func TestRoomService_IsAvailable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRoomRepoInterface(ctrl)
	svc := room_service.NewRoomService(mockRepo)

	hotelID := uuid.New()

	tests := []struct {
		name      string
		mockSetup func()
		roomReq   *payloads.RoomPayload
		want      bool
	}{
		{
			name: "repo error returns false",
			mockSetup: func() {
				mockRepo.EXPECT().GetAllRoomByHotelID(hotelID).Return(nil, errors.New("db error"))
			},
			roomReq: &payloads.RoomPayload{RoomType: room.Single, Quantity: 1},
			want:    false,
		},
		{
			name: "sufficient quantity returns true",
			mockSetup: func() {
				mockRepo.EXPECT().GetAllRoomByHotelID(hotelID).Return([]*models.Rooms{
					{RoomCategory: room.Single, AvailableQuantity: 5},
				}, nil)
			},
			roomReq: &payloads.RoomPayload{RoomType: room.Single, Quantity: 2},
			want:    true,
		},
		{
			name: "insufficient quantity returns false",
			mockSetup: func() {
				mockRepo.EXPECT().GetAllRoomByHotelID(hotelID).Return([]*models.Rooms{
					{RoomCategory: room.Single, AvailableQuantity: 1},
				}, nil)
			},
			roomReq: &payloads.RoomPayload{RoomType: room.Single, Quantity: 2},
			want:    false,
		},
		{
			name: "no matching room type returns false",
			mockSetup: func() {
				mockRepo.EXPECT().GetAllRoomByHotelID(hotelID).Return([]*models.Rooms{
					{RoomCategory: room.Double, AvailableQuantity: 10},
				}, nil)
			},
			roomReq: &payloads.RoomPayload{RoomType: room.Suite, Quantity: 1},
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			got := svc.IsAvailable(tt.roomReq, hotelID)
			if got != tt.want {
				t.Fatalf("IsAvailable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoomService_ReduceRoomQuantity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRoomRepoInterface(ctrl)
	svc := room_service.NewRoomService(mockRepo)

	hotelID := uuid.New()

	t.Run("repo GetAll error", func(t *testing.T) {
		mockRepo.EXPECT().GetAllRoomByHotelID(hotelID).Return(nil, errors.New("db"))
		err := svc.ReduceRoomQuantity(&payloads.RoomPayload{RoomType: room.Single, Quantity: 1}, hotelID)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	})

	t.Run("update room returns error", func(t *testing.T) {
		newRoom := &models.Rooms{RoomCategory: room.Single, AvailableQuantity: 3}
		mockRepo.EXPECT().GetAllRoomByHotelID(hotelID).Return([]*models.Rooms{newRoom}, nil)
		mockRepo.EXPECT().UpdateRoom(newRoom).Return(nil, errors.New("update fail"))

		err := svc.ReduceRoomQuantity(&payloads.RoomPayload{RoomType: room.Single, Quantity: 2}, hotelID)
		if err == nil {
			t.Fatalf("expected update error, got nil")
		}
	})

	t.Run("successful reduce", func(t *testing.T) {
		newRoom := &models.Rooms{RoomCategory: room.Single, AvailableQuantity: 5}
		// When GetAll is called, return room
		mockRepo.EXPECT().GetAllRoomByHotelID(hotelID).Return([]*models.Rooms{newRoom}, nil)
		// Expect UpdateRoom to be called with the modified room
		mockRepo.EXPECT().UpdateRoom(gomock.Any()).DoAndReturn(func(r *models.Rooms) (*models.Rooms, error) {
			// simulate DB returning updated room
			return r, nil
		})

		err := svc.ReduceRoomQuantity(&payloads.RoomPayload{RoomType: room.Single, Quantity: 3}, hotelID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// Ensure available quantity decreased by 3
		if newRoom.AvailableQuantity != 2 {
			t.Fatalf("expected AvailableQuantity 2, got %d", newRoom.AvailableQuantity)
		}
	})
}

func TestRoomService_IncreaseRoomQuantity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRoomRepoInterface(ctrl)
	svc := room_service.NewRoomService(mockRepo)

	hotelID := uuid.New()

	t.Run("repo GetAll error", func(t *testing.T) {
		mockRepo.EXPECT().GetAllRoomByHotelID(hotelID).Return(nil, errors.New("db"))
		_, err := svc.IncreaseRoomQuantity(&payloads.RoomPayload{RoomType: room.Single, Quantity: 1}, hotelID)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	})

	t.Run("update returns error", func(t *testing.T) {
		newRoom := &models.Rooms{RoomCategory: room.Single, AvailableQuantity: 2}
		mockRepo.EXPECT().GetAllRoomByHotelID(hotelID).Return([]*models.Rooms{newRoom}, nil)
		mockRepo.EXPECT().UpdateRoom(newRoom).Return(nil, errors.New("update fail"))

		_, err := svc.IncreaseRoomQuantity(&payloads.RoomPayload{RoomType: room.Single, Quantity: 3}, hotelID)
		if err == nil {
			t.Fatalf("expected update error, got nil")
		}
	})

	t.Run("room not found", func(t *testing.T) {
		mockRepo.EXPECT().GetAllRoomByHotelID(hotelID).Return([]*models.Rooms{
			{RoomCategory: "Other", AvailableQuantity: 1},
		}, nil)

		_, err := svc.IncreaseRoomQuantity(&payloads.RoomPayload{RoomType: "Missing", Quantity: 1}, hotelID)
		if err == nil {
			t.Fatalf("expected not found error, got nil")
		}
	})

	t.Run("successful increase", func(t *testing.T) {
		newRoom := &models.Rooms{RoomCategory: room.Single, AvailableQuantity: 2}
		mockRepo.EXPECT().GetAllRoomByHotelID(hotelID).Return([]*models.Rooms{newRoom}, nil)
		mockRepo.EXPECT().UpdateRoom(gomock.Any()).DoAndReturn(func(r *models.Rooms) (*models.Rooms, error) {
			return r, nil
		})

		got, err := svc.IncreaseRoomQuantity(&payloads.RoomPayload{RoomType: room.Single, Quantity: 5}, hotelID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got == nil {
			t.Fatalf("expected room returned, got nil")
		}
		if got.AvailableQuantity != 7 {
			t.Fatalf("expected available qty 7, got %d", got.AvailableQuantity)
		}
	})
}

func TestRoomService_GetAllRoomByHotelID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRoomRepoInterface(ctrl)
	svc := room_service.NewRoomService(mockRepo)

	hotelID := uuid.New()

	t.Run("return rooms", func(t *testing.T) {
		expected := []*models.Rooms{{RoomCategory: room.Single, AvailableQuantity: 1}}
		mockRepo.EXPECT().GetAllRoomByHotelID(hotelID).Return(expected, nil)

		got, err := svc.GetAllRoomByHotelID(hotelID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 1 {
			t.Fatalf("expected 1 room, got %d", len(got))
		}
	})

	t.Run("propagate error", func(t *testing.T) {
		mockRepo.EXPECT().GetAllRoomByHotelID(hotelID).Return(nil, errors.New("db"))
		_, err := svc.GetAllRoomByHotelID(hotelID)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	})
}
