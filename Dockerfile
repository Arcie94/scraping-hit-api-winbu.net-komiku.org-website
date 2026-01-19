# Build Frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package.json frontend/package-lock.json ./
RUN npm install
COPY frontend .
RUN npm run build

# Build Backend
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/api/main.go

# Final Stage
FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add ca-certificates tzdata

# Copy binary
COPY --from=builder /app/main .

# Copy frontend build
COPY --from=frontend-builder /app/frontend/dist ./dist

# Expose port (Internal)
EXPOSE 3000

# Run
CMD ["./main"]
