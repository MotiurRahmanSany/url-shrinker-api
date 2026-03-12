package domain

import "time"

type RefreshToken struct{
	ID int64 `json:"id"`
	UserID string `json:"user_id"`
	Token string `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	Revoked bool `json:"revoked"`
	CreatedAt time.Time `json:"created_at"`
}