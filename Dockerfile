# Stage 1: Build
FROM golang:1.24-alpine AS builder

# Set environment variables
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire source
COPY . .

# Build the Go binary
RUN go build -o main .

# Stage 2: Run (minimal image)
FROM alpine:latest

# Set working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Expose port (optional, for services)
EXPOSE 8080

# Command to run the binary
CMD ["./main"]