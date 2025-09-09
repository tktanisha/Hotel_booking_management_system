package models

import (
	"time"

	"github.com/google/uuid"
	user_role "github.com/tktanisha/booking_system/internal/enums/user"
)

type Users struct {
	Id        uuid.UUID          `json:"id"`
	Fullname  string             `json:"fullname"`
	Email     string             `json:"email"`
	Password  string             `json:"password"`
	Role      user_role.UserRole `json:"role"`
	CreatedAt time.Time          `json:"created_at"`
}

type UserContext struct {
	Id   uuid.UUID
	Role user_role.UserRole
}
