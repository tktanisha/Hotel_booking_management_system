package routes

import (
	"net/http"

	"github.com/tktanisha/booking_system/internal/api/handlers"
	"github.com/tktanisha/booking_system/internal/api/middlewares"
	"github.com/tktanisha/booking_system/internal/initializer"
)

func RegisterHotelRoutes(r *http.ServeMux) {
	hotelHandler := handlers.NewHotelHandler(initializer.HotelService)

	r.HandleFunc("POST /hotels/create", middlewares.AuthMiddleware(hotelHandler.CreateHotel))
	r.HandleFunc("GET /hotels/{hotel_id}", middlewares.AuthMiddleware(hotelHandler.GetHotelByID))
}
