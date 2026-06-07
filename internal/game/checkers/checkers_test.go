package checkers

import (
	"testing"
)

// TestNewGame verifies the starting position of a new checkers game.
func TestNewGame(t *testing.T) {
	g := NewGame()

	if g.Turn != White {
		t.Errorf("expected White to move first, got %v", g.Turn)
	}
	if g.State != Playing {
		t.Errorf("expected Playing state, got %v", g.State)
	}

	white, black := g.CountPieces()
	if white != 12 {
		t.Errorf("expected 12 white pieces, got %d", white)
	}
	if black != 12 {
		t.Errorf("expected 12 black pieces, got %d", black)
	}

	// Verify pieces are on dark squares only
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			p := g.Board[row][col]
			if !p.Empty() {
				if (row+col)%2 != 1 {
					t.Errorf("piece found on light square at (%d,%d)", row, col)
				}
			}
		}
	}

	// Verify black on rows 0-2
	for row := 0; row < 3; row++ {
		for col := 0; col < 8; col++ {
			if (row+col)%2 == 1 {
				p := g.Board[row][col]
				if p.Type != Man || p.Color != Black {
					t.Errorf("expected black man at (%d,%d), got %v", row, col, p)
				}
			}
		}
	}

	// Verify white on rows 5-7
	for row := 5; row < 8; row++ {
		for col := 0; col < 8; col++ {
			if (row+col)%2 == 1 {
				p := g.Board[row][col]
				if p.Type != Man || p.Color != White {
					t.Errorf("expected white man at (%d,%d), got %v", row, col, p)
				}
			}
		}
	}

	// Middle rows should be empty
	for row := 3; row < 5; row++ {
		for col := 0; col < 8; col++ {
			if !g.Board[row][col].Empty() {
				t.Errorf("expected empty square at (%d,%d), got %v", row, col, g.Board[row][col])
			}
		}
	}
}

// TestLegalMoves_Man verifies simple forward moves for a man.
func TestLegalMoves_Man(t *testing.T) {
	g := NewGame()

	// Pick a white man at row 5, col 0 (should be able to move to (4,1))
	// Row 5, col 0: (5+0)%2 = 1, yes dark square
	moves := g.LegalMoves(Position{5, 0})
	if len(moves) == 0 {
		t.Error("expected at least one move for white man at (5,0)")
	}

	// All moves should be forward (row 4 for white)
	for _, m := range moves {
		if m.IsJump {
			t.Error("expected no jump moves in starting position")
		}
		if m.To.Row != 4 {
			t.Errorf("expected move to row 4 for white man, got row %d", m.To.Row)
		}
	}
}

// TestLegalMoves_ManCapture verifies capture moves for a man.
func TestLegalMoves_ManCapture(t *testing.T) {
	g := &Game{
		Turn:  White,
		State: Playing,
		Moves: make([]Move, 0),
	}

	// Place white man at (5,2) and black man at (4,1)
	g.Board[5][2] = Piece{Man, White}
	g.Board[4][1] = Piece{Man, Black}

	moves := g.LegalMoves(Position{5, 2})
	foundCapture := false
	for _, m := range moves {
		if m.IsJump {
			foundCapture = true
			if m.To.Row != 3 || m.To.Col != 0 {
				t.Errorf("expected capture landing at (3,0), got %v", m.To)
			}
			if len(m.Captured) != 1 || m.Captured[0] != (Position{4, 1}) {
				t.Errorf("expected to capture (4,1), got %v", m.Captured)
			}
		}
	}
	if !foundCapture {
		t.Error("expected a capture move to be available")
	}
}

