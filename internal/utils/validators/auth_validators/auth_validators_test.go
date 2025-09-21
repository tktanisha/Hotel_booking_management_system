package auth_validators_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/tktanisha/booking_system/internal/utils/validators/auth_validators"
	"github.com/tktanisha/booking_system/internal/utils/validators/payloads"
)

func TestValidateFullName(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		errorMsg    string
	}{
		{"valid full name", "John Doe", false, ""},
		{"empty full name", "", true, "full name cannot be empty"},
		{"only first name", "John", true, "please enter at least first and last name"},
		{"extra spaces", "   ", true, "full name cannot be empty"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := auth_validators.ValidateFullName(tt.input)
			if tt.expectError {
				if err == nil || err.Error() != tt.errorMsg {
					t.Errorf("expected error %q, got %v", tt.errorMsg, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		errorMsg    string
	}{
		{"valid email", "test@example.com", false, ""},
		{"empty email", "", true, "email cannot be empty"},
		{"invalid email format", "invalidemail", true, "invalid email format"},
		{"uppercase email", "TEST@EXAMPLE.COM", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := auth_validators.ValidateEmail(tt.input)
			if tt.expectError {
				if err == nil || err.Error() != tt.errorMsg {
					t.Errorf("expected error %q, got %v", tt.errorMsg, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		errorMsg    string
	}{
		{"valid password", "Aa1!abc", false, ""},
		{"empty password", "", true, "password cannot be empty"},
		{"short password", "Aa1!", true, "password must be at least 6 characters"},
		{"no digit", "Aa!aaaa", true, "password must contain at least one digit"},
		{"no special char", "Aa1aaaa", true, "password must contain at least one special character"},
		{"no uppercase", "aa1!aaa", true, "password must contain at least one uppercase letter"},
		{"no lowercase", "AA1!AAA", true, "password must contain at least one lowercase letter"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := auth_validators.ValidatePassword(tt.input)
			if tt.expectError {
				if err == nil || err.Error() != tt.errorMsg {
					t.Errorf("expected error %q, got %v", tt.errorMsg, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestLoginValidate(t *testing.T) {
	validPayload := payloads.LoginRequest{
		Email:    "test@example.com",
		Password: "Aa1!abc",
	}

	tests := []struct {
		name        string
		body        any
		expectError bool
		errorMsg    string
	}{
		{"valid login", validPayload, false, ""},
		{"invalid JSON", "{invalid json}", true, "invalid character"},
		{"invalid email", payloads.LoginRequest{Email: "bademail", Password: "Aa1!abc"}, true, "invalid email format"},
		{"invalid password", payloads.LoginRequest{Email: "test@example.com", Password: "weak"}, true, "password must be at least 6 characters"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			switch v := tt.body.(type) {
			case string:
				body = []byte(v)
			default:
				body, _ = json.Marshal(v)
			}

			req := httptest.NewRequest("POST", "/", bytes.NewReader(body))
			_, err := auth_validators.LoginValidate(req)
			if tt.expectError {
				if err == nil || !contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error containing %q, got %v", tt.errorMsg, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestRegisterValidate(t *testing.T) {
	validPayload := payloads.RegisterRequest{
		Fullname: "John Doe",
		Email:    "test@example.com",
		Password: "Aa1!abc",
	}

	tests := []struct {
		name        string
		body        interface{}
		expectError bool
		errorMsg    string
	}{
		{"valid register", validPayload, false, ""},
		{"invalid JSON", "{invalid json", true, "invalid character"},
		{"missing full name", payloads.RegisterRequest{Fullname: "", Email: "test@example.com", Password: "Aa1!abc"}, true, "full name cannot be empty"},
		{"invalid email", payloads.RegisterRequest{Fullname: "John Doe", Email: "bademail", Password: "Aa1!abc"}, true, "invalid email format"},
		{"invalid password", payloads.RegisterRequest{Fullname: "John Doe", Email: "test@example.com", Password: "weak"}, true, "password must be at least 6 characters"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			switch v := tt.body.(type) {
			case string:
				body = []byte(v)
			default:
				body, _ = json.Marshal(v)
			}

			req := httptest.NewRequest("POST", "/", bytes.NewReader(body))
			_, err := auth_validators.RegisterValidate(req)
			if tt.expectError {
				if err == nil || !contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error containing %q, got %v", tt.errorMsg, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || bytes.Contains([]byte(s), []byte(substr)))
}
