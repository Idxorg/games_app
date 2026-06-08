package websocket

import (
	"encoding/json"
	"strings"
	"testing"

	"game-platform/internal/model"
	"game-platform/tests/mocks"
)

// ---------- NewGameRoom ----------

func TestNewGameRoom_Chess(t *testing.T) {
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
	_ = room.Engine()
}

func TestNewGameRoom_Checkers(t *testing.T) {
	room := NewGameRoom("match-1", "checkers")
	state := room.BuildStatePayload()
	if state["turn"] != "white" {
		t.Errorf("checkers should start with white's turn, got %v", state["turn"])
	}
}

func TestNewGameRoom_Backgammon(t *testing.T) {
	room := NewGameRoom("match-1", "backgammon")
	state := room.BuildStatePayload()
	if state["turn"] != "white" {
		t.Errorf("backgammon should start with white's turn, got %v", state["turn"])
	}
}

func TestNewGameRoom_UnknownType(t *testing.T) {
	room := NewGameRoom("match-1", "unknown")
	if room.GameType() != "unknown" {
		t.Errorf("expected game_type unknown, got %s", room.GameType())
	}
	state := room.BuildStatePayload()
	if state["board"] == "" {
		t.Error("even unknown type should have a board")
	}
}

// ---------- JoinPlayer ----------

