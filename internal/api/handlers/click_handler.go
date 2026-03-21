package handlers

import (
	"net/http"
	"strings"

	"github.com/MotiurRahmanSany/url-shrinker-api/internal/api/handlers/helpers"
	"github.com/MotiurRahmanSany/url-shrinker-api/internal/response"
	"github.com/MotiurRahmanSany/url-shrinker-api/internal/service"
)

type ClickHandler struct {
	cs service.ClickService
	us service.UrlService
}

func NewClickHandler(cs service.ClickService, us service.UrlService) *ClickHandler {
	return &ClickHandler{cs: cs, us: us}
}

func (h *ClickHandler) GetURLStats(w http.ResponseWriter, r *http.Request) {
	userID, ok := helpers.GetUserIDFromContext(r)

	if !ok {
		_ = response.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	code := strings.TrimSpace(r.PathValue("code"))

	if code == "" {
		_ = response.Error(w, http.StatusBadRequest, "Short Code is required", nil)
		return
	}

	url, err := h.us.GetURLByShortCode(r.Context(), code)
	if err != nil {
		writeURLError(w, err)
		return
	}

	if !isOwner(url, userID) {
		_ = response.Error(w, http.StatusForbidden, "You do not have permission to view this URL's stats", nil)
		return 
	}

	stats, err := h.cs.GetURLStats(r.Context(), url.ID)
	if err != nil {
		_ = response.Error(w, http.StatusInternalServerError, "Failed to retrieve URL stats", nil)
		return
	}

	_ = response.Success(w, http.StatusOK, "URL stats retrieved successfully", stats)
}
