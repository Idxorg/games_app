package checkers

import (
	"fmt"
	"strings"
)

// NewGame creates a new Russian checkers game with the standard starting position.
// White moves first. White men start on rows 5,6,7 (bottom from white's perspective).
// Black men start on rows 0,1,2 (top from white's perspective).
// Only dark squares (where (row+col)%2 == 1) are used.
func NewGame() *Game {
	g := &Game{
		Turn:  White,
		State: Playing,
		Moves: make([]Move, 0),
	}

	// Place black men on rows 0, 1, 2
	for row := 0; row < 3; row++ {
		for col := 0; col < 8; col++ {
			if (row+col)%2 == 1 {
				g.Board[row][col] = Piece{Man, Black}
			}
		}
	}

	// Place white men on rows 5, 6, 7
	for row := 5; row < 8; row++ {
		for col := 0; col < 8; col++ {
			if (row+col)%2 == 1 {
				g.Board[row][col] = Piece{Man, White}
			}
		}
	}

	return g
}

// diagonal directions for sliding
var diagDirs = [][2]int{
	{-1, -1}, {-1, 1},
	{1, -1}, {1, 1},
}

// AllLegalMoves returns all legal moves for the current player.
// If any capture exists, only capture moves are returned (mandatory capture).
func (g *Game) AllLegalMoves() []Move {
	captures := g.allCaptures()
	if len(captures) > 0 {
		return captures
	}
	return g.allSimpleMoves()
}

// IsForcedCapture returns true if the current player must capture.
func (g *Game) IsForcedCapture() bool {
	return len(g.allCaptures()) > 0
}

// LegalMoves returns all legal moves for a specific piece position.
// If forced capture is in effect, only capture moves for that piece are returned.
func (g *Game) LegalMoves(pos Position) []Move {
	piece := g.Board[pos.Row][pos.Col]
	if piece.Empty() || piece.Color != g.Turn {
		return nil
	}

	captures := g.capturesFrom(pos)
	if len(captures) > 0 {
		return captures
	}

	// If there are captures elsewhere, this piece can't make simple moves
	if g.IsForcedCapture() {
		return nil
	}

	return g.simpleMovesFrom(pos)
}

// simpleMovesFrom returns non-capture moves for a piece at pos.
func (g *Game) simpleMovesFrom(pos Position) []Move {
	piece := g.Board[pos.Row][pos.Col]
	if piece.Empty() {
		return nil
	}

	if piece.Type == Man {
		return g.manSimpleMoves(pos, piece)
	}
	return g.kingSimpleMoves(pos, piece)
}

// manSimpleMoves returns forward diagonal moves for a man.
func (g *Game) manSimpleMoves(pos Position, piece Piece) []Move {
	var moves []Move
	var rowDirs []int
	if piece.Color == White {
		rowDirs = []int{-1} // white moves up (decreasing row)
	} else {
		rowDirs = []int{1} // black moves down (increasing row)
	}

	for _, dr := range rowDirs {
		for _, dir := range diagDirs {
			if dir[0] != dr {
				continue
			}
			nr := pos.Row + dr
			nc := pos.Col + dir[1]
			np := Position{nr, nc}
			if np.Valid() && g.Board[nr][nc].Empty() {
				promoted := isPromotionRow(np, piece.Color)
				moves = append(moves, Move{From: pos, To: np, Piece: piece, Promoted: promoted})
			}
		}
	}
	return moves
}

// kingSimpleMoves returns all diagonal slide moves for a king.
func (g *Game) kingSimpleMoves(pos Position, piece Piece) []Move {
	var moves []Move
	for _, dir := range diagDirs {
		for dist := 1; dist < 8; dist++ {
			nr := pos.Row + dir[0]*dist
			nc := pos.Col + dir[1]*dist
			np := Position{nr, nc}
			if !np.Valid() {
				break
			}
			if g.Board[nr][nc].Empty() {
				moves = append(moves, Move{From: pos, To: np, Piece: piece})
			} else {
				break // blocked
			}
		}
	}
	return moves
}

