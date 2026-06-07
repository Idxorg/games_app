package mocks

import (
	"context"
	"fmt"
	"sync"
	"time"

	"game-platform/internal/model"
)

// MockUserRepo — in-memory mock
type MockUserRepo struct {
	mu    sync.RWMutex
	users map[string]*model.User
}

func NewMockUserRepo() *MockUserRepo {
	return &MockUserRepo{users: make(map[string]*model.User)}
}

func (m *MockUserRepo) GetBySID(_ context.Context, sid string) (*model.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	u, ok := m.users[sid]
	if !ok {
		return nil, nil
	}
	return u, nil
}

func (m *MockUserRepo) Create(_ context.Context, sid, email, name, gender, department, position, photoURL string) (*model.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	u := &model.User{SID: sid, Email: email, Name: name, Gender: gender, Department: department, Position: position, PhotoURL: photoURL, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	m.users[sid] = u
	return u, nil
}

func (m *MockUserRepo) Update(_ context.Context, user *model.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.users[user.SID]; !ok {
		return fmt.Errorf("user not found")
	}
	user.UpdatedAt = time.Now()
	m.users[user.SID] = user
	return nil
}

func (m *MockUserRepo) GetUserGroups(_ context.Context, _ string) ([]string, error) {
	return []string{"games", "tournaments"}, nil
}

// MockTournamentRepo — in-memory mock
type MockTournamentRepo struct {
	mu          sync.RWMutex
	tournaments map[string]*model.Tournament
	players     map[string]map[string]*model.TournamentPlayer
	nextID      int
}

func NewMockTournamentRepo() *MockTournamentRepo {
	return &MockTournamentRepo{tournaments: make(map[string]*model.Tournament), players: make(map[string]map[string]*model.TournamentPlayer)}
}

func (m *MockTournamentRepo) Create(_ context.Context, t *model.Tournament) (*model.Tournament, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if t.ID == "" {
		t.ID = fmt.Sprintf("t_%03d", m.nextID)
		m.nextID++
	}
	t.CreatedAt = time.Now()
	m.tournaments[t.ID] = t
	return t, nil
}

func (m *MockTournamentRepo) GetByID(_ context.Context, id string) (*model.Tournament, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	t, ok := m.tournaments[id]
	if !ok {
		return nil, nil
	}
	return t, nil
}

func (m *MockTournamentRepo) List(_ context.Context, f model.TournamentFilters) ([]model.Tournament, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []model.Tournament
	for _, t := range m.tournaments {
		if f.GameType != "" && t.GameType != f.GameType {
			continue
		}
		if f.Status != "" && t.Status != f.Status {
			continue
		}
		result = append(result, *t)
	}
	return result, nil
}

func (m *MockTournamentRepo) Update(_ context.Context, t *model.Tournament) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.tournaments[t.ID]; !ok {
		return fmt.Errorf("not found")
	}
	m.tournaments[t.ID] = t
	return nil
}

func (m *MockTournamentRepo) Delete(_ context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.tournaments, id)
	delete(m.players, id)
	return nil
}

func (m *MockTournamentRepo) AddPlayer(_ context.Context, tournamentID, sid string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.players[tournamentID] == nil {
		m.players[tournamentID] = make(map[string]*model.TournamentPlayer)
	}
	m.players[tournamentID][sid] = &model.TournamentPlayer{SID: sid}
	return nil
}

func (m *MockTournamentRepo) RemovePlayer(_ context.Context, tournamentID, sid string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.players[tournamentID] != nil {
		delete(m.players[tournamentID], sid)
	}
	return nil
}

func (m *MockTournamentRepo) GetPlayers(_ context.Context, tournamentID string) ([]model.TournamentPlayer, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	players := m.players[tournamentID]
	if players == nil {
		return nil, nil
	}
	var result []model.TournamentPlayer
	for _, p := range players {
		result = append(result, *p)
	}
	return result, nil
}

func (m *MockTournamentRepo) CountPlayerTournaments(_ context.Context, sid string) (int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	count := 0
	for _, players := range m.players {
		if _, ok := players[sid]; ok {
			count++
		}
	}
	return count, nil
}

// MockMatchRepo — in-memory mock
type MockMatchRepo struct {
	mu     sync.RWMutex
	matches map[string]*model.Match
	nextID int
}

func NewMockMatchRepo() *MockMatchRepo {
	return &MockMatchRepo{matches: make(map[string]*model.Match)}
}

func (m *MockMatchRepo) Create(_ context.Context, match *model.Match) (*model.Match, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if match.ID == "" {
		match.ID = fmt.Sprintf("m_%03d", m.nextID)
		m.nextID++
	}
	match.CreatedAt = time.Now()
	m.matches[match.ID] = match
	return match, nil
}

func (m *MockMatchRepo) GetByID(_ context.Context, id string) (*model.Match, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	match, ok := m.matches[id]
	if !ok {
		return nil, nil
	}
	return match, nil
}

func (m *MockMatchRepo) Update(_ context.Context, match *model.Match) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.matches[match.ID]; !ok {
		return fmt.Errorf("not found")
	}
	m.matches[match.ID] = match
	return nil
}

func (m *MockMatchRepo) ListByTournament(_ context.Context, tournamentID string) ([]model.Match, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []model.Match
	for _, match := range m.matches {
		if match.TournamentID == tournamentID {
			result = append(result, *match)
		}
	}
	return result, nil
}

func (m *MockMatchRepo) ListByPlayer(_ context.Context, sid string) ([]model.Match, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []model.Match
	for _, match := range m.matches {
		if match.Player1SID == sid || match.Player2SID == sid {
			result = append(result, *match)
		}
	}
	return result, nil
}

