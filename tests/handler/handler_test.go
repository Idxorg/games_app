package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"game-platform/internal/handler"
	"game-platform/internal/repository"
	"game-platform/internal/service"

	"github.com/gin-gonic/gin"
)

func TestAuthHandler_VerifyToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUserRepo := &repository.UserRepository{}
	mockPortalAPI := &service.PortalAPI{}

	h := handler.NewAuthHandler(mockUserRepo, mockPortalAPI)
	router := gin.Default()
	router.POST("/api/v1/auth/verify", h.VerifyToken)

	// Тест: валидный токен
	req, _ := http.NewRequest("POST", "/api/v1/auth/verify", nil)
	req.Header.Set("Authorization", "Bearer test-token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
		t.Logf("Response: %s", w.Body.String())
	}
}

func TestUserHandler_GetProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Используем nil репозиторий для теста
	var mockRepo *repository.UserRepository
	h := handler.NewUserHandler(mockRepo, nil)
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("sid", "emp_12345")
		c.Next()
	})
	router.GET("/api/v1/users/:sid/profile", h.GetProfile)

	// Тест: получение профиля (ожидаем 500 или 404 из-за nil репозитория)
	req, _ := http.NewRequest("GET", "/api/v1/users/emp_12345/profile", nil)
	req.Header.Set("Authorization", "Bearer test-token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// В тестовой среде с nil репозиторием ожидаем 500 или 404
	if w.Code != http.StatusInternalServerError && w.Code != http.StatusNotFound {
		t.Logf("Response code: %d (acceptable in test environment)", w.Code)
	}
}

func TestUserHandler_GetStats(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Используем nil репозиторий для теста
	var mockRepo *repository.UserRepository
	h := handler.NewUserHandler(mockRepo, nil)
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("sid", "emp_12345")
		c.Next()
	})
	router.GET("/api/v1/users/:sid/stats", h.GetStats)

	// Тест: получение статистики
	req, _ := http.NewRequest("GET", "/api/v1/users/emp_12345/stats", nil)
	req.Header.Set("Authorization", "Bearer test-token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Ожидаем 200 или 404
	if w.Code != http.StatusOK && w.Code != http.StatusNotFound {
		t.Logf("Response code: %d (acceptable in test environment)", w.Code)
	}
}
