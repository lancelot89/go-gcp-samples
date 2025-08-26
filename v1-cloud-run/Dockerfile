# Build stage
FROM golang:1.24.6 AS builder

WORKDIR /app

# Copy go mod file first for better caching
COPY go.mod go.sum* ./
RUN go mod download

# Copy source code
COPY . .

# Build with optimization flags
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -o server .

# Final stage using distroless base image
FROM gcr.io/distroless/static-debian12:nonroot

# Copy the binary from builder
COPY --from=builder /app/server /server

# Use non-root user
USER nonroot:nonroot

# Expose port (informational)
EXPOSE 8080

# Run the binary
ENTRYPOINT ["/server"]
