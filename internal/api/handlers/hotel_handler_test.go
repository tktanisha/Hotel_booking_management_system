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
	"github.com/tktanisha/booking_system/internal/api/validators/payloads"
	"github.com/tktanisha/booking_system/internal/constants"
	user_role "github.com/tktanisha/booking_system/internal/enums/user"
	hotelMocks "github.com/tktanisha/booking_system/internal/mocks"
	"github.com/tktanisha/booking_system/internal/models"
)

// Mock payload for creating a hotel
type CreateHotelRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func TestHotelHandler_CreateHotel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHotelService := hotelMocks.NewMockHotelServiceInterface(ctrl)
	handler := handlers.NewHotelHandler(mockHotelService)

	managerCtx := &models.UserContext{Id: uuid.New(), Role: user_role.RoleManager}
	userCtx := &models.UserContext{Id: uuid.New(), Role: user_role.RoleUser}

	validPayload := &payloads.CreateHotelPayload{Name: "Test Hotel", Address: "123 Street"}

	tests := []struct {
		name           string
		ctx            context.Context
		body           interface{}
		mockService    func()
		wantStatusCode int
	}{
		{
			name:           "unauthorized",
			ctx:            context.Background(),
			body:           validPayload,
			mockService:    func() {},
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name:           "forbidden",
			ctx:            context.WithValue(context.Background(), constants.UserContextKey, userCtx),
			body:           validPayload,
			mockService:    func() {},
			wantStatusCode: http.StatusForbidden,
		},
		{
			name:           "invalid payload",
			ctx:            context.WithValue(context.Background(), constants.UserContextKey, managerCtx),
			body:           map[string]interface{}{"invalid": "data"},
			mockService:    func() {},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "service error",
			ctx:  context.WithValue(context.Background(), constants.UserContextKey, managerCtx),
			body: validPayload,
			mockService: func() {
				mockHotelService.EXPECT().
					CreateHotel(managerCtx, gomock.Any()).
					Return(nil, errors.New("service failed"))
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name: "success",
			ctx:  context.WithValue(context.Background(), constants.UserContextKey, managerCtx),
			body: validPayload,
			mockService: func() {
				mockHotelService.EXPECT().
					CreateHotel(managerCtx, gomock.Any()).
					Return(&models.Hotels{Id: uuid.New(), Name: "Test Hotel"}, nil)
			},
			wantStatusCode: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockService()
			jsonBody, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/hotel", bytes.NewReader(jsonBody))
			req = req.WithContext(tt.ctx)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.CreateHotel(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("expected status %d, got %d", tt.wantStatusCode, w.Code)
			}
		})
	}
}

func TestHotelHandler_GetHotelByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHotelService := hotelMocks.NewMockHotelServiceInterface(ctrl)
	handler := handlers.NewHotelHandler(mockHotelService)

	userCtx := &models.UserContext{Id: uuid.New(), Role: user_role.RoleUser}
	hotelID := uuid.New()

	tests := []struct {
		name           string
		ctx            context.Context
		hotelIDStr     string
		mockService    func()
		wantStatusCode int
	}{
		{
			name:           "unauthorized",
			ctx:            context.Background(),
			hotelIDStr:     hotelID.String(),
			mockService:    func() {},
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name:           "invalid hotel id",
			ctx:            context.WithValue(context.Background(), constants.UserContextKey, userCtx),
			hotelIDStr:     "invalid-uuid",
			mockService:    func() {},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:       "service error",
			ctx:        context.WithValue(context.Background(), constants.UserContextKey, userCtx),
			hotelIDStr: hotelID.String(),
			mockService: func() {
				mockHotelService.EXPECT().
					GetHotelByID(hotelID).
					Return(nil, errors.New("service failed"))
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:       "success",
			ctx:        context.WithValue(context.Background(), constants.UserContextKey, userCtx),
			hotelIDStr: hotelID.String(),
			mockService: func() {
				mockHotelService.EXPECT().
					GetHotelByID(hotelID).
					Return(&models.Hotels{Id: hotelID, Name: "Test Hotel"}, nil)
			},
			wantStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockService()
			req := httptest.NewRequest(http.MethodGet, "/hotel/", nil)
			req = req.WithContext(tt.ctx)
			req.SetPathValue("hotel_id", tt.hotelIDStr)
			w := httptest.NewRecorder()

			handler.GetHotelByID(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("expected status %d, got %d", tt.wantStatusCode, w.Code)
			}
		})
	}
}
