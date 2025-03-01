# Use a valid Go version
FROM golang:1.23.3 AS builder

WORKDIR /app

# Cache module downloads
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy the rest of the source code
COPY . .

# Build the Go binary
RUN go build -v -o app ./...

# Use a minimal image for production
FROM debian:bullseye-slim

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/app .

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./app"]
