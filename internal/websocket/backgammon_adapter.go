package websocket

import (
	"fmt"

	"game-platform/internal/game/backgammon"
)

// BackgammonAdapter adapts the backgammon engine to the GameEngineAdapter interface.
type BackgammonAdapter struct {
	game *backgammon.Game
}

// NewBackgammonAdapter creates a new backgammon adapter.
func NewBackgammonAdapter() *BackgammonAdapter {
	return &BackgammonAdapter{
		game: backgammon.NewGame(),
	}
}

// GetBoard returns the board state as a string representation.
// Format: point0:color,count;point1:color,count;...;bar:white,b;black,b;off:white,b;black,b
func (a *BackgammonAdapter) GetBoard() string {
	var result string
	for i := 0; i < 24; i++ {
		if i > 0 {
			result += ";"
		}
		pt := a.game.Board[i]
		color := "."
		if pt.Count > 0 {
			color = pt.Color.String()
		}
		result += fmt.Sprintf("%d:%s,%d", i, color, pt.Count)
	}
	result += fmt.Sprintf(";bar:%d,%d", a.game.WhiteBar, a.game.BlackBar)
	result += fmt.Sprintf(";off:%d,%d", a.game.WhiteOff, a.game.BlackOff)
	return result
}

// GetTurn returns "white" or "black".
func (a *BackgammonAdapter) GetTurn() string {
	if a.game.Turn == backgammon.White {
		return "white"
	}
	return "black"
}

// GetLegalMoves returns legal moves in WS format.
func (a *BackgammonAdapter) GetLegalMoves() []WSMove {
	if a.game.State != backgammon.Moving {
		return nil
	}
	allMoves := a.game.AllLegalMoves()
	result := make([]WSMove, 0, len(allMoves))
	for _, m := range allMoves {
		result = append(result, WSMove{
			FromIdx: m.From,
			ToIdx:   m.To,
			Die:     m.Die,
			From:    pointToString(m.From),
			To:      pointToString(m.To),
		})
	}
	return result
}

// ApplyMove converts the WS move to engine format and applies it.
func (a *BackgammonAdapter) ApplyMove(move WSMove) error {
	if a.game.State != backgammon.Moving {
		return fmt.Errorf("must roll dice first")
	}

	from := move.FromIdx
	to := move.ToIdx
	die := move.Die

	if die <= 0 {
		return fmt.Errorf("die value required for backgammon move")
	}

	return a.game.MakeMove(backgammon.Move{
		From: from,
		To:   to,
		Die:  die,
	})
}

// IsGameOver returns true if the game has ended.
func (a *BackgammonAdapter) IsGameOver() bool {
	return a.game.State == backgammon.GameOver
}

// GetGameOverReason returns the reason the game ended.
func (a *BackgammonAdapter) GetGameOverReason() string {
	if a.game.State == backgammon.GameOver {
		return "checkmate" // using checkmate broadly for "win"
	}
	return ""
}

// GetWinner returns the winning color or empty for draw.
func (a *BackgammonAdapter) GetWinner() string {
	if a.game.State == backgammon.GameOver {
		return a.game.Winner.String()
	}
	return ""
}

// Resign resigns for the given color.
func (a *BackgammonAdapter) Resign(color string) {
	a.game.State = backgammon.GameOver
	if color == "white" {
		a.game.Winner = backgammon.Black
	} else {
		a.game.Winner = backgammon.White
	}
}

// RollDice rolls the dice server-side.
func (a *BackgammonAdapter) RollDice() []int {
	dice := a.game.RollDice()
	return []int{dice[0], dice[1]}
}

// GetDice returns the current dice values.
func (a *BackgammonAdapter) GetDice() [2]int {
	return a.game.Dice
}

// GetRemainingMoves returns remaining dice moves.
func (a *BackgammonAdapter) GetRemainingMoves() []int {
	return a.game.RemainingMoves
}

// GetState returns the current game state.
func (a *BackgammonAdapter) GetState() backgammon.GameState {
	return a.game.State
}

func pointToString(p int) string {
	switch p {
	case backgammon.BarPoint:
		return "bar"
	case backgammon.OffPoint:
		return "off"
	default:
		return fmt.Sprintf("%d", p)
	}
}
