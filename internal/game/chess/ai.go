package chess

import (
	"math/rand"
	"time"
)

// Piece values for evaluation.
const (
	PawnValue   = 100
	KnightValue = 320
	BishopValue = 330
	RookValue   = 500
	QueenValue  = 900
	KingValue   = 20000

	evaluationInfinity = 100000
)

// pieceValue returns the material value of a piece type.
func pieceValue(pt PieceType) int {
	switch pt {
	case Pawn:
		return PawnValue
	case Knight:
		return KnightValue
	case Bishop:
		return BishopValue
	case Rook:
		return RookValue
	case Queen:
		return QueenValue
	case King:
		return KingValue
	default:
		return 0
	}
}

// Position tables for each piece type (from white's perspective).
// Index 0 = rank 8 (far side for white), index 7 = rank 1 (starting side).

var pawnTable = [8][8]int{
	{0, 0, 0, 0, 0, 0, 0, 0},
	{50, 50, 50, 50, 50, 50, 50, 50},
	{10, 10, 20, 30, 30, 20, 10, 10},
	{5, 5, 10, 25, 25, 10, 5, 5},
	{0, 0, 0, 20, 20, 0, 0, 0},
	{5, -5, -10, 0, 0, -10, -5, 5},
	{5, 10, 10, -20, -20, 10, 10, 5},
	{0, 0, 0, 0, 0, 0, 0, 0},
}

var knightTable = [8][8]int{
	{-50, -40, -30, -30, -30, -30, -40, -50},
	{-40, -20, 0, 0, 0, 0, -20, -40},
	{-30, 0, 10, 15, 15, 10, 0, -30},
	{-30, 5, 15, 20, 20, 15, 5, -30},
	{-30, 0, 15, 20, 20, 15, 0, -30},
	{-30, 5, 10, 15, 15, 10, 5, -30},
	{-40, -20, 0, 5, 5, 0, -20, -40},
	{-50, -40, -30, -30, -30, -30, -40, -50},
}

var bishopTable = [8][8]int{
	{-20, -10, -10, -10, -10, -10, -10, -20},
	{-10, 0, 0, 0, 0, 0, 0, -10},
	{-10, 0, 5, 10, 10, 5, 0, -10},
	{-10, 5, 5, 10, 10, 5, 5, -10},
	{-10, 0, 10, 10, 10, 10, 0, -10},
	{-10, 10, 10, 10, 10, 10, 10, -10},
	{-10, 5, 0, 0, 0, 0, 5, -10},
	{-20, -10, -10, -10, -10, -10, -10, -20},
}

var rookTable = [8][8]int{
	{0, 0, 0, 0, 0, 0, 0, 0},
	{5, 10, 10, 10, 10, 10, 10, 5},
	{-5, 0, 0, 0, 0, 0, 0, -5},
	{-5, 0, 0, 0, 0, 0, 0, -5},
	{-5, 0, 0, 0, 0, 0, 0, -5},
	{-5, 0, 0, 0, 0, 0, 0, -5},
	{-5, 0, 0, 0, 0, 0, 0, -5},
	{0, 0, 0, 5, 5, 0, 0, 0},
}

var queenTable = [8][8]int{
	{-20, -10, -10, -5, -5, -10, -10, -20},
	{-10, 0, 0, 0, 0, 0, 0, -10},
	{-10, 0, 5, 5, 5, 5, 0, -10},
	{-5, 0, 5, 5, 5, 5, 0, -5},
	{0, 0, 5, 5, 5, 5, 0, -5},
	{-10, 5, 5, 5, 5, 5, 0, -10},
	{-10, 0, 5, 0, 0, 0, 0, -10},
	{-20, -10, -10, -5, -5, -10, -10, -20},
}

var kingMiddleTable = [8][8]int{
	{-30, -40, -40, -50, -50, -40, -40, -30},
	{-30, -40, -40, -50, -50, -40, -40, -30},
	{-30, -40, -40, -50, -50, -40, -40, -30},
	{-30, -40, -40, -50, -50, -40, -40, -30},
	{-20, -30, -30, -40, -40, -30, -30, -20},
	{-10, -20, -20, -20, -20, -20, -20, -10},
	{20, 20, 0, 0, 0, 0, 20, 20},
	{20, 30, 10, 0, 0, 10, 30, 20},
}

