package utils

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
)

func GetUUIDFromParams(r *http.Request, key string) (uuid.UUID, error) {
	val := r.PathValue(key)
	if val == "" {
		return uuid.Nil, errors.New("path parameter not found")
	}

	id, err := uuid.Parse(val)

	if err != nil {
		return uuid.Nil, errors.New("invalid UUID format")
	}

	return id, nil
}
