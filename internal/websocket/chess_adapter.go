package websocket

import (
	"fmt"

	"game-platform/internal/game/chess"
)

// ChessAdapter adapts the chess engine to the GameEngineAdapter interface.
type ChessAdapter struct {
	game *chess.Game
}

// NewChessAdapter creates a new chess adapter.
func NewChessAdapter() *ChessAdapter {
	return &ChessAdapter{
		game: chess.NewGame(),
	}
}

// GetBoard returns the FEN string for the current position.
func (a *ChessAdapter) GetBoard() string {
	return a.game.ToFEN()
}

// GetTurn returns "white" or "black".
func (a *ChessAdapter) GetTurn() string {
	if a.game.Turn == chess.White {
		return "white"
	}
	return "black"
}

// GetLegalMoves returns legal moves in WS format.
func (a *ChessAdapter) GetLegalMoves() []WSMove {
	allMoves := a.game.AllLegalMoves()
	result := make([]WSMove, 0, len(allMoves))
	for _, m := range allMoves {
		wsMove := WSMove{
			From: m.From.Algebraic(),
			To:   m.To.Algebraic(),
		}
		if m.Promotion != chess.NoPiece {
			wsMove.Promotion = promotionToString(m.Promotion)
		}
		result = append(result, wsMove)
	}
	return result
}

// ApplyMove converts the WS move to engine format and applies it.
func (a *ChessAdapter) ApplyMove(move WSMove) error {
	from := chess.FromAlgebraic(move.From)
	to := chess.FromAlgebraic(move.To)
	if !from.Valid() || !to.Valid() {
		return fmt.Errorf("invalid position")
	}

	var promotion chess.PieceType = chess.NoPiece
	if move.Promotion != "" {
		promotion = stringToPromotion(move.Promotion)
	}

	return a.game.MakeMove(from, to, promotion)
}

// IsGameOver returns true if the game has ended.
func (a *ChessAdapter) IsGameOver() bool {
	state := a.game.State
	return state == chess.Checkmate || state == chess.Stalemate || state == chess.Draw || state == chess.Resigned || state == chess.Timeout
}

// GetGameOverReason returns the reason the game ended.
func (a *ChessAdapter) GetGameOverReason() string {
	switch a.game.State {
	case chess.Checkmate:
		return "checkmate"
	case chess.Stalemate:
		return "stalemate"
	case chess.Draw:
		return "draw"
	case chess.Resigned:
		return "resign"
	case chess.Timeout:
		return "timeout"
	}
	return ""
}

// GetWinner returns the winning color or empty for draw.
func (a *ChessAdapter) GetWinner() string {
	switch a.game.State {
	case chess.Checkmate, chess.Resigned, chess.Timeout:
		// The winner is the opposite of the current turn
		if a.game.Turn == chess.White {
			return "black"
		}
		return "white"
	case chess.Stalemate, chess.Draw:
		return ""
	}
	return ""
}

// Resign resigns for the given color.
func (a *ChessAdapter) Resign(color string) {
	if color == "white" {
		a.game.Resign(chess.White)
	} else {
		a.game.Resign(chess.Black)
	}
}

// RollDice is a no-op for chess.
func (a *ChessAdapter) RollDice() []int {
	return nil
}

func promotionToString(pt chess.PieceType) string {
	switch pt {
	case chess.Queen:
		return "q"
	case chess.Rook:
		return "r"
	case chess.Bishop:
		return "b"
	case chess.Knight:
		return "n"
	}
	return ""
}

func stringToPromotion(s string) chess.PieceType {
	switch s {
	case "q":
		return chess.Queen
	case "r":
		return chess.Rook
	case "b":
		return chess.Bishop
	case "n":
		return chess.Knight
	}
	return chess.NoPiece
}
