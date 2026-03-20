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
	urlHandler *handlers.UrlHandler,
) *http.ServeMux {
	mux := http.NewServeMux()

	// Public Routes

	// health
	mux.HandleFunc("GET /health", healthHandler.Check)

	// auth
	mux.HandleFunc("POST /auth/register", authHandler.Register)
	mux.HandleFunc("POST /auth/login", authHandler.Login)
	mux.HandleFunc("POST /auth/refresh", authHandler.RefreshToken)

	// URL redirection (public)
	mux.HandleFunc("GET /{code}", http.HandlerFunc(urlHandler.RedirectURL))

	// Protected Routes - require authentication

	authMw := middleware.AuthMiddleware(jwtManager)
	// adminOnly := middleware.AdminOnly

	mux.Handle("GET /auth/me", authMw(http.HandlerFunc(authHandler.GetMe)))
	mux.Handle("POST /auth/logout", authMw(http.HandlerFunc(authHandler.Logout)))

	// URL routes
	mux.Handle("POST /urls", authMw(http.HandlerFunc(urlHandler.CreateShortURL)))
	mux.Handle("GET /urls", authMw(http.HandlerFunc(urlHandler.ListMyURLs)))
	mux.Handle("GET /urls/{code}", authMw(http.HandlerFunc(urlHandler.GetURLDetails)))
	mux.Handle("PATCH /urls/{code}", authMw(http.HandlerFunc(urlHandler.UpdateURL)))
	mux.Handle("DELETE /urls/{code}", authMw(http.HandlerFunc(urlHandler.DeactivateURL)))
	mux.Handle("GET /urls/{code}/stats", authMw(http.HandlerFunc(urlHandler.GetURLStats)))

	return mux
}
