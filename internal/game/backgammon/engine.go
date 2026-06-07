package backgammon

import (
	"fmt"
	"math/rand"
)

// NewGame creates a new backgammon game with the standard starting position.
//
// Standard starting position (0-indexed):
//   - White: 2 on point 23, 5 on point 12, 3 on point 7, 5 on point 5
//   - Black: 2 on point 0,  5 on point 11, 3 on point 16, 5 on point 18
func NewGame() *Game {
	g := &Game{
		Turn:  White,
		State: Rolling,
	}

	// White checkers (moves from high index → low index)
	g.Board[23] = Point{White, 2} // 2 on 24-point
	g.Board[12] = Point{White, 5} // 5 on 13-point
	g.Board[7] = Point{White, 3}  // 3 on 8-point
	g.Board[5] = Point{White, 5}  // 5 on 6-point

	// Black checkers (moves from low index → high index)
	g.Board[0] = Point{Black, 2}  // 2 on 1-point (black's 24-point)
	g.Board[11] = Point{Black, 5} // 5 on 12-point (black's 13-point)
	g.Board[16] = Point{Black, 3} // 3 on 17-point (black's 8-point)
	g.Board[18] = Point{Black, 5} // 5 on 19-point (black's 6-point)

	return g
}

// RollDice rolls two dice and populates remaining moves.
// Doubles give 4 moves of the same value.
func (g *Game) RollDice() [2]int {
	d1 := rand.Intn(6) + 1
	d2 := rand.Intn(6) + 1
	g.Dice = [2]int{d1, d2}
	if d1 == d2 {
		g.RemainingMoves = []int{d1, d1, d1, d1}
	} else {
		g.RemainingMoves = []int{d1, d2}
	}
	g.State = Moving
	return g.Dice
}

// SetDice sets specific dice values. Useful for testing.
func (g *Game) SetDice(d1, d2 int) {
	g.Dice = [2]int{d1, d2}
	if d1 == d2 {
		g.RemainingMoves = []int{d1, d1, d1, d1}
	} else {
		g.RemainingMoves = []int{d1, d2}
	}
	g.State = Moving
}

// AllLegalMoves returns all legal first moves considering dice constraints:
//   - Must use the maximum number of dice possible.
//   - If only one die can be used, must use the larger die.
func (g *Game) AllLegalMoves() []Move {
	if g.State != Moving || len(g.RemainingMoves) == 0 {
		return nil
	}

	maxDepth := g.maxMoveSequences(g.RemainingMoves, 0)
	if maxDepth == 0 {
		return nil
	}

	result := g.firstMovesFromMaxSequences(maxDepth)
	if len(result) == 0 {
		return nil
	}

	// If maxDepth == 1 with two different dice, enforce "must use larger die" rule.
	unique := uniqueInts(g.RemainingMoves)
	if maxDepth == 1 && len(unique) == 2 {
		var larger int
		if unique[0] > unique[1] {
			larger = unique[0]
		} else {
			larger = unique[1]
		}
		largerMoves := g.legalMovesForDie(larger)
		if len(largerMoves) > 0 {
			var filtered []Move
			for _, m := range result {
				if m.Die == larger {
					filtered = append(filtered, m)
				}
			}
			if len(filtered) > 0 {
				return filtered
			}
		}
	}

	return result
}

// MakeMove validates and executes a single move, then updates remaining dice.
// Returns an error if the move is not legal.
func (g *Game) MakeMove(move Move) error {
	legal := g.AllLegalMoves()
	found := false
	for _, m := range legal {
		if m.From == move.From && m.To == move.To && m.Die == move.Die {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("illegal move: from=%d to=%d die=%d", move.From, move.To, move.Die)
	}

	g.applyMove(&move)
	g.Moves = append(g.Moves, move)

	// Remove the used die from remaining moves.
	for i, d := range g.RemainingMoves {
		if d == move.Die {
			g.RemainingMoves = append(g.RemainingMoves[:i], g.RemainingMoves[i+1:]...)
			break
		}
	}

	// Check for win.
	if g.WhiteOff == 15 {
		g.State = GameOver
		g.Winner = White
		return nil
	}
	if g.BlackOff == 15 {
		g.State = GameOver
		g.Winner = Black
		return nil
	}

	// If no more dice or no legal moves, end the turn.
	if len(g.RemainingMoves) == 0 || len(g.AllLegalMoves()) == 0 {
		g.endTurn()
	}

	return nil
}

// PassTurn passes remaining dice. Only valid when no legal moves exist.
func (g *Game) PassTurn() error {
	if g.State != Moving {
		return fmt.Errorf("not in moving state")
	}
	if len(g.AllLegalMoves()) > 0 {
		return fmt.Errorf("cannot pass: legal moves available")
	}
	g.endTurn()
	return nil
}

// clone creates a deep copy of the game for simulation.
func (g *Game) clone() *Game {
	gc := *g // value copy handles Board array and primitive fields
	gc.RemainingMoves = make([]int, len(g.RemainingMoves))
	copy(gc.RemainingMoves, g.RemainingMoves)
	gc.Moves = make([]Move, len(g.Moves))
	copy(gc.Moves, g.Moves)
	return &gc
}

// applyMove mutates the board state according to the move (no validation).
func (g *Game) applyMove(move *Move) {
	color := g.Turn
	opp := color.Opposite()

	// Remove checker from source.
	if move.From == BarPoint {
		if color == White {
			g.WhiteBar--
		} else {
			g.BlackBar--
		}
	} else {
		g.Board[move.From].Count--
	}

	// Place checker at destination.
	if move.To == OffPoint {
		if color == White {
			g.WhiteOff++
		} else {
			g.BlackOff++
		}
	} else {
		// Hit: landing on a single opponent checker sends it to the bar.
		if g.Board[move.To].Color == opp && g.Board[move.To].Count == 1 {
			if opp == White {
				g.WhiteBar++
			} else {
				g.BlackBar++
			}
			g.Board[move.To].Count = 0
		}
		g.Board[move.To].Color = color
		g.Board[move.To].Count++
	}
}

// endTurn switches play to the other player and resets dice state.
func (g *Game) endTurn() {
	g.Turn = g.Turn.Opposite()
	g.RemainingMoves = nil
	g.Dice = [2]int{}
	g.State = Rolling
}
