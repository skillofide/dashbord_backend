// Package auth provides JWT validation for both HMAC-SHA256 and RS256 tokens.
package auth

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	ContextKeyUserID contextKey = "user_id"
	ContextKeyRole   contextKey = "role"
)

// Claims is the standard set of claims we expect in every token.
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"` // admin | student | instructor
	jwt.RegisteredClaims
}

// Validator validates JWTs and extracts Claims.
type Validator struct {
	publicKey *rsa.PublicKey
	secret    []byte
	useHMAC   bool
}

// NewHMACValidator creates a validator using HMAC-SHA256 (HS256).
// Suitable for simple setups or when the JWT issuer shares the secret.
func NewHMACValidator(secret string) *Validator {
	return &Validator{secret: []byte(secret), useHMAC: true}
}

// NewRS256Validator creates a validator using RSA-SHA256 (RS256).
// publicKeyPEM is the PEM-encoded RSA public key from the auth provider.
func NewRS256Validator(publicKeyPEM string) (*Validator, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKeyPEM))
	if err != nil {
		return nil, fmt.Errorf("parse RSA public key: %w", err)
	}
	return &Validator{publicKey: key}, nil
}

// Validate parses and validates a JWT string, returning its Claims.
func (v *Validator) Validate(tokenString string) (*Claims, error) {
	// Strip "Bearer " prefix if present
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	claims := &Claims{}
	keyFunc := v.keyFunc()

	token, err := jwt.ParseWithClaims(tokenString, claims, keyFunc)
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}
	if !token.Valid {
		return nil, errors.New("token is invalid")
	}
	return claims, nil
}

func (v *Validator) keyFunc() jwt.Keyfunc {
	if v.useHMAC {
		return func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return v.secret, nil
		}
	}
	return func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return v.publicKey, nil
	}
}

// WithUserID stores the user ID in a context.
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, ContextKeyUserID, userID)
}

// WithRole stores the role in a context.
func WithRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, ContextKeyRole, role)
}

// UserIDFromContext retrieves the user ID injected by the auth middleware.
func UserIDFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ContextKeyUserID).(string)
	return v, ok && v != ""
}

// RoleFromContext retrieves the role injected by the auth middleware.
func RoleFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ContextKeyRole).(string)
	return v, ok && v != ""
}

// GenerateToken creates a signed JWT token for a user using HS256.
func GenerateToken(userID, email, role, secret string, expiresAfter time.Duration) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiresAfter)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
