package middlewares_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/tktanisha/booking_system/internal/api/middlewares"
	"github.com/tktanisha/booking_system/internal/constants"
	user_role "github.com/tktanisha/booking_system/internal/enums/user"
	"github.com/tktanisha/booking_system/internal/models"
	"github.com/tktanisha/booking_system/internal/utils"
)

func TestAuthMiddleware(t *testing.T) {
	testUUID := uuid.New()
	tests := []struct {
		name           string
		authHeader     string
		setupToken     func() string
		expectedStatus int
		expectNext     bool
	}{
		{
			name:           "missing authorization header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectNext:     false,
		},
		{
			name:           "invalid authorization header format",
			authHeader:     "InvalidToken",
			expectedStatus: http.StatusUnauthorized,
			expectNext:     false,
		},
		{
			name:           "invalid jwt token",
			authHeader:     "Bearer junk.token.value",
			expectedStatus: http.StatusUnauthorized,
			expectNext:     false,
		},
		{
			name: "valid jwt token",
			setupToken: func() string {
				token, _ := utils.GenerateJWT(testUUID, user_role.RoleUser)
				return token
			},
			expectedStatus: http.StatusOK,
			expectNext:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var calledNext bool
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calledNext = true
				val := r.Context().Value(constants.UserContextKey)
				if val == nil {
					t.Errorf("expected UserCtx in context, got nil")
				}
				if userCtx, ok := val.(*models.UserContext); ok {
					if userCtx.Id != testUUID && tt.expectNext {
						t.Errorf("expected UserID=%s, got %d", testUUID, userCtx.Id)
					}
					if userCtx.Role != user_role.RoleUser && tt.expectNext {
						t.Errorf("expected Role=User, got %v", userCtx.Role)
					}
				}
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("success"))
			})

			req := httptest.NewRequest(http.MethodGet, "/protected", nil)

			// assign header if set
			if tt.setupToken != nil {
				token := tt.setupToken()
				tt.authHeader = "Bearer " + token
			}
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			rr := httptest.NewRecorder()
			handler := middlewares.AuthMiddleware(nextHandler)

			handler.ServeHTTP(rr, req)

			res := rr.Result()
			defer res.Body.Close()

			body, _ := ioutil.ReadAll(res.Body)
			fmt.Printf("%s -> %d %s\n", tt.name, res.StatusCode, string(body))

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, res.StatusCode)
			}
			if tt.expectNext && !calledNext {
				t.Errorf("expected next handler to be called, but it was not")
			}
			if !tt.expectNext && calledNext {
				t.Errorf("expected next handler NOT to be called, but it was")
			}
		})
	}
}
