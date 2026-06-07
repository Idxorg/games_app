# Coverage configuration

# Команда для запуска тестов с покрытием
test: go test -v -cover ./...

# Команда для отображения покрытия по файлам
cover: go test -coverprofile=coverage.out && go tool cover -func=coverage.out

# Команда для отображения HTML покрытия
cover-html: go test -coverprofile=coverage.out && go tool cover -html=coverage.out -o coverage.html

# Интеграционные тесты
integration: go test -tags=integration -v ./tests/integration/...

# Все тесты
all: test integration

# Очистка
clean:
	rm -f coverage.out coverage.html
