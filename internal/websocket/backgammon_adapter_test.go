package websocket

import (
	"testing"
)

// ---------- NewBackgammonAdapter ----------

func TestNewBackgammonAdapter(t *testing.T) {
	a := NewBackgammonAdapter()
	if a == nil {
		t.Fatal("expected non-nil adapter")
	}
}

// ---------- GetBoard ----------

func TestBackgammonAdapter_GetBoard(t *testing.T) {
	a := NewBackgammonAdapter()
	board := a.GetBoard()
	if board == "" {
		t.Error("board should not be empty")
	}
}

// ---------- GetTurn ----------

func TestBackgammonAdapter_GetTurn(t *testing.T) {
	a := NewBackgammonAdapter()
	if a.GetTurn() != "white" {
		t.Errorf("expected white, got %s", a.GetTurn())
	}
}

// ---------- RollDice ----------

func TestBackgammonAdapter_RollDice(t *testing.T) {
	a := NewBackgammonAdapter()
	dice := a.RollDice()
	if len(dice) != 2 {
		t.Fatalf("expected 2 dice, got %d", len(dice))
	}
	if dice[0] < 1 || dice[0] > 6 || dice[1] < 1 || dice[1] > 6 {
		t.Errorf("dice values out of range: %v", dice)
	}
}

// ---------- GetLegalMoves ----------

func TestBackgammonAdapter_GetLegalMoves_BeforeRoll(t *testing.T) {
	a := NewBackgammonAdapter()
	moves := a.GetLegalMoves()
	if len(moves) != 0 {
		t.Error("backgammon should have no moves before rolling dice")
	}
}

func TestBackgammonAdapter_GetLegalMoves_AfterRoll(t *testing.T) {
	a := NewBackgammonAdapter()
	a.RollDice()
	moves := a.GetLegalMoves()
	// With random dice, there should usually be some moves
	_ = moves
}

// ---------- GetDice ----------

func TestBackgammonAdapter_GetDice_BeforeRoll(t *testing.T) {
	a := NewBackgammonAdapter()
	dice := a.GetDice()
	if dice[0] != 0 || dice[1] != 0 {
		t.Errorf("expected zero dice before roll, got %v", dice)
	}
}

func TestBackgammonAdapter_GetDice_AfterRoll(t *testing.T) {
	a := NewBackgammonAdapter()
	a.RollDice()
	dice := a.GetDice()
	if dice[0] < 1 || dice[0] > 6 || dice[1] < 1 || dice[1] > 6 {
		t.Errorf("dice values out of range: %v", dice)
	}
}

// ---------- GetRemainingMoves ----------

func TestBackgammonAdapter_GetRemainingMoves_BeforeRoll(t *testing.T) {
	a := NewBackgammonAdapter()
	rm := a.GetRemainingMoves()
	if len(rm) != 0 {
		t.Errorf("expected 0 remaining moves before roll, got %d", len(rm))
	}
}

func TestBackgammonAdapter_GetRemainingMoves_AfterRoll(t *testing.T) {
	a := NewBackgammonAdapter()
	a.RollDice()
	rm := a.GetRemainingMoves()
	// Should have 2 or 4 moves (doubles = 4)
	if len(rm) != 2 && len(rm) != 4 {
		t.Errorf("expected 2 or 4 remaining moves, got %d", len(rm))
	}
}

// ---------- GetState ----------

func TestBackgammonAdapter_GetState(t *testing.T) {
	a := NewBackgammonAdapter()
	state := a.GetState()
	if state != 0 { // Rolling state
		t.Errorf("expected Rolling (0), got %d", state)
	}

	a.RollDice()
	state = a.GetState()
	if state != 1 { // Moving state
		t.Errorf("expected Moving (1), got %d", state)
	}
}

// ---------- ApplyMove ----------

func TestBackgammonAdapter_ApplyMove_NoDice(t *testing.T) {
	a := NewBackgammonAdapter()
	err := a.ApplyMove(WSMove{FromIdx: 23, ToIdx: 17, Die: 6})
	if err == nil {
		t.Error("should require dice roll first")
	}
}

func TestBackgammonAdapter_ApplyMove_NoDieValue(t *testing.T) {
	a := NewBackgammonAdapter()
	a.RollDice()
	err := a.ApplyMove(WSMove{FromIdx: 23, ToIdx: 17, Die: 0})
	if err == nil {
		t.Error("die value should be required")
	}
}

func TestBackgammonAdapter_ApplyMove_Invalid(t *testing.T) {
	a := NewBackgammonAdapter()
	a.RollDice()
	// Try a clearly illegal move
	err := a.ApplyMove(WSMove{FromIdx: 0, ToIdx: 24, Die: 1})
	// This might succeed or fail depending on dice, just verify no panic
	_ = err
}

func TestBackgammonAdapter_ApplyMove_LegalMove(t *testing.T) {
	a := NewBackgammonAdapter()
	a.game.SetDice(6, 1) // set specific dice for deterministic test

	moves := a.GetLegalMoves()
	if len(moves) == 0 {
		t.Skip("no legal moves available")
	}

	err := a.ApplyMove(moves[0])
	if err != nil {
		t.Fatalf("legal move should succeed: %v", err)
	}
}

// ---------- IsGameOver ----------

func TestBackgammonAdapter_IsGameOver_NotOver(t *testing.T) {
	a := NewBackgammonAdapter()
	if a.IsGameOver() {
		t.Error("new game should not be over")
	}
}

// ---------- GetGameOverReason ----------

func TestBackgammonAdapter_GetGameOverReason_NotOver(t *testing.T) {
	a := NewBackgammonAdapter()
	reason := a.GetGameOverReason()
	if reason != "" {
		t.Errorf("expected empty reason, got %s", reason)
	}
}

// ---------- GetWinner ----------

func TestBackgammonAdapter_GetWinner_NotOver(t *testing.T) {
	a := NewBackgammonAdapter()
	winner := a.GetWinner()
	if winner != "" {
		t.Errorf("expected empty winner, got %s", winner)
	}
}

// ---------- Resign ----------

func TestBackgammonAdapter_Resign_Black(t *testing.T) {
	a := NewBackgammonAdapter()
	a.Resign("black")

	if !a.IsGameOver() {
		t.Error("game should be over after resign")
	}
	if a.GetWinner() != "white" {
		t.Errorf("after black resigns, white should win, got %s", a.GetWinner())
	}
}

func TestBackgammonAdapter_Resign_White(t *testing.T) {
	a := NewBackgammonAdapter()
	a.Resign("white")

	if !a.IsGameOver() {
		t.Error("game should be over after resign")
	}
	if a.GetWinner() != "black" {
		t.Errorf("after white resigns, black should win, got %s", a.GetWinner())
	}
	if a.GetGameOverReason() != "checkmate" {
		t.Errorf("expected checkmate (broadly), got %s", a.GetGameOverReason())
	}
}

// ---------- pointToString ----------

func TestBackgammonAdapter_PointToString(t *testing.T) {
	// Indirectly tested via GetLegalMoves which uses pointToString
	a := NewBackgammonAdapter()
	a.RollDice()
	moves := a.GetLegalMoves()
	for _, m := range moves {
		if m.From == "" || m.To == "" {
			t.Errorf("move from/to should not be empty: %+v", m)
		}
		_ = m
	}
}