// pieceSquareValue returns the positional bonus for a piece at a given position.
func pieceSquareValue(pt PieceType, color Color, file, rank int) int {
	var table *[8][8]int
	switch pt {
	case Pawn:
		table = &pawnTable
	case Knight:
		table = &knightTable
	case Bishop:
		table = &bishopTable
	case Rook:
		table = &rookTable
	case Queen:
		table = &queenTable
	case King:
		table = &kingMiddleTable
	default:
		return 0
	}
	// Tables are written from white's perspective with index 0 = rank 8 (far side).
	// Board[0] = rank 1, so Position rank 0 = rank 1 → table index 7.
	if color == White {
		return table[7-rank][file]
	}
	// For black, mirror: rank 1 (from black's view, the far side) → table index 0
	return table[rank][file]
}

// Evaluate returns the evaluation score for the current position.
// Always from white's perspective (positive = white better).
func (g *Game) Evaluate() int {
	score := 0
	for r := 0; r < 8; r++ {
		for f := 0; f < 8; f++ {
			p := g.Board[r][f]
			if p.Empty() {
				continue
			}
			val := pieceValue(p.Type) + pieceSquareValue(p.Type, p.Color, f, r)
			if p.Color == White {
				score += val
			} else {
				score -= val
			}
		}
	}

	if g.State == Checkmate {
		if g.Turn == White {
			score = -KingValue
		} else {
			score = KingValue
		}
	}

	return score
}

// evaluateForSide returns evaluation from the perspective of the side to move.
func (g *Game) evaluateForSide() int {
	score := g.Evaluate()
	if g.Turn == Black {
		return -score
	}
	return score
}

// RandomMove returns a random legal move.
func (g *Game) RandomMove() Move {
	moves := g.AllLegalMoves()
	if len(moves) == 0 {
		return Move{}
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return moves[rng.Intn(len(moves))]
}

// BestMove uses minimax with alpha-beta pruning (negamax formulation) to find the best move.
// depth specifies how many plies (half-moves) to look ahead.
func (g *Game) BestMove(depth int) Move {
	moves := g.AllLegalMoves()
	if len(moves) == 0 {
		return Move{}
	}

	moves = orderMoves(moves, g)
	bestScore := -evaluationInfinity
	bestMove := moves[0]
	alpha := -evaluationInfinity
	beta := evaluationInfinity

	for _, move := range moves {
		clone := g.Clone()
		clone.executeMove(move)
		score := -clone.negamax(depth-1, -beta, -alpha)
		if score > bestScore {
			bestScore = score
			bestMove = move
		}
		if score > alpha {
			alpha = score
		}
	}

	return bestMove
}

// negamax implements the negamax variant of alpha-beta pruning.
// Returns score from the perspective of the side to move.
func (g *Game) negamax(depth int, alpha, beta int) int {
	if depth == 0 {
		return g.evaluateForSide()
	}

	moves := g.AllLegalMoves()
	if len(moves) == 0 {
		if g.isInCheck(g.Turn) {
			return -KingValue // checkmate is worst possible
		}
		return 0 // stalemate
	}

	moves = orderMoves(moves, g)

	for _, move := range moves {
		clone := g.Clone()
		clone.executeMove(move)
		score := -clone.negamax(depth-1, -beta, -alpha)
		if score >= beta {
			return beta // fail high
		}
		if score > alpha {
			alpha = score
		}
	}

	return alpha
}

// orderMoves improves alpha-beta pruning efficiency by ordering moves.
// Captures of high-value pieces are searched first.
func orderMoves(moves []Move, g *Game) []Move {
	type scoredMove struct {
		move  Move
		score int
	}
	scored := make([]scoredMove, len(moves))
	for i, m := range moves {
		s := 0
		if m.Captured != nil || m.IsEnPassant {
			captureVal := 0
			if m.Captured != nil {
				captureVal = pieceValue(m.Captured.Type)
			} else {
				captureVal = PawnValue // en passant captures a pawn
			}
			s += 10*captureVal - pieceValue(m.Piece.Type)
		}
		if m.Promotion != NoPiece {
			s += pieceValue(m.Promotion)
		}
		if m.IsCastling {
			s += 50
		}
		if m.To.File >= 2 && m.To.File <= 5 && m.To.Rank >= 2 && m.To.Rank <= 5 {
			s += 10
		}
		scored[i] = scoredMove{move: m, score: s}
	}

	for i := 1; i < len(scored); i++ {
		for j := i; j > 0 && scored[j].score > scored[j-1].score; j-- {
			scored[j], scored[j-1] = scored[j-1], scored[j]
		}
	}

	result := make([]Move, len(moves))
	for i, sm := range scored {
		result[i] = sm.move
	}
	return result
}
