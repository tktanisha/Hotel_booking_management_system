package handlers

import (
	"net/http"

	"github.com/tktanisha/booking_system/internal/api/permissions"
	validators "github.com/tktanisha/booking_system/internal/api/validators/rooms_validators"
	"github.com/tktanisha/booking_system/internal/constants"
	"github.com/tktanisha/booking_system/internal/models"
	"github.com/tktanisha/booking_system/internal/services/room_service"
	"github.com/tktanisha/booking_system/internal/utils"
)

type RoomHandler struct {
	RoomService room_service.RoomServiceInterface
}

func NewRoomHandler(roomService room_service.RoomServiceInterface) *RoomHandler {
	return &RoomHandler{
		RoomService: roomService,
	}
}

func (h *RoomHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	userContext, ok := r.Context().Value(constants.UserContextKey).(*models.UserContext)
	if !ok || userContext == nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized", "User not found in context")
		return
	}

	if !permissions.IsManager(userContext) {
		utils.WriteErrorResponse(w, http.StatusForbidden, "Forbidden", "Only managers can create rooms")
		return
	}

	payload, err := validators.ValidateCreateRoomPayload(r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	room, err := h.RoomService.CreateRoom(payload)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create room", err.Error())
		return
	}

	utils.WriteSuccessResponse(w, http.StatusCreated, "Room created successfully!", room)
}

func (h *RoomHandler) GetAllRoomByHotelID(w http.ResponseWriter, r *http.Request) {
	userContext, ok := r.Context().Value(constants.UserContextKey).(*models.UserContext)
	if !ok || userContext == nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized", "User not found in context")
		return
	}

	hotelID, err := utils.GetUUIDFromParams(r, "hotelId")
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid hotel ID", err.Error())
		return
	}

	rooms, err := h.RoomService.GetAllRoomByHotelID(hotelID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve rooms", err.Error())
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, "Rooms retrieved successfully!", rooms)
}

func (h *RoomHandler) IncreaseRoomQuantity(w http.ResponseWriter, r *http.Request) {
	userContext, ok := r.Context().Value(constants.UserContextKey).(*models.UserContext)
	if !ok || userContext == nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized", "User not found in context")
		return
	}

	if !permissions.IsManager(userContext) {
		utils.WriteErrorResponse(w, http.StatusForbidden, "Forbidden", "Only managers can increase room quantity")
		return
	}

	hotelId, err := utils.GetUUIDFromParams(r, "hotelId")
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid hotel ID", err.Error())
		return
	}

	roomPayload, err := validators.ValidateRoomPayload(r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	updatedRooms := make([]*models.Rooms, 0)
	for _, roomToInc := range roomPayload {
		updated, err := h.RoomService.IncreaseRoomQuantity(roomToInc, hotelId)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to increase room quantity", err.Error())
			return
		}
		updatedRooms = append(updatedRooms, updated)
	}

	utils.WriteSuccessResponse(w, http.StatusOK, "Room quantity increased successfully!", updatedRooms)
}
