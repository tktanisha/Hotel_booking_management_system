package booking_repo_test

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"

	"github.com/tktanisha/booking_system/internal/models"
	"github.com/tktanisha/booking_system/internal/repository/booking_repo"
)

// Assumption: db.DB is implemented using *sql.DB in production.
// sqlmock allows simulation of database queries, rows, and transactions.

func TestBookingRepo_CreateBookingWithRooms(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func(mock sqlmock.Sqlmock, bookingID uuid.UUID)
		wantErr    bool
	}{
		{
			name: "success",
			setupMocks: func(mock sqlmock.Sqlmock, bookingID uuid.UUID) {
				mock.ExpectQuery(`INSERT INTO bookings`).
					WithArgs(bookingID, uuid.Nil, uuid.Nil, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(bookingID))

				mock.ExpectExec(`INSERT INTO booked_rooms`).
					WithArgs(sqlmock.AnyArg(), bookingID, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "booking insert fails",
			setupMocks: func(mock sqlmock.Sqlmock, bookingID uuid.UUID) {
				mock.ExpectQuery(`INSERT INTO bookings`).
					WithArgs(bookingID, uuid.Nil, uuid.Nil, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(errors.New("insert failed"))
			},
			wantErr: true,
		},
		{
			name: "booked room insert fails",
			setupMocks: func(mock sqlmock.Sqlmock, bookingID uuid.UUID) {
				mock.ExpectQuery(`INSERT INTO bookings`).
					WithArgs(bookingID, uuid.Nil, uuid.Nil, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(bookingID))

				mock.ExpectExec(`INSERT INTO booked_rooms`).
					WithArgs(sqlmock.AnyArg(), bookingID, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(errors.New("room insert failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create sqlmock: %v", err)
			}
			defer db.Close()

			repo := booking_repo.NewBookingRepo(db)
			bookingID := uuid.New()

			booking := &models.Bookings{
				Id:        bookingID,
				CreatedAt: time.Now(),
			}
			bookedRooms := []*models.BookedRooms{
				{Id: uuid.New(), CreatedAt: time.Now()},
			}

			tt.setupMocks(mock, bookingID)
			_, err = repo.CreateBookingWithRooms(booking, bookedRooms)

			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet sqlmock expectations: %v", err)
			}
		})
	}
}

func TestBookingRepo_GetBookingById(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func(mock sqlmock.Sqlmock, bookingID uuid.UUID)
		wantErr    bool
	}{
		{
			name: "success",
			setupMocks: func(mock sqlmock.Sqlmock, bookingID uuid.UUID) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "hotel_id", "checkin", "checkout", "status", "created_at"}).
					AddRow(bookingID, uuid.Nil, uuid.Nil, time.Now(), time.Now(), "confirmed", time.Now())
				mock.ExpectQuery(`SELECT id, user_id, hotel_id, checkin, checkout, status, created_at FROM bookings`).
					WithArgs(bookingID).WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name: "no rows found",
			setupMocks: func(mock sqlmock.Sqlmock, bookingID uuid.UUID) {
				mock.ExpectQuery(`SELECT id, user_id, hotel_id, checkin, checkout, status, created_at FROM bookings`).
					WithArgs(bookingID).WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
		{
			name: "query error",
			setupMocks: func(mock sqlmock.Sqlmock, bookingID uuid.UUID) {
				mock.ExpectQuery(`SELECT id, user_id, hotel_id, checkin, checkout, status, created_at FROM bookings`).
					WithArgs(bookingID).WillReturnError(errors.New("query failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create sqlmock: %v", err)
			}
			defer db.Close()

			repo := booking_repo.NewBookingRepo(db)
			bookingID := uuid.New()

			tt.setupMocks(mock, bookingID)
			_, err = repo.GetBookingById(bookingID)

			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet sqlmock expectations: %v", err)
			}
		})
	}
}

func TestBookingRepo_GetBookedRoomsByBookingId(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func(mock sqlmock.Sqlmock, bookingID uuid.UUID)
		wantErr    bool
	}{
		{
			name: "success",
			setupMocks: func(mock sqlmock.Sqlmock, bookingID uuid.UUID) {
				rows := sqlmock.NewRows([]string{"id", "booking_id", "room_type", "room_quantity", "created_at"}).
					AddRow(uuid.New(), bookingID, "single", 2, time.Now())
				mock.ExpectQuery(`SELECT id, booking_id, room_type, room_quantity, created_at FROM booked_rooms`).
					WithArgs(bookingID).WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name: "query error",
			setupMocks: func(mock sqlmock.Sqlmock, bookingID uuid.UUID) {
				mock.ExpectQuery(`SELECT id, booking_id, room_type, room_quantity, created_at FROM booked_rooms`).
					WithArgs(bookingID).WillReturnError(errors.New("query failed"))
			},
			wantErr: true,
		},
		{
			name: "row scan error",
			setupMocks: func(mock sqlmock.Sqlmock, bookingID uuid.UUID) {
				rows := sqlmock.NewRows([]string{"id", "booking_id", "room_type", "room_quantity", "created_at"}).
					AddRow("invalid-uuid", bookingID, "single", 2, time.Now())
				mock.ExpectQuery(`SELECT id, booking_id, room_type, room_quantity, created_at FROM booked_rooms`).
					WithArgs(bookingID).WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create sqlmock: %v", err)
			}
			defer db.Close()

			repo := booking_repo.NewBookingRepo(db)
			bookingID := uuid.New()

			tt.setupMocks(mock, bookingID)
			_, err = repo.GetBookedRoomsByBookingId(bookingID)

			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet sqlmock expectations: %v", err)
			}
		})
	}
}

func TestBookingRepo_Save(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func(mock sqlmock.Sqlmock, bookingID uuid.UUID)
		wantErr    bool
	}{
		{
			name: "success",
			setupMocks: func(mock sqlmock.Sqlmock, bookingID uuid.UUID) {
				mock.ExpectExec(`UPDATE bookings`).
					WithArgs(bookingID, uuid.Nil, uuid.Nil, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "update error",
			setupMocks: func(mock sqlmock.Sqlmock, bookingID uuid.UUID) {
				mock.ExpectExec(`UPDATE bookings`).
					WithArgs(bookingID, uuid.Nil, uuid.Nil, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(errors.New("update failed"))
			},
			wantErr: true,
		},
		{
			name: "no rows affected",
			setupMocks: func(mock sqlmock.Sqlmock, bookingID uuid.UUID) {
				mock.ExpectExec(`UPDATE bookings`).
					WithArgs(bookingID, uuid.Nil, uuid.Nil, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create sqlmock: %v", err)
			}
			defer db.Close()

			repo := booking_repo.NewBookingRepo(db)
			booking := &models.Bookings{
				Id:        uuid.New(),
				CreatedAt: time.Now(),
			}

			tt.setupMocks(mock, booking.Id)
			err = repo.Save(booking)

			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet sqlmock expectations: %v", err)
			}
		})
	}
}
