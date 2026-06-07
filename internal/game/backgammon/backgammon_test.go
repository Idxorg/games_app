package backgammon

import (
	"fmt"
	"testing"
)

func assert(t *testing.T, cond bool, msg string, args ...interface{}) {
	t.Helper()
	if !cond {
		t.Errorf(msg, args...)
	}
}

// ---------------------------------------------------------------------------
// TestNewGame
// ---------------------------------------------------------------------------

func TestNewGame(t *testing.T) {
	g := NewGame()

	assert(t, g.Turn == White, "white should go first, got %v", g.Turn)
	assert(t, g.State == Rolling, "initial state should be Rolling, got %v", g.State)
	assert(t, g.WhiteBar == 0, "white bar should be 0")
	assert(t, g.BlackBar == 0, "black bar should be 0")
	assert(t, g.WhiteOff == 0, "white off should be 0")
	assert(t, g.BlackOff == 0, "black off should be 0")

	// White starting positions.
	assert(t, g.Board[23] == Point{White, 2}, "white: 2 on point 24 (idx 23), got %+v", g.Board[23])
	assert(t, g.Board[12] == Point{White, 5}, "white: 5 on point 13 (idx 12), got %+v", g.Board[12])
	assert(t, g.Board[7] == Point{White, 3}, "white: 3 on point 8 (idx 7), got %+v", g.Board[7])
	assert(t, g.Board[5] == Point{White, 5}, "white: 5 on point 6 (idx 5), got %+v", g.Board[5])

	// Black starting positions.
	assert(t, g.Board[0] == Point{Black, 2}, "black: 2 on point 1 (idx 0), got %+v", g.Board[0])
	assert(t, g.Board[11] == Point{Black, 5}, "black: 5 on point 12 (idx 11), got %+v", g.Board[11])
	assert(t, g.Board[16] == Point{Black, 3}, "black: 3 on point 17 (idx 16), got %+v", g.Board[16])
	assert(t, g.Board[18] == Point{Black, 5}, "black: 5 on point 19 (idx 18), got %+v", g.Board[18])

	// Total checkers.
	assert(t, g.TotalCheckers(White) == 15,
		"white should have 15 checkers, got %d", g.TotalCheckers(White))
	assert(t, g.TotalCheckers(Black) == 15,
		"black should have 15 checkers, got %d", g.TotalCheckers(Black))
}

// ---------------------------------------------------------------------------
// TestRollDice
// ---------------------------------------------------------------------------

func TestRollDice(t *testing.T) {
	g := NewGame()

	seen := make(map[[2]int]bool)
	for i := 0; i < 100; i++ {
		dice := g.RollDice()
		assert(t, dice[0] >= 1 && dice[0] <= 6,
			"die 1 out of range: %d", dice[0])
		assert(t, dice[1] >= 1 && dice[1] <= 6,
			"die 2 out of range: %d", dice[1])
		seen[dice] = true
	}
	// Should have seen many different combinations in 100 rolls.
	assert(t, len(seen) >= 10,
		"should see >= 10 unique dice combos in 100 rolls, saw %d", len(seen))
}

// ---------------------------------------------------------------------------
// TestDoublesSet
// ---------------------------------------------------------------------------

func TestDoublesSet(t *testing.T) {
	g := NewGame()
	g.SetDice(4, 4)

	assert(t, g.Dice[0] == 4 && g.Dice[1] == 4, "dice should be [4,4]")
	assert(t, len(g.RemainingMoves) == 4, "doubles should give 4 remaining moves")
	assert(t, g.State == Moving, "state should be Moving")

	g.SetDice(3, 5)
	assert(t, len(g.RemainingMoves) == 2, "non-doubles should give 2 remaining moves")
}

// ---------------------------------------------------------------------------
// TestLegalMoves_simpleForward
// ---------------------------------------------------------------------------

