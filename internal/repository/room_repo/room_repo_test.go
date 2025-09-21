package room_repo

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/tktanisha/booking_system/internal/models"
)

func TestRoomRepository_CreateRoom(t *testing.T) {
	tests := []struct {
		name          string
		room          *models.Rooms
		mockBehavior  func(mock sqlmock.Sqlmock, room *models.Rooms)
		expectedError bool
	}{
		{
			name: "Success - Room Created",
			room: &models.Rooms{
				Id:                uuid.New(),
				HotelId:           uuid.New(),
				AvailableQuantity: 5,
				RoomCategory:      "double",
				CreatedAt:         time.Now(),
			},
			mockBehavior: func(mock sqlmock.Sqlmock, room *models.Rooms) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(room.Id)
				mock.ExpectQuery(regexp.QuoteMeta(`
					INSERT INTO rooms (id, hotel_id, available_quantity, room_category, created_at)
					VALUES ($1, $2, $3, $4, $5)
					RETURNING id;
				`)).
					WithArgs(room.Id, room.HotelId, room.AvailableQuantity, room.RoomCategory, room.CreatedAt).
					WillReturnRows(rows)
			},
			expectedError: false,
		},
		{
			name: "Failure - Insert Error",
			room: &models.Rooms{
				Id:                uuid.New(),
				HotelId:           uuid.New(),
				AvailableQuantity: 2,
				RoomCategory:      "Single",
				CreatedAt:         time.Now(),
			},
			mockBehavior: func(mock sqlmock.Sqlmock, room *models.Rooms) {
				mock.ExpectQuery(regexp.QuoteMeta(`
					INSERT INTO rooms (id, hotel_id, available_quantity, room_category, created_at)
					VALUES ($1, $2, $3, $4, $5)
					RETURNING id;
				`)).
					WithArgs(room.Id, room.HotelId, room.AvailableQuantity, room.RoomCategory, room.CreatedAt).
					WillReturnError(errors.New("insert failed"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("error opening mock db: %s", err)
			}
			defer db.Close()

			tt.mockBehavior(mock, tt.room)

			repo := NewRoomRepo(db)
			_, err = repo.CreateRoom(tt.room)

			if tt.expectedError != (err != nil) {
				t.Errorf("expected error: %v, got: %v", tt.expectedError, err)
			}

		})
	}
}

func TestRoomRepository_GetAllRoomByHotelID(t *testing.T) {
	hotelID := uuid.New()

	tests := []struct {
		name          string
		mockBehavior  func(mock sqlmock.Sqlmock, hotelID uuid.UUID)
		expectedError bool
	}{
		{
			name: "Success - Rooms Found",
			mockBehavior: func(mock sqlmock.Sqlmock, hotelID uuid.UUID) {
				rows := sqlmock.NewRows([]string{
					"id", "hotel_id", "available_quantity", "room_category", "created_at",
				}).
					AddRow(uuid.New(), hotelID, 10, "Single", time.Now()).
					AddRow(uuid.New(), hotelID, 3, "Double", time.Now())
				mock.ExpectQuery(regexp.QuoteMeta(`
					SELECT id, hotel_id, available_quantity, room_category, created_at
					FROM rooms
					WHERE hotel_id = $1
				`)).WithArgs(hotelID).WillReturnRows(rows)
			},
			expectedError: false,
		},
		{
			name: "Failure - Query Error",
			mockBehavior: func(mock sqlmock.Sqlmock, hotelID uuid.UUID) {
				mock.ExpectQuery(regexp.QuoteMeta(`
					SELECT id, hotel_id, available_quantity, room_category, created_at
					FROM rooms
					WHERE hotel_id = $1
				`)).WithArgs(hotelID).WillReturnError(errors.New("query failed"))
			},
			expectedError: true,
		},
		{
			name: "Failure - Scan Error",
			mockBehavior: func(mock sqlmock.Sqlmock, hotelID uuid.UUID) {
				rows := sqlmock.NewRows([]string{
					"id", "hotel_id", "available_quantity", "room_category", "created_at",
				}).
					AddRow("invalid-uuid", hotelID, 5, "single", time.Now())
				mock.ExpectQuery(regexp.QuoteMeta(`
					SELECT id, hotel_id, available_quantity, room_category, created_at
					FROM rooms
					WHERE hotel_id = $1
				`)).WithArgs(hotelID).WillReturnRows(rows)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("error opening mock db: %s", err)
			}
			defer db.Close()

			tt.mockBehavior(mock, hotelID)

			repo := NewRoomRepo(db)
			_, err = repo.GetAllRoomByHotelID(hotelID)

			if tt.expectedError != (err != nil) {
				t.Errorf("expected error containing: %v, got: %v", tt.expectedError, err)
			}
		})
	}
}

