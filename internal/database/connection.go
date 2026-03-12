package database

import (
	"context"
	"fmt"

	"github.com/MotiurRahmanSany/url-shrinker-api/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewConnection(cfg *config.DBConfig) (*pgxpool.Pool, error) {

	sslMode := "disable"

	if cfg.EnableSSLMode {
		sslMode = "require"
	}

	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, sslMode,
	)

	pool, err := pgxpool.New(context.Background(), connStr)

	if err != nil {
		return nil, err
	}

	return pool, nil

}
