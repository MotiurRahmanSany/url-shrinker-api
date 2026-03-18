package main

import (
	"fmt"
	"net/http"
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
)

func serve(config *config.Config) {
	// connecting to database
	pool, err := database.NewConnection(config.Db)
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err.Error())
		return
	}
	defer pool.Close()
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
	// clickRepo := repository.NewClickRepository(queries)
	
	authService := service.NewAuthService(userRepo, tokenRepo, jwtManager)
	urlService := service.NewUrlService(urlRepo, redisCache)

	healthHandler := handlers.NewHealthHandler()
	authHandler := handlers.NewAuthHandler(authService)
	urlHandler := handlers.NewUrlHandler(urlService)

	mux := router.Setup(
		jwtManager,
		healthHandler,
		authHandler,
		urlHandler,
	)

	loggedMux := middleware.Logger(mux)

	fmt.Printf("Server is running on port %d\n", config.HttpPort)
	fmt.Printf("Base Url is: http:localhost:%d\n", config.HttpPort)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.HttpPort), loggedMux); err != nil {
		fmt.Printf("Error starting server: %v\n", err.Error())
	}

}
