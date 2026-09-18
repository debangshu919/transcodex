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
	Upload(ctx context.Context, bucket, key string, body io.Reader) error
	Download(ctx context.Context, bucket, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, bucket, key string) error
	ListFiles(ctx context.Context, bucket string) ([]File, error)
	GenerateDownloadLink(ctx context.Context, bucket, key string) (string, error)
}
