package repository

import (
	"context"

	"game-platform/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

// MatchRepository implements the MatchRepo interface.
type MatchRepository struct {
	db *pgxpool.Pool
}

// NewMatchRepository creates a new match repository.
func NewMatchRepository(db *pgxpool.Pool) *MatchRepository {
	return &MatchRepository{db: db}
}

// Create inserts a new match and returns it with generated fields.
func (r *MatchRepository) Create(ctx context.Context, match *model.Match) (*model.Match, error) {
	var m model.Match
	err := r.db.QueryRow(ctx, `
		INSERT INTO matches (id, tournament_id, game_type, player1_sid, player2_sid, winner_sid, score, moves, pgn_url, game_id, livekit_room_id, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, tournament_id, game_type, player1_sid, player2_sid, winner_sid, score, moves, pgn_url, game_id, livekit_room_id, status, started_at, completed_at, created_at
	`,
		match.ID, match.TournamentID, match.GameType, match.Player1SID,
		match.Player2SID, match.WinnerSID, match.Score, match.Moves,
		match.PGNURL, match.GameID, match.LiveKitRoomID, match.Status,
	).Scan(
		&m.ID, &m.TournamentID, &m.GameType, &m.Player1SID, &m.Player2SID,
		&m.WinnerSID, &m.Score, &m.Moves, &m.PGNURL, &m.GameID,
		&m.LiveKitRoomID, &m.Status, &m.StartedAt, &m.CompletedAt, &m.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// GetByID retrieves a match by its ID.
func (r *MatchRepository) GetByID(ctx context.Context, id string) (*model.Match, error) {
	var m model.Match
	err := r.db.QueryRow(ctx, `
		SELECT id, tournament_id, game_type, player1_sid, player2_sid, winner_sid, score, moves, pgn_url, game_id, livekit_room_id, status, started_at, completed_at, created_at
		FROM matches WHERE id = $1
	`, id).Scan(
		&m.ID, &m.TournamentID, &m.GameType, &m.Player1SID, &m.Player2SID,
		&m.WinnerSID, &m.Score, &m.Moves, &m.PGNURL, &m.GameID,
		&m.LiveKitRoomID, &m.Status, &m.StartedAt, &m.CompletedAt, &m.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// Update modifies an existing match.
func (r *MatchRepository) Update(ctx context.Context, match *model.Match) error {
	_, err := r.db.Exec(ctx, `
		UPDATE matches SET tournament_id = $2, game_type = $3, player1_sid = $4, player2_sid = $5,
			winner_sid = $6, score = $7, moves = $8, pgn_url = $9, game_id = $10,
			livekit_room_id = $11, status = $12, started_at = $13, completed_at = $14
		WHERE id = $1
	`,
		match.ID, match.TournamentID, match.GameType, match.Player1SID,
		match.Player2SID, match.WinnerSID, match.Score, match.Moves,
		match.PGNURL, match.GameID, match.LiveKitRoomID, match.Status,
		match.StartedAt, match.CompletedAt,
	)
	return err
}

// ListByTournament returns all matches in a tournament.
func (r *MatchRepository) ListByTournament(ctx context.Context, tournamentID string) ([]model.Match, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, tournament_id, game_type, player1_sid, player2_sid, winner_sid, score, moves, pgn_url, game_id, livekit_room_id, status, started_at, completed_at, created_at
		FROM matches WHERE tournament_id = $1
		ORDER BY created_at ASC
	`, tournamentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []model.Match
	for rows.Next() {
		var m model.Match
		if err := rows.Scan(
			&m.ID, &m.TournamentID, &m.GameType, &m.Player1SID, &m.Player2SID,
			&m.WinnerSID, &m.Score, &m.Moves, &m.PGNURL, &m.GameID,
			&m.LiveKitRoomID, &m.Status, &m.StartedAt, &m.CompletedAt, &m.CreatedAt,
		); err != nil {
			return nil, err
		}
		matches = append(matches, m)
	}
	return matches, nil
}

// ListByPlayer returns all matches involving a given player.
func (r *MatchRepository) ListByPlayer(ctx context.Context, sid string) ([]model.Match, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, tournament_id, game_type, player1_sid, player2_sid, winner_sid, score, moves, pgn_url, game_id, livekit_room_id, status, started_at, completed_at, created_at
		FROM matches WHERE player1_sid = $1 OR player2_sid = $1
		ORDER BY created_at DESC
	`, sid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []model.Match
	for rows.Next() {
		var m model.Match
		if err := rows.Scan(
			&m.ID, &m.TournamentID, &m.GameType, &m.Player1SID, &m.Player2SID,
			&m.WinnerSID, &m.Score, &m.Moves, &m.PGNURL, &m.GameID,
			&m.LiveKitRoomID, &m.Status, &m.StartedAt, &m.CompletedAt, &m.CreatedAt,
		); err != nil {
			return nil, err
		}
		matches = append(matches, m)
	}
	return matches, nil
}

// Complete marks a match as completed with the given winner, score, and moves.
func (r *MatchRepository) Complete(ctx context.Context, id, winnerID, score string, movesJSON []byte) error {
	_, err := r.db.Exec(ctx, `
		UPDATE matches SET winner_sid = $2, score = $3, moves = $4, status = 'completed', completed_at = NOW()
		WHERE id = $1
	`, id, winnerID, score, movesJSON)
	return err
}