func TestRoomRepository_UpdateRoom(t *testing.T) {
	tests := []struct {
		name          string
		room          *models.Rooms
		mockBehavior  func(mock sqlmock.Sqlmock, room *models.Rooms)
		expectedError error
	}{
		{
			name: "Success - Room Updated",
			room: &models.Rooms{
				Id:                uuid.New(),
				HotelId:           uuid.New(),
				AvailableQuantity: 8,
				RoomCategory:      "Premium",
				CreatedAt:         time.Now(),
			},
			mockBehavior: func(mock sqlmock.Sqlmock, room *models.Rooms) {
				mock.ExpectExec(regexp.QuoteMeta(`
					UPDATE rooms
					SET hotel_id=$2, available_quantity=$3, room_category=$4, created_at=$5
					WHERE id=$1
				`)).
					WithArgs(room.Id, room.HotelId, room.AvailableQuantity, room.RoomCategory, room.CreatedAt).
					WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected
			},
			expectedError: nil,
		},
		{
			name: "Failure - Room Not Found",
			room: &models.Rooms{
				Id:                uuid.New(),
				HotelId:           uuid.New(),
				AvailableQuantity: 4,
				RoomCategory:      "Standard",
				CreatedAt:         time.Now(),
			},
			mockBehavior: func(mock sqlmock.Sqlmock, room *models.Rooms) {
				mock.ExpectExec(regexp.QuoteMeta(`
					UPDATE rooms
					SET hotel_id=$2, available_quantity=$3, room_category=$4, created_at=$5
					WHERE id=$1
				`)).
					WithArgs(room.Id, room.HotelId, room.AvailableQuantity, room.RoomCategory, room.CreatedAt).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedError: errors.New("room not found"),
		},
		{
			name: "Failure - Update Error",
			room: &models.Rooms{
				Id:                uuid.New(),
				HotelId:           uuid.New(),
				AvailableQuantity: 7,
				RoomCategory:      "Suite",
				CreatedAt:         time.Now(),
			},
			mockBehavior: func(mock sqlmock.Sqlmock, room *models.Rooms) {
				mock.ExpectExec(regexp.QuoteMeta(`
					UPDATE rooms
					SET hotel_id=$2, available_quantity=$3, room_category=$4, created_at=$5
					WHERE id=$1
				`)).
					WithArgs(room.Id, room.HotelId, room.AvailableQuantity, room.RoomCategory, room.CreatedAt).
					WillReturnError(errors.New("update failed"))
			},
			expectedError: errors.New("update failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("error opening mock db: %s", err)
			}
			defer db.Close()

			tt.mockBehavior(mock, tt.room)

			repo := NewRoomRepo(db)
			result, err := repo.UpdateRoom(tt.room)

			if tt.expectedError != nil {
				if err == nil || err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error: %v, got: %v", tt.expectedError, err)
				}
				if result != nil {
					t.Errorf("expected nil, got: %+v", result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result == nil {
					t.Errorf("expected room, got nil")
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet expectations: %v", err)
			}
		})
	}
}

func containsError(got, want string) bool {
	return regexp.MustCompile(regexp.QuoteMeta(want)).MatchString(got)
}
