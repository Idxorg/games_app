# Game Platform - Корпоративная игровая платформа

## Описание

Игровая платформа для корпоративных чемпионатов по шахматам, шашкам и нардам. Интегрирована с корпоративным порталом через SSO (OIDC) и iframe embed (postMessage).

## Архитектура

```
┌──────────────────────────────────────────────────────────────┐
│           КОРПОРАТИВНЫЙ ПОРТАЛ (Matrix SSO + Groups)          │
│  - SSO: авторизация по sid                                   │
│  - База пользователей                                        │
└──────────────────────────────────────────────────────────────┘
                         │
                         ▼ (JWT + API)
┌──────────────────────────────────────────────────────────────┐
│         ИГРОВОЙ МИКРОСЕРВИС (Go + Gin + WebSocket)           │
│  - HTTP API (REST)                                          │
│  - WebSocket для игр в реальном времени                     │
│  - PostgreSQL + Redis                                       │
└──────────────────────────────────────────────────────────────┘
```

## Стек технологий

- **Backend:** Go 1.24 + Gin Framework
- **База данных:** PostgreSQL 15+
- **Кэширование:** Redis 7+
- **WebSocket:** Gorilla WebSocket
- **Docker:** Контейнеризация
- **Frontend:** React 19 + TypeScript + Vite 6 + Zustand

## Основные функции

1. **Авторизация через SSO** - интеграция с корпоративным порталом
2. **Турниры** - создание и участие в чемпионатах
3. **Рейтинговая система** - Elo для шахмат, шашек, нард
4. **Лидерборды** - таблицы лидеров по отделам и всем сотрудникам
5. **Игры в реальном времени** - WebSocket для шахмат, шашек, нард
6. **Видео-чаты** - интеграция с LiveKit для игровых сессий

## API Endpoints

### Авторизация
- `POST /api/v1/auth/verify` - Проверка JWT токена от SSO

### Пользователи
- `GET /api/v1/users/:sid/profile` - Профиль пользователя
- `GET /api/v1/users/:sid/stats` - Статистика пользователя

### Турниры
- `GET /api/v1/tournaments` - Список турниров
- `GET /api/v1/tournaments/:id` - Детали турнира
- `POST /api/v1/tournaments` - Создать турнир
- `POST /api/v1/tournaments/:id/join` - Записаться на турнир
- `POST /api/v1/tournaments/:id/leave` - Выйти из турнира

### Рейтинги
- `GET /api/v1/ratings/:game_type` - Рейтинги по игре
- `GET /api/v1/ratings/:game_type/leaderboard` - Лидерборд
- `GET /api/v1/ratings/departments/:game_type` - Рейтинги по отделам

### Игры
- `GET /api/v1/games/available` - Доступные игры
- `POST /api/v1/games/match/start` - Начать матч
- `POST /api/v1/games/match/complete` - Завершить матч
- `GET /api/v1/games/match/history/:sid` - История матчей
- `POST /api/v1/games/:game_id/video-chat` - Создать видео-чат

### WebSocket
- `WS /ws/game/:game_id` - Подключение к игре

## Quick Start (Docker)

1. **Create environment file:**
```bash
cp .env.example .env
# Edit .env — at minimum set DATABASE_URL, POSTGRES_PASSWORD, JWT_SECRET
```

2. **Start everything:**
```bash
make dev          # builds and starts all services (Ctrl+C to stop)
# or detached:
make dev-detach   # runs in background
```

3. **Verify:**
```bash
curl http://localhost:8000/health
```

The app is now running at **http://localhost:8000**. The Go server serves the React SPA from `web/dist/` with full client-side routing support.

### Useful commands

```bash
make dev          # Build and run full stack
make build        # Rebuild Docker images
make down         # Stop containers (keep data)
make clean        # Stop containers and remove volumes
make test         # Run Go tests
make migrate      # Run migrations against running DB
make frontend     # Start Vite dev server (frontend only)
make logs         # Tail application logs
make help         # Show all available targets
```

## Local Development (without Docker)

1. **Install dependencies:**
```bash
go mod tidy
cd web && npm install && cd ..
```

2. **Start PostgreSQL and Redis:**
```bash
docker compose up -d db redis
```

3. **Run migrations:**
```bash
make migrate
```

4. **Start backend + frontend dev servers:**
```bash
# Terminal 1 — Go backend:
go run cmd/server/main.go

# Terminal 2 — Vite dev server (with hot reload):
make frontend
```

The Vite dev server proxies `/api` requests to `localhost:8000`.

## Конфигурация

Переменные окружения:

```bash
# Сервер
PORT=8000
HOST=0.0.0.0

# База данных
DATABASE_URL=postgresql://game:***@localhost:5432/game_platform?sslmode=disable

# Redis
REDIS_URL=localhost:6379

# S3 (внешний)
S3_ENDPOINT=https://s3.yakbson.digital
S3_BUCKET=yakbson-games
S3_ACCESS_KEY=your-access-key
S3_SECRET_KEY=your-secret-key

# LiveKit (внешний)
LIVEKIT_URL=https://livekit.yakbson.digital
LIVEKIT_API_KEY=your-api-key
LIVEKIT_API_SECRET=your-api-secret

# Корп портал (SSO + API)
PORTAL_API_URL=https://portal.yakbson.digital/api
PORTAL_API_KEY=your-api-key
OIDC_JWKS_URL=https://portal.yakbson.digital/.well-known/jwks.json

# JWT
JWT_SECRET=your-jwt-secret
```

## Структура проекта

```
game-platform/
├── cmd/server/main.go          # Точка входа
├── internal/
│   ├── handler/                # HTTP handlers
│   ├── service/                # Бизнес-логика
│   ├── repository/             # Работа с БД
│   ├── websocket/              # WebSocket логика
│   ├── model/                  # Модели данных
│   └── middleware/             # Middleware
├── pkg/                        # Полезные утилиты
├── web/                        # Фронтенд (JS/HTML/CSS)
├── migrations/                 # Миграции БД
├── config/                     # Конфигурация
├── Dockerfile
├── docker-compose.yml
└── go.mod
```

## Разработка

### Добавление новой игры

1. Создать файл в `pkg/games/<game_name>/`
2. Реализовать логику игры
3. Добавить WebSocket обработчики в `internal/websocket/`
4. Добавить endpoint для запуска игры

### Интеграция с SSO

1. Получить публичный ключ от OIDC провайдера
2. Валидировать JWT токен в middleware
3. Синхронизировать пользователя из корп портала
4. Проверить группы доступа

## Тестирование

```bash
# Запустить все тесты
go test ./...

# Запустить тесты с покрытием
go test -cover ./...
```

## Мониторинг и логирование

- Логи записываются в stdout/stderr
- Health check endpoint: `/health`
- Prometheus метрики (опционально)

## Безопасность

- JWT аутентификация через SSO
- CORS middleware
- Rate limiting
- Валидация входных данных
- SQL injection защита (pgx)

## Масштабирование

- Горизонтальное масштабирование WebSocket серверов
- Redis для кэширования и сессий
- PostgreSQL репликация для чтения
- CDN для статических файлов

## Лицензия

Внутреннее использование Yakbson Digital

## Контакты

Разработка: Hermes AI Assistant
Технический контакт: @SergeyYakobson
