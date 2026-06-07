package chess

import (
	"fmt"
	"strings"
)

// NewGame creates a new chess game with the standard starting position.
// Board convention: Board[0] = rank 1 (white back rank), Board[7] = rank 8 (black back rank).
// Position.Rank maps directly to Board index: rank 1→Board[0], rank 8→Board[7].
func NewGame() *Game {
	g := &Game{
		Turn:            White,
		State:           Playing,
		Moves:           make([]Move, 0),
		CastlingRights:  CastlingRights{WhiteKingside: true, WhiteQueenside: true, BlackKingside: true, BlackQueenside: true},
		HalfMoveClock:   0,
		FullMoveNumber:  1,
		PositionHistory: make(map[string]int),
	}

	// Set up starting position
	// White pieces (Board[0] = rank 1)
	g.Board[0] = [8]Piece{
		{Rook, White}, {Knight, White}, {Bishop, White}, {Queen, White},
		{King, White}, {Bishop, White}, {Knight, White}, {Rook, White},
	}
	for f := 0; f < 8; f++ {
		g.Board[1][f] = Piece{Pawn, White}
	}
	// Black pieces (Board[7] = rank 8)
	g.Board[6] = [8]Piece{
		{Pawn, Black}, {Pawn, Black}, {Pawn, Black}, {Pawn, Black},
		{Pawn, Black}, {Pawn, Black}, {Pawn, Black}, {Pawn, Black},
	}
	g.Board[7] = [8]Piece{
		{Rook, Black}, {Knight, Black}, {Bishop, Black}, {Queen, Black},
		{King, Black}, {Bishop, Black}, {Knight, Black}, {Rook, Black},
	}

	// Record initial position
	g.PositionHistory[g.ToFEN()] = 1

	return g
}

// IsCheck returns true if the given color's king is currently in check.
func (g *Game) IsCheck(color Color) bool {
	return g.isInCheck(color)
}

// IsCheckmate returns true if the given color is in checkmate.
func (g *Game) IsCheckmate(color Color) bool {
	if !g.isInCheck(color) {
		return false
	}
	savedTurn := g.Turn
	g.Turn = color
	moves := g.AllLegalMoves()
	g.Turn = savedTurn
	return len(moves) == 0
}

// IsStalemate returns true if the given color is in stalemate.
func (g *Game) IsStalemate(color Color) bool {
	if g.isInCheck(color) {
		return false
	}
	savedTurn := g.Turn
	g.Turn = color
	moves := g.AllLegalMoves()
	g.Turn = savedTurn
	return len(moves) == 0
}

// IsDraw returns true if the game is a draw by 50-move rule, insufficient material, or threefold repetition.
func (g *Game) IsDraw() bool {
	if g.HalfMoveClock >= 100 {
		return true
	}
	if g.PositionHistory[g.ToFEN()] >= 3 {
		return true
	}
	return g.isInsufficientMaterial()
}

// isInsufficientMaterial checks for insufficient mating material.
func (g *Game) isInsufficientMaterial() bool {
	whitePieces := make([]PieceType, 0)
	blackPieces := make([]PieceType, 0)
	for r := 0; r < 8; r++ {
		for f := 0; f < 8; f++ {
			p := g.Board[r][f]
			if p.Empty() || p.Type == King {
				continue
			}
			if p.Color == White {
				whitePieces = append(whitePieces, p.Type)
			} else {
				blackPieces = append(blackPieces, p.Type)
			}
		}
	}

	if len(whitePieces) == 0 && len(blackPieces) == 0 {
		return true
	}
	if len(whitePieces) == 0 && len(blackPieces) == 1 && (blackPieces[0] == Bishop || blackPieces[0] == Knight) {
		return true
	}
	if len(blackPieces) == 0 && len(whitePieces) == 1 && (whitePieces[0] == Bishop || whitePieces[0] == Knight) {
		return true
	}
	if len(whitePieces) == 1 && whitePieces[0] == Bishop && len(blackPieces) == 1 && blackPieces[0] == Bishop {
		var wBishopPos, bBishopPos Position
		for r := 0; r < 8; r++ {
			for f := 0; f < 8; f++ {
				p := g.Board[r][f]
				if p.Type == Bishop {
					if p.Color == White {
						wBishopPos = Position{f, r}
					} else {
						bBishopPos = Position{f, r}
					}
				}
			}
		}
		if (wBishopPos.File+wBishopPos.Rank)%2 == (bBishopPos.File+bBishopPos.Rank)%2 {
			return true
		}
	}

	return false
}

