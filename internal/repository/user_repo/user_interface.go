package user_repo

import (
	"github.com/tktanisha/booking_system/internal/models"
)

//go:generate mockgen -source=user_interface.go -destination=../../mocks/mock_user_repo.go -package=mocks

type UserRepoInterface interface {
	CreateUser(user *models.Users) (*models.Users, error)
	FindByEmail(email string) (*models.Users, error)
}
