package s3

import (
	"context"
	"io"

	"github.com/tmeisel/glib/storage"
)

type Bucket struct {
	bucketName string
	storage    *S3
}

var _ storage.Bucket = (*Bucket)(nil)

func (b Bucket) Upload(ctx context.Context, objectKey string, reader io.Reader, metadata map[string]string) error {
	return b.storage.Upload(ctx, b.bucketName, objectKey, reader, metadata)
}

func (b Bucket) Download(ctx context.Context, objectKey string) (*storage.Object, error) {
	return b.storage.Download(ctx, b.bucketName, objectKey)
}

func (b Bucket) Delete(ctx context.Context, objectKey string) error {
	return b.storage.Delete(ctx, b.bucketName, objectKey)
}

func (b Bucket) ListObjects(ctx context.Context, input storage.ListBucketObjectsInput) (*storage.ListObjectsOutput, error) {
	return b.storage.ListObjects(ctx, storage.ListObjectsInput{
		BucketName:        b.bucketName,
		Prefix:            input.Prefix,
		ContinuationToken: input.ContinuationToken,
	})
}
