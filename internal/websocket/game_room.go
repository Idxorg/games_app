package websocket

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"game-platform/internal/model"

	"github.com/gorilla/websocket"
)

// GameRoom represents a single match room with two players.
type GameRoom struct {
	mu sync.Mutex

	matchID     string
	gameType    string
	player1SID  string
	player2SID  string
	player1Conn *websocket.Conn
	player2Conn *websocket.Conn

	engine GameEngineAdapter
	clock  *Clock

	gameOver       bool
	gameOverReason string
	winnerSID      string

	drawOfferSID string // SID of the player offering a draw

	clockStopCh chan struct{}

	// Persistence
	matchRepo   model.MatchRepo
	moveHistory []map[string]interface{} // tracked moves for match persistence

	// Test hook: when non-nil, all broadcast/error messages are also sent here.
	TestMsgCh chan []byte
}

// NewGameRoom creates a new game room.
func NewGameRoom(matchID, gameType string) *GameRoom {
	room := &GameRoom{
		matchID:     matchID,
		gameType:    gameType,
		clock:       NewDefaultClock(),
		clockStopCh: make(chan struct{}),
	}

	switch gameType {
	case "chess":
		room.engine = NewChessAdapter()
	case "checkers":
		room.engine = NewCheckersAdapter()
	case "backgammon":
		room.engine = NewBackgammonAdapter()
	default:
		slog.Error("unknown game type", "game_type", gameType)
		room.engine = NewChessAdapter() // default fallback
	}

	return room
}

// SetMatchRepo sets the match repository for persistence.
func (r *GameRoom) SetMatchRepo(repo model.MatchRepo) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.matchRepo = repo
}

// Engine returns the game engine adapter (for testing/introspection).
func (r *GameRoom) Engine() GameEngineAdapter {
	return r.engine
}

// MatchID returns the match ID.
func (r *GameRoom) MatchID() string {
	return r.matchID
}

// GameType returns the game type.
func (r *GameRoom) GameType() string {
	return r.gameType
}

// IsFull returns true if both players have joined.
func (r *GameRoom) IsFull() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.player1SID != "" && r.player2SID != ""
}

// IsGameOver returns true if the game is over.
func (r *GameRoom) IsGameOver() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.gameOver
}

// JoinPlayer adds a player to the room. Returns an error if room is full.
func (r *GameRoom) JoinPlayer(sid string, conn *websocket.Conn) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.gameOver {
		return sendError(conn, "game_over", "game has already ended")
	}

	if r.player1SID == "" {
		r.player1SID = sid
		r.player1Conn = conn
		slog.Info("player 1 joined room", "match_id", r.matchID, "sid", sid)
	} else if r.player2SID == "" {
		r.player2SID = sid
		r.player2Conn = conn
		slog.Info("player 2 joined room", "match_id", r.matchID, "sid", sid)

		// Start the clock when both players are in
		go r.runClock()

		// Broadcast initial state to both players
		r.broadcastState()
	} else {
		return sendError(conn, "room_full", "room is already full")
	}

	return nil
}

// GetPlayerColor returns "white" or "black" for the given SID.
func (r *GameRoom) GetPlayerColor(sid string) string {
	r.mu.Lock()
	defer r.mu.Unlock()
	if sid == r.player1SID {
		return "white"
	}
	if sid == r.player2SID {
		return "black"
	}
	return ""
}

