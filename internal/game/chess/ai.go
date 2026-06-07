package chess

import (
	"math"
	"math/rand"
	"time"
)

// Search statistics for analysis.
type SearchStats struct {
	Nodes    uint64
	TTHits   uint64
	Depth    int
	Duration time.Duration
}

var lastStats SearchStats

// Stats returns statistics from the last BestMove call.
func Stats() SearchStats { return lastStats }

// Piece values for evaluation.
const (
	PawnValue   = 100
	KnightValue = 320
	BishopValue = 330
	RookValue   = 500
	QueenValue  = 900
	KingValue   = 20000

	infScore = 1000000
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

// --- Piece-Square Tables ---

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

// Endgame king table: king should be active, go to center.
var kingEndgameTable = [8][8]int{
	{-50, -40, -30, -20, -20, -30, -40, -50},
	{-30, -20, -10, 0, 0, -10, -20, -30},
	{-30, -10, 20, 30, 30, 20, -10, -30},
	{-30, -10, 30, 40, 40, 30, -10, -30},
	{-30, -10, 30, 40, 40, 30, -10, -30},
	{-30, -10, 20, 30, 30, 20, -10, -30},
	{-30, -30, 0, 0, 0, 0, -30, -30},
	{-50, -30, -30, -30, -30, -30, -30, -50},
}

// pieceSquareValue returns the positional bonus for a piece at a given position.
func pieceSquareValue(pt PieceType, color Color, file, rank int, isEndgame bool) int {
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
		if isEndgame {
			table = &kingEndgameTable
		} else {
			table = &kingMiddleTable
		}
	default:
		return 0
	}
	if color == White {
		return table[7-rank][file]
	}
	return table[rank][file]
}

// --- Transposition Table ---

const ttSize = 1 << 20 // ~1M entries

const (
	ttExact = iota
	ttLower // score >= beta
	ttUpper // score <= alpha
)

type ttEntry struct {
	hash          uint32
	depth         int8
	score         int16
	flag          uint8
	bestFromSq    uint8 // packed: rank*8+file
	bestToSq      uint8
}

var tt [ttSize]ttEntry

func ttProbe(hash uint64) (ttEntry, bool) {
	idx := hash & (ttSize - 1)
	e := tt[idx]
	return e, e.hash == uint32(hash>>32)
}

func ttStore(hash uint64, depth int, score int, flag int, from, to Position) {
	idx := hash & (ttSize - 1)
	tt[idx] = ttEntry{
		hash:       uint32(hash >> 32),
		depth:      int8(depth),
		score:      int16(score),
		flag:       uint8(flag),
		bestFromSq: uint8(from.Rank*8 + from.File),
		bestToSq:   uint8(to.Rank*8 + to.File),
	}
}

func ttBestMove(hash uint64) (Position, Position, bool) {
	idx := hash & (ttSize - 1)
	e := tt[idx]
	if e.hash == uint32(hash>>32) {
		from := Position{int(e.bestFromSq % 8), int(e.bestFromSq / 8)}
		to := Position{int(e.bestToSq % 8), int(e.bestToSq / 8)}
		return from, to, true
	}
	return Position{}, Position{}, false
}

// --- Evaluation ---

// Evaluate returns the evaluation score for the current position.
// Positive = white better, negative = black better.
func (g *Game) Evaluate() int {
	isEndgame := g.isEndgame()
	score := 0
	whiteBishops := 0
	blackBishops := 0

	for r := 0; r < 8; r++ {
		for f := 0; f < 8; f++ {
			p := g.Board[r][f]
			if p.Empty() {
				continue
			}
			val := pieceValue(p.Type) + pieceSquareValue(p.Type, p.Color, f, r, isEndgame)
			if p.Color == White {
				score += val
				if p.Type == Bishop {
					whiteBishops++
				}
			} else {
				score -= val
				if p.Type == Bishop {
					blackBishops++
				}
			}
		}
	}

	// Bishop pair bonus
	if whiteBishops >= 2 {
		score += 30
	}
	if blackBishops >= 2 {
		score -= 30
	}

	// King safety (middlegame only)
	if !isEndgame {
		score += g.kingSafety(White) - g.kingSafety(Black)
	}

	// Game state
	if g.State == Checkmate {
		if g.Turn == White {
			return -KingValue
		}
		return KingValue
	}

	return score
}

