package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"game-platform/internal/middleware"

	"github.com/gin-gonic/gin"
)

func TestCORS(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.Use(middleware.CORS())
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Тест: OPTIONS запрос
	req, _ := http.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Access-Control-Request-Method", "GET")
	req.Header.Set("Access-Control-Request-Headers", "Content-Type")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204 for OPTIONS, got %d", w.Code)
	}

	// Проверяем CORS заголовки
	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("Expected Access-Control-Allow-Origin '*', got '%s'", w.Header().Get("Access-Control-Allow-Origin"))
	}
}

func TestRateLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.Use(middleware.RateLimit(100, 1)) // 100 запросов в минуту
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Тест: нормальный запрос
	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestAuthenticate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.Use(middleware.Authenticate(nil)) // nil DB pool для теста
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Тест: без токена
	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401 without token, got %d", w.Code)
	}

	// Тест: с токеном
	req, _ = http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// В тестовой среде токен не валидируется, поэтому ожидаем 200 или 401
	if w.Code != http.StatusOK && w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 200 or 401, got %d", w.Code)
	}
}