func TestLegalMoves_simpleForward(t *testing.T) {
	g := &Game{Turn: White, State: Moving}
	g.Board[12] = Point{White, 1}
	g.WhiteOff = 14
	g.RemainingMoves = []int{3}

	// White on index 12, die 3: target = 12-3 = 9 (empty, legal).
	moves := g.AllLegalMoves()
	assert(t, len(moves) == 1, "should have 1 legal move, got %d", len(moves))
	assert(t, moves[0].From == 12, "from should be 12, got %d", moves[0].From)
	assert(t, moves[0].To == 9, "to should be 9, got %d", moves[0].To)
	assert(t, moves[0].Die == 3, "die should be 3, got %d", moves[0].Die)

	// Execute the move.
	err := g.MakeMove(moves[0])
	assert(t, err == nil, "move should succeed: %v", err)
	assert(t, g.Board[12].Count == 0, "source should be empty")
	assert(t, g.Board[9].Color == White && g.Board[9].Count == 1,
		"destination should have 1 white checker")
}

// ---------------------------------------------------------------------------
// TestLegalMoves_ownPoint
// ---------------------------------------------------------------------------

func TestLegalMoves_ownPoint(t *testing.T) {
	g := &Game{Turn: White, State: Moving}
	g.Board[12] = Point{White, 2} // 2 white on point 12
	g.Board[9] = Point{White, 3}  // 3 white on point 9
	g.WhiteOff = 10
	g.RemainingMoves = []int{3}

	// White on 12, die 3: target 9. Point 9 has own pieces. Legal.
	moves := g.AllLegalMoves()
	found := false
	var moveToOwn Move
	for _, m := range moves {
		if m.From == 12 && m.To == 9 {
			found = true
			moveToOwn = m
		}
	}
	assert(t, found, "should be able to move to own point")

	err := g.MakeMove(moveToOwn)
	assert(t, err == nil, "move should succeed: %v", err)
	assert(t, g.Board[9].Color == White && g.Board[9].Count == 4,
		"point 9 should have 4 white checkers after move")
}

// ---------------------------------------------------------------------------
// TestLegalMoves_blocked
// ---------------------------------------------------------------------------

func TestLegalMoves_blocked(t *testing.T) {
	g := &Game{Turn: White, State: Moving}
	g.Board[10] = Point{White, 1}
	g.Board[7] = Point{Black, 2} // blocked by 2 black
	g.WhiteOff = 14
	g.RemainingMoves = []int{3}

	// White on index 10, die 3: target = 7. Blocked by 2 black.
	moves := g.AllLegalMoves()
	assert(t, len(moves) == 0, "should have no moves when target is blocked")
}

// ---------------------------------------------------------------------------
// TestHit
// ---------------------------------------------------------------------------

func TestHit(t *testing.T) {
	g := &Game{Turn: White, State: Moving}
	g.Board[10] = Point{White, 1}
	g.Board[7] = Point{Black, 1} // blot
	g.WhiteOff = 14
	g.BlackOff = 14
	g.RemainingMoves = []int{3}

	moves := g.AllLegalMoves()
	assert(t, len(moves) >= 1, "should have at least 1 move")

	// Find and execute the hitting move.
	found := false
	for _, m := range moves {
		if m.From == 10 && m.To == 7 {
			found = true
			err := g.MakeMove(m)
			assert(t, err == nil, "hit move should succeed: %v", err)
			break
		}
	}
	assert(t, found, "should find hitting move")
	assert(t, g.BlackBar == 1, "black should have 1 checker on bar")
	assert(t, g.Board[7].Color == White && g.Board[7].Count == 1,
		"point 7 should now be 1 white checker")
	assert(t, g.TotalCheckers(Black) == 15,
		"black should still have 15 total checkers")
}

// ---------------------------------------------------------------------------
// TestMustEnterFromBar
// ---------------------------------------------------------------------------

func TestMustEnterFromBar(t *testing.T) {
	g := &Game{Turn: White, State: Moving}
	g.WhiteBar = 1
	g.Board[23] = Point{White, 14}
	g.RemainingMoves = []int{1}

	assert(t, g.MustEnterFromBar(), "white should need to enter from bar")

	// Die 1: enters at index 23 (24-1). Point 23 has white. Legal.
	moves := g.AllLegalMoves()
	assert(t, len(moves) == 1, "should have exactly 1 move")
	assert(t, moves[0].From == BarPoint, "move should be from bar")
	assert(t, moves[0].To == 23, "should enter at index 23")

	// Can't make other board moves when on bar.
	// Even if there are other white pieces, bar entry is mandatory.
	g.Board[10] = Point{White, 1}
	moves = g.AllLegalMoves()
	for _, m := range moves {
		assert(t, m.From == BarPoint,
			"all moves should be from bar when checkers are on bar, got From=%d", m.From)
	}
}

