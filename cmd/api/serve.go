package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/MotiurRahmanSany/url-shrinker-api/internal/api/handlers"
	"github.com/MotiurRahmanSany/url-shrinker-api/internal/api/middleware"
	"github.com/MotiurRahmanSany/url-shrinker-api/internal/api/router"
	"github.com/MotiurRahmanSany/url-shrinker-api/internal/auth"
	"github.com/MotiurRahmanSany/url-shrinker-api/internal/cache"
	"github.com/MotiurRahmanSany/url-shrinker-api/internal/config"
	"github.com/MotiurRahmanSany/url-shrinker-api/internal/database"
	"github.com/MotiurRahmanSany/url-shrinker-api/internal/db"
	"github.com/MotiurRahmanSany/url-shrinker-api/internal/repository"
	"github.com/MotiurRahmanSany/url-shrinker-api/internal/service"
	"github.com/MotiurRahmanSany/url-shrinker-api/internal/worker"
)

var (
	ServerReadHeaderTimeout = 5 * time.Second
	ServerReadTimeout       = 5 * time.Second
	ServerWriteTimeout      = 10 * time.Second
	ServerIdleTimeout       = 30 * time.Second
	ServerShutdownTimeout   = 10 * time.Second
)

func serve(config *config.Config) {
	// connecting to database
	pool, err := database.NewConnection(config.Db)
	if err != nil {
		slog.Error("database connection failed", "err", err)
		return
	}
	queries := db.New(pool)

	redisCache := cache.NewRedisCache(
		fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
		config.Redis.Password,
		config.Redis.DB,
	)

	jwtManager := auth.NewJWTManager(config.JwtSecretKey, time.Minute*10)
	tokenRepo := repository.NewTokenRepository(queries)

	userRepo := repository.NewUserRepository(queries)
	urlRepo := repository.NewUrlRepository(queries)
	clickRepo := repository.NewClickRepository(queries)

	authService := service.NewAuthService(userRepo, tokenRepo, jwtManager)
	urlService := service.NewUrlService(urlRepo, clickRepo, redisCache)
	clickService := service.NewClickService(clickRepo)

	healthHandler := handlers.NewHealthHandler()
	authHandler := handlers.NewAuthHandler(authService)
	urlHandler := handlers.NewUrlHandler(urlService, clickService)
	clickHandler := handlers.NewClickHandler(clickService, urlService)

	mux := router.Setup(
		jwtManager,
		redisCache,
		healthHandler,
		authHandler,
		urlHandler,
		clickHandler,
	)

	loggedMux := middleware.Logger(mux)

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", config.HttpPort),
		Handler:           loggedMux,
		ReadHeaderTimeout: ServerReadHeaderTimeout,
		ReadTimeout:       ServerReadTimeout,
		WriteTimeout:      ServerWriteTimeout,
		IdleTimeout:       ServerIdleTimeout,
	}

	serverErrCh := make(chan error, 1)

	// Starting the server
	go func() {
		slog.Info("server starting", "port", config.HttpPort, "base_url", fmt.Sprintf("http://localhost:%d", config.HttpPort))

		if err := server.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				slog.Info("server gracefully stopped")
				return
			}
			serverErrCh <- err
		}
	}()

	rootCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	workerErrCh := make(chan error, 1)

	go func() {
		// Start background worker for cleaning up expired URLs every hour
		workerErrCh <- worker.StartExpiredURLCleanupWorker(rootCtx, urlRepo, time.Hour)
	}()

	select {
	case <-rootCtx.Done():
		slog.Info("shutdown signal received")
	case err := <-serverErrCh:
		slog.Error("server failed to start or crashed", "err", err)
		stop()
	case err := <-workerErrCh:
		slog.Error("cleanup worker failed", "err", err)
		stop()

	}

	// Graceful shutdown with a timeout context
	shutdownCtx, cancel := context.WithTimeout(context.Background(), ServerShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("graceful shutdown failed", "err", err)
	}

	// Close resources
	pool.Close()

	if err := redisCache.Close(); err != nil {
		slog.Error("redis close failed", "err", err)
	}

	slog.Info("server gracefully stopped")

}