// MakeMove validates and executes a move on the board.
func (g *Game) MakeMove(from, to Position, promotion PieceType) error {
	move, err := g.ValidateMove(from, to, promotion)
	if err != nil {
		return err
	}
	return g.executeMove(move)
}

// MakeMoveDirect executes a pre-validated move directly.
func (g *Game) MakeMoveDirect(move Move) error {
	return g.executeMove(move)
}

// executeMove applies the move to the board and updates game state.
func (g *Game) executeMove(move Move) error {
	piece := move.Piece

	// Handle en passant capture
	if move.IsEnPassant {
		g.Board[move.From.Rank][move.To.File] = Piece{}
	}

	// Handle castling rook movement
	if move.IsCastling {
		if move.CastlingSide == "K" {
			rook := g.Board[move.From.Rank][7]
			g.Board[move.From.Rank][7] = Piece{}
			g.Board[move.From.Rank][5] = rook
		} else {
			rook := g.Board[move.From.Rank][0]
			g.Board[move.From.Rank][0] = Piece{}
			g.Board[move.From.Rank][3] = rook
		}
	}

	// Move the piece (with possible promotion)
	if move.Promotion != NoPiece {
		g.Board[move.To.Rank][move.To.File] = Piece{Type: move.Promotion, Color: piece.Color}
	} else {
		g.Board[move.To.Rank][move.To.File] = piece
	}
	g.Board[move.From.Rank][move.From.File] = Piece{}

	// Update castling rights
	g.updateCastlingRights(move)

	// Update en passant target
	g.EnPassantTarget = nil
	if piece.Type == Pawn {
		diff := move.To.Rank - move.From.Rank
		if diff == 2 || diff == -2 {
			epRank := (move.From.Rank + move.To.Rank) / 2
			g.EnPassantTarget = &Position{move.To.File, epRank}
		}
	}

	// Update half-move clock
	if piece.Type == Pawn || move.Captured != nil || move.IsEnPassant {
		g.HalfMoveClock = 0
	} else {
		g.HalfMoveClock++
	}

	// Add move to history
	g.Moves = append(g.Moves, move)

	// Switch turn
	if g.Turn == Black {
		g.FullMoveNumber++
	}
	g.Turn = g.Turn.Opposite()

	// Record position for threefold repetition
	fen := g.ToFEN()
	g.PositionHistory[fen]++

	// Update game state
	g.updateState()

	return nil
}

