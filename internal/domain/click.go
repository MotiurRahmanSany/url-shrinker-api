package domain

import "time"

type Click struct {
	ID        int64     `json:"id"`
	UrlID     int64     `json:"url_id"`
	ClickedAt time.Time `json:"clicked_at"`
	IpAddress *string   `json:"ip_address,omitempty"`
	UserAgent *string   `json:"user_agent,omitempty"`
	Referer   *string   `json:"referer,omitempty"`
}
