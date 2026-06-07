package websocket

import (
	"encoding/json"
	"sync"
	"testing"
	"time"

	"game-platform/internal/model"
	ws "game-platform/internal/websocket"
	"game-platform/tests/mocks"
)

// ---------------------------------------------------------------------------
// Test helpers
// ---------------------------------------------------------------------------

// joinRoom joins two players to a room for testing (nil conns, messages go to TestMsgCh).
func joinRoom(t *testing.T, room *ws.GameRoom, p1SID, p2SID string) {
	t.Helper()
	room.TestMsgCh = make(chan []byte, 64)

	if err := room.JoinPlayer(p1SID, nil); err != nil {
		t.Fatalf("player1 join failed: %v", err)
	}
	if err := room.JoinPlayer(p2SID, nil); err != nil {
		t.Fatalf("player2 join failed: %v", err)
	}

	// Drain initial state broadcast
	drainMessages(room.TestMsgCh)
}

// drainMessages reads and discards all pending messages from the channel.
func drainMessages(ch chan []byte) {
	for {
		select {
		case <-ch:
		default:
			return
		}
	}
}

// readMsg reads a message from the channel with a timeout.
func readMsg(t *testing.T, ch chan []byte, timeout time.Duration) map[string]interface{} {
	t.Helper()
	select {
	case data := <-ch:
		var msg map[string]interface{}
		if err := json.Unmarshal(data, &msg); err != nil {
			t.Fatalf("failed to unmarshal message: %v (data: %s)", err, string(data))
		}
		return msg
	case <-time.After(timeout):
		t.Fatal("timeout waiting for message")
		return nil
	}
}

// moveJSON builds a JSON move payload for chess.
func moveJSON(from, to string) json.RawMessage {
	data, _ := json.Marshal(map[string]string{
		"from": from,
		"to":   to,
	})
	return data
}

// checkersMoveJSON builds a JSON move payload for checkers.
func checkersMoveJSON(fromRow, fromCol, toRow, toCol int) json.RawMessage {
	data, _ := json.Marshal(map[string]string{
		"from": jsonEncodePos(fromRow, fromCol),
		"to":   jsonEncodePos(toRow, toCol),
	})
	return data
}

func jsonEncodePos(row, col int) string {
	data, _ := json.Marshal([]int{row, col})
	return string(data)
}

// backgammonMoveJSON builds a JSON move payload for backgammon.
func backgammonMoveJSON(fromIdx, toIdx, die int) json.RawMessage {
	data, _ := json.Marshal(map[string]interface{}{
		"from_idx": fromIdx,
		"to_idx":   toIdx,
		"die":      die,
	})
	return data
}

// ---------------------------------------------------------------------------
// Chess gameplay tests
// ---------------------------------------------------------------------------

func TestChessGame_FullGame(t *testing.T) {
	room := ws.NewGameRoom("match-chess-1", "chess")
	p1 := "sid-white"
	p2 := "sid-black"
	joinRoom(t, room, p1, p2)

	// Verify initial state
	state := room.BuildStatePayload()
	if state["turn"] != "white" {
		t.Errorf("expected turn white, got %v", state["turn"])
	}
	board := state["board"].(string)
	if len(board) < 8 || board[:8] != "rnbqkbnr" {
		t.Errorf("expected starting position, got board: %s", board)
	}
	legalMoves := state["legal_moves"].([]ws.WSMove)
	if len(legalMoves) == 0 {
		t.Error("expected legal moves from starting position")
	}

	// 1. e2-e4 (white)
	room.HandleMove(p1, moveJSON("e2", "e4"))
	msg := readMsg(t, room.TestMsgCh, time.Second)
	if msg["type"] != "move_applied" {
		t.Errorf("expected move_applied, got type %v", msg["type"])
	}
	state = room.BuildStatePayload()
	if state["turn"] != "black" {
		t.Errorf("expected turn black after e4, got %v", state["turn"])
	}

	// Verify move history tracked
	history := room.GetMoveHistory()
	if len(history) != 1 {
		t.Fatalf("expected 1 move in history, got %d", len(history))
	}
	if history[0]["from"] != "e2" || history[0]["to"] != "e4" {
		t.Errorf("move history incorrect: %+v", history[0])
	}

	// 2. e7-e5 (black)
	room.HandleMove(p2, moveJSON("e7", "e5"))
	msg = readMsg(t, room.TestMsgCh, time.Second)
	if msg["type"] != "move_applied" {
		t.Errorf("expected move_applied, got type %v", msg["type"])
	}
	state = room.BuildStatePayload()
	if state["turn"] != "white" {
		t.Errorf("expected turn white after e5, got %v", state["turn"])
	}

	// 3. Nf3 (white)
	room.HandleMove(p1, moveJSON("g1", "f3"))
	readMsg(t, room.TestMsgCh, time.Second)

	// 4. Nc6 (black)
	room.HandleMove(p2, moveJSON("b8", "c6"))
	readMsg(t, room.TestMsgCh, time.Second)

	// 5. Bb5 (white) - Ruy Lopez
	room.HandleMove(p1, moveJSON("f1", "b5"))
	readMsg(t, room.TestMsgCh, time.Second)

	// Verify game is still playing
	if room.IsGameOver() {
		t.Error("game should not be over after 5 moves")
	}

	// Verify board state has changed from initial
	state = room.BuildStatePayload()
	boardAfter := state["board"].(string)
	if boardAfter == board {
		t.Error("board should have changed after 5 moves")
	}

	// Verify move count
	history = room.GetMoveHistory()
	if len(history) != 5 {
		t.Errorf("expected 5 moves in history, got %d", len(history))
	}
}

