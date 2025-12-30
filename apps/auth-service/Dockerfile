# multi-stage Dockerfile for Auth Service 
#=========================
# Build stage
# =========================
FROM golang:1.25-alpine AS builder

# Enable reproducible builds
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

# Install CA certs (for HTTPS calls if any)
# apk is package manager for Alpine Linux.
RUN apk add --no-cache ca-certificates

# Copy go mod files first (better layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build statically linked binary
RUN go build -o auth-service ./main


# =========================
# Runtime stage
# =========================
FROM gcr.io/distroless/static-debian12

ENV TZ=UTC

# Create non-root user
USER nonroot:nonroot

WORKDIR /app

# Copy CA certs (needed for TLS) from previous stage
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy binary from previous stage
COPY --from=builder /app/auth-service /app/auth-service

# Expose application port
EXPOSE 8080

# Read-only root filesystem is enforced in Kubernetes,
# not in Dockerfile, but container is compatible

ENTRYPOINT ["/app/auth-service"]