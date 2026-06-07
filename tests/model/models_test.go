package model_test

import (
	"testing"
	"time"

	"game-platform/internal/model"
)

func TestUser(t *testing.T) {
	now := time.Now()
	user := &model.User{
		SID:        "emp_12345",
		Email:      "ivanov@yakbson.digital",
		Name:       "Иванов Иван",
		Gender:     "male",
		Department: "IT",
		Position:   "Разработчик",
		PhotoURL:   "https://s3.yakbson.digital/avatars/emp_12345.jpg",
		LastSync:   now,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if user.SID != "emp_12345" {
		t.Errorf("Expected SID emp_12345, got %s", user.SID)
	}

	if user.Email != "ivanov@yakbson.digital" {
		t.Errorf("Expected email ivanov@yakbson.digital, got %s", user.Email)
	}

	if user.Department != "IT" {
		t.Errorf("Expected department IT, got %s", user.Department)
	}
}

func TestPlayerRating(t *testing.T) {
	rating := &model.PlayerRating{
		ID:           1,
		SID:          "emp_12345",
		GameType:     "chess",
		ELO:          1450,
		GamesPlayed:  100,
		Wins:         60,
		Draws:        20,
		Losses:       20,
	}

	if rating.ELO != 1450 {
		t.Errorf("Expected ELO 1450, got %d", rating.ELO)
	}

	if rating.Wins != 60 {
		t.Errorf("Expected wins 60, got %d", rating.Wins)
	}

	// Проверка расчета win rate
	expectedWinRate := float64(rating.Wins) / float64(rating.GamesPlayed) * 100
	if expectedWinRate != 60.0 {
		t.Errorf("Expected win rate 60.0%%, got %f", expectedWinRate)
	}
}

func TestTournament(t *testing.T) {
	startDate := time.Now()
	endDate := startDate.AddDate(0, 0, 30)

	tournament := &model.Tournament{
		ID:            "t_001",
		Name:          "Весенний чемпионат по шахматам 2026",
		GameType:      "chess",
		Status:        "active",
		StartDate:     startDate,
		EndDate:       endDate,
		MaxPlayers:    128,
		CurrentPlayers: 45,
		PrizePool:     "100,000 руб.",
		Description:   "Ежемесячный чемпионат по шахматам",
		LogoURL:       "https://s3.yakbson.digital/tournaments/t_001/logo.jpg",
		CreatedBy:     "emp_11111",
		RequiresGroup: "tournaments",
		CreatedAt:     time.Now(),
	}

	if tournament.Name != "Весенний чемпионат по шахматам 2026" {
		t.Errorf("Expected name 'Весенний чемпионат по шахматам 2026', got %s", tournament.Name)
	}

	if tournament.MaxPlayers != 128 {
		t.Errorf("Expected max_players 128, got %d", tournament.MaxPlayers)
	}

	if tournament.CurrentPlayers != 45 {
		t.Errorf("Expected current_players 45, got %d", tournament.CurrentPlayers)
	}

	// Проверка что турнир активен
	if tournament.Status != "active" {
		t.Errorf("Expected status 'active', got %s", tournament.Status)
	}
}

func TestMatch(t *testing.T) {
	startedAt := time.Now()
	completedAt := startedAt.Add(1 * time.Hour)

	match := &model.Match{
		ID:            "m_001",
		TournamentID:  "t_001",
		GameType:      "chess",
		Player1SID:    "emp_12345",
		Player2SID:    "emp_67890",
		WinnerSID:     "emp_12345",
		Score:         "1-0",
		Moves:         []byte(`[{"from":"e2","to":"e4"},{"from":"e7","to":"e5"}]`),
		PGNURL:        "https://s3.yakbson.digital/records/chess/m_001.pgn",
		GameID:        "g_chess_123456",
		LiveKitRoomID: "room_abc123",
		Status:        "completed",
		StartedAt:     &startedAt,
		CompletedAt:   &completedAt,
		CreatedAt:     time.Now(),
	}

	if match.ID != "m_001" {
		t.Errorf("Expected ID m_001, got %s", match.ID)
	}

	if match.WinnerSID != "emp_12345" {
		t.Errorf("Expected winner emp_12345, got %s", match.WinnerSID)
	}

	if match.Status != "completed" {
		t.Errorf("Expected status 'completed', got %s", match.Status)
	}

	if len(match.Moves) == 0 {
		t.Error("Expected moves data")
	}
}
