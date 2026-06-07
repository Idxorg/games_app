package websocket

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// RoomManager manages all active game rooms.
type RoomManager struct {
	mu    sync.RWMutex
	rooms map[string]*GameRoom
}

// NewRoomManager creates a new RoomManager.
func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms: make(map[string]*GameRoom),
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
func (rm *RoomManager) HandleWebSocket(c *gin.Context) {
	sid := c.GetString("sid")
	if sid == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated: no sid in context"})
		return
	}

	matchID := c.Param("match_id")
	if matchID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "match_id is required"})
		return
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
