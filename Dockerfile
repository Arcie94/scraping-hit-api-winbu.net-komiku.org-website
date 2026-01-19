# Build Stage
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Install git for fetching dependencies
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Tidy dependencies (fix versions)
RUN go mod tidy

# Debug: List files to ensure copy worked
RUN ls -R

# Build the application
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/api/main.go

# Final Stage
FROM alpine:latest
WORKDIR /app

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/main .

# Expose port (Internal)
EXPOSE 3000

# Run
CMD ["./main"]
