package service

import (
	"context"
	"errors"
	"strings"

	"github.com/MotiurRahmanSany/url-shrinker-api/internal/domain"
	"github.com/MotiurRahmanSany/url-shrinker-api/internal/repository"
)

var (
	ErrInvalidURLID = errors.New("invalid URL ID")
)

type URLClickStats struct {
	TotalClicks   int64                   `json:"total_clicks"`
	ClicksToday   int64                   `json:"clicks_today"`
	DailyTimeline []domain.ClickDailyStat `json:"daily_timeline"`
}

type ClickService interface {
    RecordClick(ctx context.Context, urlID int64, ipAddr, userAgent, referer *string) (domain.Click, error)
    GetURLStats(ctx context.Context, urlID int64) (URLClickStats, error)
}

type clickService struct {
	repo repository.ClickRepository
}

func NewClickService(repo repository.ClickRepository) ClickService {
	return &clickService{repo: repo}
}

func (s *clickService) RecordClick(
	ctx context.Context,
	urlID int64,
	ipAddr, userAgent, referer *string,
) (domain.Click, error) {

	if urlID <= 0 {
		return domain.Click{}, ErrInvalidURLID
	}

	clientIP := normalizeOptionalString(ipAddr)
	clientUserAgent := normalizeOptionalString(userAgent)
	clientReferrer := normalizeOptionalString(referer)

	return s.repo.CreateClick(
		ctx,
		urlID,
		clientIP,
		clientUserAgent,
		clientReferrer,
	)

}

func (s *clickService) GetURLStats(ctx context.Context, urlID int64) (URLClickStats, error) {
	if urlID <= 0 {
		return URLClickStats{}, ErrInvalidURLID
	}

	totalClicks, err := s.repo.CountClicksByURLID(ctx, urlID)
	if err != nil {
		return URLClickStats{}, err
	}

	todayClicks, err := s.repo.CountClicksTodayByURLID(ctx, urlID)
	if err != nil {
		return URLClickStats{}, err
	}

	dailyTimeline, err := s.repo.GetClicksByURLIDGroupedByDay(ctx, urlID)

	if err != nil {
		return URLClickStats{}, err
	}

	if dailyTimeline == nil {
		dailyTimeline = make([]domain.ClickDailyStat, 0)
	}

	return URLClickStats{
		TotalClicks:   totalClicks,
		ClicksToday:   todayClicks,
		DailyTimeline: dailyTimeline,
	}, nil
}

func normalizeOptionalString(s *string) *string {
	if s == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*s)

	if trimmed == "" {
		return nil
	}

	return &trimmed
}
