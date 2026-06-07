package websocket

import (
	"sync"
	"testing"
)

func TestNewRoomManager(t *testing.T) {
	rm := NewRoomManager()
	if rm == nil {
		t.Fatal("NewRoomManager returned nil")
	}
	if rm.RoomCount() != 0 {
		t.Errorf("expected 0 rooms, got %d", rm.RoomCount())
	}
}

func TestCreateRoom(t *testing.T) {
	rm := NewRoomManager()
	room := rm.CreateRoom("match-1", "chess")
	if room == nil {
		t.Fatal("CreateRoom returned nil")
	}
	if room.MatchID() != "match-1" {
		t.Errorf("expected match_id match-1, got %s", room.MatchID())
	}
	if room.GameType() != "chess" {
		t.Errorf("expected game_type chess, got %s", room.GameType())
	}
	if rm.RoomCount() != 1 {
		t.Errorf("expected 1 room, got %d", rm.RoomCount())
	}
}

func TestCreateRoomIdempotent(t *testing.T) {
	rm := NewRoomManager()
	room1 := rm.CreateRoom("match-1", "chess")
	room2 := rm.CreateRoom("match-1", "checkers")
	if room1 != room2 {
		t.Error("CreateRoom should return existing room for same match_id")
	}
	if room1.GameType() != "chess" {
		t.Error("room type should remain chess, not be overwritten")
	}
	if rm.RoomCount() != 1 {
		t.Errorf("expected 1 room, got %d", rm.RoomCount())
	}
}

func TestGetRoom(t *testing.T) {
	rm := NewRoomManager()

	// Non-existent room
	if rm.GetRoom("nonexistent") != nil {
		t.Error("GetRoom should return nil for non-existent room")
	}

	room := rm.CreateRoom("match-1", "chess")
	got := rm.GetRoom("match-1")
	if got != room {
		t.Error("GetRoom should return the created room")
	}
}

func TestRemoveRoom(t *testing.T) {
	rm := NewRoomManager()
	rm.CreateRoom("match-1", "chess")
	rm.CreateRoom("match-2", "checkers")

	if rm.RoomCount() != 2 {
		t.Errorf("expected 2 rooms, got %d", rm.RoomCount())
	}

	rm.RemoveRoom("match-1")
	if rm.RoomCount() != 1 {
		t.Errorf("expected 1 room after removal, got %d", rm.RoomCount())
	}
	if rm.GetRoom("match-1") != nil {
		t.Error("removed room should be nil")
	}

	// Remove non-existent (should not panic)
	rm.RemoveRoom("nonexistent")
	if rm.RoomCount() != 1 {
		t.Errorf("expected 1 room, got %d", rm.RoomCount())
	}
}

func TestConcurrentCreateRoom(t *testing.T) {
	rm := NewRoomManager()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			matchID := "match-concurrent"
			rm.CreateRoom(matchID, "chess")
		}(i)
	}

	wg.Wait()
	if rm.RoomCount() != 1 {
		t.Errorf("expected 1 room after concurrent creates, got %d", rm.RoomCount())
	}
}

func TestCreateRoomDifferentGameTypes(t *testing.T) {
	rm := NewRoomManager()

	rm.CreateRoom("chess-match", "chess")
	rm.CreateRoom("checkers-match", "checkers")
	rm.CreateRoom("backgammon-match", "backgammon")

	if rm.RoomCount() != 3 {
		t.Errorf("expected 3 rooms, got %d", rm.RoomCount())
	}

	if rm.GetRoom("chess-match").GameType() != "chess" {
		t.Error("chess room should have type chess")
	}
	if rm.GetRoom("checkers-match").GameType() != "checkers" {
		t.Error("checkers room should have type checkers")
	}
	if rm.GetRoom("backgammon-match").GameType() != "backgammon" {
		t.Error("backgammon room should have type backgammon")
	}
}
