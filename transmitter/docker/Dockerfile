# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy everything to /app
COPY . .

# Set working directory to transmitter
WORKDIR /app/transmitter

# Initialize go module with local shared dependency
RUN go mod edit -replace github.com/vlady-kotsev/lime-radio/shared=../shared
RUN go mod tidy
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/transmitter/main .
COPY --from=builder /app/transmitter/config/config.docker.toml ./config/config.toml
COPY --from=builder /app/transmitter/songs ./songs

# Create directory for database
RUN mkdir -p internal/repository/data

# Expose port
EXPOSE 8080

# Command to run
CMD ["./main"]