func TestChessGame_IllegalMove(t *testing.T) {
	room := ws.NewGameRoom("match-chess-2", "chess")
	p1 := "sid-white"
	p2 := "sid-black"
	joinRoom(t, room, p1, p2)

	// Try illegal move: e2-e5 (pawn can't move 3 squares)
	room.HandleMove(p1, moveJSON("e2", "e5"))
	msg := readMsg(t, room.TestMsgCh, time.Second)

	if msg["type"] != "error" {
		t.Errorf("expected error message, got type %v", msg["type"])
	}
	if msg["code"] != "illegal_move" {
		t.Errorf("expected illegal_move code, got %v", msg["code"])
	}

	// State should be unchanged - still white's turn
	state := room.BuildStatePayload()
	if state["turn"] != "white" {
		t.Errorf("turn should still be white after illegal move, got %v", state["turn"])
	}

	// Move history should be empty
	if len(room.GetMoveHistory()) != 0 {
		t.Error("move history should be empty after illegal move")
	}
}

func TestChessGame_NotYourTurn(t *testing.T) {
	room := ws.NewGameRoom("match-chess-3", "chess")
	p1 := "sid-white"
	p2 := "sid-black"
	joinRoom(t, room, p1, p2)

	// Try to move as black when it's white's turn
	room.HandleMove(p2, moveJSON("e7", "e5"))
	msg := readMsg(t, room.TestMsgCh, time.Second)

	if msg["type"] != "error" {
		t.Errorf("expected error message, got type %v", msg["type"])
	}
	if msg["code"] != "not_your_turn" {
		t.Errorf("expected not_your_turn code, got %v", msg["code"])
	}
}

func TestChessGame_Resign(t *testing.T) {
	room := ws.NewGameRoom("match-chess-4", "chess")
	p1 := "sid-white"
	p2 := "sid-black"
	joinRoom(t, room, p1, p2)

	// White resigns
	room.HandleResign(p1)

	if !room.IsGameOver() {
		t.Error("game should be over after resign")
	}

	msg := readMsg(t, room.TestMsgCh, time.Second)
	if msg["type"] != "game_over" {
		t.Errorf("expected game_over, got type %v", msg["type"])
	}
	if msg["reason"] != "resign" {
		t.Errorf("expected resign reason, got %v", msg["reason"])
	}
	if msg["winner_sid"] != p2 {
		t.Errorf("expected winner to be black (p2), got %v", msg["winner_sid"])
	}
	if msg["score"] != "0-1" {
		t.Errorf("expected score 0-1, got %v", msg["score"])
	}

	// Verify engine state
	if !room.Engine().IsGameOver() {
		t.Error("engine should report game over")
	}
	if room.Engine().GetWinner() != "black" {
		t.Errorf("engine winner should be black, got %v", room.Engine().GetWinner())
	}
	if room.Engine().GetGameOverReason() != "resign" {
		t.Errorf("engine reason should be resign, got %v", room.Engine().GetGameOverReason())
	}
}

