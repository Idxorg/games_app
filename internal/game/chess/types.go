package chess

// PieceType represents the type of chess piece.
type PieceType int

const (
	NoPiece PieceType = iota
	Pawn
	Knight
	Bishop
	Rook
	Queen
	King
)

// String returns the standard FEN character for the piece type.
func (p PieceType) String() string {
	switch p {
	case Pawn:
		return "p"
	case Knight:
		return "n"
	case Bishop:
		return "b"
	case Rook:
		return "r"
	case Queen:
		return "q"
	case King:
		return "k"
	}
	return ""
}

// Color represents the color of a chess piece.
type Color int

const (
	NoColor Color = iota
	White
	Black
)

// Opposite returns the other color.
func (c Color) Opposite() Color {
	if c == White {
		return Black
	}
	return White
}

// String returns "w" or "b".
func (c Color) String() string {
	if c == White {
		return "w"
	}
	return "b"
}

// Piece represents a chess piece with a type and color.
type Piece struct {
	Type  PieceType
	Color Color
}

// FENChar returns the FEN character for this piece.
func (p Piece) FENChar() rune {
	if p.Type == NoPiece {
		return '.'
	}
	c := p.Type.String()[0]
	if p.Color == White {
		return rune(c - 32) // uppercase
	}
	return rune(c)
}

// Empty returns true if this is not a piece.
func (p Piece) Empty() bool {
	return p.Type == NoPiece || p.Color == NoColor
}

// Position represents a square on the chess board.
// File 0-7 = a-h, Rank 0-7 = 1-8.
type Position struct {
	File int
	Rank int
}

// Algebraic returns the algebraic notation string (e.g., "e4").
func (p Position) Algebraic() string {
	if p.File < 0 || p.File > 7 || p.Rank < 0 || p.Rank > 7 {
		return "-"
	}
	return string(rune('a'+p.File)) + string(rune('1'+p.Rank))
}

// FromAlgebraic parses an algebraic notation string into a Position.
func FromAlgebraic(s string) Position {
	if len(s) != 2 {
		return Position{-1, -1}
	}
	f := int(s[0] - 'a')
	r := int(s[1] - '1')
	if f < 0 || f > 7 || r < 0 || r > 7 {
		return Position{-1, -1}
	}
	return Position{f, r}
}

// Valid returns true if the position is on the board.
func (p Position) Valid() bool {
	return p.File >= 0 && p.File <= 7 && p.Rank >= 0 && p.Rank <= 7
}

// Move represents a chess move.
type Move struct {
	From         Position
	To           Position
	Piece        Piece     // piece being moved
	Captured     *Piece    // piece captured (nil if none)
	Promotion    PieceType // promotion piece type (NoPiece if none)
	IsCastling   bool      // true if this is a castling move
	IsEnPassant  bool      // true if this is an en passant capture
	CastlingSide string    // "K" or "Q" for kingside/queenside
}

// GameState represents the current state of the game.
type GameState int

const (
	Playing GameState = iota
	Check
	Checkmate
	Stalemate
	Draw
	Resigned
	Timeout
)

// String returns a human-readable state.
func (s GameState) String() string {
	switch s {
	case Playing:
		return "playing"
	case Check:
		return "check"
	case Checkmate:
		return "checkmate"
	case Stalemate:
		return "stalemate"
	case Draw:
		return "draw"
	case Resigned:
		return "resigned"
	case Timeout:
		return "timeout"
	}
	return "unknown"
}

// CastlingRights tracks which castling moves are still available.
type CastlingRights struct {
	WhiteKingside  bool // white can castle kingside
	WhiteQueenside bool // white can castle queenside
	BlackKingside  bool // black can castle kingside
	BlackQueenside bool // black can castle queenside
}

// Board represents an 8x8 chess board.
type Board [8][8]Piece

// Game represents a chess game with full state.
type Game struct {
	Board            Board
	Turn             Color
	State            GameState
	Moves            []Move          // move history
	CastlingRights   CastlingRights
	EnPassantTarget  *Position       // en passant target square (nil if none)
	HalfMoveClock    int             // for 50-move rule
	FullMoveNumber   int             // increments after black moves
	PositionHistory  map[string]int  // for threefold repetition
	WhiteKingMoved   bool
	BlackKingMoved   bool
	WhiteRookMoved   [2]bool         // [0]=queenside (a1), [1]=kingside (h1)
	BlackRookMoved   [2]bool         // [0]=queenside (a8), [1]=kingside (h8)
}

// FENData holds the components of a FEN string.
type FENData struct {
	BoardState      string
	ActiveColor     string
	Castling        string
	EnPassant       string
	HalfMoveClock   int
	FullMoveNumber  int
}
