package auth

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestHMACValidation(t *testing.T) {
	secret := "my-secret-key-123"
	v := NewHMACValidator(secret)

	// Create a valid token
	now := time.Now()
	claims := &Claims{
		UserID: "user_123",
		Email:  "test@example.com",
		Role:   "student",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	// Validate valid token
	validatedClaims, err := v.Validate(tokenString)
	if err != nil {
		t.Fatalf("validation failed: %v", err)
	}
	if validatedClaims.UserID != "user_123" {
		t.Errorf("expected UserID user_123, got %s", validatedClaims.UserID)
	}
	if validatedClaims.Role != "student" {
		t.Errorf("expected Role student, got %s", validatedClaims.Role)
	}

	// Validate invalid token (wrong secret)
	vInvalid := NewHMACValidator("wrong-secret")
	_, err = vInvalid.Validate(tokenString)
	if err == nil {
		t.Error("expected validation to fail for wrong secret, but it succeeded")
	}

	// Test bearer token extraction
	validatedClaims2, err := v.Validate("Bearer " + tokenString)
	if err != nil {
		t.Fatalf("validation with Bearer prefix failed: %v", err)
	}
	if validatedClaims2.UserID != "user_123" {
		t.Errorf("expected UserID user_123, got %s", validatedClaims2.UserID)
	}
}

func TestContextHelpers(t *testing.T) {
	ctx := context.Background()
	ctx = WithUserID(ctx, "user_999")
	ctx = WithRole(ctx, "admin")

	uid, ok := UserIDFromContext(ctx)
	if !ok || uid != "user_999" {
		t.Errorf("UserIDFromContext returned %s, %t; want user_999, true", uid, ok)
	}

	role, ok := RoleFromContext(ctx)
	if !ok || role != "admin" {
		t.Errorf("RoleFromContext returned %s, %t; want admin, true", role, ok)
	}
}
