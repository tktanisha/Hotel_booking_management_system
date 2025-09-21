package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/tktanisha/booking_system/internal/api/handlers"
	authMocks "github.com/tktanisha/booking_system/internal/mocks"
	"github.com/tktanisha/booking_system/internal/models"
	"github.com/tktanisha/booking_system/internal/utils/validators/payloads"
)

func TestAuthHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := authMocks.NewMockAuthServiceInterface(ctrl)
	handler := handlers.NewAuthHandler(mockAuthService)

	tests := []struct {
		name           string
		body           payloads.LoginRequest
		mockService    func()
		wantStatusCode int
	}{
		{
			name: "invalid payload",
			body: payloads.LoginRequest{
				Email:    "",
				Password: "",
			},
			mockService:    func() {},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "service error",
			body: payloads.LoginRequest{
				Email:    "test@example.com",
				Password: "User@1234",
			},
			mockService: func() {
				mockAuthService.EXPECT().
					Login("test@example.com", "User@1234").
					Return("", nil, errors.New("invalid credentials"))
			},
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name: "success",
			body: payloads.LoginRequest{
				Email:    "success@example.com",
				Password: "User@1234",
			},
			mockService: func() {
				mockAuthService.EXPECT().
					Login("success@example.com", "User@1234").
					Return("token123", &models.Users{Id: uuid.New(), Email: "success@example.com"}, nil)
			},
			wantStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockService()
			jsonBody, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.Login(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("expected status %d, got %d", tt.wantStatusCode, w.Code)
			}
		})
	}
}

func TestAuthHandler_Register(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockAuthService := authMocks.NewMockAuthServiceInterface(ctrl)
	handler := handlers.NewAuthHandler(mockAuthService)

	tests := []struct {
		name           string
		body           payloads.RegisterRequest
		mockService    func()
		wantStatusCode int
	}{
		{
			name: "invalid payload",
			body: payloads.RegisterRequest{
				Email:    "",
				Password: "",
				Fullname: "",
			},
			mockService:    func() {},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "service error",
			body: payloads.RegisterRequest{
				Email:    "test@example.com",
				Password: "User@1234",
				Fullname: "Test User",
			},
			mockService: func() {
				mockAuthService.EXPECT().
					Register(gomock.Any()).
					Return(nil, errors.New("user already exists"))
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "success",
			body: payloads.RegisterRequest{
				Email:    "success@example.com",
				Password: "User@1234",
				Fullname: "Success User",
			},
			mockService: func() {
				mockAuthService.EXPECT().
					Register(gomock.Any()).
					Return(&models.Users{Id: uuid.New(), Email: "success@example.com"}, nil)
			},
			wantStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockService()
			jsonBody, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.Register(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("expected status %d, got %d", tt.wantStatusCode, w.Code)
			}
		})
	}
}
