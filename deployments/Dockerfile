# syntax=docker/dockerfile:1

# Stage 1: Build binary
FROM golang:1.23-alpine AS builder

WORKDIR /src

# Install build tools
RUN apk add --no-cache git ca-certificates

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY cmd/ ./cmd/
COPY internal/ ./internal/

# Build static binaries for production
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/harshad-ssh ./cmd/ssh
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/harshad-tui ./cmd/cli

# Stage 2: Minimal runtime
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

# Create dedicated non-root user
RUN addgroup -S -g 10001 harshad && \
    adduser -S -u 10001 -G harshad -h /home/harshad harshad

WORKDIR /app

COPY --from=builder /bin/harshad-ssh /app/harshad-ssh
COPY --from=builder /bin/harshad-tui /app/harshad-tui

RUN mkdir -p /app/.ssh && chown -R harshad:harshad /app

USER harshad

ENV SSH_HOST=0.0.0.0
ENV SSH_PORT=2222
ENV SSH_KEY_PATH=/app/.ssh/harshad_ed25519

EXPOSE 2222

ENTRYPOINT ["/app/harshad-ssh"]
