# Build stage
FROM golang:1.26.1-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /build/ghoper-strike-news ./cmd/bot

# Final stage
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

RUN adduser -D -g '' appuser

WORKDIR /app

COPY --from=builder /build/ghoper-strike-news /app/ghoper-strike-news

RUN mkdir -p /data && chown -R appuser:appuser /data /app

USER appuser

ENV DATABASE_URL=/data/cs2bot.db

VOLUME ["/data"]

ENTRYPOINT ["/app/ghoper-strike-news"]
