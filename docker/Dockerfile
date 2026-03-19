FROM golang:1.24 AS builder

ARG APP_VERSION=dev
ARG COMMIT_HASH=unknown
ARG BUILD_DATE=unknown

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.io,direct \
    GOTOOLCHAIN=auto

WORKDIR /builder
COPY . .
RUN go mod download && \
    go build -ldflags " \
    -s -w \
    -X MediaWarp/internal/config.appVersion=${APP_VERSION} \
    -X MediaWarp/internal/config.commitHash=${COMMIT_HASH} \
    -X MediaWarp/internal/config.buildDate=${BUILD_DATE}" \
    -o MediaWarp

FROM alpine:latest
COPY --from=builder /builder/MediaWarp /MediaWarp

RUN chmod +x /MediaWarp

EXPOSE 9000
VOLUME ["/etc/localtime", "/etc/timezone", "/config", "/logs", "/custom"]
ENTRYPOINT ["/MediaWarp"]
