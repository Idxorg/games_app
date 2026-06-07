package backgammon

// MustEnterFromBar returns true if the current player has checkers on the bar
// and must enter them before making any other move.
func (g *Game) MustEnterFromBar() bool {
	if g.Turn == White {
		return g.WhiteBar > 0
	}
	return g.BlackBar > 0
}

// CanBearOff returns true if all of the given color's checkers are in the
// home board (and none are on the bar).
//   - White home board: points 0-5
//   - Black home board: points 18-23
func (g *Game) CanBearOff(color Color) bool {
	if color == White && g.WhiteBar > 0 {
		return false
	}
	if color == Black && g.BlackBar > 0 {
		return false
	}
	if color == White {
		for i := 6; i < 24; i++ {
			if g.Board[i].Color == White && g.Board[i].Count > 0 {
				return false
			}
		}
	} else {
		for i := 0; i < 18; i++ {
			if g.Board[i].Color == Black && g.Board[i].Count > 0 {
				return false
			}
		}
	}
	return true
}

// legalMovesForDie returns all legal moves for the current player using the
// given die value. This does NOT enforce the "must use both dice" or "must
// use larger die" constraints — those are handled at a higher level.
func (g *Game) legalMovesForDie(die int) []Move {
	color := g.Turn
	opp := color.Opposite()
	var moves []Move

	// If on bar, must enter before any other move.
	if g.MustEnterFromBar() {
		target := g.barEntryPoint(color, die)
		if target >= 0 && target <= 23 && g.canLandOn(target, opp) {
			moves = append(moves, Move{From: BarPoint, To: target, Die: die})
		}
		return moves
	}

	canBear := g.CanBearOff(color)

	for i := 0; i < 24; i++ {
		if g.Board[i].Color != color || g.Board[i].Count == 0 {
			continue
		}

		target := g.moveTarget(color, i, die)

		// Bearing off.
		if isOffBoard(color, target) {
			if canBear && g.canBearOffFrom(color, i, die) {
				moves = append(moves, Move{From: i, To: OffPoint, Die: die})
			}
			continue
		}

		// Regular move.
		if g.canLandOn(target, opp) {
			moves = append(moves, Move{From: i, To: target, Die: die})
		}
	}

	return moves
}

// barEntryPoint computes the entry point index when entering from the bar.
//   - White enters in black's home board (indices 18-23): die d → index 24-d
//   - Black enters in white's home board (indices 0-5):  die d → index d-1
func (g *Game) barEntryPoint(color Color, die int) int {
	if color == White {
		return 24 - die // die 1→23, die 6→18
	}
	return die - 1 // die 1→0, die 6→5
}

// moveTarget computes the destination index for a regular move.
func (g *Game) moveTarget(color Color, from, die int) int {
	if color == White {
		return from - die
	}
	return from + die
}

// isOffBoard returns true if the target is past the bearing-off edge.
func isOffBoard(color Color, target int) bool {
	if color == White {
		return target < 0
	}
	return target > 23
}

// canLandOn returns true if the current player's checker can land on the
// target point. Valid if the point is: empty, occupied by own color, or
// has exactly 1 opponent checker (which would be hit).
func (g *Game) canLandOn(target int, opp Color) bool {
	pt := g.Board[target]
	return pt.Count == 0 || pt.Color != opp || pt.Count == 1
}

// canBearOffFrom returns true if a checker on the given point can be borne
// off with the given die value. Rules:
//   - Exact die: distance == die → always allowed.
//   - Higher die: distance < die → only from the farthest occupied point.
//   - Lower die: distance > die → not allowed.
func (g *Game) canBearOffFrom(color Color, point, die int) bool {
	dist := bearOffDistance(color, point)
	if die == dist {
		return true
	}
	if die > dist {
		return g.isFarthestInHome(color, point)
	}
	return false
}

// bearOffDistance returns the pip distance from the given point to the
// bearing-off edge.
//   - White: distance from index i is i+1 (point 0 needs 1, point 5 needs 6).
//   - Black: distance from index i is 24-i (point 23 needs 1, point 18 needs 6).
func bearOffDistance(color Color, point int) int {
	if color == White {
		return point + 1
	}
	return 24 - point
}

// isFarthestInHome returns true if the given point holds the checker that is
// farthest from the bearing-off edge among all checkers in the home board.
// This is required when using a die value higher than the exact distance.
//   - White: farthest = highest index in [0,5].
//   - Black: farthest = lowest index in [18,23].
func (g *Game) isFarthestInHome(color Color, point int) bool {
	if color == White {
		for i := point + 1; i <= 5; i++ {
			if g.Board[i].Color == White && g.Board[i].Count > 0 {
				return false
			}
		}
		return true
	}
	for i := point - 1; i >= 18; i-- {
		if g.Board[i].Color == Black && g.Board[i].Count > 0 {
			return false
		}
	}
	return true
}

// --- Dice constraint helpers ---

// uniqueInts returns the unique values from a slice, preserving order.
func uniqueInts(nums []int) []int {
	seen := make(map[int]bool)
	var result []int
	for _, n := range nums {
		if !seen[n] {
			seen[n] = true
			result = append(result, n)
		}
	}
	return result
}

// removeOne removes the first occurrence of val from the slice and returns
// a new slice (never mutates the original backing array).
func removeOne(nums []int, val int) []int {
	for i, n := range nums {
		if n == val {
			result := make([]int, 0, len(nums)-1)
			result = append(result, nums[:i]...)
			result = append(result, nums[i+1:]...)
			return result
		}
	}
	// Value not found; return a copy.
	result := make([]int, len(nums))
	copy(result, nums)
	return result
}

// maxMoveSequences computes the maximum number of moves achievable from the
// current game state with the given remaining dice. Uses recursive search
// with pruning.
func (g *Game) maxMoveSequences(remaining []int, depth int) int {
	if len(remaining) == 0 {
		return depth
	}

	maxDepth := depth
	tried := make(map[int]bool)

	for _, die := range remaining {
		if tried[die] {
			continue
		}
		tried[die] = true

		moves := g.legalMovesForDie(die)
		if len(moves) == 0 {
			continue
		}

		for _, move := range moves {
			g2 := g.clone()
			g2.applyMove(&move)
			newRemaining := removeOne(remaining, die)
			d := g2.maxMoveSequences(newRemaining, depth+1)
			if d > maxDepth {
				maxDepth = d
				// Early exit: can't beat using all remaining dice.
				if maxDepth == depth+len(remaining) {
					return maxDepth
				}
			}
		}
	}

	return maxDepth
}

// firstMovesFromMaxSequences returns all first moves that are part of a
// move sequence achieving maxDepth total moves.
func (g *Game) firstMovesFromMaxSequences(maxDepth int) []Move {
	remaining := g.RemainingMoves
	if maxDepth == 0 {
		return nil
	}

	var result []Move
	seen := make(map[Move]bool)
	tried := make(map[int]bool)

	for _, die := range remaining {
		if tried[die] {
			continue
		}
		tried[die] = true

		moves := g.legalMovesForDie(die)
		for _, move := range moves {
			if seen[move] {
				continue
			}

			g2 := g.clone()
			g2.applyMove(&move)
			newRemaining := removeOne(remaining, die)
			d := g2.maxMoveSequences(newRemaining, 1)

			if d == maxDepth {
				result = append(result, move)
				seen[move] = true
			}
		}
	}

	return result
}
