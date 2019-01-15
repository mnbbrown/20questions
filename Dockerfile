FROM golang:alpine AS builder
RUN echo "http://dl-cdn.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories \
  && apk --no-cache upgrade \
  && apk add --no-cache curl ca-certificates go git musl-dev
RUN mkdir /src
COPY . /src
WORKDIR /src
RUN go build -o 20questions

FROM alpine:latest
RUN echo "http://dl-cdn.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories \
  && apk --no-cache upgrade \
  && apk add --no-cache ca-certificates \
  && rm -rf /var/cache/apk*
COPY --from=builder /src/20questions /entrypoint
ENTRYPOINT ["/entrypoint"]
