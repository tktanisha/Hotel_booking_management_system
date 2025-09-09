package room_repo

import (
	"github.com/google/uuid"
	"github.com/tktanisha/booking_system/internal/models"
)

//go:generate mockgen -source=room_interface.go -destination=../../mocks/mock_room_repo.go -package=mocks

type RoomRepoInterface interface {
	CreateRoom(*models.Rooms) (*models.Rooms, error)
	GetAllRoomByHotelID(hotelID uuid.UUID) ([]*models.Rooms, error)
	UpdateRoom(*models.Rooms) (*models.Rooms, error)
}
