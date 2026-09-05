package router

import (
	"net/http"

	"github.com/debangshu919/transcodex/internal/http/handlers"
)

func HttpHandler() *http.ServeMux {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /healthz", handlers.Health)

	// Formats
	mux.HandleFunc("GET /formats", handlers.ListFormats)
	mux.HandleFunc("GET /formats/{format}", handlers.GetFormat)

	// API keys
	mux.HandleFunc("POST /auth/api-keys", handlers.NotImplemented)
	mux.HandleFunc("GET /auth/api-keys", handlers.NotImplemented)
	mux.HandleFunc("DELETE /auth/api-keys/{key}", handlers.NotImplemented)

	// File upload
	mux.HandleFunc("POST /files/upload", handlers.NotImplemented)
	mux.HandleFunc("GET /files", handlers.NotImplemented)
	mux.HandleFunc("GET /files/{id}", handlers.NotImplemented)
	mux.HandleFunc("DELETE /files/{id}", handlers.NotImplemented)

	// Conversions
	mux.HandleFunc("POST /conversions", handlers.NotImplemented)
	mux.HandleFunc("GET /conversions", handlers.NotImplemented)
	mux.HandleFunc("GET /conversions/{id}", handlers.NotImplemented)
	mux.HandleFunc("DELETE /conversions/{id}", handlers.NotImplemented)

	return mux
}
