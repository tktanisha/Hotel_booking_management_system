package booking_service

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/tktanisha/booking_system/internal/api/validators/payloads"
	booking_status "github.com/tktanisha/booking_system/internal/enums/booking"
	"github.com/tktanisha/booking_system/internal/models"
	"github.com/tktanisha/booking_system/internal/repository/booking_repo"
	"github.com/tktanisha/booking_system/internal/services/room_service"
)

type BookingService struct {
	BookingRepo booking_repo.BookingRepoInterface
	RoomService room_service.RoomServiceInterface
}

func NewBookingService(bookingRepo booking_repo.BookingRepoInterface, roomService room_service.RoomServiceInterface) *BookingService {
	return &BookingService{
		BookingRepo: bookingRepo,
		RoomService: roomService,
	}
}

func (b *BookingService) CancelBooking(bookingId uuid.UUID) (*models.Bookings, error) {
	booking, err := b.BookingRepo.GetBookingById(bookingId)
	if err != nil {
		return nil, err
	}

	if booking.CheckIn.Before(time.Now()) {
		return nil, errors.New("Cannot cancel booking after check-in date")
	}

	if booking.Status == booking_status.StatusCancelled {
		return nil, errors.New("Booking is already cancelled")
	}

	booking.Status = booking_status.StatusCancelled

	// increase room quantity back
	bookedRooms, err := b.BookingRepo.GetBookedRoomsByBookingId(bookingId)
	if err != nil {
		return nil, err
	}
	var roomsPayload []*payloads.RoomPayload
	for _, bookedRoom := range bookedRooms {
		roomPayload := &payloads.RoomPayload{
			RoomType: bookedRoom.RoomType,
			Quantity: bookedRoom.RoomQuantity,
		}
		roomsPayload = append(roomsPayload, roomPayload)
	}
	for _, room := range roomsPayload {
		if _, err := b.RoomService.IncreaseRoomQuantity(room, booking.HotelId); err != nil {
			return nil, err
		}
	}

	if err := b.BookingRepo.Save(booking); err != nil {
		return nil, err
	}

	return booking, nil
}

func (b *BookingService) CreateBooking(userCtx *models.UserContext, payload *payloads.BookingPayload) (*models.Bookings, error) {
	rooms := payload.Rooms
	hotelId := payload.HotelId

	for _, room := range rooms {
		if !b.RoomService.IsAvailable(room, hotelId) {
			return nil, errors.New("Not Available")
		}
	}

	booking := models.Bookings{
		Id:        uuid.New(),
		UserId:    userCtx.Id,
		HotelId:   hotelId,
		CheckIn:   payload.CheckIn,
		CheckOut:  payload.CheckOut,
		Status:    booking_status.StatusConfirmed,
		CreatedAt: time.Now(),
	}

	// Reduce available room quantity before booking
	for _, room := range rooms {
		if err := b.RoomService.ReduceRoomQuantity(room, hotelId); err != nil {
			return nil, err
		}
	}

	// Prepare booked rooms
	roomsData := make([]*models.BookedRooms, 0)
	for _, room := range rooms {
		bookedRoom := &models.BookedRooms{
			Id:           uuid.New(),
			BookingId:    booking.Id,
			RoomType:     room.RoomType,
			RoomQuantity: room.Quantity,
			CreatedAt:    time.Now(),
		}
		roomsData = append(roomsData, bookedRoom)
	}

	// Create booking + booked rooms in one repo method
	savedBooking, err := b.BookingRepo.CreateBookingWithRooms(&booking, roomsData)
	if err != nil {
		return nil, err
	}

	return savedBooking, nil
}

func (b *BookingService) CheckoutBooking(bookingId uuid.UUID) (*models.Bookings, error) {
	booking, err := b.BookingRepo.GetBookingById(bookingId)
	if err != nil {
		return nil, err
	}

	if booking.Status != booking_status.StatusConfirmed {
		return nil, errors.New("Only confirmed bookings can be checked out")
	}

	// increase room quantity back
	bookedRooms, err := b.BookingRepo.GetBookedRoomsByBookingId(bookingId)
	if err != nil {
		return nil, err
	}

	var roomsPayload []*payloads.RoomPayload
	for _, bookedRoom := range bookedRooms {
		roomPayload := &payloads.RoomPayload{
			RoomType: bookedRoom.RoomType,
			Quantity: bookedRoom.RoomQuantity,
		}
		roomsPayload = append(roomsPayload, roomPayload)
	}

	for _, room := range roomsPayload {
		if _, err := b.RoomService.IncreaseRoomQuantity(room, booking.HotelId); err != nil {
			return nil, err
		}
	}

	booking.Status = booking_status.StatusCheckedOut
	if err := b.BookingRepo.Save(booking); err != nil {
		return nil, err
	}

	return booking, nil
}
