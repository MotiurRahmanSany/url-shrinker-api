package mapper

import (
	"time"

	"github.com/MotiurRahmanSany/url-shrinker-api/internal/db"
	"github.com/MotiurRahmanSany/url-shrinker-api/internal/domain"
)

func ToDomainUrl(row db.Url) domain.Url {
	return domain.Url{
		ID:          row.ID,
		OriginalUrl: row.OriginalUrl,
		ShortCode:   row.ShortCode,
		UserID:      FromPgUUID(row.UserID),
		IsActive:    row.IsActive,
		ExpiresAt:   FromTimestamptz(row.ExpiresAt),
		MaxClicks:   FromPgInt4(row.MaxClicks),
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}
}

func ToCreateUrlParams(
	shortCode, originalUrl, userId string, expiresAt *time.Time, maxClicks *int32,
) (db.CreateURLParams, error) {
	pgUUid, err := ToPgUUID(userId)

	if err != nil {
		return db.CreateURLParams{}, err
	}
	return db.CreateURLParams{
		ShortCode:   shortCode,
		OriginalUrl: originalUrl,
		UserID:      pgUUid,
		ExpiresAt:   ToTimestamptz(expiresAt),
		MaxClicks:   ToPgInt4(maxClicks),
	}, nil
}


func ToGetURLsByUserIDParams(userId string, limit, offset *int32) (db.GetURLsByUserIDParams, error) {
	pgUUid, err := ToPgUUID(userId)

	if err != nil {
		return db.GetURLsByUserIDParams{}, err
	}
	return db.GetURLsByUserIDParams{
		UserID: pgUUid,
		Limit:  *limit,
		Offset: *offset,
	}, nil
}

func ToUpdateURLParams (id int64, originalUrl string, expiresAt *time.Time, maxClicks *int32) db.UpdateURLParams{
	return db.UpdateURLParams{
		ID:          id,
		OriginalUrl: originalUrl,
		ExpiresAt:   ToTimestamptz(expiresAt),
		MaxClicks:   ToPgInt4(maxClicks),
	}
}