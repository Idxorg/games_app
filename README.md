# Game Platform - Корпоративная игровая платформа

## Описание

Игровая платформа для корпоративных чемпионатов по шахматам, шашкам и нардам. Интегрирована с корпоративным порталом через SSO (OIDC) и использует внешние сервисы S3 и LiveKit.

## Архитектура

```
┌──────────────────────────────────────────────────────────────┐
│           КОРПОРАТИВНЫЙ ПОРТАЛ (Matrix SSO + Groups)          │
│  - SSO: авторизация по sid                                   │
│  - База пользователей                                        │
│  - S3: хранение файлов                                       │
│  - LiveKit: видео-чаты                                       │
└──────────────────────────────────────────────────────────────┘
                         │
                         ▼ (JWT + API)
┌──────────────────────────────────────────────────────────────┐
│         ИГРОВОЙ МИКРОСЕРВИС (Go + Gin + WebSocket)           │
│  - HTTP API (REST)                                          │
│  - WebSocket для игр в реальном времени                     │
│  - PostgreSQL + Redis                                       │
│  - S3 Client (интеграция с внешним S3)                      │
└──────────────────────────────────────────────────────────────┘
```

## Стек технологий

- **Backend:** Go 1.24 + Gin Framework
- **База данных:** PostgreSQL 15+
- **Кэширование:** Redis 7+
- **WebSocket:** Gorilla WebSocket
- **S3:** AWS SDK v2 (интеграция с внешним S3)
- **LiveKit:** SDK для видео-чатов (интеграция с внешним LiveKit)
- **Docker:** Контейнеризация

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

## Установка и запуск

### Локальная разработка

1. **Установить зависимости:**
```bash
go mod tidy
```

2. **Настроить окружение:**
```bash
cp .env.example .env
# Отредактировать .env с вашими данными
```

3. **Запустить PostgreSQL и Redis:**
```bash
docker-compose up -d db redis
```

4. **Выполнить миграции:**
```bash
psql -h localhost -U game -d game_platform -f migrations/001_init.sql
```

5. **Запустить сервер:**
```bash
go run cmd/server/main.go
```

### Docker

1. **Собрать и запустить:**
```bash
docker-compose up --build
```

2. **Проверить работу:**
```bash
curl http://localhost:8000/health
```

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
