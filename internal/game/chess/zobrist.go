package chess

import "math/rand"

// Zobrist hashing random numbers for incremental position hashing.
var (
	zobristPiece  [2][6][64]uint64 // [color-1][pieceType-1][squareIndex]
	zobristSide   uint64
	zobristCastle [4]uint64 // [0]=WhiteKingside, [1]=WhiteQueenside, [2]=BlackKingside, [3]=BlackQueenside
	zobristEP     [8]uint64 // [file 0..7]
)

func init() {
	rng := rand.New(rand.NewSource(42)) // fixed seed for reproducibility
	for c := 0; c < 2; c++ {
		for pt := 0; pt < 6; pt++ {
			for sq := 0; sq < 64; sq++ {
				zobristPiece[c][pt][sq] = rng.Uint64()
			}
		}
	}
	zobristSide = rng.Uint64()
	for i := 0; i < 4; i++ {
		zobristCastle[i] = rng.Uint64()
	}
	for i := 0; i < 8; i++ {
		zobristEP[i] = rng.Uint64()
	}
}

// computeZobristHash computes the full Zobrist hash from scratch.
func (g *Game) computeZobristHash() uint64 {
	var h uint64
	for r := 0; r < 8; r++ {
		for f := 0; f < 8; f++ {
			p := g.Board[r][f]
			if !p.Empty() {
				h ^= zobristPiece[p.Color-1][p.Type-1][r*8+f]
			}
		}
	}
	if g.Turn == Black {
		h ^= zobristSide
	}
	if g.CastlingRights.WhiteKingside {
		h ^= zobristCastle[0]
	}
	if g.CastlingRights.WhiteQueenside {
		h ^= zobristCastle[1]
	}
	if g.CastlingRights.BlackKingside {
		h ^= zobristCastle[2]
	}
	if g.CastlingRights.BlackQueenside {
		h ^= zobristCastle[3]
	}
	if g.EnPassantTarget != nil {
		h ^= zobristEP[g.EnPassantTarget.File]
	}
	return h
}

// initZobristHash computes and sets the Zobrist hash for the current position.
func (g *Game) initZobristHash() {
	g.ZobristHash = g.computeZobristHash()
}

// updateZobristForMove incrementally updates the Zobrist hash after a move.
// oldCastling and oldEP are the castling rights and en passant target BEFORE the move.
func (g *Game) updateZobristForMove(move Move, oldCastling CastlingRights, oldEP *Position) {
	h := g.ZobristHash

	fromSq := move.From.Rank*8 + move.From.File
	toSq := move.To.Rank*8 + move.To.File

	// XOR out moving piece from source
	h ^= zobristPiece[move.Piece.Color-1][move.Piece.Type-1][fromSq]

	// XOR out captured piece at destination
	if move.Captured != nil {
		h ^= zobristPiece[move.Captured.Color-1][move.Captured.Type-1][toSq]
	}

	// XOR in piece at destination
	if move.Promotion != NoPiece {
		h ^= zobristPiece[move.Piece.Color-1][move.Promotion-1][toSq]
	} else {
		h ^= zobristPiece[move.Piece.Color-1][move.Piece.Type-1][toSq]
	}

	// En passant capture: XOR out the captured pawn
	if move.IsEnPassant {
		epSq := move.From.Rank*8 + move.To.File
		oppColor := move.Piece.Color.Opposite()
		h ^= zobristPiece[oppColor-1][0][epSq] // Pawn is pieceType 1, index 0
	}

	// Castling rook movement
	if move.IsCastling {
		rc := move.Piece.Color - 1
		rookIdx := 3 // Rook is pieceType 4, index 3
		if move.CastlingSide == "K" {
			h ^= zobristPiece[rc][rookIdx][move.From.Rank*8+7]
			h ^= zobristPiece[rc][rookIdx][move.From.Rank*8+5]
		} else {
			h ^= zobristPiece[rc][rookIdx][move.From.Rank*8+0]
			h ^= zobristPiece[rc][rookIdx][move.From.Rank*8+3]
		}
	}

	// Toggle side to move
	h ^= zobristSide

	// Castling rights changes
	if oldCastling.WhiteKingside != g.CastlingRights.WhiteKingside {
		h ^= zobristCastle[0]
	}
	if oldCastling.WhiteQueenside != g.CastlingRights.WhiteQueenside {
		h ^= zobristCastle[1]
	}
	if oldCastling.BlackKingside != g.CastlingRights.BlackKingside {
		h ^= zobristCastle[2]
	}
	if oldCastling.BlackQueenside != g.CastlingRights.BlackQueenside {
		h ^= zobristCastle[3]
	}

	// En passant file changes
	if oldEP != nil {
		h ^= zobristEP[oldEP.File]
	}
	if g.EnPassantTarget != nil {
		h ^= zobristEP[g.EnPassantTarget.File]
	}

	g.ZobristHash = h
}
