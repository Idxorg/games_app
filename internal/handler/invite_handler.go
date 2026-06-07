package handler

import (
	"log/slog"
	"net/http"
	"time"

	"game-platform/internal/model"
	"game-platform/internal/service"

	"github.com/gin-gonic/gin"
)

// InviteHandler handles game invite endpoints.
type InviteHandler struct {
	inviteRepo model.InviteRepo
	matchRepo  model.MatchRepo
}

// NewInviteHandler creates a new InviteHandler.
func NewInviteHandler(
	inviteRepo model.InviteRepo,
	matchRepo model.MatchRepo,
) *InviteHandler {
	return &InviteHandler{
		inviteRepo: inviteRepo,
		matchRepo:  matchRepo,
	}
}

// validInviteGameTypes defines which game types support invites.
var validInviteGameTypes = map[string]bool{
	"chess":     true,
	"checkers":  true,
	"backgammon": true,
}

// createInviteRequest is the JSON body for creating an invite.
type createInviteRequest struct {
	GameType     string `json:"game_type" binding:"required"`
	RecipientSID string `json:"recipient_sid" binding:"required"`
}

// CreateInvite handles POST /api/v1/games/invite.
func (h *InviteHandler) CreateInvite(c *gin.Context) {
	var req createInviteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	// Validate game_type
	if !validInviteGameTypes[req.GameType] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid game_type",
			"valid": []string{"chess", "checkers", "backgammon"},
		})
		return
	}

	// Get inviter SID from auth context
	sid, exists := c.Get("sid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}
	inviterSID, ok := sid.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid sid in token"})
		return
	}

	// Prevent self-invite
	if inviterSID == req.RecipientSID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot invite yourself"})
		return
	}

	invite := &model.GameInvite{
		GameType:     req.GameType,
		InviterSID:   inviterSID,
		RecipientSID: req.RecipientSID,
		Status:       "pending",
		ExpiresAt:    time.Now().UTC().Add(5 * time.Minute),
	}

	if err := h.inviteRepo.Create(c.Request.Context(), invite); err != nil {
		slog.Error("failed to create invite", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create invite"})
		return
	}

	c.JSON(http.StatusCreated, invite)
}

// AcceptInvite handles POST /api/v1/games/invite/:id/accept.
func (h *InviteHandler) AcceptInvite(c *gin.Context) {
	inviteID := c.Param("id")
	if inviteID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invite id required"})
		return
	}

	// Get recipient SID from auth context
	sid, exists := c.Get("sid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}
	recipientSID, ok := sid.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid sid in token"})
		return
	}

	// Fetch invite
	invite, err := h.inviteRepo.GetByID(c.Request.Context(), inviteID)
	if err != nil || invite == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "invite not found"})
		return
	}

	// Verify caller is the recipient
	if invite.RecipientSID != recipientSID {
		c.JSON(http.StatusForbidden, gin.H{"error": "only the recipient can accept this invite"})
		return
	}

	// Check status
	if invite.Status != "pending" {
		c.JSON(http.StatusConflict, gin.H{"error": "invite is not pending", "status": invite.Status})
		return
	}

	// Check expiry
	if time.Now().UTC().After(invite.ExpiresAt) {
		c.JSON(http.StatusGone, gin.H{"error": "invite has expired"})
		return
	}

	// Create a match
	now := time.Now().UTC()
	match := &model.Match{
		ID:         service.GenerateID(),
		GameType:   invite.GameType,
		Player1SID: invite.InviterSID,
		Player2SID: invite.RecipientSID,
		Status:     "in_progress",
		StartedAt:  &now,
	}

	createdMatch, err := h.matchRepo.Create(c.Request.Context(), match)
	if err != nil {
		slog.Error("failed to create match from invite", "error", err, "invite_id", inviteID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create match"})
		return
	}

	// Update invite with match ID
	if err := h.inviteRepo.Accept(c.Request.Context(), inviteID, createdMatch.ID); err != nil {
		slog.Error("failed to accept invite", "error", err, "invite_id", inviteID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to accept invite"})
		return
	}

	invite.Status = "accepted"
	invite.MatchID = &createdMatch.ID

	c.JSON(http.StatusOK, gin.H{
		"invite":  invite,
		"match":   createdMatch,
		"message": "invite accepted, match created",
	})
}

// DeclineInvite handles POST /api/v1/games/invite/:id/decline.
func (h *InviteHandler) DeclineInvite(c *gin.Context) {
	inviteID := c.Param("id")
	if inviteID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invite id required"})
		return
	}

	// Get recipient SID from auth context
	sid, exists := c.Get("sid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}
	recipientSID, ok := sid.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid sid in token"})
		return
	}

	// Fetch invite
	invite, err := h.inviteRepo.GetByID(c.Request.Context(), inviteID)
	if err != nil || invite == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "invite not found"})
		return
	}

	// Verify caller is the recipient
	if invite.RecipientSID != recipientSID {
		c.JSON(http.StatusForbidden, gin.H{"error": "only the recipient can decline this invite"})
		return
	}

	// Check status
	if invite.Status != "pending" {
		c.JSON(http.StatusConflict, gin.H{"error": "invite is not pending", "status": invite.Status})
		return
	}

	if err := h.inviteRepo.Decline(c.Request.Context(), inviteID); err != nil {
		slog.Error("failed to decline invite", "error", err, "invite_id", inviteID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decline invite"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"invite":  invite,
		"message": "invite declined",
	})
}

// GetPendingInvites handles GET /api/v1/games/invite/pending.
func (h *InviteHandler) GetPendingInvites(c *gin.Context) {
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

	invites, err := h.inviteRepo.GetPendingByRecipient(c.Request.Context(), playerSID)
	if err != nil {
		slog.Error("failed to get pending invites", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get pending invites"})
		return
	}

	if invites == nil {
		invites = []model.GameInvite{}
	}

	c.JSON(http.StatusOK, gin.H{
		"invites": invites,
		"count":   len(invites),
	})
}
