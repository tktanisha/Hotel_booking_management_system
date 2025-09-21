package handlers

import (
	"net/http"

	"github.com/tktanisha/booking_system/internal/constants"
	"github.com/tktanisha/booking_system/internal/models"
	"github.com/tktanisha/booking_system/internal/services/hotel_service"
	"github.com/tktanisha/booking_system/internal/utils"
	"github.com/tktanisha/booking_system/internal/utils/permissions"
	"github.com/tktanisha/booking_system/internal/utils/validators/hotel_validators"
)

type HotelHandler struct {
	HotelService hotel_service.HotelServiceInterface
}

func NewHotelHandler(hotelService hotel_service.HotelServiceInterface) *HotelHandler {
	return &HotelHandler{
		HotelService: hotelService,
	}
}

func (h *HotelHandler) CreateHotel(w http.ResponseWriter, r *http.Request) {
	userContext, ok := r.Context().Value(constants.UserContextKey).(*models.UserContext)
	if !ok || userContext == nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized", "User not found in context")
		return
	}

	if !permissions.IsManager(userContext) {
		utils.WriteErrorResponse(w, http.StatusForbidden, "Forbidden", "Only managers can create hotels")
		return
	}

	payload, err := hotel_validators.ValidateHotelPayload(r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	new_hotel, err := h.HotelService.CreateHotel(userContext, payload)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create hotel", err.Error())
		return
	}

	utils.WriteSuccessResponse(w, http.StatusCreated, "Hotel created successfully", new_hotel)
}

func (h *HotelHandler) GetHotelByID(w http.ResponseWriter, r *http.Request) {
	userContext, ok := r.Context().Value(constants.UserContextKey).(*models.UserContext)
	if !ok || userContext == nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized", "User not found in context")
		return
	}

	hotelID, err := utils.GetUUIDFromParams(r, "hotel_id")
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid hotel ID", err.Error())
		return
	}

	hotel, err := h.HotelService.GetHotelByID(hotelID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve hotel", err.Error())
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, "Hotel retrieved successfully", hotel)
}
