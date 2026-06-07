package websocket

import (
	"encoding/json"
	"testing"
)

func TestNewGameRoom(t *testing.T) {
	room := NewGameRoom("match-1", "chess")
	if room.MatchID() != "match-1" {
		t.Errorf("expected match_id match-1, got %s", room.MatchID())
	}
	if room.GameType() != "chess" {
		t.Errorf("expected game_type chess, got %s", room.GameType())
	}
	if room.IsGameOver() {
		t.Error("new room should not be game over")
	}
	if room.IsFull() {
		t.Error("new room should not be full")
	}
}

func TestGameRoomChessEngine(t *testing.T) {
	room := NewGameRoom("match-1", "chess")

	// Check initial state
	state := room.BuildStatePayload()
	if state["turn"] != "white" {
		t.Errorf("chess should start with white's turn, got %v", state["turn"])
	}
	if state["board"] == "" {
		t.Error("board should not be empty")
	}
}

func TestGameRoomCheckersEngine(t *testing.T) {
	room := NewGameRoom("match-1", "checkers")

	state := room.BuildStatePayload()
	if state["turn"] != "white" {
		t.Errorf("checkers should start with white's turn, got %v", state["turn"])
	}
}

func TestGameRoomBackgammonEngine(t *testing.T) {
	room := NewGameRoom("match-1", "backgammon")

	state := room.BuildStatePayload()
	if state["turn"] != "white" {
		t.Errorf("backgammon should start with white's turn, got %v", state["turn"])
	}
}

func TestGameRoomUnknownType(t *testing.T) {
	// Unknown type should default to chess
	room := NewGameRoom("match-1", "unknown")
	if room.GameType() != "unknown" {
		t.Errorf("expected game_type unknown, got %s", room.GameType())
	}
	// Should still have a working engine
	state := room.BuildStatePayload()
	if state["board"] == "" {
		t.Error("even unknown type should have a board")
	}
}

func TestChessAdapterLegalMoves(t *testing.T) {
	adapter := NewChessAdapter()

	moves := adapter.GetLegalMoves()
	if len(moves) == 0 {
		t.Error("chess should have legal moves from starting position")
	}

	// Check that moves have from/to
	for _, m := range moves {
		if m.From == "" || m.To == "" {
			t.Errorf("move missing from/to: %+v", m)
		}
	}
}

func TestChessAdapterApplyMove(t *testing.T) {
	adapter := NewChessAdapter()

	// e2-e4 (king's pawn)
	err := adapter.ApplyMove(WSMove{From: "e2", To: "e4"})
	if err != nil {
		t.Fatalf("e2-e4 should be legal: %v", err)
	}

	if adapter.GetTurn() != "black" {
		t.Errorf("after e2-e4, turn should be black, got %s", adapter.GetTurn())
	}

	if !adapter.IsGameOver() {
		// Game should still be playing
	} else {
		t.Error("game should not be over after e2-e4")
	}
}

func TestChessAdapterIllegalMove(t *testing.T) {
	adapter := NewChessAdapter()

	// e2-e5 (pawn can't move 3 squares)
	err := adapter.ApplyMove(WSMove{From: "e2", To: "e5"})
	if err == nil {
		t.Error("e2-e5 should be illegal")
	}
}

func TestChessAdapterResign(t *testing.T) {
	adapter := NewChessAdapter()
	adapter.Resign("white")

	if !adapter.IsGameOver() {
		t.Error("game should be over after resign")
	}
	if adapter.GetGameOverReason() != "resign" {
		t.Errorf("expected resign reason, got %s", adapter.GetGameOverReason())
	}
	if adapter.GetWinner() != "black" {
		t.Errorf("after white resigns, black should win, got %s", adapter.GetWinner())
	}
}

func TestChessAdapterFEN(t *testing.T) {
	adapter := NewChessAdapter()
	fen := adapter.GetBoard()
	if fen == "" {
		t.Error("FEN should not be empty")
	}
	// Starting FEN should begin with "rnbqkbnr"
	if len(fen) < 8 || fen[:8] != "rnbqkbnr" {
		t.Errorf("starting FEN should begin with rnbqkbnr, got: %s", fen[:min(8, len(fen))])
	}
}

