# --- Build stage ---
FROM golang:1.22-alpine AS builder
WORKDIR /app

# Copy go mod trước để cache dependency layer
COPY go.mod go.sum ./
RUN go mod download

# Copy toàn bộ code
COPY . .

# Build binary
RUN go build -o server ./cmd/server

# --- Run stage ---
FROM alpine:latest
WORKDIR /app

# Copy binary từ build stage
COPY --from=builder /app/server .

# Copy file cấu hình (nếu cần)
COPY .env .env

# Cổng server
EXPOSE 8080

CMD ["./server"]
