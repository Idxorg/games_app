package handler

import (
	"crypto/subtle"
	"encoding/json"
	"net/http"
	"time"

	"game-platform/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// embedRequest represents the JSON body for the embed auth endpoint.
type embedRequest struct {
	SID        string `json:"sid"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	Department string `json:"department"`
}

// embedResponse represents the JSON response for the embed auth endpoint.
type embedResponse struct {
	Token string `json:"token"`
	SID   string `json:"sid"`
	Valid bool   `json:"valid"`
}

// EmbedHandler handles portal embed authentication.
type EmbedHandler struct {
	handoffSecret string
	jwtSecret     string
	jwtExpiryH    int // hours
	userRepo      model.UserRepo
}

// NewEmbedHandler creates a new EmbedHandler.
func NewEmbedHandler(handoffSecret, jwtSecret string, jwtExpiryHours int) *EmbedHandler {
	if jwtExpiryHours <= 0 {
		jwtExpiryHours = 24
	}
	return &EmbedHandler{
		handoffSecret: handoffSecret,
		jwtSecret:     jwtSecret,
		jwtExpiryH:    jwtExpiryHours,
	}
}

// SetUserRepo sets the user repository for upserting users during embed auth.
func (h *EmbedHandler) SetUserRepo(repo model.UserRepo) {
	h.userRepo = repo
}

// EmbedAuth handles POST /api/v1/auth/embed.
// It validates the X-Erlink-Embed-Secret header and issues a JWT token
// with claims: sid, email, name, department, groups.
func (h *EmbedHandler) EmbedAuth(c *gin.Context) {
	// 1. Verify embed handoff secret
	providedSecret := c.GetHeader("X-Erlink-Embed-Secret")
	if h.handoffSecret == "" {
		// Dev mode: no handoff secret configured, skip the check
	} else if subtle.ConstantTimeCompare([]byte(providedSecret), []byte(h.handoffSecret)) != 1 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid or missing embed secret",
		})
		return
	}

	// 2. Parse request body
	var req embedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body: " + err.Error(),
		})
		return
	}

	if req.SID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sid is required"})
		return
	}

	// 3. Upsert user in DB (if repo is configured)
	if h.userRepo != nil {
		existing, err := h.userRepo.GetBySID(c.Request.Context(), req.SID)
		if err == nil && existing != nil {
			// Update existing user
			existing.Email = req.Email
			existing.Name = req.Name
			existing.Department = req.Department
			existing.LastSync = time.Now()
			_ = h.userRepo.Update(c.Request.Context(), existing)
		} else {
			// Create new user
			_, _ = h.userRepo.Create(c.Request.Context(), req.SID, req.Email, req.Name, "", req.Department, "", "")
		}
	}

	// 4. Generate JWT with claims
	now := time.Now()
	claims := jwt.MapClaims{
		"sid":        req.SID,
		"email":      req.Email,
		"name":       req.Name,
		"department": req.Department,
		"groups":     []string{"games"},
		"iat":        now.Unix(),
		"exp":        now.Add(time.Duration(h.jwtExpiryH) * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	// 5. Return token
	resp := embedResponse{
		Token: tokenString,
		SID:   req.SID,
		Valid: true,
	}
	c.JSON(http.StatusOK, resp)
}

// decodeEmbedRequest is a test helper that parses JSON bytes into an embedRequest.
func decodeEmbedRequest(data []byte) (embedRequest, error) {
	var req embedRequest
	err := json.Unmarshal(data, &req)
	return req, err
}
