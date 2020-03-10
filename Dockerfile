FROM golang:alpine as builder

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

RUN	apk add --no-cache \
	ca-certificates

COPY . /go/src/github.com/jessfraz/s3server

RUN set -x \
	&& apk add --no-cache --virtual .build-deps \
		git \
		gcc \
		libc-dev \
		libgcc \
		make \
	&& cd /go/src/github.com/jessfraz/s3server \
	&& make static \
	&& mv s3server /usr/bin/s3server \
	&& apk del .build-deps \
	&& rm -rf /go \
	&& echo "Build complete."

FROM scratch

COPY --from=builder /usr/bin/s3server /usr/bin/s3server
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs

COPY static /static
COPY templates /templates

ENTRYPOINT [ "s3server" ]
CMD [ "--help" ]