// HandleMove processes a move from a player.
func (r *GameRoom) HandleMove(sid string, moveData json.RawMessage) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.gameOver {
		r.sendToPlayer(sid, errorPayload("game_over", "game has already ended"))
		return
	}

	// Verify it's this player's turn
	turn := r.engine.GetTurn()
	color := r.GetPlayerColorUnlocked(sid)
	if color == "" {
		r.sendToPlayer(sid, errorPayload("not_in_game", "you are not in this game"))
		return
	}
	if color != turn {
		r.sendToPlayer(sid, errorPayload("not_your_turn", "it is not your turn"))
		return
	}

	var move WSMove
	if err := json.Unmarshal(moveData, &move); err != nil {
		r.sendToPlayer(sid, errorPayload("invalid_move", "invalid move format"))
		return
	}

	if err := r.engine.ApplyMove(move); err != nil {
		slog.Info("illegal move", "match_id", r.matchID, "sid", sid, "error", err)
		r.sendToPlayer(sid, errorPayload("illegal_move", err.Error()))
		return
	}

	// Track move in history for persistence
	r.moveHistory = append(r.moveHistory, map[string]interface{}{
		"sid":   sid,
		"from":  move.From,
		"to":    move.To,
		"color": color,
	})

	// Switch clock and add increment
	r.clock.SwitchPlayer(turn, DefaultIncrement)

	// Broadcast move_applied + updated state
	r.broadcastMoveApplied(move)
	r.broadcastState()

	// Check for game over
	if r.engine.IsGameOver() {
		r.handleEngineGameOver()
	}
}

// HandleResign processes a resignation.
func (r *GameRoom) HandleResign(sid string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.gameOver {
		r.sendToPlayer(sid, errorPayload("game_over", "game has already ended"))
		return
	}

	color := r.GetPlayerColorUnlocked(sid)
	if color == "" {
		r.sendToPlayer(sid, errorPayload("not_in_game", "you are not in this game"))
		return
	}

	r.engine.Resign(color)
	r.clock.Stop()
	r.gameOver = true
	r.gameOverReason = "resign"

	// Determine winner SID
	if color == "white" {
		r.winnerSID = r.player2SID
	} else {
		r.winnerSID = r.player1SID
	}

	slog.Info("player resigned", "match_id", r.matchID, "sid", sid, "color", color)

	r.broadcastGameOver()
	r.persistGameComplete()
}

// HandleDrawOffer processes a draw offer.
func (r *GameRoom) HandleDrawOffer(sid string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.gameOver {
		r.sendToPlayer(sid, errorPayload("game_over", "game has already ended"))
		return
	}

	// If the other player already offered, this is a mutual draw
	if r.drawOfferSID != "" && r.drawOfferSID != sid {
		r.clock.Stop()
		r.gameOver = true
		r.gameOverReason = "draw"
		r.winnerSID = ""
		slog.Info("mutual draw", "match_id", r.matchID)
		r.broadcastGameOver()
		r.persistGameComplete()
		return
	}

	r.drawOfferSID = sid

	// Send draw_offer to the other player
	targetSID := r.getOpponentSID(sid)
	if targetSID != "" {
		msg, _ := json.Marshal(map[string]string{
			"type":     "draw_offer",
			"match_id": r.matchID,
		})
		r.sendToPlayer(targetSID, msg)
	}

	slog.Info("draw offer sent", "match_id", r.matchID, "from_sid", sid)
}

// HandleDrawAccept processes a draw acceptance.
func (r *GameRoom) HandleDrawAccept(sid string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.gameOver {
		r.sendToPlayer(sid, errorPayload("game_over", "game has already ended"))
		return
	}

	if r.drawOfferSID == "" {
		r.sendToPlayer(sid, errorPayload("no_draw_offer", "no pending draw offer"))
		return
	}

	// The accept must come from the opponent of the offerer
	targetSID := r.getOpponentSID(r.drawOfferSID)
	if sid != targetSID {
		r.sendToPlayer(sid, errorPayload("not_your_offer", "you cannot accept this draw offer"))
		return
	}

	r.clock.Stop()
	r.gameOver = true
	r.gameOverReason = "draw"
	r.winnerSID = ""
	r.drawOfferSID = ""

	slog.Info("draw accepted", "match_id", r.matchID)

	r.broadcastGameOver()
	r.persistGameComplete()
}

