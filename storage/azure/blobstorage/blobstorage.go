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

func NewAzBlobStorage(accountName, accountKey string) (*AzBlobStorage, error) {
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

func (s *AzBlobStorage) Upload(ctx context.Context, containerName, blobName string, reader io.Reader, metadata map[string]string) error {
	containerURL := s.serviceURL.NewContainerURL(containerName)
	bbURL := containerURL.NewBlockBlobURL(blobName)

	_, err := azblob.UploadStreamToBlockBlob(ctx, reader, bbURL, azblob.UploadStreamToBlockBlobOptions{
		Metadata: metadata,
	})

	return err
}

func (s *AzBlobStorage) Download(ctx context.Context, containerName, blobName string) (io.ReadCloser, error) {
	containerURL := s.serviceURL.NewContainerURL(containerName)
	blobURL := containerURL.NewBlockBlobURL(blobName)

	resp, err := blobURL.Download(ctx, 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false, azblob.ClientProvidedKeyOptions{})

	if err != nil {
		return nil, err
	}

	return resp.Body(azblob.RetryReaderOptions{}), nil
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

	var objects []storage.Object
	for _, blobItem := range segmentIter.Segment.BlobItems {

		isDirectory := strings.HasSuffix(blobItem.Name, "/")
		objects = append(objects, storage.Object{
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
		Objects: objects,
		Next:    next,
	}, nil
}
