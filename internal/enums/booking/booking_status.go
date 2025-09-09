package booking_status

type BookingStatus string

const (
	StatusConfirmed  BookingStatus = "confirmed"
	StatusCancelled  BookingStatus = "cancelled"
	StatusCheckedOut BookingStatus = "checked_out"
)
