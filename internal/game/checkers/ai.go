package checkers

import (
	"math"
	"math/rand"
)

// RandomMove returns a random legal move.
func (g *Game) RandomMove() Move {
	moves := g.AllLegalMoves()
	if len(moves) == 0 {
		return Move{}
	}
	return moves[rand.Intn(len(moves))]
}

// BestMove returns the best move using minimax with alpha-beta pruning.
func (g *Game) BestMove(depth int) Move {
	moves := g.AllLegalMoves()
	if len(moves) == 0 {
		return Move{}
	}

	bestScore := math.Inf(-1)
	bestMove := moves[0]

	for _, move := range moves {
		// Clone game state
		newGame := g.clone()
		newGame.MakeMove(move)

		score := newGame.minimax(depth-1, math.Inf(-1), math.Inf(1), false)

		if score > bestScore {
			bestScore = score
			bestMove = move
		}
	}

	return bestMove
}

// minimax with alpha-beta pruning.
// maximizing is true when it's the original player's turn.
func (g *Game) minimax(depth int, alpha, beta float64, maximizing bool) float64 {
	if depth == 0 || g.State != Playing {
		return g.evaluate()
	}

	moves := g.AllLegalMoves()
	if len(moves) == 0 {
		// No moves = loss for current player
		if maximizing {
			return -10000
		}
		return 10000
	}

	if maximizing {
		maxEval := math.Inf(-1)
		for _, move := range moves {
			newGame := g.clone()
			newGame.MakeMove(move)
			eval := newGame.minimax(depth-1, alpha, beta, false)
			if eval > maxEval {
				maxEval = eval
			}
			if maxEval > alpha {
				alpha = maxEval
			}
			if beta <= alpha {
				break
			}
		}
		return maxEval
	}

	minEval := math.Inf(1)
	for _, move := range moves {
		newGame := g.clone()
		newGame.MakeMove(move)
		eval := newGame.minimax(depth-1, alpha, beta, true)
		if eval < minEval {
			minEval = eval
		}
		if minEval < beta {
			beta = minEval
		}
		if beta <= alpha {
			break
		}
	}
	return minEval
}

// evaluate returns a score for the position from White's perspective.
// Positive = good for White, negative = good for Black.
func (g *Game) evaluate() float64 {
	score := 0.0

	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			p := g.Board[row][col]
			if p.Empty() {
				continue
			}

			val := 0.0
			if p.Type == King {
				val = 150
				// Kings prefer center
				centerDist := math.Abs(float64(row)-3.5) + math.Abs(float64(col)-3.5)
				val += (4.0 - centerDist) * 3
			} else {
				val = 100
				// Men get bonus for advancement
				if p.Color == White {
					val += float64(7-row) * 3 // closer to row 0 = better
				} else {
					val += float64(row) * 3 // closer to row 7 = better
				}
				// Center control for men
				centerDist := math.Abs(float64(row)-3.5) + math.Abs(float64(col)-3.5)
				val += (4.0 - centerDist) * 1
			}

			if p.Color == White {
				score += val
			} else {
				score -= val
			}
		}
	}

	// Bonus/penalty for game state
	switch g.State {
	case WhiteWin:
		score += 10000
	case BlackWin:
		score -= 10000
	case Draw:
		score = 0
	}

	return score
}

// clone creates a deep copy of the game.
func (g *Game) clone() *Game {
	clone := &Game{
		Turn:           g.Turn,
		State:          g.State,
		KingMoveCount:  g.KingMoveCount,
		TotalKingMoves: g.TotalKingMoves,
		Moves:          make([]Move, len(g.Moves)),
	}
	copy(clone.Moves, g.Moves)
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			clone.Board[row][col] = g.Board[row][col]
		}
	}
	return clone
}
