package blobstorage

import (
	"context"
	"io"

	"github.com/tmeisel/glib/storage"
)

type Bucket struct {
	bucketName  string
	blobStorage *BlobStorage
}

var _ storage.Bucket = (*Bucket)(nil)

func (b Bucket) Upload(ctx context.Context, objectKey string, reader io.Reader, metadata map[string]string) error {
	return b.blobStorage.Upload(ctx, b.bucketName, objectKey, reader, metadata)
}

func (b Bucket) Download(ctx context.Context, objectKey string) (*storage.Object, error) {
	return b.blobStorage.Download(ctx, b.bucketName, objectKey)
}

func (b Bucket) Delete(ctx context.Context, objectKey string) error {
	return b.blobStorage.Delete(ctx, b.bucketName, objectKey)
}

func (b Bucket) ListObjects(ctx context.Context, input storage.ListBucketObjectsInput) (*storage.ListObjectsOutput, error) {
	return b.blobStorage.ListObjects(ctx, storage.ListObjectsInput{
		BucketName:        b.bucketName,
		Prefix:            input.Prefix,
		ContinuationToken: input.ContinuationToken,
	})
}