func (m *MockMatchRepo) Complete(_ context.Context, id, winnerID, score string, moves []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	match, ok := m.matches[id]
	if !ok {
		return fmt.Errorf("not found")
	}
	match.WinnerSID = winnerID
	match.Score = score
	match.Moves = moves
	match.Status = "completed"
	now := time.Now()
	match.CompletedAt = &now
	return nil
}

func (m *MockMatchRepo) GetPlayerStats(_ context.Context, sid string) (*model.PlayerStats, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var stats model.PlayerStats
	for _, match := range m.matches {
		if match.Player1SID != sid && match.Player2SID != sid {
			continue
		}
		if match.Status != "completed" {
			continue
		}
		stats.GamesPlayed++
		if match.WinnerSID == "" || match.WinnerSID == "0" {
			stats.Draws++
		} else if match.WinnerSID == sid {
			stats.Wins++
		} else {
			stats.Losses++
		}
	}
	return &stats, nil
}

// MockRatingRepo — in-memory mock
type MockRatingRepo struct {
	mu      sync.RWMutex
	ratings map[string]*model.PlayerRating
}

func NewMockRatingRepo() *MockRatingRepo {
	return &MockRatingRepo{ratings: make(map[string]*model.PlayerRating)}
}

func ratingKey(sid, gameType string) string {
	return sid + ":" + gameType
}

func (m *MockRatingRepo) Get(_ context.Context, sid, gameType string) (*model.PlayerRating, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	r, ok := m.ratings[ratingKey(sid, gameType)]
	if !ok {
		return nil, nil
	}
	return r, nil
}

func (m *MockRatingRepo) Upsert(_ context.Context, rating *model.PlayerRating) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ratings[ratingKey(rating.SID, rating.GameType)] = rating
	return nil
}

func (m *MockRatingRepo) GetLeaderboard(_ context.Context, gameType string, limit int) ([]model.PlayerRating, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []model.PlayerRating
	for _, r := range m.ratings {
		if r.GameType == gameType {
			result = append(result, *r)
		}
	}
	return result, nil
}

func (m *MockRatingRepo) GetByDepartment(_ context.Context, gameType, department string) ([]model.PlayerRating, error) {
	return nil, nil
}

// MockInviteRepo — in-memory mock
type MockInviteRepo struct {
	mu      sync.RWMutex
	invites map[string]*model.GameInvite
	nextID  int
}

func NewMockInviteRepo() *MockInviteRepo {
	return &MockInviteRepo{invites: make(map[string]*model.GameInvite)}
}

func (m *MockInviteRepo) Create(_ context.Context, invite *model.GameInvite) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.nextID++
	invite.ID = fmt.Sprintf("inv_%03d", m.nextID)
	if invite.Status == "" {
		invite.Status = "pending"
	}
	if invite.CreatedAt.IsZero() {
		invite.CreatedAt = time.Now()
	}
	m.invites[invite.ID] = invite
	return nil
}

func (m *MockInviteRepo) GetByID(_ context.Context, id string) (*model.GameInvite, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	invite, ok := m.invites[id]
	if !ok {
		return nil, nil
	}
	return invite, nil
}

func (m *MockInviteRepo) GetPendingByRecipient(_ context.Context, sid string) ([]model.GameInvite, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []model.GameInvite
	for _, inv := range m.invites {
		if inv.RecipientSID == sid && inv.Status == "pending" && inv.ExpiresAt.After(time.Now()) {
			result = append(result, *inv)
		}
	}
	return result, nil
}

func (m *MockInviteRepo) Accept(_ context.Context, id string, matchID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	invite, ok := m.invites[id]
	if !ok {
		return fmt.Errorf("not found")
	}
	if invite.Status != "pending" {
		return fmt.Errorf("not pending")
	}
	invite.Status = "accepted"
	invite.MatchID = &matchID
	return nil
}

func (m *MockInviteRepo) Decline(_ context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	invite, ok := m.invites[id]
	if !ok {
		return fmt.Errorf("not found")
	}
	if invite.Status != "pending" {
		return fmt.Errorf("not pending")
	}
	invite.Status = "declined"
	return nil
}

func (m *MockInviteRepo) ExpireOld(_ context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, inv := range m.invites {
		if inv.Status == "pending" && inv.ExpiresAt.Before(time.Now()) {
			inv.Status = "expired"
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// MockRatingUpdater (implements websocket.RatingUpdater)
// ---------------------------------------------------------------------------

// MockRatingUpdater records calls to UpdateMatchRatings for test assertions.
type MockRatingUpdater struct {
	mu    sync.Mutex
	Calls []MockRatingCall
	Err   error // optional error to return
}

// MockRatingCall captures a single call to UpdateMatchRatings.
type MockRatingCall struct {
	MatchID    string
	GameType   string
	Player1SID string
	Player2SID string
	WinnerSID  string
	Score      string
}

func (m *MockRatingUpdater) UpdateMatchRatings(_ context.Context, match *model.Match) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Calls = append(m.Calls, MockRatingCall{
		MatchID:    match.ID,
		GameType:   match.GameType,
		Player1SID: match.Player1SID,
		Player2SID: match.Player2SID,
		WinnerSID:  match.WinnerSID,
		Score:      match.Score,
	})
	return m.Err
}

func (m *MockRatingUpdater) CallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.Calls)
}

// ---------------------------------------------------------------------------
// MockNotificationPublisher (implements handler.NotificationPublisher)
// ---------------------------------------------------------------------------

// MockNotificationPublisher records published events for test assertions.
type MockNotificationPublisher struct {
	mu     sync.Mutex
	Events []map[string]interface{}
	Err    error // optional error to return
}

func (m *MockNotificationPublisher) PublishEvent(event map[string]interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Events = append(m.Events, event)
	return m.Err
}

func (m *MockNotificationPublisher) EventCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.Events)
}
