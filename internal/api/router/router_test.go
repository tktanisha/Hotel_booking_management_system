package router_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tktanisha/booking_system/internal/api/router"
)

func TestMuxRouter_HandleFunc(t *testing.T) {
	mux := http.NewServeMux()
	r := router.NewMuxRouter(mux)

	called := false
	handler := func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}

	middlewareCalled := false
	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			middlewareCalled = true
			next.ServeHTTP(w, r)
		})
	}

	r.HandleFunc("/test", handler, middleware)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
	if !called {
		t.Errorf("expected handler to be called")
	}
	if !middlewareCalled {
		t.Errorf("expected middleware to be called")
	}
}

func TestNewMuxRouter(t *testing.T) {
	mux := http.NewServeMux()
	r := router.NewMuxRouter(mux)
	if r == nil {
		t.Errorf("expected non-nil router")
	}
	if r == nil {
		t.Errorf("router is nil")
	}
}