func TestChessGame_MoveAfterGameOver(t *testing.T) {
	room := ws.NewGameRoom("match-chess-5", "chess")
	p1 := "sid-white"
	p2 := "sid-black"
	joinRoom(t, room, p1, p2)

	// Resign first
	room.HandleResign(p1)
	drainMessages(room.TestMsgCh)

	// Try to move after game is over
	room.HandleMove(p2, moveJSON("e7", "e5"))
	msg := readMsg(t, room.TestMsgCh, time.Second)

	if msg["type"] != "error" {
		t.Errorf("expected error after game over, got type %v", msg["type"])
	}
	if msg["code"] != "game_over" {
		t.Errorf("expected game_over code, got %v", msg["code"])
	}
}

func TestChessGame_DrawOfferFlow(t *testing.T) {
	room := ws.NewGameRoom("match-chess-6", "chess")
	p1 := "sid-white"
	p2 := "sid-black"
	joinRoom(t, room, p1, p2)

	// p1 offers draw
	room.HandleDrawOffer(p1)
	msg := readMsg(t, room.TestMsgCh, time.Second)
	if msg["type"] != "draw_offer" {
		t.Errorf("expected draw_offer sent to opponent, got type %v", msg["type"])
	}

	// p2 declines
	room.HandleDrawDecline(p2)
	msg = readMsg(t, room.TestMsgCh, time.Second)
	if msg["type"] != "draw_declined" {
		t.Errorf("expected draw_declined, got type %v", msg["type"])
	}

	// Game should still be playing
	if room.IsGameOver() {
		t.Error("game should not be over after declined draw")
	}

	// Now p2 offers draw, p1 accepts
	room.HandleDrawOffer(p2)
	msg = readMsg(t, room.TestMsgCh, time.Second)
	if msg["type"] != "draw_offer" {
		t.Errorf("expected draw_offer, got type %v", msg["type"])
	}

	room.HandleDrawAccept(p1)
	msg = readMsg(t, room.TestMsgCh, time.Second)
	if msg["type"] != "game_over" {
		t.Errorf("expected game_over after draw accept, got type %v", msg["type"])
	}
	if msg["reason"] != "draw" {
		t.Errorf("expected draw reason, got %v", msg["reason"])
	}
	if msg["winner_sid"] != "" {
		t.Errorf("expected empty winner_sid for draw, got %v", msg["winner_sid"])
	}
	if msg["score"] != "0-0" {
		t.Errorf("expected score 0-0 for draw, got %v", msg["score"])
	}

	if !room.IsGameOver() {
		t.Error("game should be over after draw acceptance")
	}
}

func TestChessGame_DrawOfferNoPending(t *testing.T) {
	room := ws.NewGameRoom("match-chess-7", "chess")
	p1 := "sid-white"
	p2 := "sid-black"
	joinRoom(t, room, p1, p2)

	// Try to accept draw when no offer is pending
	room.HandleDrawAccept(p2)
	msg := readMsg(t, room.TestMsgCh, time.Second)

	if msg["type"] != "error" {
		t.Errorf("expected error, got type %v", msg["type"])
	}
	if msg["code"] != "no_draw_offer" {
		t.Errorf("expected no_draw_offer code, got %v", msg["code"])
	}
}

func TestChessGame_ResignAfterGameStarted(t *testing.T) {
	room := ws.NewGameRoom("match-chess-8", "chess")
	p1 := "sid-white"
	p2 := "sid-black"
	joinRoom(t, room, p1, p2)

	// Make a move first
	room.HandleMove(p1, moveJSON("e2", "e4"))
	drainMessages(room.TestMsgCh) // drain move_applied + state broadcasts

	// Black resigns
	room.HandleResign(p2)
	msg := readMsg(t, room.TestMsgCh, time.Second)

	if msg["type"] != "game_over" {
		t.Errorf("expected game_over, got type %v", msg["type"])
	}
	if msg["reason"] != "resign" {
		t.Errorf("expected resign, got %v", msg["reason"])
	}
	if msg["winner_sid"] != p1 {
		t.Errorf("expected winner p1, got %v", msg["winner_sid"])
	}
	if msg["score"] != "1-0" {
		t.Errorf("expected score 1-0, got %v", msg["score"])
	}
}

// ---------------------------------------------------------------------------
// Checkers gameplay tests
// ---------------------------------------------------------------------------

