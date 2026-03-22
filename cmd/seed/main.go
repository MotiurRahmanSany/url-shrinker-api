package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/MotiurRahmanSany/url-shrinker-api/internal/config"
	"github.com/MotiurRahmanSany/url-shrinker-api/internal/database"
	"github.com/MotiurRahmanSany/url-shrinker-api/internal/db"
	"github.com/MotiurRahmanSany/url-shrinker-api/internal/repository"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	cfg := config.GetConfig()

	adminEmail := os.Getenv("SEED_ADMIN_EMAIL")
	adminPassword := os.Getenv("SEED_ADMIN_PASSWORD")

	if adminEmail == "" || adminPassword == "" {
		slog.Error("SEED_ADMIN_EMAIL and SEED_ADMIN_PASSWORD are required")
		os.Exit(1)
	}

	pool, err := database.NewConnection(cfg.Db)
	if err != nil {
		slog.Error("database connection failed", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	queries := db.New(pool)
	userRepo := repository.NewUserRepository(queries)

	ctx := context.Background()

	_, err = userRepo.GetUserByEmail(ctx, adminEmail)
	if err == nil {
		fmt.Printf("admin user already exists: %s\n", adminEmail)
		return
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		slog.Error("failed to check existing admin user", "err", err)
		os.Exit(1)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("failed to hash admin password", "err", err)
		os.Exit(1)
	}

	created, err := userRepo.CreateUser(ctx, adminEmail, string(hash), "admin")
	if err != nil {
		slog.Error("failed to create admin user", "err", err)
		os.Exit(1)
	}

	fmt.Printf("admin user created successfully: %s (id=%s)\n", created.Email, created.ID)
}
