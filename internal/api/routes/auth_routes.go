package routes

import (
	"github.com/tktanisha/booking_system/internal/api/handlers"
	"github.com/tktanisha/booking_system/internal/api/router"
	"github.com/tktanisha/booking_system/internal/initializer"
)

func RegisterAuthRoutes(r *router.MuxRouter) {
	authHandler := handlers.NewAuthHandler(initializer.AuthService)

	// POST /auth/login
	r.HandleFunc("POST /auth/login", authHandler.Login)

	// POST /auth/register
	r.HandleFunc("POST /auth/register", authHandler.Register)
}
