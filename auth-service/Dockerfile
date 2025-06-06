# ---------- Builder Stage ----------
FROM golang:1.23.9 AS builder

WORKDIR /app

# Copy Go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service .

# ---------- Final Stage ----------
FROM alpine:latest

# Install root CA certificates (needed for HTTPS requests)
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the compiled binary and environment file
COPY --from=builder /app/auth-service .
# Optional, only if .env is used directly
COPY --from=builder /app/.env .

# Expose the application port
EXPOSE 8081

# Run the service
CMD ["./auth-service"]
