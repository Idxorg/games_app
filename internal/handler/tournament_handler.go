package handler

import (
	"net/http"
	"strconv"
	"time"

	"game-platform/internal/model"
	"game-platform/internal/repository"
	"game-platform/internal/service"

	"github.com/gin-gonic/gin"
)

// TournamentHandler handles tournament CRUD and player management endpoints.
type TournamentHandler struct {
	tournamentRepo *repository.TournamentRepository
	userRepo      *repository.UserRepository
}

// NewTournamentHandler creates a new TournamentHandler.
func NewTournamentHandler(
	tournamentRepo *repository.TournamentRepository,
	userRepo *repository.UserRepository,
) *TournamentHandler {
	return &TournamentHandler{
		tournamentRepo: tournamentRepo,
		userRepo:       userRepo,
	}
}

// createTournamentRequest represents the JSON body for creating a tournament.
type createTournamentRequest struct {
	Name          string     `json:"name" binding:"required"`
	GameType      string     `json:"game_type" binding:"required"`
	MaxPlayers    int        `json:"max_players" binding:"required,min=2"`
	StartDate     *time.Time `json:"start_date"`
	EndDate       *time.Time `json:"end_date"`
	PrizePool     string     `json:"prize_pool"`
	Description   string     `json:"description"`
	LogoURL       string     `json:"logo_url"`
	RequiresGroup string     `json:"requires_group"`
}

// updateTournamentRequest represents the JSON body for updating a tournament.
type updateTournamentRequest struct {
	Name          *string    `json:"name"`
	GameType      *string    `json:"game_type"`
	Status        *string    `json:"status"`
	MaxPlayers    *int       `json:"max_players" binding:"omitempty,min=2"`
	StartDate     *time.Time `json:"start_date"`
	EndDate       *time.Time `json:"end_date"`
	PrizePool     *string    `json:"prize_pool"`
	Description   *string    `json:"description"`
	LogoURL       *string    `json:"logo_url"`
	RequiresGroup *string    `json:"requires_group"`
}

// validTournamentStatuses enumerates allowed tournament statuses.
var validTournamentStatuses = map[string]bool{
	"upcoming":   true,
	"active":     true,
	"completed":  true,
	"cancelled":  true,
}

// ListTournaments handles GET /api/v1/tournaments.
// Supports query params: game_type, status, search, limit, offset.
func (h *TournamentHandler) ListTournaments(c *gin.Context) {
	filters := model.TournamentFilters{
		GameType:  c.Query("game_type"),
		Status:    c.Query("status"),
		Search:    c.Query("search"),
		CreatedBy: c.Query("created_by"),
	}

	if limit, err := strconv.Atoi(c.DefaultQuery("limit", "50")); err == nil && limit > 0 {
		filters.Limit = limit
	} else {
		filters.Limit = 50
	}
	if offset, err := strconv.Atoi(c.DefaultQuery("offset", "0")); err == nil && offset >= 0 {
		filters.Offset = offset
	}

	// Validate game_type if provided
	if filters.GameType != "" && !service.ValidGameType(filters.GameType) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game_type"})
		return
	}

	tournaments, err := h.tournamentRepo.List(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list tournaments"})
		return
	}

	if tournaments == nil {
		tournaments = []model.Tournament{}
	}
	c.JSON(http.StatusOK, gin.H{
		"tournaments": tournaments,
		"count":      len(tournaments),
	})
}

// CreateTournament handles POST /api/v1/tournaments.
func (h *TournamentHandler) CreateTournament(c *gin.Context) {
	var req createTournamentRequest
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

	// Get creator SID from auth context
	sid, exists := c.Get("sid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}
	creatorSID, ok := sid.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid sid in token"})
		return
	}

	tournament := &model.Tournament{
		ID:            service.GenerateID(),
		Name:          req.Name,
		GameType:      req.GameType,
		Status:        "upcoming",
		MaxPlayers:    req.MaxPlayers,
		CurrentPlayers: 0,
		PrizePool:     req.PrizePool,
		Description:   req.Description,
		LogoURL:       req.LogoURL,
		CreatedBy:     creatorSID,
		RequiresGroup: req.RequiresGroup,
	}
	if req.StartDate != nil {
		tournament.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		tournament.EndDate = *req.EndDate
	}

	created, err := h.tournamentRepo.Create(c.Request.Context(), tournament)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create tournament"})
		return
	}

	c.JSON(http.StatusCreated, created)
}