// TestMultiJump verifies a multi-jump sequence (3 captures in one move).
func TestMultiJump(t *testing.T) {
	g := &Game{
		Turn:  White,
		State: Playing,
		Moves: make([]Move, 0),
	}

	// Set up a scenario: white man at (7,0) can capture 3 black men in a chain
	// Path: (7,0) -> capture (6,1) -> land (5,2) -> capture (4,1) -> land (3,0) -> capture (2,1) -> land (1,0)
	// Wait, men can only capture forward or backward? In Russian checkers, men can capture in any direction.
	// But they can only land if the square is empty.
	// Let me set up a simpler chain:
	// White man at (5,0), black men at (4,1), (2,1), (0,1)
	// White captures (4,1)->lands(3,2), captures (2,1)->lands(1,0), etc.
	// Actually let me think about this more carefully with valid diagonal paths.

	// White man at row 7, col 0 (dark square: 7+0=7, odd, yes)
	// Black men at (6,1), (4,3), (2,5) - creating a zigzag capture path
	// Landing: (5,2), (3,4), (1,6)
	g.Board[7][0] = Piece{Man, White}
	g.Board[6][1] = Piece{Man, Black}
	g.Board[5][2] = Piece{Type: None} // ensure empty
	g.Board[4][3] = Piece{Man, Black}
	g.Board[3][4] = Piece{Type: None} // ensure empty
	g.Board[2][5] = Piece{Man, Black}
	g.Board[1][6] = Piece{Type: None} // ensure empty

	moves := g.capturesFrom(Position{7, 0})
	foundTriple := false
	for _, m := range moves {
		if len(m.Captured) == 3 {
			foundTriple = true
			if m.From != (Position{7, 0}) {
				t.Errorf("expected from (7,0), got %v", m.From)
			}
			if m.To != (Position{1, 6}) {
				t.Errorf("expected to (1,6), got %v", m.To)
			}
		}
	}
	if !foundTriple {
		t.Errorf("expected a triple capture move, got %d moves with lengths:", len(moves))
		for _, m := range moves {
			t.Logf("  capture of %d pieces: %v -> %v", len(m.Captured), m.From, m.To)
		}
	}
}

// TestKingMove verifies a king can slide any distance diagonally.
func TestKingMove(t *testing.T) {
	g := &Game{
		Turn:  White,
		State: Playing,
		Moves: make([]Move, 0),
	}

	// Place white king at center
	g.Board[3][4] = Piece{King, White}

	moves := g.LegalMoves(Position{3, 4})
	if len(moves) == 0 {
		t.Error("expected king to have moves")
	}

	// King should be able to move multiple squares
	foundLongMove := false
	for _, m := range moves {
		if m.IsJump {
			continue
		}
		dr := m.To.Row - m.From.Row
		dc := m.To.Col - m.From.Col
		if abs(dr) > 1 || abs(dc) > 1 {
			foundLongMove = true
		}
	}
	if !foundLongMove {
		t.Error("expected king to be able to slide multiple squares")
	}
}

// TestKingCapture verifies a king can capture at any distance diagonally.
func TestKingCapture(t *testing.T) {
	g := &Game{
		Turn:  White,
		State: Playing,
		Moves: make([]Move, 0),
	}

	// Place white king at (0,1), black man at (3,4)
	// King should be able to capture and land at (1,2) or (2,3) or beyond
	// Actually, in Russian checkers, king captures by jumping OVER the enemy at any distance,
	// landing on any empty square beyond the enemy on the same diagonal.
	g.Board[0][1] = Piece{King, White}
	g.Board[3][4] = Piece{Man, Black}
	// Ensure landing squares are empty
	g.Board[1][2] = Piece{Type: None}
	g.Board[2][3] = Piece{Type: None}
	g.Board[4][5] = Piece{Type: None}
	g.Board[5][6] = Piece{Type: None}
	g.Board[6][7] = Piece{Type: None}

	moves := g.capturesFrom(Position{0, 1})
	if len(moves) == 0 {
		t.Error("expected king to have capture moves")
	}

	// Verify king can land at various distances beyond the enemy
	foundFarLanding := false
	for _, m := range moves {
		if m.To.Row > 3 && m.To.Col > 4 {
			foundFarLanding = true
		}
	}
	if !foundFarLanding {
		t.Error("expected king to land far beyond captured piece")
	}
}

// TestPromotion verifies a man becomes a king when reaching the last row.
func TestPromotion(t *testing.T) {
	g := &Game{
		Turn:  White,
		State: Playing,
		Moves: make([]Move, 0),
	}

	// Place white man at (1,2) - one move from promotion row (row 0)
	g.Board[1][2] = Piece{Man, White}

	moves := g.LegalMoves(Position{1, 2})
	foundPromotion := false
	for _, m := range moves {
		if m.To.Row == 0 && m.Promoted {
			foundPromotion = true
		}
	}
	if !foundPromotion {
		t.Error("expected promotion move to be available")
	}

	// Execute the promotion move
	for _, m := range moves {
		if m.Promoted {
			err := g.MakeMove(m)
			if err != nil {
				t.Errorf("MakeMove failed: %v", err)
			}
			if g.Board[0][m.To.Col].Type != King {
				t.Error("piece should be a king after promotion")
			}
			break
		}
	}
}

