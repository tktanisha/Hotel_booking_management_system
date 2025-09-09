package factory

import (
	"errors"

	"github.com/tktanisha/booking_system/internal/enums/room"
)

func GetRoomFactory(roomType room.RoomType) (RoomFactory, error) {
	switch roomType {
	case room.Single:
		return &SingleRoomFactory{}, nil
	case room.Double:
		return &DoubleRoomFactory{}, nil
	case room.Suite:
		return &SuiteRoomFactory{}, nil
	default:
		return nil, errors.New("invalid room type")
	}
}
