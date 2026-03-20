package middleware

import (
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/MotiurRahmanSany/url-shrinker-api/internal/cache"
	"github.com/MotiurRahmanSany/url-shrinker-api/internal/response"
)

func RateLimitingMiddleware(c cache.Cache, limit int64, window time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := getClientIP(r)
			key := "rate_create:" + clientIP

			count, err := c.Increment(r.Context(), key)

			if err != nil {
				// Fail Open: do not block if cache is unavailable
				next.ServeHTTP(w, r)
				return
			}

			if count == 1 {
				_ = c.Expire(r.Context(), key, window)
			}

			if count > limit {
				_ = response.Error(w, http.StatusTooManyRequests, "Rate limit exceeded. Please try again later.", nil)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func getClientIP(r *http.Request) string {
	if xr := strings.TrimSpace(r.Header.Get("X-Real-IP")); xr != "" {
		return xr
	}

	if xff := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}

	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))

	if err == nil && host != "" {
		return host
	}

	return "unknown"
}
