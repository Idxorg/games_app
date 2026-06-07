package checkers

import (
	"math"
	"math/rand"
)

// ---------------------------------------------------------------------------
// Zobrist Hashing
// ---------------------------------------------------------------------------

// zobristPiece[Man-1=0 / King-1=1][White=0 / Black=1][row][col]
var zobristPiece [2][2][8][8]uint64

// zobristSide is XORed when it is Black's turn to move.
var zobristSide uint64

func init() {
	r := rand.New(rand.NewSource(42)) // deterministic seed for reproducibility
	for pt := 0; pt < 2; pt++ {
		for c := 0; c < 2; c++ {
			for row := 0; row < 8; row++ {
				for col := 0; col < 8; col++ {
					zobristPiece[pt][c][row][col] = r.Uint64()
				}
			}
		}
	}
	zobristSide = r.Uint64()
}

// computeHash builds a Zobrist hash from scratch (used once in NewGame).
func (g *Game) computeHash() uint64 {
	var h uint64
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			p := g.Board[row][col]
			if !p.Empty() {
				h ^= zobristPiece[p.Type-1][p.Color][row][col]
			}
		}
	}
	if g.Turn == Black {
		h ^= zobristSide
	}
	return h
}

// ---------------------------------------------------------------------------
// Transposition Table (always-replace, 100 K entries)
// ---------------------------------------------------------------------------

const ttSize = 100000

type ttFlag int

const (
	ttExact ttFlag = iota // exact score
	ttAlpha               // upper bound (fail-low)
	ttBeta                // lower bound (fail-high)
)

type ttEntry struct {
	hash     uint32 // upper 32 bits of the full 64-bit hash for validation
	depth    int8
	score    float64
	flag     ttFlag
	bestMove Move
}

var ttTable [ttSize]ttEntry

func ttProbe(hash uint64) (*ttEntry, bool) {
	idx := uint32(hash) % ttSize
	e := &ttTable[idx]
	if e.hash == uint32(hash>>32) {
		return e, true
	}
	return e, false
}

func ttStore(hash uint64, depth int, score float64, flag ttFlag, best Move) {
	idx := uint32(hash) % ttSize
	ttTable[idx] = ttEntry{
		hash:     uint32(hash >> 32),
		depth:    int8(depth),
		score:    score,
		flag:     flag,
		bestMove: best,
	}
}

// ---------------------------------------------------------------------------
// Public counters / API
// ---------------------------------------------------------------------------

// NodesSearched is the number of nodes evaluated in the most recent BestMove call.
var NodesSearched uint64

// ---------------------------------------------------------------------------
// RandomMove — unchanged
// ---------------------------------------------------------------------------

// RandomMove returns a random legal move.
func (g *Game) RandomMove() Move {
	moves := g.AllLegalMoves()
	if len(moves) == 0 {
		return Move{}
	}
	return moves[rand.Intn(len(moves))]
}

// ---------------------------------------------------------------------------
// BestMove — iterative deepening with alpha-beta
// ---------------------------------------------------------------------------

const maxSearchDepth = 64 // hard cap to prevent runaway extensions

// BestMove returns the best move using iterative deepening with alpha-beta pruning.
func (g *Game) BestMove(maxDepth int) Move {
	moves := g.AllLegalMoves()
	if len(moves) == 0 {
		return Move{}
	}

	NodesSearched = 0
	var bestMove Move

	for depth := 1; depth <= maxDepth; depth++ {
		ordered := orderMoves(moves, bestMove)

		best := math.Inf(-1)
		var cur Move

		for _, m := range ordered {
			ng := g.clone()
			if err := ng.MakeMove(m); err != nil {
				continue
			}
			s := -ng.negamax(depth-1, math.Inf(-1), -best)
			if s > best {
				best = s
				cur = m
			}
		}

		if !isZeroMove(cur) {
			bestMove = cur
		}
	}

	return bestMove
}

// isZeroMove returns true if the move is the zero value.
func isZeroMove(m Move) bool {
	return m.From.Row == 0 && m.From.Col == 0 &&
		m.To.Row == 0 && m.To.Col == 0 &&
		m.Piece.Type == None
}

// ---------------------------------------------------------------------------
// Negamax with alpha-beta, TT, move ordering, check extensions, quiescence
// ---------------------------------------------------------------------------

