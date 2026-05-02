package service

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/MotiurRahmanSany/url-shrinker-api/internal/cache"
	"github.com/MotiurRahmanSany/url-shrinker-api/internal/domain"
	"github.com/MotiurRahmanSany/url-shrinker-api/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrInvalidURL           = errors.New("invalid URL")
	ErrInvalidShortCode     = errors.New("invalid short code")
	ErrUserIDEmpty          = errors.New("user ID cannot be empty")
	ErrShortCodeTaken       = errors.New("short code already taken")
	ErrURLNotFound          = errors.New("URL not found")
	ErrURLExpiredOrInactive = errors.New("URL has expired or is inactive")
	ErrMaxClicksReached     = errors.New("maximum clicks reached for this URL")
	ErrForbidden            = errors.New("you do not have permission to perform this action")
)

var (
	defaultPageLimit      int32 = 10
	maxPageLimit          int32 = 100
	generatedCodeLength   int   = 7
	generatedCodeMaxRetry int   = 5
	defaultCacheTTL             = 24 * time.Hour
)

var (
	base62Alphabet    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	shortCodePattern  = regexp.MustCompile(`^[a-zA-Z0-9]{3,20}$`)
	reservedShortCode = map[string]struct{}{
		"health": {},
		"auth":   {},
		"urls":   {},
		"admin":  {},
		"api":    {},
	}
)

type UrlService interface {
	CreateShortURL(ctx context.Context, shortCode, orgCode string, userID string, expiresAt time.Time, maxClicks int32) (domain.Url, error)
	DeactivateURL(ctx context.Context, id int64, shortCode string) error
	GetURLByShortCode(ctx context.Context, shortCode string) (domain.Url, error)
	GetURLsByUserID(ctx context.Context, userID string, limit, offset int32) ([]domain.Url, error)
	UpdateURL(ctx context.Context, id int64, orgURL string, expiresAt time.Time, maxClicks int32) (domain.Url, error)
}

type urlService struct {
	repo      repository.UrlRepository
	clickRepo repository.ClickRepository
	cache     cache.Cache
}

func NewUrlService(repo repository.UrlRepository, clickRepo repository.ClickRepository, redisCache cache.Cache) UrlService {
	return &urlService{repo: repo, clickRepo: clickRepo, cache: redisCache}
}

func (s *urlService) CreateShortURL(
	ctx context.Context,
	shortCode,
	orgUrl,
	userID string,
	expiresAt time.Time,
	maxClicks int32,
) (domain.Url, error) {
	normalizedOriginal, err := validateOriginalURL(orgUrl)
	if err != nil {
		return domain.Url{}, err
	}

	expiresAtPtr := toOptionalTime(expiresAt)
	maxClicksPtr := toOptionalInt32(maxClicks)

	customCode := strings.TrimSpace(shortCode)

	if customCode != "" {
		if err := validateCustomShortCode(customCode); err != nil {
			return domain.Url{}, err
		}

		created, err := s.repo.CreateURL(ctx, customCode, normalizedOriginal, userID, expiresAtPtr, maxClicksPtr)
		if err != nil {
			if isUniqueViolation(err) {
				return domain.Url{}, ErrShortCodeTaken
			}
			return domain.Url{}, err

		}
		return created, nil
	}

	for i := 0; i < generatedCodeMaxRetry; i++ {
		code, err := generateBase62(generatedCodeLength)
		if err != nil {
			return domain.Url{}, err
		}

		created, err := s.repo.CreateURL(ctx, code, normalizedOriginal, userID, expiresAtPtr, maxClicksPtr)
		if err != nil {
			if isUniqueViolation(err) {
				continue
			}
			return domain.Url{}, err
		}
		return created, nil
	}

	return domain.Url{}, ErrShortCodeTaken

}

func (s *urlService) GetURLByShortCode(ctx context.Context, shortCode string) (domain.Url, error) {
	trimmed := strings.TrimSpace(shortCode)
	if trimmed == "" {
		return domain.Url{}, ErrInvalidShortCode
	}

	// Check cache first
	cacheKey := buildURLCacheKey(trimmed)

	if s.cache != nil {
		cached, err := s.cache.Get(ctx, cacheKey)
		if err == nil {
			var u domain.Url

			if jsonErr := json.Unmarshal([]byte(cached), &u); jsonErr == nil {
				if vErr := validateURLIsActive(u); vErr != nil {
					return domain.Url{}, vErr
				}

				if mErr := s.enforceMaxClickLimit(ctx, u); mErr != nil {
					return domain.Url{}, mErr
				}

				return u, nil
			}
		} else if !errors.Is(err, cache.ErrCacheMiss) {
			// Log cache error but continue to fetch from DB
			// fmt.Printf("cache error: %v\n", err)
			slog.Error("cache error", "err", err, "key", shortCode)
			// Optionally, we could return an error here if cache is critical

		}

	}

	u, err := s.repo.GetURLByShortCode(ctx, trimmed)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Url{}, ErrURLNotFound
		}
		return domain.Url{}, err
	}

	if err := validateURLIsActive(u); err != nil {
		return domain.Url{}, err
	}

	if err := s.enforceMaxClickLimit(ctx, u); err != nil {
		return domain.Url{}, err
	}

	if s.cache != nil {
		if payload, mErr := json.Marshal(u); mErr == nil {
			_ = s.cache.Set(ctx, cacheKey, string(payload), cacheTTLForURL(u.ExpiresAt))
		}
	}

	return u, nil
}

