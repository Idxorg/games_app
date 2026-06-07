package handler

import (
	"net/http"
	"strconv"

	"game-platform/internal/model"
	"game-platform/internal/repository"
	"game-platform/internal/service"

	"github.com/gin-gonic/gin"
)

// RatingHandler handles paginated rating and department leaderboard endpoints.
type RatingHandler struct {
	ratingRepo *repository.RatingRepository
	ratingSvc  *service.RatingService
	userRepo   *repository.UserRepository
}

// NewRatingHandler creates a new RatingHandler.
func NewRatingHandler(
	ratingRepo *repository.RatingRepository,
	ratingSvc *service.RatingService,
	userRepo *repository.UserRepository,
) *RatingHandler {
	return &RatingHandler{
		ratingRepo: ratingRepo,
		ratingSvc:  ratingSvc,
		userRepo:   userRepo,
	}
}

// GetRatings handles GET /api/v1/ratings/:game_type.
// Returns a paginated list of player ratings for the given game type.
func (h *RatingHandler) GetRatings(c *gin.Context) {
	gameType := c.Param("game_type")

	if !service.ValidGameType(gameType) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid game_type",
			"valid": service.AllGameTypes(),
		})
		return
	}

	limit := 50
	if l, err := strconv.Atoi(c.DefaultQuery("limit", "50")); err == nil && l > 0 {
		limit = l
	}
	// Cap at 200
	if limit > 200 {
		limit = 200
	}

	ratings, err := h.ratingSvc.GetLeaderboard(c.Request.Context(), gameType, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get ratings"})
		return
	}

	if ratings == nil {
		ratings = []model.PlayerRating{}
	}

	c.JSON(http.StatusOK, gin.H{
		"game_type": gameType,
		"ratings":   ratings,
		"count":     len(ratings),
		"limit":     limit,
	})
}

// GetLeaderboardByDepartment handles GET /api/v1/ratings/:game_type/leaderboard/department.
// Requires a "department" query parameter.
func (h *RatingHandler) GetLeaderboardByDepartment(c *gin.Context) {
	gameType := c.Param("game_type")
	department := c.Query("department")

	if !service.ValidGameType(gameType) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid game_type",
			"valid": service.AllGameTypes(),
		})
		return
	}

	if department == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "department query parameter is required"})
		return
	}

	ratings, err := h.ratingRepo.GetByDepartment(c.Request.Context(), gameType, department)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get department leaderboard"})
		return
	}

	if ratings == nil {
		ratings = []model.PlayerRating{}
	}

	c.JSON(http.StatusOK, gin.H{
		"game_type":  gameType,
		"department": department,
		"ratings":    ratings,
		"count":      len(ratings),
	})
}
