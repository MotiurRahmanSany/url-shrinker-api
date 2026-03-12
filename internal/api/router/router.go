package router

import (
	"net/http"

	"github.com/MotiurRahmanSany/url-shrinker-api/internal/api/handlers"
	"github.com/MotiurRahmanSany/url-shrinker-api/internal/api/middleware"
	"github.com/MotiurRahmanSany/url-shrinker-api/internal/auth"
)

func Setup(
	jwtManager *auth.JWTManager,
	healthHandler *handlers.HealthHandler,
	authHandler *handlers.AuthHandler,
) *http.ServeMux {
	mux := http.NewServeMux()

	// Public Routes

	// health
	mux.HandleFunc("GET /health", healthHandler.Check)

	// auth
	mux.HandleFunc("POST /auth/register", authHandler.Register)
	mux.HandleFunc("POST /auth/login", authHandler.Login)
	mux.HandleFunc("POST /auth/refresh", authHandler.RefreshToken)

	// Protected Routes - require authentication

	authMw := middleware.AuthMiddleware(jwtManager)
	// adminOnly := middleware.AdminOnly

	mux.Handle("GET /auth/me", authMw(http.HandlerFunc(authHandler.GetMe)))
	mux.Handle("POST /auth/logout", authMw(http.HandlerFunc(authHandler.Logout)))

	return mux
}
