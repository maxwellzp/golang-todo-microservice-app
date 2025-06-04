FROM golang:1.23.9 AS builder

WORKDIR /app

# First copy only the dependency files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o api-gateway .

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/api-gateway .
COPY --from=builder /app/.env .

EXPOSE 8080
CMD ["./api-gateway"]
