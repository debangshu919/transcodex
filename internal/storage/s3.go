package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/debangshu919/transcodex/internal/config"
)

type S3Storage struct {
	client *s3.Client
	bucket string
}

func NewS3Storage(cfg config.Config) *S3Storage {
	awsCfg := aws.Config{
		Region: cfg.AWSRegion,
		Credentials: credentials.NewStaticCredentialsProvider(
			cfg.AWSKey,
			cfg.AWSSecret,
			"",
		),
	}

	var client *s3.Client

	if cfg.Env == "development" {
		client = s3.NewFromConfig(awsCfg, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(
				"http://localhost.floci.io:4566",
			)
		})
	} else {
		client = s3.NewFromConfig(awsCfg)
	}

	return &S3Storage{
		client: client,
		bucket: cfg.S3Bucket,
	}
}

func (s *S3Storage) ListBuckets(
	ctx context.Context,
) ([]string, error) {

	result, err := s.client.ListBuckets(
		ctx,
		&s3.ListBucketsInput{},
	)

	if err != nil {
		return nil, fmt.Errorf("list buckets: %w", err)
	}

	buckets := make([]string, 0, len(result.Buckets))

	for _, bucket := range result.Buckets {
		buckets = append(buckets, *bucket.Name)
	}

	return buckets, nil
}

func (s *S3Storage) Upload(
	ctx context.Context,
	key string,
	body io.Reader,
) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   body,
	})

	if err != nil {
		return fmt.Errorf(
			"upload object %q to bucket %q: %w",
			key,
			s.bucket,
			err,
		)
	}

	return nil
}

func (s *S3Storage) Download(
	ctx context.Context,
	key string,
) (io.ReadCloser, error) {
	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, fmt.Errorf(
			"download object %q from bucket %q: %w",
			key,
			s.bucket,
			err,
		)
	}

	return result.Body, nil
}

func (s *S3Storage) Delete(
	ctx context.Context,
	key string,
) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return fmt.Errorf(
			"delete object %q from bucket %q: %w",
			key,
			s.bucket,
			err,
		)
	}

	return nil
}

func (s *S3Storage) ListFiles(
	ctx context.Context,
) ([]File, error) {
	var files []File

	paginator := s3.NewListObjectsV2Paginator(
		s.client,
		&s3.ListObjectsV2Input{
			Bucket: aws.String(s.bucket),
		},
	)

	for paginator.HasMorePages() {
		result, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf(
				"list files in bucket %q: %w",
				s.bucket,
				err,
			)
		}

		for _, obj := range result.Contents {
			file := File{
				Key: *obj.Key,
			}

			if obj.Size != nil {
				file.Size = *obj.Size
			}

			if obj.LastModified != nil {
				file.LastModified = obj.LastModified.String()
			}

			files = append(files, file)
		}
	}

	return files, nil
}

func (s *S3Storage) GenerateDownloadLink(
	ctx context.Context,
	key string,
) (string, error) {
	presigner := s3.NewPresignClient(s.client)

	req, err := presigner.PresignGetObject(
		ctx,
		&s3.GetObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(key),
		},
	)
	if err != nil {
		return "", fmt.Errorf(
			"generate download link for object %q from bucket %q: %w",
			key,
			s.bucket,
			err,
		)
	}

	return req.URL, nil
}
