package initializer

import (
	"github.com/tktanisha/booking_system/internal/db"
	"github.com/tktanisha/booking_system/internal/repository/booking_repo"
	"github.com/tktanisha/booking_system/internal/repository/hotel_repo"
	"github.com/tktanisha/booking_system/internal/repository/room_repo"
	"github.com/tktanisha/booking_system/internal/repository/user_repo"
	"github.com/tktanisha/booking_system/internal/services/auth_service"
	"github.com/tktanisha/booking_system/internal/services/booking_service"
	"github.com/tktanisha/booking_system/internal/services/hotel_service"
	"github.com/tktanisha/booking_system/internal/services/room_service"
)

var (
	userRepo    user_repo.UserRepoInterface
	bookingRepo booking_repo.BookingRepoInterface
	roomRepo    room_repo.RoomRepoInterface
	hotelRepo   hotel_repo.HotelRepositoryInterface

	AuthService    auth_service.AuthServiceInterface
	RoomService    room_service.RoomServiceInterface
	BookingService booking_service.BookingServiceInterface
	HotelService   hotel_service.HotelServiceInterface
)

func Initialize(db db.DB) {
	userRepo = user_repo.NewUserRepo(db)
	bookingRepo = booking_repo.NewBookingRepo(db)
	hotelRepo = hotel_repo.NewHotelRepo(db)
	roomRepo = room_repo.NewRoomRepo(db)

	AuthService = auth_service.NewAuthService(userRepo)
	RoomService = room_service.NewRoomService(roomRepo)
	BookingService = booking_service.NewBookingService(bookingRepo, RoomService)
	HotelService = hotel_service.NewHotelService(hotelRepo)
}
