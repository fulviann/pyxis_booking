# syntax=docker/dockerfile:1.6

# =========================
# Stage 1: Builder
# =========================
FROM golang:1.24-alpine AS builder

# Build deps
RUN apk add --no-cache git build-base

WORKDIR /app

# Copy go.mod & go.sum (agar cache optimal)
COPY go.mod go.sum ./

# Download dependency once (cached)
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy full source
COPY . .

# Build using cache
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -o app .

# =========================
# Stage 2: Final Image
# =========================
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 7777

CMD ["./app"]
