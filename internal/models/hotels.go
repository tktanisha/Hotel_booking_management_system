package models

import (
	"time"

	"github.com/google/uuid"
)

type Hotels struct {
	Id        uuid.UUID `json:"id"`
	ManagerId uuid.UUID `json:"manager_id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
}
