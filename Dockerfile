# =============================================================================
# Stage 1: Build Frontend (React + Vite)
# =============================================================================
FROM node:22-alpine AS build-frontend

WORKDIR /app/web

COPY web/package.json web/package-lock.json* ./
RUN npm ci

COPY web/ ./
RUN npm run build

# =============================================================================
# Stage 2: Build Backend (Go)
# =============================================================================
FROM golang:1.24-alpine AS build-backend

WORKDIR /app

# Install dependencies first (layer cache)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY pkg/ ./pkg/

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# =============================================================================
# Stage 3: Runtime
# =============================================================================
FROM alpine:3.20

RUN apk --no-cache add ca-certificates curl

WORKDIR /app

# Copy binary from backend build
COPY --from=build-backend /app/main .

# Copy frontend build output
COPY --from=build-frontend /app/web/dist ./web/dist

# Copy migrations for runtime execution
COPY migrations/ ./migrations/

# Copy migration runner
COPY migrations/run.sh ./migrations/run.sh

# Non-secret defaults only
ENV PORT=8000
ENV SERVER_MODE=release

EXPOSE 8000

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD curl -f http://localhost:${PORT}/health || exit 1

# Run migrations then start the server
CMD ["sh", "-c", "./migrations/run.sh && ./main"]
