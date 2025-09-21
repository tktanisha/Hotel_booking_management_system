package routes

import (
	"net/http"

	"github.com/tktanisha/booking_system/internal/api/handlers"
	"github.com/tktanisha/booking_system/internal/api/middlewares"
	"github.com/tktanisha/booking_system/internal/initializer"
)

func RegisterBookingRoutes(r *http.ServeMux) {
	bookingHandler := handlers.NewBookingHandler(initializer.BookingService)

	r.HandleFunc("POST /bookings/create", middlewares.AuthMiddleware(bookingHandler.CreateBooking))
	r.HandleFunc("PUT /bookings/cancel/{bookingId}", middlewares.AuthMiddleware(bookingHandler.CancelBooking))
	r.HandleFunc("POST /bookings/checkout/{bookingId}", middlewares.AuthMiddleware(bookingHandler.CheckoutBooking))
}
