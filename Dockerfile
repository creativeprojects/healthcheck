# syntax=docker/dockerfile:1

# Stage 1: Build
FROM golang:1.27-alpine AS builder

RUN adduser -D -u 10001 appuser

WORKDIR /app

# Copy source and build static binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/healthcheck .

# Stage 2: Minimal Runtime
FROM scratch

# Copy CA certificates for HTTPS support
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /app/healthcheck /healthcheck

USER 10001

ENTRYPOINT ["/healthcheck"]
