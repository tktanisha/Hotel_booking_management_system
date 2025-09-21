package routes

import "net/http"

func RegisterAllRoutes(r *http.ServeMux, routeFuncs ...func(*http.ServeMux)) {
	for _, register := range routeFuncs {
		register(r)
	}
}
