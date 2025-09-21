package room_service_test

import (
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"

	"github.com/tktanisha/booking_system/internal/enums/room"
	"github.com/tktanisha/booking_system/internal/mocks"
	"github.com/tktanisha/booking_system/internal/models"
	"github.com/tktanisha/booking_system/internal/services/room_service"
	"github.com/tktanisha/booking_system/internal/utils/validators/payloads"
)

func TestRoomService_CreateRoom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRoomRepoInterface(ctrl)
	svc := room_service.NewRoomService(mockRepo)

	test := []struct {
		name     string
		mockFunc func()
		roomReq  *payloads.CreateRoomPayload
		wantErr  bool
	}{
		{
			name:     "get factory error of roomtype",
			roomReq:  &payloads.CreateRoomPayload{HotelID: uuid.New(), RoomType: "invalid_type", Quantity: 1, Price: 456},
			mockFunc: func() {},
			wantErr:  true,
		},
		{
			name:    "room repo on creating gives error",
			roomReq: &payloads.CreateRoomPayload{HotelID: uuid.New(), RoomType: room.Single, Quantity: 1, Price: 456},
			mockFunc: func() {
				mockRepo.EXPECT().
					CreateRoom(gomock.Any()).
					Return(nil, errors.New("repository error"))
			},
			wantErr: true,
		},
		{
			name:    "success of room creation",
			roomReq: &payloads.CreateRoomPayload{HotelID: uuid.New(), RoomType: room.Single, Quantity: 1, Price: 456},
			mockFunc: func() {
				mockRepo.EXPECT().
					CreateRoom(gomock.Any()).
					DoAndReturn(func(r *models.Rooms) (*models.Rooms, error) {
						r.Id = uuid.New()
						r.CreatedAt = time.Now()
						return r, nil
					})
			},
			wantErr: false,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			_, err := svc.CreateRoom(tt.roomReq)

			if (err != nil) != tt.wantErr {
				t.Errorf("RoomService.CreateRoom() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
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
			roomReq: &payloads.RoomPayload{RoomType: room.Single, Quantity: 1},
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

	mockRepo := mocks.NewMockRoomRepoInterface(ctrl)
	svc := room_service.NewRoomService(mockRepo)

	hotelID := uuid.New()

	singleRoom := &models.Rooms{
		RoomCategory:      room.Single,
		AvailableQuantity: 5,
	}

	tests := []struct {
		name     string
		hotelID  uuid.UUID
		payload  *payloads.RoomPayload
		mockFunc func()
		wantErr  bool
	}{
		{
			name:    "repo GetAll error",
			hotelID: hotelID,
			payload: &payloads.RoomPayload{RoomType: room.Single, Quantity: 1},
			mockFunc: func() {
				mockRepo.EXPECT().GetAllRoomByHotelID(hotelID).Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name:    "update room returns error",
			hotelID: hotelID,
			payload: &payloads.RoomPayload{RoomType: room.Single, Quantity: 2},
			mockFunc: func() {
				mockRepo.EXPECT().GetAllRoomByHotelID(hotelID).Return([]*models.Rooms{singleRoom}, nil)
				mockRepo.EXPECT().UpdateRoom(singleRoom).Return(nil, errors.New("update fail"))
			},
			wantErr: true,
		},
		{
			name:    "successful reduce",
			hotelID: hotelID,
			payload: &payloads.RoomPayload{RoomType: room.Single, Quantity: 3},
			mockFunc: func() {
				mockRepo.EXPECT().GetAllRoomByHotelID(hotelID).Return([]*models.Rooms{singleRoom}, nil)
				mockRepo.EXPECT().UpdateRoom(gomock.Any()).DoAndReturn(func(r *models.Rooms) (*models.Rooms, error) {
					return r, nil
				})
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			err := svc.ReduceRoomQuantity(tt.payload, tt.hotelID)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReduceRoomQuantity() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRoomService_IncreaseRoomQuantity(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mocks.NewMockRoomRepoInterface(ctrl)
	svc := room_service.NewRoomService(mockRepo)

	hotelID := uuid.New()

	singleRoom := &models.Rooms{
		RoomCategory:      room.Single,
		AvailableQuantity: 5,
	}

	tests := []struct {
		name     string
		hotelID  uuid.UUID
		payload  *payloads.RoomPayload
		mockFunc func()
		wantErr  bool
	}{
		{
			name:    "repo GetAll error",
			hotelID: hotelID,
			payload: &payloads.RoomPayload{RoomType: room.Single, Quantity: 1},
			mockFunc: func() {
				mockRepo.EXPECT().GetAllRoomByHotelID(hotelID).Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name:    "update room returns error",
			hotelID: hotelID,
			payload: &payloads.RoomPayload{RoomType: room.Single, Quantity: 2},
			mockFunc: func() {
				mockRepo.EXPECT().GetAllRoomByHotelID(hotelID).Return([]*models.Rooms{singleRoom}, nil)
				mockRepo.EXPECT().UpdateRoom(singleRoom).Return(nil, errors.New("update fail"))
			},
			wantErr: true,
		},
		{
			name:    "successful increase",
			hotelID: hotelID,
			payload: &payloads.RoomPayload{RoomType: room.Single, Quantity: 3},
			mockFunc: func() {
				mockRepo.EXPECT().GetAllRoomByHotelID(hotelID).Return([]*models.Rooms{singleRoom}, nil)
				mockRepo.EXPECT().UpdateRoom(gomock.Any()).DoAndReturn(func(r *models.Rooms) (*models.Rooms, error) {
					return r, nil
				})
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			err := svc.ReduceRoomQuantity(tt.payload, tt.hotelID)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReduceRoomQuantity() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRoomService_GetAllRoomByHotelID(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mocks.NewMockRoomRepoInterface(ctrl)
	service := room_service.NewRoomService(mockRepo)

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
				mockRepo.EXPECT().GetAllRoomByHotelID(HotelID).Return(nil, errors.New("not found"))
			},
			wantErr: true,
		},
		{
			name:    "succesfully fetched the hotel",
			hotelId: HotelID,
			mockFunc: func() {
				mockRepo.EXPECT().GetAllRoomByHotelID(HotelID).Return([]*models.Rooms{{AvailableQuantity: 4, RoomCategory: room.Double}}, nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			_, err := service.GetAllRoomByHotelID(tt.hotelId)

			if (err != nil) != tt.wantErr {
				t.Errorf("want = %v vand got =%v", tt.wantErr, err)
			}
		})
	}

}