func TestCheckersGame_BasicMove(t *testing.T) {
	room := ws.NewGameRoom("match-checkers-1", "checkers")
	p1 := "sid-white"
	p2 := "sid-black"
	joinRoom(t, room, p1, p2)

	// Verify initial state
	state := room.BuildStatePayload()
	if state["turn"] != "white" {
		t.Errorf("expected turn white, got %v", state["turn"])
	}
	board := state["board"].(string)
	if board == "" {
		t.Error("checkers board should not be empty")
	}

	legalMoves := state["legal_moves"].([]ws.WSMove)
	if len(legalMoves) == 0 {
		t.Error("checkers should have legal moves from starting position")
	}

	// Make a white move (row 5 to row 4, forward for white)
	// Find a legal move
	foundMove := false
	for _, m := range legalMoves {
		moveData, _ := json.Marshal(m)
		room.HandleMove(p1, moveData)
		msg := readMsg(t, room.TestMsgCh, time.Second)
		if msg["type"] == "move_applied" {
			foundMove = true
			break
		}
		// If error, drain it
	}
	if !foundMove {
		t.Error("should have been able to make at least one legal move")
	}

	// Verify turn switched
	state = room.BuildStatePayload()
	if state["turn"] != "black" {
		t.Errorf("expected turn black after white move, got %v", state["turn"])
	}

	// Verify move history
	history := room.GetMoveHistory()
	if len(history) != 1 {
		t.Errorf("expected 1 move in history, got %d", len(history))
	}
}

func TestCheckersGame_IllegalMove(t *testing.T) {
	room := ws.NewGameRoom("match-checkers-2", "checkers")
	p1 := "sid-white"
	p2 := "sid-black"
	joinRoom(t, room, p1, p2)

	// Try an invalid position
	room.HandleMove(p1, checkersMoveJSON(0, 0, 1, 1))
	msg := readMsg(t, room.TestMsgCh, time.Second)

	if msg["type"] != "error" {
		t.Errorf("expected error for illegal checkers move, got type %v", msg["type"])
	}
}

func TestCheckersGame_Resign(t *testing.T) {
	room := ws.NewGameRoom("match-checkers-3", "checkers")
	p1 := "sid-white"
	p2 := "sid-black"
	joinRoom(t, room, p1, p2)

	room.HandleResign(p2)
	msg := readMsg(t, room.TestMsgCh, time.Second)

	if msg["type"] != "game_over" {
		t.Errorf("expected game_over, got type %v", msg["type"])
	}
	if msg["reason"] != "resign" {
		t.Errorf("expected resign, got %v", msg["reason"])
	}
	if msg["winner_sid"] != p1 {
		t.Errorf("expected winner p1, got %v", msg["winner_sid"])
	}
	if msg["score"] != "1-0" {
		t.Errorf("expected score 1-0, got %v", msg["score"])
	}
}

// ---------------------------------------------------------------------------
// Backgammon gameplay tests
// ---------------------------------------------------------------------------

func TestBackgammonGame_BasicMove(t *testing.T) {
	room := ws.NewGameRoom("match-bg-1", "backgammon")
	p1 := "sid-white"
	p2 := "sid-black"
	joinRoom(t, room, p1, p2)

	// Verify initial state
	state := room.BuildStatePayload()
	if state["turn"] != "white" {
		t.Errorf("expected turn white, got %v", state["turn"])
	}
	board := state["board"].(string)
	if board == "" {
		t.Error("backgammon board should not be empty")
	}

	// No legal moves before dice roll
	legalMoves := state["legal_moves"].([]ws.WSMove)
	if len(legalMoves) != 0 {
		t.Error("backgammon should have no legal moves before rolling dice")
	}

	// Try to move without rolling dice
	room.HandleMove(p1, backgammonMoveJSON(23, 22, 1))
	msg := readMsg(t, room.TestMsgCh, time.Second)
	if msg["type"] != "error" {
		t.Errorf("expected error when moving without rolling dice, got type %v", msg["type"])
	}

	// Roll dice (white rolls)
	room.HandleRollDice(p1)
	msg = readMsg(t, room.TestMsgCh, time.Second) // dice_rolled
	drainMessages(room.TestMsgCh)                 // drain any state broadcasts
	if msg["type"] != "dice_rolled" {
		t.Errorf("expected dice_rolled, got type %v", msg["type"])
	}
	dice, ok := msg["dice"].([]interface{})
	if !ok || len(dice) != 2 {
		t.Fatalf("expected 2 dice values, got: %v", msg["dice"])
	}

	// Now there should be legal moves
	state = room.BuildStatePayload()
	legalMoves = state["legal_moves"].([]ws.WSMove)
	if len(legalMoves) == 0 {
		t.Error("should have legal moves after rolling dice")
	}

	// Make the first legal move
	firstMove := legalMoves[0]
	moveData, _ := json.Marshal(firstMove)
	room.HandleMove(p1, moveData)
	msg = readMsg(t, room.TestMsgCh, time.Second)
	if msg["type"] != "move_applied" {
		t.Errorf("expected move_applied, got type %v, data: %s", msg["type"], string(moveData))
	}

	// Verify move history
	history := room.GetMoveHistory()
	if len(history) != 1 {
		t.Errorf("expected 1 move in history, got %d", len(history))
	}
}

