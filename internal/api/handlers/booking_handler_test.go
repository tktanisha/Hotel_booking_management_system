package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/tktanisha/booking_system/internal/api/handlers"
	"github.com/tktanisha/booking_system/internal/api/validators/payloads"
	"github.com/tktanisha/booking_system/internal/constants"
	"github.com/tktanisha/booking_system/internal/enums/room"
	bookingMocks "github.com/tktanisha/booking_system/internal/mocks"
	"github.com/tktanisha/booking_system/internal/models"
)

type CreateBookingRequest struct {
	HotelId  uuid.UUID               `json:"hotel_id"`
	CheckIn  time.Time               `json:"check_in"`
	CheckOut time.Time               `json:"check_out"`
	Rooms    []*payloads.RoomPayload `json:"rooms"`
}

func TestBookingHandler_CreateBooking(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingService := bookingMocks.NewMockBookingServiceInterface(ctrl)
	handler := handlers.NewBookingHandler(mockBookingService)

	userCtx := &models.UserContext{Id: uuid.New()}

	hotelID := uuid.New()
	validPayload := &payloads.BookingPayload{
		HotelId:  hotelID,
		CheckIn:  time.Now(),
		CheckOut: time.Now(),
		Rooms: []*payloads.RoomPayload{
			{RoomType: room.Single, Quantity: 2},
		},
	}

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
			name:           "invalid payload",
			ctx:            context.WithValue(context.Background(), constants.UserContextKey, userCtx),
			body:           map[string]interface{}{"invalid": "data"},
			mockService:    func() {},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "service error",
			ctx:  context.WithValue(context.Background(), constants.UserContextKey, userCtx),
			body: validPayload,
			mockService: func() {
				mockBookingService.EXPECT().
					CreateBooking(userCtx, gomock.Any()).
					Return(nil, errors.New("service failed"))
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name: "success",
			ctx:  context.WithValue(context.Background(), constants.UserContextKey, userCtx),
			body: validPayload,
			mockService: func() {
				mockBookingService.EXPECT().
					CreateBooking(userCtx, gomock.Any()).
					Return(&models.Bookings{Id: uuid.New(), HotelId: hotelID}, nil)
			},
			wantStatusCode: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockService()

			jsonBody, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/booking", bytes.NewReader(jsonBody))
			req = req.WithContext(tt.ctx)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.CreateBooking(w, req)
			if w.Code != tt.wantStatusCode {
				t.Errorf("expected status %d, got %d", tt.wantStatusCode, w.Code)
			}
		})
	}
}

func TestBookingHandler_CancelBooking(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingService := bookingMocks.NewMockBookingServiceInterface(ctrl)
	handler := handlers.NewBookingHandler(mockBookingService)

	userCtx := &models.UserContext{Id: uuid.New()}
	bookingID := uuid.New()

	tests := []struct {
		name           string
		ctx            context.Context
		bookingIDStr   string
		mockService    func()
		wantStatusCode int
	}{
		{
			name:           "unauthorized",
			ctx:            context.Background(),
			bookingIDStr:   bookingID.String(),
			mockService:    func() {},
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name:           "invalid booking id",
			ctx:            context.WithValue(context.Background(), constants.UserContextKey, userCtx),
			bookingIDStr:   "invalid-uuid",
			mockService:    func() {},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:         "service error",
			ctx:          context.WithValue(context.Background(), constants.UserContextKey, userCtx),
			bookingIDStr: bookingID.String(),
			mockService: func() {
				mockBookingService.EXPECT().
					CancelBooking(bookingID).
					Return(nil, errors.New("service failed"))
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:         "success",
			ctx:          context.WithValue(context.Background(), constants.UserContextKey, userCtx),
			bookingIDStr: bookingID.String(),
			mockService: func() {
				mockBookingService.EXPECT().
					CancelBooking(bookingID).
					Return(&models.Bookings{Id: bookingID}, nil)
			},
			wantStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockService()
			req := httptest.NewRequest(http.MethodPost, "/booking/cancel/", nil)
			req = req.WithContext(tt.ctx)
			w := httptest.NewRecorder()
			req.SetPathValue("bookingId", tt.bookingIDStr)
			handler.CancelBooking(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("expected status %d, got %d", tt.wantStatusCode, w.Code)
			}
		})
	}
}

func TestBookingHandler_CheckoutBooking(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingService := bookingMocks.NewMockBookingServiceInterface(ctrl)
	handler := handlers.NewBookingHandler(mockBookingService)

	userCtx := &models.UserContext{Id: uuid.New()}
	bookingID := uuid.New()

	tests := []struct {
		name           string
		ctx            context.Context
		bookingIDStr   string
		mockService    func()
		wantStatusCode int
	}{
		{
			name:           "unauthorized",
			ctx:            context.Background(),
			bookingIDStr:   bookingID.String(),
			mockService:    func() {},
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name:           "invalid booking id",
			ctx:            context.WithValue(context.Background(), constants.UserContextKey, userCtx),
			bookingIDStr:   "invalid-uuid",
			mockService:    func() {},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:         "service error",
			ctx:          context.WithValue(context.Background(), constants.UserContextKey, userCtx),
			bookingIDStr: bookingID.String(),
			mockService: func() {
				mockBookingService.EXPECT().
					CheckoutBooking(bookingID).
					Return(nil, errors.New("service failed"))
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:         "success",
			ctx:          context.WithValue(context.Background(), constants.UserContextKey, userCtx),
			bookingIDStr: bookingID.String(),
			mockService: func() {
				mockBookingService.EXPECT().
					CheckoutBooking(bookingID).
					Return(&models.Bookings{Id: bookingID}, nil)
			},
			wantStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockService()
			req := httptest.NewRequest(http.MethodPost, "/booking/checkout/", nil)
			req = req.WithContext(tt.ctx)
			w := httptest.NewRecorder()

			req.SetPathValue("bookingId", tt.bookingIDStr)

			handler.CheckoutBooking(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("expected status %d, got %d", tt.wantStatusCode, w.Code)
			}
		})
	}
}