func TestCheckersAdapterLegalMoves(t *testing.T) {
	adapter := NewCheckersAdapter()

	moves := adapter.GetLegalMoves()
	if len(moves) == 0 {
		t.Error("checkers should have legal moves from starting position")
	}
}

func TestCheckersAdapterResign(t *testing.T) {
	adapter := NewCheckersAdapter()
	adapter.Resign("white")

	if !adapter.IsGameOver() {
		t.Error("game should be over after resign")
	}
	if adapter.GetWinner() != "black" {
		t.Errorf("after white resigns, black should win, got %s", adapter.GetWinner())
	}
}

func TestBackgammonAdapterRollDice(t *testing.T) {
	adapter := NewBackgammonAdapter()

	// Before rolling, no legal moves (state is Rolling)
	moves := adapter.GetLegalMoves()
	if len(moves) != 0 {
		t.Error("backgammon should have no moves before rolling dice")
	}

	dice := adapter.RollDice()
	if len(dice) != 2 {
		t.Fatalf("expected 2 dice, got %d", len(dice))
	}
	if dice[0] < 1 || dice[0] > 6 || dice[1] < 1 || dice[1] > 6 {
		t.Errorf("dice values out of range: %v", dice)
	}
}

func TestBackgammonAdapterResign(t *testing.T) {
	adapter := NewBackgammonAdapter()
	adapter.Resign("black")

	if !adapter.IsGameOver() {
		t.Error("game should be over after resign")
	}
	if adapter.GetWinner() != "white" {
		t.Errorf("after black resigns, white should win, got %s", adapter.GetWinner())
	}
}

func TestNewDefaultClock(t *testing.T) {
	clock := NewDefaultClock()

	if clock.WhiteMs() != DefaultInitialTime {
		t.Errorf("expected white time %d, got %d", DefaultInitialTime, clock.WhiteMs())
	}
	if clock.BlackMs() != DefaultInitialTime {
		t.Errorf("expected black time %d, got %d", DefaultInitialTime, clock.BlackMs())
	}
	if clock.IsGameOver() {
		t.Error("new clock should not be game over")
	}
}

func TestClockTimeout(t *testing.T) {
	clock := NewClock(1000, 0) // 1 second, no increment

	// Tick away white's time
	timedOut := clock.TickForTest("white", 1000)
	if !timedOut {
		t.Error("clock should timeout after 1000ms tick")
	}
	if !clock.IsGameOver() {
		t.Error("clock should be game over")
	}
	if clock.TimeoutColor() != "white" {
		t.Errorf("timeout color should be white, got %s", clock.TimeoutColor())
	}
	if clock.WhiteMs() != 0 {
		t.Errorf("white ms should be 0, got %d", clock.WhiteMs())
	}
	// Black's time should be untouched
	if clock.BlackMs() != 1000 {
		t.Errorf("black ms should be 1000, got %d", clock.BlackMs())
	}
}

func TestClockSwitchPlayer(t *testing.T) {
	clock := NewClock(10000, 100) // 10 seconds, 100ms increment

	// Simulate white making a move: switch from white to black
	// First, we need to set lastTick
	clock.Start("white")
	clock.SwitchPlayer("white", 100)

	// White should have gotten an increment
	if clock.WhiteMs() <= 10000-100 {
		t.Errorf("white should have gotten increment, got %d", clock.WhiteMs())
	}
}

func TestClockSetTime(t *testing.T) {
	clock := NewDefaultClock()
	clock.SetTime("white", 5000)

	if clock.WhiteMs() != 5000 {
		t.Errorf("expected 5000, got %d", clock.WhiteMs())
	}
	if clock.BlackMs() != DefaultInitialTime {
		t.Errorf("black time should be unchanged")
	}
}

func TestGameStatePayloadSerialization(t *testing.T) {
	room := NewGameRoom("match-1", "chess")
	state := room.BuildStatePayload()

	data, err := json.Marshal(state)
	if err != nil {
		t.Fatalf("failed to marshal state: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to unmarshal state: %v", err)
	}

	if parsed["type"] != "state" {
		t.Errorf("expected type state, got %v", parsed["type"])
	}
	if parsed["match_id"] != "match-1" {
		t.Errorf("expected match_id match-1, got %v", parsed["match_id"])
	}
	if parsed["turn"] != "white" {
		t.Errorf("expected turn white, got %v", parsed["turn"])
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