func (g *Game) negamax(depth int, alpha, beta float64) float64 {
	NodesSearched++

	// --- Terminal state ---
	switch g.State {
	case WhiteWin:
		if g.Turn == White {
			return 10000
		}
		return -10000
	case BlackWin:
		if g.Turn == Black {
			return 10000
		}
		return -10000
	case Draw:
		return 0
	}

	// --- Transposition table probe ---
	hash := g.ZobristHash
	var ttBest Move
	if entry, ok := ttProbe(hash); ok {
		ttBest = entry.bestMove
		if entry.depth >= int8(depth) {
			switch entry.flag {
			case ttExact:
				return entry.score
			case ttAlpha:
				if entry.score <= alpha {
					return alpha
				}
			case ttBeta:
				if entry.score >= beta {
					return beta
				}
			}
		}
	}

	// --- Check extension: multi-jump captures extend search by 1 ply ---
	captures := g.allCaptures()
	multiJump := false
	for _, m := range captures {
		if len(m.Captured) > 1 {
			multiJump = true
			break
		}
	}
	if multiJump && depth >= 1 && depth < maxSearchDepth {
		depth++
	}

	// --- Leaf → quiescence search ---
	if depth <= 0 {
		return g.quiescence(alpha, beta)
	}

	// --- Generate & order moves ---
	moves := g.AllLegalMoves()
	if len(moves) == 0 {
		return -10000 // no legal moves = loss for side to move
	}
	moves = orderMoves(moves, ttBest)

	origAlpha := alpha
	bestScore := math.Inf(-1)
	var bestMove Move

	for _, m := range moves {
		ng := g.clone()
		if err := ng.MakeMove(m); err != nil {
			continue
		}
		s := -ng.negamax(depth-1, -beta, -alpha)
		if s > bestScore {
			bestScore = s
			bestMove = m
		}
		if s > alpha {
			alpha = s
		}
		if alpha >= beta {
			break
		}
	}

	// --- Transposition table store ---
	var flag ttFlag
	switch {
	case bestScore <= origAlpha:
		flag = ttAlpha
	case bestScore >= beta:
		flag = ttBeta
	default:
		flag = ttExact
	}
	ttStore(hash, depth, bestScore, flag, bestMove)

	return bestScore
}

// ---------------------------------------------------------------------------
// Quiescence search — continue through forced captures only
// ---------------------------------------------------------------------------

func (g *Game) quiescence(alpha, beta float64) float64 {
	NodesSearched++

	// Terminal
	switch g.State {
	case WhiteWin:
		if g.Turn == White {
			return 10000
		}
		return -10000
	case BlackWin:
		if g.Turn == Black {
			return 10000
		}
		return -10000
	case Draw:
		return 0
	}

	// If no forced captures, evaluate statically
	captures := g.allCaptures()
	if len(captures) == 0 {
		return g.evalSide()
	}

	// In checkers, captures are mandatory — search all of them
	captures = orderMoves(captures, Move{})

	bestScore := math.Inf(-1)
	for _, m := range captures {
		ng := g.clone()
		if err := ng.MakeMove(m); err != nil {
			continue
		}
		s := -ng.quiescence(-beta, -alpha)
		if s > bestScore {
			bestScore = s
		}
		if s > alpha {
			alpha = s
		}
		if alpha >= beta {
			break
		}
	}
	return bestScore
}

// evalSide returns evaluation from the side-to-move's perspective.
func (g *Game) evalSide() float64 {
	if g.Turn == White {
		return g.evaluate()
	}
	return -g.evaluate()
}

// ---------------------------------------------------------------------------
// Move ordering — PV move > captures (by count) > promotions > center moves
// ---------------------------------------------------------------------------

func orderMoves(moves []Move, pv Move) []Move {
	type scored struct {
		m Move
		s int
	}
	sm := make([]scored, len(moves))
	for i, m := range moves {
		sc := 0
		// PV / TT best move gets highest priority
		if m.From == pv.From && m.To == pv.To && len(m.Captured) == len(pv.Captured) {
			sc += 100000
		}
		// Captures: more captured pieces = higher priority
		if m.IsJump {
			sc += 10000 + len(m.Captured)*100
		}
		// Promotions
		if m.Promoted {
			sc += 5000
		}
		// Center preference (closer to center = better)
		cd := math.Abs(float64(m.To.Row)-3.5) + math.Abs(float64(m.To.Col)-3.5)
		sc += int((7.0 - cd) * 10)
		sm[i] = scored{m, sc}
	}
	// Insertion sort (fast for small N, typical checkers: 5-20 moves)
	for i := 1; i < len(sm); i++ {
		for j := i; j > 0 && sm[j].s > sm[j-1].s; j-- {
			sm[j], sm[j-1] = sm[j-1], sm[j]
		}
	}
	result := make([]Move, len(moves))
	for i, s := range sm {
		result[i] = s.m
	}
	return result
}

