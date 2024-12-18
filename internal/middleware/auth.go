package middleware

import (
	"context"
	"net/http"
	"strings"

	"YALP/internal/service"
	"YALP/pkg/response"
)

type ContextKey string
const(ContextUserIDKey ContextKey = "user_id")
func Auth(authSvc service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.JSONError(w, http.StatusUnauthorized, "missing authorization header")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				response.JSONError(w, http.StatusUnauthorized, "invalid token format")
				return
			}

			token := parts[1]
			userID, err := authSvc.ValidateJWT(token) // Use AuthService for validation
			if err != nil {
				response.JSONError(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			// Add user ID to context using custom type key
			ctx := context.WithValue(r.Context(), ContextUserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
