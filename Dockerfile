# Use official Golang image for building
FROM golang:1.18 AS builder

WORKDIR /app

# Copy go.mod and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o cstv-go ./cmd/main.go

# Use a minimal image for the final runtime
FROM alpine:latest

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/cstv-go .

# Set environment variables (can be overridden in docker-compose)
ENV PORT=8080
ENV LOG_LEVEL=info

# Expose application port
EXPOSE 8080

# Run the application
CMD ["./cstv-go"]
