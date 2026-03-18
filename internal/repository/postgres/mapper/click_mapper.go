package mapper

import (
	"github.com/MotiurRahmanSany/url-shrinker-api/internal/db"
	"github.com/MotiurRahmanSany/url-shrinker-api/internal/domain"
)

func ToDomainClick(row db.Click) domain.Click {
	return domain.Click{
		ID:        row.ID,
		UrlID:     row.UrlID,
		ClickedAt: row.ClickedAt.Time,
		IpAddress: FromPgText(row.IpAddress),
		UserAgent: FromPgText(row.UserAgent),
		Referer:   FromPgText(row.Referer),
	}
}

func ToCreateClickParams(
	urlID int64,
	ip *string,
	userAgent *string,
	referer *string,
) db.CreateClickParams {
	return db.CreateClickParams{
		UrlID:     urlID,
		IpAddress: ToPgText(ip),
		UserAgent: ToPgText(userAgent),
		Referer:   ToPgText(referer),
	}
}