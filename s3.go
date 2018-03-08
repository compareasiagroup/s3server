package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sirupsen/logrus"
)

type s3Provider struct {
	bucket  string
	prefix  string
	baseURL string
	client  *s3.S3
	ctx     context.Context
}

// List returns the files in an s3 bucket.
func (c *s3Provider) List(ctx context.Context, prefix string) (files []object, err error) {
	logrus.Info("about to list files")
	err = c.client.ListObjectsPagesWithContext(ctx, &s3.ListObjectsInput{
		Bucket: aws.String(c.bucket),
		Prefix: aws.String(prefix),
	}, func(p *s3.ListObjectsOutput, lastPage bool) bool {
		for _, o := range p.Contents {
			files = append(files, object{
				Name:    aws.StringValue(o.Key),
				Size:    aws.Int64Value(o.Size),
				BaseURL: c.baseURL,
			})
		}
		return true // continue paging
	})

	if err != nil {
		panic(fmt.Sprintf("failed to list objects for bucket, %s, %v", c.bucket, err))
	}

	fmt.Println("Objects in bucket:", files)

	return files, nil
}

// Prefix returns the prefix in an s3 bucket.
func (c *s3Provider) Prefix() string {
	return c.prefix
}

// BaseURL returns the baseURL in an s3 bucket.
func (c *s3Provider) BaseURL() string {
	return c.baseURL
}
