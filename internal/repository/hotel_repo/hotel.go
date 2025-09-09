package hotel_repo

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/tktanisha/booking_system/internal/db"
	"github.com/tktanisha/booking_system/internal/models"
)

type HotelRepository struct {
	db db.DB
}

func NewHotelRepo(database db.DB) *HotelRepository {
	return &HotelRepository{db: database}
}

func (hr *HotelRepository) GetHotelByID(hotelID uuid.UUID) (*models.Hotels, error) {
	query := `
		SELECT id, manager_id, name, address, created_at
		FROM hotels
		WHERE id = $1
	`

	var hotel models.Hotels
	row := hr.db.QueryRow(query, hotelID)
	if err := row.Scan(&hotel.Id, &hotel.ManagerId, &hotel.Name, &hotel.Address, &hotel.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("hotel not found")
		}
		return nil, err
	}
	return &hotel, nil
}

func (hr *HotelRepository) CreateHotel(hotel *models.Hotels) (*models.Hotels, error) {
	query := `
		INSERT INTO hotels (id, manager_id, name, address, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`

	if hotel.Id == uuid.Nil {
		hotel.Id = uuid.New()
	}
	if hotel.CreatedAt.IsZero() {
		hotel.CreatedAt = time.Now()
	}

	row := hr.db.QueryRow(query, hotel.Id, hotel.ManagerId, hotel.Name, hotel.Address, hotel.CreatedAt)
	if err := row.Scan(&hotel.Id); err != nil {
		return nil, err
	}
	return hotel, nil
}
