package domain

import "time"

type Url struct {
	ID          int64      `json:"id"`
	ShortCode   string     `json:"short_code"`
	OriginalUrl string     `json:"original_url"`
	UserID      string     `json:"user_id"`
	IsActive    bool       `json:"is_active"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"` // Nullable if URLs don't expire
	MaxClicks   *int32     `json:"max_clicks,omitempty"` // Nullable if no click limit
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
