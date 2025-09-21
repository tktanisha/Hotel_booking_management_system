package auth_service

import (
	"errors"
	"time"

	user_role "github.com/tktanisha/booking_system/internal/enums/user"
	"github.com/tktanisha/booking_system/internal/models"
	"github.com/tktanisha/booking_system/internal/repository/user_repo"
	"github.com/tktanisha/booking_system/internal/utils"
)

type AuthService struct {
	userRepo user_repo.UserRepoInterface
}

func NewAuthService(repo user_repo.UserRepoInterface) *AuthService {
	return &AuthService{
		userRepo: repo,
	}
}

func (a *AuthService) Register(user *models.Users) (*models.Users, error) {

	_, err := a.userRepo.FindByEmail(user.Email)
	if err == nil {
		return nil, errors.New("user already exists")
	}

	hash, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.CreatedAt = time.Now()
	user.Password = hash
	user.Role = user_role.RoleUser

	new_user, err := a.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return new_user, nil
}

func (a *AuthService) Login(email, password string) (string, *models.Users, error) {

	user, err := a.userRepo.FindByEmail(email)
	if err != nil {
		return "", nil, errors.New("invalid email credentials")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", nil, errors.New("invalid password credentials")
	}

	token, err := utils.GenerateJWT(user.Id, user.Role)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}