// ---------------------------------------------------------------------------
// TestEnterFromBar_blocked
// ---------------------------------------------------------------------------

func TestEnterFromBar_blocked(t *testing.T) {
	g := &Game{Turn: White, State: Moving}
	g.WhiteBar = 1
	// Block all entry points for white (indices 18-23).
	for i := 18; i <= 23; i++ {
		g.Board[i] = Point{Black, 2}
	}
	g.RemainingMoves = []int{1, 2}

	// Die 1 → index 23: blocked. Die 2 → index 22: blocked.
	moves := g.AllLegalMoves()
	assert(t, len(moves) == 0, "should have no moves when bar entry is blocked")
}

// ---------------------------------------------------------------------------
// TestCanBearOff
// ---------------------------------------------------------------------------

func TestCanBearOff_allInHome(t *testing.T) {
	g := &Game{Turn: White, State: Moving}
	g.Board[3] = Point{White, 15}
	g.Board[10] = Point{Black, 15} // black pieces outside black's home (18-23)
	g.RemainingMoves = []int{4}

	assert(t, g.CanBearOff(White), "should be able to bear off when all in home")
	assert(t, !g.CanBearOff(Black), "black should not be able to bear off (pieces outside home)")
}

func TestCanBearOff_notAllInHome(t *testing.T) {
	g := &Game{Turn: White, State: Moving}
	g.Board[3] = Point{White, 14}
	g.Board[10] = Point{White, 1} // outside home
	g.RemainingMoves = []int{4}

	assert(t, !g.CanBearOff(White), "should not bear off with pieces outside home")
}

func TestCanBearOff_onBar(t *testing.T) {
	g := &Game{Turn: White, State: Moving}
	g.Board[3] = Point{White, 14}
	g.WhiteBar = 1
	g.RemainingMoves = []int{4}

	assert(t, !g.CanBearOff(White), "should not bear off with checkers on bar")
}

// ---------------------------------------------------------------------------
// TestBearingOff_exact
// ---------------------------------------------------------------------------

func TestBearingOff_exact(t *testing.T) {
	g := &Game{Turn: White, State: Moving}
	g.Board[2] = Point{White, 1} // distance = 3
	g.Board[0] = Point{White, 1} // distance = 1
	g.WhiteOff = 13
	g.RemainingMoves = []int{3, 1}

	// Die 3: from index 2, exact bear off (distance 3).
	moves := g.AllLegalMoves()
	foundExact3 := false
	for _, m := range moves {
		if m.From == 2 && m.To == OffPoint && m.Die == 3 {
			foundExact3 = true
		}
	}
	assert(t, foundExact3, "should be able to bear off with exact die 3 from index 2")

	// Die 1: from index 0, exact bear off (distance 1).
	foundExact1 := false
	for _, m := range moves {
		if m.From == 0 && m.To == OffPoint && m.Die == 1 {
			foundExact1 = true
		}
	}
	assert(t, foundExact1, "should be able to bear off with exact die 1 from index 0")
}

// ---------------------------------------------------------------------------
// TestBearingOff_higherDie
// ---------------------------------------------------------------------------

func TestBearingOff_higherDie(t *testing.T) {
	g := &Game{Turn: White, State: Moving}
	g.Board[3] = Point{White, 1} // distance = 4
	g.Board[0] = Point{White, 1} // distance = 1
	g.WhiteOff = 13
	g.RemainingMoves = []int{6}

	// Die 6: from index 3, distance 4 < 6. Is farthest in home? Check 4,5 — no. Yes.
	// Die 6: from index 0, distance 1 < 6. Is farthest in home? Check 1,2,3,4,5 — index 3 has piece. No.
	moves := g.AllLegalMoves()
	assert(t, len(moves) == 1, "should have exactly 1 move (farthest only), got %d", len(moves))
	assert(t, moves[0].From == 3, "should bear off from farthest point 3, got %d", moves[0].From)
	assert(t, moves[0].To == OffPoint, "should be bearing off")
}

