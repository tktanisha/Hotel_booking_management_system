package handlers

import (
	"net/http"

	"github.com/tktanisha/booking_system/internal/constants"
	"github.com/tktanisha/booking_system/internal/models"
	"github.com/tktanisha/booking_system/internal/utils"
	error_handler "github.com/tktanisha/booking_system/internal/utils"
	write_response "github.com/tktanisha/booking_system/internal/utils"
	validators "github.com/tktanisha/booking_system/internal/utils/validators/booking_validators"

	"github.com/tktanisha/booking_system/internal/services/booking_service"
)

type BookingHandler struct {
	BookingService booking_service.BookingServiceInterface
}

func NewBookingHandler(bookingService booking_service.BookingServiceInterface) *BookingHandler {
	return &BookingHandler{
		BookingService: bookingService,
	}
}

func (b *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	userContext, ok := r.Context().Value(constants.UserContextKey).(*models.UserContext)
	if !ok || userContext == nil {
		error_handler.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized", "User not found in context")
		return
	}
	payload, err := validators.CreateBookingValidator(r)
	if err != nil {
		error_handler.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request Payload", err.Error())
		return
	}

	createdBooking, err := b.BookingService.CreateBooking(userContext, payload)
	if err != nil {
		error_handler.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create booking", err.Error())
		return
	}
	write_response.WriteSuccessResponse(w, http.StatusCreated, "Booking created successfully!", createdBooking)
}

func (b *BookingHandler) CancelBooking(w http.ResponseWriter, r *http.Request) {
	userContext, ok := r.Context().Value(constants.UserContextKey).(*models.UserContext)
	if !ok || userContext == nil {
		error_handler.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized", "User not found in context")
		return
	}

	bookingId, err := utils.GetUUIDFromParams(r, "bookingId")
	if err != nil {
		error_handler.WriteErrorResponse(w, http.StatusBadRequest, "Invalid booking ID", err.Error())
		return
	}

	if err := validators.ValidateCancelBooking(bookingId); err != nil {
		error_handler.WriteErrorResponse(w, http.StatusBadRequest, "Validation error", err.Error())
		return
	}

	booking, err := b.BookingService.CancelBooking(bookingId)
	if err != nil {
		error_handler.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to cancel booking", err.Error())
		return
	}

	write_response.WriteSuccessResponse(w, http.StatusOK, "Booking canceled successfully", booking)
}

func (b *BookingHandler) CheckoutBooking(w http.ResponseWriter, r *http.Request) {
	userContext, ok := r.Context().Value(constants.UserContextKey).(*models.UserContext)
	if !ok || userContext == nil {
		error_handler.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized", "User not found in context")
		return
	}
	bookingId, err := utils.GetUUIDFromParams(r, "bookingId")
	if err != nil {
		error_handler.WriteErrorResponse(w, http.StatusBadRequest, "Invalid booking ID", err.Error())
		return
	}

	booking, err := b.BookingService.CheckoutBooking(bookingId)
	if err != nil {
		error_handler.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to checkout booking", err.Error())
		return
	}
	write_response.WriteSuccessResponse(w, http.StatusOK, "Booking checked out successfully", booking)
}
