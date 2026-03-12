package service

import (
	"context"
	"errors"
	"time"

	"github.com/MotiurRahmanSany/url-shrinker-api/internal/auth"
	"github.com/MotiurRahmanSany/url-shrinker-api/internal/domain"
	"github.com/MotiurRahmanSany/url-shrinker-api/internal/repository"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrEmailAlreadyInUse = errors.New("email already in use")
var ErrInvalidRefreshToken = errors.New("invalid or expired refresh token")

type LoginResponse struct {
	User         domain.User `json:"user"`
	AccessToken  string      `json:"access_token,omitempty"`
	RefreshToken string      `json:"refresh_token,omitempty"`
}

type AuthService interface {
	Register(ctx context.Context, email, password string) (domain.User, error)
	Login(ctx context.Context, email, password string) (LoginResponse, error)
	GetMe(ctx context.Context, userID string) (domain.User, error)
	Logout(ctx context.Context, refreshToken string) error
	RefreshToken(ctx context.Context, refreshToken string) (LoginResponse, error)
}

type authService struct {
	userRepo  repository.UserRepository
	tokenRepo repository.TokenRepository
	jwt       *auth.JWTManager
}

func NewAuthService(userRepo repository.UserRepository, tokenRepo repository.TokenRepository, jwt *auth.JWTManager) AuthService {
	return &authService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		jwt:       jwt,
	}
}

func (s *authService) Register(ctx context.Context, email, password string) (domain.User, error) {
	_, err := s.userRepo.GetUserByEmail(ctx, email)
	if err == nil {
		return domain.User{}, ErrEmailAlreadyInUse
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return domain.User{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, err
	}

	return s.userRepo.CreateUser(ctx, email, string(hash), "student")
}

func (s *authService) Login(ctx context.Context, email, password string) (LoginResponse, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return LoginResponse{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))

	if err != nil {
		return LoginResponse{}, err
	}

	accessToken, err := s.jwt.Generate(user.ID, user.Role)

	if err != nil {
		return LoginResponse{}, err
	}

	refreshToken, err := auth.GenerateRefreshToken()

	if err != nil {
		return LoginResponse{}, err
	}

	// Delete all previous refresh tokens for this user
	if err := s.tokenRepo.DeleteAllUserTokens(ctx, user.ID); err != nil {
		return LoginResponse{}, err
	}

	// Storing refresh token in db
	expiresAt := time.Now().Add(24 * time.Hour) // Set refresh token expiry time
	if err := s.tokenRepo.CreateToken(ctx, user.ID, refreshToken, expiresAt); err != nil {
		return LoginResponse{}, err
	}

	res := LoginResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return res, nil

}

func (s *authService) GetMe(ctx context.Context, userID string) (domain.User, error) {
	return s.userRepo.GetUserByID(ctx, userID)
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (LoginResponse, error) {
	token, err := s.tokenRepo.GetToken(ctx, refreshToken)
	if err != nil {
		return LoginResponse{}, err
	}

	if token.Revoked || time.Now().After(token.ExpiresAt) {
		return LoginResponse{}, ErrInvalidRefreshToken
	}

	user, err := s.userRepo.GetUserByID(ctx, token.UserID)
	if err != nil {
		return LoginResponse{}, err
	}

	newAccessToken, err := s.jwt.Generate(user.ID, user.Role)
	if err != nil {
		return LoginResponse{}, err
	}

	if err := s.tokenRepo.RevokeToken(ctx, refreshToken); err != nil {
		return LoginResponse{}, err
	}

	newRefreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return LoginResponse{}, err
	}

	expiresAt := time.Now().Add(24 * time.Hour)
	if err := s.tokenRepo.CreateToken(ctx, user.ID, newRefreshToken, expiresAt); err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{
		User:         user,
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil

}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	token, err := s.tokenRepo.GetToken(ctx, refreshToken)

	if err != nil {
		return ErrInvalidRefreshToken
	}

	if token.Revoked {
		return ErrInvalidRefreshToken
	}

	return s.tokenRepo.RevokeToken(ctx, refreshToken)
}
