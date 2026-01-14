# Build stage
FROM golang:1.25.5-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./

# Download dependencies (cached if go.mod/go.sum unchanged)
RUN go mod download

# Copy source code
COPY . .

# Build binaries with optimizations
# -ldflags="-w -s" strips debug info for smaller binaries
# -trimpath removes file system paths from binaries
RUN CGO_ENABLED=0 GOOS=linux go build \
    -a -installsuffix cgo \
    -ldflags="-w -s" \
    -trimpath \
    -o /app/api ./cmd/api

RUN CGO_ENABLED=0 GOOS=linux go build \
    -a -installsuffix cgo \
    -ldflags="-w -s" \
    -trimpath \
    -o /app/worker ./cmd/worker

# Runtime stage - minimal alpine image
FROM alpine:3.19

# Install ca-certificates for HTTPS and wget for health checks
RUN apk --no-cache add ca-certificates wget && \
    addgroup -g 1000 appgroup && \
    adduser -D -u 1000 -G appgroup appuser

WORKDIR /app

# Copy binaries from builder
COPY --from=builder /app/api .
COPY --from=builder /app/worker .

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user for security
USER appuser

# Expose API port
EXPOSE 8080

# Run API server by default
CMD ["./api"]
