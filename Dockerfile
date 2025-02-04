# Stage 1: Build the Go binary
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install required build dependencies
RUN apk add --no-cache gcc musl-dev 

# Copy go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy the application source code
COPY . .

# Enable CGO for SQLite
ENV CGO_ENABLED=1

# Build the application
RUN go build -o main ./cmd/server

# Stage 2: Create a minimal image for running the application
FROM alpine:3.20

WORKDIR /app

# Install runtime dependencies required by SQLite
RUN apk add --no-cache sqlite-libs

# Copy the binary and data from the builder stage
COPY --from=builder /app/main /app/
COPY --from=builder /app/data/ /app/data/

# Ensure binary is executable
RUN chmod +x /app/main

# Expose port 8080
EXPOSE 8080

# Start the application
CMD ["./main"]