// allCaptures returns all capture moves available to the current player.
func (g *Game) allCaptures() []Move {
	var allMoves []Move
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			pos := Position{row, col}
			piece := g.Board[row][col]
			if !piece.Empty() && piece.Color == g.Turn {
				moves := g.capturesFrom(pos)
				allMoves = append(allMoves, moves...)
			}
		}
	}
	return allMoves
}

// allSimpleMoves returns all non-capture moves for the current player.
func (g *Game) allSimpleMoves() []Move {
	var allMoves []Move
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			pos := Position{row, col}
			piece := g.Board[row][col]
			if !piece.Empty() && piece.Color == g.Turn {
				moves := g.simpleMovesFrom(pos)
				allMoves = append(allMoves, moves...)
			}
		}
	}
	return allMoves
}

// capturesFrom returns all capture sequences starting from pos.
// Uses recursive multi-jump: finds all maximal capture sequences.
func (g *Game) capturesFrom(pos Position) []Move {
	piece := g.Board[pos.Row][pos.Col]
	if piece.Empty() {
		return nil
	}
	return g.findCaptureSequences(pos, pos, piece, nil, nil, false)
}

// findCaptureSequences recursively finds all capture sequences.
// origin is the original starting position of the piece.
// visited tracks positions of pieces already captured in this chain.
func (g *Game) findCaptureSequences(origin, pos Position, piece Piece, captured []Position, visited []Position, promoted bool) []Move {
	var result []Move

	if piece.Type == King {
		result = g.kingCaptureSequences(origin, pos, piece, captured, visited, promoted)
	} else {
		result = g.manCaptureSequences(origin, pos, piece, captured, visited, promoted)
	}

	return result
}

// manCaptureSequences finds capture sequences for a man.
// Men capture 1 square diagonally in any direction (Russian checkers).
func (g *Game) manCaptureSequences(origin, pos Position, piece Piece, captured []Position, visited []Position, promoted bool) []Move {
	var result []Move

	for _, dir := range diagDirs {
		mr := pos.Row + dir[0]
		mc := pos.Col + dir[1]
		mp := Position{mr, mc}
		if !mp.Valid() {
			continue
		}
		target := g.Board[mr][mc]
		if target.Empty() || target.Color == piece.Color {
			continue
		}

		// Can't recapture same piece in this chain
		if isVisited(visited, mp) {
			continue
		}

		lr := pos.Row + dir[0]*2
		lc := pos.Col + dir[1]*2
		lp := Position{lr, lc}
		if !lp.Valid() || !g.Board[lr][lc].Empty() {
			continue
		}

		// Check for promotion
		newPromoted := promoted
		newPiece := piece
		if piece.Type == Man && !promoted && isPromotionRow(lp, piece.Color) {
			newPiece = Piece{King, piece.Color}
			newPromoted = true
		}

		newCaptured := append(copyPositions(captured), mp)
		newVisited := append(copyPositions(visited), mp)

		// Temporarily apply move on board
		origFrom := g.Board[pos.Row][pos.Col]
		origMid := g.Board[mr][mc]
		origTo := g.Board[lr][lc]
		g.Board[pos.Row][pos.Col] = Piece{Type: None}
		g.Board[mr][mc] = Piece{Type: None}
		g.Board[lr][lc] = newPiece

		// Recurse for further captures
		subMoves := g.findCaptureSequences(origin, lp, newPiece, newCaptured, newVisited, newPromoted)

		// Restore
		g.Board[pos.Row][pos.Col] = origFrom
		g.Board[mr][mc] = origMid
		g.Board[lr][lc] = origTo

		if len(subMoves) == 0 {
			result = append(result, Move{
				From:     origin,
				To:       lp,
				Piece:    piece,
				Captured: newCaptured,
				IsJump:   true,
				Promoted: newPromoted,
			})
		} else {
			result = append(result, subMoves...)
		}
	}

	return result
}

