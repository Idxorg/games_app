package repository

import (
	"context"

	"game-platform/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RatingRepository implements the RatingRepo interface.
type RatingRepository struct {
	db *pgxpool.Pool
}

// NewRatingRepository creates a new rating repository.
func NewRatingRepository(db *pgxpool.Pool) *RatingRepository {
	return &RatingRepository{db: db}
}

// Get retrieves a player's rating for a specific game type.
func (r *RatingRepository) Get(ctx context.Context, sid, gameType string) (*model.PlayerRating, error) {
	var rating model.PlayerRating
	err := r.db.QueryRow(ctx, `
		SELECT id, sid, game_type, elo, games_played, wins, draws, losses
		FROM player_ratings WHERE sid = $1 AND game_type = $2
	`, sid, gameType).Scan(
		&rating.ID, &rating.SID, &rating.GameType, &rating.ELO,
		&rating.GamesPlayed, &rating.Wins, &rating.Draws, &rating.Losses,
	)
	if err != nil {
		return nil, err
	}
	return &rating, nil
}

// Upsert inserts or updates a player's rating.
func (r *RatingRepository) Upsert(ctx context.Context, rating *model.PlayerRating) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO player_ratings (sid, game_type, elo, games_played, wins, draws, losses)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (sid, game_type) DO UPDATE SET
			elo = EXCLUDED.elo,
			games_played = EXCLUDED.games_played,
			wins = EXCLUDED.wins,
			draws = EXCLUDED.draws,
			losses = EXCLUDED.losses
	`, rating.SID, rating.GameType, rating.ELO, rating.GamesPlayed, rating.Wins, rating.Draws, rating.Losses)
	return err
}

// GetLeaderboard returns the top-rated players for a game type.
func (r *RatingRepository) GetLeaderboard(ctx context.Context, gameType string, limit int) ([]model.PlayerRating, error) {
	rows, err := r.db.Query(ctx, `
		SELECT pr.id, pr.sid, pr.game_type, pr.elo, pr.games_played, pr.wins, pr.draws, pr.losses
		FROM player_ratings pr
		JOIN users u ON pr.sid = u.sid
		WHERE pr.game_type = $1
		ORDER BY pr.elo DESC
		LIMIT $2
	`, gameType, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ratings []model.PlayerRating
	for rows.Next() {
		var rating model.PlayerRating
		if err := rows.Scan(
			&rating.ID, &rating.SID, &rating.GameType, &rating.ELO,
			&rating.GamesPlayed, &rating.Wins, &rating.Draws, &rating.Losses,
		); err != nil {
			return nil, err
		}
		ratings = append(ratings, rating)
	}
	return ratings, nil
}

// GetByDepartment returns ratings for players in a specific department for a game type.
func (r *RatingRepository) GetByDepartment(ctx context.Context, gameType, department string) ([]model.PlayerRating, error) {
	rows, err := r.db.Query(ctx, `
		SELECT pr.id, pr.sid, pr.game_type, pr.elo, pr.games_played, pr.wins, pr.draws, pr.losses
		FROM player_ratings pr
		JOIN users u ON pr.sid = u.sid
		WHERE pr.game_type = $1 AND u.department = $2
		ORDER BY pr.elo DESC
	`, gameType, department)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ratings []model.PlayerRating
	for rows.Next() {
		var rating model.PlayerRating
		if err := rows.Scan(
			&rating.ID, &rating.SID, &rating.GameType, &rating.ELO,
			&rating.GamesPlayed, &rating.Wins, &rating.Draws, &rating.Losses,
		); err != nil {
			return nil, err
		}
		ratings = append(ratings, rating)
	}
	return ratings, nil
}
