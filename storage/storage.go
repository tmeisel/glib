package storage

import (
	"context"
	"io"
	"time"
)

type Service interface {
	Upload(ctx context.Context, bucketName, objectKey string, reader io.Reader, metadata map[string]string) error
	Download(ctx context.Context, bucketName, objectKey string) (io.ReadCloser, error)
	Delete(ctx context.Context, bucketName, objectKey string) error
	ListObjects(ctx context.Context, input ListObjectsInput) (*ListObjectsOutput, error)
}

type Bucket interface {
	Upload(ctx context.Context, objectKey string, reader io.Reader, metadata map[string]string) error
	Download(ctx context.Context, objectKey string) (io.ReadCloser, error)
	Delete(ctx context.Context, objectKey string) error
	ListObjects(ctx context.Context, prefix string, continuationToken *string) (*ListObjectsOutput, error)
}

type ListObjectsInput struct {
	BucketName        string
	Prefix            string
	ContinuationToken *string
}

type ListObjectsOutput struct {
	Objects []Object
	Next    *string
}

type Object struct {
	Name         string
	Size         int64
	LastModified time.Time
	Metadata     map[string]string
	IsDirectory  bool
}
