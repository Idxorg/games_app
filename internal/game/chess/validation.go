package chess

import "fmt"

// pseudoLegalMoves returns all pseudo-legal moves for a piece at pos (ignoring pins/check).
func (g *Game) pseudoLegalMoves(pos Position) []Move {
	piece := g.Board[pos.Rank][pos.File]
	if piece.Empty() {
		return nil
	}
	if piece.Color != g.Turn {
		return nil
	}

	var moves []Move
	switch piece.Type {
	case Pawn:
		moves = g.pawnMoves(pos, piece)
	case Knight:
		moves = g.knightMoves(pos, piece)
	case Bishop:
		moves = g.bishopMoves(pos, piece)
	case Rook:
		moves = g.rookMoves(pos, piece)
	case Queen:
		moves = g.queenMoves(pos, piece)
	case King:
		moves = g.kingMoves(pos, piece)
	}
	return moves
}

// pawnMoves generates all pseudo-legal pawn moves.
func (g *Game) pawnMoves(pos Position, piece Piece) []Move {
	var moves []Move
	dir := 1 // white moves up
	startRank := 1
	promoRank := 7
	if piece.Color == Black {
		dir = -1
		startRank = 6
		promoRank = 0
	}

	// Single push
	toRank := pos.Rank + dir
	if toRank >= 0 && toRank <= 7 {
		to := Position{pos.File, toRank}
		if g.Board[to.Rank][to.File].Empty() {
			if to.Rank == promoRank {
				moves = append(moves, g.promotionMoves(pos, to, piece)...)
			} else {
				moves = append(moves, Move{From: pos, To: to, Piece: piece})
			}
			// Double push from starting position
			if pos.Rank == startRank {
				to2 := Position{pos.File, pos.Rank + 2*dir}
				if g.Board[to2.Rank][to2.File].Empty() {
					moves = append(moves, Move{From: pos, To: to2, Piece: piece})
				}
			}
		}
	}

	// Captures (including en passant)
	for _, df := range []int{-1, 1} {
		cf := pos.File + df
		cr := pos.Rank + dir
		if cf < 0 || cf > 7 || cr < 0 || cr > 7 {
			continue
		}
		target := Position{cf, cr}
		targetPiece := g.Board[cr][cf]

		if !targetPiece.Empty() && targetPiece.Color != piece.Color {
			if cr == promoRank {
				moves = append(moves, g.promotionMoves(pos, target, piece)...)
			} else {
				moves = append(moves, Move{From: pos, To: target, Piece: piece, Captured: &targetPiece})
			}
		}

		// En passant
		if g.EnPassantTarget != nil && target == *g.EnPassantTarget {
			captured := g.Board[pos.Rank][cf]
			moves = append(moves, Move{
				From:        pos,
				To:          target,
				Piece:       piece,
				Captured:    &captured,
				IsEnPassant: true,
			})
		}
	}

	return moves
}

// promotionMoves generates promotion moves for all possible pieces.
func (g *Game) promotionMoves(from, to Position, piece Piece) []Move {
	var moves []Move
	for _, pt := range []PieceType{Queen, Rook, Bishop, Knight} {
		moves = append(moves, Move{
			From:      from,
			To:        to,
			Piece:     piece,
			Promotion: pt,
		})
	}
	return moves
}

// knightMoves generates all pseudo-legal knight moves.
func (g *Game) knightMoves(pos Position, piece Piece) []Move {
	var moves []Move
	offsets := []struct{ f, r int }{
		{-2, -1}, {-2, 1}, {-1, -2}, {-1, 2},
		{1, -2}, {1, 2}, {2, -1}, {2, 1},
	}
	for _, o := range offsets {
		f, r := pos.File+o.f, pos.Rank+o.r
		if f < 0 || f > 7 || r < 0 || r > 7 {
			continue
		}
		target := g.Board[r][f]
		if target.Empty() || target.Color != piece.Color {
			m := Move{From: pos, To: Position{f, r}, Piece: piece}
			if !target.Empty() {
				cp := target
				m.Captured = &cp
			}
			moves = append(moves, m)
		}
	}
	return moves
}

