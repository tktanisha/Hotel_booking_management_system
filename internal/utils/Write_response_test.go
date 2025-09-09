package utils_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tktanisha/booking_system/internal/utils"
)

func TestWriteSuccessResponse(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		message        string
		data           interface{}
		expectedStatus bool
	}{
		{
			name:           "with data",
			statusCode:     http.StatusOK,
			message:        "Success",
			data:           map[string]string{"key": "value"},
			expectedStatus: true,
		},
		{
			name:           "without data",
			statusCode:     http.StatusCreated,
			message:        "Created",
			data:           nil,
			expectedStatus: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			utils.WriteSuccessResponse(w, tt.statusCode, tt.message, tt.data)

			resp := w.Result()
			if resp.StatusCode != tt.statusCode {
				t.Errorf("expected status code %d, got %d", tt.statusCode, resp.StatusCode)
			}

			contentType := resp.Header.Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("expected Content-Type application/json, got %s", contentType)
			}

			var body utils.WriteResponse
			if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
				t.Fatalf("failed to decode response body: %v", err)
			}

			if body.Status != tt.expectedStatus {
				t.Errorf("expected Status %v, got %v", tt.expectedStatus, body.Status)
			}

			if body.StatusCode != tt.statusCode {
				t.Errorf("expected StatusCode %d, got %d", tt.statusCode, body.StatusCode)
			}

			if body.Message != tt.message {
				t.Errorf("expected Message %s, got %s", tt.message, body.Message)
			}

			if tt.data != nil {
				if body.Data == nil {
					t.Errorf("expected Data to be %v, got nil", tt.data)
				}
			} else if body.Data != nil {
				t.Errorf("expected Data to be nil, got %v", body.Data)
			}
		})
	}
}
