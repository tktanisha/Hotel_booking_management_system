package utils_test

import (
	"testing"

	"github.com/tktanisha/booking_system/internal/utils"
)

func TestHashPassword(t *testing.T) {
	password := "mySecret123"

	hashed, err := utils.HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if hashed == "" {
		t.Fatalf("expected a hashed password, got empty string")
	}

	// Ensure that the hashed password is not the same as the original
	if hashed == password {
		t.Errorf("hashed password should not match original password")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "mySecret123"
	wrongPassword := "wrongPassword"

	hashed, err := utils.HashPassword(password)
	if err != nil {
		t.Fatalf("unexpected error hashing password: %v", err)
	}

	t.Run("correct password", func(t *testing.T) {
		if !utils.CheckPasswordHash(password, hashed) {
			t.Errorf("expected password to match hash")
		}
	})

	t.Run("incorrect password", func(t *testing.T) {
		if utils.CheckPasswordHash(wrongPassword, hashed) {
			t.Errorf("expected password not to match hash")
		}
	})
}
