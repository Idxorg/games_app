package repository

import (
	"context"
	"fmt"
	"strings"

	"game-platform/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TournamentRepository implements the TournamentRepo interface.
type TournamentRepository struct {
	db *pgxpool.Pool
}

// NewTournamentRepository creates a new tournament repository.
func NewTournamentRepository(db *pgxpool.Pool) *TournamentRepository {
	return &TournamentRepository{db: db}
}

// Create inserts a new tournament and returns it with generated fields.
func (r *TournamentRepository) Create(ctx context.Context, tournament *model.Tournament) (*model.Tournament, error) {
	var t model.Tournament
	err := r.db.QueryRow(ctx, `
		INSERT INTO tournaments (id, name, game_type, status, start_date, end_date, max_players, current_players, prize_pool, description, logo_url, created_by, requires_group)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, name, game_type, status, start_date, end_date, max_players, current_players, prize_pool, description, logo_url, created_by, requires_group, created_at
	`,
		tournament.ID, tournament.Name, tournament.GameType, tournament.Status,
		tournament.StartDate, tournament.EndDate, tournament.MaxPlayers,
		tournament.CurrentPlayers, tournament.PrizePool, tournament.Description,
		tournament.LogoURL, tournament.CreatedBy, tournament.RequiresGroup,
	).Scan(
		&t.ID, &t.Name, &t.GameType, &t.Status, &t.StartDate, &t.EndDate,
		&t.MaxPlayers, &t.CurrentPlayers, &t.PrizePool, &t.Description,
		&t.LogoURL, &t.CreatedBy, &t.RequiresGroup, &t.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// GetByID retrieves a tournament by its ID.
func (r *TournamentRepository) GetByID(ctx context.Context, id string) (*model.Tournament, error) {
	var t model.Tournament
	err := r.db.QueryRow(ctx, `
		SELECT id, name, game_type, status, start_date, end_date, max_players, current_players, prize_pool, description, logo_url, created_by, requires_group, created_at
		FROM tournaments WHERE id = $1
	`, id).Scan(
		&t.ID, &t.Name, &t.GameType, &t.Status, &t.StartDate, &t.EndDate,
		&t.MaxPlayers, &t.CurrentPlayers, &t.PrizePool, &t.Description,
		&t.LogoURL, &t.CreatedBy, &t.RequiresGroup, &t.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// List returns tournaments matching the given filters.
func (r *TournamentRepository) List(ctx context.Context, filters model.TournamentFilters) ([]model.Tournament, error) {
	query := `
		SELECT id, name, game_type, status, start_date, end_date, max_players, current_players, prize_pool, description, logo_url, created_by, requires_group, created_at
		FROM tournaments
	`
	var conditions []string
	var args []interface{}
	argNum := 1

	if filters.GameType != "" {
		conditions = append(conditions, fmt.Sprintf("game_type = $%d", argNum))
		args = append(args, filters.GameType)
		argNum++
	}
	if filters.Status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argNum))
		args = append(args, filters.Status)
		argNum++
	}
	if filters.Search != "" {
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", argNum))
		args = append(args, "%"+filters.Search+"%")
		argNum++
	}
	if filters.CreatedBy != "" {
		conditions = append(conditions, fmt.Sprintf("created_by = $%d", argNum))
		args = append(args, filters.CreatedBy)
		argNum++
	}
	if filters.RequiresGroup != "" {
		conditions = append(conditions, fmt.Sprintf("requires_group = $%d", argNum))
		args = append(args, filters.RequiresGroup)
		argNum++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY created_at DESC"

	if filters.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argNum)
		args = append(args, filters.Limit)
		argNum++
	}
	if filters.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argNum)
		args = append(args, filters.Offset)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tournaments []model.Tournament
	for rows.Next() {
		var t model.Tournament
		if err := rows.Scan(
			&t.ID, &t.Name, &t.GameType, &t.Status, &t.StartDate, &t.EndDate,
			&t.MaxPlayers, &t.CurrentPlayers, &t.PrizePool, &t.Description,
			&t.LogoURL, &t.CreatedBy, &t.RequiresGroup, &t.CreatedAt,
		); err != nil {
			return nil, err
		}
		tournaments = append(tournaments, t)
	}
	return tournaments, nil
}

// Update modifies an existing tournament.
func (r *TournamentRepository) Update(ctx context.Context, tournament *model.Tournament) error {
	_, err := r.db.Exec(ctx, `
		UPDATE tournaments SET name = $2, game_type = $3, status = $4, start_date = $5, end_date = $6,
			max_players = $7, current_players = $8, prize_pool = $9, description = $10,
			logo_url = $11, created_by = $12, requires_group = $13
		WHERE id = $1
	`,
		tournament.ID, tournament.Name, tournament.GameType, tournament.Status,
		tournament.StartDate, tournament.EndDate, tournament.MaxPlayers,
		tournament.CurrentPlayers, tournament.PrizePool, tournament.Description,
		tournament.LogoURL, tournament.CreatedBy, tournament.RequiresGroup,
	)
	return err
}

// Delete removes a tournament by its ID.
func (r *TournamentRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM tournaments WHERE id = $1`, id)
	return err
}

// AddPlayer adds a player to a tournament.
func (r *TournamentRepository) AddPlayer(ctx context.Context, tournamentID, sid string) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO tournament_players (tournament_id, sid) VALUES ($1, $2)
		ON CONFLICT (tournament_id, sid) DO NOTHING
	`, tournamentID, sid)
	return err
}

// RemovePlayer removes a player from a tournament.
func (r *TournamentRepository) RemovePlayer(ctx context.Context, tournamentID, sid string) error {
	_, err := r.db.Exec(ctx, `
		DELETE FROM tournament_players WHERE tournament_id = $1 AND sid = $2
	`, tournamentID, sid)
	return err
}

// GetPlayers returns all players in a tournament.
func (r *TournamentRepository) GetPlayers(ctx context.Context, tournamentID string) ([]model.TournamentPlayer, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, tournament_id, sid, rank, points, wins, draws, losses, joined_at
		FROM tournament_players WHERE tournament_id = $1
		ORDER BY points DESC, rank ASC
	`, tournamentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []model.TournamentPlayer
	for rows.Next() {
		var p model.TournamentPlayer
		if err := rows.Scan(
			&p.ID, &p.TournamentID, &p.SID, &p.Rank, &p.Points,
			&p.Wins, &p.Draws, &p.Losses, &p.JoinedAt,
		); err != nil {
			return nil, err
		}
		players = append(players, p)
	}
	return players, nil
}
