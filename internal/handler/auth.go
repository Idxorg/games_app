package handler

import (
	"net/http"

	"game-platform/internal/repository"
	"game-platform/internal/service"

	"github.com/gin-gonic/gin"
)

// AuthHandler handler для аутентификации
type AuthHandler struct {
	userRepo  *repository.UserRepository
	portalAPI *service.PortalAPI
}

// NewAuthHandler создает новый AuthHandler
func NewAuthHandler(userRepo *repository.UserRepository, portalAPI *service.PortalAPI) *AuthHandler {
	return &AuthHandler{userRepo: userRepo, portalAPI: portalAPI}
}

// VerifyToken проверяет JWT токен от SSO
func (h *AuthHandler) VerifyToken(c *gin.Context) {
	// TODO: Реализовать валидацию JWT
	c.JSON(http.StatusOK, gin.H{"valid": true, "message": "Token verified"})
}

// UserHandler handler для пользователей
type UserHandler struct {
	userRepo   *repository.UserRepository
	ratingRepo interface{} // TODO: RatingRepository
}

// NewUserHandler создает новый UserHandler
func NewUserHandler(userRepo *repository.UserRepository, ratingRepo interface{}) *UserHandler {
	return &UserHandler{userRepo: userRepo, ratingRepo: ratingRepo}
}

// GetProfile получает профиль пользователя
func (h *UserHandler) GetProfile(c *gin.Context) {
	sid := c.Param("sid")

	user, err := h.userRepo.GetBySID(c.Request.Context(), sid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetStats получает статистику пользователя
func (h *UserHandler) GetStats(c *gin.Context) {
	sid := c.Param("sid")

	// TODO: Получить статистику
	c.JSON(http.StatusOK, gin.H{
		"sid": sid,
		"games_played": 0,
		"tournaments_joined": 0,
	})
}
