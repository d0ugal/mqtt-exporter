# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git for version detection
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with version information
RUN VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev") && \
    COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown") && \
    BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -ldflags="-s -w \
        -X github.com/d0ugal/mqtt-exporter/internal/version.Version=$VERSION \
        -X github.com/d0ugal/mqtt-exporter/internal/version.Commit=$COMMIT \
        -X github.com/d0ugal/mqtt-exporter/internal/version.BuildDate=$BUILD_DATE" \
    -o mqtt-exporter ./cmd/main.go

# Final stage
FROM alpine:3.22.1

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/mqtt-exporter .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./mqtt-exporter"]
