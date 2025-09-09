package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/tktanisha/booking_system/internal/api/validators/auth_validators"
	"github.com/tktanisha/booking_system/internal/models"
	"github.com/tktanisha/booking_system/internal/services/auth_service"
	error_handler "github.com/tktanisha/booking_system/internal/utils"
	write_response "github.com/tktanisha/booking_system/internal/utils"
)

type AuthHandler struct {
	AuthService auth_service.AuthServiceInterface
}

func NewAuthHandler(authService auth_service.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	payload, err := auth_validators.LoginValidate(r)
	if err != nil {
		error_handler.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	token, user, err := h.AuthService.Login(payload.Email, payload.Password)
	if err != nil {
		error_handler.WriteErrorResponse(w, http.StatusUnauthorized, "Login failed", err.Error())
		return
	}

	write_response.WriteSuccessResponse(w, http.StatusOK, "Login successful", map[string]interface{}{
		"token": token,
		"user":  user,
	})
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	payload, err := auth_validators.RegisterValidate(r)
	if err != nil {
		error_handler.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	user := models.Users{
		Id:       uuid.New(),
		Email:    payload.Email,
		Password: payload.Password,
		Fullname: payload.Fullname,
	}

	new_user, err := h.AuthService.Register(&user)
	if err != nil {
		error_handler.WriteErrorResponse(w, http.StatusBadRequest, "Registration failed", err.Error())
		return
	}

	write_response.WriteSuccessResponse(w, http.StatusOK, "Registration successful", new_user)
}
