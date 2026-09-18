package storage

import (
	"context"
	"io"
)

type File struct {
	Key          string
	Size         int64
	LastModified string
}

type Storage interface {
	// ListBuckets(ctx context.Context) ([]string, error)
	Upload(ctx context.Context, key string, body io.Reader) error
	Download(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
	ListFiles(ctx context.Context) ([]File, error)
	GenerateDownloadLink(ctx context.Context, key string) (string, error)
}
