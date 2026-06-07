package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"game-platform/internal/model"
	"game-platform/internal/service"

	"github.com/gin-gonic/gin"
)

// GameHandler handles game-related endpoints (available games, match management).
type GameHandler struct {
	matchRepo model.MatchRepo
	ratingSvc *service.RatingService
}

// NewGameHandler creates a new GameHandler.
func NewGameHandler(
	matchRepo model.MatchRepo,
	ratingSvc *service.RatingService,
) *GameHandler {
	return &GameHandler{
		matchRepo: matchRepo,
		ratingSvc: ratingSvc,
	}
}

// gameInfo describes a supported game for the /games/available endpoint.
type gameInfo struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MinPlayers  int    `json:"min_players"`
	MaxPlayers  int    `json:"max_players"`
}

// allGames lists every supported game with metadata.
var allGames = []gameInfo{
	{Type: "chess", Name: "Chess", Description: "Classic two-player strategy board game", MinPlayers: 2, MaxPlayers: 2},
	{Type: "checkers", Name: "Checkers", Description: "Traditional diagonal movement board game", MinPlayers: 2, MaxPlayers: 2},
	{Type: "backgammon", Name: "Backgammon", Description: "Classic dice-based board game", MinPlayers: 2, MaxPlayers: 2},
	{Type: "snake", Name: "Snake", Description: "Competitive snake game", MinPlayers: 1, MaxPlayers: 4},
	{Type: "mines", Name: "Mines", Description: "Competitive minesweeper", MinPlayers: 1, MaxPlayers: 4},
	{Type: "arena", Name: "Arena", Description: "Multi-player arena battle", MinPlayers: 2, MaxPlayers: 8},
	{Type: "poker", Name: "Poker", Description: "Card game with various variants", MinPlayers: 2, MaxPlayers: 10},
}

// AvailableGames handles GET /api/v1/games/available.
func (h *GameHandler) AvailableGames(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"games": allGames,
		"count": len(allGames),
	})
}

// startMatchRequest represents the JSON body for starting a new match.
type startMatchRequest struct {
	GameType      string `json:"game_type" binding:"required"`
	Player2SID    string `json:"player2_sid" binding:"required"`
	TournamentID  string `json:"tournament_id"`
	LiveKitRoomID string `json:"livekit_room_id"`
}

// StartMatch handles POST /api/v1/games/match/start.
func (h *GameHandler) StartMatch(c *gin.Context) {
	var req startMatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	// Validate game_type
	if !service.ValidGameType(req.GameType) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid game_type",
			"valid": service.AllGameTypes(),
		})
		return
	}

	// Get player1 SID from auth context
	sid, exists := c.Get("sid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}
	player1SID, ok := sid.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid sid in token"})
		return
	}

	// Prevent self-play
	if player1SID == req.Player2SID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot play against yourself"})
		return
	}

	now := time.Now().UTC()
	match := &model.Match{
		ID:            service.GenerateID(),
		TournamentID:  req.TournamentID,
		GameType:      req.GameType,
		Player1SID:    player1SID,
		Player2SID:    req.Player2SID,
		Status:        "in_progress",
		LiveKitRoomID: req.LiveKitRoomID,
		StartedAt:     &now,
	}

	created, err := h.matchRepo.Create(c.Request.Context(), match)
	if err != nil {
		slog.Error("failed to create match", "error", err, "game_type", req.GameType, "player1", player1SID, "player2", req.Player2SID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to start match"})
		return
	}

	c.JSON(http.StatusCreated, created)
}

// completeMatchRequest represents the JSON body for completing a match.
type completeMatchRequest struct {
	MatchID   string          `json:"match_id" binding:"required"`
	WinnerSID string          `json:"winner_sid"`
	Score     string          `json:"score" binding:"required"`
	Moves     json.RawMessage `json:"moves"`
}

// CompleteMatch handles POST /api/v1/games/match/complete.
// It records the match result and updates Elo ratings for both players.
func (h *GameHandler) CompleteMatch(c *gin.Context) {
	var req completeMatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	// Fetch the match
	match, err := h.matchRepo.GetByID(c.Request.Context(), req.MatchID)
	if err != nil || match == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "match not found"})
		return
	}

	if match.Status == "completed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "match already completed"})
		return
	}

	// Verify that the caller is one of the players
	callerSID, _ := c.Get("sid")
	caller, _ := callerSID.(string)
	if caller != match.Player1SID && caller != match.Player2SID {
		c.JSON(http.StatusForbidden, gin.H{"error": "only match participants can complete a match"})
		return
	}

	// If winner is provided, validate it's one of the players
	if req.WinnerSID != "" && req.WinnerSID != match.Player1SID && req.WinnerSID != match.Player2SID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "winner must be one of the match participants"})
		return
	}

	// Complete the match in the repository
	if err := h.matchRepo.Complete(c.Request.Context(), req.MatchID, req.WinnerSID, req.Score, req.Moves); err != nil {
		slog.Error("failed to complete match", "error", err, "match_id", req.MatchID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to complete match"})
		return
	}

	// Update match object for Elo calculation
	match.WinnerSID = req.WinnerSID
	match.Score = req.Score
	match.Status = "completed"

	// Update Elo ratings
	if h.ratingSvc != nil {
		if err := h.ratingSvc.UpdateMatchRatings(c.Request.Context(), match); err != nil {
			slog.Warn("failed to update ratings for match", "match_id", req.MatchID, "error", err)
			// Don't fail the request — match is already completed in DB
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "match completed",
		"match_id": req.MatchID,
	})
}

// GetMatchHistory handles GET /api/v1/games/match/history.
// Returns match history for the authenticated user, with optional game_type filter.
func (h *GameHandler) GetMatchHistory(c *gin.Context) {
	sid, exists := c.Get("sid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}
	playerSID, ok := sid.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid sid in token"})
		return
	}

	matches, err := h.matchRepo.ListByPlayer(c.Request.Context(), playerSID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get match history"})
		return
	}

	// Filter by game_type if requested
	gameType := c.Query("game_type")
	if gameType != "" {
		filtered := make([]model.Match, 0)
		for _, m := range matches {
			if m.GameType == gameType {
				filtered = append(filtered, m)
			}
		}
		matches = filtered
	}

	if matches == nil {
		matches = []model.Match{}
	}
	c.JSON(http.StatusOK, gin.H{
		"matches": matches,
		"count":   len(matches),
	})
}

// GetLeaderboard handles GET /api/v1/ratings/:game_type/leaderboard.
func (h *GameHandler) GetLeaderboard(c *gin.Context) {
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

	leaderboard, err := h.ratingSvc.GetLeaderboard(c.Request.Context(), gameType, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get leaderboard"})
		return
	}

	if leaderboard == nil {
		leaderboard = []model.PlayerRating{}
	}
	c.JSON(http.StatusOK, gin.H{
		"game_type":   gameType,
		"leaderboard": leaderboard,
		"count":       len(leaderboard),
	})
}

// GetUserRating handles GET /api/v1/ratings/:game_type/me.
func (h *GameHandler) GetUserRating(c *gin.Context) {
	gameType := c.Param("game_type")

	if !service.ValidGameType(gameType) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid game_type",
			"valid": service.AllGameTypes(),
		})
		return
	}

	sid, exists := c.Get("sid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}
	playerSID, ok := sid.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid sid in token"})
		return
	}

	rating, err := h.ratingSvc.GetRating(c.Request.Context(), playerSID, gameType)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "rating not found"})
		return
	}

	c.JSON(http.StatusOK, rating)
}
