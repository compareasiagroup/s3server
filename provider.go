package main

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sirupsen/logrus"
)

type cloud interface {
	List(ctx context.Context, prefix string) ([]object, error)
	Prefix() string
	BaseURL() string
}

func newProvider(provider, bucket, s3Region, s3AccessKey, s3SecretKey string) (cloud, error) {
	conf := initAwsConfig(s3Region, s3AccessKey, s3SecretKey)

	p := s3Provider{}
	p.client = s3.New(session.New(conf))
	p.bucket, p.prefix = cleanBucketName(p.bucket)
	p.baseURL = p.bucket + ".s3.amazonaws.com"

	logrus.Info("baseURL: %q", p.baseURL)
	return &p, nil
}

// cleanBucketName returns the bucket and prefix
// for a given s3bucket.
func cleanBucketName(bucket string) (string, string) {
	bucket = strings.TrimPrefix(bucket, "s3://")
	parts := strings.SplitN(bucket, "/", 2)
	if len(parts) == 1 {
		return bucket, "/"
	}

	return parts[0], parts[1]
}

func initAwsConfig(region, accessKey, secretKey string) *aws.Config {
	awsConfig := aws.NewConfig()
	creds := credentials.NewChainCredentials([]credentials.Provider{
		&credentials.StaticProvider{
			Value: credentials.Value{
				AccessKeyID:     accessKey,
				SecretAccessKey: secretKey,
			},
		},
		&credentials.EnvProvider{},
		&credentials.SharedCredentialsProvider{},
	})
	awsConfig.WithCredentials(creds)
	awsConfig.WithRegion(region)
	return awsConfig
}
