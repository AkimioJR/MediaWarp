FROM golang:1.24 AS builder
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.io,direct

WORKDIR /builder
COPY . .
RUN go mod download && \
    go build -ldflags " \
    -s -w \
    -X MediaWarp/internal/config.commitHash=$(git rev-parse HEAD) \
    -X MediaWarp/internal/config.buildDate=$(date -u '+%Y-%m-%d %H:%M:%S')" \
    -o MediaWarp

FROM alpine:latest
COPY --from=builder /builder/MediaWarp /MediaWarp

RUN chmod +x /MediaWarp

EXPOSE 9000
VOLUME ["/etc/localtime", "/etc/timezone", "/config", "/logs", "/custom"]
ENTRYPOINT ["/MediaWarp"]
