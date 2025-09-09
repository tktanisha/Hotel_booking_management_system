package booking_repo

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/tktanisha/booking_system/internal/db"
	"github.com/tktanisha/booking_system/internal/models"
)

type BookingRepo struct {
	db db.DB
}

func NewBookingRepo(database db.DB) *BookingRepo {
	return &BookingRepo{db: database}
}

func (r *BookingRepo) CreateBookingWithRooms(
	booking *models.Bookings,
	bookedRooms []*models.BookedRooms,
) (*models.Bookings, error) {

	// Step 1: Insert Booking
	bookingQuery := `
        INSERT INTO bookings (id, user_id, hotel_id, checkin, checkout, status, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id;
    `
	row := r.db.QueryRow(bookingQuery,
		booking.Id,
		booking.UserId,
		booking.HotelId,
		booking.CheckIn,
		booking.CheckOut,
		booking.Status,
		booking.CreatedAt,
	)
	if err := row.Scan(&booking.Id); err != nil {
		return nil, err
	}

	// Step 2: Insert Booked Rooms
	bookedRoomsQuery := `
        INSERT INTO booked_rooms (id, booking_id, room_type, room_quantity, created_at)
        VALUES ($1, $2, $3, $4, $5)
    `
	for _, room := range bookedRooms {
		_, err := r.db.Exec(bookedRoomsQuery,
			room.Id,
			booking.Id,
			room.RoomType,
			room.RoomQuantity,
			room.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
	}

	return booking, nil
}

func (r *BookingRepo) GetBookingById(bookingId uuid.UUID) (*models.Bookings, error) {
	query := `SELECT id, user_id, hotel_id, checkin, checkout, status, created_at FROM bookings WHERE id = $1`
	row := r.db.QueryRow(query, bookingId)

	var booking models.Bookings
	if err := row.Scan(&booking.Id, &booking.UserId, &booking.HotelId, &booking.CheckIn, &booking.CheckOut, &booking.Status, &booking.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("booking not found")
		}
		return nil, err
	}
	return &booking, nil
}

func (r *BookingRepo) GetBookedRoomsByBookingId(bookingId uuid.UUID) ([]*models.BookedRooms, error) {
	query := `SELECT id, booking_id, room_type, room_quantity, created_at FROM booked_rooms WHERE booking_id = $1`
	rows, err := r.db.Query(query, bookingId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var bookedRooms []*models.BookedRooms
	for rows.Next() {
		var room models.BookedRooms
		if err := rows.Scan(&room.Id, &room.BookingId, &room.RoomType, &room.RoomQuantity, &room.CreatedAt); err != nil {
			return nil, err
		}
		bookedRooms = append(bookedRooms, &room)
	}
	return bookedRooms, nil
}

func (r *BookingRepo) Save(booking *models.Bookings) error {
	query := `
		UPDATE bookings
		SET user_id=$2, hotel_id=$3, checkin=$4, checkout=$5, status=$6, created_at=$7
		WHERE id=$1
	`
	result, err := r.db.Exec(query, booking.Id, booking.UserId, booking.HotelId, booking.CheckIn, booking.CheckOut, booking.Status, booking.CreatedAt)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("booking not found")
	}
	return nil
}
