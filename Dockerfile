FROM golang:alpine as builder
MAINTAINER Jessica Frazelle <jess@linux.com>

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

RUN	apk add --no-cache \
	bash \
	ca-certificates

COPY . /go/src/github.com/genuinetools/ghb0t

RUN set -x \
	&& apk add --no-cache --virtual .build-deps \
		git \
		gcc \
		libc-dev \
		libgcc \
		make \
	&& cd /go/src/github.com/genuinetools/ghb0t \
	&& make static \
	&& mv ghb0t /usr/bin/ghb0t \
	&& apk del .build-deps \
	&& rm -rf /go \
	&& echo "Build complete."

FROM alpine:latest

COPY --from=builder /usr/bin/ghb0t /usr/bin/ghb0t
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs

ENTRYPOINT [ "ghb0t" ]
CMD [ "--help" ]
