package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"

	"github.com/MotiurRahmanSany/url-shrinker-api/internal/api/middleware"
	"github.com/MotiurRahmanSany/url-shrinker-api/internal/response"
	"github.com/MotiurRahmanSany/url-shrinker-api/internal/service"
)

type RegisterUserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token"`
}

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(s service.AuthService) *AuthHandler {
	return &AuthHandler{
		service: s,
	}
}

func isValidEmail(email string) bool {
	// Simple net/mail validation
	_, err := mail.ParseAddress(email)
	return err == nil
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterUserReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = response.Error(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	if req.Email == "" || req.Password == "" {
		_ = response.Error(w, http.StatusBadRequest, "Email and password are required", nil)
		return
	}

	if !isValidEmail(req.Email) {
		_ = response.Error(w, http.StatusBadRequest, "Invalid email format", nil)
		return
	}

	if len(req.Password) < 6 {
		_ = response.Error(w, http.StatusBadRequest, "Password must be at least 6 characters long", nil)
		return
	}

	user, err := h.service.Register(r.Context(), req.Email, req.Password)

	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyInUse) {
			_ = response.Error(w, http.StatusConflict, "Email already in use", nil)
			return
		}
		_ = response.Error(w, http.StatusInternalServerError, "Failed to register user", nil)
		return
	}

	_ = response.Success(w, http.StatusCreated, "User registered successfully", user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req RegisterUserReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = response.Error(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	loginResponse, err := h.service.Login(r.Context(), req.Email, req.Password)

	if err != nil {
		_ = response.Error(w, http.StatusUnauthorized, "Invalid credentials", nil)
		return
	}

	_ = response.Success(w, http.StatusOK, "User logged in successfully", loginResponse)
}

func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserContextKey).(string)
	if !ok || userID == "" {
		_ = response.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	user, err := h.service.GetMe(r.Context(), userID)
	if err != nil {
		_ = response.Error(w, http.StatusInternalServerError, "Failed to retrieve user", nil)
		return
	}

	_ = response.Success(w, http.StatusOK, "User retrieved successfully", user)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req RefreshTokenReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = response.Error(w, http.StatusBadRequest, "Invalid request body", nil)
		return

	}

	if req.RefreshToken == "" {
		_ = response.Error(w, http.StatusBadRequest, "Refresh token is required", nil)
		return
	}

	loginResponse, err := h.service.RefreshToken(r.Context(), req.RefreshToken)

	if err != nil {
		if errors.Is(err, service.ErrInvalidRefreshToken) {
			_ = response.Error(w, http.StatusUnauthorized, "Invalid or expired refresh token", nil)
			return
		}
		_ = response.Error(w, http.StatusInternalServerError, "Failed to refresh token", nil)
		return
	}

	_ = response.Success(w, http.StatusOK, "Token refreshed successfully", loginResponse)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req RefreshTokenReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = response.Error(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	if req.RefreshToken == "" {
		_ = response.Error(w, http.StatusBadRequest, "Refresh token is required", nil)
		return
	}

	if err := h.service.Logout(r.Context(), req.RefreshToken); err != nil {
		if errors.Is(err, service.ErrInvalidRefreshToken) {
			_ = response.Error(w, http.StatusUnauthorized, "Invalid or expired refresh token", nil)
			return
		}

		_ = response.Error(w, http.StatusInternalServerError, "Failed to logout", nil)
		return
	}

	_ = response.Success(w, http.StatusOK, "User logged out successfully", nil)
}