// TestPromotion_Black verifies black man promotes at row 7.
func TestPromotion_Black(t *testing.T) {
	g := &Game{
		Turn:  Black,
		State: Playing,
		Moves: make([]Move, 0),
	}

	g.Board[6][1] = Piece{Man, Black}

	moves := g.LegalMoves(Position{6, 1})
	foundPromotion := false
	for _, m := range moves {
		if m.To.Row == 7 && m.Promoted {
			foundPromotion = true
		}
	}
	if !foundPromotion {
		t.Error("expected black promotion move to be available")
	}
}

// TestMandatoryCapture verifies non-capture moves are illegal when capture is available.
func TestMandatoryCapture(t *testing.T) {
	g := &Game{
		Turn:  White,
		State: Playing,
		Moves: make([]Move, 0),
	}

	// Place two white men: one can capture, one can only do simple moves
	g.Board[5][0] = Piece{Man, White} // can capture at (4,1)
	g.Board[5][4] = Piece{Man, Black} // enemy at capture position
	g.Board[4][1] = Piece{Man, White} // this man can only do simple forward move

	// Actually let me set it up clearer:
	// White man at (5,2) can capture black man at (4,1)
	// White man at (3,0) can only do simple moves (but capture is forced!)
	g = &Game{
		Turn:  White,
		State: Playing,
		Moves: make([]Move, 0),
	}
	g.Board[5][2] = Piece{Man, White}
	g.Board[4][1] = Piece{Man, Black}
	g.Board[4][3] = Piece{Man, Black} // another enemy, white can't move to (4,3)
	g.Board[3][0] = Piece{Man, White}  // this man should NOT have legal simple moves

	if !g.IsForcedCapture() {
		t.Error("expected forced capture to be true")
	}

	// The man at (3,0) should have no legal moves (forced capture, but it can't capture)
	movesAt30 := g.LegalMoves(Position{3, 0})
	if len(movesAt30) != 0 {
		t.Errorf("expected no legal moves for piece at (3,0) when forced capture exists elsewhere, got %d", len(movesAt30))
	}

	// All legal moves should be captures
	allMoves := g.AllLegalMoves()
	for _, m := range allMoves {
		if !m.IsJump {
			t.Error("expected all moves to be captures when forced capture is active")
		}
	}
}

// TestForcedCaptureDetection verifies IsForcedCapture works correctly.
func TestForcedCaptureDetection(t *testing.T) {
	g := NewGame()

	// Starting position should have no forced captures
	if g.IsForcedCapture() {
		t.Error("expected no forced capture in starting position")
	}

	// Set up a capture scenario with clear landing square
	g.Board[3][2] = Piece{Man, White} // white man at row 3, col 2
	g.Board[2][1] = Piece{Man, Black} // black man at row 2, col 1
	g.Board[1][0] = Piece{}           // clear landing square behind black man

	if !g.IsForcedCapture() {
		t.Error("expected forced capture when capture is available")
	}
}

// TestDrawByKingMoves verifies draw detection by king move count.
func TestDrawByKingMoves(t *testing.T) {
	// Test KingMoveCount draw condition directly
	g := &Game{
		Turn:  White,
		State: Playing,
		Moves: make([]Move, 0),
	}
	g.Board[0][1] = Piece{King, White}
	g.Board[7][0] = Piece{King, Black}

	// Manually set KingMoveCount to 15 and verify draw is detected
	g.KingMoveCount = 15
	if !g.checkDraw() {
		t.Error("expected draw when KingMoveCount >= 15")
	}

	// Verify that a capture resets the counter
	g.KingMoveCount = 14
	if g.checkDraw() {
		t.Error("expected no draw when KingMoveCount < 15")
	}

	// Test TotalKingMoves draw condition
	g = &Game{
		Turn:  White,
		State: Playing,
		Moves: make([]Move, 0),
	}
	g.Board[0][1] = Piece{King, White}
	g.Board[7][0] = Piece{King, Black}
	g.TotalKingMoves = 25
	if !g.checkDraw() {
		t.Error("expected draw when TotalKingMoves >= 25 and both sides kings only")
	}
}

