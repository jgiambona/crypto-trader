FROM golang:1.10-alpine as builder
WORKDIR /app
COPY . /app/
RUN dep ensure && \
    go build

FROM debian:8
MAINTAINER ffimnsr <ffimnsr@outlook.com>
COPY --from=builder /app/trader /app/
EXPOSE 80

