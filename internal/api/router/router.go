package router

import "net/http"

type MuxRouter struct {
	mux *http.ServeMux
}

func NewMuxRouter(mux *http.ServeMux) *MuxRouter {
	return &MuxRouter{mux: mux}
}

func (r *MuxRouter) HandleFunc(pattern string, handler http.HandlerFunc, middlewares ...func(http.Handler) http.Handler) {
	var h http.Handler = handler

	// Apply middlewares in reverse order (first provided, executed first)
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	r.mux.Handle(pattern, h)
}
