package handler

import (
	"fmt"
	"net/http"
	"strings"

	"game-platform/internal/model"
	"game-platform/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthHandler handler для аутентификации
type AuthHandler struct {
	userRepo  model.UserRepo
	portalAPI *service.PortalAPI
	jwtSecret string
}

// NewAuthHandler создает новый AuthHandler
func NewAuthHandler(userRepo model.UserRepo, portalAPI *service.PortalAPI) *AuthHandler {
	return &AuthHandler{userRepo: userRepo, portalAPI: portalAPI}
}

// SetJWTSecret sets the JWT secret used for token validation.
// Must be called before VerifyToken if authentication is required.
func (h *AuthHandler) SetJWTSecret(secret string) {
	h.jwtSecret = secret
}

// VerifyToken проверяет JWT токен от SSO
func (h *AuthHandler) VerifyToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid authorization header"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// If no secret configured, accept any token (dev/test mode)
	if h.jwtSecret == "" {
		c.JSON(http.StatusOK, gin.H{
			"valid":   true,
			"message": "Token accepted (dev mode — no JWT secret configured)",
		})
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(h.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
		return
	}

	// Build response with extracted claims
	result := gin.H{
		"valid": true,
	}
	if sid, ok := claims["sid"].(string); ok {
		result["sid"] = sid
	}
	if email, ok := claims["email"].(string); ok {
		result["email"] = email
	}
	if name, ok := claims["name"].(string); ok {
		result["name"] = name
	}
	if groups, ok := claims["groups"]; ok {
		result["groups"] = groups
	}

	// Optionally verify with Portal API if configured
	if h.portalAPI != nil && h.portalAPI.BaseURL() != "" {
		if sid, ok := claims["sid"].(string); ok {
			if _, err := h.portalAPI.GetUser(c.Request.Context(), sid); err != nil {
				// Portal verification failed — still accept JWT but log warning
				result["portal_verified"] = false
			} else {
				result["portal_verified"] = true
			}
		}
	}

	c.JSON(http.StatusOK, result)
}

// UserHandler handler для пользователей
type UserHandler struct {
	userRepo       model.UserRepo
	matchRepo      model.MatchRepo
	tournamentRepo model.TournamentRepo
	ratingRepo     interface{} // TODO: RatingRepository
}

// NewUserHandler создает новый UserHandler
func NewUserHandler(userRepo model.UserRepo, ratingRepo interface{}, matchRepo model.MatchRepo, tournamentRepo model.TournamentRepo) *UserHandler {
	return &UserHandler{userRepo: userRepo, ratingRepo: ratingRepo, matchRepo: matchRepo, tournamentRepo: tournamentRepo}
}

// GetProfile получает профиль пользователя
func (h *UserHandler) GetProfile(c *gin.Context) {
	sid := c.Param("sid")

	user, err := h.userRepo.GetBySID(c.Request.Context(), sid)
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetStats получает статистику пользователя
func (h *UserHandler) GetStats(c *gin.Context) {
	sid := c.Param("sid")

	stats := model.PlayerStats{}

	if h.matchRepo != nil {
		matchStats, err := h.matchRepo.GetPlayerStats(c.Request.Context(), sid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get match stats"})
			return
		}
		if matchStats != nil {
			stats.GamesPlayed = matchStats.GamesPlayed
			stats.Wins = matchStats.Wins
			stats.Draws = matchStats.Draws
			stats.Losses = matchStats.Losses
		}
	}

	if h.tournamentRepo != nil {
		tournaments, err := h.tournamentRepo.CountPlayerTournaments(c.Request.Context(), sid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get tournament count"})
			return
		}
		stats.TournamentsJoined = tournaments
	}

	c.JSON(http.StatusOK, gin.H{
		"sid":                sid,
		"games_played":       stats.GamesPlayed,
		"wins":              stats.Wins,
		"draws":             stats.Draws,
		"losses":            stats.Losses,
		"tournaments_joined": stats.TournamentsJoined,
	})
}
