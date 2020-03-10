package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sirupsen/logrus"
)

type s3Provider struct {
	bucket    string
	prefix    string
	client    *s3.S3
	ctx       context.Context
	proxyPath string
}

// Prefix returns the prefix in an s3 bucket.
func (c *s3Provider) Prefix() string {
	return c.prefix
}

// ProxyPath returns the proxyPath in an s3 bucket.
func (c *s3Provider) ProxyPath() string {
	return c.proxyPath
}

// List returns the files in an s3 bucket.
func (c *s3Provider) List(ctx context.Context, prefix string) (files []object, err error) {
	err = c.client.ListObjectsPagesWithContext(ctx, &s3.ListObjectsInput{
		Bucket: aws.String(c.bucket),
		Prefix: aws.String(prefix),
	}, func(p *s3.ListObjectsOutput, lastPage bool) bool {
		for _, o := range p.Contents {
			files = append(files, object{
				Name:      aws.StringValue(o.Key),
				Size:      aws.Int64Value(o.Size),
				ProxyPath: c.proxyPath,
			})
		}
		return true // continue paging
	})

	if err != nil {
		panic(fmt.Sprintf("failed to list objects for bucket, %s, %v", c.bucket, err))
	}
	return files, nil
}

// ServeHTTP gets files for the c.proxyPath from the S3 bucket instead of through staticFileHandler
func (c *s3Provider) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	key := strings.TrimPrefix(req.URL.Path, c.proxyPath)
	logrus.Infof("Proxying: %q from %q", key, c.bucket)
	input := &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	}

	if v := req.Header.Get("If-None-Match"); v != "" {
		input.IfNoneMatch = aws.String(v)
	}

	var is304 bool
	resp, err := c.client.GetObjectWithContext(req.Context(), input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch awsErr.Code() {
		case s3.ErrCodeNoSuchKey:
			http.Error(rw, "Page Not Found", 404)
			return
		case "NotModified":
			is304 = true
			// continue so other headers get set appropriately
		default:
			log.Printf("Error: %v %v", awsErr.Code(), awsErr.Message())
			http.Error(rw, "Internal Error", 500)
			return
		}
	} else if err != nil {
		log.Printf("not aws error %v %s", err, err)
		http.Error(rw, "Internal Error", 500)
		return
	}

	var contentType string
	// if resp.ContentType != nil {
	// 	contentType = *resp.ContentType
	// 	log.Printf("resp.contentType: %v ", contentType)
	// }

	if contentType == "" {
		ext := path.Ext(key)
		contentType = mime.TypeByExtension(ext)
		log.Printf("mime.contentType: %v ", contentType)
	}

	if resp.ETag != nil && *resp.ETag != "" {
		rw.Header().Set("Etag", *resp.ETag)
	}

	if contentType != "" {
		rw.Header().Set("Content-Type", contentType)
	}
	if resp.ContentLength != nil && *resp.ContentLength > 0 {
		rw.Header().Set("Content-Length", fmt.Sprintf("%d", *resp.ContentLength))
	}

	if is304 {
		rw.WriteHeader(304)
	} else {
		io.Copy(rw, resp.Body)
		resp.Body.Close()
	}
}