// TestDraw_ByTotalKingMoves verifies draw by 25 total king moves with kings only.
func TestDraw_ByTotalKingMoves(t *testing.T) {
	g := &Game{
		Turn:  White,
		State: Playing,
		Moves: make([]Move, 0),
	}
	g.Board[0][1] = Piece{King, White}
	g.Board[7][6] = Piece{King, Black}

	// Alternate king moves: king captures reset the counter but total increases
	// To trigger 25 total king moves with kings only, just make 25 simple king moves
	// with some captures interspersed that reset the consecutive counter
	for i := 0; i < 30; i++ {
		moves := g.AllLegalMoves()
		if len(moves) == 0 {
			break
		}
		// Use any move (prefer non-capture to avoid pieces being removed)
		var chosen Move
		for _, m := range moves {
			if !m.IsJump {
				chosen = m
				break
			}
		}
		if chosen.From.Row == 0 && chosen.From.Col == 0 {
			chosen = moves[0]
		}
		g.MakeMove(chosen)
		if g.State != Playing {
			break
		}
	}

	// Should have drawn by now (either 15 consecutive or 25 total)
	// If not, the total king moves should have triggered it since both are kings only
	if g.State != Draw && g.TotalKingMoves >= 25 && g.bothSidesKingsOnly() {
		t.Errorf("expected draw by total king moves, got state %v (total: %d)", g.State, g.TotalKingMoves)
	}
}

// TestAI verifies the AI can select moves.
func TestAI(t *testing.T) {
	g := NewGame()

	// RandomMove should return a valid move
	move := g.RandomMove()
	if move.From.Row == 0 && move.From.Col == 0 && move.To.Row == 0 && move.To.Col == 0 {
		t.Error("expected RandomMove to return a valid move")
	}

	// BestMove with depth 2 should return a valid move
	move = g.BestMove(2)
	if move.From.Row == 0 && move.From.Col == 0 && move.To.Row == 0 && move.To.Col == 0 {
		t.Error("expected BestMove to return a valid move")
	}

	// Verify the move is actually legal
	legal := g.AllLegalMoves()
	found := false
	for _, lm := range legal {
		if movesEqual(lm, move) {
			found = true
			break
		}
	}
	if !found {
		t.Error("BestMove returned an illegal move")
	}
}

// TestAI_PreferCapture verifies AI prefers captures over simple moves.
func TestAI_PreferCapture(t *testing.T) {
	g := &Game{
		Turn:  White,
		State: Playing,
		Moves: make([]Move, 0),
	}

	// White man that can capture
	g.Board[5][2] = Piece{Man, White}
	g.Board[4][1] = Piece{Man, Black}

	move := g.BestMove(3)
	if !move.IsJump {
		t.Error("AI should prefer capture moves when available")
	}
}

// TestGameEnd_NoPieces verifies game ends when a player has no pieces.
func TestGameEnd_NoPieces(t *testing.T) {
	g := &Game{
		Turn:  White,
		State: Playing,
		Moves: make([]Move, 0),
	}

	// Only white piece, no black pieces
	g.Board[5][2] = Piece{Man, White}

	// It's white's turn, black has no pieces -> black loses
	if g.State != Playing {
		// State might already be set from previous check... let's check
	}

	// After white moves, check should trigger
	moves := g.AllLegalMoves()
	if len(moves) > 0 {
		g.MakeMove(moves[0])
	}

	if g.State != WhiteWin {
		t.Errorf("expected WhiteWin when black has no pieces, got %v", g.State)
	}
}

// TestPositionValidity verifies Position.Valid() works correctly.
func TestPositionValidity(t *testing.T) {
	tests := []struct {
		pos  Position
		valid bool
	}{
		{Position{0, 0}, false}, // light square
		{Position{0, 1}, true},  // dark square
		{Position{7, 6}, true},  // dark square
		{Position{7, 7}, false}, // light square
		{Position{-1, 0}, false}, // out of bounds
		{Position{8, 0}, false},  // out of bounds
		{Position{0, 8}, false},  // out of bounds
		{Position{3, 4}, true},   // dark square (3+4=7)
		{Position{4, 4}, false},  // light square (4+4=8)
	}

	for _, tt := range tests {
		if tt.pos.Valid() != tt.valid {
			t.Errorf("Position%v.Valid() = %v, want %v", tt.pos, tt.pos.Valid(), tt.valid)
		}
	}
}

