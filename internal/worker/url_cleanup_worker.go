package worker

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/MotiurRahmanSany/url-shrinker-api/internal/repository"
)

var ErrInvalidCleanupInterval = errors.New("invalid cleanup interval, must be a positive duration")

func StartExpiredURLCleanupWorker(
	ctx context.Context,
	urlRepo repository.UrlRepository,
	interval time.Duration,
) error {
	if interval <= 0 {
		return ErrInvalidCleanupInterval
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	slog.Info("expired URL cleanup worker started", "interval", interval.String())

	for {
		select {
		case <-ctx.Done():
			slog.Info("expired URL cleanup worker stopping due to context cancellation")
			return nil
		case <-ticker.C:
			deletedCount, err := urlRepo.DeleteExpiredURLs(ctx)
			if err != nil {
				slog.Error("error deleting expired URLs", "error", err)
				continue
			}
			if deletedCount > 0 {
				slog.Info("expired URL cleanup completed", "deleted_count", deletedCount)
			}
		}
	}

}