// HandleDrawDecline processes a draw decline.
func (r *GameRoom) HandleDrawDecline(sid string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.drawOfferSID == "" {
		r.sendToPlayer(sid, errorPayload("no_draw_offer", "no pending draw offer"))
		return
	}

	targetSID := r.getOpponentSID(r.drawOfferSID)
	if sid != targetSID {
		r.sendToPlayer(sid, errorPayload("not_your_offer", "you cannot decline this draw offer"))
		return
	}

	// Notify the offerer that their draw was declined
	msg, _ := json.Marshal(map[string]string{
		"type":     "draw_declined",
		"match_id": r.matchID,
	})
	r.sendToPlayer(r.drawOfferSID, msg)
	r.drawOfferSID = ""
}

// HandleRollDice handles a dice roll (backgammon only).
func (r *GameRoom) HandleRollDice(sid string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.gameOver {
		r.sendToPlayer(sid, errorPayload("game_over", "game has already ended"))
		return
	}

	if r.gameType != "backgammon" {
		r.sendToPlayer(sid, errorPayload("invalid_action", "dice roll only for backgammon"))
		return
	}

	dice := r.engine.RollDice()
	msg, _ := json.Marshal(map[string]interface{}{
		"type":     "dice_rolled",
		"match_id": r.matchID,
		"dice":     dice,
	})
	r.broadcast(msg)
}

// GetMoveHistory returns the move history (for persistence/testing).
func (r *GameRoom) GetMoveHistory() []map[string]interface{} {
	r.mu.Lock()
	defer r.mu.Unlock()
	result := make([]map[string]interface{}, len(r.moveHistory))
	copy(result, r.moveHistory)
	return result
}

// Cleanup cleans up room resources when the room is destroyed.
func (r *GameRoom) Cleanup() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.clock.Stop()
	if r.clockStopCh != nil {
		select {
		case <-r.clockStopCh:
		default:
			close(r.clockStopCh)
		}
	}

	if r.player1Conn != nil {
		r.player1Conn.Close()
	}
	if r.player2Conn != nil {
		r.player2Conn.Close()
	}
}

// runClock ticks the clock every 100ms and checks for timeout.
func (r *GameRoom) runClock() {
	ticker := time.NewTicker(ClockTickInterval)
	defer ticker.Stop()

	r.clock.Start(r.engine.GetTurn())

	for {
		select {
		case <-ticker.C:
			r.mu.Lock()
			if r.gameOver {
				r.mu.Unlock()
				return
			}
			activeColor := r.engine.GetTurn()
			r.mu.Unlock()

			if r.clock.Tick(activeColor) {
				r.mu.Lock()
				if !r.gameOver {
					r.gameOver = true
					r.gameOverReason = "timeout"
					timeoutColor := r.clock.TimeoutColor()
					if timeoutColor == "white" {
						r.winnerSID = r.player2SID
					} else {
						r.winnerSID = r.player1SID
					}
					r.broadcastGameOver()
					r.persistGameComplete()
					slog.Info("timeout", "match_id", r.matchID, "timeout_color", timeoutColor)
				}
				r.mu.Unlock()
				return
			}

		case <-r.clockStopCh:
			return
		}
	}
}

// handleEngineGameOver is called when the engine detects game over (checkmate, etc.)
func (r *GameRoom) handleEngineGameOver() {
	r.gameOver = true
	r.gameOverReason = r.engine.GetGameOverReason()
	winner := r.engine.GetWinner()

	if winner == "white" {
		r.winnerSID = r.player1SID
	} else if winner == "black" {
		r.winnerSID = r.player2SID
	}

	r.clock.Stop()
	r.broadcastGameOver()
	r.persistGameComplete()
	slog.Info("game over (engine)", "match_id", r.matchID, "reason", r.gameOverReason, "winner", winner)
}

// persistGameComplete persists the game result to the match repo if available.
// Must be called with r.mu held.
func (r *GameRoom) persistGameComplete() {
	if r.matchRepo == nil {
		return
	}

	score := "0-0"
	if r.winnerSID != "" {
		if r.winnerSID == r.player1SID {
			score = "1-0"
		} else {
			score = "0-1"
		}
	}

	movesJSON, _ := json.Marshal(r.moveHistory)
	go func() {
		err := r.matchRepo.Complete(nil, r.matchID, r.winnerSID, score, movesJSON)
		if err != nil {
			slog.Error("failed to persist game complete", "match_id", r.matchID, "error", err)
		}
	}()
}

