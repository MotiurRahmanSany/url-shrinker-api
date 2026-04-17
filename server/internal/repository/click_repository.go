package repository

import (
	"context"

	"github.com/MotiurRahmanSany/url-shrinker-api/internal/db"
	"github.com/MotiurRahmanSany/url-shrinker-api/internal/domain"
	"github.com/MotiurRahmanSany/url-shrinker-api/internal/repository/postgres/mapper"
)

type ClickRepository interface {
	CountClicksByURLID(ctx context.Context, urlID int64) (int64, error)
	CountClicksTodayByURLID(ctx context.Context, urlID int64) (int64, error)
	CreateClick(ctx context.Context, urlId int64, ipAddr, userAgent, referrer *string) (domain.Click, error)
	GetClicksByURLIDGroupedByDay(ctx context.Context, urlID int64) ([]domain.ClickDailyStat, error)
}

type clickRepository struct {
	q *db.Queries
}

func NewClickRepository(q *db.Queries) ClickRepository {
	return &clickRepository{q: q}
}

func (r *clickRepository) CountClicksByURLID(ctx context.Context, urlID int64) (int64, error) {
	return r.q.CountClicksByURLID(ctx, urlID)
}

func (r *clickRepository) CountClicksTodayByURLID(ctx context.Context, urlID int64) (int64, error) {
	return r.q.CountClicksTodayByURLID(ctx, urlID)
}

func (r *clickRepository) GetClicksByURLIDGroupedByDay(ctx context.Context, urlID int64) ([]domain.ClickDailyStat, error) {
	rows, err := r.q.GetClicksByURLIDGroupedByDay(ctx, urlID)
	if err != nil {
		return nil, err
	}
	var stats []domain.ClickDailyStat
	for _, row := range rows {
		stats = append(stats, mapper.ToDomainClickDailyStat(row))
	}
	return stats, nil
}

func (r *clickRepository) CreateClick(ctx context.Context, urlId int64, ipAddr, userAgent, referrer *string) (domain.Click, error) {
	params := mapper.ToCreateClickParams(urlId, ipAddr, userAgent, referrer)
	row, err := r.q.CreateClick(ctx, params)
	if err != nil {
		return domain.Click{}, err
	}
	return mapper.ToDomainClick(row), nil
}
