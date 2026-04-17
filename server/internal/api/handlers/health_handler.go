package handlers

import (
	"net/http"

	"github.com/MotiurRahmanSany/url-shrinker-api/internal/response"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	_ = response.Success(w, http.StatusOK, "Server is running and healthy!", nil)
}