// GetTournament handles GET /api/v1/tournaments/:id.
func (h *TournamentHandler) GetTournament(c *gin.Context) {
	id := c.Param("id")

	tournament, err := h.tournamentRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tournament not found"})
		return
	}

	c.JSON(http.StatusOK, tournament)
}

// UpdateTournament handles PUT /api/v1/tournaments/:id.
func (h *TournamentHandler) UpdateTournament(c *gin.Context) {
	id := c.Param("id")

	// Check tournament exists
	existing, err := h.tournamentRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tournament not found"})
		return
	}

	var req updateTournamentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	// Validate game_type if being changed
	if req.GameType != nil && !service.ValidGameType(*req.GameType) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid game_type",
			"valid": service.AllGameTypes(),
		})
		return
	}

	// Validate status if being changed
	if req.Status != nil && !validTournamentStatuses[*req.Status] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "invalid status",
			"valid":  []string{"upcoming", "active", "completed", "cancelled"},
		})
		return
	}

	// Apply partial updates
	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.GameType != nil {
		existing.GameType = *req.GameType
	}
	if req.Status != nil {
		existing.Status = *req.Status
	}
	if req.MaxPlayers != nil {
		existing.MaxPlayers = *req.MaxPlayers
	}
	if req.StartDate != nil {
		existing.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		existing.EndDate = *req.EndDate
	}
	if req.PrizePool != nil {
		existing.PrizePool = *req.PrizePool
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.LogoURL != nil {
		existing.LogoURL = *req.LogoURL
	}
	if req.RequiresGroup != nil {
		existing.RequiresGroup = *req.RequiresGroup
	}

	if err := h.tournamentRepo.Update(c.Request.Context(), existing); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update tournament"})
		return
	}

	c.JSON(http.StatusOK, existing)
}

// DeleteTournament handles DELETE /api/v1/tournaments/:id.
func (h *TournamentHandler) DeleteTournament(c *gin.Context) {
	id := c.Param("id")

	// Verify existence first
	_, err := h.tournamentRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tournament not found"})
		return
	}

	if err := h.tournamentRepo.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete tournament"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "tournament deleted"})
}

// JoinTournament handles POST /api/v1/tournaments/:id/join.
func (h *TournamentHandler) JoinTournament(c *gin.Context) {
	id := c.Param("id")

	// Get SID from auth context
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

	// Check tournament exists and is joinable
	tournament, err := h.tournamentRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tournament not found"})
		return
	}
	if tournament.Status != "upcoming" && tournament.Status != "active" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot join a tournament that is not upcoming or active"})
		return
	}
	if tournament.CurrentPlayers >= tournament.MaxPlayers {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tournament is full"})
		return
	}

	// Check group restriction if set
	if tournament.RequiresGroup != "" && h.userRepo != nil {
		groups, err := h.userRepo.GetUserGroups(c.Request.Context(), playerSID)
		if err != nil || !containsGroup(groups, tournament.RequiresGroup) {
			c.JSON(http.StatusForbidden, gin.H{"error": "you do not have the required group access to join this tournament"})
			return
		}
	}

	if err := h.tournamentRepo.AddPlayer(c.Request.Context(), id, playerSID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to join tournament"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "joined tournament"})
}

// LeaveTournament handles POST /api/v1/tournaments/:id/leave.
func (h *TournamentHandler) LeaveTournament(c *gin.Context) {
	id := c.Param("id")

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

	// Verify tournament exists and is still joinable
	tournament, err := h.tournamentRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tournament not found"})
		return
	}
	if tournament.Status == "completed" || tournament.Status == "cancelled" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot leave a completed or cancelled tournament"})
		return
	}

	if err := h.tournamentRepo.RemovePlayer(c.Request.Context(), id, playerSID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to leave tournament"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "left tournament"})
}

// GetTournamentPlayers handles GET /api/v1/tournaments/:id/players.
func (h *TournamentHandler) GetTournamentPlayers(c *gin.Context) {
	id := c.Param("id")

	// Verify tournament exists
	if _, err := h.tournamentRepo.GetByID(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tournament not found"})
		return
	}

	players, err := h.tournamentRepo.GetPlayers(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get tournament players"})
		return
	}

	if players == nil {
		players = []model.TournamentPlayer{}
	}
	c.JSON(http.StatusOK, gin.H{
		"players": players,
		"count":   len(players),
	})
}

// containsGroup is a helper to check if a group is in a list.
func containsGroup(groups []string, target string) bool {
	for _, g := range groups {
		if g == target {
			return true
		}
	}
	return false
}
