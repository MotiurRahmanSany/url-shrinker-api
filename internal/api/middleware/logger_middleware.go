package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

// ANSI escape codes for terminal colours.
// These are standard codes understood by all modern terminals.
const (
	colourReset  = "\033[0m"
	colourGreen  = "\033[32m" // 2xx success
	colourYellow = "\033[33m" // 3xx redirect
	colourRed    = "\033[31m" // 4xx/5xx error
)

// responseWriter wraps http.ResponseWriter to capture the status code.
// The problem: once a handler calls w.WriteHeader(404), you cannot read
// that 404 back from the standard ResponseWriter — there's no getter.
// Solution: intercept WriteHeader, save the code, then pass it through.
type responseWriter struct {
	http.ResponseWriter     // embed the real writer (inherits all its methods)
	statusCode          int // we add this field to capture the code
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	// Default to 200 because handlers that write a body without calling
	// WriteHeader() explicitly are implicitly sending 200 OK.
	return &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
}

// WriteHeader overrides the embedded method.
// When a handler calls w.WriteHeader(code), Go calls THIS method because
// our struct has a WriteHeader method — it shadows the embedded one.
// We save the code, then call the real underlying WriteHeader.
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// statusColour returns the ANSI colour code appropriate for the HTTP status.
func statusColour(code int) string {
	switch {
	case code >= 200 && code < 300:
		return colourGreen
	case code >= 300 && code < 400:
		return colourYellow
	default: // 4xx and 5xx
		return colourRed
	}
}

// clientIP extracts the real client IP from the request.
// In production, apps run behind reverse proxies (Nginx, load balancers).
// The proxy sets special headers with the original client IP, because
// from the Go server's view, r.RemoteAddr is just the proxy's IP.
func clientIP(r *http.Request) string {
	// X-Real-IP: set by Nginx — single clean IP, most reliable
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	// X-Forwarded-For: standard header, may contain multiple IPs
	// e.g. "203.0.113.5, 10.0.0.1, 10.0.0.2" (client, proxy1, proxy2)
	// The first IP is always the original client.
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		return strings.Split(forwarded, ",")[0]
	}
	// Direct connection — no proxy — RemoteAddr is the real client
	return r.RemoteAddr
}

// Logger is the HTTP logging middleware. It wraps every request and logs:
// - method, path, status code, duration, client IP
// Applied to the whole mux in serve.go so it catches ALL requests.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Capture arrival time BEFORE the handler runs.
		start := time.Now()

		// Wrap the real ResponseWriter with our capturing wrapper.
		wrapped := newResponseWriter(w)

		// Run the actual handler chain. After this returns,
		// wrapped.statusCode holds whatever status the handler wrote.
		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)
		status := wrapped.statusCode

		if os.Getenv("APP_ENV") == "production" {
			// In production: emit structured JSON via slog — no ANSI codes.
			// Log aggregators (Datadog, Loki, CloudWatch) parse JSON fields,
			// so clean values with no escape sequences are essential here.
			slog.Info("request",
				"status_code", status,
				"status_text", http.StatusText(status),
				"method", r.Method,
				"path", r.URL.Path,
				"duration_ms", duration.Milliseconds(),
				"ip", clientIP(r),
			)
		} else {
			// In development: write a colorized line directly with fmt.Printf.
			//
			// Why not slog here? slog.TextHandler escapes any string value that
			// contains non-printable bytes (like the ANSI ESC \033) — it turns
			// them into \x1b, printing the literal text instead of the colour.
			// fmt.Printf writes the bytes as-is, so the terminal sees the real
			// escape sequences and renders the colours.
			colour := statusColour(status)
			fmt.Printf("%s | %s%d %s%s | %-7s | %-30s | %.6fms | %s\n",
				start.Format("2006/01/02 15:04:05"),
				colour, status, http.StatusText(status), colourReset,
				r.Method,
				r.URL.Path,
				float64(duration.Nanoseconds())/1_000_000.0,
				clientIP(r),
			)
		}
	})
}
