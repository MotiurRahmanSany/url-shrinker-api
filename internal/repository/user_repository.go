package repository

import (
	"context"

	"github.com/MotiurRahmanSany/url-shrinker-api/internal/db"
	"github.com/MotiurRahmanSany/url-shrinker-api/internal/domain"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	CreateUser(ctx context.Context, email, passwordHash, role string) (domain.User, error)
	GetUserByID(ctx context.Context, id string) (domain.User, error)
}

type userRepository struct {
	q *db.Queries
}

func NewUserRepository(q *db.Queries) UserRepository {
	return &userRepository{q: q}
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	row, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	user := domain.User{
		ID:           row.ID.String(),
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		Role:         row.Role,
		IsActive:     row.IsActive,
		CreatedAt:    row.CreatedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
	}
	return user, nil
}

func (r *userRepository) CreateUser(ctx context.Context, email, passwordHash, role string) (domain.User, error) {
	row, err := r.q.CreateUser(ctx, db.CreateUserParams{
		Email:        email,
		PasswordHash: passwordHash,
		Role:         role,
	})
	if err != nil {
		return domain.User{}, err

	}
	createdUser := domain.User{
		ID:        row.ID.String(),
		Email:     row.Email,
		Role:      row.Role,
		IsActive:  row.IsActive,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}
	return createdUser, nil

}

func (r *userRepository) GetUserByID(ctx context.Context, id string) (domain.User, error) {
	var pgID pgtype.UUID
	err := pgID.Scan(id)
	if err != nil {
		return domain.User{}, err
	}
	row, err := r.q.GetUserByID(ctx, pgID)
	if err != nil {
		return domain.User{}, err
	}

	user := domain.User{
		ID:        row.ID.String(),
		Email:     row.Email,
		Role:      row.Role,
		IsActive:  row.IsActive,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}
	return user, nil
}
