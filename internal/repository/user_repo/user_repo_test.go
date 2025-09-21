package user_repo

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/tktanisha/booking_system/internal/models"
)

func TestUserRepo_CreateUser(t *testing.T) {
	tests := []struct {
		name          string
		user          *models.Users
		mockBehavior  func(mock sqlmock.Sqlmock, user *models.Users)
		expectedError bool
	}{
		{
			name: "Success - Create User with Predefined ID & CreatedAt",
			user: &models.Users{
				Id:        uuid.New(),
				Fullname:  "John Doe",
				Email:     "john@example.com",
				Password:  "password123",
				Role:      "user",
				CreatedAt: time.Now(),
			},
			mockBehavior: func(mock sqlmock.Sqlmock, user *models.Users) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(user.Id)
				mock.ExpectQuery(regexp.QuoteMeta(`
					INSERT INTO users (id, full_name, email, pass_word, role, created_at)
					VALUES ($1, $2, $3, $4, $5, $6)
					RETURNING id;
				`)).
					WithArgs(user.Id, user.Fullname, user.Email, user.Password, user.Role, user.CreatedAt).
					WillReturnRows(rows)
			},
			expectedError: false,
		},
		{
			name: "Success - Auto-generate ID & CreatedAt",
			user: &models.Users{
				Fullname: "Alice",
				Email:    "alice@example.com",
				Password: "securepass",
				Role:     "admin",
			},
			mockBehavior: func(mock sqlmock.Sqlmock, user *models.Users) {
				mock.ExpectQuery(regexp.QuoteMeta(`
					INSERT INTO users (id, full_name, email, pass_word, role, created_at)
					VALUES ($1, $2, $3, $4, $5, $6)
					RETURNING id;
				`)).
					WithArgs(sqlmock.AnyArg(), user.Fullname, user.Email, user.Password, user.Role, sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
			},
			expectedError: false,
		},
		{
			name: "Failure - Insert Error",
			user: &models.Users{
				Id:        uuid.New(),
				Fullname:  "Failed User",
				Email:     "fail@example.com",
				Password:  "pass",
				Role:      "guest",
				CreatedAt: time.Now(),
			},
			mockBehavior: func(mock sqlmock.Sqlmock, user *models.Users) {
				mock.ExpectQuery(regexp.QuoteMeta(`
					INSERT INTO users (id, full_name, email, pass_word, role, created_at)
					VALUES ($1, $2, $3, $4, $5, $6)
					RETURNING id;
				`)).
					WithArgs(user.Id, user.Fullname, user.Email, user.Password, user.Role, user.CreatedAt).
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

			tt.mockBehavior(mock, tt.user)

			repo := NewUserRepo(db)
			_, err = repo.CreateUser(tt.user)

			if tt.expectedError != (err != nil) {
				t.Errorf("expected error: %v, got: %v", tt.expectedError, err)
			}

		})
	}
}

func TestUserRepo_FindByEmail(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		mockBehavior  func(mock sqlmock.Sqlmock, email string)
		expectedError bool
	}{
		{
			name:  "Success - User Found",
			email: "found@example.com",
			mockBehavior: func(mock sqlmock.Sqlmock, email string) {
				rows := sqlmock.NewRows([]string{
					"id", "full_name", "email", "pass_word", "role", "created_at",
				}).AddRow(uuid.New(), "Jane Doe", email, "pass123", "user", time.Now())
				mock.ExpectQuery(regexp.QuoteMeta(`
					SELECT id, full_name, email, pass_word, role, created_at FROM users WHERE email = $1
				`)).WithArgs(email).WillReturnRows(rows)
			},
			expectedError: false,
		},
		{
			name:  "Failure - User Not Found",
			email: "missing@example.com",
			mockBehavior: func(mock sqlmock.Sqlmock, email string) {
				mock.ExpectQuery(regexp.QuoteMeta(`
					SELECT id, full_name, email, pass_word, role, created_at FROM users WHERE email = $1
				`)).WithArgs(email).WillReturnError(sql.ErrNoRows)
			},
			expectedError: true,
		},
		{
			name:  "Failure - Query Error",
			email: "error@example.com",
			mockBehavior: func(mock sqlmock.Sqlmock, email string) {
				mock.ExpectQuery(regexp.QuoteMeta(`
					SELECT id, full_name, email, pass_word, role, created_at FROM users WHERE email = $1
				`)).WithArgs(email).WillReturnError(errors.New("query failed"))
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

			tt.mockBehavior(mock, tt.email)

			repo := NewUserRepo(db)
			_, err = repo.FindByEmail(tt.email)

			if err != nil != tt.expectedError {
				t.Errorf("expected error: %v, got: %v", tt.expectedError, err)

			}
		})
	}
}
