package s3storage

import (
	"bytes"
	"context"
	"crypto/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tmeisel/glib/storage"
	"github.com/tmeisel/glib/storage/aws"
)

var (
	accessKeyID     = os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
)

const (
	bucketName = "glib-s3-testing"
)

func TestNewWithCredentialProvider(t *testing.T) {
	if accessKeyID == "" || secretAccessKey == "" {
		t.Skip("missing AWS_ACCESS_KEY_ID or AWS_SECRET_ACCESS_KEY")
	}

	cp := aws.NewStaticCredentialProvider(aws.StaticCredentials{
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
	})

	ctx := context.Background()

	s3, err := NewWithCredentialProvider(ctx, "eu-central-1", cp)
	require.NoError(t, err)
	require.NotNil(t, s3)

	output, err := s3.ListObjects(ctx, storage.ListObjectsInput{
		BucketName: bucketName,
		Prefix:     "",
	})

	require.NoError(t, err)
	require.NotNil(t, output)
}

func TestAzBlobStorage_Upload(t *testing.T) {
	if accessKeyID == "" || secretAccessKey == "" {
		t.Skip("missing AWS_ACCESS_KEY_ID or AWS_SECRET_ACCESS_KEY")
	}

	ctx := context.Background()
	cp := aws.NewStaticCredentialProvider(aws.StaticCredentials{
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
	})

	s3, err := NewWithCredentialProvider(ctx, "eu-central-1", cp)

	require.NoError(t, err)
	require.NotNil(t, s3)

	data := make([]byte, 256)
	if _, err := rand.Read(data); err != nil {
		t.Fatal(err)
	}

	type testCase struct {
		BlobName string
		Meta     map[string]string
	}

	for name, tc := range map[string]testCase{
		"file only": {
			BlobName: "01/01.dat",
		},
		"file with meta": {
			BlobName: "02/02.dat",
			Meta: map[string]string{
				"contenttype":     "application/octet-stream",
				"contentlanguage": "en-us",
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			err := s3.Upload(context.Background(), bucketName, tc.BlobName, bytes.NewReader(data), tc.Meta)
			require.NoError(t, err)

			obj, err := s3.Download(context.Background(), bucketName, tc.BlobName)
			require.NoError(t, err)

			require.NotNil(t, obj)
			assert.Equal(t, tc.BlobName, obj.Name)

			if tc.Meta == nil {
				tc.Meta = map[string]string{}
			}
			assert.Equal(t, tc.Meta, obj.Metadata)

			err = s3.Delete(context.Background(), bucketName, tc.BlobName)
			require.NoError(t, err)

			_, err = s3.Download(context.Background(), bucketName, tc.BlobName)
			require.Error(t, err)
		})
	}

}
