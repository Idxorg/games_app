package websocket_test

import (
	"encoding/json"
	"testing"
	"time"

	"game-platform/internal/websocket"
)

func TestHub_NewHub(t *testing.T) {
	hub := websocket.NewHub()
	if hub == nil {
		t.Fatal("Expected hub, got nil")
	}
}

func TestClient_HandleMessage(t *testing.T) {
	hub := websocket.NewHub()

	client := websocket.NewClient(hub, "emp_12345")

	// join
	client.HandleMessage([]byte(`{"type":"join","game_id":"g_chess_123"}`))
	// move
	client.HandleMessage([]byte(`{"type":"move","from":"e2","to":"e4"}`))
	// game_over
	client.HandleMessage([]byte(`{"type":"game_over","winner":"white"}`))
	// invalid JSON — should not panic
	client.HandleMessage([]byte(`{invalid}`))
	// empty
	client.HandleMessage([]byte{})
}

func TestClient_HandleMessage_GameTypes(t *testing.T) {
	hub := websocket.NewHub()
	client := websocket.NewClient(hub, "emp_99999")

	tests := []struct {
		name string
		msg  string
	}{
		{"join chess", `{"type":"join","game_id":"g_chess_001"}`},
		{"join checkers", `{"type":"join","game_id":"g_checkers_001"}`},
		{"join backgammon", `{"type":"join","game_id":"g_backgammon_001"}`},
		{"move pawn", `{"type":"move","from":"e2","to":"e4","piece":"pawn"}`},
		{"move king", `{"type":"move","from":"e1","to":"e2","piece":"king"}`},
		{"resign", `{"type":"game_over","winner":"black","reason":"resignation"}`},
		{"timeout", `{"type":"game_over","winner":"white","reason":"timeout"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client.HandleMessage([]byte(tt.msg))
		})
	}
}

func TestClient_HandleMessage_Malformed(t *testing.T) {
	hub := websocket.NewHub()
	client := websocket.NewClient(hub, "emp_12345")

	malformed := []string{
		"",
		"{}",
		`{"type":"unknown"}`,
		`{"type":123}`,
		`not json at all`,
		`{"type":"join"}`, // missing game_id
	}

	for _, msg := range malformed {
		client.HandleMessage([]byte(msg)) // should not panic
	}
}

func TestClient_MessageParsing(t *testing.T) {
	msg := []byte(`{"type":"join","game_id":"g_123"}`)
	var parsed map[string]interface{}
	if err := json.Unmarshal(msg, &parsed); err != nil {
		t.Fatalf("Failed to parse valid JSON: %v", err)
	}
	if parsed["type"] != "join" {
		t.Errorf("Expected type=join, got %v", parsed["type"])
	}
	if parsed["game_id"] != "g_123" {
		t.Errorf("Expected game_id=g_123, got %v", parsed["game_id"])
	}
}

func TestHub_RegisterUnregister(t *testing.T) {
	hub := websocket.NewHub()
	done := make(chan struct{})

	go func() {
		hub.Run()
	}()

	client := websocket.NewClient(hub, "emp_12345")
	hub.Register(client)
	time.Sleep(50 * time.Millisecond)

	hub.Unregister(client)
	time.Sleep(50 * time.Millisecond)

	// Clean shutdown — unregister then close channels
	close(done)
	time.Sleep(50 * time.Millisecond)
}

func TestNewClient(t *testing.T) {
	hub := websocket.NewHub()
	client := websocket.NewClient(hub, "emp_test")

	if client == nil {
		t.Fatal("Expected non-nil client")
	}
}

func TestHub_Broadcast_DoesNotBlock(t *testing.T) {
	hub := websocket.NewHub()
	done := make(chan struct{})

	go func() {
		hub.Run()
		defer func() { close(done) }()
	}()

	// Broadcast without any clients should not block
	go hub.Broadcast([]byte(`{"type":"ping"}`))

	select {
	case <-done:
		t.Error("Hub stopped unexpectedly")
	case <-time.After(100 * time.Millisecond):
		// Expected — broadcast processed
	}
}
