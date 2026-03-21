package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/MotiurRahmanSany/url-shrinker-api/internal/api/handlers/helpers"
	"github.com/MotiurRahmanSany/url-shrinker-api/internal/domain"
	"github.com/MotiurRahmanSany/url-shrinker-api/internal/response"
	"github.com/MotiurRahmanSany/url-shrinker-api/internal/service"
)

/*
demo for create
{
	"original_url": "https://www.example.com",
	"custom_short_code": "mycode123", // optional
	"expires_at": "2024-12-31T23:59:59Z", // optional
	"max_clicks": 100 // optional
}
*/

type CreateShortURLRequest struct {
	OriginalURL     string     `json:"original_url"`
	CustomShortCode string     `json:"custom_short_code,omitempty"`
	ExpiresAt       *time.Time `json:"expires_at,omitempty"`
	MaxClicks       *int32     `json:"max_clicks,omitempty"`
}

/*
demo for update
{
	"original_url": "https://www.newexample.com",
	"expires_at": "2024-12-31T23:59:59Z",
	"max_clicks": 200

}

*/

type UpdateURLRequest struct {
	OriginalURL *string    `json:"original_url,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	MaxClicks   *int32     `json:"max_clicks,omitempty"`
}

type UrlHandler struct {
	urlService   service.UrlService
	clickService service.ClickService
}

func NewUrlHandler(urlService service.UrlService, clickService service.ClickService) *UrlHandler {
	return &UrlHandler{urlService: urlService, clickService: clickService}
}

func (h *UrlHandler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	userID, ok := helpers.GetUserIDFromContext(r)

	if !ok {
		_ = response.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	var req CreateShortURLRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = response.Error(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	req.OriginalURL = strings.TrimSpace(req.OriginalURL)

	if req.OriginalURL == "" {
		_ = response.Error(w, http.StatusBadRequest, "original_url is required", nil)
		return
	}

	var expiresAt time.Time
	if req.ExpiresAt != nil {
		expiresAt = *req.ExpiresAt
	}

	var maxClicks int32
	if req.MaxClicks != nil {
		maxClicks = *req.MaxClicks
	}

	url, err := h.urlService.CreateShortURL(
		r.Context(),
		strings.TrimSpace(req.CustomShortCode),
		req.OriginalURL,
		userID,
		expiresAt,
		maxClicks,
	)

	if err != nil {
		writeURLError(w, err)
		return
	}

	_ = response.Success(w, http.StatusCreated, "Short URL created successfully", url)
}

func (h *UrlHandler) RedirectURL(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimSpace(r.PathValue("code"))

	if code == "" {
		_ = response.Error(w, http.StatusBadRequest, "Short code is required", nil)
		return
	}

	url, err := h.urlService.GetURLByShortCode(r.Context(), code)
	if err != nil {
		writeURLError(w, err)
		return
	}

	ip := r.Header.Get("X-Real-IP")

	if strings.TrimSpace(ip) == "" {
		ip = r.RemoteAddr
	}
	ua := r.UserAgent()
	ref := r.Referer()

	var ipPtr, uaPtr, refPtr *string
	if strings.TrimSpace(ip) != "" {
		ipPtr = &ip
	}
	if strings.TrimSpace(ua) != "" {
		uaPtr = &ua
	}
	if strings.TrimSpace(ref) != "" {
		refPtr = &ref
	}

	if _, err := h.clickService.RecordClick(r.Context(), url.ID, ipPtr, uaPtr, refPtr); err != nil {
		// Log the error but do not prevent the redirect, as click recording failure should not affect user experience
		// In a real application, you would want to log this error to your logging system for later analysis
		slog.WarnContext(r.Context(), "click recording failed; redirecting anyway",
			"short_code", code,
			"url_id", url.ID,
			"error", err,
		)
	}

	// Intentionally using 302 Found to allow clients to update the URL in their cache if needed
	// as 301 Moved Permanently can lead to aggressive caching by browsers and CDNs, which might cause issues if the original URL changes in the future.
	http.Redirect(w, r, url.OriginalUrl, http.StatusFound) // status is 302 here
}

func (h *UrlHandler) ListMyURLs(w http.ResponseWriter, r *http.Request) {
	userID, ok := helpers.GetUserIDFromContext(r)

	if !ok {
		_ = response.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	page, limit, offset, err := parsePagination(r)
	if err != nil {
		_ = response.Error(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	urls, err := h.urlService.GetURLsByUserID(r.Context(), userID, limit, offset)
	if err != nil {
		writeURLError(w, err)
		return
	}

	total := int64(len(urls)) // In a real implementation, you would want to get the total count from the database for pagination metadata

	data := response.NewPaginatedData(urls, page, limit, total)

	_ = response.Success(w, http.StatusOK, "URLs retrieved successfully", data)
}

func (h *UrlHandler) GetURLDetails(w http.ResponseWriter, r *http.Request) {
	userID, ok := helpers.GetUserIDFromContext(r)

	if !ok {
		_ = response.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	code := strings.TrimSpace(r.PathValue("code"))

	if code == "" {
		_ = response.Error(w, http.StatusBadRequest, "Short code is required", nil)
		return
	}

	url, err := h.urlService.GetURLByShortCode(r.Context(), code)
	if err != nil {
		writeURLError(w, err)
		return
	}

	if !isOwner(url, userID) {
		_ = response.Error(w, http.StatusForbidden, "You do not have permission to view this URL", nil)
		return
	}

	_ = response.Success(w, http.StatusOK, "URL details retrieved successfully", url)
}

func (h *UrlHandler) DeactivateURL(w http.ResponseWriter, r *http.Request) {
	userID, ok := helpers.GetUserIDFromContext(r)

	if !ok {
		_ = response.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	code := strings.TrimSpace(r.PathValue("code"))

	if code == "" {
		_ = response.Error(w, http.StatusBadRequest, "Short code is required", nil)
		return
	}

	url, err := h.urlService.GetURLByShortCode(r.Context(), code)
	if err != nil {
		writeURLError(w, err)
		return
	}

	if !isOwner(url, userID) {
		_ = response.Error(w, http.StatusForbidden, "You do not have permission to view this URL", nil)
		return
	}

	if err := h.urlService.DeactivateURL(r.Context(), url.ID, url.ShortCode); err != nil {
		writeURLError(w, err)
		return
	}

	_ = response.Success(w, http.StatusOK, "URL deactivated successfully", nil)
}

func (h *UrlHandler) UpdateURL(w http.ResponseWriter, r *http.Request) {
	userID, ok := helpers.GetUserIDFromContext(r)

	if !ok {
		_ = response.Error(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	code := strings.TrimSpace(r.PathValue("code"))

	if code == "" {
		_ = response.Error(w, http.StatusBadRequest, "Short code is required", nil)
		return
	}

	var req UpdateURLRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = response.Error(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	if req.OriginalURL == nil && req.ExpiresAt == nil && req.MaxClicks == nil {
		_ = response.Error(w, http.StatusBadRequest, "At least one field (original_url, expires_at, max_clicks) must be provided for update", nil)
		return
	}

	currUrl, err := h.urlService.GetURLByShortCode(r.Context(), code)
	if err != nil {
		writeURLError(w, err)
		return
	}

	if !isOwner(currUrl, userID) {
		_ = response.Error(w, http.StatusForbidden, "You do not have permission to update this URL", nil)
		return
	}

	newOriginalUrl := currUrl.OriginalUrl
	if req.OriginalURL != nil {
		newOriginalUrl = strings.TrimSpace(*req.OriginalURL)
		if newOriginalUrl == "" {
			_ = response.Error(w, http.StatusBadRequest, "original_url cannot be empty", nil)
			return
		}
	}

	var newExpiresAt time.Time
	if currUrl.ExpiresAt != nil {
		newExpiresAt = *currUrl.ExpiresAt
	}

	if req.ExpiresAt != nil {
		newExpiresAt = *req.ExpiresAt
	}

	var newMaxClicks int32
	if currUrl.MaxClicks != nil {
		newMaxClicks = *currUrl.MaxClicks
	}

	if req.MaxClicks != nil {
		newMaxClicks = *req.MaxClicks
	}

	updatedUrl, err := h.urlService.UpdateURL(
		r.Context(),
		currUrl.ID,
		newOriginalUrl,
		newExpiresAt,
		newMaxClicks,
	)

	if err != nil {
		writeURLError(w, err)
		return
	}

	_ = response.Success(w, http.StatusOK, "URL updated successfully", updatedUrl)
}

func isOwner(u domain.Url, requesterUserID string) bool {
	return strings.TrimSpace(u.UserID) != "" && u.UserID == requesterUserID
}

func parsePagination(r *http.Request) (page, limit, offset int32, err error) {
	page = 1
	limit = 10

	q := r.URL.Query()

	if p := strings.TrimSpace(q.Get("page")); p != "" {
		pi, convErr := strconv.Atoi(p)
		if convErr != nil || pi <= 0 {
			return 0, 0, 0, errors.New("page must be a positive integer")
		}
		page = int32(pi)
	}

	if l := strings.TrimSpace(q.Get("limit")); l != "" {
		li, convErr := strconv.Atoi(l)
		if convErr != nil || li <= 0 {
			return 0, 0, 0, errors.New("limit must be a positive integer")
		}
		if li > 100 {
			li = 100
		}
		limit = int32(li)
	}
	offset = (page - 1) * limit

	return page, limit, offset, nil
}

func writeURLError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidURL), errors.Is(err, service.ErrInvalidShortCode):
		_ = response.Error(w, http.StatusBadRequest, err.Error(), nil)
	case errors.Is(err, service.ErrURLNotFound):
		_ = response.Error(w, http.StatusNotFound, "URL not found", nil)
	case errors.Is(err, service.ErrShortCodeTaken):
		_ = response.Error(w, http.StatusConflict, "Short code already in use", nil)
	case errors.Is(err, service.ErrURLExpiredOrInactive):
		_ = response.Error(w, http.StatusGone, "URL is no longer active", nil)
	default:
		_ = response.Error(w, http.StatusInternalServerError, "Internal server error", err.Error())
	}
}