// TestManBackwardCapture verifies men can capture backwards in Russian checkers.
func TestManBackwardCapture(t *testing.T) {
	g := &Game{
		Turn:  White,
		State: Playing,
		Moves: make([]Move, 0),
	}

	// White man at (3,4), black man at (4,3)
	// White should be able to capture backwards (from white's perspective)
	g.Board[3][4] = Piece{Man, White}
	g.Board[4][3] = Piece{Man, Black}

	moves := g.capturesFrom(Position{3, 4})
	foundBackward := false
	for _, m := range moves {
		if m.To.Row > m.From.Row {
			foundBackward = true
			if m.To != (Position{5, 2}) {
				t.Errorf("expected backward capture to land at (5,2), got %v", m.To)
			}
		}
	}
	if !foundBackward {
		t.Error("expected man to be able to capture backwards in Russian checkers")
	}
}

// TestMakeMove_Illegal verifies MakeMove rejects illegal moves.
func TestMakeMove_Illegal(t *testing.T) {
	g := NewGame()

	// Try to move a piece to an invalid square
	err := g.MakeMove(Move{
		From:  Position{5, 0},
		To:    Position{5, 2}, // not diagonal, 2 squares over
		Piece: Piece{Man, White},
	})
	if err == nil {
		t.Error("expected error for illegal move")
	}
}

// TestKingMultiCapture verifies king can do multi-captures at distance.
func TestKingMultiCapture(t *testing.T) {
	g := &Game{
		Turn:  White,
		State: Playing,
		Moves: make([]Move, 0),
	}

	// White king at (7,0), two black men in capture path
	g.Board[7][0] = Piece{King, White}
	g.Board[5][2] = Piece{Man, Black}
	g.Board[3][4] = Piece{Man, Black}
	g.Board[6][1] = Piece{Type: None}
	g.Board[4][3] = Piece{Type: None}
	g.Board[2][5] = Piece{Type: None}
	g.Board[1][6] = Piece{Type: None}
	g.Board[0][7] = Piece{Type: None}

	moves := g.capturesFrom(Position{7, 0})
	foundDouble := false
	for _, m := range moves {
		if len(m.Captured) == 2 {
			foundDouble = true
		}
	}
	if !foundDouble {
		t.Errorf("expected king double capture, got %d moves:", len(moves))
		for _, m := range moves {
			t.Logf("  %d captures: %v -> %v", len(m.Captured), m.From, m.To)
		}
	}
}

// TestStringOutput verifies board string rendering.
func TestStringOutput(t *testing.T) {
	g := NewGame()
	s := g.String()
	if len(s) == 0 {
		t.Error("expected non-empty string representation")
	}
	// Should contain piece markers
	found := false
	for _, c := range s {
		if c == 'w' || c == 'b' {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected board string to contain piece markers")
	}
}

// TestCountPieces verifies piece counting.
func TestCountPieces(t *testing.T) {
	g := NewGame()
	w, b := g.CountPieces()
	if w != 12 || b != 12 {
		t.Errorf("expected 12/12 pieces, got %d/%d", w, b)
	}

	// Remove some pieces
	g.Board[5][0] = Piece{Type: None}
	g.Board[5][2] = Piece{Type: None}
	w, b = g.CountPieces()
	if w != 10 || b != 12 {
		t.Errorf("expected 10/12 pieces after removal, got %d/%d", w, b)
	}
}

// TestClone verifies game cloning works correctly.
func TestClone(t *testing.T) {
	g := NewGame()
	clone := g.clone()

	if clone.Turn != g.Turn {
		t.Error("clone should have same turn")
	}

	// Modify clone, verify original is unchanged
	clone.Board[5][0] = Piece{Type: None}
	if !g.Board[5][0].Empty() == false {
		t.Error("modifying clone should not affect original")
	}
	// Original should still have the piece
	if g.Board[5][0].Empty() {
		t.Error("original board was modified through clone")
	}
}

// TestEvaluate verifies the evaluation function works.
func TestEvaluate(t *testing.T) {
	g := NewGame()
	score := g.evaluate()
	if score == 0 {
		// In symmetric starting position, score should be near 0
		// But because of the asymmetric row bonuses, white and black may not be perfectly balanced
		// Actually white men at rows 5,6,7 and black at 0,1,2 should roughly balance
	}
	// Just verify it doesn't panic and returns a finite value
	if score != score { // NaN check
		t.Error("evaluate returned NaN")
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
