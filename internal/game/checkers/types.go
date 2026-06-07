package checkers

import "fmt"

// PieceType represents the type of a checkers piece.
type PieceType int

const (
	None PieceType = iota // no piece (empty square)
	Man
	King
)

// Color represents the color of a checkers piece.
type Color int

const (
	White Color = iota
	Black
)

func (c Color) String() string {
	if c == White {
		return "White"
	}
	return "Black"
}

// Piece represents a single piece on the board.
type Piece struct {
	Type  PieceType
	Color Color
}

func (p Piece) String() string {
	if p.Type == None {
		return "Empty"
	}
	if p.Type == King {
		return "King(" + p.Color.String() + ")"
	}
	return "Man(" + p.Color.String() + ")"
}

// Empty returns true if the square has no piece.
func (p Piece) Empty() bool {
	return p.Type == None
}

// Position represents a square on the 8x8 board.
type Position struct {
	Row int // 0..7 (0 = top of board from display perspective)
	Col int // 0..7
}

// Position on a checkers board is valid if within bounds and on a dark square.
// Convention: dark squares have (row+col) % 2 == 1.
func (p Position) Valid() bool {
	return p.Row >= 0 && p.Row < 8 && p.Col >= 0 && p.Col < 8 && (p.Row+p.Col)%2 == 1
}

func (p Position) String() string {
	return fmt.Sprintf("(%d,%d)", p.Row, p.Col)
}

// Board represents an 8x8 checkers board.
type Board [8][8]Piece

// Game represents a checkers game state.
type Game struct {
	Board          Board
	Turn           Color
	State          GameState
	Moves          []Move // history of moves played
	KingMoveCount  int    // consecutive king moves without capture (for draw rule)
	TotalKingMoves int    // total king moves (for draw rule)
}

// GameState represents the current state of the game.
type GameState int

const (
	Playing GameState = iota
	Draw
	WhiteWin
	BlackWin
)

func (s GameState) String() string {
	switch s {
	case Playing:
		return "Playing"
	case Draw:
		return "Draw"
	case WhiteWin:
		return "WhiteWin"
	case BlackWin:
		return "BlackWin"
	}
	return "Unknown"
}

// Move represents a single move in checkers.
// For multi-jump sequences, a Move captures multiple pieces in one turn.
type Move struct {
	From     Position   // starting position
	To       Position   // final position
	Piece    Piece      // piece that moved
	Captured []Position // positions of captured pieces (may be multiple for multi-jump)
	IsJump   bool       // true if this move involves at least one capture
	Promoted bool       // true if a man was promoted to king during this move
}
