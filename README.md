s3server
========

Static server for s3 files.

## Usage

```console
$ s3server -h

 Server to index & view files from s3 bucket.
 Version: v0.2.3
 Build: e5a60e2

  -bucket string
        bucket path from which to serve files
  -cert string
        path to ssl certificate
  -interval string
        interval to generate new index.html's at (default "5m")
  -key string
        path to ssl key
  -p string
        port for server to run on (default "8080")
  -provider string
        cloud provider (ex. s3, gcs) (default "s3")
  -s3key string
        s3 access key
  -s3region string
        aws region for the bucket (default "us-west-2")
  -s3secret string
        s3 access secret
  -v    print version and exit (shorthand)
  -version
        print version and exit
```

**run with the docker image**

```console
# On AWS S3
$ docker run -d \
    --restart always \
    -e AWS_ACCESS_KEY_ID \
    -e AWS_SECRET_ACCESS_KEY \
    -p 8080:8080 \
    --name s3server \
    --tmpfs /tmp \
    ${ECR_REGISTRY_URL}/${ECR_REPO}:${VERSION} -bucket s3://hugthief/gifs -s3Region ap-southeast-1
```
