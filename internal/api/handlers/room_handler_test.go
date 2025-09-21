package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/tktanisha/booking_system/internal/api/handlers"
	"github.com/tktanisha/booking_system/internal/constants"
	"github.com/tktanisha/booking_system/internal/enums/room"
	user_role "github.com/tktanisha/booking_system/internal/enums/user"
	roomMocks "github.com/tktanisha/booking_system/internal/mocks"
	"github.com/tktanisha/booking_system/internal/models"
	"github.com/tktanisha/booking_system/internal/utils/validators/payloads"
)

func TestRoomHandler_CreateRoom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoomService := roomMocks.NewMockRoomServiceInterface(ctrl)
	handler := handlers.NewRoomHandler(mockRoomService)

	managerCtx := &models.UserContext{Id: uuid.New(), Role: user_role.RoleManager}
	userCtx := &models.UserContext{Id: uuid.New(), Role: user_role.RoleUser}

	roomPayload := &payloads.CreateRoomPayload{
		HotelID:  uuid.New(),
		RoomType: room.Double,
		Price:    150.0,
		Quantity: 5,
	}

	tests := []struct {
		name           string
		ctx            context.Context
		payload        interface{}
		mockService    func()
		wantStatusCode int
	}{
		{
			name:           "unauthorized",
			ctx:            context.Background(),
			payload:        roomPayload,
			mockService:    func() {},
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name:           "forbidden",
			ctx:            context.WithValue(context.Background(), constants.UserContextKey, userCtx),
			payload:        roomPayload,
			mockService:    func() {},
			wantStatusCode: http.StatusForbidden,
		},
		{
			name:           "invalid payload",
			ctx:            context.WithValue(context.Background(), constants.UserContextKey, managerCtx),
			payload:        map[string]interface{}{"invalid": "data"},
			mockService:    func() {},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:    "service error",
			ctx:     context.WithValue(context.Background(), constants.UserContextKey, managerCtx),
			payload: roomPayload,
			mockService: func() {
				mockRoomService.EXPECT().
					CreateRoom(roomPayload).
					Return(nil, errors.New("service failed"))
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:    "success",
			ctx:     context.WithValue(context.Background(), constants.UserContextKey, managerCtx),
			payload: roomPayload,
			mockService: func() {
				mockRoomService.EXPECT().
					CreateRoom(roomPayload).
					Return(&models.Rooms{Id: uuid.New(), RoomCategory: room.Double, AvailableQuantity: 5}, nil)
			},
			wantStatusCode: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockService()
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/rooms", bytes.NewReader(body))
			req = req.WithContext(tt.ctx)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.CreateRoom(w, req)
			if w.Code != tt.wantStatusCode {
				t.Errorf("expected status %d, got %d", tt.wantStatusCode, w.Code)
			}
		})
	}
}

func TestRoomHandler_GetAllRoomByHotelID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoomService := roomMocks.NewMockRoomServiceInterface(ctrl)
	handler := handlers.NewRoomHandler(mockRoomService)

	userCtx := &models.UserContext{Id: uuid.New()}
	hotelID := uuid.New()

	tests := []struct {
		name           string
		ctx            context.Context
		pathHotelID    string
		mockService    func()
		wantStatusCode int
	}{
		{
			name:           "unauthorized",
			ctx:            context.Background(),
			pathHotelID:    hotelID.String(),
			mockService:    func() {},
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name:           "invalid hotel id",
			ctx:            context.WithValue(context.Background(), constants.UserContextKey, userCtx),
			pathHotelID:    "invalid-uuid",
			mockService:    func() {},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:        "service error",
			ctx:         context.WithValue(context.Background(), constants.UserContextKey, userCtx),
			pathHotelID: hotelID.String(),
			mockService: func() {
				mockRoomService.EXPECT().
					GetAllRoomByHotelID(hotelID).
					Return(nil, errors.New("service failed"))
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:        "success",
			ctx:         context.WithValue(context.Background(), constants.UserContextKey, userCtx),
			pathHotelID: hotelID.String(),
			mockService: func() {
				mockRoomService.EXPECT().
					GetAllRoomByHotelID(hotelID).
					Return([]*models.Rooms{{Id: uuid.New(), RoomCategory: room.Double, AvailableQuantity: 5}}, nil)
			},
			wantStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockService()
			req := httptest.NewRequest(http.MethodGet, "/rooms/", nil)
			req = req.WithContext(tt.ctx)
			req.SetPathValue("hotelId", tt.pathHotelID)
			w := httptest.NewRecorder()

			handler.GetAllRoomByHotelID(w, req)
			if w.Code != tt.wantStatusCode {
				t.Errorf("expected %d, got %d", tt.wantStatusCode, w.Code)
			}
		})
	}
}

func TestRoomHandler_IncreaseRoomQuantity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoomService := roomMocks.NewMockRoomServiceInterface(ctrl)
	handler := handlers.NewRoomHandler(mockRoomService)

	managerCtx := &models.UserContext{Id: uuid.New(), Role: user_role.RoleManager}
	userCtx := &models.UserContext{Id: uuid.New(), Role: user_role.RoleUser}
	hotelID := uuid.New()

	roomPayload := []*payloads.RoomPayload{{RoomType: room.Double, Quantity: 3}}

	tests := []struct {
		name           string
		ctx            context.Context
		pathHotelID    string
		payload        interface{}
		mockService    func()
		wantStatusCode int
	}{
		{
			name:           "unauthorized",
			ctx:            context.Background(),
			pathHotelID:    hotelID.String(),
			payload:        roomPayload,
			mockService:    func() {},
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name:           "forbidden",
			ctx:            context.WithValue(context.Background(), constants.UserContextKey, userCtx),
			pathHotelID:    hotelID.String(),
			payload:        roomPayload,
			mockService:    func() {},
			wantStatusCode: http.StatusForbidden,
		},
		{
			name:           "invalid hotel id",
			ctx:            context.WithValue(context.Background(), constants.UserContextKey, managerCtx),
			pathHotelID:    "invalid-uuid",
			payload:        roomPayload,
			mockService:    func() {},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:        "service error",
			ctx:         context.WithValue(context.Background(), constants.UserContextKey, managerCtx),
			pathHotelID: hotelID.String(),
			payload:     roomPayload,
			mockService: func() {
				mockRoomService.EXPECT().
					IncreaseRoomQuantity(roomPayload[0], hotelID).
					Return(nil, errors.New("service failed"))
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:        "success",
			ctx:         context.WithValue(context.Background(), constants.UserContextKey, managerCtx),
			pathHotelID: hotelID.String(),
			payload:     roomPayload,
			mockService: func() {
				mockRoomService.EXPECT().
					IncreaseRoomQuantity(roomPayload[0], hotelID).
					Return(&models.Rooms{Id: uuid.New(), RoomCategory: room.Double, AvailableQuantity: 5}, nil)
			},
			wantStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockService()
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPut, "/rooms/"+tt.pathHotelID+"/increase", bytes.NewReader(body))
			req = req.WithContext(tt.ctx)
			req.Header.Set("Content-Type", "application/json")

			req.SetPathValue("hotelId", tt.pathHotelID)

			w := httptest.NewRecorder()
			handler.IncreaseRoomQuantity(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("expected %d, got %d", tt.wantStatusCode, w.Code)
			}
		})
	}
}
