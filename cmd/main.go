package main

import (
	"fmt"
	"net/http"

	"github.com/tktanisha/booking_system/internal/api/router"
	"github.com/tktanisha/booking_system/internal/api/routes"
	"github.com/tktanisha/booking_system/internal/config"
	"github.com/tktanisha/booking_system/internal/db"
	"github.com/tktanisha/booking_system/internal/initializer"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Database connection
	dbURL := config.GetDBURL()
	database, err := db.InitPostgres(dbURL)
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		return
	}
	defer database.Close()

	// Running migrations
	// err = db.RunMigrations(database, "./internal/db/tables.sql")
	// if err != nil {
	// 	fmt.Printf("Failed to run migrations: %v\n", err)
	// 	return
	// }

	// Initializing services
	initializer.Initialize(database)

	// Setting routes
	mux := http.NewServeMux()
	r := router.NewMuxRouter(mux)
	routes.RegisterAllRoutes(r,
		routes.RegisterAuthRoutes,
		routes.RegisterBookingRoutes,
		routes.RegisterHotelRoutes,
		routes.RegisterRoomRoutes,
	)

	// Starting server
	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", mux)
}
