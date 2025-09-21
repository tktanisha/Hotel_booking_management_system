package auth_validators

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/tktanisha/booking_system/internal/utils/validators/payloads"
)

func LoginValidate(r *http.Request) (*payloads.LoginRequest, error) {
	var payload payloads.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return nil, err
	}
	if err := ValidateEmail(payload.Email); err != nil {
		return nil, err
	}
	if err := ValidatePassword(payload.Password); err != nil {
		return nil, err
	}
	return &payload, nil
}

func ValidateFullName(input string) error {
	if strings.TrimSpace(input) == "" {
		return errors.New("full name cannot be empty")
	}
	if len(strings.Fields(input)) < 2 {
		return errors.New("please enter at least first and last name")
	}
	return nil
}

func ValidateEmail(input string) error {
	if strings.TrimSpace(input) == "" {
		return errors.New("email cannot be empty")
	}
	regex := "^[a-z0-9._%+\\-]+@[a-z0-9.\\-]+\\.[a-z]{2,}$"
	matched, _ := regexp.MatchString(regex, strings.ToLower(input))
	if !matched {
		return errors.New("invalid email format")
	}
	return nil
}

func ValidatePassword(input string) error {
	if strings.TrimSpace(input) == "" {
		return errors.New("password cannot be empty")
	}
	if len(input) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	hasDigit := regexp.MustCompile("[0-9]").MatchString
	if !hasDigit(input) {
		return errors.New("password must contain at least one digit")
	}
	hasSpecialChar := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString
	if !hasSpecialChar(input) {
		return errors.New("password must contain at least one special character")
	}
	hasUpperCase := regexp.MustCompile(`[A-Z]`).MatchString
	if !hasUpperCase(input) {
		return errors.New("password must contain at least one uppercase letter")
	}
	hasLowerCase := regexp.MustCompile(`[a-z]`).MatchString
	if !hasLowerCase(input) {
		return errors.New("password must contain at least one lowercase letter")
	}
	return nil
}
