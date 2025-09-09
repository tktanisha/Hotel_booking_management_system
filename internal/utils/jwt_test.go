package utils_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	user_role "github.com/tktanisha/booking_system/internal/enums/user"
	"github.com/tktanisha/booking_system/internal/utils"
)

func TestGenerateAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	role := user_role.RoleManager

	// Generate a JWT token
	token, err := utils.GenerateJWT(userID, role)
	if err != nil {
		t.Fatalf("expected no error generating JWT, got %v", err)
	}
	if token == "" {
		t.Fatalf("expected a token string, got empty string")
	}

	// Validate the token
	claims, err := utils.ValidateJWT(token)
	if err != nil {
		t.Fatalf("expected no error validating JWT, got %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("expected userID %v, got %v", userID, claims.UserID)
	}
	if claims.Role != role {
		t.Errorf("expected role %v, got %v", role, claims.Role)
	}

	// Check expiry is roughly 24 hours from issued time
	expectedExpiry := claims.IssuedAt.Time.Add(24 * time.Hour)
	if claims.ExpiresAt.Time.Sub(expectedExpiry) > time.Second {
		t.Errorf("expected expiry %v, got %v", expectedExpiry, claims.ExpiresAt.Time)
	}
}

func TestValidateJWT_InvalidToken(t *testing.T) {
	invalidToken := "this.is.not.a.valid.token"

	_, err := utils.ValidateJWT(invalidToken)
	if err == nil {
		t.Fatalf("expected error validating invalid token, got nil")
	}
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	userID := uuid.New()
	role := user_role.RoleUser

	// Manually create expired token
	expiredClaims := &utils.Claims{
		UserID:           userID,
		Role:             role,
		RegisteredClaims: utils.Claims{}.RegisteredClaims,
	}
	expiredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(-1 * time.Hour))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	tokenStr, _ := token.SignedString([]byte("secret_key"))

	_, err := utils.ValidateJWT(tokenStr)
	if err == nil {
		t.Fatalf("expected error validating expired token, got nil")
	}
}
