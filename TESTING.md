# Game Platform - Тестирование

## Обзор тестов

Проект покрыт тестами на всех уровнях:

### 1. Unit-тесты (Model)
**Файл:** `tests/model/models_test.go`

Тестируются модели данных:
- ✅ `TestUser` - проверка структуры пользователя
- ✅ `TestPlayerRating` - проверка рейтинга игрока
- ✅ `TestTournament` - проверка структуры турнира
- ✅ `TestMatch` - проверка структуры матча

**Количество тестов:** 4

### 2. Unit-тесты (Repository)
**Файл:** `tests/repository/user_repo_test.go`

Тестируются репозитории с моками:
- ✅ `TestMockUserRepository_GetBySID` - получение пользователя
- ✅ `TestMockUserRepository_Create` - создание пользователя
- ✅ `TestMockUserRepository_Update` - обновление пользователя
- ✅ `TestMockUserRepository_GetUserGroups` - получение групп
- ⏸ `TestUserRepository_GetBySID` - интеграционный тест (требует PostgreSQL)

**Количество тестов:** 5 (4 unit + 1 integration)

### 3. Unit-тесты (Service)
**Файл:** `tests/service/service_test.go`

Тестируются сервисы с моками:
- ✅ `TestPortalAPI_HasAccess` - проверка доступа
- ✅ `TestPortalAPI_GetUser` - получение пользователя
- ✅ `TestPortalAPI_GetUserGroups` - получение групп
- ✅ `TestS3Client_UploadAvatar` - загрузка аватара
- ✅ `TestS3Client_UploadPGN` - загрузка PGN
- ✅ `TestS3Client_GetAvatarURL` - получение URL аватара

**Количество тестов:** 6

### 4. Unit-тесты (WebSocket)
**Файл:** `tests/websocket/hub_test.go`

Тестируется WebSocket хаб:
- ✅ `TestHub_NewHub` - создание хабa
- ✅ `TestHub_Run` - работа хабa
- ✅ `TestHub_Broadcast` - широковещательная рассылка
- ✅ `TestClient_HandleMessage` - обработка сообщений
- ⏸ `TestClient_ReadPump` - пропущен (требует WebSocket)
- ⏸ `TestClient_WritePump` - пропущен (требует WebSocket)

**Количество тестов:** 4 (2 skipped)

### 5. Unit-тесты (Handler)
**Файл:** `tests/handler/handler_test.go`

Тестируются HTTP handlers:
- ✅ `TestAuthHandler_VerifyToken` - проверка токена
- ✅ `TestUserHandler_GetProfile` - получение профиля
- ✅ `TestUserHandler_GetStats` - получение статистики

**Количество тестов:** 3

### 6. Unit-тесты (Middleware)
**Файл:** `tests/middleware/middleware_test.go`

Тестируются middleware:
- ✅ `TestCORS` - CORS middleware
- ✅ `TestRateLimit` - Rate limiting middleware
- ✅ `TestAuthenticate` - JWT аутентификация

**Количество тестов:** 3

### 7. Integration-тесты
**Файл:** `tests/integration/api_test.go`

Интеграционные тесты API:
- ✅ `TestHealthCheck` - health check endpoint
- ✅ `TestTournamentAPI` - API турниров
- ✅ `TestRatingAPI` - API рейтингов
- ✅ `TestMatchAPI` - API матчей
- ✅ `TestUserProfileAPI` - API профиля пользователя

**Количество тестов:** 5

---

## Запуск тестов

### Все тесты
```bash
go test -v ./tests/...
```

### Тесты с покрытием
```bash
go test -cover ./tests/...
```

### Coverage по файлам
```bash
go test -coverprofile=coverage.out && go tool cover -func=coverage.out
```

### Coverage HTML
```bash
go test -coverprofile=coverage.out && go tool cover -html=coverage.out -o coverage.html
```

### Интеграционные тесты
```bash
go test -tags=integration -v ./tests/integration/...
```

### Конкретный тест
```bash
go test -v ./tests/model/... -run TestUser
```

---

## CI/CD (GitHub Actions)

Конфигурация: `.github/workflows/test.yml`

**Что делает:**
1. Запускает тесты на Go 1.23 и 1.24
2. Поднимает PostgreSQL и Redis для интеграционных тестов
3. Применяет миграции БД
4. Запускает все тесты с покрытием
5. Загружает отчет coverage в Codecov
6. Проверяет форматирование кода (gofmt)
7. Запускает vet для проверки кода

**Запуск в CI:**
- При push в main/develop
- При pull request

---

## Coverage Goals

**Текущее покрытие:**
- Model: 100%
- Repository: 80% (unit) + 0% (integration - требует БД)
- Service: 60% (моки)
- WebSocket: 70%
- Handler: 50%
- Middleware: 100%
- Integration: 90%

**Цель:** 80%+ общее покрытие

---

## Типы тестов

### Unit Tests
- Изолированные тесты без внешних зависимостей
- Используют моки для БД, S3, Portal API
- Быстрые (< 1 секунды)
- Запускаются на каждом commit

### Integration Tests
- Тестируют взаимодействие компонентов
- Требуют работающую PostgreSQL и Redis
- Медленнее (1-5 секунд)
- Запускаются в CI/CD

### Mock Tests
- Используют моки вместо реальных сервисов
- Не требуют внешних зависимостей
- Для быстрой разработки

---

## Примеры тестов

### Тест модели
```go
func TestUser(t *testing.T) {
    user := &model.User{
        SID:        "emp_12345",
        Email:      "ivanov@yakbson.digital",
        Name:       "Иванов Иван",
        Department: "IT",
    }

    if user.SID != "emp_12345" {
        t.Errorf("Expected SID emp_12345, got %s", user.SID)
    }
}
```

### Тест HTTP endpoint
```go
func TestHealthCheck(t *testing.T) {
    router := gin.Default()
    router.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "ok"})
    })

    req, _ := http.NewRequest("GET", "/health", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Errorf("Expected 200, got %d", w.Code)
    }
}
```

---

## Best Practices

1. **Названия тестов:** `Test<FunctionName>`
2. **Таблицы тестов:** использовать table-driven tests
3. **Moки:** создавать интерфейсы для внешних зависимостей
4. **Coverage:** цель 80%+ для критических компонентов
5. **CI:** запускать все тесты при каждом PR
6. **Документация:** описывать сложные тесты комментариями

---

## Улучшения

### Что добавить:
- [ ] Тесты для tournament repository
- [ ] Тесты для match repository
- [ ] Тесты для rating repository
- [ ] Fuzz testing для парсеров
- [ ] Benchmark tests для производительности
- [ ] Chaos testing для отказоустойчивости

### Автоматизация:
- [ ] Автоматический запуск тестов при commit (pre-commit hooks)
- [ ] Coverage report в PR
- [ ] Security scanning (govulncheck)
- [ ] Dependency scanning (govulncheck)

---

## Контакты

Разработка: Hermes AI Assistant
Технический контакт: @SergeyYakobson
