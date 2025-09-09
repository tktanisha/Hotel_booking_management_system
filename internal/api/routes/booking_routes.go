package routes

import (
	"github.com/tktanisha/booking_system/internal/api/handlers"
	"github.com/tktanisha/booking_system/internal/api/middlewares"
	"github.com/tktanisha/booking_system/internal/api/router"
	"github.com/tktanisha/booking_system/internal/initializer"
)

func RegisterBookingRoutes(r *router.MuxRouter) {
	bookingHandler := handlers.NewBookingHandler(initializer.BookingService)

	r.HandleFunc("POST /bookings/create", bookingHandler.CreateBooking, middlewares.AuthMiddleware)
	r.HandleFunc("PUT /bookings/cancel/{bookingId}", bookingHandler.CancelBooking, middlewares.AuthMiddleware)
	r.HandleFunc("POST /bookings/checkout/{bookingId}", bookingHandler.CheckoutBooking, middlewares.AuthMiddleware)
}
