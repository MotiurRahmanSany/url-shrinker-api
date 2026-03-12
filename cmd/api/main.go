package main

import (
	"log/slog"
	"net/http"
)

func main() {
	slog.Info("Setting up API server...")

	mux := http.NewServeMux()

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		slog.Error("Failed to start API server", "error", err) 
	}
}
