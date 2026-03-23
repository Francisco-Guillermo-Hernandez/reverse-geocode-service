# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy the entire project
COPY . .

# Build the statically linked Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server_bin ./cmd/server

# Final stage
FROM alpine:latest

# Working directory for the application
WORKDIR /app

# Install CA certificates just in case, though not strictly needed here
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /app/server_bin .

# Copy the data directory since the application expects data/countries.csv and data/geocode.json
COPY --from=builder /app/data ./data

# Expose the application port
EXPOSE 8080

# Run the binary
CMD ["./server_bin"]
