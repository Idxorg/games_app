package chess

import (
    "fmt"
    "testing"
)

func TestDebugPromotion(t *testing.T) {
    g := NewGame()
    g.Board = Board{}
    for i := 0; i < 4; i++ {
        g.Board[6][i] = Piece{Pawn, White}
    }
    g.Board[0][4] = Piece{King, White}
    g.Board[7][4] = Piece{King, Black}
    g.Turn = White

    fmt.Printf("Before: Board[7][3] = %+v\n", g.Board[7][3])
    
    pt := Knight
    from := Position{3, 6}
    to := Position{3, 7}
    
    legal := g.LegalMoves(from)
    fmt.Printf("Legal moves from d7: %d moves\n", len(legal))
    for _, m := range legal {
        fmt.Printf("  %s -> %s promo=%v\n", m.From.Algebraic(), m.To.Algebraic(), m.Promotion)
    }
    
    err := g.MakeMove(from, to, pt)
    fmt.Printf("MakeMove error: %v\n", err)
    fmt.Printf("After: Board[7][3] = %+v\n", g.Board[7][3])
    fmt.Printf("After FEN: %s\n", g.ToFEN())
}
