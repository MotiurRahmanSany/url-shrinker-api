package main

import (
	"log/slog"
	"net/http"
	"os"
)

func main() {
	slog.Info("Setting up API server...")
	// setup logger before any other logging calls
	setupLogger()

	mux := http.NewServeMux()

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		slog.Error("Failed to start API server", "error", err)
	}
}

func setupLogger() {
	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo, // Only Info and above (no Debug noise in production)
	}

	if os.Getenv("APP_ENV") == "production" {
		// JSON format — for log aggregation tools (Datadog, Loki, CloudWatch, etc.)
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		// Human-readable text format — for local development
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	// SetDefault makes this handler the global default
	// After this, all slog.Info/Warn/Error calls anywhere in the program use it
	slog.SetDefault(slog.New(handler))
}
