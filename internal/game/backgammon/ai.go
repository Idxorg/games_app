package backgammon

import (
	"math"
	"math/rand"
)

// RandomMove picks a random legal move from the available moves.
// Returns a zero Move if no legal moves exist.
func RandomMove(g *Game) Move {
	moves := g.AllLegalMoves()
	if len(moves) == 0 {
		return Move{}
	}
	return moves[rand.Intn(len(moves))]
}

// BestMove picks the best move using a one-step greedy evaluation.
// Returns a zero Move if no legal moves exist.
func BestMove(g *Game) Move {
	moves := g.AllLegalMoves()
	if len(moves) == 0 {
		return Move{}
	}

	bestScore := -math.MaxFloat64
	bestMove := moves[0]
	color := g.Turn

	for _, m := range moves {
		g2 := g.clone()
		g2.applyMove(&m)
		score := Evaluate(g2, color)
		if score > bestScore {
			bestScore = score
			bestMove = m
		}
	}

	return bestMove
}

// Evaluate computes a heuristic score for the position from the perspective
// of the given color. Higher score is better.
//
// Components:
//   - Pip count: fewer pips is better.
//   - Blots (single exposed checkers): fewer is better.
//   - Bar count: fewer on bar is better.
//   - Borne off: more is better.
//   - Home board points made (2+ checkers): more is better.
//   - Hitting bonus: opponent on bar is good.
func Evaluate(g *Game, color Color) float64 {
	opp := color.Opposite()
	score := 0.0

	// Pip count (lower is better for us).
	score -= float64(pipCount(g, color)) * 1.0
	score += float64(pipCount(g, opp)) * 1.0

	// Blots: our blots are bad, opponent blots are good.
	score -= float64(countBlots(g, color)) * 3.0
	score += float64(countBlots(g, opp)) * 3.0

	// Bar: own bar is very bad, opponent bar is good.
	score -= float64(barCount(g, color)) * 4.0
	score += float64(barCount(g, opp)) * 4.0

	// Borne off: more is better.
	score += float64(offCount(g, color)) * 5.0
	score -= float64(offCount(g, opp)) * 5.0

	// Home board points made (2+ checkers) are strategically strong.
	score += float64(homeBoardPoints(g, color)) * 2.0
	score -= float64(homeBoardPoints(g, opp)) * 2.0

	// Hitting bonus: extra weight for opponent on bar.
	score += float64(barCount(g, opp)) * 2.0

	return score
}

// pipCount returns the total pip count for a color.
// Each checker on the board contributes its distance to bear off.
// Checkers on the bar contribute 25 pips (approximate).
func pipCount(g *Game, color Color) int {
	count := 0

	bar := barCount(g, color)
	count += bar * 25

	for i := 0; i < 24; i++ {
		if g.Board[i].Color == color && g.Board[i].Count > 0 {
			if color == White {
				count += g.Board[i].Count * (i + 1)
			} else {
				count += g.Board[i].Count * (24 - i)
			}
		}
	}

	return count
}

// countBlots returns the number of single (exposed) checkers for a color.
func countBlots(g *Game, color Color) int {
	count := 0
	for i := 0; i < 24; i++ {
		if g.Board[i].Color == color && g.Board[i].Count == 1 {
			count++
		}
	}
	return count
}

// barCount returns the number of checkers on the bar for a color.
func barCount(g *Game, color Color) int {
	if color == White {
		return g.WhiteBar
	}
	return g.BlackBar
}

// offCount returns the number of checkers already borne off for a color.
func offCount(g *Game, color Color) int {
	if color == White {
		return g.WhiteOff
	}
	return g.BlackOff
}

// homeBoardPoints returns the number of home board points with 2+ checkers
// (made points) for a color.
func homeBoardPoints(g *Game, color Color) int {
	count := 0
	if color == White {
		for i := 0; i < 6; i++ {
			if g.Board[i].Color == color && g.Board[i].Count >= 2 {
				count++
			}
		}
	} else {
		for i := 18; i < 24; i++ {
			if g.Board[i].Color == color && g.Board[i].Count >= 2 {
				count++
			}
		}
	}
	return count
}