// kingCaptureSequences finds capture sequences for a king.
// Kings capture at any distance along a diagonal, landing on any empty square beyond the enemy.
func (g *Game) kingCaptureSequences(origin, pos Position, piece Piece, captured []Position, visited []Position, promoted bool) []Move {
	var result []Move

	for _, dir := range diagDirs {
		for dist := 1; dist < 8; dist++ {
			mr := pos.Row + dir[0]*dist
			mc := pos.Col + dir[1]*dist
			mp := Position{mr, mc}
			if !mp.Valid() {
				break
			}
			target := g.Board[mr][mc]
			if target.Empty() {
				continue // empty, keep sliding
			}
			if target.Color == piece.Color {
				break // blocked by own piece
			}

			// Found enemy piece; check if already captured
			if isVisited(visited, mp) {
				break
			}

			// Check all landing squares beyond the enemy
			for landDist := dist + 1; landDist < 8; landDist++ {
				lr := pos.Row + dir[0]*landDist
				lc := pos.Col + dir[1]*landDist
				lp := Position{lr, lc}
				if !lp.Valid() {
					break
				}
				if !g.Board[lr][lc].Empty() {
					break // blocked after enemy
				}

				newCaptured := append(copyPositions(captured), mp)
				newVisited := append(copyPositions(visited), mp)

				// Temporarily apply move
				origFrom := g.Board[pos.Row][pos.Col]
				origMid := g.Board[mr][mc]
				origTo := g.Board[lr][lc]
				g.Board[pos.Row][pos.Col] = Piece{Type: None}
				g.Board[mr][mc] = Piece{Type: None}
				g.Board[lr][lc] = piece

				// Recurse for further captures
				subMoves := g.findCaptureSequences(origin, lp, piece, newCaptured, newVisited, promoted)

				// Restore
				g.Board[pos.Row][pos.Col] = origFrom
				g.Board[mr][mc] = origMid
				g.Board[lr][lc] = origTo

				if len(subMoves) == 0 {
					result = append(result, Move{
						From:     origin,
						To:       lp,
						Piece:    piece,
						Captured: newCaptured,
						IsJump:   true,
						Promoted: false,
					})
				} else {
					result = append(result, subMoves...)
				}
			}
			break // after processing enemy at this distance, can't skip over
		}
	}

	return result
}

// isPromotionRow returns true if the position is the promotion row for the given color.
func isPromotionRow(pos Position, color Color) bool {
	if color == White {
		return pos.Row == 0
	}
	return pos.Row == 7
}

func isVisited(visited []Position, p Position) bool {
	for _, v := range visited {
		if v == p {
			return true
		}
	}
	return false
}

func copyPositions(positions []Position) []Position {
	cp := make([]Position, len(positions))
	copy(cp, positions)
	return cp
}

// GetCaptureSequences returns all capture sequences for a piece at pos.
// Each sequence is a slice containing a single Move (representing the full chain).
func (g *Game) GetCaptureSequences(pos Position) [][]Move {
	moves := g.capturesFrom(pos)
	if len(moves) == 0 {
		return nil
	}
	sequences := make([][]Move, len(moves))
	for i, m := range moves {
		sequences[i] = []Move{m}
	}
	return sequences
}