// broadcastState sends the full game state to both players.
func (r *GameRoom) broadcastState() {
	state := r.BuildStatePayload()
	data, err := json.Marshal(state)
	if err != nil {
		slog.Error("failed to marshal state", "error", err)
		return
	}
	r.broadcast(data)
}

// broadcastMoveApplied sends move_applied + state to both players.
func (r *GameRoom) broadcastMoveApplied(move WSMove) {
	moveApplied := map[string]interface{}{
		"type":     "move_applied",
		"match_id": r.matchID,
		"move":     move,
		"state":    r.BuildStatePayload(),
	}
	data, err := json.Marshal(moveApplied)
	if err != nil {
		slog.Error("failed to marshal move_applied", "error", err)
		return
	}
	r.broadcast(data)
}

// broadcastGameOver sends game_over to both players.
func (r *GameRoom) broadcastGameOver() {
	score := "0-0"
	if r.winnerSID != "" {
		if r.winnerSID == r.player1SID {
			score = "1-0"
		} else {
			score = "0-1"
		}
	}

	msg := map[string]interface{}{
		"type":       "game_over",
		"match_id":   r.matchID,
		"winner_sid": r.winnerSID,
		"score":      score,
		"reason":     r.gameOverReason,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		slog.Error("failed to marshal game_over", "error", err)
		return
	}
	r.broadcast(data)
}

// BuildStatePayload builds the state message payload.
func (r *GameRoom) BuildStatePayload() map[string]interface{} {
	return map[string]interface{}{
		"type":        "state",
		"match_id":    r.matchID,
		"game_type":   r.gameType,
		"board":       r.engine.GetBoard(),
		"turn":        r.engine.GetTurn(),
		"player1_sid": r.player1SID,
		"player2_sid": r.player2SID,
		"legal_moves": r.engine.GetLegalMoves(),
		"clock": map[string]int64{
			"white_ms": r.clock.WhiteMs(),
			"black_ms": r.clock.BlackMs(),
		},
	}
}

// sendToPlayer sends a message to a specific player by SID.
// Must be called with r.mu held.
func (r *GameRoom) sendToPlayer(sid string, message []byte) {
	var conn *websocket.Conn
	if sid == r.player1SID {
		conn = r.player1Conn
	} else if sid == r.player2SID {
		conn = r.player2Conn
	}

	if conn != nil {
		conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		conn.WriteMessage(websocket.TextMessage, message)
	}

	// Test hook: capture message
	if r.TestMsgCh != nil {
		select {
		case r.TestMsgCh <- message:
		default:
		}
	}
}

// broadcast sends a message to both players.
// Must be called with r.mu held.
func (r *GameRoom) broadcast(message []byte) {
	r.sendToPlayer(r.player1SID, message)
	r.sendToPlayer(r.player2SID, message)
}

// GetPlayerColorUnlocked returns the color for a SID (unlocked version).
// Must be called with r.mu held.
func (r *GameRoom) GetPlayerColorUnlocked(sid string) string {
	if sid == r.player1SID {
		return "white"
	}
	if sid == r.player2SID {
		return "black"
	}
	return ""
}

// getOpponentSID returns the SID of the other player.
// Must be called with r.mu held.
func (r *GameRoom) getOpponentSID(sid string) string {
	if sid == r.player1SID {
		return r.player2SID
	}
	if sid == r.player2SID {
		return r.player1SID
	}
	return ""
}

// sendError sends an error message to a connection. Handles nil conn gracefully.
func sendError(conn *websocket.Conn, code, message string) error {
	msg := errorPayload(code, message)
	if conn != nil {
		conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		return conn.WriteMessage(websocket.TextMessage, msg)
	}
	return fmt.Errorf("%s: %s", code, message)
}

func errorPayload(code, message string) []byte {
	data, _ := json.Marshal(map[string]string{
		"type":    "error",
		"code":    code,
		"message": message,
	})
	return data
}