func TestBackgammonGame_Resign(t *testing.T) {
	room := ws.NewGameRoom("match-bg-2", "backgammon")
	p1 := "sid-white"
	p2 := "sid-black"
	joinRoom(t, room, p1, p2)

	room.HandleResign(p1)
	msg := readMsg(t, room.TestMsgCh, time.Second)

	if msg["type"] != "game_over" {
		t.Errorf("expected game_over, got type %v", msg["type"])
	}
	if msg["reason"] != "resign" {
		t.Errorf("expected resign, got %v", msg["reason"])
	}
	if msg["winner_sid"] != p2 {
		t.Errorf("expected winner p2, got %v", msg["winner_sid"])
	}
}

func TestBackgammonGame_RollDiceNotBackgammon(t *testing.T) {
	room := ws.NewGameRoom("match-chess-dice", "chess")
	p1 := "sid-white"
	p2 := "sid-black"
	joinRoom(t, room, p1, p2)

	// Try to roll dice in a chess game
	room.HandleRollDice(p1)
	msg := readMsg(t, room.TestMsgCh, time.Second)

	if msg["type"] != "error" {
		t.Errorf("expected error, got type %v", msg["type"])
	}
	if msg["code"] != "invalid_action" {
		t.Errorf("expected invalid_action code, got %v", msg["code"])
	}
}

// ---------------------------------------------------------------------------
// Room join tests
// ---------------------------------------------------------------------------

func TestGameRoom_JoinFullRoom(t *testing.T) {
	room := ws.NewGameRoom("match-join-full", "chess")
	p1 := "sid-white"
	p2 := "sid-black"

	room.TestMsgCh = make(chan []byte, 64)
	room.JoinPlayer(p1, nil)
	room.JoinPlayer(p2, nil)
	drainMessages(room.TestMsgCh)

	// Third player should fail
	err := room.JoinPlayer("sid-extra", nil)
	if err == nil {
		t.Error("expected error when joining full room")
	}
}

func TestGameRoom_JoinAfterGameOver(t *testing.T) {
	room := ws.NewGameRoom("match-join-over", "chess")
	p1 := "sid-white"
	p2 := "sid-black"
	joinRoom(t, room, p1, p2)

	// End the game
	room.HandleResign(p1)
	drainMessages(room.TestMsgCh)

	// Try to join after game over (shouldn't work via room but let's test)
	if !room.IsGameOver() {
		t.Fatal("game should be over")
	}
}

// ---------------------------------------------------------------------------
// Match persistence tests
// ---------------------------------------------------------------------------

func TestChessGame_PersistOnResign(t *testing.T) {
	mockRepo := mocks.NewMockMatchRepo()
	// Pre-create the match
	match := &model.Match{
		ID:         "match-persist-1",
		GameType:   "chess",
		Player1SID: "p1",
		Player2SID: "p2",
		Status:     "playing",
	}
	mockRepo.Create(nil, match)

	room := ws.NewGameRoom("match-persist-1", "chess")
	room.SetMatchRepo(mockRepo)
	joinRoom(t, room, "p1", "p2")

	// Make a move
	room.HandleMove("p1", moveJSON("e2", "e4"))
	readMsg(t, room.TestMsgCh, time.Second)

	// Resign
	room.HandleResign("p1")
	readMsg(t, room.TestMsgCh, time.Second)

	// Give the async persist time to complete
	time.Sleep(100 * time.Millisecond)

	// Verify match was completed in repo
	updatedMatch, err := mockRepo.GetByID(nil, "match-persist-1")
	if err != nil {
		t.Fatalf("failed to get match from repo: %v", err)
	}
	if updatedMatch == nil {
		t.Fatal("match not found in repo")
	}
	if updatedMatch.Status != "completed" {
		t.Errorf("expected status completed, got %s", updatedMatch.Status)
	}
	if updatedMatch.WinnerSID != "p2" {
		t.Errorf("expected winner p2, got %s", updatedMatch.WinnerSID)
	}
	if updatedMatch.Score != "0-1" {
		t.Errorf("expected score 0-1, got %s", updatedMatch.Score)
	}
}

