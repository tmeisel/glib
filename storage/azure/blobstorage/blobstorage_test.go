package blobstorage

import (
	"bytes"
	"context"
	"crypto/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	accessKey   = os.Getenv("AZURE_ACCESS_KEY")
	accountName = os.Getenv("AZURE_ACCOUNT_NAME")
)

const (
	containerName = "testing"
)

func TestNew(t *testing.T) {
	if accessKey == "" || accountName == "" {
		t.Skip("missing AZURE_ACCESS_KEY or AZURE_ACCOUNT_NAME")
	}

	bs, err := New(accountName, accessKey)
	require.NoError(t, err)
	assert.NotNil(t, bs)
}

func TestNewFromConf(t *testing.T) {
	if accessKey == "" || accountName == "" {
		t.Skip("missing AZURE_ACCESS_KEY or AZURE_ACCOUNT_NAME")
	}

	bs, err := NewFromConf(AccessKeyConf{
		AccessKey:   accessKey,
		AccountName: accountName,
	})

	require.NoError(t, err)
	assert.NotNil(t, bs)
}

func TestAzBlobStorage_Upload(t *testing.T) {
	if accessKey == "" || accountName == "" {
		t.Skip("missing AZURE_ACCESS_KEY or AZURE_ACCOUNT_NAME")
	}

	bs, err := New(accountName, accessKey)
	require.NoError(t, err)

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
			err := bs.Upload(context.Background(), containerName, tc.BlobName, bytes.NewReader(data), tc.Meta)
			require.NoError(t, err)

			obj, err := bs.Download(context.Background(), containerName, tc.BlobName)
			require.NoError(t, err)

			require.NotNil(t, obj)
			assert.Equal(t, tc.BlobName, obj.Name)

			if tc.Meta == nil {
				tc.Meta = map[string]string{}
			}
			assert.Equal(t, tc.Meta, obj.Metadata)

			err = bs.Delete(context.Background(), containerName, tc.BlobName)
			require.NoError(t, err)

			_, err = bs.Download(context.Background(), containerName, tc.BlobName)
			require.Error(t, err)
		})
	}

}
