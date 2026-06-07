package model

import "time"

// User представляет пользователя из корп портала
type User struct {
    SID        string    `json:"sid"`
    Email      string    `json:"email"`
    Name       string    `json:"name"`
    Gender     string    `json:"gender"`
    Department string    `json:"department"`
    Position   string    `json:"position"`
    PhotoURL   string    `json:"photo_url"`
    LastSync   time.Time `json:"last_sync"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
}

// UserGroups представляет группу пользователя
type UserGroup struct {
    ID       int       `json:"id"`
    SID      string    `json:"sid"`
    GroupName string   `json:"group_name"`
    GrantedAt time.Time `json:"granted_at"`
}

// PlayerRating представляет рейтинг игрока
type PlayerRating struct {
    ID           int    `json:"id"`
    SID          string `json:"sid"`
    GameType     string `json:"game_type"`
    ELO          int    `json:"elo"`
    GamesPlayed  int    `json:"games_played"`
    Wins         int    `json:"wins"`
    Draws        int    `json:"draws"`
    Losses       int    `json:"losses"`
}

// Tournament представляет турнир
type Tournament struct {
    ID            string    `json:"id"`
    Name          string    `json:"name"`
    GameType      string    `json:"game_type"`
    Status        string    `json:"status"`
    StartDate     time.Time `json:"start_date"`
    EndDate       time.Time `json:"end_date"`
    MaxPlayers    int       `json:"max_players"`
    CurrentPlayers int      `json:"current_players"`
    PrizePool     string    `json:"prize_pool"`
    Description   string    `json:"description"`
    LogoURL       string    `json:"logo_url"`
    CreatedBy     string    `json:"created_by"`
    RequiresGroup string    `json:"requires_group"`
    CreatedAt     time.Time `json:"created_at"`
}

// TournamentPlayer представляет участника турнира
type TournamentPlayer struct {
    ID        int       `json:"id"`
    TournamentID string  `json:"tournament_id"`
    SID       string    `json:"sid"`
    Rank      int       `json:"rank"`
    Points    int       `json:"points"`
    Wins      int       `json:"wins"`
    Draws     int       `json:"draws"`
    Losses    int       `json:"losses"`
    JoinedAt  time.Time `json:"joined_at"`
}

// GameInvite представляет приглашение на игру
type GameInvite struct {
	ID            string    `json:"id"`
	GameType      string    `json:"game_type"`
	InviterSID    string    `json:"inviter_sid"`
	RecipientSID  string    `json:"recipient_sid"`
	Status        string    `json:"status"`
	MatchID       *string   `json:"match_id,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	ExpiresAt     time.Time `json:"expires_at"`
}

// PlayerStats represents aggregated player statistics.
type PlayerStats struct {
	GamesPlayed      int `json:"games_played"`
	Wins            int `json:"wins"`
	Draws           int `json:"draws"`
	Losses          int `json:"losses"`
	TournamentsJoined int `json:"tournaments_joined"`
}

// Match представляет матч
type Match struct {
    ID            string    `json:"id"`
    TournamentID  string    `json:"tournament_id"`
    GameType      string    `json:"game_type"`
    Player1SID    string    `json:"player1_sid"`
    Player2SID    string    `json:"player2_sid"`
    WinnerSID     string    `json:"winner_sid"`
    Score         string    `json:"score"`
    Moves         []byte    `json:"moves"`  // JSONB
    PGNURL        string    `json:"pgn_url"`
    GameID        string    `json:"game_id"`
    LiveKitRoomID string    `json:"livekit_room_id"`
    Status        string    `json:"status"`
    StartedAt     *time.Time `json:"started_at"`
    CompletedAt   *time.Time `json:"completed_at"`
    CreatedAt     time.Time `json:"created_at"`
}