// updateCastlingRights updates castling rights based on the move.
func (g *Game) updateCastlingRights(move Move) {
	piece := move.Piece

	if piece.Type == King {
		if piece.Color == White {
			g.CastlingRights.WhiteKingside = false
			g.CastlingRights.WhiteQueenside = false
			g.WhiteKingMoved = true
		} else {
			g.CastlingRights.BlackKingside = false
			g.CastlingRights.BlackQueenside = false
			g.BlackKingMoved = true
		}
	}

	// White kingside rook at h1 = Position{7, 0}
	if move.From == (Position{7, 0}) || move.To == (Position{7, 0}) {
		g.CastlingRights.WhiteKingside = false
		if move.From == (Position{7, 0}) {
			g.WhiteRookMoved[1] = true
		}
	}
	// White queenside rook at a1 = Position{0, 0}
	if move.From == (Position{0, 0}) || move.To == (Position{0, 0}) {
		g.CastlingRights.WhiteQueenside = false
		if move.From == (Position{0, 0}) {
			g.WhiteRookMoved[0] = true
		}
	}
	// Black kingside rook at h8 = Position{7, 7}
	if move.From == (Position{7, 7}) || move.To == (Position{7, 7}) {
		g.CastlingRights.BlackKingside = false
		if move.From == (Position{7, 7}) {
			g.BlackRookMoved[1] = true
		}
	}
	// Black queenside rook at a8 = Position{0, 7}
	if move.From == (Position{0, 7}) || move.To == (Position{0, 7}) {
		g.CastlingRights.BlackQueenside = false
		if move.From == (Position{0, 7}) {
			g.BlackRookMoved[0] = true
		}
	}
}

// updateState checks the current game state after a move.
func (g *Game) updateState() {
	inCheck := g.isInCheck(g.Turn)
	allMoves := g.AllLegalMoves()

	if inCheck {
		if len(allMoves) == 0 {
			g.State = Checkmate
			return
		}
		g.State = Check
		return
	}

	if len(allMoves) == 0 {
		g.State = Stalemate
		return
	}

	if g.IsDraw() {
		g.State = Draw
		return
	}

	g.State = Playing
}

// Resign sets the game state to resigned for the given color.
func (g *Game) Resign(color Color) {
	g.State = Resigned
	_ = color
}

// ToFEN generates the FEN string for the current game state.
// Board[0]=rank 1, Board[7]=rank 8. FEN starts from rank 8.
func (g *Game) ToFEN() string {
	var sb strings.Builder

	// Piece placement: iterate from rank 8 (Board[7]) down to rank 1 (Board[0])
	for r := 7; r >= 0; r-- {
		empty := 0
		for f := 0; f < 8; f++ {
			p := g.Board[r][f]
			if p.Empty() {
				empty++
			} else {
				if empty > 0 {
					sb.WriteByte(byte('0' + empty))
					empty = 0
				}
				sb.WriteRune(p.FENChar())
			}
		}
		if empty > 0 {
			sb.WriteByte(byte('0' + empty))
		}
		if r > 0 {
			sb.WriteByte('/')
		}
	}

	// Active color
	sb.WriteByte(' ')
	sb.WriteString(g.Turn.String())

	// Castling
	sb.WriteByte(' ')
	castling := ""
	if g.CastlingRights.WhiteKingside {
		castling += "K"
	}
	if g.CastlingRights.WhiteQueenside {
		castling += "Q"
	}
	if g.CastlingRights.BlackKingside {
		castling += "k"
	}
	if g.CastlingRights.BlackQueenside {
		castling += "q"
	}
	if castling == "" {
		castling = "-"
	}
	sb.WriteString(castling)

	// En passant target
	sb.WriteByte(' ')
	if g.EnPassantTarget != nil {
		sb.WriteString(g.EnPassantTarget.Algebraic())
	} else {
		sb.WriteByte('-')
	}

	// Half-move clock and full move number
	fmt.Fprintf(&sb, " %d %d", g.HalfMoveClock, g.FullMoveNumber)

	return sb.String()
}

