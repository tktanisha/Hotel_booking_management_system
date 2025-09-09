package auth_service

import "github.com/tktanisha/booking_system/internal/models"

//go:generate mockgen -source=auth_service_interface.go -destination=../../mocks/mock_auth_service.go -package=mocks

type AuthServiceInterface interface {
	Login(email, password string) (string, *models.Users, error)
	Register(user *models.Users) (*models.Users, error)
}