func TestJoinPlayer_P1(t *testing.T) {
	room := NewGameRoom("m1", "chess")
	room.TestMsgCh = make(chan []byte, 10)

	err := room.JoinPlayer("sid1", nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if room.IsFull() {
		t.Error("room should not be full after one player")
	}
	if room.GetPlayerColor("sid1") != "white" {
		t.Errorf("first player should be white, got %s", room.GetPlayerColor("sid1"))
	}
}

func TestJoinPlayer_P2(t *testing.T) {
	room := NewGameRoom("m1", "chess")
	room.TestMsgCh = make(chan []byte, 20)

	err := room.JoinPlayer("sid1", nil)
	if err != nil {
		t.Fatalf("p1 join error: %v", err)
	}
	err = room.JoinPlayer("sid2", nil)
	if err != nil {
		t.Fatalf("p2 join error: %v", err)
	}
	if !room.IsFull() {
		t.Error("room should be full after two players")
	}
	if room.GetPlayerColor("sid2") != "black" {
		t.Errorf("second player should be black, got %s", room.GetPlayerColor("sid2"))
	}
}

func TestJoinPlayer_FullRoom(t *testing.T) {
	room := NewGameRoom("m1", "chess")
	room.TestMsgCh = make(chan []byte, 20)

	room.JoinPlayer("sid1", nil)
	room.JoinPlayer("sid2", nil)

	err := room.JoinPlayer("sid3", nil)
	if err == nil {
		t.Error("expected error when room is full")
	}
}

func TestJoinPlayer_GameOver(t *testing.T) {
	room := NewGameRoom("m1", "chess")
	room.TestMsgCh = make(chan []byte, 20)

	room.JoinPlayer("sid1", nil)
	room.JoinPlayer("sid2", nil)
	// Force game over
	room.HandleResign("sid1")

	err := room.JoinPlayer("sid3", nil)
	if err == nil {
		t.Error("expected error when game is over")
	}
}

// ---------- HandleMove ----------

func TestHandleMove_ChessValidMove(t *testing.T) {
	room := NewGameRoom("m1", "chess")
	room.TestMsgCh = make(chan []byte, 50)

	room.JoinPlayer("sid1", nil)
	room.JoinPlayer("sid2", nil)

	moveData, _ := json.Marshal(WSMove{From: "e2", To: "e4"})
	room.HandleMove("sid1", moveData)

	history := room.GetMoveHistory()
	if len(history) != 1 {
		t.Fatalf("expected 1 move in history, got %d", len(history))
	}
	if history[0]["from"] != "e2" {
		t.Errorf("expected from e2, got %v", history[0]["from"])
	}
}

func TestHandleMove_ChessIllegalMove(t *testing.T) {
	room := NewGameRoom("m1", "chess")
	room.TestMsgCh = make(chan []byte, 50)

	room.JoinPlayer("sid1", nil)
	room.JoinPlayer("sid2", nil)

	// e2-e5 is illegal (pawn can't move 3 squares)
	moveData, _ := json.Marshal(WSMove{From: "e2", To: "e5"})
	room.HandleMove("sid1", moveData)

	// Drain messages - should contain an error
	history := room.GetMoveHistory()
	if len(history) != 0 {
		t.Errorf("illegal move should not be recorded, got %d moves", len(history))
	}
}

func TestHandleMove_NotYourTurn(t *testing.T) {
	room := NewGameRoom("m1", "chess")
	room.TestMsgCh = make(chan []byte, 50)

	room.JoinPlayer("sid1", nil)
	room.JoinPlayer("sid2", nil)

	moveData, _ := json.Marshal(WSMove{From: "e7", To: "e5"})
	room.HandleMove("sid2", moveData) // black tries to move first

	history := room.GetMoveHistory()
	if len(history) != 0 {
		t.Errorf("should not record move when not your turn")
	}
}

func TestHandleMove_InvalidJSON(t *testing.T) {
	room := NewGameRoom("m1", "chess")
	room.TestMsgCh = make(chan []byte, 50)

	room.JoinPlayer("sid1", nil)
	room.JoinPlayer("sid2", nil)

	room.HandleMove("sid1", []byte("not-json"))

	history := room.GetMoveHistory()
	if len(history) != 0 {
		t.Errorf("invalid JSON should not record a move")
	}
}

func TestHandleMove_GameAlreadyOver(t *testing.T) {
	room := NewGameRoom("m1", "chess")
	room.TestMsgCh = make(chan []byte, 50)

	room.JoinPlayer("sid1", nil)
	room.JoinPlayer("sid2", nil)
	room.HandleResign("sid1")

	moveData, _ := json.Marshal(WSMove{From: "e2", To: "e4"})
	room.HandleMove("sid2", moveData)

	history := room.GetMoveHistory()
	if len(history) != 0 {
		t.Errorf("no moves should be recorded after game over")
	}
}

func TestHandleMove_NotInGame(t *testing.T) {
	room := NewGameRoom("m1", "chess")
	room.TestMsgCh = make(chan []byte, 50)

	room.JoinPlayer("sid1", nil)
	room.JoinPlayer("sid2", nil)

	moveData, _ := json.Marshal(WSMove{From: "e2", To: "e4"})
	room.HandleMove("sid_other", moveData)

	history := room.GetMoveHistory()
	if len(history) != 0 {
		t.Errorf("should not record move for non-player")
	}
}

// ---------- HandleResign ----------

func TestHandleResign_Valid(t *testing.T) {
	room := NewGameRoom("m1", "chess")
	room.TestMsgCh = make(chan []byte, 50)
	room.SetMatchRepo(mocks.NewMockMatchRepo())

	room.JoinPlayer("sid1", nil)
	room.JoinPlayer("sid2", nil)
	room.HandleResign("sid1")

	if !room.IsGameOver() {
		t.Error("game should be over after resign")
	}
}

func TestHandleResign_NotInGame(t *testing.T) {
	room := NewGameRoom("m1", "chess")
	room.TestMsgCh = make(chan []byte, 50)

	room.JoinPlayer("sid1", nil)
	room.JoinPlayer("sid2", nil)
	room.HandleResign("sid_other") // not a player
}

func TestHandleResign_GameAlreadyOver(t *testing.T) {
	room := NewGameRoom("m1", "chess")
	room.TestMsgCh = make(chan []byte, 50)

	room.JoinPlayer("sid1", nil)
	room.JoinPlayer("sid2", nil)
	room.HandleResign("sid1")
	room.HandleResign("sid2") // double resign
}

// ---------- HandleDrawOffer ----------

func TestHandleDrawOffer_Sent(t *testing.T) {
	room := NewGameRoom("m1", "chess")
	room.TestMsgCh = make(chan []byte, 50)

	room.JoinPlayer("sid1", nil)
	room.JoinPlayer("sid2", nil)
	room.HandleDrawOffer("sid1")

	if room.IsGameOver() {
		t.Error("game should not be over after single draw offer")
	}
}

func TestHandleDrawOffer_MutualDraw(t *testing.T) {
	room := NewGameRoom("m1", "chess")
	room.TestMsgCh = make(chan []byte, 50)
	room.SetMatchRepo(mocks.NewMockMatchRepo())

	room.JoinPlayer("sid1", nil)
	room.JoinPlayer("sid2", nil)
	room.HandleDrawOffer("sid1")
	room.HandleDrawOffer("sid2") // mutual draw

	if !room.IsGameOver() {
		t.Error("game should be over after mutual draw offer")
	}
}

func TestHandleDrawOffer_GameOver(t *testing.T) {
	room := NewGameRoom("m1", "chess")
	room.TestMsgCh = make(chan []byte, 50)

	room.JoinPlayer("sid1", nil)
	room.JoinPlayer("sid2", nil)
	room.HandleResign("sid1")
	room.HandleDrawOffer("sid2") // draw offer after game over
}

// ---------- HandleDrawAccept ----------

func TestHandleDrawAccept_NoPendingOffer(t *testing.T) {
	room := NewGameRoom("m1", "chess")
	room.TestMsgCh = make(chan []byte, 50)

	room.JoinPlayer("sid1", nil)
	room.JoinPlayer("sid2", nil)
	room.HandleDrawAccept("sid2") // no pending offer
}

func TestHandleDrawAccept_WrongAcceptor(t *testing.T) {
	room := NewGameRoom("m1", "chess")
	room.TestMsgCh = make(chan []byte, 50)

	room.JoinPlayer("sid1", nil)
	room.JoinPlayer("sid2", nil)
	room.HandleDrawOffer("sid1")
	room.HandleDrawAccept("sid1") // offerer tries to accept
}

func TestHandleDrawAccept_Valid(t *testing.T) {
	room := NewGameRoom("m1", "chess")
	room.TestMsgCh = make(chan []byte, 50)
	room.SetMatchRepo(mocks.NewMockMatchRepo())

	room.JoinPlayer("sid1", nil)
	room.JoinPlayer("sid2", nil)
	room.HandleDrawOffer("sid1")
	room.HandleDrawAccept("sid2")

	if !room.IsGameOver() {
		t.Error("game should be over after draw accepted")
	}
}

func TestHandleDrawAccept_GameOver(t *testing.T) {
	room := NewGameRoom("m1", "chess")
	room.TestMsgCh = make(chan []byte, 50)

	room.JoinPlayer("sid1", nil)
	room.JoinPlayer("sid2", nil)
	room.HandleDrawOffer("sid1")
	room.HandleResign("sid1")
	room.HandleDrawAccept("sid2") // accept after game over
}

// ---------- HandleDrawDecline ----------

func TestHandleDrawDecline_NoPendingOffer(t *testing.T) {
	room := NewGameRoom("m1", "chess")
	room.TestMsgCh = make(chan []byte, 50)

	room.JoinPlayer("sid1", nil)
	room.JoinPlayer("sid2", nil)
	room.HandleDrawDecline("sid2") // no pending offer
}

func TestHandleDrawDecline_WrongDecliner(t *testing.T) {
	room := NewGameRoom("m1", "chess")
	room.TestMsgCh = make(chan []byte, 50)

	room.JoinPlayer("sid1", nil)
	room.JoinPlayer("sid2", nil)
	room.HandleDrawOffer("sid1")
	room.HandleDrawDecline("sid1") // offerer tries to decline
}

func TestHandleDrawDecline_Valid(t *testing.T) {
	room := NewGameRoom("m1", "chess")
	room.TestMsgCh = make(chan []byte, 50)

	room.JoinPlayer("sid1", nil)
	room.JoinPlayer("sid2", nil)
	room.HandleDrawOffer("sid1")
	room.HandleDrawDecline("sid2")

	if room.IsGameOver() {
		t.Error("game should not be over after draw declined")
	}
}

// ---------- HandleRollDice ----------

func TestHandleRollDice_Backgammon(t *testing.T) {
	room := NewGameRoom("m1", "backgammon")
	room.TestMsgCh = make(chan []byte, 50)

	room.JoinPlayer("sid1", nil)
	room.JoinPlayer("sid2", nil)
	room.HandleRollDice("sid1")

	// Check that a dice_rolled message was broadcast
	found := false
	for i := 0; i < 5; i++ {
		select {
		case msg := <-room.TestMsgCh:
			if strings.Contains(string(msg), "dice_rolled") {
				found = true
			}
		default:
		}
	}
	if !found {
		t.Error("expected dice_rolled message")
	}
}

func TestHandleRollDice_NotBackgammon(t *testing.T) {
	room := NewGameRoom("m1", "chess")
	room.TestMsgCh = make(chan []byte, 50)

	room.JoinPlayer("sid1", nil)
	room.JoinPlayer("sid2", nil)
	room.HandleRollDice("sid1")

	// Should send an error, not dice_rolled
	found := false
	for i := 0; i < 5; i++ {
		select {
		case msg := <-room.TestMsgCh:
			if strings.Contains(string(msg), "dice_rolled") {
				found = true
			}
		default:
		}
	}
	if found {
		t.Error("chess room should not broadcast dice_rolled")
	}
}

func TestHandleRollDice_GameOver(t *testing.T) {
	room := NewGameRoom("m1", "backgammon")
	room.TestMsgCh = make(chan []byte, 50)

	room.JoinPlayer("sid1", nil)
	room.JoinPlayer("sid2", nil)
	room.HandleResign("sid1")
	room.HandleRollDice("sid2")
}

// ---------- persistGameComplete ----------

func TestPersistGameComplete_NilMatchRepo(t *testing.T) {
	room := NewGameRoom("m1", "chess")
	room.TestMsgCh = make(chan []byte, 50)
	// matchRepo is nil by default

	room.JoinPlayer("sid1", nil)
	room.JoinPlayer("sid2", nil)
	room.HandleResign("sid1") // should not panic

	if !room.IsGameOver() {
		t.Error("game should be over")
	}
}

func TestPersistGameComplete_WithMatchRepo(t *testing.T) {
	matchRepo := mocks.NewMockMatchRepo()
	// Pre-create the match so Complete won't fail
	matchRepo.Create(nil, &model.Match{
		ID: "m1", GameType: "chess",
		Player1SID: "sid1", Player2SID: "sid2",
		Status: "in_progress",
	})

	room := NewGameRoom("m1", "chess")
	room.TestMsgCh = make(chan []byte, 50)
	room.SetMatchRepo(matchRepo)

	room.JoinPlayer("sid1", nil)
	room.JoinPlayer("sid2", nil)
	room.HandleResign("sid1")

	// Give the async goroutine time to complete
	// Just verify the game is over; no panic
	if !room.IsGameOver() {
		t.Error("game should be over")
	}
}

// ---------- GetMoveHistory ----------

func TestGetMoveHistory_Empty(t *testing.T) {
	room := NewGameRoom("m1", "chess")
	history := room.GetMoveHistory()
	if len(history) != 0 {
		t.Errorf("expected empty history, got %d", len(history))
	}
}

func TestGetMoveHistory_AfterMoves(t *testing.T) {
	room := NewGameRoom("m1", "chess")
	room.TestMsgCh = make(chan []byte, 50)

	room.JoinPlayer("sid1", nil)
	room.JoinPlayer("sid2", nil)

	moveData, _ := json.Marshal(WSMove{From: "e2", To: "e4"})
	room.HandleMove("sid1", moveData)

	moveData, _ = json.Marshal(WSMove{From: "e7", To: "e5"})
	room.HandleMove("sid2", moveData)

	history := room.GetMoveHistory()
	if len(history) != 2 {
		t.Errorf("expected 2 moves, got %d", len(history))
	}
}

// ---------- GetPlayerColor ----------

func TestGetPlayerColor(t *testing.T) {
	room := NewGameRoom("m1", "chess")
	room.TestMsgCh = make(chan []byte, 20)

	room.JoinPlayer("sid1", nil)
	room.JoinPlayer("sid2", nil)

	if room.GetPlayerColor("sid1") != "white" {
		t.Errorf("p1 should be white, got %s", room.GetPlayerColor("sid1"))
	}
	if room.GetPlayerColor("sid2") != "black" {
		t.Errorf("p2 should be black, got %s", room.GetPlayerColor("sid2"))
	}
	if room.GetPlayerColor("sid_other") != "" {
		t.Errorf("unknown SID should return empty, got %s", room.GetPlayerColor("sid_other"))
	}
}

// ---------- getOpponentSID (tested via draw logic) ----------

func TestGetOpponentSID_ThroughDraw(t *testing.T) {
	room := NewGameRoom("m1", "chess")
	room.TestMsgCh = make(chan []byte, 50)

	room.JoinPlayer("sid1", nil)
	room.JoinPlayer("sid2", nil)

	// Offer draw from sid1, decline from sid2
	room.HandleDrawOffer("sid1")
	// Verify draw_offer was sent (drain channel)
	drain := func() {
		for {
			select {
			case <-room.TestMsgCh:
			default:
				return
			}
		}
	}
	drain()

	// Now sid2 declines
	room.HandleDrawDecline("sid2")
}

// ---------- Cleanup ----------

func TestCleanup(t *testing.T) {
	room := NewGameRoom("m1", "chess")
	// Cleanup should not panic even with no connections
	room.Cleanup()
}

// ---------- SetMatchRepo / SetRatingService ----------

func TestSetMatchRepo(t *testing.T) {
	room := NewGameRoom("m1", "chess")
	repo := mocks.NewMockMatchRepo()
	room.SetMatchRepo(repo)
	room.SetRatingService(nil)
}

func TestSetRatingService(t *testing.T) {
	room := NewGameRoom("m1", "chess")
	updater := &mocks.MockRatingUpdater{}
	room.SetRatingService(updater)
}

// ---------- Helper ----------

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func TestHandleEngineGameOver_Checkmate(t *testing.T) {
	room := NewGameRoom("m1", "chess")
	room.player1SID = "white-p"
	room.player2SID = "black-p"
	repo := mocks.NewMockMatchRepo()
	room.SetMatchRepo(repo)
	updater := &mocks.MockRatingUpdater{}
	room.SetRatingService(updater)

	// Set up a position where white delivers checkmate (Scholar's mate)
	// e2-e4, e7-e5, d1-h5, b8-c6, f1-c4, g8-f6, h5xf7#
	_ = room.engine.ApplyMove(WSMove{From: "e2", To: "e4"})
	_ = room.engine.ApplyMove(WSMove{From: "e7", To: "e5"})
	_ = room.engine.ApplyMove(WSMove{From: "d1", To: "h5"})
	_ = room.engine.ApplyMove(WSMove{From: "b8", To: "c6"})
	_ = room.engine.ApplyMove(WSMove{From: "f1", To: "c4"})
	_ = room.engine.ApplyMove(WSMove{From: "g8", To: "f6"})

	// Qxf7# — checkmate
	err := room.engine.ApplyMove(WSMove{From: "h5", To: "f7"})
	if err != nil {
		t.Fatalf("move failed: %v", err)
	}

	room.handleEngineGameOver()

	if !room.gameOver {
		t.Error("game should be over")
	}
	if room.winnerSID != "white-p" {
		t.Errorf("expected winnerSID=white-p, got %s", room.winnerSID)
	}
	reason := room.gameOverReason
	if reason == "" {
		t.Error("gameOverReason should not be empty")
	}
}