func TestBearingOff_higherDie_black(t *testing.T) {
	g := &Game{Turn: Black, State: Moving}
	g.Board[20] = Point{Black, 1} // distance = 4
	g.Board[23] = Point{Black, 1} // distance = 1
	g.BlackOff = 13
	g.RemainingMoves = []int{6}

	// Die 6: from index 20, distance 4 < 6. Is farthest? Check 19,18 — no. Yes.
	// Die 6: from index 23, distance 1 < 6. Is farthest? Check 22,21,20 — index 20 has piece. No.
	moves := g.AllLegalMoves()
	assert(t, len(moves) == 1, "should have exactly 1 move for black, got %d", len(moves))
	assert(t, moves[0].From == 20, "black should bear off from farthest point 20, got %d", moves[0].From)
}

// ---------------------------------------------------------------------------
// TestBearingOff_cannotUseLowerDie
// ---------------------------------------------------------------------------

func TestBearingOff_cannotUseLowerDie(t *testing.T) {
	g := &Game{Turn: White, State: Moving}
	g.Board[5] = Point{White, 1} // distance = 6
	g.WhiteOff = 14
	g.RemainingMoves = []int{3}

	// Die 3 < distance 6. Can't bear off. But can move within home: 5-3=2.
	moves := g.AllLegalMoves()
	assert(t, len(moves) == 1, "should have 1 move (regular, not bear off), got %d", len(moves))
	assert(t, moves[0].To == 2, "should move to index 2, got %d", moves[0].To)
	assert(t, moves[0].To != OffPoint, "should NOT bear off with lower die")
}

// ---------------------------------------------------------------------------
// TestDoubles
// ---------------------------------------------------------------------------

func TestDoubles(t *testing.T) {
	g := NewGame()
	g.SetDice(3, 3)

	assert(t, len(g.RemainingMoves) == 4, "doubles should give 4 remaining moves")

	moves := g.AllLegalMoves()
	assert(t, len(moves) > 0, "should have legal moves")

	// Execute first move.
	err := g.MakeMove(moves[0])
	assert(t, err == nil, "first move should succeed: %v", err)
	assert(t, len(g.RemainingMoves) == 3, "should have 3 remaining after first move")

	// Execute second move.
	moves = g.AllLegalMoves()
	if len(moves) > 0 {
		err = g.MakeMove(moves[0])
		assert(t, err == nil, "second move should succeed: %v", err)
		assert(t, len(g.RemainingMoves) == 2, "should have 2 remaining after second move")
	}

	// Execute third move.
	moves = g.AllLegalMoves()
	if len(moves) > 0 {
		err = g.MakeMove(moves[0])
		assert(t, err == nil, "third move should succeed: %v", err)
		assert(t, len(g.RemainingMoves) == 1, "should have 1 remaining after third move")
	}

	// Execute fourth move.
	moves = g.AllLegalMoves()
	if len(moves) > 0 {
		err = g.MakeMove(moves[0])
		assert(t, err == nil, "fourth move should succeed: %v", err)
		assert(t, len(g.RemainingMoves) == 0, "should have 0 remaining after fourth move")
	}
}

// ---------------------------------------------------------------------------
// TestNoLegalMoves
// ---------------------------------------------------------------------------

func TestNoLegalMoves(t *testing.T) {
	g := &Game{Turn: White, State: Moving}
	g.WhiteBar = 1
	// Block all white entry points (indices 18-23).
	for i := 18; i <= 23; i++ {
		g.Board[i] = Point{Black, 2}
	}
	g.RemainingMoves = []int{1, 2}

	moves := g.AllLegalMoves()
	assert(t, len(moves) == 0, "should have no legal moves")

	err := g.PassTurn()
	assert(t, err == nil, "pass should succeed when no moves: %v", err)
	assert(t, g.Turn == Black, "should be black's turn after pass")
	assert(t, g.State == Rolling, "should be rolling state after pass")
}

func TestCannotPassWithLegalMoves(t *testing.T) {
	g := NewGame()
	g.SetDice(3, 1)

	err := g.PassTurn()
	assert(t, err != nil, "should not be able to pass when legal moves exist")
	assert(t, g.Turn == White, "should still be white's turn")
}

// ---------------------------------------------------------------------------
// TestMustUseBothDice
// ---------------------------------------------------------------------------

