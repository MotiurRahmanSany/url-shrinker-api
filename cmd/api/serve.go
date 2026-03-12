package main

import (
	"fmt"
	"net/http"

	"github.com/MotiurRahmanSany/url-shrinker-api/internal/api/middleware"
	"github.com/MotiurRahmanSany/url-shrinker-api/internal/cache"
	"github.com/MotiurRahmanSany/url-shrinker-api/internal/config"
	"github.com/MotiurRahmanSany/url-shrinker-api/internal/database"
)

func serve(config *config.Config) {
	// connecting to database
	pool, err := database.NewConnection(config.Db)
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err.Error())
		return
	}
	defer pool.Close()
	// queries := db.New(pool)

	_ = cache.NewRedisCache(
		fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
		config.Redis.Password,
		config.Redis.DB,
	)

	mux := http.NewServeMux()

	loggedMux := middleware.Logger(mux)

	fmt.Printf("Server is running on port %d\n", config.HttpPort)
	fmt.Printf("Base Url is: http:localhost:%d\n", config.HttpPort)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.HttpPort), loggedMux); err != nil {
		fmt.Printf("Error starting server: %v\n", err.Error())
	}

}
