package routes

import (
	"net/http"

	"github.com/tktanisha/booking_system/internal/api/handlers"
	"github.com/tktanisha/booking_system/internal/api/middlewares"
	"github.com/tktanisha/booking_system/internal/initializer"
)

func RegisterRoomRoutes(r *http.ServeMux) {
	roomHandler := handlers.NewRoomHandler(initializer.RoomService)

	r.HandleFunc("POST /rooms/create", middlewares.AuthMiddleware(roomHandler.CreateRoom))
	r.HandleFunc("GET /rooms/{hotelId}", middlewares.AuthMiddleware(roomHandler.GetAllRoomByHotelID))
	r.HandleFunc("PUT /rooms/increase-quantity/{hotelId}", middlewares.AuthMiddleware(roomHandler.IncreaseRoomQuantity))
}