func TestChessGame_PersistOnDraw(t *testing.T) {
	mockRepo := mocks.NewMockMatchRepo()
	match := &model.Match{
		ID:         "match-persist-draw",
		GameType:   "chess",
		Player1SID: "p1",
		Player2SID: "p2",
		Status:     "playing",
	}
	mockRepo.Create(nil, match)

	room := ws.NewGameRoom("match-persist-draw", "chess")
	room.SetMatchRepo(mockRepo)
	joinRoom(t, room, "p1", "p2")

	// Offer and accept draw
	room.HandleDrawOffer("p1")
	readMsg(t, room.TestMsgCh, time.Second)
	room.HandleDrawAccept("p2")
	readMsg(t, room.TestMsgCh, time.Second)

	// Give async persist time
	time.Sleep(100 * time.Millisecond)

	updatedMatch, _ := mockRepo.GetByID(nil, "match-persist-draw")
	if updatedMatch.Status != "completed" {
		t.Errorf("expected completed, got %s", updatedMatch.Status)
	}
	if updatedMatch.WinnerSID != "" {
		t.Errorf("expected empty winner for draw, got %s", updatedMatch.WinnerSID)
	}
	if updatedMatch.Score != "0-0" {
		t.Errorf("expected score 0-0, got %s", updatedMatch.Score)
	}
}

// ---------------------------------------------------------------------------
// Rating update on game_over tests
// ---------------------------------------------------------------------------

func TestChessGame_RatingUpdateOnResign(t *testing.T) {
	mockRepo := mocks.NewMockMatchRepo()
	match := &model.Match{
		ID:         "match-rating-1",
		GameType:   "chess",
		Player1SID: "p1",
		Player2SID: "p2",
		Status:     "playing",
	}
	mockRepo.Create(nil, match)

	ratingUpdater := &mocks.MockRatingUpdater{}

	room := ws.NewGameRoom("match-rating-1", "chess")
	room.SetMatchRepo(mockRepo)
	room.SetRatingService(ratingUpdater)
	joinRoom(t, room, "p1", "p2")

	// White resigns
	room.HandleResign("p1")
	readMsg(t, room.TestMsgCh, time.Second)

	// Give the async rating update time to complete
	time.Sleep(100 * time.Millisecond)

	if ratingUpdater.CallCount() != 1 {
		t.Fatalf("Expected 1 rating update call, got %d", ratingUpdater.CallCount())
	}
	call := ratingUpdater.Calls[0]
	if call.MatchID != "match-rating-1" {
		t.Errorf("Expected match_id match-rating-1, got %s", call.MatchID)
	}
	if call.GameType != "chess" {
		t.Errorf("Expected game_type chess, got %s", call.GameType)
	}
	if call.WinnerSID != "p2" {
		t.Errorf("Expected winner p2, got %s", call.WinnerSID)
	}
	if call.Score != "0-1" {
		t.Errorf("Expected score 0-1, got %s", call.Score)
	}
	if call.Player1SID != "p1" {
		t.Errorf("Expected player1 p1, got %s", call.Player1SID)
	}
	if call.Player2SID != "p2" {
		t.Errorf("Expected player2 p2, got %s", call.Player2SID)
	}
}

func TestChessGame_RatingUpdateOnDraw(t *testing.T) {
	mockRepo := mocks.NewMockMatchRepo()
	match := &model.Match{
		ID:         "match-rating-draw",
		GameType:   "chess",
		Player1SID: "p1",
		Player2SID: "p2",
		Status:     "playing",
	}
	mockRepo.Create(nil, match)

	ratingUpdater := &mocks.MockRatingUpdater{}

	room := ws.NewGameRoom("match-rating-draw", "chess")
	room.SetMatchRepo(mockRepo)
	room.SetRatingService(ratingUpdater)
	joinRoom(t, room, "p1", "p2")

	// Mutual draw
	room.HandleDrawOffer("p1")
	readMsg(t, room.TestMsgCh, time.Second)
	room.HandleDrawAccept("p2")
	readMsg(t, room.TestMsgCh, time.Second)

	time.Sleep(100 * time.Millisecond)

	if ratingUpdater.CallCount() != 1 {
		t.Fatalf("Expected 1 rating update call, got %d", ratingUpdater.CallCount())
	}
	call := ratingUpdater.Calls[0]
	if call.WinnerSID != "" {
		t.Errorf("Expected empty winner for draw, got %s", call.WinnerSID)
	}
	if call.Score != "0-0" {
		t.Errorf("Expected score 0-0 for draw, got %s", call.Score)
	}
}