func TestMustUseBothDice(t *testing.T) {
	g := &Game{Turn: White, State: Moving}
	g.Board[5] = Point{White, 1}
	g.Board[4] = Point{White, 1}
	g.WhiteOff = 13
	g.RemainingMoves = []int{3, 5}

	// Both pieces in home board. Die 3: 5→2, 4→1. Die 5: 5→0.
	// All combinations use both dice:
	//   {5,2,3} → {4,off,5} (exact: dist=5)
	//   {4,1,3} → {5,0,5} (regular)
	//   {5,0,5} → {4,1,3} (regular)
	//   {4,off,5} → {5,2,3} (regular)
	moves := g.AllLegalMoves()
	assert(t, len(moves) == 4,
		"should have 4 moves when both dice are usable, got %d: %s", len(moves), printMoves(moves))

	has3 := false
	has5 := false
	for _, m := range moves {
		if m.Die == 3 {
			has3 = true
		}
		if m.Die == 5 {
			has5 = true
		}
	}
	assert(t, has3 && has5, "should include moves for both dice")
}

// ---------------------------------------------------------------------------
// TestMustUseLargerDie
// ---------------------------------------------------------------------------

func TestMustUseLargerDie(t *testing.T) {
	g := &Game{Turn: White, State: Moving}
	g.Board[10] = Point{White, 1}
	g.Board[1] = Point{Black, 2} // will block after first move
	g.WhiteOff = 14
	g.RemainingMoves = []int{3, 6}

	// Die 3: from 10→7. Then die 6: from 7→1. Blocked (2 black).
	// Die 6: from 10→4. Then die 3: from 4→1. Blocked.
	// Both paths: 1 move then blocked. maxDepth=1. Must use larger die (6).
	moves := g.AllLegalMoves()
	assert(t, len(moves) > 0, "should have moves for larger die")
	for _, m := range moves {
		assert(t, m.Die == 6,
			"must use larger die 6, got die %d (move: %+v)", m.Die, m)
	}
}

// ---------------------------------------------------------------------------
// TestWin
// ---------------------------------------------------------------------------

func TestWin(t *testing.T) {
	g := &Game{Turn: White, State: Moving}
	g.Board[0] = Point{White, 2} // 2 white on index 0 (dist 1)
	g.WhiteOff = 13
	g.RemainingMoves = []int{1, 1} // doubles of 1

	// Bear off both checkers with exact die 1.
	moves := g.AllLegalMoves()
	assert(t, len(moves) > 0, "should have moves")

	err := g.MakeMove(moves[0])
	assert(t, err == nil, "first bear off should succeed: %v", err)
	assert(t, g.WhiteOff == 14, "should have 14 off after first move")

	if g.State == GameOver {
		// Only 14 off, shouldn't be game over yet.
		t.Fatalf("game shouldn't be over at 14 off")
	}

	moves = g.AllLegalMoves()
	assert(t, len(moves) > 0, "should have second bear off move")
	err = g.MakeMove(moves[0])
	assert(t, err == nil, "second bear off should succeed: %v", err)
	assert(t, g.State == GameOver, "game should be over")
	assert(t, g.Winner == White, "white should win")
	assert(t, g.WhiteOff == 15, "white should have all 15 off")
}

// ---------------------------------------------------------------------------
// TestIllegalMove
// ---------------------------------------------------------------------------

func TestIllegalMove(t *testing.T) {
	g := NewGame()
	g.SetDice(3, 1)

	err := g.MakeMove(Move{From: 23, To: 20, Die: 6}) // die 6 not rolled
	assert(t, err != nil, "should reject move with wrong die")
	assert(t, g.State == Moving, "state should still be Moving")

	err = g.MakeMove(Move{From: 23, To: 18, Die: 5}) // die 5 not rolled
	assert(t, err != nil, "should reject move with unrolled die")
}

// ---------------------------------------------------------------------------
// TestTurnAlternation
// ---------------------------------------------------------------------------

func TestTurnAlternation(t *testing.T) {
	g := NewGame()
	g.SetDice(3, 1)
	assert(t, g.Turn == White, "white's turn")

	// Use all moves.
	for len(g.AllLegalMoves()) > 0 {
		moves := g.AllLegalMoves()
		g.MakeMove(moves[0])
	}

	assert(t, g.Turn == Black, "should be black's turn after white finishes")
	assert(t, g.State == Rolling, "should be rolling state")

	g.SetDice(2, 5)
	assert(t, g.Turn == Black, "still black's turn")
}

