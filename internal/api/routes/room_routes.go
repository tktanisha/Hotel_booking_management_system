package routes

import (
	"github.com/tktanisha/booking_system/internal/api/handlers"
	"github.com/tktanisha/booking_system/internal/api/middlewares"
	"github.com/tktanisha/booking_system/internal/api/router"
	"github.com/tktanisha/booking_system/internal/initializer"
)

func RegisterRoomRoutes(r *router.MuxRouter) {
	roomHandler := handlers.NewRoomHandler(initializer.RoomService)

	r.HandleFunc("POST /rooms/create", roomHandler.CreateRoom, middlewares.AuthMiddleware)
	r.HandleFunc("GET /rooms/{hotelId}", roomHandler.GetAllRoomByHotelID, middlewares.AuthMiddleware)
	r.HandleFunc("PUT /rooms/increase-quantity/{hotelId}", roomHandler.IncreaseRoomQuantity, middlewares.AuthMiddleware)
}