// slidingMoves generates moves along straight lines (for bishop, rook, queen).
func (g *Game) slidingMoves(pos Position, piece Piece, directions []struct{ f, r int }) []Move {
	var moves []Move
	for _, d := range directions {
		f, r := pos.File+d.f, pos.Rank+d.r
		for f >= 0 && f <= 7 && r >= 0 && r <= 7 {
			target := g.Board[r][f]
			if target.Empty() {
				moves = append(moves, Move{From: pos, To: Position{f, r}, Piece: piece})
			} else {
				if target.Color != piece.Color {
					cp := target
					moves = append(moves, Move{From: pos, To: Position{f, r}, Piece: piece, Captured: &cp})
				}
				break
			}
			f += d.f
			r += d.r
		}
	}
	return moves
}

// bishopMoves generates all pseudo-legal bishop moves.
func (g *Game) bishopMoves(pos Position, piece Piece) []Move {
	diagonals := []struct{ f, r int }{
		{-1, -1}, {-1, 1}, {1, -1}, {1, 1},
	}
	return g.slidingMoves(pos, piece, diagonals)
}

// rookMoves generates all pseudo-legal rook moves.
func (g *Game) rookMoves(pos Position, piece Piece) []Move {
	straights := []struct{ f, r int }{
		{-1, 0}, {1, 0}, {0, -1}, {0, 1},
	}
	return g.slidingMoves(pos, piece, straights)
}

// queenMoves generates all pseudo-legal queen moves.
func (g *Game) queenMoves(pos Position, piece Piece) []Move {
	allDirs := []struct{ f, r int }{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}
	return g.slidingMoves(pos, piece, allDirs)
}

// kingMoves generates all pseudo-legal king moves including castling.
func (g *Game) kingMoves(pos Position, piece Piece) []Move {
	var moves []Move
	// Normal king moves
	for dr := -1; dr <= 1; dr++ {
		for df := -1; df <= 1; df++ {
			if dr == 0 && df == 0 {
				continue
			}
			f, r := pos.File+df, pos.Rank+dr
			if f < 0 || f > 7 || r < 0 || r > 7 {
				continue
			}
			target := g.Board[r][f]
			if target.Empty() || target.Color != piece.Color {
				m := Move{From: pos, To: Position{f, r}, Piece: piece}
				if !target.Empty() {
					cp := target
					m.Captured = &cp
				}
				moves = append(moves, m)
			}
		}
	}

	// Castling
	moves = append(moves, g.castlingMoves(pos, piece)...)
	return moves
}

// castlingMoves generates castling moves if legal (pseudo-legal, not checking if king passes through check here).
func (g *Game) castlingMoves(pos Position, piece Piece) []Move {
	var moves []Move
	if piece.Type != King {
		return moves
	}

	if piece.Color == White {
		// Kingside: e1-g1, rook h1
		if g.CastlingRights.WhiteKingside && !g.WhiteKingMoved && !g.WhiteRookMoved[1] {
			if g.Board[0][5].Empty() && g.Board[0][6].Empty() {
				moves = append(moves, Move{
					From:         pos,
					To:           Position{6, 0},
					Piece:        piece,
					IsCastling:   true,
					CastlingSide: "K",
				})
			}
		}
		// Queenside: e1-c1, rook a1
		if g.CastlingRights.WhiteQueenside && !g.WhiteKingMoved && !g.WhiteRookMoved[0] {
			if g.Board[0][3].Empty() && g.Board[0][2].Empty() && g.Board[0][1].Empty() {
				moves = append(moves, Move{
					From:         pos,
					To:           Position{2, 0},
					Piece:        piece,
					IsCastling:   true,
					CastlingSide: "Q",
				})
			}
		}
	} else {
		// Kingside: e8-g8, rook h8
		if g.CastlingRights.BlackKingside && !g.BlackKingMoved && !g.BlackRookMoved[1] {
			if g.Board[7][5].Empty() && g.Board[7][6].Empty() {
				moves = append(moves, Move{
					From:         pos,
					To:           Position{6, 7},
					Piece:        piece,
					IsCastling:   true,
					CastlingSide: "K",
				})
			}
		}
		// Queenside: e8-c8, rook a8
		if g.CastlingRights.BlackQueenside && !g.BlackKingMoved && !g.BlackRookMoved[0] {
			if g.Board[7][3].Empty() && g.Board[7][2].Empty() && g.Board[7][1].Empty() {
				moves = append(moves, Move{
					From:         pos,
					To:           Position{2, 7},
					Piece:        piece,
					IsCastling:   true,
					CastlingSide: "Q",
				})
			}
		}
	}

	return moves
}