// ---------------------------------------------------------------------------
// Evaluation — material, advancement, center, back row, trapped kings,
//               mobility, endgame knowledge
// ---------------------------------------------------------------------------

// evaluate returns a score from White's perspective.
// Positive = good for White, negative = good for Black.
func (g *Game) evaluate() float64 {
	score := 0.0

	whiteMen, whiteKings := 0, 0
	blackMen, blackKings := 0, 0
	var wkPos, bkPos Position

	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			p := g.Board[row][col]
			if p.Empty() {
				continue
			}

			val := 0.0

			if p.Type == King {
				val = 150

				if p.Color == White {
					whiteKings++
					wkPos = Position{row, col}
				} else {
					blackKings++
					bkPos = Position{row, col}
				}

				// Center control for kings
				cd := math.Abs(float64(row)-3.5) + math.Abs(float64(col)-3.5)
				val += (4.0 - cd) * 5

				// Trapped king penalty
				val -= trappedKingPenalty(row, col, g)
			} else {
				val = 100

				if p.Color == White {
					whiteMen++
				} else {
					blackMen++
				}

				// Advancement bonus (closer to promotion = better)
				if p.Color == White {
					val += float64(7-row) * 5
				} else {
					val += float64(row) * 5
				}

				// Back row defense bonus (prevents easy king captures)
				if p.Color == White && row == 7 {
					val += 15
				}
				if p.Color == Black && row == 0 {
					val += 15
				}

				// Center control for men
				cd := math.Abs(float64(row)-3.5) + math.Abs(float64(col)-3.5)
				val += (4.0 - cd) * 2
			}

			if p.Color == White {
				score += val
			} else {
				score -= val
			}
		}
	}

	// --- Mobility: number of legal moves available ---
	saved := g.Turn
	g.Turn = White
	wm := len(g.AllLegalMoves())
	g.Turn = Black
	bm := len(g.AllLegalMoves())
	g.Turn = saved
	score += float64(wm-bm) * 3

	// --- Endgame knowledge (only when few pieces remain) ---
	totalPieces := whiteMen + whiteKings + blackMen + blackKings
	if totalPieces <= 6 {
		score += endgameBonus(whiteMen, whiteKings, blackMen, blackKings, wkPos, bkPos)
	}

	// --- Game state bonuses ---
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

// trappedKingPenalty returns a penalty for a king stuck on the edge/corner
// with few or no moves.
func trappedKingPenalty(row, col int, g *Game) float64 {
	pen := 0.0

	// Corner is worst, edge is bad
	if (row == 0 || row == 7) && (col == 0 || col == 7) {
		pen = 20
	} else if row == 0 || row == 7 || col == 0 || col == 7 {
		pen = 10
	}

	// Check actual mobility of this king
	pos := Position{row, col}
	p := g.Board[row][col]
	n := len(g.kingSimpleMoves(pos, p))
	if n == 0 {
		pen += 30
	} else if n == 1 {
		pen += 15
	}

	return pen
}

// endgameBonus adds endgame-specific evaluation bonuses.
func endgameBonus(wm, wk, bm, bk int, wkPos, bkPos Position) float64 {
	s := 0.0
	wt := wm + wk
	bt := bm + bk

	// Material advantage is amplified in endgame
	diff := wt - bt
	s += float64(diff) * 50

	// 1 king vs 1 king: opposition / center distance
	if wt == 1 && bt == 1 && wk == 1 && bk == 1 {
		wd := math.Abs(float64(wkPos.Row)-3.5) + math.Abs(float64(wkPos.Col)-3.5)
		bd := math.Abs(float64(bkPos.Row)-3.5) + math.Abs(float64(bkPos.Col)-3.5)
		// Closer to center is better
		s += (bd - wd) * 10
		if wd < bd {
			s += 20 // explicit opposition bonus
		}
	}

	// 2 pieces vs 1 piece: strong advantage
	if (wt == 2 && bt == 1) || (wt == 1 && bt == 2) {
		s += float64(wt-bt) * 80
	}

	return s
}

// ---------------------------------------------------------------------------
// Clone — deep copy with Zobrist hash
// ---------------------------------------------------------------------------

// clone creates a deep copy of the game.
func (g *Game) clone() *Game {
	c := &Game{
		Board:          g.Board,
		Turn:           g.Turn,
		State:          g.State,
		KingMoveCount:  g.KingMoveCount,
		TotalKingMoves: g.TotalKingMoves,
		Moves:          make([]Move, len(g.Moves)),
		ZobristHash:    g.ZobristHash,
	}
	copy(c.Moves, g.Moves)
	return c
}
