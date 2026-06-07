package backgammon

import (
	"math"
	"math/rand"
	"time"
)

// All 21 possible dice combinations (unordered pairs).
var allDiceCombos = [][2]int{
	{1, 1}, {2, 2}, {3, 3}, {4, 4}, {5, 5}, {6, 6}, // doubles (6 combos, weight 1 each)
	{1, 2}, {1, 3}, {1, 4}, {1, 5}, {1, 6}, // mixed (15 combos, weight 2 each)
	{2, 3}, {2, 4}, {2, 5}, {2, 6},
	{3, 4}, {3, 5}, {3, 6},
	{4, 5}, {4, 6},
	{5, 6},
}

// diceProbability returns the probability of a given dice combination.
// Doubles: 1/36, Mixed: 2/36.
func diceProbability(d1, d2 int) float64 {
	if d1 == d2 {
		return 1.0 / 36.0
	}
	return 2.0 / 36.0
}

// All possible dice sequences for a given dice combo.
func diceSequences(d1, d2 int) [][]int {
	if d1 == d2 {
		return [][]int{{d1, d1, d1, d1}} // doubles: 4 moves
	}
	return [][]int{{d1, d2}, {d2, d1}} // mixed: try both orders
}

// RandomMove picks a random legal move from the available moves.
func RandomMove(g *Game) Move {
	moves := g.AllLegalMoves()
	if len(moves) == 0 {
		return Move{}
	}
	return moves[rand.Intn(len(moves))]
}

// BestMove picks the best move using expectiminimax.
// maxDepth is the number of player turns to look ahead (each turn = roll + move).
// If maxDepth <= 0, defaults to 2.
func BestMove(g *Game, maxDepth int) Move {
	if maxDepth <= 0 {
		maxDepth = 2
	}
	return bestMoveForDice(g, g.RemainingMoves, maxDepth, math.Inf(-1), math.Inf(1))
}

// bestMoveForDice evaluates all legal moves and returns the best one.
// Uses expectiminimax with the remaining dice already set.
func bestMoveForDice(g *Game, remaining []int, maxDepth int, alpha, beta float64) Move {
	moves := g.AllLegalMoves()
	if len(moves) == 0 {
		return Move{}
	}

	bestScore := -math.MaxFloat64
	bestMove := moves[0]

	for _, move := range moves {
		g2 := g.clone()
		g2.applyMove(&move)
		g2.Moves = append(g2.Moves, move)

		// Remove used die
		newRemaining := removeOne(remaining, move.Die)

		// Evaluate: continue with remaining dice or complete the turn
		var score float64
		if len(newRemaining) > 0 {
			score = bestMoveScore(g2, newRemaining, maxDepth, alpha, beta)
		} else {
			score = chanceNode(g2, maxDepth-1, alpha, beta)
		}

		if score > bestScore {
			bestScore = score
			bestMove = move
		}
		if bestScore > alpha {
			alpha = bestScore
		}
	}

	return bestMove
}

// bestMoveScore recursively evaluates remaining dice in a turn.
func bestMoveScore(g *Game, remaining []int, maxDepth int, alpha, beta float64) float64 {
	moves := g.AllLegalMoves()
	if len(moves) == 0 {
		// No more moves with remaining dice — turn ends
		return chanceNode(g, maxDepth-1, alpha, beta)
	}

	bestScore := -math.MaxFloat64
	for _, move := range moves {
		g2 := g.clone()
		g2.applyMove(&move)

		newRemaining := removeOne(remaining, move.Die)
		var score float64
		if len(newRemaining) > 0 {
			score = bestMoveScore(g2, newRemaining, maxDepth, alpha, beta)
		} else {
			score = chanceNode(g2, maxDepth-1, alpha, beta)
		}

		if score > bestScore {
			bestScore = score
		}
		if bestScore > alpha {
			alpha = bestScore
		}
		if alpha >= beta {
			break
		}
	}

	return bestScore
}

