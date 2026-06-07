package websocket

// WSMove represents a move in the WebSocket protocol format.
type WSMove struct {
	From      string `json:"from,omitempty"`       // e.g. "e2"
	To        string `json:"to,omitempty"`         // e.g. "e4"
	Promotion string `json:"promotion,omitempty"`  // e.g. "q"
	Die       int    `json:"die,omitempty"`         // backgammon die value
	FromIdx   int    `json:"from_idx,omitempty"`     // backgammon: point index or -1 for bar
	ToIdx     int    `json:"to_idx,omitempty"`       // backgammon: point index or 24 for off
	// Checkers: captured positions
	Captured []WSPosition `json:"captured,omitempty"`
}

// WSPosition represents a board position in WS format.
type WSPosition struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

// GameEngineAdapter is the interface all game adapters must implement.
type GameEngineAdapter interface {
	// GetBoard returns the board state as a serializable string (FEN for chess, custom for others).
	GetBoard() string

	// GetTurn returns the active color: "white" or "black".
	GetTurn() string

	// GetLegalMoves returns legal moves in WS format for the current turn.
	GetLegalMoves() []WSMove

	// ApplyMove applies a move from WS format and returns an error if illegal.
	ApplyMove(move WSMove) error

	// IsGameOver returns true if the game has ended.
	IsGameOver() bool

	// GetGameOverReason returns the reason string ("checkmate", "stalemate", "timeout", "resign", "draw").
	GetGameOverReason() string

	// GetWinner returns the winning color ("white", "black", or "" for draw).
	GetWinner() string

	// Resign sets the game as resigned by the given color.
	Resign(color string)

	// RollDice rolls dice (backgammon only, no-op for others).
	RollDice() []int
}
