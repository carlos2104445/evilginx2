# Multi-stage build for Go application
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o evilginx2 .

# Build React frontend
FROM node:18-alpine AS frontend-builder

WORKDIR /app/web

# Copy package files
COPY web/package*.json ./

# Install dependencies
RUN npm ci --only=production

# Copy web source
COPY web/ ./

# Build frontend
RUN npm run build

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S evilginx && \
    adduser -S -D -H -u 1001 -h /app -s /sbin/nologin -G evilginx -g evilginx evilginx

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/evilginx2 .

# Copy frontend build
COPY --from=frontend-builder /app/web/dist ./web/dist

# Copy configuration files
COPY --chown=evilginx:evilginx phishlets/ ./phishlets/
COPY --chown=evilginx:evilginx redirectors/ ./redirectors/

# Create directories for data
RUN mkdir -p /app/data /app/logs && \
    chown -R evilginx:evilginx /app

# Switch to non-root user
USER evilginx

# Expose ports
EXPOSE 8080 8443

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/health || exit 1

# Run the application
ENTRYPOINT ["./evilginx2"]
CMD ["-p", "./phishlets", "-t", "./redirectors", "-debug"]
