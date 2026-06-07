package handler

import (
	"log/slog"
	"net/http"
	"time"

	"game-platform/internal/model"
	"game-platform/internal/service"

	"github.com/gin-gonic/gin"
)

// NotificationPublisher is the interface for publishing events to a notification system.
type NotificationPublisher interface {
	PublishEvent(event map[string]interface{}) error
}

// InviteHandler handles game invite endpoints.
type InviteHandler struct {
	inviteRepo model.InviteRepo
	matchRepo  model.MatchRepo
	userRepo   model.UserRepo
	notifier   NotificationPublisher
}

// SetUserRepo sets the user repository (for resolving inviter names).
func (h *InviteHandler) SetUserRepo(repo model.UserRepo) {
	h.userRepo = repo
}

// SetNotificationsClient sets the notification publisher.
func (h *InviteHandler) SetNotificationsClient(client NotificationPublisher) {
	h.notifier = client
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

	// Resolve inviter name
	inviterName := ""
	if h.userRepo != nil {
		if u, err := h.userRepo.GetBySID(c.Request.Context(), inviterSID); err == nil && u != nil {
			inviterName = u.Name
		}
	}

	// Publish notification
	if h.notifier != nil {
		_ = h.notifier.PublishEvent(map[string]interface{}{
			"type":          "game.invite",
			"invite_id":     invite.ID,
			"game_type":     invite.GameType,
			"inviter_sid":   invite.InviterSID,
			"inviter_name":  inviterName,
			"recipient_sid": invite.RecipientSID,
		})
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":            invite.ID,
		"game_type":     invite.GameType,
		"inviter_sid":   invite.InviterSID,
		"inviter_name":  inviterName,
		"recipient_sid": invite.RecipientSID,
		"status":        invite.Status,
		"created_at":    invite.CreatedAt,
		"expires_at":    invite.ExpiresAt,
	})
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

	// Publish notification on accept
	if h.notifier != nil {
		_ = h.notifier.PublishEvent(map[string]interface{}{
			"type":          "game.invite.accepted",
			"invite_id":     invite.ID,
			"match_id":      createdMatch.ID,
			"game_type":     invite.GameType,
			"inviter_sid":   invite.InviterSID,
			"recipient_sid": invite.RecipientSID,
		})
	}

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

	// Publish notification on decline
	if h.notifier != nil {
		_ = h.notifier.PublishEvent(map[string]interface{}{
			"type":          "game.invite.declined",
			"invite_id":     inviteID,
			"inviter_sid":   invite.InviterSID,
			"recipient_sid": invite.RecipientSID,
		})
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

	// Resolve inviter names for each invite
	type inviteWithNames struct {
		model.GameInvite
		InviterName string `json:"inviter_name"`
	}
	result := make([]inviteWithNames, len(invites))
	for i, inv := range invites {
		name := ""
		if h.userRepo != nil {
			if u, err := h.userRepo.GetBySID(c.Request.Context(), inv.InviterSID); err == nil && u != nil {
				name = u.Name
			}
		}
		result[i] = inviteWithNames{GameInvite: inv, InviterName: name}
	}

	c.JSON(http.StatusOK, gin.H{
		"invites": result,
		"count":   len(result),
	})
}
