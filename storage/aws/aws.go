package aws

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

type StaticCredentials struct {
	AccessKeyID     string `envconfig:"ACCESS_KEY_ID"`
	SecretAccessKey string `envconfig:"SECRET_ACCESS_KEY"`
}

func NewStaticCredentialProvider(c StaticCredentials) aws.CredentialsProvider {
	return credentials.NewStaticCredentialsProvider(c.AccessKeyID, c.SecretAccessKey, "")
}
