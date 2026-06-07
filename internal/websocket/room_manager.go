package websocket

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"game-platform/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

// RoomManager manages all active game rooms.
type RoomManager struct {
	mu        sync.RWMutex
	rooms     map[string]*GameRoom
	jwtSecret string
	matchRepo model.MatchRepo
}

// NewRoomManager creates a new RoomManager.
func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms: make(map[string]*GameRoom),
	}
}

// NewRoomManagerWithDeps creates a RoomManager with dependencies (JWT secret, match repo).
func NewRoomManagerWithDeps(jwtSecret string, matchRepo model.MatchRepo) *RoomManager {
	return &RoomManager{
		rooms:     make(map[string]*GameRoom),
		jwtSecret: jwtSecret,
		matchRepo: matchRepo,
	}
}

// CreateRoom creates a new room for the given match and game type.
func (rm *RoomManager) CreateRoom(matchID, gameType string) *GameRoom {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if _, exists := rm.rooms[matchID]; exists {
		return rm.rooms[matchID]
	}

	room := NewGameRoom(matchID, gameType)
	if rm.matchRepo != nil {
		room.SetMatchRepo(rm.matchRepo)
	}
	rm.rooms[matchID] = room
	slog.Info("room created", "match_id", matchID, "game_type", gameType)
	return room
}

// GetRoom returns a room by match ID, or nil if not found.
func (rm *RoomManager) GetRoom(matchID string) *GameRoom {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return rm.rooms[matchID]
}

// RemoveRoom removes and cleans up a room.
func (rm *RoomManager) RemoveRoom(matchID string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if room, exists := rm.rooms[matchID]; exists {
		room.Cleanup()
		delete(rm.rooms, matchID)
		slog.Info("room removed", "match_id", matchID)
	}
}

// RoomCount returns the number of active rooms.
func (rm *RoomManager) RoomCount() int {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return len(rm.rooms)
}

// RouteMessage routes a client message to the appropriate room.
func (rm *RoomManager) RouteMessage(matchID string, sid string, conn *websocket.Conn, msgType string, payload []byte) bool {
	room := rm.GetRoom(matchID)
	if room == nil {
		// Room doesn't exist yet — could create on join
		if msgType == "join" {
			return rm.handleJoin(matchID, sid, conn, payload)
		}
		sendError(conn, "room_not_found", "room not found for match_id")
		return false
	}

	switch msgType {
	case "join":
		return rm.handleJoin(matchID, sid, conn, payload)
	case "move":
		room.HandleMove(sid, payload)
	case "resign":
		room.HandleResign(sid)
	case "draw_offer":
		room.HandleDrawOffer(sid)
	case "draw_accept":
		room.HandleDrawAccept(sid)
	case "draw_decline":
		room.HandleDrawDecline(sid)
	case "roll_dice":
		room.HandleRollDice(sid)
	default:
		sendError(conn, "unknown_message_type", "unknown message type: "+msgType)
	}

	return true
}

// handleJoin processes a join request, creating the room if needed.
func (rm *RoomManager) handleJoin(matchID, sid string, conn *websocket.Conn, payload []byte) bool {
	// Parse game_type from the join payload
	var joinMsg struct {
		GameType string `json:"game_type"`
	}
	if err := json.Unmarshal(payload, &joinMsg); err != nil || joinMsg.GameType == "" {
		sendError(conn, "invalid_join", "game_type is required")
		return false
	}

	room := rm.GetRoom(matchID)
	if room == nil {
		room = rm.CreateRoom(matchID, joinMsg.GameType)
	}

	if err := room.JoinPlayer(sid, conn); err != nil {
		slog.Error("join player error", "match_id", matchID, "sid", sid, "error", err)
		return false
	}

	// If game is over, schedule cleanup
	if room.IsGameOver() {
		rm.RemoveRoom(matchID)
	}

	return true
}

// HandleWebSocket handles WebSocket upgrade and routes messages to rooms.
// It supports JWT auth via:
//   - Bearer token in Authorization header (set by middleware)
//   - Query parameter ?token= (for WebSocket clients that can't set headers)
func (rm *RoomManager) HandleWebSocket(c *gin.Context) {
	sid := c.GetString("sid")
	if sid == "" {
		// Try query param token authentication
		sid = rm.AuthenticateQueryToken(c)
		if sid == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated: no valid token"})
			return
		}
	}

	matchID := c.Param("match_id")
	if matchID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "match_id is required"})
		return
	}

	// Verify match exists if matchRepo is available
	if rm.matchRepo != nil {
		match, err := rm.matchRepo.GetByID(context.Background(), matchID)
		if err != nil || match == nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "match not found"})
			return
		}
		// Verify the connecting player is part of this match
		if match.Player1SID != sid && match.Player2SID != sid {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "not a player in this match"})
			return
		}
	}

	conn, err := websocket.Upgrade(c.Writer, c.Request, nil, 4096, 0)
	if err != nil {
		slog.Error("websocket upgrade error", "sid", sid, "match_id", matchID, "error", err)
		return
	}

	slog.Info("websocket connected", "sid", sid, "match_id", matchID)

	client := NewWSClient(rm, sid, matchID, conn)
	go client.writePump()
	go client.readPump()
}

// ValidateTokenString validates a raw JWT token string and returns the SID.
// Returns empty string if invalid.
func (rm *RoomManager) ValidateTokenString(tokenString string) string {
	if tokenString == "" || rm.jwtSecret == "" {
		return ""
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(rm.jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return ""
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return ""
	}

	sid, _ := claims["sid"].(string)
	return sid
}

// AuthenticateQueryToken validates a JWT token from the ?token= query parameter.
// Returns the SID on success, empty string on failure.
func (rm *RoomManager) AuthenticateQueryToken(c *gin.Context) string {
	return rm.ValidateTokenString(c.Query("token"))
}

// SetMatchRepo sets the match repository (for dependency injection after creation).
func (rm *RoomManager) SetMatchRepo(repo model.MatchRepo) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.matchRepo = repo
}

// SetJWTSecret sets the JWT secret (for dependency injection after creation).
func (rm *RoomManager) SetJWTSecret(secret string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.jwtSecret = secret
}

// GenerateTestToken creates a JWT token for testing purposes.
func GenerateTestToken(jwtSecret, sid string) string {
	claims := jwt.MapClaims{
		"sid":   sid,
		"email": sid + "@test.com",
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(jwtSecret))
	return signed
}
