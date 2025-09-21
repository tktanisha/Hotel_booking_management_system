package room_repo

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/tktanisha/booking_system/internal/db"
	"github.com/tktanisha/booking_system/internal/models"
)

type RoomRepository struct {
	db db.DB
}

func NewRoomRepo(database db.DB) *RoomRepository {
	return &RoomRepository{db: database}
}

func (rr *RoomRepository) CreateRoom(room *models.Rooms) (*models.Rooms, error) {
	query := `
		INSERT INTO rooms (id, hotel_id, available_quantity, room_category, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`

	row := rr.db.QueryRow(query, room.Id, room.HotelId, room.AvailableQuantity, room.RoomCategory, room.CreatedAt)
	if err := row.Scan(&room.Id); err != nil {
		return nil, err
	}
	return room, nil
}

func (rr *RoomRepository) GetAllRoomByHotelID(hotelID uuid.UUID) ([]*models.Rooms, error) {
	query := `
		SELECT id, hotel_id, available_quantity, room_category, created_at
		FROM rooms
		WHERE hotel_id = $1
	`

	rows, err := rr.db.Query(query, hotelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []*models.Rooms
	for rows.Next() {
		room := &models.Rooms{}
		if err := rows.Scan(&room.Id, &room.HotelId, &room.AvailableQuantity, &room.RoomCategory, &room.CreatedAt); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		fmt.Println("Error after iterating rows:", err)
		return nil, err
	}

	return rooms, nil
}

func (rr *RoomRepository) UpdateRoom(room *models.Rooms) (*models.Rooms, error) {
	query := `
		UPDATE rooms
		SET hotel_id=$2, available_quantity=$3, room_category=$4, created_at=$5
		WHERE id=$1
	`

	result, err := rr.db.Exec(query, room.Id, room.HotelId, room.AvailableQuantity, room.RoomCategory, room.CreatedAt)
	if err != nil {
		return nil, err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, errors.New("room not found")
	}

	return room, nil
}
