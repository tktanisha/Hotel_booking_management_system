package routes

import (
	"github.com/tktanisha/booking_system/internal/api/handlers"
	"github.com/tktanisha/booking_system/internal/api/middlewares"
	"github.com/tktanisha/booking_system/internal/api/router"
	"github.com/tktanisha/booking_system/internal/initializer"
)

func RegisterHotelRoutes(r *router.MuxRouter) {
	hotelHandler := handlers.NewHotelHandler(initializer.HotelService)

	r.HandleFunc("POST /hotels/create", hotelHandler.CreateHotel, middlewares.AuthMiddleware)
	r.HandleFunc("GET /hotels/{hotel_id}", hotelHandler.GetHotelByID, middlewares.AuthMiddleware)
}
