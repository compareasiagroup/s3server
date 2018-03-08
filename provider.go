package main

import (
	"context"
	"fmt"
	"net/url"
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
	url, err := parseS3URI(bucket)
	if err != nil {
		return nil, err
	}
	conf := initAwsConfig(s3Region, s3AccessKey, s3SecretKey)

	p := s3Provider{
		client: s3.New(session.New(conf)),
		bucket: url.Host,
		prefix: strings.TrimPrefix(url.Path, "/"),
	}
	logrus.Info(p.bucket, p.prefix)
	p.baseURL = p.bucket + ".s3.amazonaws.com"
	return &p, nil
}

func parseS3URI(s3URI string) (*url.URL, error) {
	s3URL, err := url.Parse(s3URI)
	if err != nil {
		return nil, fmt.Errorf("could not parse s3URI %q", s3URI)
	}
	if s3URL.Scheme != "s3" {
		return nil, fmt.Errorf("s3URI argument does not have valid protocol, should be 's3'")
	}
	if s3URL.Host == "" {
		return nil, fmt.Errorf("s3URI is missing bucket name")
	}
	return s3URL, nil
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
