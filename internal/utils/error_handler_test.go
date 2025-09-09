package utils_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tktanisha/booking_system/internal/utils"
)

func TestWriteErrorResponse(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		errMsg         string
		details        string
		expectedStatus bool
	}{
		{
			name:           "with details",
			statusCode:     http.StatusBadRequest,
			errMsg:         "Invalid request",
			details:        "Missing required field",
			expectedStatus: false,
		},
		{
			name:           "without details",
			statusCode:     http.StatusUnauthorized,
			errMsg:         "Unauthorized",
			details:        "",
			expectedStatus: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			utils.WriteErrorResponse(w, tt.statusCode, tt.errMsg, tt.details)

			resp := w.Result()
			if resp.StatusCode != tt.statusCode {
				t.Errorf("expected status code %d, got %d", tt.statusCode, resp.StatusCode)
			}

			contentType := resp.Header.Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("expected Content-Type application/json, got %s", contentType)
			}

			var body utils.ErrorResponse
			if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
				t.Fatalf("failed to decode response body: %v", err)
			}

			if body.Status != tt.expectedStatus {
				t.Errorf("expected Status %v, got %v", tt.expectedStatus, body.Status)
			}

			if body.StatusCode != tt.statusCode {
				t.Errorf("expected StatusCode %d, got %d", tt.statusCode, body.StatusCode)
			}

			if body.Error != tt.errMsg {
				t.Errorf("expected Error %s, got %s", tt.errMsg, body.Error)
			}

			if body.Details != tt.details {
				t.Errorf("expected Details %s, got %s", tt.details, body.Details)
			}
		})
	}
}
