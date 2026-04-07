FROM golang:1.24 AS builder

ARG APP_VERSION=dev
ARG COMMIT_HASH=unknown
ARG BUILD_DATE=unknown

ENV GO111MODULE=on \
    GOTOOLCHAIN=auto

WORKDIR /builder
COPY . .
RUN go mod download && \
    go build -ldflags " \
    -s -w \
    -X github.com/AkimioJR/MediaWarp/internal/config.appVersion=${APP_VERSION} \
    -X github.com/AkimioJR/MediaWarp/internal/config.commitHash=${COMMIT_HASH} \
    -X github.com/AkimioJR/MediaWarp/internal/config.buildDate=${BUILD_DATE}" \
    -o MediaWarp

FROM alpine:latest
COPY --from=builder /builder/MediaWarp /MediaWarp

RUN chmod +x /MediaWarp

EXPOSE 9000
VOLUME ["/etc/localtime", "/etc/timezone", "/config", "/logs", "/custom"]
ENTRYPOINT ["/MediaWarp"]