// FromFEN parses a FEN string and sets up the board state.
// FEN starts from rank 8 (index 0) down to rank 1 (index 7).
// Board[0]=rank 1, Board[7]=rank 8, so FEN rank index i maps to Board[7-i].
func (g *Game) FromFEN(fen string) error {
	parts := strings.Split(fen, " ")
	if len(parts) != 6 {
		return fmt.Errorf("invalid FEN: expected 6 fields, got %d", len(parts))
	}

	g.Board = Board{}
	g.PositionHistory = make(map[string]int)
	g.Moves = make([]Move, 0)
	g.State = Playing

	// Parse piece placement
	ranks := strings.Split(parts[0], "/")
	if len(ranks) != 8 {
		return fmt.Errorf("invalid FEN: expected 8 ranks, got %d", len(ranks))
	}

	for ri, rankStr := range ranks {
		r := 7 - ri // FEN ranks[0] = rank 8 → Board[7]
		f := 0
		for _, ch := range rankStr {
			if ch >= '1' && ch <= '8' {
				f += int(ch - '0')
			} else {
				if f > 7 {
					return fmt.Errorf("invalid FEN: file overflow in rank %d", r)
				}
				g.Board[r][f] = pieceFromFENChar(ch)
				f++
			}
		}
	}

	// Parse active color
	switch parts[1] {
	case "w":
		g.Turn = White
	case "b":
		g.Turn = Black
	default:
		return fmt.Errorf("invalid FEN: invalid active color '%s'", parts[1])
	}

	// Parse castling rights
	g.CastlingRights = CastlingRights{}
	if parts[2] != "-" {
		for _, ch := range parts[2] {
			switch ch {
			case 'K':
				g.CastlingRights.WhiteKingside = true
			case 'Q':
				g.CastlingRights.WhiteQueenside = true
			case 'k':
				g.CastlingRights.BlackKingside = true
			case 'q':
				g.CastlingRights.BlackQueenside = true
			}
		}
	}

	// Parse en passant target
	if parts[3] != "-" {
		g.EnPassantTarget = new(Position)
		*g.EnPassantTarget = FromAlgebraic(parts[3])
	} else {
		g.EnPassantTarget = nil
	}

	fmt.Sscanf(parts[4], "%d", &g.HalfMoveClock)
	fmt.Sscanf(parts[5], "%d", &g.FullMoveNumber)

	g.WhiteKingMoved = !g.CastlingRights.WhiteKingside && !g.CastlingRights.WhiteQueenside
	g.BlackKingMoved = !g.CastlingRights.BlackKingside && !g.CastlingRights.BlackQueenside

	g.updateState()
	g.PositionHistory[g.ToFEN()] = 1

	return nil
}

// pieceFromFENChar converts a FEN character to a Piece.
func pieceFromFENChar(ch rune) Piece {
	var pt PieceType
	var c Color
	switch ch {
	case 'P':
		pt, c = Pawn, White
	case 'N':
		pt, c = Knight, White
	case 'B':
		pt, c = Bishop, White
	case 'R':
		pt, c = Rook, White
	case 'Q':
		pt, c = Queen, White
	case 'K':
		pt, c = King, White
	case 'p':
		pt, c = Pawn, Black
	case 'n':
		pt, c = Knight, Black
	case 'b':
		pt, c = Bishop, Black
	case 'r':
		pt, c = Rook, Black
	case 'q':
		pt, c = Queen, Black
	case 'k':
		pt, c = King, Black
	default:
		return Piece{}
	}
	return Piece{Type: pt, Color: c}
}

// Clone creates a deep copy of the game state.
func (g *Game) Clone() *Game {
	clone := &Game{
		Board:           g.Board,
		Turn:            g.Turn,
		State:           g.State,
		CastlingRights:  g.CastlingRights,
		HalfMoveClock:   g.HalfMoveClock,
		FullMoveNumber:  g.FullMoveNumber,
		WhiteKingMoved:  g.WhiteKingMoved,
		BlackKingMoved:  g.BlackKingMoved,
		WhiteRookMoved:  g.WhiteRookMoved,
		BlackRookMoved:  g.BlackRookMoved,
	}
	if g.EnPassantTarget != nil {
		ep := *g.EnPassantTarget
		clone.EnPassantTarget = &ep
	}
	clone.Moves = make([]Move, len(g.Moves))
	copy(clone.Moves, g.Moves)
	clone.PositionHistory = make(map[string]int)
	for k, v := range g.PositionHistory {
		clone.PositionHistory[k] = v
	}
	return clone
}
