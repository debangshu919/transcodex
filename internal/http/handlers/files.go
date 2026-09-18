package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"

	"github.com/debangshu919/transcodex/internal/storage"
	"github.com/google/uuid"
)

type FileHandler struct {
	storage *storage.S3Storage
	db      *sql.DB
	logger  *slog.Logger
}

func NewFileHandler(storage *storage.S3Storage, db *sql.DB, logger *slog.Logger) *FileHandler {
	return &FileHandler{
		storage: storage,
		db:      db,
		logger:  logger,
	}
}

func (h *FileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	// TODO: Resolve the bucket name from the request
	bucket := "meow"

	// Upload the file
	file, header, err := r.FormFile("file")
	if err != nil {
		h.logger.Error("File is required")
		http.Error(w, "file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Generate a unique object key.
	ext := filepath.Ext(header.Filename)
	key := fmt.Sprintf("uploads/%s%s", uuid.NewString(), ext)

	// Stream the file directly to S3.
	if err := h.storage.Upload(
		r.Context(),
		bucket,
		key,
		file,
	); err != nil {
		h.logger.Error("failed to upload file", "error", err)
		http.Error(
			w,
			"failed to upload file",
			http.StatusInternalServerError,
		)
		return
	}

	// Return response.
	response := map[string]any{
		"key":      key,
		"filename": header.Filename,
		"size":     header.Size,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(response)
}

func ListFiles() {
	// List all files

}

func DeleteFile() {
	// Delete the file
}

func GetFile() {
	// Generate a download link
	// Redirect to the download link
}
