package websocket_test

import (
	"testing"
	"time"

	"game-platform/internal/websocket"
)

func TestHub_NewHub(t *testing.T) {
	hub := websocket.NewHub()

	if hub == nil {
		t.Fatal("Expected hub, got nil")
	}

	t.Log("Hub created successfully")
}

func TestHub_Run(t *testing.T) {
	hub := websocket.NewHub()
	
	done := make(chan bool)
	
	go func() {
		hub.Run()
		done <- true
	}()

	// Тест: регистрация клиента
	client := websocket.NewClient(hub, "emp_12345")

	hub.Register(client)

	// Даем время на обработку
	time.Sleep(100 * time.Millisecond)

	// Тест: отписка клиента
	hub.Unregister(client)

	// Даем время на обработку
	time.Sleep(100 * time.Millisecond)

	// Останавливаем хаб
	close(done)

	t.Log("Hub run test completed")
}

func TestHub_Broadcast(t *testing.T) {
	hub := websocket.NewHub()

	// Тест: широковещательная рассылка
	testMessage := []byte(`{"type":"game_state","board":[]}`)
	hub.Broadcast(testMessage)

	// Даем время на обработку
	time.Sleep(100 * time.Millisecond)

	t.Log("Broadcast message sent")
}

func TestClient_HandleMessage(t *testing.T) {
	hub := websocket.NewHub()
	
	client := websocket.NewClient(hub, "emp_12345")

	// Тест: сообщение join
	joinMsg := []byte(`{"type":"join","game_id":"g_chess_123"}`)
	client.HandleMessage(joinMsg)
	t.Log("Join message handled")

	// Тест: сообщение move
	moveMsg := []byte(`{"type":"move","from":"e2","to":"e4"}`)
	client.HandleMessage(moveMsg)
	t.Log("Move message handled")

	// Тест: сообщение game_over
	gameOverMsg := []byte(`{"type":"game_over","winner":"white"}`)
	client.HandleMessage(gameOverMsg)
	t.Log("Game over message handled")

	// Тест: некорректное сообщение
	badMsg := []byte(`{"invalid":}`)
	client.HandleMessage(badMsg) // Не должно паниковать
	t.Log("Invalid message handled without panic")
}

func TestClient_ReadPump(t *testing.T) {
	// Этот тест требует реального WebSocket соединения
	// Пропускаем в unit-тестах
	t.Skip("Skipping readPump test: requires WebSocket connection")
}

func TestClient_WritePump(t *testing.T) {
	// Этот тест требует реального WebSocket соединения
	// Пропускаем в unit-тестах
	t.Skip("Skipping writePump test: requires WebSocket connection")
}
