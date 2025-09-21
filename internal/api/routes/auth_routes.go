package routes

import (
	"net/http"

	"github.com/tktanisha/booking_system/internal/api/handlers"
	"github.com/tktanisha/booking_system/internal/initializer"
)

func RegisterAuthRoutes(r *http.ServeMux) {
	authHandler := handlers.NewAuthHandler(initializer.AuthService)

	r.HandleFunc("POST /auth/login", authHandler.Login)
	r.HandleFunc("POST /auth/register", authHandler.Register)
}
