package model

import (
	"context"
	"time"
)

// TournamentFilters represents filtering options for tournament listing.
type TournamentFilters struct {
	GameType      string
	Status        string
	Search        string
	CreatedBy     string
	RequiresGroup string
	Limit         int
	Offset        int
}

// UserRepo defines the user repository interface.
type UserRepo interface {
	GetBySID(ctx context.Context, sid string) (*User, error)
	Create(ctx context.Context, sid, email, name, gender, department, position, photoURL string) (*User, error)
	Update(ctx context.Context, user *User) error
	GetUserGroups(ctx context.Context, sid string) ([]string, error)
}

// TournamentRepo defines the tournament repository interface.
type TournamentRepo interface {
	List(ctx context.Context, filters TournamentFilters) ([]Tournament, error)
	GetByID(ctx context.Context, id string) (*Tournament, error)
	Create(ctx context.Context, tournament *Tournament) (*Tournament, error)
	Update(ctx context.Context, tournament *Tournament) error
	Delete(ctx context.Context, id string) error
	AddPlayer(ctx context.Context, tournamentID, sid string) error
	RemovePlayer(ctx context.Context, tournamentID, sid string) error
	GetPlayers(ctx context.Context, tournamentID string) ([]TournamentPlayer, error)
	CountPlayerTournaments(ctx context.Context, sid string) (int, error)
}

// MatchRepo defines the match repository interface.
type MatchRepo interface {
	Create(ctx context.Context, match *Match) (*Match, error)
	GetByID(ctx context.Context, id string) (*Match, error)
	Update(ctx context.Context, match *Match) error
	ListByTournament(ctx context.Context, tournamentID string) ([]Match, error)
	ListByPlayer(ctx context.Context, sid string) ([]Match, error)
	Complete(ctx context.Context, id, winnerID, score string, movesJSON []byte) error
	GetPlayerStats(ctx context.Context, sid string) (*PlayerStats, error)
}

// RatingRepo defines the rating repository interface.
type RatingRepo interface {
	Get(ctx context.Context, sid, gameType string) (*PlayerRating, error)
	Upsert(ctx context.Context, rating *PlayerRating) error
	GetLeaderboard(ctx context.Context, gameType string, limit int) ([]PlayerRating, error)
	GetByDepartment(ctx context.Context, gameType, department string) ([]PlayerRating, error)
}

// InviteRepo defines the invite repository interface.
type InviteRepo interface {
	Create(ctx context.Context, invite *GameInvite) error
	GetByID(ctx context.Context, id string) (*GameInvite, error)
	GetPendingByRecipient(ctx context.Context, sid string) ([]GameInvite, error)
	Accept(ctx context.Context, id string, matchID string) error
	Decline(ctx context.Context, id string) error
	ExpireOld(ctx context.Context) error
}

// S3Service defines the S3 storage service interface.
type S3Service interface {
	UploadAvatar(ctx context.Context, sid string, data []byte) (string, error)
	UploadPGN(ctx context.Context, gameType, matchID, pgn string) (string, error)
	GetAvatarURL(ctx context.Context, sid string) (string, error)
}

// PortalService defines the corporate portal service interface.
type PortalService interface {
	GetUser(ctx context.Context, sid string) (*User, error)
	HasAccess(ctx context.Context, sid, group string) (bool, error)
	GetUserGroups(ctx context.Context, sid string) ([]string, error)
}

// CacheService defines the cache service interface.
type CacheService interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string, ttl time.Duration) error
	Del(ctx context.Context, key string) error
	Increment(ctx context.Context, key string) (int64, error)
}
