package user_repo

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/tktanisha/booking_system/internal/db"
	"github.com/tktanisha/booking_system/internal/models"
	"time"
)

type UserRepo struct {
	db db.DB
}

func NewUserRepo(database db.DB) *UserRepo {
	return &UserRepo{db: database}
}

func (r *UserRepo) CreateUser(user *models.Users) (*models.Users, error) {
	query := `
		INSERT INTO users (id, full_name, email, pass_word, role, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id;
	`
	if user.Id == uuid.Nil {
		user.Id = uuid.New()
	}
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}

	row := r.db.QueryRow(query, user.Id, user.Fullname, user.Email, user.Password, user.Role, user.CreatedAt)
	if err := row.Scan(&user.Id); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepo) FindByEmail(email string) (*models.Users, error) {
	query := `SELECT id, full_name, email, pass_word, role, created_at FROM users WHERE email = $1`
	row := r.db.QueryRow(query, email)

	var user models.Users
	if err := row.Scan(&user.Id, &user.Fullname, &user.Email, &user.Password, &user.Role, &user.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}
