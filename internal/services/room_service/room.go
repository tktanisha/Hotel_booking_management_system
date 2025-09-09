package room_service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/tktanisha/booking_system/internal/api/validators/payloads"
	"github.com/tktanisha/booking_system/internal/models"
	"github.com/tktanisha/booking_system/internal/repository/room_repo"
	"github.com/tktanisha/booking_system/internal/services/room_service/factory"
)

type RoomService struct {
	RoomRepo room_repo.RoomRepoInterface
}

func NewRoomService(roomRepo room_repo.RoomRepoInterface) *RoomService {
	return &RoomService{
		RoomRepo: roomRepo,
	}
}

func (r *RoomService) CreateRoom(payload *payloads.CreateRoomPayload) (*models.Rooms, error) {
	factory, err := factory.GetRoomFactory(payload.RoomType)
	if err != nil {
		return nil, err
	}

	room := factory.Create(payload)

	room, err = r.RoomRepo.CreateRoom(room)
	if err != nil {
		return nil, err
	}

	return room, nil
}

func (r *RoomService) IsAvailable(room *payloads.RoomPayload, hotelId uuid.UUID) bool {
	rooms, err := r.RoomRepo.GetAllRoomByHotelID(hotelId)
	if err != nil {
		return false
	}
	for _, currentRoom := range rooms {
		if room.RoomType == currentRoom.RoomCategory && currentRoom.AvailableQuantity >= room.Quantity {
			return true
		}
	}
	return false
}

func (r *RoomService) ReduceRoomQuantity(room *payloads.RoomPayload, hotelId uuid.UUID) error {
	rooms, err := r.RoomRepo.GetAllRoomByHotelID(hotelId)
	if err != nil {
		return err
	}

	for _, currentRoom := range rooms {
		if currentRoom.RoomCategory == room.RoomType {
			currentRoom.AvailableQuantity -= room.Quantity
			if _, err := r.RoomRepo.UpdateRoom(currentRoom); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *RoomService) IncreaseRoomQuantity(room *payloads.RoomPayload, hotelId uuid.UUID) (*models.Rooms, error) {
	rooms, err := r.RoomRepo.GetAllRoomByHotelID(hotelId)
	if err != nil {
		return nil, err
	}

	for _, currentRoom := range rooms {
		if currentRoom.RoomCategory == room.RoomType {
			currentRoom.AvailableQuantity += room.Quantity
			if _, err := r.RoomRepo.UpdateRoom(currentRoom); err != nil {
				return nil, err
			}
			return currentRoom, nil
		}
	}

	return nil, errors.New("room type not found")
}

func (r *RoomService) GetAllRoomByHotelID(hotelID uuid.UUID) ([]*models.Rooms, error) {
	return r.RoomRepo.GetAllRoomByHotelID(hotelID)
}
