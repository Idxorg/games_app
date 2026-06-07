package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Client клиент WebSocket
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
	sid  string
}

// NewClient создает нового клиента
func NewClient(hub *Hub, sid string) *Client {
	return &Client{
		hub: hub,
		sid: sid,
		send: make(chan []byte, 256),
	}
}

// Hub центральный хаб WebSocket соединений
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

// NewHub создает новый хаб
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Register регистрирует клиента
func (h *Hub) Register(client *Client) {
	h.mu.Lock()
	h.clients[client] = true
	h.mu.Unlock()
}

// Unregister отписывает клиента
func (h *Hub) Unregister(client *Client) {
	h.mu.Lock()
	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.send)
	}
	h.mu.Unlock()
}

// Broadcast отправляет сообщение всем клиентам
func (h *Hub) Broadcast(message []byte) {
	h.broadcast <- message
}

// Run запускает хаб
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// HandleWebSocket обрабатывает WebSocket соединения
func (h *Hub) HandleWebSocket(c *gin.Context) {
	_ = c.Param("game_id") // game_id используется для маршрутизации
	
	conn, err := websocket.Upgrade(c.Writer, c.Request, nil, 4096, 0)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := NewClient(h, c.GetString("sid"))
	client.conn = conn

	client.hub.register <- client

	// Запуск goroutines для чтения и записи
	go client.writePump()
	go client.readPump()
}

// readPump читает сообщения от клиента
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(2 * time.Minute))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(2 * time.Minute))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Обработка сообщений от клиента
		c.HandleMessage(message)
	}
}

// HandleMessage обрабатывает сообщения от клиента
func (c *Client) HandleMessage(message []byte) {
	var msg map[string]interface{}
	if err := json.Unmarshal(message, &msg); err != nil {
		return
	}

	switch msg["type"] {
	case "join":
		c.handleJoin(msg)
	case "move":
		c.handleMove(msg)
	case "game_over":
		c.handleGameOver(msg)
	}
}

// handleJoin обрабатывает присоединение к игре
func (c *Client) handleJoin(msg map[string]interface{}) {
	// TODO: Добавить игрока в комнату игры
	log.Printf("Player %s joined game", c.sid)
}

// handleMove обрабатывает ход в игре
func (c *Client) handleMove(msg map[string]interface{}) {
	// TODO: Обработать ход
	log.Printf("Player %s made move: %v", c.sid, msg)
}

// handleGameOver обрабатывает завершение игры
func (c *Client) handleGameOver(msg map[string]interface{}) {
	// TODO: Завершить игру
	log.Printf("Game over: %v", msg)
}

// writePump отправляет сообщения клиенту
func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()

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
		}
	}
}
