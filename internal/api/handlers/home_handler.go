package handlers

import (
	"net/http"

	"github.com/MotiurRahmanSany/url-shrinker-api/internal/response"
)

type HomeHandler struct{}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

func (h *HomeHandler) Welcome(w http.ResponseWriter, r *http.Request) {
	_ = response.Success(w, http.StatusOK, "Welcome you!, This is the URL shrinker API :)", nil)
}
