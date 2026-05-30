// Package auth provides JWT validation middleware for the API gateway.
package middleware

import (
	"context"
	"net/http"
	"strings"

	"go.uber.org/zap"

	pkgauth "github.com/skillofide/pkg/auth"
)

type contextKey string

const (
	CtxUserID contextKey = "user_id"
	CtxRole   contextKey = "role"
)

// Auth returns an HTTP middleware that validates JWTs and injects user claims into the request context.
// Requests to public paths (listed in publicPaths) bypass authentication.
func Auth(validator *pkgauth.Validator, log *zap.Logger, publicPaths ...string) func(http.Handler) http.Handler {
	publicSet := make(map[string]bool, len(publicPaths))
	for _, p := range publicPaths {
		publicSet[p] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Allow public paths
			if publicSet[r.URL.Path] {
				next.ServeHTTP(w, r)
				return
			}

			// Extract token
			token := r.Header.Get("Authorization")
			if token == "" {
				token = r.URL.Query().Get("token") // fallback for WebSocket
			}
			if token == "" {
				http.Error(w, `{"error":"missing authorization token"}`, http.StatusUnauthorized)
				return
			}

			// Validate
			claims, err := validator.Validate(strings.TrimPrefix(token, "Bearer "))
			if err != nil {
				log.Warn("invalid token", zap.Error(err))
				http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
				return
			}

			// Inject into context
			ctx := context.WithValue(r.Context(), CtxUserID, claims.UserID)
			ctx = context.WithValue(ctx, CtxRole, claims.Role)

			// Also set headers for downstream services (WebSocket, etc.)
			r.Header.Set("X-User-ID", claims.UserID)
			r.Header.Set("X-User-Role", claims.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// UserIDFromContext extracts the user ID from the request context.
func UserIDFromContext(ctx context.Context) string {
	v, _ := ctx.Value(CtxUserID).(string)
	return v
}

// RoleFromContext extracts the user role from the request context.
func RoleFromContext(ctx context.Context) string {
	v, _ := ctx.Value(CtxRole).(string)
	return v
}

// CORS returns a middleware that adds CORS headers for the GraphQL endpoint.
func CORS(allowedOrigins string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if allowedOrigins == "*" || strings.Contains(allowedOrigins, origin) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			} else {
				w.Header().Set("Access-Control-Allow-Origin", allowedOrigins)
			}
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
