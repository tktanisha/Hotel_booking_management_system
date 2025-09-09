package permissions

import (
	"testing"

	user_role "github.com/tktanisha/booking_system/internal/enums/user"
	"github.com/tktanisha/booking_system/internal/models"
)

func TestIsManager(t *testing.T) {
	tests := []struct {
		name     string
		userCtx  *models.UserContext
		expected bool
	}{
		{
			name:     "Nil UserContext",
			userCtx:  nil,
			expected: false,
		},
		{
			name: "Non-Manager Role",
			userCtx: &models.UserContext{
				Role: user_role.RoleUser,
			},
			expected: false,
		},
		{
			name: "Manager Role",
			userCtx: &models.UserContext{
				Role: user_role.RoleManager,
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsManager(tt.userCtx)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
