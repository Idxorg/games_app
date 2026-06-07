package chess

import (
    "fmt"
    "testing"
)

func TestDebugPromotion2(t *testing.T) {
    g := NewGame()
    g.Board = Board{}
    for i := 0; i < 4; i++ {
        g.Board[6][i] = Piece{Pawn, White}
    }
    g.Board[0][4] = Piece{King, White}
    g.Board[7][4] = Piece{King, Black}
    g.Turn = White

    pt := Knight
    from := Position{3, 6}
    to := Position{3, 7}
    
    // Get the validated move
    move, err := g.ValidateMove(from, to, pt)
    fmt.Printf("Validated move: From=%s To=%s Promo=%v Piece=%+v\n", 
        move.From.Algebraic(), move.To.Algebraic(), move.Promotion, move.Piece)
    fmt.Printf("Error: %v\n", err)
    
    if err != nil {
        return
    }
    
    // Now call executeMove directly
    fmt.Printf("Before executeMove: Board[7][3]=%+v Board[6][3]=%+v\n", g.Board[7][3], g.Board[6][3])
    err = g.executeMove(move)
    fmt.Printf("executeMove error: %v\n", err)
    fmt.Printf("After executeMove: Board[7][3]=%+v Board[6][3]=%+v\n", g.Board[7][3], g.Board[6][3])
    fmt.Printf("FEN: %s\n", g.ToFEN())
}
