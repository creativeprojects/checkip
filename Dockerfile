# syntax=docker/dockerfile:1

# Stage 1: Build
FROM golang:1.27-alpine AS builder

RUN adduser -D -u 10001 appuser

WORKDIR /app

# Copy module files for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build static binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /app/server .

# Stage 2: Minimal Runtime
FROM scratch

# Copy CA certificates for HTTPS support
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy healthcheck binary
COPY --from=ghcr.io/creativeprojects/healthcheck:latest /healthcheck /healthcheck

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /app/server /server

USER 10001

# Expose port and run
EXPOSE 8133
ENTRYPOINT ["/server", "-host", "0.0.0.0", "-port", "8133"]