// ---------------------------------------------------------------------------
// TestBlackMovement
// ---------------------------------------------------------------------------

func TestBlackMovement(t *testing.T) {
	g := &Game{Turn: Black, State: Moving}
	g.Board[5] = Point{Black, 1}
	g.BlackOff = 14
	g.RemainingMoves = []int{3}

	// Black on index 5, die 3: target = 5+3 = 8.
	moves := g.AllLegalMoves()
	assert(t, len(moves) == 1, "black should have 1 move, got %d", len(moves))
	assert(t, moves[0].From == 5, "from should be 5")
	assert(t, moves[0].To == 8, "to should be 8 (black moves forward)")
}

// ---------------------------------------------------------------------------
// TestHit_byBlack
// ---------------------------------------------------------------------------

func TestHit_byBlack(t *testing.T) {
	g := &Game{Turn: Black, State: Moving}
	g.Board[5] = Point{Black, 1}
	g.Board[8] = Point{White, 1} // white blot
	g.BlackOff = 14
	g.WhiteOff = 14
	g.RemainingMoves = []int{3}

	moves := g.AllLegalMoves()
	found := false
	for _, m := range moves {
		if m.From == 5 && m.To == 8 {
			found = true
			err := g.MakeMove(m)
			assert(t, err == nil, "black hit should succeed: %v", err)
			break
		}
	}
	assert(t, found, "black should have hitting move")
	assert(t, g.WhiteBar == 1, "white should be on bar after black hits")
	assert(t, g.Board[8].Color == Black && g.Board[8].Count == 1,
		"point 8 should be black after hit")
}

// ---------------------------------------------------------------------------
// TestBarReentry_forBlack
// ---------------------------------------------------------------------------

func TestBarReentry_forBlack(t *testing.T) {
	g := &Game{Turn: Black, State: Moving}
	g.BlackBar = 1
	g.Board[0] = Point{Black, 14}
	g.RemainingMoves = []int{6}

	assert(t, g.MustEnterFromBar(), "black should need to enter from bar")

	// Black die 6: enters at index 6-1 = 5.
	moves := g.AllLegalMoves()
	assert(t, len(moves) == 1, "should have 1 bar entry move")
	assert(t, moves[0].From == BarPoint, "should enter from bar")
	assert(t, moves[0].To == 5, "black should enter at index 5")

	err := g.MakeMove(moves[0])
	assert(t, err == nil, "entry should succeed: %v", err)
	assert(t, g.BlackBar == 0, "black bar should be empty after entry")
}

// ---------------------------------------------------------------------------
// TestAI_RandomMove
// ---------------------------------------------------------------------------

func TestAI_RandomMove(t *testing.T) {
	g := NewGame()
	g.SetDice(3, 5)

	move := RandomMove(g)
	legal := g.AllLegalMoves()
	found := false
	for _, m := range legal {
		if m.From == move.From && m.To == move.To && m.Die == move.Die {
			found = true
			break
		}
	}
	assert(t, found, "random move should be legal: got %+v", move)
}

// ---------------------------------------------------------------------------
// TestAI_BestMove
// ---------------------------------------------------------------------------

func TestAI_BestMove(t *testing.T) {
	g := &Game{Turn: White, State: Moving}
	g.Board[10] = Point{White, 1}
	g.Board[3] = Point{White, 1}
	g.Board[7] = Point{Black, 1} // blot that can be hit
	g.WhiteOff = 13
	g.BlackOff = 14
	g.RemainingMoves = []int{3}

	best := BestMove(g)
	assert(t, best.From == 10 && best.To == 7,
		"AI should prefer hitting (from=10 to=7), got from=%d to=%d",
		best.From, best.To)
}

// ---------------------------------------------------------------------------
// TestEvaluate
// ---------------------------------------------------------------------------

func TestEvaluate(t *testing.T) {
	g := NewGame()

	whiteScore := Evaluate(g, White)
	blackScore := Evaluate(g, White)

	// Score should be deterministic.
	assert(t, whiteScore == blackScore, "eval should be deterministic for same state")

	// A position where all white pieces are borne off should score very high.
	g2 := &Game{}
	g2.WhiteOff = 15
	g2.Turn = White
	winScore := Evaluate(g2, White)
	assert(t, winScore > whiteScore,
		"winning position should score higher: win=%.2f start=%.2f", winScore, whiteScore)
}

