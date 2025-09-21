package booking_validators

import (
	"errors"

	"github.com/google/uuid"
)

func ValidateCancelBooking(bookingId uuid.UUID) error {
	if bookingId == uuid.Nil {
		return errors.New("booking_id is required")
	}
	return nil
}
