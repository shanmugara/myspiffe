# Multi-stage Dockerfile for myspiffe
# Builder stage: compile the Go binary statically
FROM golang:1.25-alpine AS builder

# Install git (for fetching modules) and build dependencies
RUN apk add --no-cache git

WORKDIR /src

# Copy go.mod and go.sum first to leverage Docker layer cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source and build
COPY . .
# Build the main binary (output to /out/myspiffe)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags='-s -w' -o /out/myspiffe ./

# Runtime stage: small image with ca-certificates
FROM alpine:3.19
RUN apk add --no-cache ca-certificates

# Create non-root user
RUN addgroup -S app && adduser -S -G app app

# Copy binary from builder
COPY --from=builder /out/myspiffe /usr/local/bin/myspiffe
RUN chmod +x /usr/local/bin/myspiffe

ENV SPIFFE_ENDPOINT_SOCKET="unix:///tmp/spire-agent/public/api.sock"

USER app
WORKDIR /home/app

# Default command
ENTRYPOINT ["/usr/local/bin/myspiffe"]

