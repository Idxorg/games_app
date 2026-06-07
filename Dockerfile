# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Установка зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копирование кода
COPY . .

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Final stage
FROM alpine:3.19

WORKDIR /app

# Установка ca-certificates для HTTPS
RUN apk --no-cache add ca-certificates

# Копирование бинарника
COPY --from=builder /app/main .
COPY --from=builder /app/web ./web

# Переменные окружения
ENV PORT=8000
ENV DATABASE_URL=postgresql://game:***@db:5432/game_platform?sslmode=disable
ENV REDIS_URL=redis:6379
ENV S3_ENDPOINT=https://s3.yakbson.digital
ENV S3_BUCKET=yakbson-games
ENV S3_ACCESS_KEY=your-access-key
ENV S3_SECRET_KEY=your-secret-key
ENV PORTAL_API_URL=https://portal.yakbson.digital/api
ENV PORTAL_API_KEY=your-api-key
ENV LIVEKIT_URL=https://livekit.yakbson.digital
ENV LIVEKIT_API_KEY=your-api-key
ENV LIVEKIT_API_SECRET=your-api-secret

# Порт
EXPOSE 8000

# Запуск приложения
CMD ["./main"]
