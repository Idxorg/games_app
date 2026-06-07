package repository

import (
	"context"

	"game-platform/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// InviteRepository implements the InviteRepo interface.
type InviteRepository struct {
	db *pgxpool.Pool
}

// NewInviteRepository creates a new invite repository.
func NewInviteRepository(db *pgxpool.Pool) *InviteRepository {
	return &InviteRepository{db: db}
}

// Create inserts a new game invite.
func (r *InviteRepository) Create(ctx context.Context, invite *model.GameInvite) error {
	return r.db.QueryRow(ctx, `
		INSERT INTO game_invites (game_type, inviter_sid, recipient_sid, status, match_id, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, expires_at
	`,
		invite.GameType, invite.InviterSID, invite.RecipientSID,
		invite.Status, invite.MatchID, invite.ExpiresAt,
	).Scan(&invite.ID, &invite.CreatedAt, &invite.ExpiresAt)
}

// GetByID retrieves an invite by its ID.
func (r *InviteRepository) GetByID(ctx context.Context, id string) (*model.GameInvite, error) {
	var invite model.GameInvite
	err := r.db.QueryRow(ctx, `
		SELECT id, game_type, inviter_sid, recipient_sid, status, match_id, created_at, expires_at
		FROM game_invites WHERE id = $1
	`, id).Scan(
		&invite.ID, &invite.GameType, &invite.InviterSID, &invite.RecipientSID,
		&invite.Status, &invite.MatchID, &invite.CreatedAt, &invite.ExpiresAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &invite, nil
}

// GetPendingByRecipient returns all pending (non-expired) invites for a recipient.
func (r *InviteRepository) GetPendingByRecipient(ctx context.Context, sid string) ([]model.GameInvite, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, game_type, inviter_sid, recipient_sid, status, match_id, created_at, expires_at
		FROM game_invites
		WHERE recipient_sid = $1 AND status = 'pending' AND expires_at > NOW()
		ORDER BY created_at DESC
	`, sid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invites []model.GameInvite
	for rows.Next() {
		var invite model.GameInvite
		if err := rows.Scan(
			&invite.ID, &invite.GameType, &invite.InviterSID, &invite.RecipientSID,
			&invite.Status, &invite.MatchID, &invite.CreatedAt, &invite.ExpiresAt,
		); err != nil {
			return nil, err
		}
		invites = append(invites, invite)
	}
	return invites, nil
}

// Accept marks an invite as accepted and sets the match ID.
func (r *InviteRepository) Accept(ctx context.Context, id string, matchID string) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE game_invites
		SET status = 'accepted', match_id = $2, expires_at = NOW()
		WHERE id = $1 AND status = 'pending'
	`, id, matchID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

// Decline marks an invite as declined.
func (r *InviteRepository) Decline(ctx context.Context, id string) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE game_invites SET status = 'declined'
		WHERE id = $1 AND status = 'pending'
	`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

// ExpireOld marks all expired pending invites as expired.
func (r *InviteRepository) ExpireOld(ctx context.Context) error {
	_, err := r.db.Exec(ctx, `
		UPDATE game_invites SET status = 'expired'
		WHERE status = 'pending' AND expires_at < NOW()
	`)
	return err
}
