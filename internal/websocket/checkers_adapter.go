package websocket

import (
	"fmt"

	"game-platform/internal/game/checkers"
)

// CheckersAdapter adapts the checkers engine to the GameEngineAdapter interface.
type CheckersAdapter struct {
	game *checkers.Game
}

// NewCheckersAdapter creates a new checkers adapter.
func NewCheckersAdapter() *CheckersAdapter {
	return &CheckersAdapter{
		game: checkers.NewGame(),
	}
}

// GetBoard returns a string representation of the board.
// For checkers we encode the board as a simple string: each row is encoded
// as 8 chars: w=white man, b=black man, W=white king, B=black king, .=empty.
func (a *CheckersAdapter) GetBoard() string {
	var result string
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			p := a.game.Board[row][col]
			if p.Empty() {
				result += "."
			} else {
				switch {
				case p.Type == checkers.Man && p.Color == checkers.White:
					result += "w"
				case p.Type == checkers.Man && p.Color == checkers.Black:
					result += "b"
				case p.Type == checkers.King && p.Color == checkers.White:
					result += "W"
				case p.Type == checkers.King && p.Color == checkers.Black:
					result += "B"
				}
			}
		}
		if row < 7 {
			result += "/"
		}
	}
	return result
}

// GetTurn returns "white" or "black".
func (a *CheckersAdapter) GetTurn() string {
	if a.game.Turn == checkers.White {
		return "white"
	}
	return "black"
}

// GetLegalMoves returns legal moves in WS format.
func (a *CheckersAdapter) GetLegalMoves() []WSMove {
	allMoves := a.game.AllLegalMoves()
	result := make([]WSMove, 0, len(allMoves))
	for _, m := range allMoves {
		wsMove := WSMove{
			From: fmt.Sprintf("%d,%d", m.From.Row, m.From.Col),
			To:   fmt.Sprintf("%d,%d", m.To.Row, m.To.Col),
		}
		if len(m.Captured) > 0 {
			wsMove.Captured = make([]WSPosition, len(m.Captured))
			for i, cap := range m.Captured {
				wsMove.Captured[i] = WSPosition{Row: cap.Row, Col: cap.Col}
			}
		}
		result = append(result, wsMove)
	}
	return result
}

// ApplyMove converts the WS move to engine format and applies it.
func (a *CheckersAdapter) ApplyMove(move WSMove) error {
	from := parseCheckersPos(move.From)
	to := parseCheckersPos(move.To)
	if !from.Valid() || !to.Valid() {
		return fmt.Errorf("invalid checkers position")
	}

	// Build the move - we need to find it in legal moves
	legalMoves := a.game.AllLegalMoves()
	for _, lm := range legalMoves {
		if lm.From == from && lm.To == to {
			return a.game.MakeMove(lm)
		}
	}

	// If no exact match found, just try From->To
	return a.game.MakeMove(checkers.Move{
		From:  from,
		To:    to,
		Piece: a.game.Board[from.Row][from.Col],
	})
}

// IsGameOver returns true if the game has ended.
func (a *CheckersAdapter) IsGameOver() bool {
	return a.game.State != checkers.Playing
}

// GetGameOverReason returns the reason the game ended.
func (a *CheckersAdapter) GetGameOverReason() string {
	switch a.game.State {
	case checkers.WhiteWin:
		return "checkmate" // using checkmate broadly for "win"
	case checkers.BlackWin:
		return "checkmate"
	case checkers.Draw:
		return "draw"
	}
	return ""
}

// GetWinner returns the winning color or empty for draw.
func (a *CheckersAdapter) GetWinner() string {
	switch a.game.State {
	case checkers.WhiteWin:
		return "white"
	case checkers.BlackWin:
		return "black"
	case checkers.Draw:
		return ""
	}
	return ""
}

// Resign resigns for the given color.
func (a *CheckersAdapter) Resign(color string) {
	// Checkers engine doesn't have a Resign method; we set state directly
	if color == "white" {
		a.game.State = checkers.BlackWin
	} else {
		a.game.State = checkers.WhiteWin
	}
}

// RollDice is a no-op for checkers.
func (a *CheckersAdapter) RollDice() []int {
	return nil
}

func parseCheckersPos(s string) checkers.Position {
	var row, col int
	if _, err := fmt.Sscanf(s, "%d,%d", &row, &col); err != nil {
		return checkers.Position{}
	}
	return checkers.Position{Row: row, Col: col}
}