// isEndgame returns true if the position is an endgame.
func (g *Game) isEndgame() bool {
	whiteMat := 0
	blackMat := 0
	for r := 0; r < 8; r++ {
		for f := 0; f < 8; f++ {
			p := g.Board[r][f]
			if p.Empty() || p.Type == King {
				continue
			}
			v := pieceValue(p.Type)
			if p.Color == White {
				whiteMat += v
			} else {
				blackMat += v
			}
		}
	}
	return whiteMat <= 1300 || blackMat <= 1300
}

// kingSafety evaluates pawn shield in front of the king.
func (g *Game) kingSafety(color Color) int {
	// Find king position
	kr, kf := -1, -1
	for r := 0; r < 8; r++ {
		for f := 0; f < 8; f++ {
			if g.Board[r][f].Type == King && g.Board[r][f].Color == color {
				kr, kf = r, f
				break
			}
		}
		if kr >= 0 {
			break
		}
	}
	if kr < 0 {
		return 0
	}

	score := 0
	// Pawn shield: pawns on rank+1 (or rank-1 for black) in front of king
	pawnDir := 1
	if color == Black {
		pawnDir = -1
	}
	shieldRank := kr + pawnDir
	if shieldRank >= 0 && shieldRank < 8 {
		for df := -1; df <= 1; df++ {
			sf := kf + df
			if sf >= 0 && sf < 8 {
				p := g.Board[shieldRank][sf]
				if p.Type == Pawn && p.Color == color {
					score += 10
				} else if p.Empty() {
					// Check second rank too
					secondRank := kr + 2*pawnDir
					if secondRank >= 0 && secondRank < 8 {
						p2 := g.Board[secondRank][sf]
						if p2.Type == Pawn && p2.Color == color {
							score += 5 // doubled pawn cover
						} else {
							score -= 3 // semi-open file near king
						}
					} else {
						score -= 3
					}
				}
			}
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

// --- SEE (Simple Static Exchange Evaluation) ---

func (g *Game) seeValue(pos Position, color Color) int {
	// Simple SEE: estimate if capturing at pos is winning
	p := g.Board[pos.Rank][pos.File]
	if p.Empty() {
		return 0
	}
	attackers := g.countAttackers(pos, color)
	if attackers == 0 {
		return 0
	}
	attackCost := g.cheapestAttackerValue(color)
	if attackCost == 0 {
		return 0
	}
	return pieceValue(p.Type) - attackCost
}

func (g *Game) countAttackers(pos Position, color Color) int {
	count := 0
	for r := 0; r < 8; r++ {
		for f := 0; f < 8; f++ {
			p := g.Board[r][f]
			if p.Empty() || p.Color != color {
				continue
			}
			if g.canAttack(Position{f, r}, pos, p) {
				count++
			}
		}
	}
	return count
}

func (g *Game) cheapestAttackerValue(color Color) int {
	cheapest := QueenValue + 1
	for r := 0; r < 8; r++ {
		for f := 0; f < 8; f++ {
			p := g.Board[r][f]
			if p.Empty() || p.Color != color || p.Type == King {
				continue
			}
			v := pieceValue(p.Type)
			if v < cheapest {
				cheapest = v
			}
		}
	}
	if cheapest == QueenValue+1 {
		return 0
	}
	return cheapest
}

func (g *Game) canAttack(from, to Position, piece Piece) bool {
	dr := to.Rank - from.Rank
	dc := to.File - from.File
	adr, adc := iabs(dr), iabs(dc)

	switch piece.Type {
	case Pawn:
		dir := 1
		if piece.Color == Black {
			dir = -1
		}
		return dr == dir && adc == 1
	case Knight:
		return (adr == 2 && adc == 1) || (adr == 1 && adc == 2)
	case Bishop:
		return adr == adc && adr > 0 && g.isPathClear(from, to)
	case Rook:
		return (dr == 0 || dc == 0) && (adr+adc > 0) && g.isPathClear(from, to)
	case Queen:
		return ((dr == 0 || dc == 0) || (adr == adc)) && (adr+adc > 0) && g.isPathClear(from, to)
	case King:
		return adr <= 1 && adc <= 1 && (adr+adc > 0)
	}
	return false
}

func (g *Game) isPathClear(from, to Position) bool {
	dr := to.Rank - from.Rank
	dc := to.File - from.File
	steps := imax(iabs(dr), iabs(dc))
	if steps <= 1 {
		return true
	}
	sr := isgn(dr)
	sc := isgn(dc)
	for i := 1; i < steps; i++ {
		if !g.Board[from.Rank+i*sr][from.File+i*sc].Empty() {
			return false
		}
	}
	return true
}

// --- Move Ordering ---

// scoreMove assigns a score to a move for ordering.
func (g *Game) scoreMove(m Move, pvMove Move) int {
	score := 0
	// PV move first
	if m.From == pvMove.From && m.To == pvMove.To && m.Promotion == pvMove.Promotion {
		return 1000000
	}
	// TT move
	if ttFrom, ttTo, ok := ttBestMove(g.ZobristHash); ok && m.From == ttFrom && m.To == ttTo {
		score += 500000
	}
	// Captures: MVV-LVA
	if m.Captured != nil {
		score += 10*pieceValue(m.Captured.Type) - pieceValue(m.Piece.Type)
	} else if m.IsEnPassant {
		score += 10*PawnValue - PawnValue
	}
	// Promotions
	if m.Promotion != NoPiece {
		score += pieceValue(m.Promotion)
	}
	// Castling
	if m.IsCastling {
		score += 50
	}
	// Center preference
	if m.To.File >= 2 && m.To.File <= 5 && m.To.Rank >= 2 && m.To.Rank <= 5 {
		score += 10
	}
	return score
}

func orderMoves(moves []Move, g *Game, pvMove Move) []Move {
	scored := make([]struct {
		move  Move
		score int
	}, len(moves))
	for i, m := range moves {
		scored[i].move = m
		scored[i].score = g.scoreMove(m, pvMove)
	}
	// Insertion sort
	for i := 1; i < len(scored); i++ {
		item := scored[i]
		j := i
		for j > 0 && scored[j-1].score < item.score {
			scored[j] = scored[j-1]
			j--
		}
		scored[j] = item
	}
	result := make([]Move, len(moves))
	for i, s := range scored {
		result[i] = s.move
	}
	return result
}

// --- Search ---

// RandomMove returns a random legal move.
func (g *Game) RandomMove() Move {
	moves := g.AllLegalMoves()
	if len(moves) == 0 {
		return Move{}
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return moves[rng.Intn(len(moves))]
}

// BestMoveFixed uses negamax with alpha-beta at a fixed depth (backward compatible).
func (g *Game) BestMoveFixed(depth int) Move {
	return g.bestMove(depth, Move{})
}

// BestMove uses iterative deepening with alpha-beta pruning.
func (g *Game) BestMove(maxDepth int, timeLimit time.Duration) Move {
	if g.ZobristHash == 0 {
		g.initZobristHash()
	}
	lastStats = SearchStats{Depth: maxDepth}
	start := time.Now()

	var bestMove Move
	for depth := 1; depth <= maxDepth; depth++ {
		move := g.bestMove(depth, bestMove)
		if move.From.Valid() {
			bestMove = move
		}
		lastStats.Depth = depth
		lastStats.Duration = time.Since(start)

		if timeLimit > 0 && time.Since(start) > timeLimit/2 && depth < maxDepth {
			break
		}
	}
	return bestMove
}

func (g *Game) bestMove(depth int, pvMove Move) Move {
	lastStats.Nodes = 0

	moves := g.AllLegalMoves()
	if len(moves) == 0 {
		return Move{}
	}

	moves = orderMoves(moves, g, pvMove)
	bestScore := -infScore
	bestMove := moves[0]
	alpha := -infScore
	beta := infScore

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

	ttStore(g.ZobristHash, depth, bestScore, ttExact, bestMove.From, bestMove.To)
	return bestMove
}

// negamax with alpha-beta, TT, quiescence, null move pruning, LMR, check extensions.
func (g *Game) negamax(depth int, alpha, beta int) int {
	lastStats.Nodes++

	if g.ZobristHash == 0 {
		g.initZobristHash()
	}

	// Leaf node: quiescence search
	if depth <= 0 {
		return g.quiescence(alpha, beta)
	}

	// Check for game over
	if g.State == Checkmate {
		return -KingValue
	}
	if g.State == Stalemate || g.State == Draw {
		return 0
	}

	// TT probe
	if entry, ok := ttProbe(g.ZobristHash); ok && entry.depth >= int8(depth) {
		score := int(entry.score)
		switch entry.flag {
		case ttExact:
			lastStats.TTHits++
			return score
		case ttLower:
			if score > alpha {
				alpha = score
			}
		case ttUpper:
			if score < beta {
				beta = score
			}
		}
		if alpha >= beta {
			lastStats.TTHits++
			return score
		}
	}

	moves := g.AllLegalMoves()
	if len(moves) == 0 {
		if g.isInCheck(g.Turn) {
			return -KingValue
		}
		return 0
	}

	// Get TT/PV best move for ordering
	pvFrom, pvTo, _ := ttBestMove(g.ZobristHash)
	pvMove := Move{From: pvFrom, To: pvTo}
	moves = orderMoves(moves, g, pvMove)

	// Check extension
	inCheck := g.isInCheck(g.Turn)
	if inCheck {
		depth++
	}

	bestScore := -infScore
	bestFlag := ttUpper

	for i, move := range moves {
		clone := g.Clone()
		clone.executeMove(move)

		// Null Move Pruning: skip for first 3 moves, when in check, endgame, or capturing
		if i >= 3 && !inCheck && depth >= 3 && !g.isEndgame() && move.Captured == nil && !move.IsEnPassant {
			nullScore := -clone.negamax(depth-3, -beta, -beta+1)
			if nullScore >= beta {
				lastStats.Nodes++
				return beta
			}
		}

		// Late Move Reductions
		newDepth := depth - 1
		if i >= 4 && !inCheck && depth >= 3 && move.Captured == nil && !move.IsEnPassant && move.Promotion == NoPiece {
			reducedScore := -clone.negamax(newDepth-1, -alpha-1, -alpha)
			if reducedScore <= alpha {
				lastStats.Nodes++
				continue // this move is probably bad
			}
		}

		score := -clone.negamax(newDepth, -beta, -alpha)
		if score > bestScore {
			bestScore = score
		}
		if score > alpha {
			alpha = score
			bestFlag = ttExact
		}
		if alpha >= beta {
			bestFlag = ttLower
			break
		}
	}

	ttStore(g.ZobristHash, depth, bestScore, bestFlag, moves[0].From, moves[0].To)
	return bestScore
}

// quiescence search: continue searching captures at depth 0.
func (g *Game) quiescence(alpha, beta int) int {
	lastStats.Nodes++
	standPat := g.evaluateForSide()

	if standPat >= beta {
		return beta
	}
	if standPat > alpha {
		alpha = standPat
	}

	// Generate capture moves only
	moves := g.AllLegalMoves()
	captures := make([]Move, 0, len(moves))
	for _, m := range moves {
		if m.Captured != nil || m.IsEnPassant || m.Promotion != NoPiece {
			captures = append(captures, m)
		}
	}

	captures = orderMoves(captures, g, Move{})

	for _, move := range captures {
		// Delta pruning: skip if capture can't improve alpha
		targetVal := PawnValue
		if move.Captured != nil {
			targetVal = pieceValue(move.Captured.Type)
		}
		if standPat+targetVal+200 < alpha {
			continue
		}

		clone := g.Clone()
		clone.executeMove(move)
		score := -clone.quiescence(-beta, -alpha)
		if score >= beta {
			return beta
		}
		if score > alpha {
			alpha = score
		}
	}

	return alpha
}

// --- Helpers ---

func iabs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func imax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func isgn(x int) int {
	if x > 0 {
		return 1
	}
	if x < 0 {
		return -1
	}
	return 0
}

// Verify we use math (imported for math.Inf in future use — but if unused, remove)
var _ = math.MaxInt
