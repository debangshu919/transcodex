package storage

import (
	"context"
	"io"
)

type Storage interface {
	ListBuckets(ctx context.Context) ([]string, error)
	Upload(ctx context.Context, bucket, key string, body io.Reader) error
	Download(ctx context.Context, bucket, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, bucket, key string) error
}