// MakeMove validates and executes a move.
func (g *Game) MakeMove(move Move) error {
	if g.State != Playing {
		return fmt.Errorf("game is over: %s", g.State)
	}

	piece := g.Board[move.From.Row][move.From.Col]
	if piece.Empty() || piece.Color != g.Turn {
		return fmt.Errorf("no valid piece at %s for %s", move.From, g.Turn)
	}

	// Validate move is legal
	legal := g.AllLegalMoves()
	found := false
	for _, lm := range legal {
		if movesEqual(lm, move) {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("illegal move from %s to %s", move.From, move.To)
	}

	// Execute the move
	g.Board[move.To.Row][move.To.Col] = piece
	g.Board[move.From.Row][move.From.Col] = Piece{Type: None}

	// Remove captured pieces
	for _, cap := range move.Captured {
		g.Board[cap.Row][cap.Col] = Piece{Type: None}
	}

	// Handle promotion
	finalPiece := piece
	if move.Promoted {
		finalPiece = Piece{King, piece.Color}
		g.Board[move.To.Row][move.To.Col] = finalPiece
	}

	// Track king moves for draw detection
	if finalPiece.Type == King {
		if len(move.Captured) == 0 {
			g.KingMoveCount++
		} else {
			g.KingMoveCount = 0
		}
		g.TotalKingMoves++
	}

	g.Moves = append(g.Moves, move)

	// Switch turn
	g.Turn = opponent(g.Turn)

	// Check for draw
	if g.checkDraw() {
		g.State = Draw
		return nil
	}

	// Check for win/loss
	if g.hasNoPieces(g.Turn) || g.hasNoMoves(g.Turn) {
		if g.Turn == White {
			g.State = BlackWin
		} else {
			g.State = WhiteWin
		}
	}

	return nil
}

// movesEqual checks if two moves are equivalent.
func movesEqual(a, b Move) bool {
	if a.From != b.From || a.To != b.To {
		return false
	}
	if len(a.Captured) != len(b.Captured) {
		return false
	}
	for i := range a.Captured {
		if a.Captured[i] != b.Captured[i] {
			return false
		}
	}
	return true
}

// opponent returns the other color.
func opponent(c Color) Color {
	if c == White {
		return Black
	}
	return White
}

// hasNoPieces returns true if the given color has no pieces on the board.
func (g *Game) hasNoPieces(color Color) bool {
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			p := g.Board[row][col]
			if !p.Empty() && p.Color == color {
				return false
			}
		}
	}
	return true
}

// hasNoMoves returns true if the given color has no legal moves.
func (g *Game) hasNoMoves(color Color) bool {
	savedTurn := g.Turn
	g.Turn = color
	moves := g.AllLegalMoves()
	g.Turn = savedTurn
	return len(moves) == 0
}

// checkDraw checks draw conditions per Russian checkers rules.
func (g *Game) checkDraw() bool {
	// 15 consecutive king moves without capture
	if g.KingMoveCount >= 15 {
		return true
	}
	// 25 total king moves if both sides have only kings
	if g.TotalKingMoves >= 25 && g.bothSidesKingsOnly() {
		return true
	}
	return false
}

// bothSidesKingsOnly returns true if all remaining pieces are kings.
func (g *Game) bothSidesKingsOnly() bool {
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			p := g.Board[row][col]
			if !p.Empty() && p.Type == Man {
				return false
			}
		}
	}
	return true
}

// CountPieces returns the number of pieces for each color.
func (g *Game) CountPieces() (white int, black int) {
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			p := g.Board[row][col]
			if !p.Empty() {
				if p.Color == White {
					white++
				} else {
					black++
				}
			}
		}
	}
	return
}

// MandatoryCapture enforces mandatory capture rule.
func (g *Game) MandatoryCapture() bool {
	return g.IsForcedCapture()
}

// String returns a string representation of the board.
func (g *Game) String() string {
	var sb strings.Builder
	sb.WriteString("  ")
	for col := 0; col < 8; col++ {
		sb.WriteString(fmt.Sprintf("%2d", col))
	}
	sb.WriteString("\n")
	for row := 0; row < 8; row++ {
		sb.WriteString(fmt.Sprintf("%2d", row))
		for col := 0; col < 8; col++ {
			p := g.Board[row][col]
			if p.Empty() {
				sb.WriteString(" .")
			} else {
				switch {
				case p.Type == Man && p.Color == White:
					sb.WriteString(" w")
				case p.Type == Man && p.Color == Black:
					sb.WriteString(" b")
				case p.Type == King && p.Color == White:
					sb.WriteString(" W")
				case p.Type == King && p.Color == Black:
					sb.WriteString(" B")
				}
			}
		}
		sb.WriteString("\n")
	}
	return sb.String()
}
