package websocket

import (
	"encoding/json"
	"testing"
)

// ---------- WSClient ----------

func TestNewWSClient(t *testing.T) {
	rm := NewRoomManager()
	client := NewWSClient(rm, "sid1", "m1", nil)
	if client == nil {
		t.Fatal("expected non-nil client")
	}
	if client.sid != "sid1" {
		t.Errorf("expected sid1, got %s", client.sid)
	}
	if client.matchID != "m1" {
		t.Errorf("expected m1, got %s", client.matchID)
	}
}

// ---------- handleMessage ----------

func TestHandleMessage_Join(t *testing.T) {
	rm := NewRoomManager()
	client := NewWSClient(rm, "sid1", "m1", nil)

	msg, _ := json.Marshal(map[string]interface{}{
		"type":       "join",
		"game_type":  "chess",
	})
	client.handleMessage(msg)

	if rm.GetRoom("m1") == nil {
		t.Error("room should be created after join message")
	}
}

func TestHandleMessage_Move(t *testing.T) {
	rm := NewRoomManager()
	client := NewWSClient(rm, "sid1", "m1", nil)

	movePayload, _ := json.Marshal(WSMove{From: "e2", To: "e4"})
	msg, _ := json.Marshal(map[string]interface{}{
		"type":    "move",
		"payload": movePayload,
	})
	client.handleMessage(msg)
	// Room doesn't exist, so move is routed to non-existent room
}

func TestHandleMessage_Resign(t *testing.T) {
	rm := NewRoomManager()
	client := NewWSClient(rm, "sid1", "m1", nil)

	msg, _ := json.Marshal(map[string]interface{}{
		"type": "resign",
	})
	client.handleMessage(msg)
}

func TestHandleMessage_DrawOffer(t *testing.T) {
	rm := NewRoomManager()
	client := NewWSClient(rm, "sid1", "m1", nil)

	msg, _ := json.Marshal(map[string]interface{}{
		"type": "draw_offer",
	})
	client.handleMessage(msg)
}

func TestHandleMessage_DrawAccept(t *testing.T) {
	rm := NewRoomManager()
	client := NewWSClient(rm, "sid1", "m1", nil)

	msg, _ := json.Marshal(map[string]interface{}{
		"type": "draw_accept",
	})
	client.handleMessage(msg)
}

func TestHandleMessage_DrawDecline(t *testing.T) {
	rm := NewRoomManager()
	client := NewWSClient(rm, "sid1", "m1", nil)

	msg, _ := json.Marshal(map[string]interface{}{
		"type": "draw_decline",
	})
	client.handleMessage(msg)
}

func TestHandleMessage_RollDice(t *testing.T) {
	rm := NewRoomManager()
	client := NewWSClient(rm, "sid1", "m1", nil)

	msg, _ := json.Marshal(map[string]interface{}{
		"type": "roll_dice",
	})
	client.handleMessage(msg)
}

func TestHandleMessage_InvalidJSON(t *testing.T) {
	rm := NewRoomManager()
	client := NewWSClient(rm, "sid1", "m1", nil)

	client.handleMessage([]byte("not json"))
	// Should not panic
}

func TestHandleMessage_UnknownType(t *testing.T) {
	rm := NewRoomManager()
	client := NewWSClient(rm, "sid1", "m1", nil)

	msg, _ := json.Marshal(map[string]interface{}{
		"type": "unknown_type",
	})
	client.handleMessage(msg)
	// Should not panic
}
