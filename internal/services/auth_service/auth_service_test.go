package auth_service_test

import (
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	user_role "github.com/tktanisha/booking_system/internal/enums/user"
	"github.com/tktanisha/booking_system/internal/mocks"
	"github.com/tktanisha/booking_system/internal/models"
	"github.com/tktanisha/booking_system/internal/services/auth_service"
	"github.com/tktanisha/booking_system/internal/utils"
)

func TestAuthService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepoInterface(ctrl)
	svc := auth_service.NewAuthService(mockRepo)

	tests := []struct {
		name      string
		input     *models.Users
		mockSetup func()
		wantErr   bool
	}{
		{
			name:  "User already exists",
			input: &models.Users{Email: "existing@example.com", Password: "password"},
			mockSetup: func() {
				mockRepo.EXPECT().
					FindByEmail("existing@example.com").
					Return(&models.Users{}, nil)
			},
			wantErr: true,
		},
		{
			name:  "Successful registration",
			input: &models.Users{Email: "new@example.com", Password: "password"},
			mockSetup: func() {
				mockRepo.EXPECT().
					FindByEmail("new@example.com").
					Return(nil, errors.New("not found"))
				mockRepo.EXPECT().
					CreateUser(gomock.Any()).
					DoAndReturn(func(user *models.Users) (*models.Users, error) {
						user.Id = uuid.New()
						user.CreatedAt = time.Now()
						return user, nil
					})
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			_, err := svc.Register(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepoInterface(ctrl)
	svc := auth_service.NewAuthService(mockRepo)

	// Pre-hash password for success case
	hashedPass, _ := utils.HashPassword("correct")

	tests := []struct {
		name      string
		email     string
		password  string
		mockSetup func()
		wantErr   bool
	}{
		{
			name:     "User not found",
			email:    "notfound@example.com",
			password: "password",
			mockSetup: func() {
				mockRepo.EXPECT().
					FindByEmail("notfound@example.com").
					Return(nil, errors.New("not found"))
			},
			wantErr: true,
		},
		{
			name:     "Wrong password",
			email:    "user@example.com",
			password: "wrong",
			mockSetup: func() {
				mockRepo.EXPECT().
					FindByEmail("user@example.com").
					Return(&models.Users{Password: hashedPass}, nil)
			},
			wantErr: true,
		},
		{
			name:     "Successful login",
			email:    "user@example.com",
			password: "correct",
			mockSetup: func() {
				mockRepo.EXPECT().
					FindByEmail("user@example.com").
					Return(&models.Users{
						Id:       uuid.New(),
						Password: hashedPass,
						Role:     user_role.RoleUser,
					}, nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			_, _, err := svc.Login(tt.email, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
