package repository

import (
	"context"
	"time"

	"github.com/MotiurRahmanSany/url-shrinker-api/internal/db"
	"github.com/MotiurRahmanSany/url-shrinker-api/internal/domain"
	"github.com/MotiurRahmanSany/url-shrinker-api/internal/repository/postgres/mapper"
)

type UrlRepository interface {
	CreateURL(ctx context.Context, shortCode, originalUrl, userID string, expiresAt *time.Time, maxClicks *int32) (domain.Url, error)
	GetURLByShortCode(ctx context.Context, shortCode string) (domain.Url, error)
	GetURLsByUserID(ctx context.Context, userID string, limit, offset *int32) ([]domain.Url, error)
	UpdateURL(ctx context.Context, id int64, originalUrl string, expiresAt *time.Time, maxClicks *int32) (domain.Url, error)
	DeactivateURL(ctx context.Context, urlID int64) error
	DeleteExpiredURLs(ctx context.Context) (int64, error)
}

type urlRepository struct {
	q *db.Queries
}

func NewUrlRepository(q *db.Queries) UrlRepository {
	return &urlRepository{q: q}
}

func (r *urlRepository) CreateURL(ctx context.Context,
	shortCode,
	originalUrl,
	userID string,
	expiresAt *time.Time,
	maxClicks *int32,
) (domain.Url, error) {
	params, err := mapper.ToCreateUrlParams(shortCode, originalUrl, userID, expiresAt, maxClicks)

	if err != nil {
		return domain.Url{}, err
	}

	row, err := r.q.CreateURL(ctx, params)

	if err != nil {
		return domain.Url{}, err
	}

	return mapper.ToDomainUrl(row), nil

}

func (r *urlRepository) GetURLByShortCode(ctx context.Context, shortCode string) (domain.Url, error) {
	row, err := r.q.GetURLByShortCode(ctx, shortCode)

	if err != nil {
		return domain.Url{}, err
	}

	return mapper.ToDomainUrl(row), nil
}

func (r *urlRepository) GetURLsByUserID(ctx context.Context, userID string, limit, offset *int32) ([]domain.Url, error) {
	params, err := mapper.ToGetURLsByUserIDParams(userID, limit, offset)

	if err != nil {
		return nil, err
	}

	rows, err := r.q.GetURLsByUserID(ctx, params)

	if err != nil {
		return nil, err
	}

	var urls []domain.Url
	for _, row := range rows {
		urls = append(urls, mapper.ToDomainUrl(row))
	}

	return urls, nil
}

func (r *urlRepository) UpdateURL(ctx context.Context, id int64, originalUrl string, expiresAt *time.Time, maxClicks *int32) (domain.Url, error) {
	params := mapper.ToUpdateURLParams(id, originalUrl, expiresAt, maxClicks)

	row, err := r.q.UpdateURL(ctx, params)

	if err != nil {
		return domain.Url{}, err
	}

	return mapper.ToDomainUrl(row), nil
}

func (r *urlRepository) DeactivateURL(ctx context.Context, urlID int64) error {
	return r.q.DeactivateURL(ctx, urlID)
}


func (r *urlRepository) DeleteExpiredURLs(ctx context.Context) (int64, error) {
	return r.q.DeleteExpiredURLs(ctx)
}