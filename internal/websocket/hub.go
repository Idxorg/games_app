package websocket

import (
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WSClient represents a single WebSocket client connected to a game room.
// It replaces the old Client struct, routing messages through RoomManager.
type WSClient struct {
	mu       sync.Mutex
	manager  *RoomManager
	conn     *websocket.Conn
	send     chan []byte
	sid      string
	matchID  string
}

// NewWSClient creates a new WebSocket client.
func NewWSClient(manager *RoomManager, sid, matchID string, conn *websocket.Conn) *WSClient {
	return &WSClient{
		manager: manager,
		conn:    conn,
		send:    make(chan []byte, 256),
		sid:     sid,
		matchID: matchID,
	}
}

// readPump reads messages from the WebSocket connection and routes them.
func (c *WSClient) readPump() {
	defer func() {
		c.manager.RemoveRoom(c.matchID) // cleanup room on disconnect
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(5 * time.Minute))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(5 * time.Minute))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error("websocket read error", "sid", c.sid, "match_id", c.matchID, "error", err)
			}
			slog.Info("websocket disconnected", "sid", c.sid, "match_id", c.matchID)
			break
		}

		c.handleMessage(message)
	}
}

// handleMessage dispatches incoming messages to the room manager.
func (c *WSClient) handleMessage(message []byte) {
	var msg struct {
		Type   string          `json:"type"`
		Payload json.RawMessage `json:"payload"`
	}
	if err := json.Unmarshal(message, &msg); err != nil {
		slog.Warn("invalid message format", "sid", c.sid, "error", err)
		return
	}

	switch msg.Type {
	case "join":
		// For join, the entire message is the payload (contains game_type)
		c.manager.RouteMessage(c.matchID, c.sid, c.conn, "join", message)
	case "move":
		c.manager.RouteMessage(c.matchID, c.sid, c.conn, "move", msg.Payload)
	case "resign":
		c.manager.RouteMessage(c.matchID, c.sid, c.conn, "resign", nil)
	case "draw_offer":
		c.manager.RouteMessage(c.matchID, c.sid, c.conn, "draw_offer", nil)
	case "draw_accept":
		c.manager.RouteMessage(c.matchID, c.sid, c.conn, "draw_accept", nil)
	case "draw_decline":
		c.manager.RouteMessage(c.matchID, c.sid, c.conn, "draw_decline", nil)
	case "roll_dice":
		c.manager.RouteMessage(c.matchID, c.sid, c.conn, "roll_dice", nil)
	default:
		slog.Warn("unknown message type", "sid", c.sid, "type", msg.Type)
	}
}

// writePump keeps the connection alive with pings.
func (c *WSClient) writePump() {
	defer func() {
		c.conn.Close()
	}()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
