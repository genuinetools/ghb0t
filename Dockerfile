FROM alpine:latest
MAINTAINER Jessica Frazelle <jess@linux.com>

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

RUN	apk add --no-cache \
	ca-certificates

COPY . /go/src/github.com/jfrazelle/ghb0t

RUN apk add --no-cache --virtual .build-deps \
		go \
		git \
		gcc \
		libc-dev \
		libgcc \
	set -x \
	&& apk update \
	&& apk add $buildDeps \
	&& cd /go/src/github.com/jfrazelle/ghb0t \
	&& go build -o /usr/bin/ghb0t . \
	&& apk del .build-deps \
	&& rm -rf /go \
	&& echo "Build complete."


ENTRYPOINT [ "ghb0t" ]
