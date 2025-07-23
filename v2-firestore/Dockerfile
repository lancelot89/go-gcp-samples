# Use the official Go image to build the application
FROM golang:1.24-alpine AS builder

# Set necessary environment variables for Go modules
ENV GO111MODULE=on
ENV CGO_ENABLED=0

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./ 

# Download all dependencies. Dependencies will be cached if go.mod and go.sum are not changed
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application
RUN go build -o /go-gcp-samples-v2-firestore ./cmd/server

# Use a minimal image for the final stage
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the compiled application from the builder stage
COPY --from=builder /go-gcp-samples-v2-firestore .

# Expose the port the application listens on
EXPOSE 8080

# Set environment variables for the application
ENV PORT=8080

# Command to run the application
CMD ["./go-gcp-samples-v2-firestore"]
