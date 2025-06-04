# Build stage
FROM golang:1.23.9 AS builder

WORKDIR /app

# First copy only the dependency files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o notification-service .

# Runtime stage
FROM alpine:latest

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/notification-service .
COPY --from=builder /app/.env .

# Install CA certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Expose the service port
EXPOSE 8083

# Command to run the application
CMD ["./notification-service"]