// ---------------------------------------------------------------------------
// TestGamePlay_simple
// ---------------------------------------------------------------------------

func TestGamePlay_simple(t *testing.T) {
	g := NewGame()
	assert(t, g.State == Rolling, "initial state should be Rolling")

	// White rolls.
	g.SetDice(3, 1)
	assert(t, g.State == Moving, "state should be Moving after set dice")
	assert(t, len(g.RemainingMoves) == 2, "should have 2 remaining moves")

	// Make all white moves.
	moveCount := 0
	for len(g.AllLegalMoves()) > 0 {
		moves := g.AllLegalMoves()
		err := g.MakeMove(moves[0])
		assert(t, err == nil,
			"white move %d should succeed: %v", moveCount+1, err)
		moveCount++
	}
	t.Logf("White made %d move(s) with dice [3,1]", moveCount)

	assert(t, g.Turn == Black, "should be black's turn")
	assert(t, g.State == Rolling, "should be rolling for black")

	// Black rolls.
	g.SetDice(2, 4)
	moveCount = 0
	for len(g.AllLegalMoves()) > 0 {
		moves := g.AllLegalMoves()
		err := g.MakeMove(moves[0])
		assert(t, err == nil,
			"black move %d should succeed: %v", moveCount+1, err)
		moveCount++
	}
	t.Logf("Black made %d move(s) with dice [2,4]", moveCount)

	assert(t, g.Turn == White, "should be white's turn again")
	assert(t, g.TotalCheckers(White) == 15, "white should still have 15 checkers")
	assert(t, g.TotalCheckers(Black) == 15, "black should still have 15 checkers")
}

// ---------------------------------------------------------------------------
// TestGamePlay_withHit
// ---------------------------------------------------------------------------

func TestGamePlay_withHit(t *testing.T) {
	g := &Game{Turn: White, State: Moving}
	g.Board[10] = Point{White, 1}
	g.Board[7] = Point{Black, 1}
	g.WhiteOff = 14
	g.BlackOff = 14
	g.RemainingMoves = []int{3}

	// White hits black on point 7.
	moves := g.AllLegalMoves()
	hitMove := Move{}
	for _, m := range moves {
		if m.To == 7 {
			hitMove = m
			break
		}
	}
	assert(t, hitMove.From == 10, "hit move from should be 10")
	err := g.MakeMove(hitMove)
	assert(t, err == nil, "hit move should succeed: %v", err)
	assert(t, g.BlackBar == 1, "black on bar")

	// Now it's black's turn. Black must enter from bar.
	// endTurn was called by MakeMove since no remaining moves.
	assert(t, g.Turn == Black, "should be black's turn")
	assert(t, g.State == Rolling, "should be rolling")

	g.SetDice(2, 4)
	assert(t, g.MustEnterFromBar(), "black must enter from bar")

	moves = g.AllLegalMoves()
	for _, m := range moves {
		assert(t, m.From == BarPoint,
			"black must enter from bar, got from=%d", m.From)
		// Die 2: enters at index 1. Die 4: enters at index 3.
		if m.Die == 2 {
			assert(t, m.To == 1, "black die 2 enters at index 1, got %d", m.To)
		}
		if m.Die == 4 {
			assert(t, m.To == 3, "black die 4 enters at index 3, got %d", m.To)
		}
	}
}

// ---------------------------------------------------------------------------
// TestBearingOff_gameFlow
// ---------------------------------------------------------------------------

func TestBearingOff_gameFlow(t *testing.T) {
	g := &Game{Turn: White, State: Moving}
	// All white in home board: various positions.
	g.Board[0] = Point{White, 3}
	g.Board[1] = Point{White, 2}
	g.Board[2] = Point{White, 2}
	g.Board[3] = Point{White, 2}
	g.Board[4] = Point{White, 2}
	g.Board[5] = Point{White, 4}
	g.WhiteOff = 0
	g.BlackOff = 0
	g.BlackBar = 0
	g.RemainingMoves = []int{6, 6} // doubles

	assert(t, g.CanBearOff(White), "all white in home, should be able to bear off")

	moves := g.AllLegalMoves()
	assert(t, len(moves) > 0, "should have bear off moves")

	// Execute moves and verify pieces are borne off.
	initialOff := g.WhiteOff
	for len(g.AllLegalMoves()) > 0 {
		moves := g.AllLegalMoves()
		g.MakeMove(moves[0])
	}
	assert(t, g.WhiteOff > initialOff,
		"should have borne off some checkers (%d before, %d after)",
		initialOff, g.WhiteOff)
}

