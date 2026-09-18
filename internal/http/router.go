package router

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/debangshu919/transcodex/internal/http/handlers"
	"github.com/debangshu919/transcodex/internal/storage"
)

func HttpHandler(db *sql.DB, storage *storage.S3Storage, logger *slog.Logger) *http.ServeMux {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /healthz", handlers.Health(logger))

	// Formats
	mux.HandleFunc("GET /formats", handlers.ListFormats)
	mux.HandleFunc("GET /formats/{format}", handlers.GetFormat)

	// API keys
	mux.HandleFunc("POST /auth/api-keys", handlers.NotImplemented)
	mux.HandleFunc("GET /auth/api-keys", handlers.NotImplemented)
	mux.HandleFunc("DELETE /auth/api-keys/{key}", handlers.NotImplemented)

	// File upload
	fh := handlers.NewFileHandler(storage, db, logger)
	mux.HandleFunc("POST /files/upload", fh.UploadFile)
	mux.HandleFunc("GET /files", fh.ListFiles)
	mux.HandleFunc("GET /files/{id}", fh.GetFile)
	mux.HandleFunc("DELETE /files/{id}", fh.DeleteFile)

	// Conversions
	mux.HandleFunc("POST /conversions", handlers.NotImplemented)
	mux.HandleFunc("GET /conversions", handlers.NotImplemented)
	mux.HandleFunc("GET /conversions/{id}", handlers.NotImplemented)
	mux.HandleFunc("DELETE /conversions/{id}", handlers.NotImplemented)

	return mux
}