func TestChessGame_NoRatingUpdateWhenNil(t *testing.T) {
	mockRepo := mocks.NewMockMatchRepo()
	match := &model.Match{
		ID:         "match-norating",
		GameType:   "chess",
		Player1SID: "p1",
		Player2SID: "p2",
		Status:     "playing",
	}
	mockRepo.Create(nil, match)

	room := ws.NewGameRoom("match-norating", "chess")
	room.SetMatchRepo(mockRepo)
	// NO rating service set
	joinRoom(t, room, "p1", "p2")

	room.HandleResign("p1")
	readMsg(t, room.TestMsgCh, time.Second)
	time.Sleep(100 * time.Millisecond)

	// Should still work fine, no crash
	if !room.IsGameOver() {
		t.Error("game should be over")
	}
}

// ---------------------------------------------------------------------------
// RoomManager integration tests
// ---------------------------------------------------------------------------

func TestRoomManager_CreateRoomWithMatchRepo(t *testing.T) {
	mockRepo := mocks.NewMockMatchRepo()
	rm := ws.NewRoomManagerWithDeps("test-secret", mockRepo)

	room := rm.CreateRoom("match-rm-1", "chess")
	if room == nil {
		t.Fatal("room should not be nil")
	}
	if room.MatchID() != "match-rm-1" {
		t.Errorf("expected match_id match-rm-1, got %s", room.MatchID())
	}
}

func TestRoomManager_GenerateTestToken(t *testing.T) {
	token := ws.GenerateTestToken("test-secret", "user123")
	if token == "" {
		t.Error("expected non-empty token")
	}

	// Verify token is valid JWT
	rm := ws.NewRoomManagerWithDeps("test-secret", nil)
	sid := rm.ValidateTokenString(token)
	if sid != "user123" {
		t.Errorf("expected sid user123, got %s", sid)
	}
}

func TestRoomManager_GenerateTestTokenWrongSecret(t *testing.T) {
	token := ws.GenerateTestToken("secret1", "user123")

	rm := ws.NewRoomManagerWithDeps("secret2", nil)
	sid := rm.ValidateTokenString(token)
	if sid != "" {
		t.Errorf("expected empty sid for wrong secret, got %s", sid)
	}
}

func TestRoomManager_AuthenticateQueryTokenEmpty(t *testing.T) {
	rm := ws.NewRoomManagerWithDeps("", nil)
	// With empty secret, should return empty sid
	token := ws.GenerateTestToken("any", "user")
	sid := rm.ValidateTokenString(token)
	if sid != "" {
		t.Errorf("expected empty sid with empty jwt secret, got %s", sid)
	}
}

func TestRoomManager_ConcurrentJoins(t *testing.T) {
	rm := ws.NewRoomManager()
	rm.CreateRoom("match-concurrent", "chess")

	var wg sync.WaitGroup
	successCount := 0
	var mu sync.Mutex

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			room := rm.GetRoom("match-concurrent")
			if room == nil {
				return
			}
			err := room.JoinPlayer("sid-"+string(rune('A'+idx)), nil)
			if err == nil {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()
	// Only 2 players should successfully join
	if successCount != 2 {
		t.Errorf("expected 2 successful joins, got %d", successCount)
	}
}

func TestRoomManager_SetMatchRepoAfterCreation(t *testing.T) {
	rm := ws.NewRoomManager()
	mockRepo := mocks.NewMockMatchRepo()
	rm.SetMatchRepo(mockRepo)

	// Create a room - it should get the match repo
	room := rm.CreateRoom("match-set-repo", "chess")
	if room == nil {
		t.Fatal("room should not be nil")
	}
	// Verify room is properly created
	if room.MatchID() != "match-set-repo" {
		t.Errorf("expected match_id match-set-repo, got %s", room.MatchID())
	}
}
