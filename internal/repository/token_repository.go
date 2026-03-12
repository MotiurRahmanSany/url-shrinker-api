package repository

import (
	"context"
	"time"

	"github.com/MotiurRahmanSany/url-shrinker-api/internal/db"
	"github.com/MotiurRahmanSany/url-shrinker-api/internal/domain"
	"github.com/jackc/pgx/v5/pgtype"
)

type TokenRepository interface {
	CreateToken(ctx context.Context, userID, token string, expiresAt time.Time) error
	GetToken(ctx context.Context, token string) (domain.RefreshToken, error)
	RevokeToken(ctx context.Context, token string) error
	DeleteAllUserTokens(ctx context.Context, userID string) error
}

type tokenRepository struct {
	q *db.Queries
}

func NewTokenRepository(q *db.Queries) TokenRepository {
	return &tokenRepository{q: q}
}

func (r *tokenRepository) CreateToken(ctx context.Context, userID, token string, expiresAt time.Time) error {
	var pgID pgtype.UUID

	if err := pgID.Scan(userID); err != nil {
		return err
	}

	var pgExp pgtype.Timestamp

	if err := pgExp.Scan(expiresAt); err != nil {
		return err
	}

	_, err := r.q.CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
		UserID:    pgID,
		Token:     token,
		ExpiresAt: pgExp,
	})
	return err
}

func (r *tokenRepository) GetToken(ctx context.Context, token string) (domain.RefreshToken, error) {
	row, err := r.q.GetRefreshToken(ctx, token)
	if err != nil {
		return domain.RefreshToken{}, err
	}

	refreshToken := domain.RefreshToken{
		ID:        row.ID,
		UserID:    row.UserID.String(),
		Token:     row.Token,
		ExpiresAt: row.ExpiresAt.Time,
		Revoked:   row.Revoked,
		CreatedAt: row.CreatedAt.Time,
	}

	return refreshToken, nil
}

func (r *tokenRepository) RevokeToken(ctx context.Context, token string) error {
	return r.q.RevokeRefreshToken(ctx, token)
}

func (r *tokenRepository) DeleteAllUserTokens(ctx context.Context, userID string) error {
	var pgID pgtype.UUID

	if err := pgID.Scan(userID); err != nil {
		return err
	}

	return r.q.DeleteAllUserTokens(ctx, pgID)
}