// chanceNode averages over all possible dice rolls for the next player.
// This is the "expecti" part of expectiminimax.
func chanceNode(g *Game, remainingDepth int, alpha, beta float64) float64 {
	if remainingDepth <= 0 {
		return Evaluate(g, g.Turn)
	}

	// Check for game over
	if g.State == GameOver {
		if g.Winner == g.Turn {
			return 100000 // last player won
		}
		return -100000
	}

	// Check if current player has won (all off)
	if g.Turn == White && g.WhiteOff == 15 {
		return 100000
	}
	if g.Turn == Black && g.BlackOff == 15 {
		return 100000
	}

	// Average over all 21 dice combinations
	totalScore := 0.0
	totalProb := 0.0

	for _, dice := range allDiceCombos {
		prob := diceProbability(dice[0], dice[1])
		sequences := diceSequences(dice[0], dice[1])

		bestSeqScore := -math.MaxFloat64
		for _, seq := range sequences {
			g2 := g.clone()
			g2.Turn = g.Turn
			g2.SetDice(dice[0], dice[1])
			g2.RemainingMoves = make([]int, len(seq))
			copy(g2.RemainingMoves, seq)

			score := bestMoveScore(g2, g2.RemainingMoves, remainingDepth, alpha, beta)
			if score > bestSeqScore {
				bestSeqScore = score
			}
		}

		totalScore += bestSeqScore * prob
		totalProb += prob
	}

	if totalProb == 0 {
		return Evaluate(g, g.Turn)
	}
	return totalScore / totalProb
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
//   - Prime/blockade detection: consecutive blocked points
//   - Anchor detection: made points in opponent's home board
//   - Backgame potential
func Evaluate(g *Game, color Color) float64 {
	opp := color.Opposite()
	score := 0.0

	// --- Pip count (lower is better for us) ---
	myPip := pipCount(g, color)
	oppPip := pipCount(g, opp)
	score += float64(oppPip-myPip) * 1.0

	// --- Blots: our blots are bad, opponent blots are good ---
	score -= float64(countBlots(g, color)) * 4.0
	score += float64(countBlots(g, opp)) * 3.0

	// --- Bar: own bar is very bad, opponent bar is good ---
	score -= float64(barCount(g, color)) * 15.0
	score += float64(barCount(g, opp)) * 10.0

	// --- Borne off: more is better ---
	score += float64(offCount(g, color)) * 8.0
	score -= float64(offCount(g, opp)) * 8.0

	// --- Home board points made (2+ checkers) ---
	score += float64(homeBoardPoints(g, color)) * 4.0
	score -= float64(homeBoardPoints(g, opp)) * 4.0

	// --- Prime detection: consecutive made points in our home area ---
	score += float64(primeLength(g, color)) * 6.0
	score -= float64(primeLength(g, opp)) * 6.0

	// --- Anchor detection: made points in opponent's home board ---
	score += float64(anchors(g, color)) * 3.0
	score -= float64(anchors(g, opp)) * 3.0

	// --- Builder bonus: checkers on points adjacent to own made points ---
	score += float64(builders(g, color)) * 2.0
	score -= float64(builders(g, opp)) * 2.0

	// --- Hit potential: opponent blots within striking range of our checkers ---
	score += float64(hitPotential(g, color)) * 2.0
	score -= float64(hitPotential(g, opp)) * 2.0

	// --- Backgame potential: anchored points deep in opponent's home board ---
	score += float64(backgamePotential(g, color)) * 2.0
	score -= float64(backgamePotential(g, opp)) * 2.0

	// --- Race mode: when both sides are bearing off, pip count dominates ---
	if canRace(g) {
		raceBonus := float64(oppPip-myPip) * 3.0 // extra weight in race
		score += raceBonus
	}

	return score
}

// canRace returns true if both sides are close to bearing off (all in home or off).
func canRace(g *Game) bool {
	return g.CanBearOff(White) && g.CanBearOff(Black)
}

// pipCount returns the total pip count for a color.
func pipCount(g *Game, color Color) int {
	count := barCount(g, color) * 25

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

// homeBoardPoints returns the number of home board points with 2+ checkers.
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

// primeLength returns the length of the longest consecutive prime (made points)
// for a color anywhere on the board.
func primeLength(g *Game, color Color) int {
	best := 0
	current := 0
	start := 0
	if color == Black {
		start = 0
	} else {
		start = 0
	}
	for i := start; i < 24; i++ {
		if g.Board[i].Color == color && g.Board[i].Count >= 2 {
			current++
			if current > best {
				best = current
			}
		} else {
			current = 0
		}
	}
	return best
}

// anchors returns the number of made points (2+) in the opponent's home board.
func anchors(g *Game, color Color) int {
	count := 0
	if color == White {
		// White anchors in Black's home (18-23)
		for i := 18; i < 24; i++ {
			if g.Board[i].Color == White && g.Board[i].Count >= 2 {
				count++
			}
		}
	} else {
		// Black anchors in White's home (0-5)
		for i := 0; i < 6; i++ {
			if g.Board[i].Color == Black && g.Board[i].Count >= 2 {
				count++
			}
		}
	}
	return count
}

// builders returns the number of checkers that are alone on a point adjacent
// to a made point of the same color (potential to make a new point).
func builders(g *Game, color Color) int {
	count := 0
	for i := 0; i < 24; i++ {
		if g.Board[i].Color != color || g.Board[i].Count == 0 {
			continue
		}
		// Check if adjacent to a made point
		for _, adj := range []int{i - 1, i + 1} {
			if adj >= 0 && adj < 24 && g.Board[adj].Color == color && g.Board[adj].Count >= 2 {
				count++
				break
			}
		}
	}
	return count
}

// hitPotential returns the number of opponent blots within range of our checkers
// (6-point range, the maximum die distance).
func hitPotential(g *Game, color Color) int {
	opp := color.Opposite()
	potential := 0
	for i := 0; i < 24; i++ {
		if g.Board[i].Color != opp || g.Board[i].Count != 1 {
			continue
		}
		// Check if any of our checkers can reach this blot
		for j := 0; j < 24; j++ {
			if g.Board[j].Color != color || g.Board[j].Count == 0 {
				continue
			}
			dist := 0
			if color == White {
				dist = j - i // White moves from high to low
			} else {
				dist = i - j // Black moves from low to high
			}
			if dist >= 1 && dist <= 6 {
				potential++
				break
			}
		}
	}
	return potential
}

// backgamePotential returns a bonus for having anchored points deep in
// opponent's home board (valuable when behind in the race).
func backgamePotential(g *Game, color Color) int {
	count := 0
	if color == White {
		// Deep anchors in Black's home (points 20-23 are "deep")
		for i := 20; i < 24; i++ {
			if g.Board[i].Color == White && g.Board[i].Count >= 2 {
				count++
			}
		}
	} else {
		// Deep anchors in White's home (points 0-3 are "deep")
		for i := 0; i < 4; i++ {
			if g.Board[i].Color == Black && g.Board[i].Count >= 2 {
				count++
			}
		}
	}
	return count
}

// --- Time helper to avoid import issues ---
var _ = time.Now
