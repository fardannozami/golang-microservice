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

# Build main application binary
RUN CGO_ENABLED=0 GOOS=linux go build -o inventory-service ./cmd

# Build seeder binary
RUN CGO_ENABLED=0 GOOS=linux go build -o seed ./cmd/seed

# Stage 2: Minimal runtime image
FROM alpine:latest

WORKDIR /app

# Install SSL certs (buat koneksi HTTPS / Postgres, biasanya dibutuhin)
RUN apk --no-cache add ca-certificates

# Copy binaries from builder
COPY --from=builder /app/inventory-service .
COPY --from=builder /app/seed .

# Copy entrypoint script
COPY scripts/entrypoint.sh .
RUN chmod +x entrypoint.sh

# Expose port (match sama config default = 9090)
EXPOSE 9090

# Environment variable untuk mengontrol apakah seeder dijalankan
ENV RUN_SEEDER=false

ENTRYPOINT ["/app/entrypoint.sh"]
