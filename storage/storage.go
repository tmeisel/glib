package storage

import (
	"context"
	"io"
	"time"
)

type Service interface {
	Upload(ctx context.Context, bucketName, objectKey string, reader io.Reader, metadata map[string]string) error
	Download(ctx context.Context, bucketName, objectKey string) (*Object, error)
	Delete(ctx context.Context, bucketName, objectKey string) error
	ListObjects(ctx context.Context, input ListObjectsInput) (*ListObjectsOutput, error)
}

type Bucket interface {
	Upload(ctx context.Context, objectKey string, reader io.Reader, metadata map[string]string) error
	Download(ctx context.Context, objectKey string) (*Object, error)
	Delete(ctx context.Context, objectKey string) error
	ListObjects(ctx context.Context, input ListBucketObjectsInput) (*ListObjectsOutput, error)
}

type ListObjectsInput struct {
	BucketName        string
	Prefix            string
	ContinuationToken *string
}

type ListObjectsOutput struct {
	BucketName string
	Prefix     string
	Objects    []ListObject
	Next       *string
}

func (o *ListObjectsOutput) More() bool {
	return o.Next != nil
}

func (o *ListObjectsOutput) GetNextInput() ListObjectsInput {
	return ListObjectsInput{
		BucketName:        o.BucketName,
		Prefix:            o.Prefix,
		ContinuationToken: o.Next,
	}
}

type ListBucketObjectsInput struct {
	Prefix            string
	ContinuationToken *string
}

type ListObject struct {
	Name         string
	Size         int64
	LastModified time.Time
	Metadata     map[string]string
	IsDirectory  bool
}

type Object struct {
	ListObject
	Reader io.Reader
}
