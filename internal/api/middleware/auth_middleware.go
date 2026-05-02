package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/MotiurRahmanSany/url-shrinker-api/internal/auth"
	"github.com/MotiurRahmanSany/url-shrinker-api/internal/response"
)

type contextKey string

const (
	UserContextKey contextKey = "userID"
	RoleContextKey contextKey = "role"
)

func AuthMiddleware(jwtManager *auth.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || (!strings.HasPrefix(authHeader, "Bearer ")) {
				_ = response.Error(w, http.StatusUnauthorized, "Invalid authorization header", nil)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			claims, err := jwtManager.Verify(tokenStr)
			if err != nil {
				_ = response.Error(w, http.StatusUnauthorized, "Invalid token", nil)
				return

			}

			ctx := context.WithValue(r.Context(), UserContextKey, claims.UserID)
			ctx = context.WithValue(ctx, RoleContextKey, claims.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value(RoleContextKey).(string)
		if !ok || role != "admin" {
			_ = response.Error(w, http.StatusForbidden, "Forbidden: Admin access required", nil)
			return
		}

		next.ServeHTTP(w, r)
	})
}