func (s *urlService) GetURLsByUserID(ctx context.Context, userID string, limit, offset int32) ([]domain.Url, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, ErrUserIDEmpty
	}

	if limit <= 0 {
		limit = defaultPageLimit
	}
	if limit > maxPageLimit {
		limit = maxPageLimit
	}
	if offset < 0 {
		offset = 0
	}

	return s.repo.GetURLsByUserID(ctx, userID, &limit, &offset)
}

func (s *urlService) UpdateURL(ctx context.Context,
	id int64,
	orgURL string,
	expiresAt time.Time,
	maxClicks int32,
) (domain.Url, error) {
	normalizedOriginal, err := validateOriginalURL(orgURL)

	if err != nil {
		return domain.Url{}, err
	}

	expiresAtPtr := toOptionalTime(expiresAt)
	maxClicksPtr := toOptionalInt32(maxClicks)

	updatedURL, err := s.repo.UpdateURL(ctx, id, normalizedOriginal, expiresAtPtr, maxClicksPtr)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Url{}, ErrURLNotFound
		}
		return domain.Url{}, err
	}

	if s.cache != nil {
		cacheKey := buildURLCacheKey(updatedURL.ShortCode)
		_ = s.cache.Delete(ctx, cacheKey) // Invalidate cache after update
	}

	return updatedURL, nil
}

func (s *urlService) DeactivateURL(ctx context.Context, id int64, shortCode string) error {
	if err := s.repo.DeactivateURL(ctx, id); err != nil {
		return err
	}

	// Invalidate cache

	if s.cache != nil && strings.TrimSpace(shortCode) != "" {
		cacheKey := buildURLCacheKey(shortCode)
		_ = s.cache.Delete(ctx, cacheKey)
	}

	return nil
}
func validateOriginalURL(rawUrl string) (string, error) {
	rawUrl = strings.TrimSpace(rawUrl)
	if rawUrl == "" {
		return "", ErrInvalidURL
	}
	u, err := url.ParseRequestURI(rawUrl)
	if err != nil {
		return "", ErrInvalidURL
	}
	scheme := strings.ToLower(u.Scheme)
	if scheme != "http" && scheme != "https" {
		return "", ErrInvalidURL
	}
	if u.Host == "" {
		return "", ErrInvalidURL
	}

	return u.String(), nil
}

// / this
func validateURLIsActive(u domain.Url) error {
	if !u.IsActive {
		return ErrURLExpiredOrInactive
	}

	if u.ExpiresAt != nil && time.Now().After(*u.ExpiresAt) {
		return ErrURLExpiredOrInactive
	}

	return nil
}

func validateCustomShortCode(code string) error {
	code = strings.TrimSpace(code)
	if !shortCodePattern.MatchString(code) {
		return ErrInvalidShortCode
	}

	if _, blocked := reservedShortCode[strings.ToLower(code)]; blocked {
		return ErrInvalidShortCode
	}

	return nil
}

func generateBase62(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("length must be greater than 0")
	}

	var sb strings.Builder
	sb.Grow(length)

	max := big.NewInt(int64(len(base62Alphabet)))

	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		sb.WriteByte(base62Alphabet[n.Int64()])
	}
	return sb.String(), nil
}

func buildURLCacheKey(shortCode string) string {
	return "url:" + shortCode
}

func cacheTTLForURL(expiresAt *time.Time) time.Duration {
	if expiresAt == nil {
		return defaultCacheTTL
	}
	ttl := time.Until(*expiresAt)

	if ttl <= 0 {
		return time.Minute
	}
	if ttl > defaultCacheTTL {
		return defaultCacheTTL
	}

	return ttl
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}

func toOptionalTime(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	tt := t.UTC() // store in UTC to avoid timezone issues

	return &tt
}

func toOptionalInt32(v int32) *int32 {
	if v <= 0 {
		return nil
	}
	return &v
}

func (s *urlService) enforceMaxClickLimit(ctx context.Context, u domain.Url) error {
	if u.MaxClicks == nil || *u.MaxClicks <= 0 {
		return nil
	}

	if s.clickRepo == nil {
		return nil
	}

	totalClicks, err := s.clickRepo.CountClicksByURLID(ctx, u.ID)
	if err != nil {
		return err
	}

	if totalClicks >= int64(*u.MaxClicks) {
		return ErrMaxClicksReached
	}

	return nil
}
