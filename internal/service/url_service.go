package service

import (
	"context"
	"time"

	"github.com/MotiurRahmanSany/url-shrinker-api/internal/cache"
	"github.com/MotiurRahmanSany/url-shrinker-api/internal/domain"
	"github.com/MotiurRahmanSany/url-shrinker-api/internal/repository"
)

type UrlService interface {
	CreateShortURL(ctx context.Context, shortCode, orgCode string, userID string, expiresAt time.Time, maxClicks int32) (domain.Url, error)
	DeactivateURL(ctx context.Context, id int64) error
	GetURLByShortCode(ctx context.Context, shortCode string) (domain.Url, error)
	GetURLsByUserID(ctx context.Context, userID string, limit, offset int32) ([]domain.Url, error)
	UpdateURL(ctx context.Context, id int64, orgURL string, expiresAt time.Time, maxClicks int32) (domain.Url, error)
}

type urlService struct {
	repo repository.UrlRepository
	cache  cache.Cache
}

func NewUrlService(repo repository.UrlRepository, redisCache cache.Cache) UrlService {
	return &urlService{repo: repo, cache: redisCache}
}