// isSquareAttacked returns true if the square at pos is attacked by any piece of the given color.
func (g *Game) isSquareAttacked(pos Position, byColor Color) bool {
	// Check knight attacks
	knightOffsets := []struct{ f, r int }{
		{-2, -1}, {-2, 1}, {-1, -2}, {-1, 2},
		{1, -2}, {1, 2}, {2, -1}, {2, 1},
	}
	for _, o := range knightOffsets {
		f, r := pos.File+o.f, pos.Rank+o.r
		if f >= 0 && f <= 7 && r >= 0 && r <= 7 {
			p := g.Board[r][f]
			if p.Type == Knight && p.Color == byColor {
				return true
			}
		}
	}

	// Check pawn attacks
	pawnDir := 1 // pawns that attack downward (for checking if white attacks)
	if byColor == White {
		pawnDir = -1 // white pawns attack from below
	}
	for _, df := range []int{-1, 1} {
		f := pos.File + df
		r := pos.Rank + pawnDir
		if f >= 0 && f <= 7 && r >= 0 && r <= 7 {
			p := g.Board[r][f]
			if p.Type == Pawn && p.Color == byColor {
				return true
			}
		}
	}

	// Check king attacks
	for dr := -1; dr <= 1; dr++ {
		for df := -1; df <= 1; df++ {
			if dr == 0 && df == 0 {
				continue
			}
			f, r := pos.File+df, pos.Rank+dr
			if f >= 0 && f <= 7 && r >= 0 && r <= 7 {
				p := g.Board[r][f]
				if p.Type == King && p.Color == byColor {
					return true
				}
			}
		}
	}

	// Check sliding pieces (bishop/rook/queen)
	// Diagonal directions (bishop, queen)
	diagonals := []struct{ f, r int }{{-1, -1}, {-1, 1}, {1, -1}, {1, 1}}
	for _, d := range diagonals {
		f, r := pos.File+d.f, pos.Rank+d.r
		for f >= 0 && f <= 7 && r >= 0 && r <= 7 {
			p := g.Board[r][f]
			if !p.Empty() {
				if p.Color == byColor && (p.Type == Bishop || p.Type == Queen) {
					return true
				}
				break
			}
			f += d.f
			r += d.r
		}
	}

	// Straight directions (rook, queen)
	straights := []struct{ f, r int }{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	for _, d := range straights {
		f, r := pos.File+d.f, pos.Rank+d.r
		for f >= 0 && f <= 7 && r >= 0 && r <= 7 {
			p := g.Board[r][f]
			if !p.Empty() {
				if p.Color == byColor && (p.Type == Rook || p.Type == Queen) {
					return true
				}
				break
			}
			f += d.f
			r += d.r
		}
	}

	return false
}

// findKing finds the position of the king of the given color.
func (g *Game) findKing(color Color) (Position, bool) {
	for r := 0; r < 8; r++ {
		for f := 0; f < 8; f++ {
			p := g.Board[r][f]
			if p.Type == King && p.Color == color {
				return Position{f, r}, true
			}
		}
	}
	return Position{}, false
}

// isInCheck returns true if the given color's king is in check.
func (g *Game) isInCheck(color Color) bool {
	kingPos, ok := g.findKing(color)
	if !ok {
		return false
	}
	return g.isSquareAttacked(kingPos, color.Opposite())
}

// wouldBeInCheck simulates a move and checks if the moving side's king would be in check.
func (g *Game) wouldBeInCheck(move Move) bool {
	// Save state
	origFrom := g.Board[move.From.Rank][move.From.File]
	origTo := g.Board[move.To.Rank][move.To.File]
	var origEP *Piece
	if move.IsEnPassant {
		epPiece := g.Board[move.From.Rank][move.To.File]
		origEP = &epPiece
		g.Board[move.From.Rank][move.To.File] = Piece{}
	}
	var origRookFrom, origRookTo Piece
	if move.IsCastling {
		if move.CastlingSide == "K" {
			origRookFrom = g.Board[move.From.Rank][7]
			origRookTo = g.Board[move.From.Rank][5]
			g.Board[move.From.Rank][7] = Piece{}
			g.Board[move.From.Rank][5] = origRookFrom
		} else {
			origRookFrom = g.Board[move.From.Rank][0]
			origRookTo = g.Board[move.From.Rank][3]
			g.Board[move.From.Rank][0] = Piece{}
			g.Board[move.From.Rank][3] = origRookFrom
		}
	}

	// Apply move
	piece := move.Piece
	if move.Promotion != NoPiece {
		piece = Piece{Type: move.Promotion, Color: move.Piece.Color}
	}
	g.Board[move.From.Rank][move.From.File] = Piece{}
	g.Board[move.To.Rank][move.To.File] = piece

	inCheck := g.isInCheck(move.Piece.Color)

	// Restore state
	g.Board[move.From.Rank][move.From.File] = origFrom
	g.Board[move.To.Rank][move.To.File] = origTo
	if move.IsEnPassant {
		g.Board[move.From.Rank][move.To.File] = *origEP
	}
	if move.IsCastling {
		g.Board[move.From.Rank][7] = origRookFrom
		g.Board[move.From.Rank][5] = origRookTo
		g.Board[move.From.Rank][0] = origRookFrom
		g.Board[move.From.Rank][3] = origRookTo
		if move.CastlingSide == "K" {
			g.Board[move.From.Rank][7] = origRookFrom
			g.Board[move.From.Rank][5] = origRookTo
		} else {
			g.Board[move.From.Rank][0] = origRookFrom
			g.Board[move.From.Rank][3] = origRookTo
		}
	}

	return inCheck
}

// validateCastling checks additional castling constraints (king not in check, doesn't pass through check).
func (g *Game) validateCastling(move Move) bool {
	kingPos := move.From
	color := move.Piece.Color
	opponent := color.Opposite()

	// King must not be in check
	if g.isSquareAttacked(kingPos, opponent) {
		return false
	}

	if move.CastlingSide == "K" {
		// King passes through f-file
		passThrough := Position{5, kingPos.Rank}
		// King lands on g-file (move.To already validated)
		if g.isSquareAttacked(passThrough, opponent) || g.isSquareAttacked(move.To, opponent) {
			return false
		}
	} else {
		// King passes through d-file
		passThrough := Position{3, kingPos.Rank}
		if g.isSquareAttacked(passThrough, opponent) || g.isSquareAttacked(move.To, opponent) {
			return false
		}
	}

	return true
}

// LegalMoves returns all legal moves for the piece at the given position.
func (g *Game) LegalMoves(pos Position) []Move {
	pseudoMoves := g.pseudoLegalMoves(pos)
	var legal []Move
	for _, m := range pseudoMoves {
		// For castling, check additional constraints
		if m.IsCastling {
			if !g.validateCastling(m) {
				continue
			}
		}
		// Check that the move doesn't leave own king in check
		if !g.wouldBeInCheck(m) {
			legal = append(legal, m)
		}
	}
	return legal
}

// AllLegalMoves returns all legal moves for the side to move.
func (g *Game) AllLegalMoves() []Move {
	var moves []Move
	for r := 0; r < 8; r++ {
		for f := 0; f < 8; f++ {
			p := g.Board[r][f]
			if !p.Empty() && p.Color == g.Turn {
				pos := Position{f, r}
				moves = append(moves, g.LegalMoves(pos)...)
			}
		}
	}
	return moves
}

// moveError represents errors during move execution.
type moveError struct {
	msg string
}

func (e *moveError) Error() string {
	return e.msg
}

// ValidateMove checks if a move is legal and returns a fully populated Move or an error.
func (g *Game) ValidateMove(from, to Position, promotion PieceType) (Move, error) {
	if !from.Valid() || !to.Valid() {
		return Move{}, fmt.Errorf("invalid position")
	}
	if g.State == Checkmate || g.State == Stalemate || g.State == Draw || g.State == Resigned || g.State == Timeout {
		return Move{}, fmt.Errorf("game is over: %s", g.State)
	}

	piece := g.Board[from.Rank][from.File]
	if piece.Empty() {
		return Move{}, fmt.Errorf("no piece at %s", from.Algebraic())
	}
	if piece.Color != g.Turn {
		return Move{}, fmt.Errorf("not your turn")
	}

	legal := g.LegalMoves(from)
	for _, m := range legal {
		if m.To == to {
			// If promotion is specified, match it
			if promotion != NoPiece {
				if m.Promotion == promotion {
					return m, nil
				}
				continue
			}
			// If no promotion specified but move requires promotion
			if m.Promotion != NoPiece {
				// Default to queen promotion
				for _, m2 := range legal {
					if m2.To == to && m2.Promotion == Queen {
						return m2, nil
					}
				}
				return Move{}, fmt.Errorf("promotion required")
			}
			return m, nil
		}
	}

	return Move{}, fmt.Errorf("illegal move from %s to %s", from.Algebraic(), to.Algebraic())
}
