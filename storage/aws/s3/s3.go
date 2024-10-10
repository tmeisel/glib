package s3storage

import (
	"context"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/tmeisel/glib/storage"
)

type S3Storage struct {
	uploader *manager.Uploader
	client   *s3.Client
}

var _ storage.Service = &S3Storage{}

func NewS3Storage(ctx context.Context, region string) (*S3Storage, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)

	return &S3Storage{
		uploader: manager.NewUploader(client),
		client:   client,
	}, nil
}

func (s *S3Storage) Upload(ctx context.Context, bucketName, objectKey string, reader io.Reader, metadata map[string]string) error {
	_, err := s.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:   &bucketName,
		Key:      &objectKey,
		Body:     reader,
		Metadata: metadata,
	})
	return err
}

func (s *S3Storage) Download(ctx context.Context, bucketName, objectKey string) (io.ReadCloser, error) {
	objectOutput, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &objectKey,
	})

	if err != nil {
		return nil, err
	}

	return io.NopCloser(objectOutput.Body), nil
}

func (s *S3Storage) Delete(ctx context.Context, bucketName, objectKey string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &bucketName,
		Key:    &objectKey,
	})

	return err
}

func (s *S3Storage) ListObjects(ctx context.Context, input storage.ListObjectsInput) (*storage.ListObjectsOutput, error) {
	resp, err := s.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket:            &input.BucketName,
		Prefix:            &input.Prefix,
		ContinuationToken: input.ContinuationToken,
	})

	if err != nil {
		return nil, err
	}

	var objects []storage.Object
	for _, item := range resp.Contents {
		objects = append(objects, storage.Object{
			Name:         *item.Key,
			Size:         *item.Size,
			LastModified: *item.LastModified,
			Metadata:     nil,
			IsDirectory:  strings.HasSuffix(*item.Key, "/"),
		})
	}

	var next *string
	if resp.NextContinuationToken != nil {
		next = resp.NextContinuationToken
	}

	return &storage.ListObjectsOutput{
		Objects: objects,
		Next:    next,
	}, nil
}
