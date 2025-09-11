# Stage 1: Build the Go binary
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git (sering dipake buat ambil dependency)
RUN apk add --no-cache git

# Copy go mod files and download deps first (cache-friendly)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o inventory-service ./cmd

# Stage 2: Minimal runtime image
FROM alpine:latest

WORKDIR /app

# Install SSL certs (buat koneksi HTTPS / Postgres, biasanya dibutuhin)
RUN apk --no-cache add ca-certificates

# Copy binary from builder
COPY --from=builder /app/inventory-service .

# Expose port (match sama config default = 9090)
EXPOSE 9090

CMD ["./inventory-service"]