// ---------------------------------------------------------------------------
// TestMustUseMaxDice_doubles
// ---------------------------------------------------------------------------

func TestMustUseMaxDice_doubles(t *testing.T) {
	g := &Game{Turn: White, State: Moving}
	g.Board[5] = Point{White, 3}
	g.Board[1] = Point{White, 1}
	g.WhiteOff = 11
	g.RemainingMoves = []int{4, 4, 4, 4}

	// Die 4: from 5→1 (exact), from 1→-3 (bear off, over-die).
	// Can use all 4: 5→1, 5→1, 5→1, 1→off. That's 4 moves.
	// Actually: 3 on index 5. Move all to 1: 5→1, 5→1, 5→1 (3 moves).
	// Then 4 on index 1: bear off (over-die, farthest). 1 move. Total 4.
	moves := g.AllLegalMoves()
	assert(t, len(moves) > 0, "should have legal moves with doubles")

	// All first moves should use die 4 (only option).
	for _, m := range moves {
		assert(t, m.Die == 4, "all moves should use die 4, got %d", m.Die)
	}

	// Execute one move, verify 3 remaining.
	g.MakeMove(moves[0])
	assert(t, len(g.RemainingMoves) == 3, "should have 3 remaining after first move")
}

// ---------------------------------------------------------------------------
// TestPartialDiceSequence
// ---------------------------------------------------------------------------

func TestPartialDiceSequence(t *testing.T) {
	g := &Game{Turn: White, State: Moving}
	g.Board[3] = Point{White, 1}
	g.Board[1] = Point{White, 1}
	g.WhiteOff = 13
	g.RemainingMoves = []int{2, 4}

	// Die 2: 3→1, 1→-1 (bear off exact)
	// Die 4: 3→-1 (bear off exact), 1→-3 (over-die, check farthest)
	// After 3→1, die 4: 1→-3. dist=2, die=4>2. isFarthestInHome(White,1)? Check 2,3,4,5. None. Yes. Bear off.
	// After 1→off (die 2), die 4: 3→-1. dist=4, die=4. Exact. Bear off.
	// After 3→off (die 4), die 2: 1→-1. dist=2, die=2. Exact. Bear off.
	// After 1→-3 (die 4): dist=2, die=4>2. isFarthestInHome(White,1)? Check 2,3,4,5. Index 3 has piece. NO. Can't bear off from 1 with die 4.
	// So die 4 can't bear off from index 1 (index 3 is farther). Die 4 only from index 3.

	// All 3 first moves achieve 2 total:
	//   {3,1,2}→{1,off,4}: 2 moves ✓
	//   {1,off,2}→{3,off,4}: 2 moves ✓
	//   {3,off,4}→{1,off,2}: 2 moves ✓
	moves := g.AllLegalMoves()
	assert(t, len(moves) == 3,
		"should have 3 first moves (all achieve max 2), got %d: %v", len(moves), moves)
}

// ---------------------------------------------------------------------------
// Benchmark
// ---------------------------------------------------------------------------

func BenchmarkAllLegalMoves(b *testing.B) {
	g := NewGame()
	g.SetDice(3, 5)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.AllLegalMoves()
	}
}

func BenchmarkAllLegalMoves_doubles(b *testing.B) {
	g := NewGame()
	g.SetDice(4, 4)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.AllLegalMoves()
	}
}

func BenchmarkBestMove(b *testing.B) {
	g := NewGame()
	g.SetDice(3, 5)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BestMove(g)
	}
}

// Helper for pretty printing moves in test output.
func printMoves(moves []Move) string {
	s := "["
	for i, m := range moves {
		if i > 0 {
			s += ", "
		}
		s += fmt.Sprintf("{%d→%d d%d}", m.From, m.To, m.Die)
	}
	return s + "]"
}
