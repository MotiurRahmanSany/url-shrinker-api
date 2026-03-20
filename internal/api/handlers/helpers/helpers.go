package helpers

import (
	"net/http"
	"strings"

	"github.com/MotiurRahmanSany/url-shrinker-api/internal/api/middleware"
)

func GetUserIDFromContext(r *http.Request) (string, bool) {
	userID, ok := r.Context().Value(middleware.UserContextKey).(string)
	if !ok || strings.TrimSpace(userID) == "" {
		return "", false
	}
	return userID, true
}
