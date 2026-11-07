# syntax=docker/dockerfile:1

# Stage 1: Build
FROM golang:1.27-alpine AS builder

WORKDIR /app

# Copy module files for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build static binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOAMD64="v3" \
    go build -ldflags="-s -w" -o /app/server .

# Stage 2: Minimal Runtime
FROM scratch

# Copy CA certificates for HTTPS support
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy binary
COPY --from=builder /app/server /server

# Expose port and run
EXPOSE 8080
ENTRYPOINT ["/server"]
