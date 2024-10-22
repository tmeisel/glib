package azblobstorage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/Azure/azure-storage-blob-go/azblob"

	"github.com/tmeisel/glib/storage"
)

type AzBlobStorage struct {
	serviceURL azblob.ServiceURL
}

type AccessKeyConf struct {
	AccessKey   string `envconfig:"ACCESS_KEY"`
	AccountName string `envconfig:"ACCOUNT_NAME"`
}

func New(accountName, accountKey string) (*AzBlobStorage, error) {
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return nil, err
	}

	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	u, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net", accountName))

	return &AzBlobStorage{
		serviceURL: azblob.NewServiceURL(*u, p),
	}, nil
}

func NewFromConf(conf AccessKeyConf) (*AzBlobStorage, error) {
	return New(conf.AccountName, conf.AccessKey)
}

func (s *AzBlobStorage) GetBucket(name string) Bucket {
	return Bucket{
		bucketName:  name,
		blobStorage: s,
	}
}

func (s *AzBlobStorage) Upload(ctx context.Context, containerName, blobName string, reader io.Reader, metadata map[string]string) error {
	containerURL := s.serviceURL.NewContainerURL(containerName)
	bbURL := containerURL.NewBlockBlobURL(blobName)

	_, err := azblob.UploadStreamToBlockBlob(ctx, reader, bbURL, azblob.UploadStreamToBlockBlobOptions{
		Metadata: metadata,
	})

	return err
}

func (s *AzBlobStorage) Download(ctx context.Context, containerName, blobName string) (*storage.Object, error) {
	containerURL := s.serviceURL.NewContainerURL(containerName)
	blobURL := containerURL.NewBlockBlobURL(blobName)

	resp, err := blobURL.Download(ctx, 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false, azblob.ClientProvidedKeyOptions{})

	if err != nil {
		return nil, err
	}

	return &storage.Object{
		ListObject: storage.ListObject{
			Name:         blobName,
			Size:         resp.ContentLength(),
			LastModified: resp.LastModified(),
			Metadata:     resp.NewMetadata(),
			IsDirectory:  false,
		},
		Reader: resp.Body(azblob.RetryReaderOptions{MaxRetryRequests: 20}),
	}, nil
}

func (s *AzBlobStorage) Delete(ctx context.Context, containerName, blobName string) error {
	containerURL := s.serviceURL.NewContainerURL(containerName)
	blobURL := containerURL.NewBlobURL(blobName)

	_, err := blobURL.Delete(ctx, azblob.DeleteSnapshotsOptionInclude, azblob.BlobAccessConditions{})
	return err
}

func (s *AzBlobStorage) ListObjects(ctx context.Context, input storage.ListObjectsInput) (*storage.ListObjectsOutput, error) {
	containerURL := s.serviceURL.NewContainerURL(input.BucketName)
	options := azblob.ListBlobsSegmentOptions{Prefix: input.Prefix}
	marker := azblob.Marker{Val: input.ContinuationToken}
	segmentIter, err := containerURL.ListBlobsHierarchySegment(ctx, marker, "", options)
	if err != nil {
		return nil, err
	}

	var objects []storage.ListObject
	for _, blobItem := range segmentIter.Segment.BlobItems {

		isDirectory := strings.HasSuffix(blobItem.Name, "/")
		objects = append(objects, storage.ListObject{
			Name:         blobItem.Name,
			Size:         *blobItem.Properties.ContentLength,
			LastModified: blobItem.Properties.LastModified,
			Metadata:     blobItem.Metadata,
			IsDirectory:  isDirectory,
		})
	}

	var next *string
	if segmentIter.NextMarker.Val != nil {
		next = segmentIter.NextMarker.Val
	}

	return &storage.ListObjectsOutput{
		BucketName: input.BucketName,
		Prefix:     input.Prefix,
		Objects:    objects,
		Next:       next,
	}, nil
}
