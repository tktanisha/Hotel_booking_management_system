package routes

import (
	"github.com/tktanisha/booking_system/internal/api/router"
)

func RegisterAllRoutes(r *router.MuxRouter, routeFuncs ...func(*router.MuxRouter)) {
	for _, register := range routeFuncs {
		register(r)
	}
}
