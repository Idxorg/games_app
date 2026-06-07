package backgammon

// Color represents a player color.
type Color int

const (
	White Color = iota
	Black
)

// Opposite returns the other color.
func (c Color) Opposite() Color {
	if c == White {
		return Black
	}
	return White
}

// String returns a string representation.
func (c Color) String() string {
	if c == White {
		return "white"
	}
	return "black"
}

// Point represents a backgammon point with checkers of one color.
type Point struct {
	Color Color
	Count int
}

// Special point constants for moves.
const (
	BarPoint = -1 // source: entering from the bar
	OffPoint = 24 // destination: bearing off
)

// Move represents a single checker move in backgammon.
type Move struct {
	From int // source point index (0-23) or BarPoint (-1)
	To   int // destination point index (0-23) or OffPoint (24)
	Die  int // die value used for this move
}

// GameState represents the current phase of the game.
type GameState int

const (
	Rolling  GameState = iota // waiting for dice roll
	Moving                     // player is making moves with rolled dice
	GameOver                   // game has ended
)

// String returns a human-readable state name.
func (s GameState) String() string {
	switch s {
	case Rolling:
		return "rolling"
	case Moving:
		return "moving"
	case GameOver:
		return "game_over"
	}
	return "unknown"
}

// Game represents a complete backgammon game with full state.
//
// Board indexing (0-based):
//   - Points 0-23 represent the 24 points on the board.
//   - White moves from high index to low (23→0), home board is 0-5.
//   - Black moves from low index to high (0→23), home board is 18-23.
//   - 1-indexed display: point i (0-based) = point i+1.
type Game struct {
	Board          [24]Point
	WhiteBar       int
	BlackBar       int
	WhiteOff       int
	BlackOff       int
	Turn           Color
	Dice           [2]int
	RemainingMoves []int
	State          GameState
	Moves          []Move
	Winner         Color
}

// TotalCheckers returns the total number of checkers for a color
// (board + bar + borne off). Should always be 15.
func (g *Game) TotalCheckers(color Color) int {
	count := 0
	if color == White {
		count = g.WhiteBar + g.WhiteOff
	} else {
		count = g.BlackBar + g.BlackOff
	}
	for i := 0; i < 24; i++ {
		if g.Board[i].Color == color {
			count += g.Board[i].Count
		}
	}
	return count
}
